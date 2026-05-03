<template>
  <div class="recharge-page">
    <div class="container">
      <h2 class="page-title">账户充值</h2>

      <div class="recharge-container">
        <div class="recharge-info">
          <div class="balance-card">
            <div class="balance-info">
              <div class="balance-label">当前余额</div>
              <div class="balance-value">¥{{ userInfo.balance || 0 }}</div>
            </div>
            <div class="points-info">
              <div class="points-label">当前积分</div>
              <div class="points-value">{{ userInfo.points || 0 }} 积分</div>
            </div>
          </div>
        </div>

        <el-card class="recharge-form-card">
          <template #header>
            <span class="card-title">选择充值金额</span>
          </template>

          <div class="amount-options">
            <div
              v-for="amount in amountOptions"
              :key="amount"
              :class="{ active: selectedAmount === amount }"
              class="amount-item"
              @click="selectedAmount = amount"
            >
              <div class="amount-value">¥{{ amount }}</div>
              <div v-if="amount >= 100" class="amount-bonus">送 {{ Math.floor(amount / 100) * 10 }} 积分</div>
            </div>
          </div>

          <div class="custom-amount">
            <span class="custom-label">自定义金额：</span>
            <el-input-number
              v-model="customAmount"
              :min="1"
              :precision="0"
              placeholder="请输入充值金额"
              @change="selectedAmount = 0"
            />
            <span class="custom-unit">元</span>
          </div>

          <div class="payment-methods">
            <div class="method-label">支付方式：</div>
            <div class="method-list">
              <div
                v-for="method in paymentMethods"
                :key="method.value"
                :class="{ active: selectedPayment === method.value }"
                class="method-item"
                @click="selectedPayment = method.value"
              >
                <el-icon :size="24" class="method-icon">{{ method.icon }}</el-icon>
                <span class="method-name">{{ method.name }}</span>
              </div>
            </div>
          </div>

          <div class="recharge-summary">
            <div class="summary-info">
              <span>充值金额：</span>
              <span class="summary-amount">¥{{ rechargeAmount }}</span>
            </div>
            <el-button type="primary" size="large" class="recharge-btn" @click="handleRecharge" :loading="loading">
              立即充值
            </el-button>
          </div>
        </el-card>

        <el-card class="recharge-record" style="margin-top: 20px;">
          <template #header>
            <span class="card-title">充值记录</span>
          </template>
          <el-table :data="rechargeRecords" stripe v-loading="recordLoading">
            <el-table-column prop="order_no" label="订单号" width="250" />
            <el-table-column prop="amount" label="充值金额" width="120">
              <template #default="{ row }">
                <span class="amount">¥{{ row.amount }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="payment_method" label="支付方式" width="120" />
            <el-table-column prop="status" label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="row.status === 'success' ? 'success' : 'warning'">
                  {{ row.status === 'success' ? '成功' : '处理中' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="created_at" label="时间" width="180" />
          </el-table>
        </el-card>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Wallet, Money } from '@element-plus/icons-vue'
import { getUserInfo, recharge } from '@/api'

const router = useRouter()
const loading = ref(false)
const recordLoading = ref(false)
const userInfo = ref({})
const selectedAmount = ref(100)
const customAmount = ref(0)
const selectedPayment = ref('alipay')
const rechargeRecords = ref([])

const amountOptions = [50, 100, 200, 500, 1000, 2000]

const paymentMethods = [
  { value: 'alipay', name: '支付宝', icon: Wallet },
  { value: 'wechat', name: '微信支付', icon: Money }
]

const rechargeAmount = computed(() => {
  if (selectedAmount.value > 0) {
    return selectedAmount.value
  }
  return customAmount.value || 0
})

const loadUserInfo = async () => {
  const token = localStorage.getItem('token')
  if (!token) {
    ElMessage.warning('请先登录')
    router.push('/login')
    return
  }
  try {
    const res = await getUserInfo()
    userInfo.value = res.data || {}
  } catch (e) {
    console.error(e)
  }
}

const handleRecharge = async () => {
  if (rechargeAmount.value <= 0) {
    ElMessage.warning('请选择充值金额')
    return
  }

  loading.value = true
  try {
    await recharge({
      amount: rechargeAmount.value,
      payment_method: selectedPayment.value
    })
    ElMessage.success('充值成功')
    loadUserInfo()
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadUserInfo()
})
</script>

<style scoped>
.recharge-page {
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

.balance-card {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 12px;
  padding: 30px;
  margin-bottom: 20px;
  display: flex;
  justify-content: space-between;
  color: #fff;
}

.balance-label,
.points-label {
  font-size: 14px;
  opacity: 0.9;
  margin-bottom: 10px;
}

.balance-value {
  font-size: 32px;
  font-weight: bold;
}

.points-value {
  font-size: 24px;
  font-weight: bold;
  color: #ffd700;
}

.card-title {
  font-weight: bold;
  font-size: 16px;
}

.amount-options {
  display: flex;
  flex-wrap: wrap;
  gap: 15px;
  margin-bottom: 25px;
}

.amount-item {
  width: calc(16.66% - 13px);
  padding: 20px 10px;
  border: 2px solid #e4e7ed;
  border-radius: 8px;
  text-align: center;
  cursor: pointer;
  transition: all 0.3s;
}

.amount-item:hover {
  border-color: #409EFF;
}

.amount-item.active {
  border-color: #409EFF;
  background: #ecf5ff;
}

.amount-value {
  font-size: 20px;
  font-weight: bold;
  color: #e4393c;
}

.amount-bonus {
  font-size: 12px;
  color: #67c23a;
  margin-top: 5px;
}

.custom-amount {
  display: flex;
  align-items: center;
  margin-bottom: 25px;
  padding: 15px;
  background: #fafafa;
  border-radius: 8px;
}

.custom-label {
  margin-right: 15px;
  color: #666;
}

.custom-unit {
  margin-left: 10px;
  color: #333;
}

.payment-methods {
  margin-bottom: 25px;
}

.method-label {
  color: #666;
  margin-bottom: 15px;
}

.method-list {
  display: flex;
  gap: 20px;
}

.method-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 15px 30px;
  border: 2px solid #e4e7ed;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s;
}

.method-item:hover {
  border-color: #409EFF;
}

.method-item.active {
  border-color: #409EFF;
  background: #ecf5ff;
}

.method-icon {
  color: #409EFF;
}

.recharge-summary {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  background: #fafafa;
  border-radius: 8px;
}

.summary-info {
  font-size: 16px;
}

.summary-amount {
  font-size: 24px;
  font-weight: bold;
  color: #e4393c;
}

.recharge-btn {
  padding: 0 50px;
  font-size: 16px;
  height: 46px;
}

.amount {
  color: #e4393c;
  font-weight: bold;
}
</style>
