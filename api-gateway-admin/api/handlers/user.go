package handlers

import (
	"auth-service/config"
	"auth-service/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// BanUser godoc
// @Summary Ban a user
// @Description Bans a user which doesn't allow user to use. Only admins are allowed to use this function.
// @Tags banning
// @Accept json
// @Produce json
// @Param id path string true "id or email of the user"
// @Param data query string true "Search with" Enums(id, email)
// @Success 200 {object} string "User is banned"
// @Failure 400 {object} string "Invalid request payload"
// @Failure 500 {object} string "Server error"
// @Security BearerAuth
// @Router /ban/{id} [put]
func (h *HTTPHandler) BanUser(c *gin.Context) {
	id_or_email := c.Param("id")
	data := c.Query("data")
	if data == "email" {
		err := h.US.BanUser(&models.BanUserReq{Email: id_or_email})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Couldn't ban user": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"User is banned!": id_or_email})
	} else if data == "id" {
		if err := config.IsValidUUID(id_or_email); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err := h.US.BanUser(&models.BanUserReq{ID: id_or_email})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Couldn't ban user": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"User is banned": id_or_email})
	}
}

// UnbanUser godoc
// @Summary Unban a user
// @Description Unbans a user which allows user to use. Only admins are allowed to use this function.
// @Tags banning
// @Accept json
// @Produce json
// @Param id path string true "id or email of the user"
// @Param data query string true "Search with" Enums(id, email)
// @Success 200 {object} string "User is unbanned"
// @Failure 400 {object} string "Invalid request payload"
// @Failure 500 {object} string "Server error"
// @Security BearerAuth
// @Router /unban/{id} [put]
func (h *HTTPHandler) UnbanUser(c *gin.Context) {
	id_or_email := c.Param("id")
	data := c.Query("data")
	if data == "email" {
		err := h.US.UnbanUser(&models.UnbanUserReq{Email: id_or_email})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Couldn't unban user": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"User is unbanned!": id_or_email})
	} else if data == "id" {
		if err := config.IsValidUUID(id_or_email); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err := h.US.UnbanUser(&models.UnbanUserReq{ID: id_or_email})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Couldn't unban user": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"User is unbanned": id_or_email})
	}
}
