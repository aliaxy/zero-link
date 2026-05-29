<template>
  <div class="layout">
    <aside class="sidebar">
      <div class="sidebar-header">
        <span class="logo">zero·link</span>
      </div>

      <nav class="sidebar-nav">
        <router-link class="nav-item" :class="{ active: $route.path.startsWith('/links') }" to="/links">
          <el-icon><Link /></el-icon>
          <span>Links</span>
        </router-link>
      </nav>

      <div class="sidebar-footer">
        <div class="user-info">
          <span class="username">{{ auth.admin?.username }}</span>
          <span class="user-label">Administrator</span>
        </div>
        <button class="logout-btn" @click="handleLogout" title="Sign out">
          <el-icon><SwitchButton /></el-icon>
        </button>
      </div>
    </aside>

    <main class="main">
      <router-view v-slot="{ Component }">
        <transition name="page" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>
    </main>
  </div>
</template>

<script setup>
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'

const auth = useAuthStore()
const router = useRouter()

function handleLogout() {
  auth.logout()
  router.push('/login')
}
</script>

<style scoped>
.layout {
  display: flex;
  height: 100vh;
  background: var(--color-bg);
}

.sidebar {
  width: 220px;
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  background: var(--color-surface-frosted);
  backdrop-filter: blur(20px) saturate(180%);
  -webkit-backdrop-filter: blur(20px) saturate(180%);
  border-right: 1px solid var(--color-separator);
  position: sticky;
  top: 0;
  height: 100vh;
  z-index: 10;
}

.sidebar-header {
  padding: 24px 20px 20px;
  border-bottom: 1px solid var(--color-separator);
}

.logo {
  font-family: var(--font-mono);
  font-size: 15px;
  font-weight: 600;
  color: var(--color-accent);
  letter-spacing: -0.02em;
}

.sidebar-nav {
  flex: 1;
  padding: 12px 10px;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 9px 12px;
  border-radius: var(--radius-md);
  font-size: 14px;
  font-weight: 500;
  color: var(--color-label-secondary);
  text-decoration: none;
  transition: background var(--duration-fast) var(--ease-out),
              color var(--duration-fast) var(--ease-out);
}

.nav-item:hover {
  background: rgba(0, 0, 0, 0.04);
  color: var(--color-label);
}

.nav-item.active {
  background: rgba(0, 122, 255, 0.1);
  color: var(--color-accent);
}

.nav-item .el-icon {
  font-size: 16px;
}

.sidebar-footer {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 16px 16px 20px;
  border-top: 1px solid var(--color-separator);
}

.user-info {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
}

.username {
  font-size: 13px;
  font-weight: 600;
  color: var(--color-label);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.user-label {
  font-size: 11px;
  color: var(--color-label-tertiary);
  margin-top: 1px;
}

.logout-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 30px;
  height: 30px;
  border: none;
  background: transparent;
  border-radius: var(--radius-sm);
  color: var(--color-label-tertiary);
  cursor: pointer;
  transition: background var(--duration-fast), color var(--duration-fast);
  flex-shrink: 0;
}

.logout-btn:hover {
  background: rgba(255, 59, 48, 0.08);
  color: var(--color-destructive);
}

.main {
  flex: 1;
  min-width: 0;
  overflow-y: auto;
  padding: 32px 36px;
}
</style>
