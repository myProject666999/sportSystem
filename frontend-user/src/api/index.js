import request from '@/utils/request'

export function userLogin(data) {
  return request.post('/user/login', data)
}

export function userRegister(data) {
  return request.post('/user/register', data)
}

export function getUserInfo() {
  return request.get('/user/info')
}

export function updateUserInfo(data) {
  return request.put('/user/info', data)
}

export function updatePassword(data) {
  return request.put('/user/password', data)
}

export function recharge(data) {
  return request.post('/user/recharge', data)
}

export function getCategories() {
  return request.get('/categories')
}

export function getProducts(params) {
  return request.get('/products', { params })
}

export function getProductDetail(id) {
  return request.get(`/products/${id}`)
}

export function getBanners() {
  return request.get('/banners')
}

export function getNotices(params) {
  return request.get('/notices', { params })
}

export function getNoticeDetail(id) {
  return request.get(`/notices/${id}`)
}

export function getNews(params) {
  return request.get('/news', { params })
}

export function getNewsDetail(id) {
  return request.get(`/news/${id}`)
}

export function getVideos(params) {
  return request.get('/videos', { params })
}

export function getVideoDetail(id) {
  return request.get(`/videos/${id}`)
}

export function getAds(params) {
  return request.get('/ads', { params })
}

export function getCartList() {
  return request.get('/user/cart')
}

export function addToCart(data) {
  return request.post('/user/cart', data)
}

export function updateCart(id, data) {
  return request.put(`/user/cart/${id}`, data)
}

export function removeFromCart(id) {
  return request.delete(`/user/cart/${id}`)
}

export function getFavorites(params) {
  return request.get('/user/favorites', { params })
}

export function addFavorite(data) {
  return request.post('/user/favorites', data)
}

export function removeFavorite(id) {
  return request.delete(`/user/favorites/${id}`)
}

export function checkFavorite(productId) {
  return request.get(`/user/favorites/check/${productId}`)
}

export function getAddresses() {
  return request.get('/user/addresses')
}

export function getDefaultAddress() {
  return request.get('/user/addresses/default')
}

export function createAddress(data) {
  return request.post('/user/addresses', data)
}

export function updateAddress(id, data) {
  return request.put(`/user/addresses/${id}`, data)
}

export function deleteAddress(id) {
  return request.delete(`/user/addresses/${id}`)
}

export function createOrder(data) {
  return request.post('/user/orders', data)
}

export function getOrders(params) {
  return request.get('/user/orders', { params })
}

export function getOrderDetail(id) {
  return request.get(`/user/orders/${id}`)
}

export function payOrder(id) {
  return request.post(`/user/orders/${id}/pay`)
}

export function cancelOrder(id) {
  return request.post(`/user/orders/${id}/cancel`)
}

export function confirmReceive(id) {
  return request.post(`/user/orders/${id}/receive`)
}

export function addComment(data) {
  return request.post('/user/comments', data)
}
