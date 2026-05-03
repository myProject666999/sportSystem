package controllers

import (
	"sportSystem/database"
	"sportSystem/models"
	"sportSystem/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetNewsList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	keyword := c.Query("keyword")

	offset := (page - 1) * pageSize

	var news []models.News
	var total int64

	query := database.DB.Model(&models.News{}).Where("status = 1")
	if keyword != "" {
		query = query.Where("title LIKE ?", "%"+keyword+"%")
	}

	query.Count(&total)
	query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&news)

	utils.SuccessPage(c, news, total, page, pageSize)
}

func GetNewsDetail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var news models.News
	if err := database.DB.First(&news, id).Error; err != nil {
		utils.NotFound(c, "新闻不存在")
		return
	}

	database.DB.Model(&models.News{}).Where("id = ?", id).Update("views", database.DB.Raw("views + 1"))

	utils.Success(c, news)
}

func GetAllNews(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	keyword := c.Query("keyword")
	status := c.Query("status")

	offset := (page - 1) * pageSize

	var news []models.News
	var total int64

	query := database.DB.Model(&models.News{})
	if keyword != "" {
		query = query.Where("title LIKE ?", "%"+keyword+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)
	query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&news)

	utils.SuccessPage(c, news, total, page, pageSize)
}

func CreateNews(c *gin.Context) {
	var req struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
		Cover   string `json:"cover"`
		Status  int    `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	news := &models.News{
		Title:   req.Title,
		Content: req.Content,
		Cover:   req.Cover,
		Status:  req.Status,
	}

	if err := database.DB.Create(news).Error; err != nil {
		utils.ServerError(c, "创建失败")
		return
	}

	utils.SuccessWithMessage(c, "创建成功", news)
}

func UpdateNews(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
		Cover   string `json:"cover"`
		Status  int    `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Content != "" {
		updates["content"] = req.Content
	}
	if req.Cover != "" {
		updates["cover"] = req.Cover
	}
	updates["status"] = req.Status

	if len(updates) > 0 {
		if err := database.DB.Model(&models.News{}).Where("id = ?", id).Updates(updates).Error; err != nil {
			utils.ServerError(c, "更新失败")
			return
		}
	}

	var news models.News
	database.DB.First(&news, id)
	utils.SuccessWithMessage(c, "更新成功", news)
}

func DeleteNews(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := database.DB.Delete(&models.News{}, id).Error; err != nil {
		utils.ServerError(c, "删除失败")
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}

func GetVideoList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	keyword := c.Query("keyword")

	offset := (page - 1) * pageSize

	var videos []models.Video
	var total int64

	query := database.DB.Model(&models.Video{}).Where("status = 1")
	if keyword != "" {
		query = query.Where("title LIKE ?", "%"+keyword+"%")
	}

	query.Count(&total)
	query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&videos)

	utils.SuccessPage(c, videos, total, page, pageSize)
}

func GetVideoDetail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var video models.Video
	if err := database.DB.First(&video, id).Error; err != nil {
		utils.NotFound(c, "视频不存在")
		return
	}

	database.DB.Model(&models.Video{}).Where("id = ?", id).Update("views", database.DB.Raw("views + 1"))

	utils.Success(c, video)
}

func GetAllVideos(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	keyword := c.Query("keyword")
	status := c.Query("status")

	offset := (page - 1) * pageSize

	var videos []models.Video
	var total int64

	query := database.DB.Model(&models.Video{})
	if keyword != "" {
		query = query.Where("title LIKE ?", "%"+keyword+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)
	query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&videos)

	utils.SuccessPage(c, videos, total, page, pageSize)
}

