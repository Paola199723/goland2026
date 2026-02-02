<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-500 to-blue-700 px-4">
    <div class="bg-white rounded-lg shadow-lg p-8 w-full max-w-md">
      <div class="text-center mb-8">
        <h1 class="text-3xl font-bold text-gray-800">Bienvenido</h1>
        <p class="text-gray-600 mt-2">Prueba Goland 2026 - Challenge Tracker</p>
      </div>

      <form @submit.prevent="handleLogin" class="space-y-6">
        <div>
          <label for="email" class="block text-sm font-medium text-gray-700 mb-2">
            Correo Electrónico
          </label>
          <input
            id="email"
            v-model="form.email"
            type="email"
            required
            placeholder="correo@example.com"
            class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>

        <div>
          <label for="password" class="block text-sm font-medium text-gray-700 mb-2">
            Contraseña
          </label>
          <input
            id="password"
            v-model="form.password"
            type="password"
            required
            placeholder="••••••••"
            class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>

        <div v-if="authStore.error" class="bg-red-50 border border-red-200 text-red-600 px-4 py-3 rounded-lg text-sm">
          {{ authStore.error }}
        </div>

        <button
          type="submit"
          :disabled="authStore.isLoading"
          class="w-full bg-blue-600 hover:bg-blue-700 disabled:bg-gray-400 text-white font-semibold py-2 px-4 rounded-lg transition duration-200"
        >
          <span v-if="!authStore.isLoading">Ingresar</span>
          <span v-else>Cargando...</span>
        </button>
      </form>

      <p class="text-center text-gray-600 text-sm mt-6">
        Ingresa tu correo y contraseña para acceder
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { authService } from '../services/api'
import { useAuthStore } from '../stores/authStore'

const authStore = useAuthStore()

const form = ref({
  email: '',
  password: '',
})

const handleLogin = async () => {
  authStore.setError('')
  authStore.setLoading(true)

  try {
    const response = await authService.login(form.value.email, form.value.password)
    
    // Guardar token y información del usuario
    authStore.setAuth(response.token, response.email, response.user_id, response.expires_at)
    
    // Emitir evento para cambiar de página en App.vue
    emit('login')
  } catch (err: any) {
    const errorMessage = err.response?.data?.message || err.response?.data?.error || 'Error al iniciar sesión'
    authStore.setError(errorMessage)
  } finally {
    authStore.setLoading(false)
  }
}

const emit = defineEmits<{
  login: []
}>()
</script>

