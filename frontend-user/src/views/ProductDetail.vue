<template>
  <div class="product-detail-page">
    <div class="container">
      <el-breadcrumb separator="/" class="breadcrumb">
        <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
        <el-breadcrumb-item :to="{ path: '/products' }">商品列表</el-breadcrumb-item>
        <el-breadcrumb-item>{{ product.name }}</el-breadcrumb-item>
      </el-breadcrumb>

      <div class="product-card" v-loading="loading">
        <el-row :gutter="40">
          <el-col :span="10">
            <div class="product-image">
              <el-image :src="product.image" fit="cover" style="width: 100%; height: 400px;" />
            </div>
          </el-col>
          <el-col :span="14">
            <div class="product-info">
              <h1 class="product-name">{{ product.name }}</h1>
              <div class="product-price">
                <span class="price-label">售价</span>
                <span class="price">¥{{ product.price }}</span>
                <el-tag v-if="product.points_price" type="success" size="large" style="margin-left: 15px;">
                  可用 {{ product.points_price }} 积分兑换
                </el-tag>
              </div>
              <div class="product-meta">
                <div class="meta-item">
                  <span class="meta-label">库存</span>
                  <span class="meta-value">{{ product.stock }} 件</span>
                </div>
                <div class="meta-item">
                  <span class="meta-label">销量</span>
                  <span class="meta-value">{{ product.sales }} 件</span>
                </div>
                <div class="meta-item">
                  <span class="meta-label">分类</span>
                  <span class="meta-value">{{ product.category_name }}</span>
                </div>
              </div>
              <div class="quantity-section">
                <span class="label">数量</span>
                <el-input-number v-model="quantity" :min="1" :max="product.stock" />
              </div>
              <div class="action-section">
                <el-button type="danger" size="large" @click="buyNow">
                  <el-icon><ShoppingCart /></el-icon>
                  立即购买
                </el-button>
                <el-button type="primary" size="large" @click="handleAddToCart">
                  <el-icon><Plus /></el-icon>
                  加入购物车
                </el-button>
                <el-button :type="isFavorited ? 'success' : ''" size="large" @click="toggleFavorite">
                  <el-icon><Star /></el-icon>
                  {{ isFavorited ? '已收藏' : '收藏' }}
                </el-button>
              </div>
            </div>
          </el-col>
        </el-row>
      </div>

      <el-card class="detail-card">
        <template #header>
          <div class="tab-header">
            <span :class="{ active: activeTab === 'detail' }" @click="activeTab = 'detail'">商品详情</span>
            <span :class="{ active: activeTab === 'comments' }" @click="activeTab = 'comments'">用户评论 ({{ comments.length }})</span>
          </div>
        </template>

        <div v-if="activeTab === 'detail'" class="content-section">
          <div class="description" v-html="product.description"></div>
        </div>

        <div v-else class="comments-section">
          <div class="comment-form">
            <el-input
              v-model="commentContent"
              type="textarea"
              :rows="3"
              placeholder="发表您的评论..."
            />
            <div class="comment-actions">
              <el-button type="primary" @click="submitComment" :disabled="!commentContent.trim()">
                发表评论
              </el-button>
            </div>
          </div>

          <div class="comment-list">
            <div v-for="comment in comments" :key="comment.id" class="comment-item">
              <div class="comment-header">
                <el-avatar :size="36">
                  <el-icon><User /></el-icon>
                </el-avatar>
                <div class="comment-info">
                  <span class="comment-user">{{ comment.username }}</span>
                  <span class="comment-time">{{ comment.created_at }}</span>
                </div>
              </div>
              <div class="comment-content">{{ comment.content }}</div>
            </div>
            <el-empty v-if="comments.length === 0" description="暂无评论" />
          </div>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { ShoppingCart, Plus, Star, User } from '@element-plus/icons-vue'
import { getProductDetail, addToCart, addFavorite, removeFavorite } from '@/api'

const route = useRoute()
const router = useRouter()
const loading = ref(false)
const product = ref({})
const quantity = ref(1)
const activeTab = ref('detail')
const commentContent = ref('')
const comments = ref([])
const isFavorited = ref(false)

