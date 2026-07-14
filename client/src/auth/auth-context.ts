import { createContext } from 'react'
import type { LoginResult } from './types'

export type AuthState = {
  isAuthenticated: boolean
  isLoading: boolean
  checkAuth: () => Promise<void>
  login: (password: string) => Promise<LoginResult>
  logout: () => Promise<void>
}

export const AuthContext = createContext<AuthState | null>(null)
