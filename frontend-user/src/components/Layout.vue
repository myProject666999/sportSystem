<template>
  <div class="layout">
    <el-header class="header">
      <div class="container header-content">
        <router-link to="/" class="logo">
          <el-icon class="logo-icon"><Trophy /></el-icon>
          <span class="logo-text">体育商城</span>
        </router-link>
        
        <el-menu
          :default-active="activeMenu"
          mode="horizontal"
          background-color="#1E3050"
          text-color="#fff"
          active-text-color="#409EFF"
          router
          class="nav-menu"
        >
          <el-menu-item index="/">首页</el-menu-item>
          <el-menu-item index="/products">商品中心</el-menu-item>
          <el-menu-item index="/news">新闻资讯</el-menu-item>
          <el-menu-item index="/videos">赛事视频</el-menu-item>
          <el-menu-item index="/notices">公告通知</el-menu-item>
        </el-menu>
        
        <div class="search-box">
          <el-input
            v-model="searchKeyword"
            placeholder="搜索商品..."
            class="search-input"
            @keyup.enter="handleSearch"
          >
            <template #append>
              <el-button type="primary" @click="handleSearch">
                <el-icon><Search /></el-icon>
              </el-button>
            </template>
          </el-input>
        </div>
        
        <div class="user-actions">
          <template v-if="isLoggedIn">
            <router-link to="/cart" class="action-item">
              <el-badge :value="cartCount" :hidden="cartCount === 0" class="cart-badge">
                <el-icon class="action-icon"><ShoppingCart /></el-icon>
              </el-badge>
            </router-link>
            <el-dropdown @command="handleCommand">
              <span class="user-info">
                <el-avatar :size="32" class="user-avatar">
                  <el-icon><User /></el-icon>
                </el-avatar>
                <span class="username">{{ userInfo.username }}</span>
              </span>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="profile">个人中心</el-dropdown-item>
                  <el-dropdown-item command="orders">我的订单</el-dropdown-item>
                  <el-dropdown-item command="favorites">我的收藏</el-dropdown-item>
                  <el-dropdown-item command="addresses">地址管理</el-dropdown-item>
                  <el-dropdown-item command="recharge">账户充值</el-dropdown-item>
                  <el-dropdown-item divided command="logout">退出登录</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
          <template v-else>
            <router-link to="/login" class="action-btn">登录</router-link>
            <router-link to="/register" class="action-btn primary">注册</router-link>
          </template>
        </div>
      </div>
    </el-header>
    
    <el-main class="main-content">
      <slot />
    </el-main>
    
    <el-footer class="footer">
      <div class="container footer-content">
        <div class="footer-section">
          <h4>关于我们</h4>
          <p>体育商城是专业的体育用品购物平台</p>
          <p>为您提供优质的运动装备和服务</p>
        </div>
        <div class="footer-section">
          <h4>帮助中心</h4>
          <p>购物指南</p>
          <p>支付方式</p>
          <p>配送方式</p>
        </div>
        <div class="footer-section">
          <h4>联系我们</h4>
          <p>客服热线：400-888-8888</p>
          <p>工作时间：9:00-21:00</p>
          <p>电子邮箱：service@sport.com</p>
        </div>
      </div>
      <div class="copyright">
        <p>© 2024 体育商城 版权所有</p>
      </div>
    </el-footer>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Trophy, Search, ShoppingCart, User } from '@element-plus/icons-vue'
import { getCartList } from '@/api'

const router = useRouter()
const route = useRoute()

const searchKeyword = ref('')
const cartCount = ref(0)

const isLoggedIn = computed(() => !!localStorage.getItem('token'))
const userInfo = computed(() => {
  const info = localStorage.getItem('userInfo')
  return info ? JSON.parse(info) : {}
})

const activeMenu = computed(() => {
  const path = route.path
  if (path.startsWith('/product')) return '/products'
  if (path.startsWith('/news')) return '/news'
  if (path.startsWith('/video')) return '/videos'
  if (path.startsWith('/notice')) return '/notices'
  return path
})

const handleSearch = () => {
  if (searchKeyword.value.trim()) {
    router.push({ path: '/products', query: { keyword: searchKeyword.value }})
  }
}

const handleCommand = (command) => {
  switch (command) {
    case 'profile':
      router.push('/profile')
      break
    case 'orders':
      router.push('/orders')
      break
    case 'favorites':
      router.push('/favorites')
      break
    case 'addresses':
      router.push('/addresses')
      break
    case 'recharge':
      router.push('/recharge')
      break
    case 'logout':
      ElMessageBox.confirm('确定要退出登录吗？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        localStorage.removeItem('token')
        localStorage.removeItem('userInfo')
        ElMessage.success('已退出登录')
        router.push('/')
      }).catch(() => {})
      break
  }
}

const fetchCartCount = async () => {
  if (!isLoggedIn.value) return
  try {
    const res = await getCartList()
    cartCount.value = res.data ? res.data.length : 0
  } catch (e) {
    console.error(e)
  }
}

onMounted(() => {
  fetchCartCount()
})
</script>

<style scoped>
.layout {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.header {
  background: linear-gradient(135deg, #1E3050 0%, #2E4A70 100%);
  padding: 0;
  height: 70px;
  line-height: 70px;
}

.header-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 100%;
}

.logo {
  display: flex;
  align-items: center;
  color: #fff;
  text-decoration: none;
}

.logo-icon {
  font-size: 32px;
  margin-right: 10px;
}

.logo-text {
  font-size: 22px;
  font-weight: bold;
}

.nav-menu {
  border: none;
  background: transparent;
  flex: 1;
  margin-left: 30px;
}

.nav-menu :deep(.el-menu-item) {
  border-bottom: none;
  height: 70px;
  line-height: 70px;
  margin: 0 5px;
}

.search-box {
  width: 300px;
}

.search-input {
  margin-right: 20px;
}

.user-actions {
  display: flex;
  align-items: center;
  gap: 15px;
}

.action-item {
  position: relative;
  color: #fff;
  font-size: 24px;
  padding: 5px;
}

.action-icon {
  font-size: 24px;
}

.cart-badge {
  position: relative;
  top: -10px;
}

.user-info {
  display: flex;
  align-items: center;
  color: #fff;
  cursor: pointer;
}

.user-avatar {
  margin-right: 8px;
  background: rgba(255, 255, 255, 0.2);
}

.username {
  max-width: 100px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.action-btn {
  padding: 8px 20px;
  color: #fff;
  text-decoration: none;
  border-radius: 4px;
  transition: all 0.3s;
}

.action-btn:hover {
  background: rgba(255, 255, 255, 0.1);
}

.action-btn.primary {
  background: #409EFF;
}

.action-btn.primary:hover {
  background: #66b1ff;
}

.main-content {
  flex: 1;
  padding: 20px;
  background: #f5f5f5;
}

.footer {
  background: #2c3e50;
  color: #bdc3c7;
  padding: 0;
  height: auto;
  line-height: normal;
}

.footer-content {
  display: flex;
  justify-content: space-between;
  padding: 30px 0;
}

.footer-section {
  flex: 1;
}

.footer-section h4 {
  color: #fff;
  margin-bottom: 15px;
  font-size: 16px;
}

.footer-section p {
  margin: 8px 0;
  font-size: 14px;
}

.copyright {
  border-top: 1px solid #34495e;
  padding: 20px 0;
  text-align: center;
  font-size: 14px;
}
</style>
