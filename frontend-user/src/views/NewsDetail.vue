<template>
  <div class="news-detail-page">
    <div class="container">
      <div class="breadcrumb">
        <span @click="$router.push('/')">首页</span>
        <span class="separator">/</span>
        <span @click="$router.push('/news')">新闻资讯</span>
        <span class="separator">/</span>
        <span class="active">{{ news.title }}</span>
      </div>

      <div class="news-content" v-loading="loading">
        <div class="news-main">
          <h1 class="news-title">{{ news.title }}</h1>
          <div class="news-meta">
            <span class="meta-item">
              <el-icon><Clock /></el-icon>
              {{ news.created_at }}
            </span>
            <span class="meta-item">
              <el-icon><View /></el-icon>
              {{ news.views || 0 }} 阅读
            </span>
          </div>
          <div class="news-body" v-html="news.content"></div>

          <div class="news-navigation">
            <el-button type="primary" link @click="$router.push('/news')">
              <el-icon><ArrowLeft /></el-icon>
              返回列表
            </el-button>
          </div>
        </div>

        <div class="news-sidebar">
          <el-card class="sidebar-card" shadow="never">
            <template #header>
              <span class="sidebar-title">热门推荐</span>
            </template>
            <div class="hot-list">
              <div
                v-for="(item, index) in hotList"
                :key="item.id"
                class="hot-item"
                @click="$router.push(`/news/${item.id}`)"
              >
                <span class="hot-index" :class="{ 'top-three': index < 3 }">{{ index + 1 }}</span>
                <span class="hot-title">{{ item.title }}</span>
              </div>
            </div>
          </el-card>

          <el-card class="sidebar-card" shadow="never" style="margin-top: 20px;">
            <template #header>
              <span class="sidebar-title">猜你喜欢</span>
            </template>
            <div class="product-list">
              <div
                v-for="product in recommendProducts"
                :key="product.id"
                class="product-item"
                @click="$router.push(`/product/${product.id}`)"
              >
                <el-image :src="product.image" fit="cover" class="product-img" />
                <div class="product-info">
                  <div class="product-name">{{ product.name }}</div>
                  <div class="product-price">¥{{ product.price }}</div>
                </div>
              </div>
            </div>
          </el-card>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Clock, View, ArrowLeft } from '@element-plus/icons-vue'
import { getNewsDetail, getNews, getProducts } from '@/api'

const route = useRoute()
const router = useRouter()
const loading = ref(false)
const news = ref({})
const hotList = ref([])
const recommendProducts = ref([])

const loadNewsDetail = async () => {
  loading.value = true
  try {
    const res = await getNewsDetail(route.params.id)
    news.value = res.data || {}
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const loadHotNews = async () => {
  try {
    const res = await getNews({ page_size: 8 })
    hotList.value = res.data?.list || []
  } catch (e) {
    console.error(e)
  }
}

const loadRecommendProducts = async () => {
  try {
    const res = await getProducts({ page_size: 6, status: 1 })
    recommendProducts.value = res.data?.list || []
  } catch (e) {
    console.error(e)
  }
}

onMounted(() => {
  loadNewsDetail()
  loadHotNews()
  loadRecommendProducts()
})
</script>

<style scoped>
.news-detail-page {
  padding: 30px 0;
  background: #f5f5f5;
  min-height: calc(100vh - 140px);
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 15px;
}

.breadcrumb {
  color: #666;
  font-size: 14px;
  margin-bottom: 20px;
}

.breadcrumb span {
  cursor: pointer;
}

.breadcrumb span:hover {
  color: #409EFF;
}

.breadcrumb .separator {
  margin: 0 8px;
  cursor: default;
}

.breadcrumb .active {
  color: #333;
  cursor: default;
}

.breadcrumb .active:hover {
  color: #333;
}

.news-content {
  display: flex;
  gap: 20px;
}

.news-main {
  flex: 1;
  background: #fff;
  border-radius: 8px;
  padding: 30px;
}

.news-title {
  font-size: 24px;
  color: #333;
  font-weight: bold;
  line-height: 1.4;
  margin-bottom: 20px;
}

.news-meta {
  display: flex;
  gap: 25px;
  padding-bottom: 20px;
  border-bottom: 1px solid #eee;
  margin-bottom: 25px;
  color: #999;
  font-size: 14px;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 5px;
}

.news-body {
  font-size: 16px;
  line-height: 1.8;
  color: #333;
}

.news-body :deep(img) {
  max-width: 100%;
  height: auto;
  border-radius: 8px;
  margin: 15px 0;
}

.news-body :deep(p) {
  margin-bottom: 15px;
  text-indent: 2em;
}

.news-body :deep(h1),
.news-body :deep(h2),
.news-body :deep(h3) {
  margin: 25px 0 15px;
  font-weight: bold;
}

.news-navigation {
  margin-top: 40px;
  padding-top: 20px;
  border-top: 1px solid #eee;
}

.news-sidebar {
  width: 300px;
  flex-shrink: 0;
}

.sidebar-card {
  border-radius: 8px;
}

.sidebar-title {
  font-weight: bold;
  font-size: 16px;
}

.hot-list {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.hot-item {
  display: flex;
  align-items: center;
  gap: 12px;
  cursor: pointer;
  padding: 5px 0;
  transition: all 0.3s;
}

.hot-item:hover {
  color: #409EFF;
}

.hot-index {
  width: 20px;
  height: 20px;
  background: #e4e7ed;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  color: #909399;
  flex-shrink: 0;
}

.hot-index.top-three {
  background: #e4393c;
  color: #fff;
}

.hot-title {
  flex: 1;
  font-size: 14px;
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.product-list {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.product-item {
  display: flex;
  gap: 12px;
  cursor: pointer;
  padding: 5px 0;
  transition: all 0.3s;
}

.product-item:hover .product-name {
  color: #409EFF;
}

.product-img {
  width: 70px;
  height: 70px;
  border-radius: 4px;
  flex-shrink: 0;
}

.product-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.product-name {
  font-size: 14px;
  color: #333;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  margin-bottom: 8px;
  transition: color 0.3s;
}

.product-price {
  color: #e4393c;
  font-weight: bold;
  font-size: 16px;
}
</style>
