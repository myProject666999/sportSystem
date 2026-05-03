<template>
  <div class="video-detail-page">
    <div class="container">
      <div class="breadcrumb">
        <span @click="$router.push('/')">首页</span>
        <span class="separator">/</span>
        <span @click="$router.push('/videos')">赛事视频</span>
        <span class="separator">/</span>
        <span class="active">{{ video.title }}</span>
      </div>

      <div class="video-content" v-loading="loading">
        <div class="video-main">
          <div class="video-player">
            <div class="video-placeholder" v-if="!video.url">
              <el-image :src="video.cover" fit="cover" class="placeholder-img" />
              <div class="play-overlay">
                <el-icon :size="64"><VideoPlay /></el-icon>
                <p class="play-tip">点击播放视频</p>
              </div>
            </div>
            <video
              v-else
              :src="video.url"
              controls
              class="video-element"
              :poster="video.cover"
            ></video>
          </div>

          <div class="video-info">
            <h1 class="video-title">{{ video.title }}</h1>
            <div class="video-meta">
              <span class="meta-item">
                <el-icon><View /></el-icon>
                {{ video.views || 0 }} 播放
              </span>
              <span class="meta-item">
                <el-icon><Clock /></el-icon>
                {{ video.created_at }}
              </span>
            </div>
          </div>

          <div class="video-description" v-if="video.description">
            <h3 class="section-title">视频简介</h3>
            <p class="description-text">{{ video.description }}</p>
          </div>

          <div class="recommend-section">
            <h3 class="section-title">更多精彩视频</h3>
            <div class="recommend-list">
              <div
                v-for="item in recommendVideos"
                :key="item.id"
                class="recommend-item"
                @click="$router.push(`/video/${item.id}`)"
              >
                <div class="recommend-cover">
                  <el-image :src="item.cover" fit="cover" />
                  <div class="recommend-play">
                    <el-icon :size="20"><VideoPlay /></el-icon>
                  </div>
                </div>
                <div class="recommend-info">
                  <div class="recommend-title">{{ item.title }}</div>
                  <div class="recommend-views">
                    <el-icon><View /></el-icon>
                    {{ item.views || 0 }} 播放
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="video-sidebar">
          <el-card class="sidebar-card" shadow="never">
            <template #header>
              <span class="sidebar-title">热门视频</span>
            </template>
            <div class="hot-list">
              <div
                v-for="(item, index) in hotList"
                :key="item.id"
                class="hot-item"
                @click="$router.push(`/video/${item.id}`)"
              >
                <span class="hot-index" :class="{ 'top-three': index < 3 }">{{ index + 1 }}</span>
                <div class="hot-info">
                  <div class="hot-title">{{ item.title }}</div>
                  <div class="hot-views">
                    <el-icon><View /></el-icon>
                    {{ item.views || 0 }}
                  </div>
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
import { View, Clock, VideoPlay } from '@element-plus/icons-vue'
import { getVideoDetail, getVideos } from '@/api'

const route = useRoute()
const router = useRouter()
const loading = ref(false)
const video = ref({})
const hotList = ref([])
const recommendVideos = ref([])

const loadVideoDetail = async () => {
  loading.value = true
  try {
    const res = await getVideoDetail(route.params.id)
    video.value = res.data || {}
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const loadHotVideos = async () => {
  try {
    const res = await getVideos({ page_size: 8, status: 1 })
    hotList.value = res.data?.list || []
  } catch (e) {
    console.error(e)
  }
}

const loadRecommendVideos = async () => {
  try {
    const res = await getVideos({ page_size: 4, status: 1 })
    recommendVideos.value = res.data?.list || []
  } catch (e) {
    console.error(e)
  }
}

onMounted(() => {
  loadVideoDetail()
  loadHotVideos()
  loadRecommendVideos()
})
</script>

<style scoped>
.video-detail-page {
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

.video-content {
  display: flex;
  gap: 20px;
}

.video-main {
  flex: 1;
}

.video-player {
  width: 100%;
  background: #000;
  border-radius: 8px;
  overflow: hidden;
}

.video-placeholder {
  position: relative;
  width: 100%;
  height: 500px;
  cursor: pointer;
}

.placeholder-img {
  width: 100%;
  height: 100%;
}

.play-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.3);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #fff;
  transition: all 0.3s;
}

.video-placeholder:hover .play-overlay {
  background: rgba(0, 0, 0, 0.5);
}

.play-tip {
  margin-top: 15px;
  font-size: 16px;
}

.video-element {
  width: 100%;
  display: block;
}

.video-info {
  background: #fff;
  padding: 25px;
  margin-top: 20px;
  border-radius: 8px;
}

.video-title {
  font-size: 22px;
  color: #333;
  font-weight: bold;
  line-height: 1.4;
  margin-bottom: 15px;
}

.video-meta {
  display: flex;
  gap: 25px;
  color: #999;
  font-size: 14px;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 5px;
}

.video-description {
  background: #fff;
  padding: 25px;
  margin-top: 20px;
  border-radius: 8px;
}

.section-title {
  font-size: 16px;
  color: #333;
  font-weight: bold;
  margin-bottom: 15px;
}

.description-text {
  color: #666;
  line-height: 1.8;
  text-indent: 2em;
}

.recommend-section {
  background: #fff;
  padding: 25px;
  margin-top: 20px;
  border-radius: 8px;
}

.recommend-list {
  display: flex;
  gap: 20px;
}

.recommend-item {
  flex: 1;
  cursor: pointer;
  transition: all 0.3s;
}

.recommend-item:hover .recommend-title {
  color: #409EFF;
}

.recommend-cover {
  width: 100%;
  height: 120px;
  border-radius: 8px;
  overflow: hidden;
  position: relative;
  margin-bottom: 10px;
}

.recommend-cover .el-image {
  width: 100%;
  height: 100%;
}

.recommend-play {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 40px;
  height: 40px;
  background: rgba(0, 0, 0, 0.5);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
}

.recommend-title {
  font-size: 14px;
  color: #333;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  margin-bottom: 5px;
  transition: color 0.3s;
}

.recommend-views {
  color: #999;
  font-size: 12px;
  display: flex;
  align-items: center;
  gap: 4px;
}

.video-sidebar {
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
  gap: 12px;
  cursor: pointer;
  padding: 5px 0;
  transition: all 0.3s;
}

.hot-item:hover .hot-title {
  color: #409EFF;
}

.hot-index {
  width: 24px;
  height: 24px;
  background: #e4e7ed;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  color: #909399;
  flex-shrink: 0;
}

.hot-index.top-three {
  background: #e4393c;
  color: #fff;
}

.hot-info {
  flex: 1;
}

.hot-title {
  font-size: 14px;
  color: #333;
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  margin-bottom: 5px;
  transition: color 0.3s;
}

.hot-views {
  color: #999;
  font-size: 12px;
  display: flex;
  align-items: center;
  gap: 4px;
}
</style>
