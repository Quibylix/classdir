import { createContext, useEffect, useState } from 'react'
import { api } from '../shared/api/client'
import { AUTH_LOGIN, AUTH_LOGOUT, AUTH_CHECK } from '../shared/cfg/routes'

export type LoginResult = 'ok' | 'invalid' | 'error'

export type AuthState = {
  isAuthenticated: boolean
  isLoading: boolean
  checkAuth: () => Promise<void>
  login: (password: string) => Promise<LoginResult>
  logout: () => Promise<void>
}

export const AuthContext = createContext<AuthState | null>(null)

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [isAuthenticated, setIsAuthenticated] = useState(false)
  const [isLoading, setIsLoading] = useState(true)

  useEffect(() => { checkAuth() }, [])

  async function checkAuth(): Promise<void> {
    setIsLoading(true)
    try {
      const res = await api(AUTH_CHECK)
      setIsAuthenticated(res.ok)
    } catch {
      setIsAuthenticated(false)
    } finally {
      setIsLoading(false)
    }
  }

  async function login(password: string): Promise<LoginResult> {
    try {
      const res = await api(AUTH_LOGIN, {
        method: 'POST',
        body: JSON.stringify({ password }),
      })
      if (res.ok) {
        setIsAuthenticated(true)
        return 'ok'
      }
      return 'invalid'
    } catch {
      return 'error'
    }
  }

  async function logout(): Promise<void> {
    await api(AUTH_LOGOUT, { method: 'POST' })
    setIsAuthenticated(false)
  }

  return (
    <AuthContext.Provider value={{ isAuthenticated, isLoading, checkAuth, login, logout }}>
      {children}
    </AuthContext.Provider>
  )
}
