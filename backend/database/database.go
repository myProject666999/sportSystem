package database

import (
	"log"
	"sportSystem/config"
	"sportSystem/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() error {
	dsn := config.GetDSN()

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return err
	}

	DB = db

	err = DB.AutoMigrate(
		&models.User{},
		&models.Admin{},
		&models.Category{},
		&models.Product{},
		&models.ProductComment{},
		&models.Favorite{},
		&models.Cart{},
		&models.Address{},
		&models.Order{},
		&models.OrderItem{},
		&models.Recharge{},
		&models.News{},
		&models.Video{},
		&models.Banner{},
		&models.Notice{},
		&models.Advertisement{},
	)
	if err != nil {
		return err
	}

	initData()

	log.Println("数据库初始化成功")
	return nil
}

func initData() {
	var adminCount int64
	DB.Model(&models.Admin{}).Count(&adminCount)
	if adminCount == 0 {
		admin := &models.Admin{
			Username: "admin",
			Role:     "super",
		}
		admin.SetPassword("admin123")
		DB.Create(admin)
		log.Println("默认管理员账号创建成功: admin / admin123")
	}

	var categoryCount int64
	DB.Model(&models.Category{}).Count(&categoryCount)
	if categoryCount == 0 {
		categories := []models.Category{
			{Name: "运动器材", Sort: 1},
			{Name: "运动服饰", Sort: 2},
			{Name: "运动鞋", Sort: 3},
			{Name: "配件", Sort: 4},
		}
		DB.Create(&categories)
		log.Println("默认分类创建成功")
	}
}
