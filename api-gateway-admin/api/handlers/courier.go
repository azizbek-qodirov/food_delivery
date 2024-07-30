package handlers

import (
	"auth-service/config"
	"auth-service/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AddCourier godoc
// @Summary Add a courier
// @Description Adds a courier to the system. Only admins are allowed to use this function.
// @Tags courier
// @Accept json
// @Produce json
// @Param data body models.AddCourierReq true "Courier data"
// @Success 200 {object} string "Courier is added"
// @Failure 400 {object} string "Invalid request payload"
// @Failure 500 {object} string "Server error"
// @Security BearerAuth
// @Router /add-courier [post]
func (h *HTTPHandler) AddCourier(c *gin.Context) {
	var req models.AddCourierReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Invalid request payload": err.Error()})
		return
	}
	if !config.IsValidEmail(req.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	if err := h.US.IsEmailExists(req.Email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.IsValidPassword(req.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := config.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error", "err": err.Error()})
	}

	req.Password = string(hashedPassword)

	err = h.US.AddCourier(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Couldn't add courier": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Courier is added": req.Email})
}

// DeleteCourier godoc
// @Summary Delete a courier
// @Description Deletes a courier from the system. Only admins are allowed to use this function.
// @Tags courier
// @Accept json
// @Produce json
// @Param id path string true "id or email of the courier"
// @Param data query string true "Search with" Enums(id, email)
// @Success 200 {object} string "Courier is deleted"
// @Failure 400 {object} string "Invalid request payload"
// @Failure 500 {object} string "Server error"
// @Security BearerAuth
// @Router /delete-courier/{id} [delete]
func (h *HTTPHandler) DeleteCourier(c *gin.Context) {
	id_or_email := c.Param("id")
	data := c.Query("data")
	if data == "email" {
		if !config.IsValidEmail(id_or_email) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
			return
		}
		err := h.US.DeleteCourier(&models.DeleteCourierReq{Email: id_or_email})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Couldn't delete courier": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"Courier is deleted!": id_or_email})
	} else if data == "id" {
		err := h.US.DeleteCourier(&models.DeleteCourierReq{ID: id_or_email})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Couldn't delete courier": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"Courier is deleted": id_or_email})
	}
}
