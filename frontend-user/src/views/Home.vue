<template>
  <Layout>
    <div class="home-page">
      <el-carousel height="400px" class="banner-carousel" indicator-position="outside">
        <el-carousel-item v-for="banner in banners" :key="banner.id">
          <div class="banner-item" @click="goToLink(banner.link)">
            <img :src="banner.image" :alt="banner.title" class="banner-image" />
            <div class="banner-overlay">
              <h3>{{ banner.title }}</h3>
            </div>
          </div>
        </el-carousel-item>
      </el-carousel>

      <div class="container">
        <div class="section notice-section">
          <div class="section-header">
            <h3>公告通知</h3>
            <router-link to="/notices" class="more-link">查看更多 ></router-link>
          </div>
          <el-row :gutter="20">
            <el-col :span="12" v-for="notice in notices.slice(0, 4)" :key="notice.id">
              <div class="notice-item" @click="goToNotice(notice.id)">
                <el-icon class="notice-icon"><Bell /></el-icon>
                <span class="notice-title">{{ notice.title }}</span>
                <span class="notice-date">{{ formatDate(notice.created_at) }}</span>
              </div>
            </el-col>
          </el-row>
        </div>

        <div class="section category-section">
          <div class="section-header">
            <h3>商品分类</h3>
          </div>
          <el-row :gutter="20">
            <el-col :span="6" v-for="category in categories" :key="category.id">
              <div class="category-card card-hover" @click="goToCategory(category.id)">
                <el-avatar :size="60" class="category-icon">
                  <el-icon :size="30"><ShoppingBag /></el-icon>
                </el-avatar>
                <h4>{{ category.name }}</h4>
              </div>
            </el-col>
          </el-row>
        </div>

        <div class="section hot-section">
          <div class="section-header">
            <h3>热门商品</h3>
            <router-link to="/products" class="more-link">查看更多 ></router-link>
          </div>
          <el-row :gutter="20">
            <el-col :span="6" v-for="product in hotProducts" :key="product.id">
              <div class="product-card card-hover" @click="goToProduct(product.id)">
                <div class="product-image">
                  <img :src="product.image || 'https://picsum.photos/300/300?random=' + product.id" :alt="product.name" />
                  <div class="product-overlay">
                    <el-button type="primary" size="small">立即查看</el-button>
                  </div>
                </div>
                <div class="product-info">
                  <h4 class="product-name">{{ product.name }}</h4>
                  <div class="product-price">
                    <span class="price">¥{{ product.price }}</span>
                    <span class="sales">已售{{ product.sales }}件</span>
                  </div>
                </div>
              </div>
            </el-col>
          </el-row>
        </div>

        <div class="section news-section">
          <div class="section-header">
            <h3>新闻资讯</h3>
            <router-link to="/news" class="more-link">查看更多 ></router-link>
          </div>
          <el-row :gutter="20">
            <el-col :span="8" v-for="news in newsList.slice(0, 3)" :key="news.id">
              <div class="news-card card-hover" @click="goToNews(news.id)">
                <div class="news-image">
                  <img :src="news.cover || 'https://picsum.photos/400/250?random=' + news.id" :alt="news.title" />
                </div>
                <div class="news-info">
                  <h4 class="news-title">{{ news.title }}</h4>
                  <div class="news-meta">
                    <el-icon><View /></el-icon>
                    <span>{{ news.views }} 阅读</span>
                    <span class="news-date">{{ formatDate(news.created_at) }}</span>
                  </div>
                </div>
              </div>
            </el-col>
          </el-row>
        </div>

        <div class="section video-section">
          <div class="section-header">
            <h3>赛事视频</h3>
            <router-link to="/videos" class="more-link">查看更多 ></router-link>
          </div>
          <el-row :gutter="20">
            <el-col :span="6" v-for="video in videos.slice(0, 4)" :key="video.id">
              <div class="video-card card-hover" @click="goToVideo(video.id)">
                <div class="video-image">
                  <img :src="video.cover || 'https://picsum.photos/300/200?random=' + video.id" :alt="video.title" />
                  <div class="video-play">
                    <el-icon :size="40"><VideoPlay /></el-icon>
                  </div>
                </div>
                <div class="video-info">
                  <h4 class="video-title">{{ video.title }}</h4>
                  <div class="video-meta">
                    <el-icon><View /></el-icon>
                    <span>{{ video.views }} 播放</span>
                  </div>
                </div>
              </div>
            </el-col>
          </el-row>
        </div>
      </div>
    </div>
  </Layout>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import {
  Bell,
  ShoppingBag,
  View,
  VideoPlay
} from '@element-plus/icons-vue'
import Layout from '@/components/Layout.vue'
import {
  getBanners,
  getNotices,
  getCategories,
  getProducts,
  getNews,
  getVideos
} from '@/api'

const router = useRouter()

const banners = ref([])
const notices = ref([])
const categories = ref([])
const hotProducts = ref([])
const newsList = ref([])
const videos = ref([])

