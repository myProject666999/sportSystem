<template>
  <div class="favorites-page">
    <div class="container">
      <h2 class="page-title">我的收藏</h2>

      <div class="product-list" v-loading="loading">
        <div v-if="favorites.length === 0">
          <el-empty description="暂无收藏商品">
            <el-button type="primary" @click="$router.push('/products')">去逛逛</el-button>
          </el-empty>
        </div>
        <el-row :gutter="20" v-else>
          <el-col :span="6" v-for="item in favorites" :key="item.id">
            <div class="product-card" @click="$router.push(`/product/${item.product_id}`)">
              <div class="product-img-wrapper">
                <el-image :src="item.product?.image" fit="cover" class="product-img" />
              </div>
              <div class="product-info">
                <div class="product-name">{{ item.product?.name }}</div>
                <div class="product-price">
                  <span class="price">¥{{ item.product?.price }}</span>
                </div>
              </div>
              <div class="product-actions">
                <el-button type="danger" size="small" @click.stop="handleRemoveFavorite(item)">
                  <el-icon><Delete /></el-icon>
                  取消收藏
                </el-button>
              </div>
            </div>
          </el-col>
        </el-row>

        <div class="pagination-box" v-if="favorites.length > 0">
          <el-pagination
            v-model:current-page="pagination.page"
            v-model:page-size="pagination.pageSize"
            :page-sizes="[12, 24, 48]"
            :total="pagination.total"
            layout="prev, pager, next"
            @current-change="loadFavorites"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Delete } from '@element-plus/icons-vue'
import { getFavorites, removeFavorite } from '@/api'

const router = useRouter()
const loading = ref(false)
const favorites = ref([])

const pagination = ref({
  page: 1,
  pageSize: 12,
  total: 0
})

const loadFavorites = async () => {
  const token = localStorage.getItem('token')
  if (!token) {
    ElMessage.warning('请先登录')
    router.push('/login')
    return
  }
  loading.value = true
  try {
    const res = await getFavorites({
      page: pagination.value.page,
      page_size: pagination.value.pageSize
    })
    favorites.value = res.data?.list || []
    pagination.value.total = res.data?.total || 0
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const handleRemoveFavorite = async (item) => {
  try {
    await ElMessageBox.confirm('确定要取消收藏该商品吗？', '提示', {
      type: 'warning'
    })
    await removeFavorite(item.id)
    ElMessage.success('已取消收藏')
    loadFavorites()
  } catch (e) {
    if (e !== 'cancel') console.error(e)
  }
}

onMounted(() => {
  loadFavorites()
})
</script>

<style scoped>
.favorites-page {
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

.product-list {
  background: #fff;
  border-radius: 8px;
  padding: 20px;
}

.product-card {
  background: #fff;
  border: 1px solid #eee;
  border-radius: 8px;
  overflow: hidden;
  cursor: pointer;
  transition: all 0.3s;
  margin-bottom: 20px;
}

.product-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  transform: translateY(-2px);
}

.product-img-wrapper {
  width: 100%;
  height: 200px;
  overflow: hidden;
}

.product-img {
  width: 100%;
  height: 100%;
  transition: transform 0.3s;
}

.product-card:hover .product-img {
  transform: scale(1.05);
}

.product-info {
  padding: 15px;
}

.product-name {
  color: #333;
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  height: 42px;
  margin-bottom: 10px;
}

.product-price {
  display: flex;
  align-items: baseline;
}

.price {
  color: #e4393c;
  font-size: 20px;
  font-weight: bold;
}

.product-actions {
  padding: 0 15px 15px;
}

.pagination-box {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}
</style>
