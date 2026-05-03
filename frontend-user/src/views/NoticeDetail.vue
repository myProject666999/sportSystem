<template>
  <div class="notice-detail-page">
    <div class="container">
      <div class="breadcrumb">
        <span @click="$router.push('/')">首页</span>
        <span class="separator">/</span>
        <span @click="$router.push('/notices')">系统公告</span>
        <span class="separator">/</span>
        <span class="active">{{ notice.title }}</span>
      </div>

      <div class="notice-content" v-loading="loading">
        <div class="notice-main">
          <h1 class="notice-title">{{ notice.title }}</h1>
          <div class="notice-meta">
            <span class="meta-item">
              <el-icon><Clock /></el-icon>
              {{ notice.created_at }}
            </span>
          </div>
          <div class="notice-body" v-html="notice.content"></div>

          <div class="notice-navigation">
            <el-button type="primary" link @click="$router.push('/notices')">
              <el-icon><ArrowLeft /></el-icon>
              返回公告列表
            </el-button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Clock, ArrowLeft } from '@element-plus/icons-vue'
import { getNoticeDetail } from '@/api'

const route = useRoute()
const router = useRouter()
const loading = ref(false)
const notice = ref({})

const loadNoticeDetail = async () => {
  loading.value = true
  try {
    const res = await getNoticeDetail(route.params.id)
    notice.value = res.data || {}
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadNoticeDetail()
})
</script>

<style scoped>
.notice-detail-page {
  padding: 30px 0;
  background: #f5f5f5;
  min-height: calc(100vh - 140px);
}

.container {
  max-width: 1000px;
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

.notice-content {
  background: #fff;
  border-radius: 8px;
  padding: 40px;
}

.notice-main {
  max-width: 800px;
  margin: 0 auto;
}

.notice-title {
  font-size: 26px;
  color: #333;
  font-weight: bold;
  line-height: 1.4;
  text-align: center;
  margin-bottom: 20px;
}

.notice-meta {
  display: flex;
  justify-content: center;
  gap: 25px;
  color: #999;
  font-size: 14px;
  padding-bottom: 30px;
  border-bottom: 1px solid #eee;
  margin-bottom: 30px;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 5px;
}

.notice-body {
  font-size: 16px;
  line-height: 2;
  color: #333;
}

.notice-body :deep(img) {
  max-width: 100%;
  height: auto;
  border-radius: 8px;
  margin: 15px 0;
}

.notice-body :deep(p) {
  margin-bottom: 20px;
  text-indent: 2em;
}

.notice-body :deep(h1),
.notice-body :deep(h2),
.notice-body :deep(h3) {
  margin: 25px 0 15px;
  font-weight: bold;
}

.notice-navigation {
  margin-top: 50px;
  padding-top: 30px;
  border-top: 1px solid #eee;
  text-align: center;
}
</style>
