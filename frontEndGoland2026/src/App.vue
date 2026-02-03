<template>
  <div class="min-h-screen bg-gray-50">
    <LoginPage v-if="!authStore.isAuthenticated" @login="handleLogin" />
    <DashboardPage v-else @logout="handleLogout" />
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import DashboardPage from './pages/DashboardPage.vue'
import LoginPage from './pages/LoginPage.vue'
import { authService } from './services/api'
import { useAuthStore } from './stores/authStore'

const authStore = useAuthStore()

const handleLogin = () => {
  // La autenticación se maneja en el store
}

const handleLogout = () => {
  authStore.logout()
}

onMounted(async () => {
  try {
    // Si hay un token en localStorage, comprobar con el backend si sigue válido
    if (authStore.token) {
      try {
        const ok = await authService.verify()
        if (!ok) {
          authStore.logout()
        }
      } catch (err) {
        // Si la verificación falla, logout
        console.error('Token verification failed:', err)
        authStore.logout()
      }
    }
  } catch (err) {
    console.error('App mount error:', err)
  }
})
</script>

