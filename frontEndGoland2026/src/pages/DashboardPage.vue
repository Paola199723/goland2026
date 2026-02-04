<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Header -->
    <header class="bg-white shadow">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
        <div class="flex justify-between items-center">
          <div class="flex items-center space-x-3">
            <div class="w-10 h-10 bg-blue-600 rounded-full flex items-center justify-center text-white font-bold">
              👤
            </div>
            <div>
              <p class="text-sm text-gray-500">Usuario</p>
              <p class="text-lg font-semibold text-gray-900">{{ authStore.email }}</p>
            </div>
          </div>
          <button
            @click="handleLogout"
            class="bg-red-500 hover:bg-red-600 text-white px-4 py-2 rounded-lg transition"
          >
            Cerrar Sesión
          </button>
        </div>
      </div>
    </header>

    <!-- Main Content -->
    <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <!-- Welcome Section -->
      <div class="bg-white rounded-lg shadow-md p-6 mb-8">
        <h1 class="text-3xl font-bold text-gray-900 mb-2">
          Bienvenido a la Prueba
        </h1>
        <p class="text-gray-600">
          Aquí puedes visualizar los últimos cambios de ratings y targets de valores bursátiles
        </p>
      </div>

      <!-- Challenges Table -->
      <div class="bg-white rounded-lg shadow-md overflow-hidden">
        <div class="px-6 py-4 border-b border-gray-200">
          <h2 class="text-xl font-semibold text-gray-900">Actualizaciones de Challenges</h2>
        </div>

        <div v-if="authStore.isLoading" class="p-8 text-center">
          <p class="text-gray-500">Cargando datos...</p>
        </div>

        <div v-else-if="authStore.challenges.length === 0" class="p-8 text-center">
          <p class="text-gray-500">No hay desafíos disponibles</p>
        </div>

        <div v-else class="overflow-x-auto">
          <table class="w-full">
            <thead class="bg-gray-50 border-b border-gray-200">
              <tr>
                <th class="px-6 py-3 text-left text-sm font-semibold text-gray-900">Ticker</th>
                <th class="px-6 py-3 text-left text-sm font-semibold text-gray-900">Empresa</th>
                <th class="px-6 py-3 text-left text-sm font-semibold text-gray-900">Casa de Bolsa</th>
                <th class="px-6 py-3 text-left text-sm font-semibold text-gray-900">Acción</th>
                <th class="px-6 py-3 text-left text-sm font-semibold text-gray-900">Rating Anterior</th>
                <th class="px-6 py-3 text-left text-sm font-semibold text-gray-900">Rating Nuevo</th>
                <th class="px-6 py-3 text-left text-sm font-semibold text-gray-900">Target Anterior</th>
                <th class="px-6 py-3 text-left text-sm font-semibold text-gray-900">Target Nuevo</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-200">
              <tr v-for="(challenge, index) in authStore.challenges" :key="index" class="hover:bg-gray-50">
                <td class="px-6 py-4 text-sm font-semibold text-blue-600">{{ challenge.ticker }}</td>
                <td class="px-6 py-4 text-sm text-gray-900">{{ challenge.company }}</td>
                <td class="px-6 py-4 text-sm text-gray-600">{{ challenge.brokerage || 'N/A' }}</td>
                <td class="px-6 py-4 text-sm text-gray-900">{{ challenge.action }}</td>
                <td class="px-6 py-4 text-sm text-gray-600">{{ challenge.rating_from }}</td>
                <td class="px-6 py-4 text-sm text-gray-600">{{ challenge.rating_to }}</td>
                <td class="px-6 py-4 text-sm text-gray-600">{{ challenge.target_from }}</td>
                <td class="px-6 py-4 text-sm text-gray-600">{{ challenge.target_to }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Pagination (cursor-based: next_page) -->
      <div class="mt-8 flex justify-center items-center space-x-4">
        <button
          @click="goPrevious"
          :disabled="currentPage <= 1"
          class="px-4 py-2 rounded-lg font-medium transition bg-white text-gray-900 border border-gray-300 hover:bg-gray-50 disabled:opacity-50"
        >
          Anterior
        </button>

        <div class="text-sm text-gray-600">Página {{ currentPage }}</div>

        <button
          @click="goNext"
          :disabled="!authStore.nextPageCursor"
          class="px-4 py-2 rounded-lg font-medium transition bg-blue-600 text-white disabled:opacity-50"
        >
          Siguiente
        </button>
      </div>

      <!-- Recommendations Section -->
      <div class="bg-white rounded-lg shadow-md p-6 mt-8">
        <div class="flex items-center mb-4">
          <h2 class="text-xl font-semibold text-gray-900 mr-2">Recomendaciones del día</h2>
          <button @click="refreshRecommendations" class="p-2 rounded-full hover:bg-blue-100 transition" :title="'Refrescar recomendaciones'">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-blue-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582M20 20v-5h-.581M5.21 17.293A8.001 8.001 0 0112 4a8 8 0 017.418 5.293M18.36 6.64A8.001 8.001 0 014 20a8 8 0 01-3.418-5.293" />
            </svg>
          </button>
        </div>
        <div class="bg-gray-50 border-2 border-dashed border-gray-300 rounded-lg p-8 text-center text-gray-500">
          <div>
    <h2>Recomendación del día</h2>
      <div v-if="data && data.recommended_stock">
        <div class="overflow-x-auto">
          <table class="min-w-full divide-y divide-gray-200 rounded-lg shadow">
            <thead class="bg-blue-600">
              <tr>
                <th colspan="2" class="px-6 py-4 text-left text-lg font-bold text-white rounded-t-lg">Recomendación principal</th>
              </tr>
            </thead>
            <tbody class="bg-white">
              <tr>
                <td class="px-6 py-4 font-semibold text-gray-700">Empresa</td>
                <td class="px-6 py-4">{{ data.recommended_stock.company }} ({{ data.recommended_stock.ticker }})</td>
              </tr>
              <tr>
                <td class="px-6 py-4 font-semibold text-gray-700">Acción sugerida</td>
                <td class="px-6 py-4">
                  <span :class="{
                    'text-green-600 font-bold': data.recommended_stock.action === 'BUY',
                    'text-yellow-600 font-bold': data.recommended_stock.action === 'HOLD',
                    'text-orange-600 font-bold': data.recommended_stock.action === 'WATCH',
                    'text-red-600 font-bold': data.recommended_stock.action === 'SELL',
                  }">{{ data.recommended_stock.action }}</span>
                </td>
              </tr>
              <tr>
                <td class="px-6 py-4 font-semibold text-gray-700">Score</td>
                <td class="px-6 py-4">{{ data.recommended_stock.score.toFixed(2) }}</td>
              </tr>
              <tr>
                <td class="px-6 py-4 font-semibold text-gray-700">Confianza</td>
                <td class="px-6 py-4">{{ data.recommended_stock.confidence }}</td>
              </tr>
              <tr>
                <td class="px-6 py-4 font-semibold text-gray-700 align-top">Razones</td>
                <td class="px-6 py-4">
                  <ul class="list-disc list-inside text-left">
                    <li v-for="reason in data.recommended_stock.reasons" :key="reason">{{ reason }}</li>
                  </ul>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
        <div class="overflow-x-auto mt-8">
          <h3 class="text-lg font-semibold text-gray-900 mb-2">Ranking de acciones</h3>
          <table class="min-w-full divide-y divide-gray-200 rounded-lg shadow">
            <thead class="bg-gray-100">
              <tr>
                <th class="px-6 py-3 text-left text-sm font-semibold text-gray-900">Ticker</th>
                <th class="px-6 py-3 text-left text-sm font-semibold text-gray-900">Score</th>
                <th class="px-6 py-3 text-left text-sm font-semibold text-gray-900">Acción</th>
              </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200">
              <tr v-for="item in data.ranking" :key="item.ticker">
                <td class="px-6 py-4 font-semibold text-blue-600">{{ item.ticker }}</td>
                <td class="px-6 py-4">{{ item.score.toFixed(2) }}</td>
                <td class="px-6 py-4">
                  <span :class="{
                    'text-green-600 font-bold': item.action === 'BUY',
                    'text-yellow-600 font-bold': item.action === 'HOLD',
                    'text-orange-600 font-bold': item.action === 'WATCH',
                    'text-red-600 font-bold': item.action === 'SELL',
                  }">{{ item.action }}</span>
                </td>
              </tr>
            </tbody>
          </table>
          <p class="disclaimer mt-4 text-xs text-gray-500">{{ data.disclaimer }}</p>
        </div>
      </div>
    <div v-else-if="data && data.error">
      <p class="text-red-500">Error: {{ data.error }}</p>
    </div>
    <div v-else>
      <p>Cargando recomendaciones...</p>
    </div>
        </div>
        </div>
      </div>
    </main>
  </div>