const loadProduct = async () => {
  loading.value = true
  try {
    const res = await getProductDetail(route.params.id)
    product.value = res.data || {}
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const buyNow = () => {
  const token = localStorage.getItem('token')
  if (!token) {
    ElMessage.warning('请先登录')
    router.push(`/login?redirect=${route.fullPath}`)
    return
  }
  ElMessage.success('跳转到确认订单页')
}

const handleAddToCart = async () => {
  const token = localStorage.getItem('token')
  if (!token) {
    ElMessage.warning('请先登录')
    router.push(`/login?redirect=${route.fullPath}`)
    return
  }
  try {
    await addToCart({ product_id: product.value.id, quantity: quantity.value })
    ElMessage.success('已加入购物车')
  } catch (e) {
    console.error(e)
  }
}

const toggleFavorite = async () => {
  const token = localStorage.getItem('token')
  if (!token) {
    ElMessage.warning('请先登录')
    router.push(`/login?redirect=${route.fullPath}`)
    return
  }
  try {
    if (isFavorited.value) {
      await removeFavorite(product.value.id)
      isFavorited.value = false
      ElMessage.success('已取消收藏')
    } else {
      await addFavorite({ product_id: product.value.id })
      isFavorited.value = true
      ElMessage.success('已收藏')
    }
  } catch (e) {
    console.error(e)
  }
}

const submitComment = async () => {
  const token = localStorage.getItem('token')
  if (!token) {
    ElMessage.warning('请先登录')
    router.push(`/login?redirect=${route.fullPath}`)
    return
  }
  if (!commentContent.value.trim()) return
  ElMessage.success('评论发表成功')
  comments.value.unshift({
    id: Date.now(),
    username: '当前用户',
    content: commentContent.value,
    created_at: new Date().toLocaleString()
  })
  commentContent.value = ''
}

onMounted(() => {
  loadProduct()
})
</script>

<style scoped>
.product-detail-page {
  padding: 20px 0;
  background: #f5f5f5;
  min-height: calc(100vh - 140px);
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 15px;
}

.breadcrumb {
  margin-bottom: 20px;
}

.product-card {
  background: #fff;
  border-radius: 8px;
  padding: 30px;
  margin-bottom: 20px;
}

.product-image {
  border: 1px solid #eee;
  border-radius: 8px;
  overflow: hidden;
}

.product-name {
  font-size: 24px;
  color: #333;
  margin: 0 0 20px 0;
  line-height: 1.4;
}

.product-price {
  display: flex;
  align-items: center;
  padding: 15px 20px;
  background: #fff5f5;
  border-radius: 4px;
  margin-bottom: 20px;
}

.price-label {
  color: #999;
  margin-right: 10px;
}

.price {
  font-size: 28px;
  color: #e4393c;
  font-weight: bold;
}

.product-meta {
  display: flex;
  gap: 40px;
  margin-bottom: 25px;
  padding: 0 10px;
}

.meta-item {
  display: flex;
  align-items: center;
}

.meta-label {
  color: #999;
  margin-right: 10px;
}

.meta-value {
  color: #333;
}

.quantity-section {
  display: flex;
  align-items: center;
  margin-bottom: 30px;
}

.quantity-section .label {
  color: #999;
  margin-right: 15px;
}

.action-section {
  display: flex;
  gap: 15px;
}

.action-section .el-button {
  height: 46px;
  font-size: 16px;
  padding: 0 30px;
}

.detail-card {
  border-radius: 8px;
}

.tab-header {
  display: flex;
  gap: 30px;
}

.tab-header span {
  font-size: 16px;
  cursor: pointer;
  color: #666;
  padding-bottom: 5px;
}

.tab-header span.active {
  color: #409EFF;
  font-weight: bold;
  border-bottom: 2px solid #409EFF;
}

.content-section {
  padding: 20px 0;
}

.description {
  line-height: 1.8;
  color: #666;
}

.comments-section {
  padding: 20px 0;
}

.comment-form {
  margin-bottom: 30px;
  padding-bottom: 20px;
  border-bottom: 1px solid #eee;
}

.comment-actions {
  margin-top: 15px;
  text-align: right;
}

.comment-item {
  padding: 20px 0;
  border-bottom: 1px solid #f0f0f0;
}

.comment-header {
  display: flex;
  align-items: center;
  margin-bottom: 12px;
}

.comment-info {
  margin-left: 12px;
}

.comment-user {
  display: block;
  color: #333;
  font-weight: 500;
}

.comment-time {
  font-size: 12px;
  color: #999;
}

.comment-content {
  color: #666;
  line-height: 1.6;
  padding-left: 48px;
}
</style>
