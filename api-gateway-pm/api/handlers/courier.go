package handlers

import (
	"net/http"

	pb "gateway-admin/genprotos"

	"github.com/gin-gonic/gin"
)

// AddProduct godoc
// @Summary Add a product
// @Description Adds a product to the system. Only admins are allowed to use this function.
// @Tags product
// @Accept json
// @Produce json
// @Param data body pb.ProductCReq true "Product data"
// @Success 200 {object} string "Product is added"
// @Failure 400 {object} string "Invalid request payload"
// @Failure 500 {object} string "Server error"
// @Security BearerAuth
// @Router /add-product [post]
func (h *HTTPHandler) AddProduct(c *gin.Context) {
	var req pb.ProductCReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Invalid request payload": err.Error()})
		return
	}
}

// DeleteProduct godoc
// @Summary Delete a product
// @Description Deletes a product from the system. Only admins are allowed to use this function.
// @Tags product
// @Accept json
// @Produce json
// @Param id path string true "id or email of the product"
// @Param data query string true "Search with" Enums(id, email)
// @Success 200 {object} string "Product is deleted"
// @Failure 400 {object} string "Invalid request payload"
// @Failure 500 {object} string "Server error"
// @Security BearerAuth
// @Router /delete-product/{id} [delete]
func (h *HTTPHandler) DeleteProduct(c *gin.Context) {

}
