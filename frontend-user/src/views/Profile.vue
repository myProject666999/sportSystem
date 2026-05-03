<template>
  <div class="profile-page">
    <div class="container">
      <div class="profile-container">
        <div class="profile-sidebar">
          <div class="user-info">
            <el-avatar :size="80" class="user-avatar">
              <el-icon :size="40"><User /></el-icon>
            </el-avatar>
            <div class="user-name">{{ userInfo.username }}</div>
            <div class="user-points">积分: <span class="points">{{ userInfo.points || 0 }}</span></div>
          </div>
          <el-menu
            :default-active="activeMenu"
            background-color="#fff"
            text-color="#666"
            active-text-color="#409EFF"
            @select="handleMenuSelect"
          >
            <el-menu-item index="info">
              <el-icon><User /></el-icon>
              个人信息
            </el-menu-item>
            <el-menu-item index="password">
              <el-icon><Lock /></el-icon>
              修改密码
            </el-menu-item>
            <el-menu-item index="orders" @click="$router.push('/orders')">
              <el-icon><Document /></el-icon>
              我的订单
            </el-menu-item>
            <el-menu-item index="favorites" @click="$router.push('/favorites')">
              <el-icon><Star /></el-icon>
              我的收藏
            </el-menu-item>
            <el-menu-item index="addresses" @click="$router.push('/addresses')">
              <el-icon><Location /></el-icon>
              收货地址
            </el-menu-item>
            <el-menu-item index="recharge" @click="$router.push('/recharge')">
              <el-icon><Wallet /></el-icon>
              账户充值
            </el-menu-item>
          </el-menu>
        </div>

        <div class="profile-content">
          <div v-if="activeMenu === 'info'">
            <h3 class="section-title">个人信息</h3>
            <el-form :model="infoForm" :rules="infoRules" ref="infoFormRef" label-width="100px">
              <el-form-item label="用户名" prop="username">
                <el-input v-model="infoForm.username" :disabled="true" />
              </el-form-item>
              <el-form-item label="手机号码" prop="phone">
                <el-input v-model="infoForm.phone" placeholder="请输入手机号码" />
              </el-form-item>
              <el-form-item label="邮箱" prop="email">
                <el-input v-model="infoForm.email" placeholder="请输入邮箱" />
              </el-form-item>
              <el-form-item>
                <el-button type="primary" @click="saveInfo" :loading="infoLoading">保存修改</el-button>
              </el-form-item>
            </el-form>
          </div>

          <div v-else-if="activeMenu === 'password'">
            <h3 class="section-title">修改密码</h3>
            <el-form :model="passwordForm" :rules="passwordRules" ref="passwordFormRef" label-width="120px">
              <el-form-item label="原密码" prop="oldPassword">
                <el-input v-model="passwordForm.oldPassword" type="password" placeholder="请输入原密码" show-password />
              </el-form-item>
              <el-form-item label="新密码" prop="newPassword">
                <el-input v-model="passwordForm.newPassword" type="password" placeholder="请输入新密码" show-password />
              </el-form-item>
              <el-form-item label="确认密码" prop="confirmPassword">
                <el-input v-model="passwordForm.confirmPassword" type="password" placeholder="请再次输入新密码" show-password />
              </el-form-item>
              <el-form-item>
                <el-button type="primary" @click="savePassword" :loading="passwordLoading">确认修改</el-button>
              </el-form-item>
            </el-form>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { User, Lock, Document, Star, Location, Wallet } from '@element-plus/icons-vue'
import { getUserInfo, updateUserInfo, updatePassword } from '@/api'

const router = useRouter()
const activeMenu = ref('info')
const infoLoading = ref(false)
const passwordLoading = ref(false)
const infoFormRef = ref(null)
const passwordFormRef = ref(null)

const userInfo = ref({})

const infoForm = reactive({
  username: '',
  phone: '',
  email: ''
})

const passwordForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const validateConfirmPassword = (rule, value, callback) => {
  if (value === '') {
    callback(new Error('请再次输入密码'))
  } else if (value !== passwordForm.newPassword) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

const infoRules = {
  phone: [
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号码', trigger: 'blur' }
  ],
  email: [
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
  ]
}

const passwordRules = {
  oldPassword: [
    { required: true, message: '请输入原密码', trigger: 'blur' }
  ],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请再次输入新密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' }
  ]
}

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
    infoForm.username = userInfo.value.username || ''
    infoForm.phone = userInfo.value.phone || ''
    infoForm.email = userInfo.value.email || ''
  } catch (e) {
    console.error(e)
  }
}

const handleMenuSelect = (index) => {
  activeMenu.value = index
}

const saveInfo = async () => {
  const valid = await infoFormRef.value.validate().catch(() => false)
  if (!valid) return

  infoLoading.value = true
  try {
    await updateUserInfo({
      phone: infoForm.phone,
      email: infoForm.email
    })
    ElMessage.success('修改成功')
    loadUserInfo()
  } catch (e) {
    console.error(e)
  } finally {
    infoLoading.value = false
  }
}

const savePassword = async () => {
  const valid = await passwordFormRef.value.validate().catch(() => false)
  if (!valid) return

  passwordLoading.value = true
  try {
    await updatePassword({
      old_password: passwordForm.oldPassword,
      new_password: passwordForm.newPassword
    })
    ElMessage.success('密码修改成功')
    passwordForm.oldPassword = ''
    passwordForm.newPassword = ''
    passwordForm.confirmPassword = ''
  } catch (e) {
    console.error(e)
  } finally {
    passwordLoading.value = false
  }
}

onMounted(() => {
  loadUserInfo()
})
</script>

<style scoped>
.profile-page {
  padding: 30px 0;
  background: #f5f5f5;
  min-height: calc(100vh - 140px);
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 15px;
}

.profile-container {
  display: flex;
  gap: 20px;
}

.profile-sidebar {
  width: 240px;
  background: #fff;
  border-radius: 8px;
  overflow: hidden;
}

.user-info {
  padding: 30px 20px;
  text-align: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #fff;
}

.user-avatar {
  background: rgba(255, 255, 255, 0.3);
  margin-bottom: 15px;
}

.user-name {
  font-size: 18px;
  font-weight: bold;
  margin-bottom: 10px;
}

.user-points {
  font-size: 14px;
  opacity: 0.9;
}

.points {
  color: #ffd700;
  font-weight: bold;
}

.profile-sidebar .el-menu {
  border-right: none;
}

.profile-content {
  flex: 1;
  background: #fff;
  border-radius: 8px;
  padding: 30px;
}

.section-title {
  font-size: 18px;
  color: #333;
  margin-bottom: 25px;
  padding-bottom: 15px;
  border-bottom: 1px solid #eee;
}
</style>