func CreateVideo(c *gin.Context) {
	var req struct {
		Title  string `json:"title" binding:"required"`
		URL    string `json:"url" binding:"required"`
		Cover  string `json:"cover"`
		Status int    `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	video := &models.Video{
		Title:  req.Title,
		URL:    req.URL,
		Cover:  req.Cover,
		Status: req.Status,
	}

	if err := database.DB.Create(video).Error; err != nil {
		utils.ServerError(c, "创建失败")
		return
	}

	utils.SuccessWithMessage(c, "创建成功", video)
}

func UpdateVideo(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req struct {
		Title  string `json:"title"`
		URL    string `json:"url"`
		Cover  string `json:"cover"`
		Status int    `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.URL != "" {
		updates["url"] = req.URL
	}
	if req.Cover != "" {
		updates["cover"] = req.Cover
	}
	updates["status"] = req.Status

	if len(updates) > 0 {
		if err := database.DB.Model(&models.Video{}).Where("id = ?", id).Updates(updates).Error; err != nil {
			utils.ServerError(c, "更新失败")
			return
		}
	}

	var video models.Video
	database.DB.First(&video, id)
	utils.SuccessWithMessage(c, "更新成功", video)
}

func DeleteVideo(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := database.DB.Delete(&models.Video{}, id).Error; err != nil {
		utils.ServerError(c, "删除失败")
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}

func GetBannerList(c *gin.Context) {
	var banners []models.Banner
	database.DB.Where("status = 1").Order("sort ASC, created_at DESC").Find(&banners)
	utils.Success(c, banners)
}

func GetAllBanners(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	offset := (page - 1) * pageSize

	var banners []models.Banner
	var total int64

	database.DB.Model(&models.Banner{}).Count(&total)
	database.DB.Offset(offset).Limit(pageSize).Order("sort ASC, created_at DESC").Find(&banners)

	utils.SuccessPage(c, banners, total, page, pageSize)
}

func CreateBanner(c *gin.Context) {
	var req struct {
		Title  string `json:"title"`
		Image  string `json:"image" binding:"required"`
		Link   string `json:"link"`
		Sort   int    `json:"sort"`
		Status int    `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	banner := &models.Banner{
		Title:  req.Title,
		Image:  req.Image,
		Link:   req.Link,
		Sort:   req.Sort,
		Status: req.Status,
	}

	if err := database.DB.Create(banner).Error; err != nil {
		utils.ServerError(c, "创建失败")
		return
	}

	utils.SuccessWithMessage(c, "创建成功", banner)
}

func UpdateBanner(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req struct {
		Title  string `json:"title"`
		Image  string `json:"image"`
		Link   string `json:"link"`
		Sort   int    `json:"sort"`
		Status int    `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Image != "" {
		updates["image"] = req.Image
	}
	if req.Link != "" {
		updates["link"] = req.Link
	}
	if req.Sort != 0 {
		updates["sort"] = req.Sort
	}
	updates["status"] = req.Status

	if len(updates) > 0 {
		if err := database.DB.Model(&models.Banner{}).Where("id = ?", id).Updates(updates).Error; err != nil {
			utils.ServerError(c, "更新失败")
			return
		}
	}

	var banner models.Banner
	database.DB.First(&banner, id)
	utils.SuccessWithMessage(c, "更新成功", banner)
}

func DeleteBanner(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := database.DB.Delete(&models.Banner{}, id).Error; err != nil {
		utils.ServerError(c, "删除失败")
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}

func GetNoticeList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	offset := (page - 1) * pageSize

	var notices []models.Notice
	var total int64

	query := database.DB.Model(&models.Notice{}).Where("status = 1")
	query.Count(&total)
	query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&notices)

	utils.SuccessPage(c, notices, total, page, pageSize)
}

func GetNoticeDetail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var notice models.Notice
	if err := database.DB.First(&notice, id).Error; err != nil {
		utils.NotFound(c, "公告不存在")
		return
	}

	utils.Success(c, notice)
}

func GetAllNotices(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	offset := (page - 1) * pageSize

	var notices []models.Notice
	var total int64

	database.DB.Model(&models.Notice{}).Count(&total)
	database.DB.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&notices)

	utils.SuccessPage(c, notices, total, page, pageSize)
}

func CreateNotice(c *gin.Context) {
	var req struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
		Status  int    `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	notice := &models.Notice{
		Title:   req.Title,
		Content: req.Content,
		Status:  req.Status,
	}

	if err := database.DB.Create(notice).Error; err != nil {
		utils.ServerError(c, "创建失败")
		return
	}

	utils.SuccessWithMessage(c, "创建成功", notice)
}

func UpdateNotice(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
		Status  int    `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Content != "" {
		updates["content"] = req.Content
	}
	updates["status"] = req.Status

	if len(updates) > 0 {
		if err := database.DB.Model(&models.Notice{}).Where("id = ?", id).Updates(updates).Error; err != nil {
			utils.ServerError(c, "更新失败")
			return
		}
	}

	var notice models.Notice
	database.DB.First(&notice, id)
	utils.SuccessWithMessage(c, "更新成功", notice)
}

func DeleteNotice(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := database.DB.Delete(&models.Notice{}, id).Error; err != nil {
		utils.ServerError(c, "删除失败")
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}

func GetAdList(c *gin.Context) {
	position := c.Query("position")

	var ads []models.Advertisement
	query := database.DB.Where("status = 1")
	if position != "" {
		query = query.Where("position = ?", position)
	}
	query.Order("sort ASC, created_at DESC").Find(&ads)

	utils.Success(c, ads)
}

func GetAllAds(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	offset := (page - 1) * pageSize

	var ads []models.Advertisement
	var total int64

	database.DB.Model(&models.Advertisement{}).Count(&total)
	database.DB.Offset(offset).Limit(pageSize).Order("sort ASC, created_at DESC").Find(&ads)

	utils.SuccessPage(c, ads, total, page, pageSize)
}

func CreateAd(c *gin.Context) {
	var req struct {
		Title    string `json:"title"`
		Image    string `json:"image" binding:"required"`
		Link     string `json:"link"`
		Position string `json:"position"`
		Sort     int    `json:"sort"`
		Status   int    `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	ad := &models.Advertisement{
		Title:    req.Title,
		Image:    req.Image,
		Link:     req.Link,
		Position: req.Position,
		Sort:     req.Sort,
		Status:   req.Status,
	}

	if err := database.DB.Create(ad).Error; err != nil {
		utils.ServerError(c, "创建失败")
		return
	}

	utils.SuccessWithMessage(c, "创建成功", ad)
}

func UpdateAd(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req struct {
		Title    string `json:"title"`
		Image    string `json:"image"`
		Link     string `json:"link"`
		Position string `json:"position"`
		Sort     int    `json:"sort"`
		Status   int    `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Image != "" {
		updates["image"] = req.Image
	}
	if req.Link != "" {
		updates["link"] = req.Link
	}
	if req.Position != "" {
		updates["position"] = req.Position
	}
	if req.Sort != 0 {
		updates["sort"] = req.Sort
	}
	updates["status"] = req.Status

	if len(updates) > 0 {
		if err := database.DB.Model(&models.Advertisement{}).Where("id = ?", id).Updates(updates).Error; err != nil {
			utils.ServerError(c, "更新失败")
			return
		}
	}

	var ad models.Advertisement
	database.DB.First(&ad, id)
	utils.SuccessWithMessage(c, "更新成功", ad)
}

func DeleteAd(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := database.DB.Delete(&models.Advertisement{}, id).Error; err != nil {
		utils.ServerError(c, "删除失败")
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}
