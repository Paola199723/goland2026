<template>
  <div class="min-h-screen bg-gray-50">
    <LoginPage v-if="!authStore.isAuthenticated" @login="handleLogin" />
    <DashboardPage v-else @logout="handleLogout" />
  </div>
</template>

<script setup lang="ts">
import DashboardPage from './pages/DashboardPage.vue'
import LoginPage from './pages/LoginPage.vue'
import { useAuthStore } from './stores/authStore'

const authStore = useAuthStore()

const handleLogin = () => {
  // La autenticación se maneja en el store
}

const handleLogout = () => {
  authStore.logout()
}
</script>

