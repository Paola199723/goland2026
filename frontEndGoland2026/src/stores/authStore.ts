import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Challenge } from '../services/api'

// Función helper para verificar si el token está expirado
const isTokenExpired = (): boolean => {
  const expiresAt = localStorage.getItem('expiresAt')
  if (!expiresAt) return true
  return new Date(expiresAt) <= new Date()
}

export const useAuthStore = defineStore('auth', () => {
  const email = ref(localStorage.getItem('userEmail') || '')
  const token = ref(localStorage.getItem('authToken') || '')
  const userId = ref(Number(localStorage.getItem('userId') || 0))
  const expiresAt = ref(localStorage.getItem('expiresAt') || '')
  // Solo autenticado si hay token Y NO está expirado
  const isAuthenticated = ref(!!token.value && !isTokenExpired())
  const challenges = ref<Challenge[]>([])
  const currentPage = ref(1)
  const totalPages = ref(1)
  const nextPageCursor = ref('')
  const isLoading = ref(false)
  const error = ref('')

  const setAuth = (authToken: string, userEmail: string, userID: number, expirationTime: string) => {
    email.value = userEmail
    token.value = authToken
    userId.value = userID
    expiresAt.value = expirationTime
    isAuthenticated.value = true

    // Guardar en localStorage
    localStorage.setItem('authToken', authToken)
    localStorage.setItem('userEmail', userEmail)
    localStorage.setItem('userId', userID.toString())
    localStorage.setItem('expiresAt', expirationTime)
  }

  const setChallenges = (data: Challenge[], total: number, nextPage?: string) => {
    challenges.value = data
    totalPages.value = total
    if (nextPage) {
      nextPageCursor.value = nextPage
    }
  }

  const setCurrentPage = (page: number) => {
    currentPage.value = page
  }

  const setLoading = (loading: boolean) => {
    isLoading.value = loading
  }

  const setError = (err: string) => {
    error.value = err
  }

  const logout = () => {
    email.value = ''
    token.value = ''
    userId.value = 0
    expiresAt.value = ''
    isAuthenticated.value = false
    challenges.value = []
    currentPage.value = 1
    totalPages.value = 1
    nextPageCursor.value = ''
    error.value = ''

    // Limpiar localStorage
    localStorage.removeItem('authToken')
    localStorage.removeItem('userEmail')
    localStorage.removeItem('userId')
    localStorage.removeItem('expiresAt')
  }

  return {
    email,
    token,
    userId,
    expiresAt,
    isAuthenticated,
    challenges,
    currentPage,
    totalPages,
    nextPageCursor,
    isLoading,
    error,
    setAuth,
    setChallenges,
    setCurrentPage,
    setLoading,
    setError,
    logout,
  }
})

