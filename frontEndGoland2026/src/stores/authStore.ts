import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Challenge } from '../services/api'

export const useAuthStore = defineStore('auth', () => {
  const email = ref('')
  const token = ref('')
  const isAuthenticated = ref(false)
  const challenges = ref<Challenge[]>([])
  const currentPage = ref(1)
  const totalPages = ref(1)
  const nextPageCursor = ref('')
  const isLoading = ref(false)
  const error = ref('')

  const setAuth = (userEmail: string, authToken: string) => {
    email.value = userEmail
    token.value = authToken
    isAuthenticated.value = true
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
    isAuthenticated.value = false
    challenges.value = []
    currentPage.value = 1
    totalPages.value = 1
    nextPageCursor.value = ''
    error.value = ''
  }

  return {
    email,
    token,
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

