<template>
  <div class="news-page">
    <div class="container">
      <div class="page-header">
        <h2 class="page-title">新闻资讯</h2>
        <el-input
          v-model="keyword"
          placeholder="搜索新闻"
          class="search-input"
          clearable
          @clear="search"
          @keyup.enter="search"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
      </div>

      <div class="news-content">
        <div class="news-list" v-loading="loading">
          <div v-if="newsList.length === 0">
            <el-empty description="暂无新闻" />
          </div>
          <div v-else>
            <div
              v-for="news in newsList"
              :key="news.id"
              class="news-item"
              @click="$router.push(`/news/${news.id}`)"
            >
              <div class="news-img" v-if="news.cover">
                <el-image :src="news.cover" fit="cover" />
              </div>
              <div class="news-info" :class="{ 'full-width': !news.cover }">
                <h3 class="news-title">{{ news.title }}</h3>
                <p class="news-summary">{{ news.summary || news.content?.replace(/<[^>]+>/g, '').substring(0, 150) }}...</p>
                <div class="news-meta">
                  <span class="meta-item">
                    <el-icon><View /></el-icon>
                    {{ news.views || 0 }} 阅读
                  </span>
                  <span class="meta-item">
                    <el-icon><Clock /></el-icon>
                    {{ news.created_at }}
                  </span>
                </div>
              </div>
            </div>

            <div class="pagination-box">
              <el-pagination
                v-model:current-page="pagination.page"
                v-model:page-size="pagination.pageSize"
                :page-sizes="[10, 20, 50]"
                :total="pagination.total"
                layout="prev, pager, next"
                @current-change="loadNews"
              />
            </div>
          </div>
        </div>

        <div class="news-sidebar">
          <el-card class="sidebar-card" shadow="never">
            <template #header>
              <span class="sidebar-title">热门新闻</span>
            </template>
            <div class="hot-list">
              <div
                v-for="(news, index) in hotNews"
                :key="news.id"
                class="hot-item"
                @click="$router.push(`/news/${news.id}`)"
              >
                <span class="hot-index" :class="{ 'top-three': index < 3 }">{{ index + 1 }}</span>
                <span class="hot-title">{{ news.title }}</span>
              </div>
            </div>
          </el-card>

          <el-card class="sidebar-card" shadow="never" style="margin-top: 20px;">
            <template #header>
              <span class="sidebar-title">最新视频</span>
            </template>
            <div class="video-list">
              <div
                v-for="video in latestVideos"
                :key="video.id"
                class="video-item"
                @click="$router.push(`/video/${video.id}`)"
              >
                <div class="video-cover">
                  <el-image :src="video.cover" fit="cover" />
                  <div class="video-play">
                    <el-icon :size="24"><VideoPlay /></el-icon>
                  </div>
                </div>
                <div class="video-title">{{ video.title }}</div>
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
import { Search, View, Clock, VideoPlay } from '@element-plus/icons-vue'
import { getNews, getVideos } from '@/api'

const loading = ref(false)
const keyword = ref('')
const newsList = ref([])
const hotNews = ref([])
const latestVideos = ref([])

const pagination = ref({
  page: 1,
  pageSize: 10,
  total: 0
})

const loadNews = async () => {
  loading.value = true
  try {
    const res = await getNews({
      page: pagination.value.page,
      page_size: pagination.value.pageSize,
      keyword: keyword.value
    })
    newsList.value = res.data?.list || []
    pagination.value.total = res.data?.total || 0
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const loadHotNews = async () => {
  try {
    const res = await getNews({ page_size: 8 })
    hotNews.value = res.data?.list || []
  } catch (e) {
    console.error(e)
  }
}

const loadLatestVideos = async () => {
  try {
    const res = await getVideos({ page_size: 4, status: 1 })
    latestVideos.value = res.data?.list || []
  } catch (e) {
    console.error(e)
  }
}

const search = () => {
  pagination.value.page = 1
  loadNews()
}

onMounted(() => {
  loadNews()
  loadHotNews()
  loadLatestVideos()
})
</script>

<style scoped>
.news-page {
  padding: 30px 0;
  background: #f5f5f5;
  min-height: calc(100vh - 140px);
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 15px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.page-title {
  font-size: 22px;
  color: #333;
}

.search-input {
  width: 300px;
}

.news-content {
  display: flex;
  gap: 20px;
}

.news-list {
  flex: 1;
  background: #fff;
  border-radius: 8px;
  padding: 20px;
}

.news-item {
  display: flex;
  gap: 20px;
  padding: 20px 0;
  border-bottom: 1px solid #eee;
  cursor: pointer;
  transition: all 0.3s;
}

.news-item:hover {
  background: #fafafa;
  margin: 0 -20px;
  padding-left: 20px;
  padding-right: 20px;
}

.news-item:last-child {
  border-bottom: none;
}

.news-img {
  width: 200px;
  height: 140px;
  flex-shrink: 0;
  border-radius: 8px;
  overflow: hidden;
}

.news-img .el-image {
  width: 100%;
  height: 100%;
}

.news-info {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.news-info.full-width {
  width: 100%;
}

.news-title {
  font-size: 18px;
  color: #333;
  font-weight: bold;
  margin-bottom: 10px;
  line-height: 1.4;
}

.news-item:hover .news-title {
  color: #409EFF;
}

.news-summary {
  color: #666;
  font-size: 14px;
  line-height: 1.6;
  margin-bottom: 15px;
  flex: 1;
  overflow: hidden;
}

.news-meta {
  display: flex;
  gap: 20px;
  color: #999;
  font-size: 13px;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 5px;
}

.pagination-box {
  display: flex;
  justify-content: center;
  margin-top: 30px;
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

.video-list {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.video-item {
  cursor: pointer;
  transition: all 0.3s;
}

.video-item:hover .video-title {
  color: #409EFF;
}

.video-cover {
  width: 100%;
  height: 150px;
  border-radius: 8px;
  overflow: hidden;
  position: relative;
  margin-bottom: 10px;
}

.video-cover .el-image {
  width: 100%;
  height: 100%;
}

.video-play {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 50px;
  height: 50px;
  background: rgba(0, 0, 0, 0.5);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
}

.video-title {
  font-size: 14px;
  color: #333;
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  transition: color 0.3s;
}
</style>
