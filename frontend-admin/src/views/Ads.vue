<template>
  <div class="ads-page">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>广告管理</span>
          <el-button type="primary" :icon="Plus" @click="handleAdd">新增广告</el-button>
        </div>
      </template>

      <el-table :data="ads" stripe v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="图片" width="150">
          <template #default="{ row }">
            <el-image :src="row.image" fit="cover" style="width: 120px; height: 60px; border-radius: 4px" />
          </template>
        </el-table-column>
        <el-table-column prop="title" label="标题" min-width="200" />
        <el-table-column prop="position" label="位置" width="120">
          <template #default="{ row }">
            <el-tag :type="getPositionType(row.position)">{{ getPositionText(row.position) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="sort" label="排序" width="100" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-button type="danger" link @click="handleDelete(row)">删除</el-button>
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

    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑广告' : '新增广告'" width="600px">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="标题" prop="title">
          <el-input v-model="form.title" placeholder="请输入广告标题" />
        </el-form-item>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="位置" prop="position">
              <el-select v-model="form.position" placeholder="请选择位置" style="width: 100%">
                <el-option label="首页顶部" value="top" />
                <el-option label="首页侧边" value="sidebar" />
                <el-option label="首页底部" value="bottom" />
                <el-option label="商品列表页" value="product" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="排序" prop="sort">
              <el-input-number v-model="form.sort" :min="0" :max="999" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="图片" prop="image">
          <el-upload
            v-model:file-list="fileList"
            action="#"
            list-type="picture-card"
            :auto-upload="false"
            :limit="1"
          >
            <el-icon><Plus /></el-icon>
          </el-upload>
        </el-form-item>
        <el-form-item label="跳转链接" prop="link">
          <el-input v-model="form.link" placeholder="请输入跳转链接(可选)" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-switch v-model="form.status" :active-value="1" :inactive-value="0" active-text="启用" inactive-text="禁用" />
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
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { getAds, createAd, updateAd, deleteAd } from '@/api'

const loading = ref(false)
const submitLoading = ref(false)
const ads = ref([])
const dialogVisible = ref(false)
const isEdit = ref(false)
const formRef = ref(null)
const fileList = ref([])

const pagination = ref({
  page: 1,
  pageSize: 10,
  total: 0
})

const form = reactive({
  id: null,
  title: '',
  image: '',
  link: '',
  position: 'top',
  sort: 0,
  status: 1
})

const rules = {
  title: [
    { required: true, message: '请输入广告标题', trigger: 'blur' }
  ],
  position: [
    { required: true, message: '请选择位置', trigger: 'change' }
  ]
}

const getPositionType = (position) => {
  const types = {
    top: 'primary',
    sidebar: 'success',
    bottom: 'warning',
    product: 'info'
  }
  return types[position] || 'info'
}

const getPositionText = (position) => {
  const texts = {
    top: '首页顶部',
    sidebar: '首页侧边',
    bottom: '首页底部',
    product: '商品列表页'
  }
  return texts[position] || position
}

const loadData = async () => {
  loading.value = true
  try {
    const res = await getAds({
      page: pagination.value.page,
      page_size: pagination.value.pageSize
    })
    ads.value = res.data?.list || []
    pagination.value.total = res.data?.total || 0
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const handleAdd = () => {
  isEdit.value = false
  form.id = null
  form.title = ''
  form.image = ''
  form.link = ''
  form.position = 'top'
  form.sort = 0
  form.status = 1
  fileList.value = []
  dialogVisible.value = true
}

const handleEdit = (row) => {
  isEdit.value = true
  form.id = row.id
  form.title = row.title
  form.image = row.image
  form.link = row.link
  form.position = row.position
  form.sort = row.sort
  form.status = row.status
  fileList.value = row.image ? [{ url: row.image }] : []
  dialogVisible.value = true
}

const handleSubmit = async () => {
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  submitLoading.value = true
  try {
    const data = {
      title: form.title,
      image: form.image,
      link: form.link,
      position: form.position,
      sort: form.sort,
      status: form.status
    }
    
    if (isEdit.value) {
      await updateAd(form.id, data)
    } else {
      await createAd(data)
    }
    ElMessage.success(isEdit.value ? '编辑成功' : '添加成功')
    dialogVisible.value = false
    loadData()
  } catch (e) {
    console.error(e)
  } finally {
    submitLoading.value = false
  }
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(`确定要删除广告"${row.title}"吗？`, '提示', {
      type: 'warning'
    })
    await deleteAd(row.id)
    ElMessage.success('删除成功')
    loadData()
  } catch (e) {
    if (e !== 'cancel') console.error(e)
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

.pagination-box {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>
