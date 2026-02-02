import axios from 'axios'

const API_BASE_URL = 'http://localhost:8080/api'

export interface LoginRequest {
  email: string
  password: string
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

export interface LoginResponse {
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

export const authService = {
  login: async (email: string, password: string): Promise<LoginResponse> => {
    const response = await apiClient.post<LoginResponse>('/login', {
      email,
      password,
    })
    return response.data
  },
}

export const challengeService = {
  getChallenges: async (page: number): Promise<LoginResponse> => {
    const response = await apiClient.get<LoginResponse>('/challenges', {
      params: {
        page,
      },
    })
    return response.data
  },
}

        token,
      },
    })
    return response.data
  },
}
