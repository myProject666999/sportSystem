<template>
  <Layout>
    <div class="products-page container">
      <el-row :gutter="20">
        <el-col :span="5">
          <div class="filter-box">
            <h3 class="filter-title">商品分类</h3>
            <el-menu
              :default-active="activeCategory"
              background-color="#fff"
              text-color="#333"
              active-text-color="#409EFF"
              @select="handleCategorySelect"
            >
              <el-menu-item index="">全部商品</el-menu-item>
              <el-menu-item
                v-for="category in categories"
                :key="category.id"
                :index="String(category.id)"
              >
                {{ category.name }}
              </el-menu-item>
            </el-menu>
          </div>
        </el-col>
        <el-col :span="19">
          <div class="products-header">
            <div class="sort-bar">
              <span class="sort-label">排序：</span>
              <el-radio-group v-model="sortBy" size="small">
                <el-radio-button label="">默认</el-radio-button>
                <el-radio-button label="sales">销量</el-radio-button>
                <el-radio-button label="price_asc">价格升序</el-radio-button>
                <el-radio-button label="price_desc">价格降序</el-radio-button>
              </el-radio-group>
            </div>
            <div class="search-bar">
              <el-input
                v-model="searchKeyword"
                placeholder="搜索商品"
                style="width: 200px"
                size="small"
                clearable
                @clear="fetchProducts"
                @keyup.enter="fetchProducts"
              >
                <template #append>
                  <el-button icon="Search" size="small" @click="fetchProducts"></el-button>
                </template>
              </el-input>
            </div>
          </div>

          <el-row :gutter="20" v-loading="loading">
            <el-col :span="6" v-for="product in products" :key="product.id">
              <div class="product-card card-hover" @click="goToProduct(product.id)">
                <div class="product-image">
                  <img :src="product.image || 'https://picsum.photos/300/300?random=' + product.id" :alt="product.name" />
                </div>
                <div class="product-info">
                  <h4 class="product-name">{{ product.name }}</h4>
                  <div class="product-price">
                    <span class="price">¥{{ product.price }}</span>
                    <span class="points" v-if="product.points_price">积分{{ product.points_price }}</span>
                  </div>
                  <div class="product-meta">
                    <span>库存{{ product.stock }}</span>
                    <span>已售{{ product.sales }}</span>
                  </div>
                </div>
              </div>
            </el-col>
          </el-row>

          <div class="pagination-wrap" v-if="total > 0">
            <el-pagination
              v-model:current-page="page"
              v-model:page-size="pageSize"
              :page-sizes="[12, 24, 48]"
              :total="total"
              layout="total, sizes, prev, pager, next, jumper"
              @size-change="fetchProducts"
              @current-change="fetchProducts"
            />
          </div>

          <el-empty v-if="!loading && products.length === 0" description="暂无商品" />
        </el-col>
      </el-row>
    </div>
  </Layout>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import Layout from '@/components/Layout.vue'
import { getCategories, getProducts } from '@/api'

const router = useRouter()
const route = useRoute()

const loading = ref(false)
const categories = ref([])
const products = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(12)
const sortBy = ref('')
const searchKeyword = ref('')
const categoryId = ref('')

const activeCategory = computed(() => {
  return categoryId.value || ''
})

const handleCategorySelect = (index) => {
  categoryId.value = index
  page.value = 1
  fetchProducts()
}

const goToProduct = (id) => {
  router.push(`/product/${id}`)
}

const fetchCategories = async () => {
  try {
    const res = await getCategories()
    categories.value = res.data || []
  } catch (e) {
    console.error(e)
  }
}

const fetchProducts = async () => {
  loading.value = true
  try {
    const params = {
      page: page.value,
      page_size: pageSize.value
    }
    if (categoryId.value) {
      params.category_id = categoryId.value
    }
    if (searchKeyword.value) {
      params.keyword = searchKeyword.value
    }

    const res = await getProducts(params)
    products.value = res.data?.list || []
    total.value = res.data?.total || 0
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

watch(() => route.query, (query) => {
  if (query.category_id) {
    categoryId.value = query.category_id
  }
  if (query.keyword) {
    searchKeyword.value = query.keyword
  }
  page.value = 1
  fetchProducts()
}, { immediate: true })

onMounted(() => {
  fetchCategories()
})
</script>

<style scoped>
.products-page {
  padding: 20px 0;
}

.filter-box {
  background: #fff;
  border-radius: 8px;
  padding: 15px;
}

.filter-title {
  margin: 0 0 15px 0;
  font-size: 16px;
  color: #333;
  padding-bottom: 10px;
  border-bottom: 1px solid #f0f0f0;
}

.filter-box :deep(.el-menu) {
  border-right: none;
}

.filter-box :deep(.el-menu-item) {
  height: 40px;
  line-height: 40px;
  margin-bottom: 5px;
  border-radius: 4px;
}

.products-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding: 15px 20px;
  background: #fff;
  border-radius: 8px;
}

.sort-label {
  margin-right: 10px;
  color: #666;
}

.product-card {
  background: #fff;
  border-radius: 8px;
  overflow: hidden;
  cursor: pointer;
  border: 1px solid #f0f0f0;
  margin-bottom: 20px;
}

.product-image {
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

.product-info {
  padding: 15px;
}

.product-name {
  font-size: 14px;
  color: #333;
  margin: 0 0 10px 0;
  height: 40px;
  overflow: hidden;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.product-price {
  display: flex;
  align-items: center;
  margin-bottom: 8px;
}

.price {
  font-size: 18px;
  color: #f56c6c;
  font-weight: bold;
  margin-right: 10px;
}

.points {
  font-size: 12px;
  color: #67c23a;
}

.product-meta {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  color: #999;
}

.pagination-wrap {
  margin-top: 20px;
  text-align: center;
}
</style>
