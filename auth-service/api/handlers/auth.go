package handlers

import (
	"auth-service/api/token"
	"auth-service/models"
	"context"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/golang-jwt/jwt"
	_ "github.com/swaggo/swag"
	"golang.org/x/crypto/bcrypt"
)

// Register godoc
// @Summary Register a new user
// @Description Register a new user with email, username, and password
// @Tags registration
// @Accept json
// @Produce json
// @Param user body models.RegisterReqSwag true "User registration request"
// @Success 201 {object} token.Tokens "JWT tokens"
// @Failure 400 {object} string "Invalid request payload"
// @Failure 500 {object} string "Server error"
// @Router /register [post]
func (h *HTTPHandler) Register(c *gin.Context) {
	var req models.RegisterReq

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Invalid request payload": err.Error()})
		return
	}

	// Check if email is valid
	if !isValidEmail(req.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	// Check if email already exists
	exists, err := h.US.EmailExists(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error", "err": err.Error()})
		return
	}

	if exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already registered", "email": req.Email})
		return
	}

	// Check if password meets requirements
	if len(req.Password) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 8 characters long"})
		return
	}
	if !regexp.MustCompile(`[0-9]`).MatchString(req.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must contain at least one number"})
		return
	}
	if !regexp.MustCompile(`[A-Z]`).MatchString(req.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must contain at least one uppercase letter"})
		return
	}
	if !regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(req.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must contain at least one special character"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error", "err": err.Error()})
		return
	}
	req.Password = string(hashedPassword)

	err = h.US.Register(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error", "err": err.Error()})
		return
	}

	err = h.SendConfirmationCode(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error sending confirmation code", "err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Your account has been registered. Please check your email for a confirmation link. You have 3 minutes to confirm your account."})
}

// ConfirmRegistration godoc
// @Summary Confirm registration with code
// @Description Confirms a user's registration using the code sent to their email.
// @Tags registration
// @Accept json
// @Produce json
// @Param confirmation body models.ConfirmRegistrationReq true "Confirmation request"
// @Success 200 {object} token.Tokens "JWT tokens"
// @Failure 400 {object} string "Invalid request payload"
// @Failure 401 {object} string "Incorrect verification code"
// @Failure 404 {object} string "Verification code expired or email not found"
// @Router /confirm-registration [post]
func (h *HTTPHandler) ConfirmRegistration(c *gin.Context) {
	var req models.ConfirmRegistrationReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Invalid request payload": err.Error()})
		return
	}

	storedCode, err := rdb.Get(context.Background(), req.Email).Result()
	if err == redis.Nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Verification code expired or email not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "This email not found in confirmation requests!"})
		return
	}

	if storedCode != req.Code {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect verification code"})
		return
	}

	// Mark the user as confirmed in the database
	err = h.US.UM.ConfirmUser(&models.ConfirmUserReq{Email: req.Email})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error confirming user", "details": err.Error()})
		return
	}

	// Fetch the confirmed user
	user, err := h.US.GetProfile(&models.GetProfileReq{Email: req.Email})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching user", "details": err.Error()})
		return
	}

	// Generate JWT tokens
	tokens := token.GenerateJWTToken(user.ID, user.Email, user.Role)

	// Delete the confirmation code from Redis
	rdb.Del(context.Background(), req.Email)

	c.JSON(http.StatusOK, tokens)
}

// Login godoc
// @Summary Login a user
// @Description Authenticate user with email and password
// @Tags login
// @Accept json
// @Produce json
// @Param credentials body models.LoginReq true "User login credentials"
// @Success 200 {object} token.Tokens "JWT tokens"
// @Failure 400 {object} string "Invalid request payload"
// @Failure 401 {object} string "Invalid email or password"
// @Router /login [post]
func (h *HTTPHandler) Login(c *gin.Context) {
	req := models.LoginReq{}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Invalid request payload": err.Error()})
		return
	}

	user, err := h.US.GetProfile(&models.GetProfileReq{Email: req.Email})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"Invalid email or password": err.Error()})
		return
	}

	if !user.IsConfirmed {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Your account is not verified. Please check your email for a confirmation link."})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password!"})
		return
	}

	tokens := token.GenerateJWTToken(user.ID, user.Email, user.Role)

	c.JSON(http.StatusOK, tokens)
}

