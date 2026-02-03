import axios from 'axios'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL + '/api'

export interface LoginRequest {
  email: string
  password: string
}

export interface AuthTokenResponse {
  token: string
  expires_at: string
  user_id: number
  email: string
}

export interface Challenge {
  ticker: string
  target_from: string
  target_to: string
  company: string
  action: string
  brokerage: string
  rating_from: string
  rating_to: string
  time: string
}

export interface ChallengesResponse {
  email: string
  challenges: Challenge[]
  next_page?: string
  total_pages: number
}

const apiClient = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Interceptor para agregar token a cada request
apiClient.interceptors.request.use((config) => {
  const token = localStorage.getItem('authToken')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// Interceptor para manejar errores de token expirado
apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 400 && error.response?.data?.message?.includes('expired')) {
      // Token expirado - limpiar storage y redirigir a login
      localStorage.removeItem('authToken')
      localStorage.removeItem('userEmail')
      localStorage.removeItem('userId')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

export const authService = {
  login: async (email: string, password: string): Promise<AuthTokenResponse> => {
    const response = await apiClient.post<AuthTokenResponse>('/auth/login', {
      email,
      password,
    })
    
    // Guardar token en localStorage
    localStorage.setItem('authToken', response.data.token)
    localStorage.setItem('userEmail', response.data.email)
    localStorage.setItem('userId', response.data.user_id.toString())
    
    return response.data
  },

  verify: async (): Promise<boolean> => {
    try {
      await apiClient.get('/auth/verify')
      return true
    } catch {
      return false
    }
  },

  logout: () => {
    localStorage.removeItem('authToken')
    localStorage.removeItem('userEmail')
    localStorage.removeItem('userId')
  },
}

export const challengeService = {
  // Enviar next_page (cursor) en lugar de page cuando esté disponible
  getChallenges: async (nextPage?: string): Promise<ChallengesResponse> => {
    const params: Record<string, any> = {}
    if (nextPage) params.next_page = nextPage

    const response = await apiClient.get<ChallengesResponse>('/challenges', {
      params,
    })
    return response.data
  },
}