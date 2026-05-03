package controllers

import (
	"sportSystem/database"
	"sportSystem/models"
	"sportSystem/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type UserRegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

type UserLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UpdateUserRequest struct {
	Email  string `json:"email"`
	Phone  string `json:"phone"`
	Avatar string `json:"avatar"`
}

type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

func UserRegister(c *gin.Context) {
	var req UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var existingUser models.User
	if database.DB.Where("username = ?", req.Username).First(&existingUser).Error == nil {
		utils.BadRequest(c, "用户名已存在")
		return
	}

	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Phone:    req.Phone,
		Points:   100,
	}
	if err := user.SetPassword(req.Password); err != nil {
		utils.ServerError(c, "密码加密失败")
		return
	}

	if err := database.DB.Create(user).Error; err != nil {
		utils.ServerError(c, "注册失败")
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Username, "user")
	if err != nil {
		utils.ServerError(c, "生成token失败")
		return
	}

	utils.Success(c, gin.H{
		"token": token,
		"user":  user,
	})
}

func UserLogin(c *gin.Context) {
	var req UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var user models.User
	if err := database.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		utils.BadRequest(c, "用户名或密码错误")
		return
	}

	if !user.CheckPassword(req.Password) {
		utils.BadRequest(c, "用户名或密码错误")
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Username, "user")
	if err != nil {
		utils.ServerError(c, "生成token失败")
		return
	}

	utils.Success(c, gin.H{
		"token": token,
		"user":  user,
	})
}

func GetUserInfo(c *gin.Context) {
	userID := c.GetUint("user_id")

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	utils.Success(c, user)
}

func UpdateUserInfo(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req UpdateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	updates := make(map[string]interface{})
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}

	if len(updates) > 0 {
		if err := database.DB.Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
			utils.ServerError(c, "更新失败")
			return
		}
	}

	var user models.User
	database.DB.First(&user, userID)
	utils.Success(c, user)
}

func UpdateUserPassword(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req UpdatePasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	if !user.CheckPassword(req.OldPassword) {
		utils.BadRequest(c, "原密码错误")
		return
	}

	if err := user.SetPassword(req.NewPassword); err != nil {
		utils.ServerError(c, "密码加密失败")
		return
	}

	if err := database.DB.Save(&user).Error; err != nil {
		utils.ServerError(c, "更新密码失败")
		return
	}

	utils.SuccessWithMessage(c, "密码修改成功", nil)
}

func Recharge(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req struct {
		Amount float64 `json:"amount" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	if req.Amount <= 0 {
		utils.BadRequest(c, "充值金额必须大于0")
		return
	}

	orderNo := "RC" + time.Now().Format("20060102150405") + strconv.Itoa(int(time.Now().UnixNano())%10000)
	recharge := &models.Recharge{
		UserID:  userID,
		OrderNo: orderNo,
		Amount:  req.Amount,
		Status:  1,
	}

	if err := database.DB.Create(recharge).Error; err != nil {
		utils.ServerError(c, "创建充值订单失败")
		return
	}

	now := time.Now()
	recharge.PayTime = &now
	if err := database.DB.Model(&models.Recharge{}).Where("id = ?", recharge.ID).Updates(
		map[string]interface{}{"status": 1, "pay_time": now},
	).Error; err != nil {
		utils.ServerError(c, "更新充值状态失败")
		return
	}

	if err := database.DB.Model(&models.User{}).Where("id = ?", userID).
		Update("balance", database.DB.Raw("balance + ?", req.Amount)).Error; err != nil {
		utils.ServerError(c, "充值失败")
		return
	}

	utils.SuccessWithMessage(c, "充值成功", gin.H{
		"order_no": orderNo,
		"amount":   req.Amount,
	})
}
