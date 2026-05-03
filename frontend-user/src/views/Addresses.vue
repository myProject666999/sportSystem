<template>
  <div class="addresses-page">
    <div class="container">
      <div class="page-header">
        <h2 class="page-title">地址管理</h2>
        <el-button type="primary" @click="handleAdd">
          <el-icon><Plus /></el-icon>
          新增地址
        </el-button>
      </div>

      <div class="address-list" v-loading="loading">
        <div v-if="addresses.length === 0">
          <el-empty description="暂无收货地址">
            <el-button type="primary" @click="handleAdd">新增地址</el-button>
          </el-empty>
        </div>
        <div v-else>
          <div v-for="address in addresses" :key="address.id" class="address-card">
            <div class="address-info">
              <div class="address-header">
                <span class="receiver">
                  <span class="name">{{ address.name }}</span>
                  <span class="phone">{{ address.phone }}</span>
                </span>
                <el-tag v-if="address.is_default" type="success" size="small">默认地址</el-tag>
              </div>
              <div class="address-detail">{{ address.province }}{{ address.city }}{{ address.district }}{{ address.detail }}</div>
            </div>
            <div class="address-actions">
              <el-button type="text" @click="handleEdit(address)">编辑</el-button>
              <el-button type="text" v-if="!address.is_default" @click="setDefault(address)">设为默认</el-button>
              <el-button type="text" class="delete-btn" @click="handleDelete(address)">删除</el-button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑地址' : '新增地址'" width="600px">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="收货人" prop="name">
          <el-input v-model="form.name" placeholder="请输入收货人姓名" />
        </el-form-item>
        <el-form-item label="手机号码" prop="phone">
          <el-input v-model="form.phone" placeholder="请输入手机号码" />
        </el-form-item>
        <el-form-item label="所在地区" prop="province">
          <el-cascader
            v-model="region"
            :options="regionOptions"
            :props="{ value: 'label', label: 'label', children: 'children' }"
            placeholder="请选择省/市/区"
            style="width: 100%;"
            @change="handleRegionChange"
          />
        </el-form-item>
        <el-form-item label="详细地址" prop="detail">
          <el-input v-model="form.detail" type="textarea" :rows="2" placeholder="请输入详细地址" />
        </el-form-item>
        <el-form-item label="设为默认">
          <el-switch v-model="form.is_default" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { getAddresses, createAddress, updateAddress, deleteAddress } from '@/api'

const router = useRouter()
const loading = ref(false)
const submitLoading = ref(false)
const addresses = ref([])
const dialogVisible = ref(false)
const isEdit = ref(false)
const formRef = ref(null)
const region = ref([])

const regionOptions = [
  {
    label: '北京市',
    children: [
      { label: '北京市', children: [{ label: '东城区' }, { label: '西城区' }, { label: '朝阳区' }, { label: '海淀区' }] }
    ]
  },
  {
    label: '上海市',
    children: [
      { label: '上海市', children: [{ label: '黄浦区' }, { label: '徐汇区' }, { label: '静安区' }, { label: '浦东新区' }] }
    ]
  },
  {
    label: '广东省',
    children: [
      { label: '广州市', children: [{ label: '天河区' }, { label: '越秀区' }, { label: '海珠区' }] },
      { label: '深圳市', children: [{ label: '南山区' }, { label: '福田区' }, { label: '罗湖区' }] }
    ]
  }
]

const form = reactive({
  id: null,
  name: '',
  phone: '',
  province: '',
  city: '',
  district: '',
  detail: '',
  is_default: false
})

const rules = {
  name: [{ required: true, message: '请输入收货人姓名', trigger: 'blur' }],
  phone: [
    { required: true, message: '请输入手机号码', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号码', trigger: 'blur' }
  ],
  detail: [{ required: true, message: '请输入详细地址', trigger: 'blur' }]
}

const loadAddresses = async () => {
  const token = localStorage.getItem('token')
  if (!token) {
    ElMessage.warning('请先登录')
    router.push('/login')
    return
  }
  loading.value = true
  try {
    const res = await getAddresses()
    addresses.value = res.data || []
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const handleRegionChange = (val) => {
  if (val && val.length === 3) {
    form.province = val[0]
    form.city = val[1]
    form.district = val[2]
  }
}

const handleAdd = () => {
  isEdit.value = false
  form.id = null
  form.name = ''
  form.phone = ''
  form.province = ''
  form.city = ''
  form.district = ''
  form.detail = ''
  form.is_default = false
  region.value = []
  dialogVisible.value = true
}

const handleEdit = (address) => {
  isEdit.value = true
  form.id = address.id
  form.name = address.name
  form.phone = address.phone
  form.province = address.province
  form.city = address.city
  form.district = address.district
  form.detail = address.detail
  form.is_default = address.is_default
  region.value = [address.province, address.city, address.district]
  dialogVisible.value = true
}

const handleSubmit = async () => {
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  submitLoading.value = true
  try {
    const data = {
      name: form.name,
      phone: form.phone,
      province: form.province,
      city: form.city,
      district: form.district,
      detail: form.detail,
      is_default: form.is_default
    }

    if (isEdit.value) {
      await updateAddress(form.id, data)
    } else {
      await createAddress(data)
    }
    ElMessage.success(isEdit.value ? '编辑成功' : '添加成功')
    dialogVisible.value = false
    loadAddresses()
  } catch (e) {
    console.error(e)
  } finally {
    submitLoading.value = false
  }
}

const setDefault = async (address) => {
  try {
    await updateAddress(address.id, { is_default: true })
    ElMessage.success('已设为默认地址')
    loadAddresses()
  } catch (e) {
    console.error(e)
  }
}

const handleDelete = async (address) => {
  try {
    await ElMessageBox.confirm('确定要删除该地址吗？', '提示', {
      type: 'warning'
    })
    await deleteAddress(address.id)
    ElMessage.success('删除成功')
    loadAddresses()
  } catch (e) {
    if (e !== 'cancel') console.error(e)
  }
}

onMounted(() => {
  loadAddresses()
})
</script>

<style scoped>
.addresses-page {
  padding: 30px 0;
  background: #f5f5f5;
  min-height: calc(100vh - 140px);
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 15px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.page-title {
  font-size: 22px;
  color: #333;
  margin: 0;
}

.address-list {
  background: #fff;
  border-radius: 8px;
  padding: 20px;
}

.address-card {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  padding: 20px;
  border: 1px solid #eee;
  border-radius: 4px;
  margin-bottom: 15px;
}

.address-card:hover {
  border-color: #409EFF;
}

.address-info {
  flex: 1;
}

.address-header {
  display: flex;
  align-items: center;
  margin-bottom: 10px;
}

.receiver {
  display: flex;
  align-items: center;
  margin-right: 15px;
}

.receiver .name {
  font-weight: bold;
  margin-right: 15px;
  font-size: 16px;
}

.receiver .phone {
  color: #666;
}

.address-detail {
  color: #666;
  line-height: 1.6;
}

.address-actions {
  display: flex;
  gap: 10px;
}

.delete-btn {
  color: #F56C6C !important;
}
</style>
