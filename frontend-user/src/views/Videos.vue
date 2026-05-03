<template>
  <div class="videos-page">
    <div class="container">
      <h2 class="page-title">赛事视频</h2>

      <div class="videos-content">
        <div class="videos-grid" v-loading="loading">
          <div v-if="videoList.length === 0">
            <el-empty description="暂无视频" />
          </div>
          <el-row :gutter="20" v-else>
            <el-col :span="6" v-for="video in videoList" :key="video.id">
              <div class="video-card" @click="$router.push(`/video/${video.id}`)">
                <div class="video-cover">
                  <el-image :src="video.cover" fit="cover" class="cover-img" />
                  <div class="video-play">
                    <el-icon :size="32"><VideoPlay /></el-icon>
                  </div>
                  <div class="video-duration" v-if="video.duration">
                    {{ video.duration }}
                  </div>
                </div>
                <div class="video-info">
                  <h3 class="video-title">{{ video.title }}</h3>
                  <div class="video-meta">
                    <span class="meta-item">
                      <el-icon><View /></el-icon>
                      {{ video.views || 0 }}
                    </span>
                    <span class="meta-item">
                      <el-icon><Clock /></el-icon>
                      {{ video.created_at }}
                    </span>
                  </div>
                </div>
              </div>
            </el-col>
          </el-row>

          <div class="pagination-box" v-if="videoList.length > 0">
            <el-pagination
              v-model:current-page="pagination.page"
              v-model:page-size="pagination.pageSize"
              :page-sizes="[12, 24, 48]"
              :total="pagination.total"
              layout="prev, pager, next"
              @current-change="loadVideos"
            />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { View, Clock, VideoPlay } from '@element-plus/icons-vue'
import { getVideos } from '@/api'

const loading = ref(false)
const videoList = ref([])

const pagination = ref({
  page: 1,
  pageSize: 12,
  total: 0
})

const loadVideos = async () => {
  loading.value = true
  try {
    const res = await getVideos({
      page: pagination.value.page,
      page_size: pagination.value.pageSize,
      status: 1
    })
    videoList.value = res.data?.list || []
    pagination.value.total = res.data?.total || 0
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadVideos()
})
</script>

<style scoped>
.videos-page {
  padding: 30px 0;
  background: #f5f5f5;
  min-height: calc(100vh - 140px);
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 15px;
}

.page-title {
  font-size: 22px;
  color: #333;
  margin-bottom: 20px;
}

.videos-grid {
  background: #fff;
  border-radius: 8px;
  padding: 20px;
}

.video-card {
  background: #fff;
  border-radius: 8px;
  overflow: hidden;
  cursor: pointer;
  transition: all 0.3s;
  margin-bottom: 20px;
  border: 1px solid #eee;
}

.video-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  transform: translateY(-2px);
}

.video-cover {
  width: 100%;
  height: 160px;
  position: relative;
  overflow: hidden;
}

.cover-img {
  width: 100%;
  height: 100%;
  transition: transform 0.3s;
}

.video-card:hover .cover-img {
  transform: scale(1.05);
}

.video-play {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 60px;
  height: 60px;
  background: rgba(0, 0, 0, 0.5);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  transition: all 0.3s;
}

.video-card:hover .video-play {
  background: rgba(64, 158, 255, 0.8);
  transform: translate(-50%, -50%) scale(1.1);
}

.video-duration {
  position: absolute;
  bottom: 10px;
  right: 10px;
  background: rgba(0, 0, 0, 0.6);
  color: #fff;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
}

.video-info {
  padding: 15px;
}

.video-title {
  font-size: 14px;
  color: #333;
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  margin-bottom: 10px;
  font-weight: 500;
}

.video-card:hover .video-title {
  color: #409EFF;
}

.video-meta {
  display: flex;
  gap: 15px;
  color: #999;
  font-size: 12px;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 4px;
}

.pagination-box {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}
</style>
