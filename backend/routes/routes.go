package routes

import (
	"sportSystem/controllers"
	"sportSystem/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(middlewares.Cors())

	api := r.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})

		api.POST("/user/register", controllers.UserRegister)
		api.POST("/user/login", controllers.UserLogin)

		api.GET("/categories", controllers.GetCategoryList)
		api.GET("/products", controllers.GetProductList)
		api.GET("/products/:id", controllers.GetProductDetail)
		api.GET("/products/:id/comments", controllers.GetProductComments)

		api.GET("/banners", controllers.GetBannerList)
		api.GET("/notices", controllers.GetNoticeList)
		api.GET("/notices/:id", controllers.GetNoticeDetail)
		api.GET("/news", controllers.GetNewsList)
		api.GET("/news/:id", controllers.GetNewsDetail)
		api.GET("/videos", controllers.GetVideoList)
		api.GET("/videos/:id", controllers.GetVideoDetail)
		api.GET("/ads", controllers.GetAdList)

		user := api.Group("/user")
		user.Use(middlewares.UserAuth())
		{
			user.GET("/info", controllers.GetUserInfo)
			user.PUT("/info", controllers.UpdateUserInfo)
			user.PUT("/password", controllers.UpdateUserPassword)
			user.POST("/recharge", controllers.Recharge)

			user.GET("/cart", controllers.GetCartList)
			user.POST("/cart", controllers.AddToCart)
			user.PUT("/cart/:id", controllers.UpdateCart)
			user.DELETE("/cart/:id", controllers.RemoveFromCart)

			user.GET("/favorites", controllers.GetFavorites)
			user.POST("/favorites", controllers.AddFavorite)
			user.DELETE("/favorites/:id", controllers.RemoveFavorite)
			user.GET("/favorites/check/:product_id", controllers.CheckFavorite)

			user.GET("/addresses", controllers.GetAddressList)
			user.GET("/addresses/default", controllers.GetDefaultAddress)
			user.POST("/addresses", controllers.CreateAddress)
			user.PUT("/addresses/:id", controllers.UpdateAddress)
			user.DELETE("/addresses/:id", controllers.DeleteAddress)

			user.POST("/orders", controllers.CreateOrder)
			user.GET("/orders", controllers.GetOrderList)
			user.GET("/orders/:id", controllers.GetOrderDetail)
			user.POST("/orders/:id/pay", controllers.PayOrder)
			user.POST("/orders/:id/cancel", controllers.CancelOrder)
			user.POST("/orders/:id/receive", controllers.ConfirmReceive)

			user.POST("/comments", controllers.AddProductComment)
		}

		admin := api.Group("/admin")
		{
			admin.POST("/login", controllers.AdminLogin)

			adminAuth := admin.Group("")
			adminAuth.Use(middlewares.AdminAuth())
			{
				adminAuth.GET("/info", controllers.GetAdminInfo)
				adminAuth.GET("/dashboard", controllers.GetDashboardStats)

				adminAuth.GET("/users", controllers.GetUserList)

				adminAuth.GET("/admins", controllers.GetAdminList)
				adminAuth.POST("/admins", controllers.CreateAdmin)
				adminAuth.PUT("/admins/:id", controllers.UpdateAdmin)
				adminAuth.DELETE("/admins/:id", controllers.DeleteAdmin)

				adminAuth.GET("/categories/all", controllers.GetAllCategories)
				adminAuth.POST("/categories", controllers.CreateCategory)
				adminAuth.PUT("/categories/:id", controllers.UpdateCategory)
				adminAuth.DELETE("/categories/:id", controllers.DeleteCategory)

				adminAuth.GET("/products/all", controllers.GetProductList)
				adminAuth.POST("/products", controllers.CreateProduct)
				adminAuth.PUT("/products/:id", controllers.UpdateProduct)
				adminAuth.DELETE("/products/:id", controllers.DeleteProduct)

				adminAuth.GET("/orders", controllers.GetAdminOrderList)
				adminAuth.GET("/orders/:id", controllers.GetAdminOrderDetail)
				adminAuth.POST("/orders/:id/ship", controllers.ShipOrder)
				adminAuth.PUT("/orders/:id/logistics", controllers.UpdateLogistics)

				adminAuth.GET("/news/all", controllers.GetAllNews)
				adminAuth.POST("/news", controllers.CreateNews)
				adminAuth.PUT("/news/:id", controllers.UpdateNews)
				adminAuth.DELETE("/news/:id", controllers.DeleteNews)

				adminAuth.GET("/videos/all", controllers.GetAllVideos)
				adminAuth.POST("/videos", controllers.CreateVideo)
				adminAuth.PUT("/videos/:id", controllers.UpdateVideo)
				adminAuth.DELETE("/videos/:id", controllers.DeleteVideo)

				adminAuth.GET("/banners/all", controllers.GetAllBanners)
				adminAuth.POST("/banners", controllers.CreateBanner)
				adminAuth.PUT("/banners/:id", controllers.UpdateBanner)
				adminAuth.DELETE("/banners/:id", controllers.DeleteBanner)

				adminAuth.GET("/notices/all", controllers.GetAllNotices)
				adminAuth.POST("/notices", controllers.CreateNotice)
				adminAuth.PUT("/notices/:id", controllers.UpdateNotice)
				adminAuth.DELETE("/notices/:id", controllers.DeleteNotice)

				adminAuth.GET("/ads/all", controllers.GetAllAds)
				adminAuth.POST("/ads", controllers.CreateAd)
				adminAuth.PUT("/ads/:id", controllers.UpdateAd)
				adminAuth.DELETE("/ads/:id", controllers.DeleteAd)
			}
		}
	}

	return r
}
