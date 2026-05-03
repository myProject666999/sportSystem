package controllers

import (
	"sportSystem/database"
	"sportSystem/models"
	"sportSystem/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetAddressList(c *gin.Context) {
	userID := c.GetUint("user_id")

	var addresses []models.Address
	database.DB.Where("user_id = ?", userID).Order("is_default DESC, created_at DESC").Find(&addresses)

	utils.Success(c, addresses)
}

func GetDefaultAddress(c *gin.Context) {
	userID := c.GetUint("user_id")

	var address models.Address
	if err := database.DB.Where("user_id = ? AND is_default = 1", userID).First(&address).Error; err != nil {
		database.DB.Where("user_id = ?", userID).Order("created_at DESC").First(&address)
	}

	utils.Success(c, address)
}

func CreateAddress(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req struct {
		Name      string `json:"name" binding:"required"`
		Phone     string `json:"phone" binding:"required"`
		Province  string `json:"province"`
		City      string `json:"city"`
		District  string `json:"district"`
		Detail    string `json:"detail" binding:"required"`
		IsDefault int    `json:"is_default"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	if req.IsDefault == 1 {
		database.DB.Model(&models.Address{}).Where("user_id = ?", userID).Update("is_default", 0)
	}

	address := &models.Address{
		UserID:    userID,
		Name:      req.Name,
		Phone:     req.Phone,
		Province:  req.Province,
		City:      req.City,
		District:  req.District,
		Detail:    req.Detail,
		IsDefault: req.IsDefault,
	}

	if err := database.DB.Create(address).Error; err != nil {
		utils.ServerError(c, "创建地址失败")
		return
	}

	utils.SuccessWithMessage(c, "创建成功", address)
}

func UpdateAddress(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, _ := strconv.Atoi(c.Param("id"))
	var req struct {
		Name      string `json:"name"`
		Phone     string `json:"phone"`
		Province  string `json:"province"`
		City      string `json:"city"`
		District  string `json:"district"`
		Detail    string `json:"detail"`
		IsDefault int    `json:"is_default"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var address models.Address
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&address).Error; err != nil {
		utils.NotFound(c, "地址不存在")
		return
	}

	if req.IsDefault == 1 {
		database.DB.Model(&models.Address{}).Where("user_id = ?", userID).Update("is_default", 0)
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.Province != "" {
		updates["province"] = req.Province
	}
	if req.City != "" {
		updates["city"] = req.City
	}
	if req.District != "" {
		updates["district"] = req.District
	}
	if req.Detail != "" {
		updates["detail"] = req.Detail
	}
	updates["is_default"] = req.IsDefault

	database.DB.Model(&address).Updates(updates)
	database.DB.First(&address, id)

	utils.SuccessWithMessage(c, "更新成功", address)
}

func DeleteAddress(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, _ := strconv.Atoi(c.Param("id"))

	result := database.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Address{})
	if result.RowsAffected == 0 {
		utils.NotFound(c, "地址不存在")
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}

type CreateOrderRequest struct {
	AddressID uint   `json:"address_id" binding:"required"`
	Items     []struct {
		ProductID uint `json:"product_id" binding:"required"`
		Quantity  int  `json:"quantity" binding:"required"`
	} `json:"items" binding:"required"`
	PaymentType string `json:"payment_type"`
}

func CreateOrder(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req CreateOrderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	if len(req.Items) == 0 {
		utils.BadRequest(c, "请选择商品")
		return
	}

	var address models.Address
	if err := database.DB.Where("id = ? AND user_id = ?", req.AddressID, userID).First(&address).Error; err != nil {
		utils.NotFound(c, "地址不存在")
		return
	}

	totalPrice := 0.0
	totalPoints := 0
	var orderItems []models.OrderItem

	for _, item := range req.Items {
		var product models.Product
		if err := database.DB.First(&product, item.ProductID).Error; err != nil {
			utils.NotFound(c, "商品不存在")
			return
		}

		if product.Status != 1 {
			utils.BadRequest(c, "商品"+product.Name+"已下架")
			return
		}

		if product.Stock < item.Quantity {
			utils.BadRequest(c, "商品"+product.Name+"库存不足")
			return
		}

		orderItem := models.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     product.Price,
			Points:    product.PointsPrice * item.Quantity,
		}

		totalPrice += product.Price * float64(item.Quantity)
		totalPoints += product.PointsPrice * item.Quantity
		orderItems = append(orderItems, orderItem)
	}

	orderNo := "ORD" + time.Now().Format("20060102150405") + strconv.Itoa(int(time.Now().UnixNano())%10000)

	order := &models.Order{
		OrderNo:         orderNo,
		UserID:          userID,
		TotalPrice:      totalPrice,
		TotalPoints:     totalPoints,
		Status:          0,
		ReceiverName:    address.Name,
		ReceiverPhone:   address.Phone,
		ReceiverAddress: address.Province + address.City + address.District + address.Detail,
	}

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		for i := range orderItems {
			orderItems[i].OrderID = order.ID
		}

		if err := tx.Create(&orderItems).Error; err != nil {
			return err
		}

		for _, item := range req.Items {
			if err := tx.Model(&models.Product{}).Where("id = ?", item.ProductID).
				Updates(map[string]interface{}{
					"stock": gorm.Expr("stock - ?", item.Quantity),
					"sales": gorm.Expr("sales + ?", item.Quantity),
				}).Error; err != nil {
				return err
			}
		}

		var cartProductIDs []uint
		for _, item := range req.Items {
			cartProductIDs = append(cartProductIDs, item.ProductID)
		}
		tx.Where("user_id = ? AND product_id IN ?", userID, cartProductIDs).Delete(&models.Cart{})

		return nil
	})

	if err != nil {
		utils.ServerError(c, "创建订单失败")
		return
	}

	utils.SuccessWithMessage(c, "订单创建成功", gin.H{
		"order_no":    orderNo,
		"total_price": totalPrice,
	})
}

