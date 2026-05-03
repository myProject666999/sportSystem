package controllers

import (
	"sportSystem/database"
	"sportSystem/models"
	"sportSystem/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetCartList(c *gin.Context) {
	userID := c.GetUint("user_id")

	var cartItems []models.Cart
	database.DB.Where("user_id = ?", userID).Preload("Product").Find(&cartItems)

	utils.Success(c, cartItems)
}

func AddToCart(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req struct {
		ProductID uint `json:"product_id" binding:"required"`
		Quantity  int  `json:"quantity"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var product models.Product
	if err := database.DB.First(&product, req.ProductID).Error; err != nil {
		utils.NotFound(c, "商品不存在")
		return
	}

	if product.Status != 1 {
		utils.BadRequest(c, "商品已下架")
		return
	}

	var existingCart models.Cart
	if database.DB.Where("user_id = ? AND product_id = ?", userID, req.ProductID).First(&existingCart).Error == nil {
		newQuantity := existingCart.Quantity + req.Quantity
		if newQuantity <= 0 {
			newQuantity = 1
		}
		database.DB.Model(&existingCart).Update("quantity", newQuantity)
		utils.SuccessWithMessage(c, "已更新购物车", existingCart)
		return
	}

	quantity := req.Quantity
	if quantity <= 0 {
		quantity = 1
	}

	cartItem := &models.Cart{
		UserID:    userID,
		ProductID: req.ProductID,
		Quantity:  quantity,
	}

	if err := database.DB.Create(cartItem).Error; err != nil {
		utils.ServerError(c, "添加购物车失败")
		return
	}

	utils.SuccessWithMessage(c, "添加成功", cartItem)
}

func UpdateCart(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, _ := strconv.Atoi(c.Param("id"))
	var req struct {
		Quantity int `json:"quantity" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	if req.Quantity <= 0 {
		utils.BadRequest(c, "数量必须大于0")
		return
	}

	var cartItem models.Cart
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&cartItem).Error; err != nil {
		utils.NotFound(c, "购物车项目不存在")
		return
	}

	database.DB.Model(&cartItem).Update("quantity", req.Quantity)
	utils.SuccessWithMessage(c, "更新成功", cartItem)
}

func RemoveFromCart(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, _ := strconv.Atoi(c.Param("id"))

	result := database.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Cart{})
	if result.RowsAffected == 0 {
		utils.NotFound(c, "购物车项目不存在")
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}

func GetFavorites(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	offset := (page - 1) * pageSize

	var favorites []models.Favorite
	var total int64

	database.DB.Model(&models.Favorite{}).Where("user_id = ?", userID).Count(&total)
	database.DB.Where("user_id = ?", userID).Preload("Product").
		Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&favorites)

	utils.SuccessPage(c, favorites, total, page, pageSize)
}

func AddFavorite(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req struct {
		ProductID uint `json:"product_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var product models.Product
	if err := database.DB.First(&product, req.ProductID).Error; err != nil {
		utils.NotFound(c, "商品不存在")
		return
	}

	var existingFavorite models.Favorite
	if database.DB.Where("user_id = ? AND product_id = ?", userID, req.ProductID).First(&existingFavorite).Error == nil {
		utils.SuccessWithMessage(c, "已收藏", existingFavorite)
		return
	}

	favorite := &models.Favorite{
		UserID:    userID,
		ProductID: req.ProductID,
	}

	if err := database.DB.Create(favorite).Error; err != nil {
		utils.ServerError(c, "收藏失败")
		return
	}

	utils.SuccessWithMessage(c, "收藏成功", favorite)
}

func RemoveFavorite(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, _ := strconv.Atoi(c.Param("id"))

	result := database.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Favorite{})
	if result.RowsAffected == 0 {
		utils.NotFound(c, "收藏不存在")
		return
	}

	utils.SuccessWithMessage(c, "取消收藏成功", nil)
}

func CheckFavorite(c *gin.Context) {
	userID := c.GetUint("user_id")
	productID, _ := strconv.Atoi(c.Param("product_id"))

	var favorite models.Favorite
	isFavorited := database.DB.Where("user_id = ? AND product_id = ?", userID, productID).First(&favorite).Error == nil

	utils.Success(c, gin.H{
		"is_favorited": isFavorited,
	})
}