// ForgotPassword godoc
// @Summary Forgot passwrod
// @Description Sends a confirmation code to email recovery password
// @Tags password_recovery
// @Accept json
// @Produce json
// @Param credentials body models.ForgotPasswordReq true "User login credentials"
// @Success 200 {object} string ""
// @Failure 401 {object} string "Unauthorized"
// @Failure 404 {object} string "Page not found"
// @Failure 500 {object} string "Server error"
// @Security BearerAuth
// @Router /forgot-password [POST]
func (h *HTTPHandler) ForgotPassword(c *gin.Context) {
	var req models.ForgotPasswordReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Invalid request payload": err.Error()})
		return
	}

	if !isValidEmail(req.Email) {
		c.JSON(http.StatusBadRequest, "Invalid email format")
	}

	user, err := h.US.GetProfile(&models.GetProfileReq{Email: req.Email})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "details": err.Error()})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	err = h.SendConfirmationCode(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error sending confirmation code to email", "err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Confirmation code sent to your email. Please use your code within 3 minutes."})
}

// RecoverPassword godoc
// @Summary Recover password (Use this one after sending verification code)
// @Description Verifies the code and updates the password
// @Tags password_recovery
// @Accept json
// @Produce json
// @Param request body models.RecoverPasswordReq true "Recover Password Request"
// @Success 200 {object} string "Password successfully updated"
// @Failure 400 {object} string "Invalid request payload"
// @Failure 401 {object} string "Incorrect verification code"
// @Failure 404 {object} string "Verification code expired or email not found"
// @Failure 500 {object} string "Error updating password"
// @Router /recover-password [post]
func (h *HTTPHandler) RecoverPassword(c *gin.Context) {
	var req models.RecoverPasswordReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Invalid request payload": err.Error()})
		return
	}

	if req.Email == "" || req.Code == "" || req.NewPassword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email, code, and new password are required fields."})
		return
	}

	if len(req.NewPassword) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 8 characters long"})
		return
	}
	if !regexp.MustCompile(`[0-9]`).MatchString(req.NewPassword) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must contain at least one number"})
		return
	}
	if !regexp.MustCompile(`[A-Z]`).MatchString(req.NewPassword) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must contain at least one uppercase letter"})
		return
	}
	if !regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(req.NewPassword) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must contain at least one special character"})
		return
	}

	storedCode, err := rdb.Get(context.Background(), req.Email).Result()
	if err == redis.Nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Verification code expired or email not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "This email not found in a recovery requests!"})
		return
	}

	if storedCode != req.Code {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect verification code"})
		return
	}

	err = h.US.UM.UpdatePassword(&models.UpdatePasswordReq{Email: req.Email, NewPassword: req.NewPassword})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating password", "details": err.Error()})
		return
	}
	rdb.Del(context.Background(), req.Email)

	c.JSON(http.StatusOK, gin.H{"message": "Password successfully updated"})
}

// GetProfile godoc
// @Summary Get user profile
// @Description Get the profile of the authenticated user
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} models.GetProfileResp
// @Failure 401 {object} string "Unauthorized"
// @Failure 404 {object} string "User not found"
// @Failure 500 {object} string "Server error"
// @Security BearerAuth
// @Router /profile [get]
func (h *HTTPHandler) Profile(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	email := claims.(jwt.MapClaims)["email"].(string)
	user, err := h.US.GetProfile(&models.GetProfileReq{Email: email})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Server error": err.Error()})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *HTTPHandler) GetByID(c *gin.Context) {
	id := &models.GetProfileByIdReq{ID: c.Param("id")}
	user, err := h.US.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Couldn't get the user": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

// ChangeUserRole godoc
// @Summary Change a user's role
// @Description Changes the role of a user or admin. Only admins are allowed to use this function.
// @Tags admin-panel
// @Accept json
// @Produce json
// @Param id path string true "id or email of the user"
// @Param data query string true "Search with" Enums(id, email)
// @Param role query string true "New role of the user" Enums(admin, user, banned)
// @Success 200 {object} string "User role updated"
// @Failure 400 {object} string "Invalid request payload"
// @Failure 500 {object} string "Server error"
// @Security BearerAuth
// @Router /change-role/{id} [put]
func (h *HTTPHandler) ChangeUserRole(c *gin.Context) {
	id_or_email := c.Param("id")
	data := c.Query("data")
	role := c.Query("role")
	if data == "email" {
		err := h.US.ChangeRole(&models.ChangeRoleReq{Email: id_or_email, Role: role})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Couldn't change the role": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"User role updated": role})
	}
	if data == "id" {
		err := h.US.ChangeRole(&models.ChangeRoleReq{ID: id_or_email, Role: role})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Couldn't change the role": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"User role updated": role})
	}
}

func isValidEmail(email string) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`).MatchString(email)
}