func GetOrderList(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")

	offset := (page - 1) * pageSize

	var orders []models.Order
	var total int64

	query := database.DB.Model(&models.Order{}).Where("user_id = ?", userID)
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)
	query.Preload("User").Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&orders)

	for i := range orders {
		var orderItems []models.OrderItem
		database.DB.Where("order_id = ?", orders[i].ID).Preload("Product").Find(&orderItems)
		orders[i].UpdatedAt = orders[i].CreatedAt
	}

	utils.SuccessPage(c, orders, total, page, pageSize)
}

func GetOrderDetail(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, _ := strconv.Atoi(c.Param("id"))

	var order models.Order
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).Preload("User").First(&order).Error; err != nil {
		utils.NotFound(c, "订单不存在")
		return
	}

	var orderItems []models.OrderItem
	database.DB.Where("order_id = ?", order.ID).Preload("Product").Find(&orderItems)

	utils.Success(c, gin.H{
		"order": order,
		"items": orderItems,
	})
}

func PayOrder(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, _ := strconv.Atoi(c.Param("id"))
	var req struct {
		PaymentType string `json:"payment_type"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var order models.Order
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&order).Error; err != nil {
		utils.NotFound(c, "订单不存在")
		return
	}

	if order.Status != 0 {
		utils.BadRequest(c, "订单状态错误")
		return
	}

	now := time.Now()
	order.Status = 1
	order.PayTime = &now

	if err := database.DB.Save(&order).Error; err != nil {
		utils.ServerError(c, "支付失败")
		return
	}

	utils.SuccessWithMessage(c, "支付成功", order)
}

func CancelOrder(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, _ := strconv.Atoi(c.Param("id"))

	var order models.Order
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&order).Error; err != nil {
		utils.NotFound(c, "订单不存在")
		return
	}

	if order.Status != 0 {
		utils.BadRequest(c, "只能取消待付款订单")
		return
	}

	order.Status = -1

	if err := database.DB.Save(&order).Error; err != nil {
		utils.ServerError(c, "取消失败")
		return
	}

	utils.SuccessWithMessage(c, "取消成功", order)
}

func ConfirmReceive(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, _ := strconv.Atoi(c.Param("id"))

	var order models.Order
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&order).Error; err != nil {
		utils.NotFound(c, "订单不存在")
		return
	}

	if order.Status != 2 {
		utils.BadRequest(c, "只能确认已发货的订单")
		return
	}

	now := time.Now()
	order.Status = 3
	order.ReceiveTime = &now

	if err := database.DB.Save(&order).Error; err != nil {
		utils.ServerError(c, "确认失败")
		return
	}

	utils.SuccessWithMessage(c, "确认收货成功", order)
}

func GetAdminOrderList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")
	keyword := c.Query("keyword")

	offset := (page - 1) * pageSize

	var orders []models.Order
	var total int64

	query := database.DB.Model(&models.Order{}).Preload("User")
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if keyword != "" {
		query = query.Where("order_no LIKE ? OR receiver_name LIKE ? OR receiver_phone LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	query.Count(&total)
	query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&orders)

	utils.SuccessPage(c, orders, total, page, pageSize)
}

func GetAdminOrderDetail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var order models.Order
	if err := database.DB.Preload("User").First(&order, id).Error; err != nil {
		utils.NotFound(c, "订单不存在")
		return
	}

	var orderItems []models.OrderItem
	database.DB.Where("order_id = ?", order.ID).Preload("Product").Find(&orderItems)

	utils.Success(c, gin.H{
		"order": order,
		"items": orderItems,
	})
}

func ShipOrder(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req struct {
		ExpressCompany string `json:"express_company" binding:"required"`
		ExpressNo      string `json:"express_no" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var order models.Order
	if err := database.DB.First(&order, id).Error; err != nil {
		utils.NotFound(c, "订单不存在")
		return
	}

	if order.Status != 1 {
		utils.BadRequest(c, "只能发货待发货的订单")
		return
	}

	now := time.Now()
	order.Status = 2
	order.ShipTime = &now
	order.ExpressCompany = req.ExpressCompany
	order.ExpressNo = req.ExpressNo

	if err := database.DB.Save(&order).Error; err != nil {
		utils.ServerError(c, "发货失败")
		return
	}

	utils.SuccessWithMessage(c, "发货成功", order)
}

func UpdateLogistics(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req struct {
		ExpressCompany string `json:"express_company"`
		ExpressNo      string `json:"express_no"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var order models.Order
	if err := database.DB.First(&order, id).Error; err != nil {
		utils.NotFound(c, "订单不存在")
		return
	}

	updates := make(map[string]interface{})
	if req.ExpressCompany != "" {
		updates["express_company"] = req.ExpressCompany
	}
	if req.ExpressNo != "" {
		updates["express_no"] = req.ExpressNo
	}

	if len(updates) > 0 {
		database.DB.Model(&order).Updates(updates)
	}

	database.DB.First(&order, id)
	utils.SuccessWithMessage(c, "更新成功", order)
}
