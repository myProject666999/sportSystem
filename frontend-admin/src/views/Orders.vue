<template>
  <div class="orders-page">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>订单管理</span>
          <div class="search-box">
            <el-input
              v-model="searchQuery"
              placeholder="搜索订单号/用户名"
              clearable
              @keyup.enter="handleSearch"
              style="width: 200px; margin-right: 10px"
            />
            <el-select v-model="status" placeholder="订单状态" clearable style="width: 130px; margin-right: 10px">
              <el-option label="待支付" value="pending" />
              <el-option label="已支付" value="paid" />
              <el-option label="已发货" value="shipped" />
              <el-option label="已完成" value="completed" />
              <el-option label="已取消" value="cancelled" />
            </el-select>
            <el-button type="primary" :icon="Search" @click="handleSearch">搜索</el-button>
          </div>
        </div>
      </template>

      <el-table :data="orders" stripe v-loading="loading">
        <el-table-column prop="order_no" label="订单号" width="200" />
        <el-table-column prop="username" label="用户" width="120" />
        <el-table-column prop="total_amount" label="订单金额" width="120">
          <template #default="{ row }">
            <span class="price">¥{{ row.total_amount }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="订单状态" width="120">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ getStatusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="shipping_company" label="物流公司" width="120">
          <template #default="{ row }">
            {{ row.shipping_company || '暂无' }}
          </template>
        </el-table-column>
        <el-table-column prop="tracking_number" label="物流单号" width="180">
          <template #default="{ row }">
            {{ row.tracking_number || '暂无' }}
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="250" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="viewDetail(row)">详情</el-button>
            <el-button
              v-if="row.status === 'paid'"
              type="success"
              link
              @click="handleShip(row)"
            >发货</el-button>
            <el-button
              v-if="row.status === 'shipped'"
              type="warning"
              link
              @click="handleUpdateLogistics(row)"
            >修改物流</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-box">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="loadData"
          @current-change="loadData"
        />
      </div>
    </el-card>

    <el-dialog v-model="detailVisible" title="订单详情" width="800px">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="订单号">{{ currentOrder.order_no }}</el-descriptions-item>
        <el-descriptions-item label="订单状态">
          <el-tag :type="getStatusType(currentOrder.status)">{{ getStatusText(currentOrder.status) }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="用户">{{ currentOrder.username || '-' }}</el-descriptions-item>
        <el-descriptions-item label="订单金额">
          <span class="price">¥{{ currentOrder.total_amount }}</span>
        </el-descriptions-item>
        <el-descriptions-item label="物流公司">{{ currentOrder.shipping_company || '-' }}</el-descriptions-item>
        <el-descriptions-item label="物流单号">{{ currentOrder.tracking_number || '-' }}</el-descriptions-item>
        <el-descriptions-item label="收货地址" :span="2">
          {{ currentOrder.address_info || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ currentOrder.created_at }}</el-descriptions-item>
        <el-descriptions-item label="支付时间">{{ currentOrder.paid_at || '-' }}</el-descriptions-item>
      </el-descriptions>

      <div style="margin-top: 20px;">
        <h4 style="margin-bottom: 15px;">订单商品</h4>
        <el-table :data="currentOrder.items || []" border size="small">
          <el-table-column prop="product_name" label="商品名称" />
          <el-table-column prop="quantity" label="数量" width="80" />
          <el-table-column prop="price" label="单价" width="100">
            <template #default="{ row }">¥{{ row.price }}</template>
          </el-table-column>
          <el-table-column prop="subtotal" label="小计" width="100">
            <template #default="{ row }">¥{{ row.subtotal }}</template>
          </el-table-column>
        </el-table>
      </div>
    </el-dialog>

    <el-dialog v-model="shipVisible" title="发货" width="500px">
      <el-form :model="shipForm" :rules="shipRules" ref="shipFormRef" label-width="100px">
        <el-form-item label="物流公司" prop="shipping_company">
          <el-select v-model="shipForm.shipping_company" placeholder="请选择物流公司" style="width: 100%">
            <el-option label="顺丰速运" value="顺丰速运" />
            <el-option label="中通快递" value="中通快递" />
            <el-option label="圆通速递" value="圆通速递" />
            <el-option label="申通快递" value="申通快递" />
            <el-option label="韵达速递" value="韵达速递" />
            <el-option label="邮政EMS" value="邮政EMS" />
          </el-select>
        </el-form-item>
        <el-form-item label="物流单号" prop="tracking_number">
          <el-input v-model="shipForm.tracking_number" placeholder="请输入物流单号" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="shipVisible = false">取消</el-button>
        <el-button type="primary" @click="submitShip" :loading="shipLoading">确认发货</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Search } from '@element-plus/icons-vue'
import { getOrders, getOrderDetail, shipOrder, updateLogistics } from '@/api'

const loading = ref(false)
const orders = ref([])
const searchQuery = ref('')
const status = ref(null)
const detailVisible = ref(false)
const shipVisible = ref(false)
const currentOrder = ref({})
const shipFormRef = ref(null)
const shipLoading = ref(false)

const pagination = ref({
  page: 1,
  pageSize: 10,
  total: 0
})

const shipForm = reactive({
  order_id: null,
  shipping_company: '',
  tracking_number: ''
})

const shipRules = {
  shipping_company: [
    { required: true, message: '请选择物流公司', trigger: 'change' }
  ],
  tracking_number: [
    { required: true, message: '请输入物流单号', trigger: 'blur' }
  ]
}

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

const loadData = async () => {
  loading.value = true
  try {
    const res = await getOrders({
      page: pagination.value.page,
      page_size: pagination.value.pageSize,
      keyword: searchQuery.value,
      status: status.value
    })
    orders.value = res.data?.list || []
    pagination.value.total = res.data?.total || 0
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.value.page = 1
  loadData()
}

const viewDetail = async (row) => {
  try {
    const res = await getOrderDetail(row.id)
    currentOrder.value = res.data || row
    detailVisible.value = true
  } catch (e) {
    console.error(e)
    currentOrder.value = row
    detailVisible.value = true
  }
}

const handleShip = (row) => {
  shipForm.order_id = row.id
  shipForm.shipping_company = ''
  shipForm.tracking_number = ''
  shipVisible.value = true
}

const handleUpdateLogistics = (row) => {
  shipForm.order_id = row.id
  shipForm.shipping_company = row.shipping_company || ''
  shipForm.tracking_number = row.tracking_number || ''
  shipVisible.value = true
}

const submitShip = async () => {
  const valid = await shipFormRef.value.validate().catch(() => false)
  if (!valid) return

  shipLoading.value = true
  try {
    const data = {
      shipping_company: shipForm.shipping_company,
      tracking_number: shipForm.tracking_number
    }
    
    const order = orders.value.find(o => o.id === shipForm.order_id)
    if (order && order.status === 'paid') {
      await shipOrder(shipForm.order_id, data)
    } else {
      await updateLogistics(shipForm.order_id, data)
    }
    
    ElMessage.success('操作成功')
    shipVisible.value = false
    loadData()
  } catch (e) {
    console.error(e)
  } finally {
    shipLoading.value = false
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.search-box {
  display: flex;
  align-items: center;
}

.pagination-box {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.price {
  color: #F56C6C;
  font-weight: bold;
}
</style>
