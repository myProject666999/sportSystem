package controllers

import (
	"sportSystem/database"
	"sportSystem/models"
	"sportSystem/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetCategoryList(c *gin.Context) {
	var categories []models.Category
	database.DB.Where("status = 1").Order("sort ASC").Find(&categories)
	utils.Success(c, categories)
}

func GetAllCategories(c *gin.Context) {
	var categories []models.Category
	database.DB.Order("sort ASC").Find(&categories)
	utils.Success(c, categories)
}

func CreateCategory(c *gin.Context) {
	var req struct {
		Name   string `json:"name" binding:"required"`
		Sort   int    `json:"sort"`
		Status int    `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	category := &models.Category{
		Name:   req.Name,
		Sort:   req.Sort,
		Status: req.Status,
	}

	if err := database.DB.Create(category).Error; err != nil {
		utils.ServerError(c, "创建失败")
		return
	}

	utils.SuccessWithMessage(c, "创建成功", category)
}

func UpdateCategory(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req struct {
		Name   string `json:"name"`
		Sort   int    `json:"sort"`
		Status int    `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Sort != 0 {
		updates["sort"] = req.Sort
	}
	updates["status"] = req.Status

	if len(updates) > 0 {
		if err := database.DB.Model(&models.Category{}).Where("id = ?", id).Updates(updates).Error; err != nil {
			utils.ServerError(c, "更新失败")
			return
		}
	}

	var category models.Category
	database.DB.First(&category, id)
	utils.SuccessWithMessage(c, "更新成功", category)
}

func DeleteCategory(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := database.DB.Delete(&models.Category{}, id).Error; err != nil {
		utils.ServerError(c, "删除失败")
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}

func GetProductList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	categoryID, _ := strconv.Atoi(c.DefaultQuery("category_id", "0"))
	keyword := c.Query("keyword")
	status := c.Query("status")

	offset := (page - 1) * pageSize

	var products []models.Product
	var total int64

	query := database.DB.Model(&models.Product{}).Preload("Category")

	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}
	if keyword != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	} else {
		query = query.Where("status = 1")
	}

	query.Count(&total)
	query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&products)

	utils.SuccessPage(c, products, total, page, pageSize)
}

func GetProductDetail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var product models.Product
	if err := database.DB.Preload("Category").First(&product, id).Error; err != nil {
		utils.NotFound(c, "商品不存在")
		return
	}

	database.DB.Model(&models.Product{}).Where("id = ?", id).Update("sales", database.DB.Raw("sales + 1"))

	var comments []models.ProductComment
	database.DB.Where("product_id = ? AND status = 1", id).Preload("User").Order("created_at DESC").Limit(10).Find(&comments)

	utils.Success(c, gin.H{
		"product":  product,
		"comments": comments,
	})
}

func CreateProduct(c *gin.Context) {
	var req struct {
		CategoryID  uint    `json:"category_id" binding:"required"`
		Name        string  `json:"name" binding:"required"`
		Description string  `json:"description"`
		Price       float64 `json:"price" binding:"required"`
		PointsPrice int     `json:"points_price"`
		Stock       int     `json:"stock"`
		Image       string  `json:"image"`
		Images      string  `json:"images"`
		Status      int     `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	product := &models.Product{
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		PointsPrice: req.PointsPrice,
		Stock:       req.Stock,
		Image:       req.Image,
		Images:      req.Images,
		Status:      req.Status,
	}

	if err := database.DB.Create(product).Error; err != nil {
		utils.ServerError(c, "创建失败")
		return
	}

	utils.SuccessWithMessage(c, "创建成功", product)
}

func UpdateProduct(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req struct {
		CategoryID  uint    `json:"category_id"`
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		PointsPrice int     `json:"points_price"`
		Stock       int     `json:"stock"`
		Image       string  `json:"image"`
		Images      string  `json:"images"`
		Status      int     `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	updates := make(map[string]interface{})
	if req.CategoryID != 0 {
		updates["category_id"] = req.CategoryID
	}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Price != 0 {
		updates["price"] = req.Price
	}
	if req.PointsPrice != 0 {
		updates["points_price"] = req.PointsPrice
	}
	if req.Stock != 0 {
		updates["stock"] = req.Stock
	}
	if req.Image != "" {
		updates["image"] = req.Image
	}
	if req.Images != "" {
		updates["images"] = req.Images
	}
	updates["status"] = req.Status

	if len(updates) > 0 {
		if err := database.DB.Model(&models.Product{}).Where("id = ?", id).Updates(updates).Error; err != nil {
			utils.ServerError(c, "更新失败")
			return
		}
	}

	var product models.Product
	database.DB.Preload("Category").First(&product, id)
	utils.SuccessWithMessage(c, "更新成功", product)
}

func DeleteProduct(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := database.DB.Delete(&models.Product{}, id).Error; err != nil {
		utils.ServerError(c, "删除失败")
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}

func AddProductComment(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req struct {
		ProductID uint   `json:"product_id" binding:"required"`
		Content   string `json:"content" binding:"required"`
		Rating    int    `json:"rating"`
		Images    string `json:"images"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	comment := &models.ProductComment{
		ProductID: req.ProductID,
		UserID:    userID,
		Content:   req.Content,
		Rating:    req.Rating,
		Images:    req.Images,
	}

	if comment.Rating == 0 {
		comment.Rating = 5
	}

	if err := database.DB.Create(comment).Error; err != nil {
		utils.ServerError(c, "评论失败")
		return
	}

	utils.SuccessWithMessage(c, "评论成功", comment)
}

func GetProductComments(c *gin.Context) {
	productID, _ := strconv.Atoi(c.Param("id"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	offset := (page - 1) * pageSize

	var comments []models.ProductComment
	var total int64

	database.DB.Model(&models.ProductComment{}).Where("product_id = ? AND status = 1", productID).Count(&total)
	database.DB.Where("product_id = ? AND status = 1", productID).Preload("User").
		Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&comments)

	utils.SuccessPage(c, comments, total, page, pageSize)
}