</template>
<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { challengeService } from '../services/api';
import { useAuthStore } from '../stores/authStore';

/* ========================
   Stores & emits
======================== */
const authStore = useAuthStore()

const emit = defineEmits<{
  logout: []
}>()

/* ========================
   State
======================== */
const currentPage = ref(1)
const prevCursors: string[] = []

const data = ref<any>(null)

/* ========================
   Pagination logic
======================== */
const goToPage = async (nextPageCursor?: string) => {
  authStore.setLoading(true)
  try {
    const response = await challengeService.getChallenges(nextPageCursor)
    authStore.setChallenges(
      response.challenges,
      response.total_pages,
      response.next_page
    )
  } catch (error) {
    console.error('Error loading page:', error)
    authStore.setError('Error al cargar los desafíos')
  } finally {
    authStore.setLoading(false)
  }
}

const goNext = async () => {
  prevCursors.push(authStore.nextPageCursor || '')
  await goToPage(authStore.nextPageCursor)
  currentPage.value += 1
}

const goPrevious = async () => {
  if (prevCursors.length === 0) return
  const prev = prevCursors.pop() || ''
  await goToPage(prev === '' ? undefined : prev)
  if (currentPage.value > 1) currentPage.value -= 1
}

/* ========================
   Auth
======================== */
const handleLogout = () => {
  authStore.logout()
  emit('logout')
}

/* ========================
   Lifecycle
======================== */
onMounted(async () => {
  // Cargar challenges
  if (authStore.challenges.length === 0) {
    await goToPage()
  }

  await fetchRecommendations()
/* ========================
   Recomendaciones
======================== */
async function fetchRecommendations() {
  const token = authStore.token || localStorage.getItem('authToken')
  if (!token) {
    data.value = { error: 'No autenticado. Inicia sesión para ver recomendaciones.' }
    return
  }
  try {
    const res = await fetch('http://localhost:8081/api/recommendation/today', {
      headers: {
        Authorization: `Bearer ${token}`
      }
    })
    data.value = await res.json()
  } catch (err) {
    data.value = { error: 'No se pudo obtener la recomendación.' }
  }
}

function refreshRecommendations() {
  fetchRecommendations()
}
})
</script>
