import { useState } from 'react'
import { api } from '../../shared/api/client'

export type LoginResult = 'ok' | 'invalid' | 'error'

export function useAuth() {
  const [isAuthenticated, setIsAuthenticated] = useState(false)
  const [isLoading, setIsLoading] = useState(true)

  async function checkAuth(): Promise<void> {
    setIsLoading(true)
    try {
      const res = await api('/api/v1/auth/check')
      setIsAuthenticated(res.ok)
    } catch {
      setIsAuthenticated(false)
    } finally {
      setIsLoading(false)
    }
  }

  async function login(password: string): Promise<LoginResult> {
    try {
      const res = await api('/api/v1/auth/login', {
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
    await api('/api/v1/auth/logout', { method: 'POST' })
    setIsAuthenticated(false)
  }

  return { isAuthenticated, isLoading, checkAuth, login, logout }
}
