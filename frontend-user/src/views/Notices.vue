<template>
  <div class="notices-page">
    <div class="container">
      <h2 class="page-title">系统公告</h2>

      <div class="notices-list" v-loading="loading">
        <div v-if="notices.length === 0">
          <el-empty description="暂无公告" />
        </div>
        <div v-else>
          <div
            v-for="notice in notices"
            :key="notice.id"
            class="notice-item"
            @click="$router.push(`/notice/${notice.id}`)"
          >
            <div class="notice-header">
              <span class="notice-title">{{ notice.title }}</span>
              <span class="notice-date">{{ notice.created_at }}</span>
            </div>
            <div class="notice-content">
              {{ notice.content?.replace(/<[^>]+>/g, '').substring(0, 200) }}...
            </div>
            <div class="notice-footer">
              <span class="view-more">查看详情</span>
            </div>
          </div>

          <div class="pagination-box">
            <el-pagination
              v-model:current-page="pagination.page"
              v-model:page-size="pagination.pageSize"
              :page-sizes="[10, 20, 50]"
              :total="pagination.total"
              layout="prev, pager, next"
              @current-change="loadNotices"
            />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getNotices } from '@/api'

const loading = ref(false)
const notices = ref([])

const pagination = ref({
  page: 1,
  pageSize: 10,
  total: 0
})

const loadNotices = async () => {
  loading.value = true
  try {
    const res = await getNotices({
      page: pagination.value.page,
      page_size: pagination.value.pageSize
    })
    notices.value = res.data?.list || []
    pagination.value.total = res.data?.total || 0
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadNotices()
})
</script>

<style scoped>
.notices-page {
  padding: 30px 0;
  background: #f5f5f5;
  min-height: calc(100vh - 140px);
}

.container {
  max-width: 1000px;
  margin: 0 auto;
  padding: 0 15px;
}

.page-title {
  font-size: 22px;
  color: #333;
  margin-bottom: 20px;
}

.notices-list {
  background: #fff;
  border-radius: 8px;
  padding: 20px;
}

.notice-item {
  padding: 20px;
  border-bottom: 1px solid #eee;
  cursor: pointer;
  transition: all 0.3s;
}

.notice-item:hover {
  background: #fafafa;
}

.notice-item:last-child {
  border-bottom: none;
}

.notice-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.notice-title {
  font-size: 18px;
  color: #333;
  font-weight: bold;
}

.notice-item:hover .notice-title {
  color: #409EFF;
}

.notice-date {
  color: #999;
  font-size: 14px;
}

.notice-content {
  color: #666;
  font-size: 14px;
  line-height: 1.8;
  margin-bottom: 10px;
}

.notice-footer {
  text-align: right;
}

.view-more {
  color: #409EFF;
  font-size: 14px;
  cursor: pointer;
}

.view-more:hover {
  text-decoration: underline;
}

.pagination-box {
  display: flex;
  justify-content: center;
  margin-top: 30px;
}
</style>
