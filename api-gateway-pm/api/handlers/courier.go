package handlers

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	pb "gateway-admin/genprotos"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
)

// AddProduct godoc
// @Summary Add a product
// @Description Adds a product image to the system. Only admins are allowed to use this function.
// @Tags product
// @Accept multipart/mixed
// @Produce json
// @Param image formData file false "Product image"
// @Success 200 {object} string "Product image is added"
// @Failure 400 {object} string "Invalid request payload"
// @Failure 500 {object} string "Server error"
// @Security BearerAuth
// @Router /add-product [POST]
func (h *HTTPHandler) AddProductImage(c *gin.Context) {
	image, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, "image file is required: "+err.Error())
		return
	}
	defer image.Close()

	filename := header.Filename
	filepath := filename
	outfile, err := os.Create(filepath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "failed to create file"+err.Error())
		return
	}
	defer outfile.Close()

	_, err = io.Copy(outfile, image)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "failed to save image")
		return
	}

	objectName := filepath
	_, err = h.MinIO.FPutObject(context.Background(), "another-bucket", objectName, filepath, minio.PutObjectOptions{ContentType: "image/png"})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Message": err.Error(),
		})
		log.Println(err)
		return
	}

	minioURL := fmt.Sprintf("http://localhost:9000/another-bucket/%s", objectName)

	c.JSON(http.StatusCreated, gin.H{
		"url": minioURL,
	})

	if err := c.SaveUploadedFile(header, filepath); err != nil {
		c.JSON(500, gin.H{
			"error": "Couln't find matching file format",
		})
		log.Println(err)
		return
	}
	c.JSON(200, gin.H{"message": "success", "url": minioURL})
	
}

// AddProduct godoc
// @Summary Add a product
// @Description Adds a product to the system. Only admins are allowed to use this function.
// @Tags product
// @Accept multipart/mixed
// @Produce json
// @Param data body pb.ProductCReqForSwagger true "Product data"
// @Success 200 {object} string "Product is added"
// @Failure 400 {object} string "Invalid request payload"
// @Failure 500 {object} string "Server error"
// @Security BearerAuth
// @Router /add-product [POST]
func (h *HTTPHandler) AddProduct(c *gin.Context) {
	req := pb.ProductCReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, "invalid request payload")
		return
	}

	_, err := h.ProductManager.Create(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "failed to add product")
		return
	}

	c.JSON(http.StatusOK, "product is added")
}

// UpdateProduct godoc
// @Summary Update a product
// @Description Updates a product in the system. Only admins are allowed to use this function.
// @Tags product
// @Accept json
// @Produce json
// @Param id path string true "id or email of the product"
// @Param data query string true "Search with" Enums(id, email)
// @Param data body pb.ProductUReq true "Product data"
// @Success 200 {object} string "Product is updated"
// @Failure 400 {object} string "Invalid request payload"
// @Failure 500 {object} string "Server error"
// @Security BearerAuth
// @Router /update-product/{id} [PUT]
func (h *HTTPHandler) UpdateProduct(c *gin.Context) {
	var (
		req pb.ProductUReq
	)
	req.Id = c.Param("id")

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid request payload")
		return
	}

	_, err = h.ProductManager.Update(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "failed to update product")
		return
	}

	c.JSON(http.StatusOK, "product is updated")
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
// @Router /delete-product/{id} [DELETE]
func (h *HTTPHandler) DeleteProduct(c *gin.Context) {
	res, err := h.ProductManager.Delete(context.Background(), &pb.ByID{
		Id: c.Param("id"),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, "failed to delete product")
		return
	}

	c.JSON(http.StatusOK, res)
}

// GetProduct godoc
// @Summary Get a product
// @Description Gets a product from the system.
// @Tags product
// @Accept json
// @Produce json
// @Param id path string true "id or email of the product"
// @Param data query string true "Search with" Enums(id, email)
// @Success 200 {object} pb.ProductGRes "Product data"
// @Failure 400 {object} string "Invalid request payload"
// @Failure 500 {object} string "Server error"
// @Router /get-product/{id} [GET]
func (h *HTTPHandler) GetProduct(c *gin.Context) {
	res, err := h.ProductManager.Get(context.Background(), &pb.ByID{
		Id: c.Param("id"),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, "failed to get product")
		return
	}

	c.JSON(http.StatusOK, res)
}

// GetProducts godoc
// @Summary Get all products
// @Description Gets all products from the system.
// @Tags product
// @Accept json
// @Produce json
// @Success 200 {array} pb.ProductGARes "Product data"
// @Failure 400 {object} string "Invalid request payload"
// @Failure 500 {object} string "Server error"
// @Router /get-products [GET]
func (h *HTTPHandler) GetAllProducts(c *gin.Context) {
	res, err := h.ProductManager.GetAll(context.Background(), &pb.ProductGAReq{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, "failed to get products")
		return
	}

	c.JSON(http.StatusOK, res)
}

// UpdateProductRating godoc
// @Summary Update product rating
// @Description Updates the rating of a product.
// @Tags product
// @Accept json
// @Produce json
// @Param product_id path string true "Product ID"
// @Param data body pb.ProductRatingUReq true "Rating data"
// @Success 200 {object} string "Product rating updated"
// @Failure 400 {object} string "Invalid request payload"
// @Failure 500 {object} string "Server error"
// @Router /update-product-rating/{product_id} [PUT]
func (h *HTTPHandler) UpdateProductRating(c *gin.Context) {
	var req pb.ProductRatingUReq
	req.ProductId = c.Param("product_id")

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, "invalid request payload")
		return
	}

	_, err := h.ProductManager.UpdateRating(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "failed to update product rating")
		return
	}

	c.JSON(http.StatusOK, "product rating updated")
}

// UpdateProductCount godoc
// @Summary Update product count
// @Description Updates the count of a product.
// @Tags product
// @Accept json
// @Produce json
// @Param product_id path string true "Product ID"
// @Param data body pb.ProductCountUReq true "Count data"
// @Success 200 {object} string "Product count updated"
// @Failure 400 {object} string "Invalid request payload"
// @Failure 500 {object} string "Server error"
// @Router /update-product-count/{product_id} [PUT]
func (h *HTTPHandler) UpdateProductCount(c *gin.Context) {
	var req pb.ProductCountUReq
	req.ProductId = c.Param("product_id")

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, "invalid request payload")
		return
	}

	_, err := h.ProductManager.UpdateCount(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "failed to update product count")
		return
	}

	c.JSON(http.StatusOK, "product count updated")
}

// UpdateProductImage godoc
// @Summary Update product image URL
// @Description Updates the image URL of a product.
// @Tags product
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param data body pb.ProductImageUReq true "Image URL data"
// @Success 200 {object} string "Product image URL updated"
// @Failure 400 {object} string "Invalid request payload"
// @Failure 500 {object} string "Server error"
// @Router /update-product-image/{id} [PUT]
func (h *HTTPHandler) UpdateProductImage(c *gin.Context) {
	var req pb.ProductImageUReq
	req.Id = c.Param("id")

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, "invalid request payload")
		return
	}

	_, err := h.ProductManager.UpdateImg(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "failed to update product image URL")
		return
	}

	c.JSON(http.StatusOK, "product image URL updated")
}
