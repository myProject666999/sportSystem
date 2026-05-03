<template>
  <div class="order-detail-page">
    <div class="container">
      <el-breadcrumb separator="/" class="breadcrumb">
        <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
        <el-breadcrumb-item :to="{ path: '/orders' }">我的订单</el-breadcrumb-item>
        <el-breadcrumb-item>订单详情</el-breadcrumb-item>
      </el-breadcrumb>

      <div class="order-status-card">
        <div class="status-info">
          <div class="status-icon">
            <el-icon :size="48"><Checked /></el-icon>
          </div>
          <div class="status-text">
            <h2>{{ getStatusText(order.status) }}</h2>
            <p v-if="order.status === 'shipped'">
              物流公司: {{ order.shipping_company }}
              <span style="margin: 0 10px;">|</span>
              物流单号: {{ order.tracking_number }}
            </p>
            <p v-else-if="order.status === 'pending'">
              请在 30 分钟内完成支付，超时订单将自动取消
            </p>
          </div>
          <div class="status-actions">
            <el-button v-if="order.status === 'pending'" type="danger" size="large" @click="handlePayOrder">立即支付</el-button>
            <el-button v-if="order.status === 'pending'" @click="handleCancelOrder">取消订单</el-button>
            <el-button v-if="order.status === 'shipped'" type="primary" size="large" @click="handleConfirmReceive">确认收货</el-button>
          </div>
        </div>
      </div>

      <el-row :gutter="20">
        <el-col :span="16">
          <el-card class="info-card">
            <template #header>
              <span class="card-title">订单信息</span>
            </template>
            <el-descriptions :column="2" border>
              <el-descriptions-item label="订单号">{{ order.order_no }}</el-descriptions-item>
              <el-descriptions-item label="订单状态">
                <el-tag :type="getStatusType(order.status)">{{ getStatusText(order.status) }}</el-tag>
              </el-descriptions-item>
              <el-descriptions-item label="创建时间">{{ order.created_at }}</el-descriptions-item>
              <el-descriptions-item label="支付时间">{{ order.paid_at || '未支付' }}</el-descriptions-item>
              <el-descriptions-item label="物流公司">{{ order.shipping_company || '-' }}</el-descriptions-item>
              <el-descriptions-item label="物流单号">{{ order.tracking_number || '-' }}</el-descriptions-item>
            </el-descriptions>
          </el-card>

          <el-card class="info-card" style="margin-top: 20px;">
            <template #header>
              <span class="card-title">收货地址</span>
            </template>
            <div class="address-info">
              <div class="receiver">
                <span class="name">{{ order.address_name }}</span>
                <span class="phone">{{ order.address_phone }}</span>
              </div>
              <div class="address">{{ order.address_info }}</div>
            </div>
          </el-card>

          <el-card class="info-card" style="margin-top: 20px;">
            <template #header>
              <span class="card-title">商品清单</span>
            </template>
            <el-table :data="order.items || []" stripe>
              <el-table-column label="商品信息" min-width="300">
                <template #default="{ row }">
                  <div class="product-info">
                    <el-image :src="row.product_image" fit="cover" class="product-img" />
                    <div class="product-detail">
                      <div class="product-name">{{ row.product_name }}</div>
                    </div>
                  </div>
                </template>
              </el-table-column>
              <el-table-column prop="price" label="单价" width="120">
                <template #default="{ row }">¥{{ row.price }}</template>
              </el-table-column>
              <el-table-column prop="quantity" label="数量" width="100" />
              <el-table-column label="小计" width="120">
                <template #default="{ row }">
                  <span class="subtotal">¥{{ (row.price * row.quantity).toFixed(2) }}</span>
                </template>
              </el-table-column>
            </el-table>
          </el-card>
        </el-col>

        <el-col :span="8">
          <el-card class="summary-card">
            <template #header>
              <span class="card-title">订单汇总</span>
            </template>
            <div class="summary-list">
              <div class="summary-item">
                <span class="label">商品总价</span>
                <span class="value">¥{{ order.total_amount }}</span>
              </div>
              <div class="summary-item">
                <span class="label">运费</span>
                <span class="value">¥0.00</span>
              </div>
              <div class="summary-item total">
                <span class="label">应付金额</span>
                <span class="value">¥{{ order.total_amount }}</span>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Checked } from '@element-plus/icons-vue'
