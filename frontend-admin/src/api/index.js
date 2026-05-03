import request from '@/utils/request'

export function adminLogin(data) {
  return request.post('/admin/login', data)
}

export function getAdminInfo() {
  return request.get('/admin/info')
}

export function getDashboardStats() {
  return request.get('/admin/dashboard')
}

export function getUsers(params) {
  return request.get('/admin/users', { params })
}

export function getAdmins(params) {
  return request.get('/admin/admins', { params })
}

export function createAdmin(data) {
  return request.post('/admin/admins', data)
}

export function updateAdmin(id, data) {
  return request.put(`/admin/admins/${id}`, data)
}

export function deleteAdmin(id) {
  return request.delete(`/admin/admins/${id}`)
}

export function getCategories() {
  return request.get('/admin/categories/all')
}

export function createCategory(data) {
  return request.post('/admin/categories', data)
}

export function updateCategory(id, data) {
  return request.put(`/admin/categories/${id}`, data)
}

export function deleteCategory(id) {
  return request.delete(`/admin/categories/${id}`)
}

export function getProducts(params) {
  return request.get('/admin/products/all', { params })
}

export function createProduct(data) {
  return request.post('/admin/products', data)
}

export function updateProduct(id, data) {
  return request.put(`/admin/products/${id}`, data)
}

export function deleteProduct(id) {
  return request.delete(`/admin/products/${id}`)
}

export function getOrders(params) {
  return request.get('/admin/orders', { params })
}

export function getOrderDetail(id) {
  return request.get(`/admin/orders/${id}`)
}

export function shipOrder(id, data) {
  return request.post(`/admin/orders/${id}/ship`, data)
}

export function updateLogistics(id, data) {
  return request.put(`/admin/orders/${id}/logistics`, data)
}

export function getNews(params) {
  return request.get('/admin/news/all', { params })
}

export function createNews(data) {
  return request.post('/admin/news', data)
}

export function updateNews(id, data) {
  return request.put(`/admin/news/${id}`, data)
}

export function deleteNews(id) {
  return request.delete(`/admin/news/${id}`)
}

export function getVideos(params) {
  return request.get('/admin/videos/all', { params })
}

export function createVideo(data) {
  return request.post('/admin/videos', data)
}

export function updateVideo(id, data) {
  return request.put(`/admin/videos/${id}`, data)
}

export function deleteVideo(id) {
  return request.delete(`/admin/videos/${id}`)
}

export function getBanners(params) {
  return request.get('/admin/banners/all', { params })
}

export function createBanner(data) {
  return request.post('/admin/banners', data)
}

export function updateBanner(id, data) {
  return request.put(`/admin/banners/${id}`, data)
}

export function deleteBanner(id) {
  return request.delete(`/admin/banners/${id}`)
}

export function getNotices(params) {
  return request.get('/admin/notices/all', { params })
}

export function createNotice(data) {
  return request.post('/admin/notices', data)
}

export function updateNotice(id, data) {
  return request.put(`/admin/notices/${id}`, data)
}

export function deleteNotice(id) {
  return request.delete(`/admin/notices/${id}`)
}

export function getAds(params) {
  return request.get('/admin/ads/all', { params })
}

export function createAd(data) {
  return request.post('/admin/ads', data)
}

export function updateAd(id, data) {
  return request.put(`/admin/ads/${id}`, data)
}

export function deleteAd(id) {
  return request.delete(`/admin/ads/${id}`)
}
