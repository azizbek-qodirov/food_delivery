package handlers

import (
	"auth-service/config"
	"auth-service/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AddProductManager godoc
// @Summary Add a product-manager
// @Description Adds a product-manager to the system. Only admins are allowed to use this function.
// @Tags product-manager
// @Accept json
// @Produce json
// @Param data body models.AddProductManagerReq true "ProductManager data"
// @Success 200 {object} string "ProductManager is added"
// @Failure 400 {object} string "Invalid request payload"
// @Failure 500 {object} string "Server error"
// @Security BearerAuth
// @Router /add-product-manager [post]
func (h *HTTPHandler) AddProductManager(c *gin.Context) {
	var req models.AddProductManagerReq
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

	err = h.US.AddProductManager(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Couldn't add product-manager": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ProductManager is added": req.Email})
}

// DeleteProductManager godoc
// @Summary Delete a product-manager
// @Description Deletes a product-manager from the system. Only admins are allowed to use this function.
// @Tags product-manager
// @Accept json
// @Produce json
// @Param id path string true "id or email of the product-manager"
// @Param data query string true "Search with" Enums(id, email)
// @Success 200 {object} string "ProductManager is deleted"
// @Failure 400 {object} string "Invalid request payload"
// @Failure 500 {object} string "Server error"
// @Security BearerAuth
// @Router /delete-product-manager/{id} [delete]
func (h *HTTPHandler) DeleteProductManager(c *gin.Context) {
	id_or_email := c.Param("id")
	data := c.Query("data")
	if data == "email" {
		if !config.IsValidEmail(id_or_email) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
			return
		}
		err := h.US.DeleteProductManager(&models.DeleteProductManagerReq{Email: id_or_email})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Couldn't delete product-manager": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"ProductManager is deleted!": id_or_email})
	} else if data == "id" {
		err := h.US.DeleteProductManager(&models.DeleteProductManagerReq{ID: id_or_email})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Couldn't delete product-manager": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"ProductManager is deleted": id_or_email})
	}
}
