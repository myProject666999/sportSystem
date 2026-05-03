<template>
  <div class="orders-page">
    <div class="container">
      <h2 class="page-title">我的订单</h2>

      <div class="tabs">
        <div
          v-for="tab in tabs"
          :key="tab.value"
          :class="{ active: activeTab === tab.value }"
          @click="activeTab = tab.value; pagination.page = 1; loadOrders()"
        >
          {{ tab.label }}
        </div>
      </div>

      <div class="order-list" v-loading="loading">
        <div v-if="orders.length === 0">
          <el-empty description="暂无订单">
            <el-button type="primary" @click="$router.push('/products')">去购物</el-button>
          </el-empty>
        </div>
        <div v-else>
          <div v-for="order in orders" :key="order.id" class="order-card">
            <div class="order-header">
              <div class="order-info">
                <span class="order-no">订单号: {{ order.order_no }}</span>
                <span class="order-time">{{ order.created_at }}</span>
              </div>
              <div class="order-status">
                <el-tag :type="getStatusType(order.status)" size="large">{{ getStatusText(order.status) }}</el-tag>
              </div>
            </div>

            <div class="order-items" @click="$router.push(`/order/${order.id}`)">
              <div v-for="item in order.items" :key="item.id" class="order-item">
                <el-image :src="item.product_image" fit="cover" class="item-img" />
                <div class="item-info">
                  <div class="item-name">{{ item.product_name }}</div>
                  <div class="item-price">¥{{ item.price }} × {{ item.quantity }}</div>
                </div>
              </div>
            </div>

            <div class="order-footer">
              <div class="order-amount">
                <span>共 {{ order.items?.length || 0 }} 件商品</span>
                <span class="amount">订单金额: <strong>¥{{ order.total_amount }}</strong></span>
              </div>
              <div class="order-actions">
                <template v-if="order.status === 'pending'">
                  <el-button type="text" @click="handleCancelOrder(order)">取消订单</el-button>
                  <el-button type="danger" size="small" @click="handlePayOrder(order)">立即支付</el-button>
                </template>
                <template v-else-if="order.status === 'shipped'">
                  <el-button type="primary" size="small" @click="handleConfirmReceive(order)">确认收货</el-button>
                </template>
                <el-button type="text" @click="$router.push(`/order/${order.id}`)">查看详情</el-button>
              </div>
            </div>
          </div>
        </div>

        <div class="pagination-box" v-if="orders.length > 0">
          <el-pagination
            v-model:current-page="pagination.page"
            v-model:page-size="pagination.pageSize"
            :page-sizes="[10, 20, 50]"
            :total="pagination.total"
            layout="prev, pager, next"
            @current-change="loadOrders"
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
import { getOrders, cancelOrder, payOrder, confirmReceive } from '@/api'

const router = useRouter()
const loading = ref(false)
const orders = ref([])
const activeTab = ref('all')

const tabs = [
  { label: '全部', value: 'all' },
  { label: '待支付', value: 'pending' },
  { label: '待发货', value: 'paid' },
  { label: '待收货', value: 'shipped' },
  { label: '已完成', value: 'completed' },
]

const pagination = ref({
  page: 1,
  pageSize: 10,
  total: 0
})

const getStatusType = (status) => {
  const types = {
    pending: 'warning',
    paid: 'primary',
    shipped: 'info',
    completed: 'success',
    cancelled: 'danger'
  }
  return types[status] || 'info'
}

const getStatusText = (status) => {
  const texts = {
    pending: '待支付',
    paid: '已支付',
    shipped: '已发货',
    completed: '已完成',
    cancelled: '已取消'
  }
  return texts[status] || status
}

const loadOrders = async () => {
  const token = localStorage.getItem('token')
  if (!token) {
    ElMessage.warning('请先登录')
    router.push('/login')
    return
  }
  loading.value = true
  try {
    const params = {
      page: pagination.value.page,
      page_size: pagination.value.pageSize
    }
    if (activeTab.value !== 'all') {
      params.status = activeTab.value
    }
    const res = await getOrders(params)
    orders.value = res.data?.list || []
    pagination.value.total = res.data?.total || 0
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const handleCancelOrder = async (order) => {
  try {
    await ElMessageBox.confirm('确定要取消该订单吗？', '提示', {
      type: 'warning'
    })
    await cancelOrder(order.id)
    ElMessage.success('订单已取消')
    loadOrders()
  } catch (e) {
    if (e !== 'cancel') console.error(e)
  }
}

const handlePayOrder = async (order) => {
  try {
    await ElMessageBox.confirm('确认支付该订单？', '提示', {
      type: 'info'
    })
    await payOrder(order.id)
    ElMessage.success('支付成功')
    loadOrders()
  } catch (e) {
    if (e !== 'cancel') console.error(e)
  }
}

const handleConfirmReceive = async (order) => {
  try {
    await ElMessageBox.confirm('确认已收到货物？', '提示', {
      type: 'warning'
    })
    await confirmReceive(order.id)
    ElMessage.success('已确认收货')
    loadOrders()
  } catch (e) {
    if (e !== 'cancel') console.error(e)
  }
}

onMounted(() => {
  loadOrders()
})
</script>

<style scoped>
.orders-page {
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

.tabs {
  display: flex;
  background: #fff;
  padding: 0 20px;
  margin-bottom: 20px;
  border-radius: 4px;
}

.tabs > div {
  padding: 15px 25px;
  font-size: 15px;
  color: #666;
  cursor: pointer;
  border-bottom: 2px solid transparent;
}

.tabs > div.active {
  color: #409EFF;
  border-bottom-color: #409EFF;
  font-weight: 500;
}

.order-card {
  background: #fff;
  margin-bottom: 15px;
  border-radius: 4px;
  overflow: hidden;
}

.order-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px 20px;
  background: #fafafa;
  border-bottom: 1px solid #eee;
}

.order-info {
  display: flex;
  gap: 20px;
}

.order-no {
  color: #333;
}

.order-time {
  color: #999;
  font-size: 14px;
}

.order-items {
  padding: 15px 20px;
  cursor: pointer;
}

.order-item {
  display: flex;
  align-items: center;
  padding: 10px 0;
  border-bottom: 1px dashed #eee;
}

.order-item:last-child {
  border-bottom: none;
}

.item-img {
  width: 80px;
  height: 80px;
  border-radius: 4px;
  border: 1px solid #eee;
}

.item-info {
  margin-left: 15px;
  flex: 1;
}

.item-name {
  color: #333;
  line-height: 1.4;
}

.item-price {
  color: #999;
  margin-top: 8px;
}

.order-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px 20px;
  background: #fafafa;
  border-top: 1px solid #eee;
}

.order-amount {
  display: flex;
  gap: 30px;
  color: #666;
}

.amount {
  font-size: 15px;
}

.amount strong {
  color: #e4393c;
  font-size: 18px;
}

.order-actions {
  display: flex;
  gap: 10px;
  align-items: center;
}

.pagination-box {
  display: flex;
  justify-content: center;
  margin-top: 30px;
}
</style>
