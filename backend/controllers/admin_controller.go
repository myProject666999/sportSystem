package controllers

import (
	"sportSystem/database"
	"sportSystem/models"
	"sportSystem/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type AdminLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func AdminLogin(c *gin.Context) {
	var req AdminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var admin models.Admin
	if err := database.DB.Where("username = ?", req.Username).First(&admin).Error; err != nil {
		utils.BadRequest(c, "用户名或密码错误")
		return
	}

	if !admin.CheckPassword(req.Password) {
		utils.BadRequest(c, "用户名或密码错误")
		return
	}

	token, err := utils.GenerateToken(admin.ID, admin.Username, admin.Role)
	if err != nil {
		utils.ServerError(c, "生成token失败")
		return
	}

	utils.Success(c, gin.H{
		"token": token,
		"admin": admin,
	})
}

func GetAdminInfo(c *gin.Context) {
	adminID := c.GetUint("admin_id")

	var admin models.Admin
	if err := database.DB.First(&admin, adminID).Error; err != nil {
		utils.NotFound(c, "管理员不存在")
		return
	}

	utils.Success(c, admin)
}

func GetDashboardStats(c *gin.Context) {
	var userCount, orderCount, productCount, newsCount, todayOrders int64
	var totalSales float64

	database.DB.Model(&models.User{}).Count(&userCount)
	database.DB.Model(&models.Order{}).Count(&orderCount)
	database.DB.Model(&models.Product{}).Count(&productCount)
	database.DB.Model(&models.News{}).Count(&newsCount)

	today := time.Now().Format("2006-01-02")
	database.DB.Model(&models.Order{}).Where("DATE(created_at) = ? AND status >= 1", today).Count(&todayOrders)
	database.DB.Model(&models.Order{}).Where("status >= 1").Select("COALESCE(SUM(total_price), 0)").Scan(&totalSales)

	utils.Success(c, gin.H{
		"user_count":    userCount,
		"order_count":   orderCount,
		"product_count": productCount,
		"news_count":    newsCount,
		"today_orders":  todayOrders,
		"total_sales":   totalSales,
	})
}

func GetUserList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	keyword := c.Query("keyword")

	offset := (page - 1) * pageSize

	var users []models.User
	var total int64

	query := database.DB.Model(&models.User{})
	if keyword != "" {
		query = query.Where("username LIKE ? OR email LIKE ? OR phone LIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	query.Count(&total)
	query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&users)

	utils.SuccessPage(c, users, total, page, pageSize)
}

func GetAdminList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	offset := (page - 1) * pageSize

	var admins []models.Admin
	var total int64

	database.DB.Model(&models.Admin{}).Count(&total)
	database.DB.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&admins)

	utils.SuccessPage(c, admins, total, page, pageSize)
}

func CreateAdmin(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Role     string `json:"role"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var existingAdmin models.Admin
	if database.DB.Where("username = ?", req.Username).First(&existingAdmin).Error == nil {
		utils.BadRequest(c, "用户名已存在")
		return
	}

	admin := &models.Admin{
		Username: req.Username,
		Role:     req.Role,
	}
	if admin.Role == "" {
		admin.Role = "admin"
	}

	if err := admin.SetPassword(req.Password); err != nil {
		utils.ServerError(c, "密码加密失败")
		return
	}

	if err := database.DB.Create(admin).Error; err != nil {
		utils.ServerError(c, "创建管理员失败")
		return
	}

	utils.SuccessWithMessage(c, "创建成功", admin)
}

func UpdateAdmin(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var admin models.Admin
	if err := database.DB.First(&admin, id).Error; err != nil {
		utils.NotFound(c, "管理员不存在")
		return
	}

	updates := make(map[string]interface{})
	if req.Username != "" {
		var existingAdmin models.Admin
		if database.DB.Where("username = ? AND id != ?", req.Username, id).First(&existingAdmin).Error == nil {
			utils.BadRequest(c, "用户名已存在")
			return
		}
		updates["username"] = req.Username
	}
	if req.Role != "" {
		updates["role"] = req.Role
	}

	if len(updates) > 0 {
		if err := database.DB.Model(&models.Admin{}).Where("id = ?", id).Updates(updates).Error; err != nil {
			utils.ServerError(c, "更新失败")
			return
		}
	}

	if req.Password != "" {
		if err := admin.SetPassword(req.Password); err != nil {
			utils.ServerError(c, "密码加密失败")
			return
		}
		database.DB.Save(&admin)
	}

	database.DB.First(&admin, id)
	utils.SuccessWithMessage(c, "更新成功", admin)
}

func DeleteAdmin(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if id == 1 {
		utils.BadRequest(c, "不能删除超级管理员")
		return
	}

	if err := database.DB.Delete(&models.Admin{}, id).Error; err != nil {
		utils.ServerError(c, "删除失败")
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}
