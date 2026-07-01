import { useState } from 'react'
import { api } from '../../shared/api/client'
import { AUTH_LOGIN, AUTH_LOGOUT, AUTH_CHECK } from '../../shared/cfg/routes'

export type LoginResult = 'ok' | 'invalid' | 'error'

export function useAuth() {
  const [isAuthenticated, setIsAuthenticated] = useState(false)
  const [isLoading, setIsLoading] = useState(true)

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

  return { isAuthenticated, isLoading, checkAuth, login, logout }
}
