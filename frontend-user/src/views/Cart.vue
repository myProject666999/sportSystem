<template>
  <div class="cart-page">
    <div class="container">
      <h2 class="page-title">我的购物车</h2>

      <template v-if="loading">
        <el-skeleton :rows="5" animated />
      </template>

      <template v-else-if="cartItems.length === 0">
        <el-empty description="购物车是空的，快去选购商品吧">
          <el-button type="primary" @click="$router.push('/products')">去购物</el-button>
        </el-empty>
      </template>

      <template v-else>
        <el-table :data="cartItems" stripe>
          <el-table-column width="55">
            <template #header>
              <el-checkbox v-model="selectAll" @change="toggleSelectAll" />
            </template>
            <template #default="{ row }">
              <el-checkbox v-model="row.selected" @change="updateSelectStatus" />
            </template>
          </el-table-column>
          <el-table-column label="商品信息" min-width="300">
            <template #default="{ row }">
              <div class="product-info" @click="$router.push(`/product/${row.product_id}`)">
                <el-image :src="row.product?.image" fit="cover" class="product-img" />
                <div class="product-detail">
                  <div class="product-name">{{ row.product?.name }}</div>
                </div>
              </div>
            </template>
          </el-table-column>
          <el-table-column label="单价" width="120">
            <template #default="{ row }">
              <span class="price">¥{{ row.product?.price }}</span>
            </template>
          </el-table-column>
          <el-table-column label="数量" width="140">
            <template #default="{ row }">
              <el-input-number
                v-model="row.quantity"
                :min="1"
                :max="row.product?.stock || 999"
                @change="updateQuantity(row)"
                size="small"
              />
            </template>
          </el-table-column>
          <el-table-column label="小计" width="120">
            <template #default="{ row }">
              <span class="subtotal">¥{{ (row.quantity * (row.product?.price || 0)).toFixed(2) }}</span>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="100" fixed="right">
            <template #default="{ row }">
              <el-button type="danger" link @click="removeItem(row)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>

        <div class="cart-footer">
          <div class="footer-left">
            <el-checkbox v-model="selectAll" @change="toggleSelectAll">全选</el-checkbox>
            <el-button type="text" style="margin-left: 20px;" @click="clearSelected">删除选中</el-button>
          </div>
          <div class="footer-right">
            <div class="total-info">
              <span class="total-label">已选 <span class="total-count">{{ selectedCount }}</span> 件商品</span>
              <span class="total-amount">
                合计: <span class="amount-price">¥{{ selectedTotal.toFixed(2) }}</span>
              </span>
            </div>
            <el-button type="danger" size="large" @click="checkout" :disabled="selectedCount === 0">
              结算
            </el-button>
          </div>
        </div>
      </template>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getCartList, updateCart, removeFromCart } from '@/api'

const router = useRouter()
const loading = ref(false)
const cartItems = ref([])
const selectAll = ref(false)

const selectedCount = computed(() => {
  return cartItems.value.filter(item => item.selected).length
})

const selectedTotal = computed(() => {
  return cartItems.value
    .filter(item => item.selected)
    .reduce((sum, item) => sum + (item.quantity * (item.product?.price || 0)), 0)
})

const loadCart = async () => {
  const token = localStorage.getItem('token')
  if (!token) {
    ElMessage.warning('请先登录')
    router.push('/login')
    return
  }
  loading.value = true
  try {
    const res = await getCartList()
    cartItems.value = (res.data?.list || res.data || []).map(item => ({
      ...item,
      selected: false
    }))
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const toggleSelectAll = (val) => {
  cartItems.value.forEach(item => {
    item.selected = val
  })
}

const updateSelectStatus = () => {
  selectAll.value = cartItems.value.every(item => item.selected)
}

const updateQuantity = async (item) => {
  try {
    await updateCart(item.id, { quantity: item.quantity })
  } catch (e) {
    console.error(e)
  }
}

const removeItem = async (item) => {
  try {
    await ElMessageBox.confirm('确定要删除该商品吗？', '提示', {
      type: 'warning'
    })
    await removeFromCart(item.id)
    ElMessage.success('已删除')
    loadCart()
  } catch (e) {
    if (e !== 'cancel') console.error(e)
  }
}

const clearSelected = async () => {
  const selectedItems = cartItems.value.filter(item => item.selected)
  if (selectedItems.length === 0) {
    ElMessage.warning('请选择要删除的商品')
    return
  }
  try {
    await ElMessageBox.confirm(`确定要删除选中的 ${selectedItems.length} 件商品吗？`, '提示', {
      type: 'warning'
    })
    for (const item of selectedItems) {
      await removeFromCart(item.id)
    }
    ElMessage.success('已删除')
    loadCart()
  } catch (e) {
    if (e !== 'cancel') console.error(e)
  }
}

const checkout = () => {
  const selectedIds = cartItems.value
    .filter(item => item.selected)
    .map(item => item.id)
  
  if (selectedIds.length === 0) {
    ElMessage.warning('请选择要结算的商品')
    return
  }
  router.push('/checkout')
}

onMounted(() => {
  loadCart()
})
</script>

<style scoped>
.cart-page {
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

.product-info {
  display: flex;
  align-items: center;
  cursor: pointer;
}

.product-img {
  width: 80px;
  height: 80px;
  border-radius: 4px;
  border: 1px solid #eee;
}

.product-detail {
  margin-left: 15px;
}

.product-name {
  color: #333;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.price {
  color: #e4393c;
  font-weight: bold;
}

.subtotal {
  color: #e4393c;
  font-weight: bold;
  font-size: 16px;
}

.cart-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: #fff;
  padding: 15px 20px;
  margin-top: 20px;
  border-radius: 4px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
}

.footer-left {
  display: flex;
  align-items: center;
}

.footer-right {
  display: flex;
  align-items: center;
  gap: 20px;
}

.total-info {
  display: flex;
  align-items: center;
  gap: 30px;
}

.total-label {
  color: #666;
}

.total-count {
  color: #e4393c;
  font-weight: bold;
}

.total-amount {
  font-size: 16px;
}

.amount-price {
  color: #e4393c;
  font-size: 22px;
  font-weight: bold;
}
</style>