import { getOrderDetail, payOrder, cancelOrder, confirmReceive } from '@/api'

const route = useRoute()
const router = useRouter()
const order = ref({})

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

const loadOrder = async () => {
  const token = localStorage.getItem('token')
  if (!token) {
    ElMessage.warning('请先登录')
    router.push('/login')
    return
  }
  try {
    const res = await getOrderDetail(route.params.id)
    order.value = res.data || {}
  } catch (e) {
    console.error(e)
  }
}

const handlePayOrder = async () => {
  try {
    await ElMessageBox.confirm('确认支付该订单？', '提示', {
      type: 'info'
    })
    await payOrder(route.params.id)
    ElMessage.success('支付成功')
    loadOrder()
  } catch (e) {
    if (e !== 'cancel') console.error(e)
  }
}

const handleCancelOrder = async () => {
  try {
    await ElMessageBox.confirm('确定要取消该订单吗？', '提示', {
      type: 'warning'
    })
    await cancelOrder(route.params.id)
    ElMessage.success('订单已取消')
    loadOrder()
  } catch (e) {
    if (e !== 'cancel') console.error(e)
  }
}

const handleConfirmReceive = async () => {
  try {
    await ElMessageBox.confirm('确认已收到货物？', '提示', {
      type: 'warning'
    })
    await confirmReceive(route.params.id)
    ElMessage.success('已确认收货')
    loadOrder()
  } catch (e) {
    if (e !== 'cancel') console.error(e)
  }
}

onMounted(() => {
  loadOrder()
})
</script>

<style scoped>
.order-detail-page {
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
  margin-bottom: 20px;
}

.order-status-card {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 8px;
  padding: 30px;
  margin-bottom: 20px;
}

.status-info {
  display: flex;
  align-items: center;
  justify-content: space-between;
  color: #fff;
}

.status-icon {
  background: rgba(255, 255, 255, 0.2);
  width: 80px;
  height: 80px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.status-text {
  flex: 1;
  margin-left: 30px;
}

.status-text h2 {
  margin: 0 0 10px 0;
  font-size: 24px;
}

.status-text p {
  margin: 0;
  opacity: 0.9;
}

.status-actions {
  display: flex;
  gap: 10px;
}

.info-card {
  border-radius: 8px;
}

.card-title {
  font-weight: bold;
  font-size: 16px;
}

.address-info {
  padding: 10px 0;
}

.receiver {
  margin-bottom: 10px;
}

.receiver .name {
  font-weight: bold;
  margin-right: 20px;
}

.receiver .phone {
  color: #666;
}

.address {
  color: #666;
}

.product-info {
  display: flex;
  align-items: center;
}

.product-img {
  width: 70px;
  height: 70px;
  border-radius: 4px;
  border: 1px solid #eee;
}

.product-detail {
  margin-left: 15px;
}

.product-name {
  color: #333;
  line-height: 1.4;
}

.subtotal {
  color: #e4393c;
  font-weight: bold;
}

.summary-card {
  border-radius: 8px;
}

.summary-list {
  padding: 10px 0;
}

.summary-item {
  display: flex;
  justify-content: space-between;
  padding: 12px 0;
  border-bottom: 1px solid #f0f0f0;
}

.summary-item.total {
  border-bottom: none;
  padding-top: 15px;
  margin-top: 5px;
}

.summary-item.total .label {
  font-size: 16px;
  font-weight: bold;
}

.summary-item.total .value {
  color: #e4393c;
  font-size: 20px;
  font-weight: bold;
}

.summary-item .label {
  color: #666;
}
</style>