const fetchData = async () => {
  try {
    const [bannerRes, noticeRes, categoryRes, productRes, newsRes, videoRes] = await Promise.all([
      getBanners(),
      getNotices({ page: 1, page_size: 8 }),
      getCategories(),
      getProducts({ page: 1, page_size: 8 }),
      getNews({ page: 1, page_size: 6 }),
      getVideos({ page: 1, page_size: 8 })
    ])

    banners.value = bannerRes.data || []
    notices.value = noticeRes.data?.list || []
    categories.value = categoryRes.data || []
    hotProducts.value = productRes.data?.list || []
    newsList.value = newsRes.data?.list || []
    videos.value = videoRes.data?.list || []
  } catch (e) {
    console.error(e)
  }
}

const goToLink = (link) => {
  if (link) {
    router.push(link)
  }
}

const goToNotice = (id) => {
  router.push(`/notice/${id}`)
}

const goToCategory = (id) => {
  router.push({ path: '/products', query: { category_id: id }})
}

const goToProduct = (id) => {
  router.push(`/product/${id}`)
}

const goToNews = (id) => {
  router.push(`/news/${id}`)
}

const goToVideo = (id) => {
  router.push(`/video/${id}`)
}

const formatDate = (date) => {
  if (!date) return ''
  return new Date(date).toLocaleDateString('zh-CN')
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.home-page {
  padding-bottom: 40px;
}

.banner-carousel {
  margin-bottom: 30px;
}

.banner-item {
  position: relative;
  width: 100%;
  height: 100%;
  cursor: pointer;
}

.banner-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.banner-overlay {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  background: linear-gradient(transparent, rgba(0, 0, 0, 0.7));
  padding: 20px;
  color: #fff;
}

.section {
  background: #fff;
  border-radius: 8px;
  padding: 20px;
  margin-bottom: 30px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 15px;
  border-bottom: 2px solid #f0f0f0;
}

.section-header h3 {
  font-size: 20px;
  color: #333;
  margin: 0;
}

.more-link {
  color: #409EFF;
  text-decoration: none;
  font-size: 14px;
}

.more-link:hover {
  text-decoration: underline;
}

.notice-item {
  display: flex;
  align-items: center;
  padding: 12px 0;
  border-bottom: 1px dashed #f0f0f0;
  cursor: pointer;
  transition: all 0.3s;
}

.notice-item:hover {
  color: #409EFF;
}

.notice-icon {
  color: #409EFF;
  margin-right: 10px;
}

.notice-title {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.notice-date {
  color: #999;
  font-size: 12px;
  margin-left: 15px;
}

.category-card {
  text-align: center;
  padding: 30px 20px;
  cursor: pointer;
  border-radius: 8px;
  background: #fafafa;
  transition: all 0.3s;
}

.category-icon {
  background: linear-gradient(135deg, #409EFF, #67C23A);
  margin-bottom: 15px;
}

.category-icon :deep(svg) {
  color: #fff;
}

.category-card h4 {
  margin: 0;
  color: #333;
}

.product-card {
  background: #fff;
  border-radius: 8px;
  overflow: hidden;
  cursor: pointer;
  border: 1px solid #f0f0f0;
}

.product-image {
  position: relative;
  height: 200px;
  overflow: hidden;
}

.product-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s;
}

.product-card:hover .product-image img {
  transform: scale(1.05);
}

.product-overlay {
  position: absolute;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: opacity 0.3s;
}

.product-card:hover .product-overlay {
  opacity: 1;
}

.product-info {
  padding: 15px;
}

.product-name {
  font-size: 14px;
  color: #333;
  margin: 0 0 10px 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.product-price {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.price {
  font-size: 18px;
  color: #f56c6c;
  font-weight: bold;
}

.sales {
  font-size: 12px;
  color: #999;
}

.news-card {
  background: #fff;
  border-radius: 8px;
  overflow: hidden;
  cursor: pointer;
  border: 1px solid #f0f0f0;
}

.news-image {
  height: 180px;
  overflow: hidden;
}

.news-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s;
}

.news-card:hover .news-image img {
  transform: scale(1.05);
}

.news-info {
  padding: 15px;
}

.news-title {
  font-size: 15px;
  color: #333;
  margin: 0 0 10px 0;
  height: 40px;
  overflow: hidden;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.news-meta {
  display: flex;
  align-items: center;
  gap: 10px;
  color: #999;
  font-size: 12px;
}

.news-date {
  margin-left: auto;
}

.video-card {
  background: #fff;
  border-radius: 8px;
  overflow: hidden;
  cursor: pointer;
  border: 1px solid #f0f0f0;
}

.video-image {
  position: relative;
  height: 150px;
  overflow: hidden;
}

.video-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.video-play {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.3);
  color: #fff;
  opacity: 0;
  transition: opacity 0.3s;
}

.video-card:hover .video-play {
  opacity: 1;
}

.video-info {
  padding: 12px;
}

.video-title {
  font-size: 14px;
  color: #333;
  margin: 0 0 8px 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.video-meta {
  display: flex;
  align-items: center;
  gap: 5px;
  color: #999;
  font-size: 12px;
}
</style>
