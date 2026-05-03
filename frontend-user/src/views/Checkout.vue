<template>
  <div class="checkout-page">
    <div class="container">
      <h2 class="page-title">确认订单</h2>

      <div class="checkout-steps">
        <div class="step-item active">
          <span class="step-number">1</span>
          <span class="step-name">确认订单</span>
        </div>
        <div class="step-divider"></div>
        <div class="step-item">
          <span class="step-number">2</span>
          <span class="step-name">支付</span>
        </div>
        <div class="step-divider"></div>
        <div class="step-item">
          <span class="step-number">3</span>
          <span class="step-name">完成</span>
        </div>
      </div>

      <div class="checkout-content" v-loading="loading">
        <el-card class="section-card" shadow="never">
          <template #header>
            <span class="section-title">收货地址</span>
          </template>
          <div v-if="addresses.length === 0" class="no-address">
            <el-empty description="暂无收货地址" :image-size="60">
              <el-button type="primary" @click="$router.push('/addresses')">添加地址</el-button>
            </el-empty>
          </div>
          <div v-else>
            <div class="address-list">
              <div
                v-for="address in addresses"
                :key="address.id"
                :class="{ active: selectedAddressId === address.id }"
                class="address-item"
                @click="selectedAddressId = address.id"
              >
                <div class="address-header">
                  <span class="receiver">{{ address.receiver_name }}</span>
                  <span class="phone">{{ address.receiver_phone }}</span>
                  <el-tag v-if="address.is_default" type="success" size="small">默认</el-tag>
                </div>
                <div class="address-detail">
                  {{ address.province }}{{ address.city }}{{ address.district }}{{ address.detail }}
                </div>
              </div>
            </div>
          </div>
        </el-card>

        <el-card class="section-card" shadow="never" style="margin-top: 20px;">
          <template #header>
            <span class="section-title">商品清单</span>
          </template>
          <el-table :data="cartItems" style="width: 100%">
            <el-table-column prop="product_name" label="商品" width="400">
              <template #default="{ row }">
                <div class="product-info">
                  <el-image :src="row.product?.image" fit="cover" class="product-img" />
                  <div class="product-text">
                    <div class="product-name">{{ row.product?.name }}</div>
                  </div>
                </div>
              </template>
            </el-table-column>
            <el-table-column prop="price" label="单价" width="150">
              <template #default="{ row }">
                <span class="price">¥{{ row.product?.price }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="quantity" label="数量" width="120" align="center">
              <template #default="{ row }">
                <span>{{ row.quantity }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="subtotal" label="小计" width="150">
              <template #default="{ row }">
                <span class="subtotal">¥{{ (row.product?.price * row.quantity).toFixed(2) }}</span>
              </template>
            </el-table-column>
          </el-table>
        </el-card>

        <el-card class="section-card" shadow="never" style="margin-top: 20px;">
          <template #header>
            <span class="section-title">订单备注</span>
          </template>
          <el-input
            v-model="orderRemark"
            type="textarea"
            :rows="2"
            placeholder="选填，有什么特殊需求可以告诉卖家（限100字）"
            maxlength="100"
            show-word-limit
          />
        </el-card>

        <el-card class="summary-card" shadow="never" style="margin-top: 20px;">
          <div class="summary-content">
            <div class="summary-row">
              <span class="label">商品金额：</span>
              <span class="value">¥{{ totalAmount.toFixed(2) }}</span>
            </div>
            <div class="summary-row">
              <span class="label">运费：</span>
              <span class="value">¥{{ shippingFee.toFixed(2) }}</span>
            </div>
            <div class="summary-row total-row">
              <span class="label">应付金额：</span>
              <span class="total-value">¥{{ (totalAmount + shippingFee).toFixed(2) }}</span>
            </div>
          </div>
          <div class="summary-actions">
            <el-button @click="$router.back()">返回</el-button>
            <el-button type="primary" size="large" @click="submitOrder" :loading="submitting">
              提交订单
            </el-button>
          </div>
        </el-card>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { getCartList, getAddresses, createOrder } from '@/api'

const router = useRouter()
const route = useRoute()
const loading = ref(false)
const submitting = ref(false)
const cartItems = ref([])
const addresses = ref([])
const selectedAddressId = ref(null)
const orderRemark = ref('')
const shippingFee = ref(0)

const totalAmount = computed(() => {
  return cartItems.value.reduce((sum, item) => {
    return sum + (item.product?.price || 0) * item.quantity
  }, 0)
})

const loadData = async () => {
  loading.value = true
  const token = localStorage.getItem('token')
  if (!token) {
    ElMessage.warning('请先登录')
    router.push('/login')
    return
  }
  try {
    const cartRes = await getCartList()
    cartItems.value = cartRes.data?.list || cartRes.data || []

    const addressRes = await getAddresses()
    addresses.value = addressRes.data || []

    const defaultAddress = addresses.value.find(a => a.is_default)
    if (defaultAddress) {
      selectedAddressId.value = defaultAddress.id
    } else if (addresses.value.length > 0) {
      selectedAddressId.value = addresses.value[0].id
    }
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const submitOrder = async () => {
  if (!selectedAddressId.value) {
    ElMessage.warning('请选择收货地址')
    return
  }

  if (cartItems.value.length === 0) {
    ElMessage.warning('购物车为空')
    return
  }

  submitting.value = true
  try {
    const res = await createOrder({
      address_id: selectedAddressId.value,
      remark: orderRemark.value,
      use_points: false
    })
    ElMessage.success('订单创建成功')
    router.push(`/order/${res.data?.id}`)
  } catch (e) {
    console.error(e)
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.checkout-page {
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

.checkout-steps {
  display: flex;
  align-items: center;
  justify-content: center;
  background: #fff;
  padding: 25px;
  border-radius: 8px;
  margin-bottom: 20px;
}

.step-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.step-number {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: #e4e7ed;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #909399;
  font-weight: bold;
}

.step-name {
  color: #909399;
  font-size: 14px;
}

.step-item.active .step-number {
  background: #409EFF;
  color: #fff;
}

.step-item.active .step-name {
  color: #409EFF;
}

.step-divider {
  width: 80px;
  height: 2px;
  background: #e4e7ed;
  margin: 0 30px;
}

.section-card {
  border-radius: 8px;
}

.section-title {
  font-weight: bold;
  font-size: 16px;
}

.address-list {
  display: flex;
  flex-wrap: wrap;
  gap: 15px;
}

.address-item {
  width: calc(50% - 10px);
  padding: 15px;
  border: 2px solid #e4e7ed;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s;
}

.address-item:hover {
  border-color: #b3d8ff;
}

.address-item.active {
  border-color: #409EFF;
  background: #ecf5ff;
}

.address-header {
  display: flex;
  align-items: center;
  gap: 15px;
  margin-bottom: 10px;
}

.receiver {
  font-weight: bold;
  color: #333;
}

.phone {
  color: #666;
}

.address-detail {
  color: #666;
  font-size: 14px;
  line-height: 1.5;
}

.no-address {
  padding: 20px 0;
}

.product-info {
  display: flex;
  align-items: center;
  gap: 15px;
}

.product-img {
  width: 80px;
  height: 80px;
  border-radius: 4px;
}

.product-text {
  flex: 1;
}

.product-name {
  color: #333;
  line-height: 1.5;
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

.summary-card {
  border-radius: 8px;
}

.summary-content {
  padding: 20px 0;
  border-bottom: 1px solid #eee;
}

.summary-row {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  margin-bottom: 15px;
}

.summary-row .label {
  color: #666;
  margin-right: 30px;
}

.summary-row .value {
  color: #333;
  font-size: 14px;
}

.total-row {
  margin-bottom: 0;
}

.total-row .label {
  font-size: 16px;
}

.total-value {
  color: #e4393c;
  font-size: 24px;
  font-weight: bold;
}

.summary-actions {
  display: flex;
  justify-content: flex-end;
  gap: 15px;
  padding-top: 20px;
}
</style>
