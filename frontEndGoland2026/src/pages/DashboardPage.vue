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
        <h2 class="text-xl font-semibold text-gray-900 mb-4">Recomendaciones</h2>
        <div class="bg-gray-50 border-2 border-dashed border-gray-300 rounded-lg p-8 text-center text-gray-500">
          <p>Sección de recomendaciones (pendiente de implementar)</p>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { challengeService } from '../services/api'
import { useAuthStore } from '../stores/authStore'

const authStore = useAuthStore()

// Local page counter to display page number when using cursor-based pagination
const currentPage = ref(1)

// Stack to keep previous cursors for "Anterior" functionality
const prevCursors: string[] = []

const goToPage = async (nextPageCursor?: string) => {
  authStore.setLoading(true)

  try {
    const response = await challengeService.getChallenges(nextPageCursor)

    // Si recibimos next_page, pushear el cursor actual para poder volver atrás
    if (authStore.nextPageCursor && !nextPageCursor) {
      // no-op
    }

    // Actualizar store con resultados
    authStore.setChallenges(response.challenges, response.total_pages, response.next_page)

  } catch (error) {
    console.error('Error loading page:', error)
    authStore.setError('Error al cargar los desafíos')
  } finally {
    authStore.setLoading(false)
  }
}

const goNext = async () => {
  // guardar cursor actual para poder retroceder
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

const handleLogout = () => {
  authStore.logout()
  emit('logout')
}

// Cargar desafíos cuando se monta el componente
onMounted(async () => {
  if (authStore.challenges.length === 0) {
    // Cargar la primera página (sin cursor)
    await goToPage()
  }
})

const emit = defineEmits<{
  logout: []
}>()
</script>
