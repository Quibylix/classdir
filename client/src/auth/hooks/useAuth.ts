import { useState } from 'react'
import { api } from '../../shared/api/client'

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

  async function login(password: string): Promise<boolean> {
    const res = await api('/api/v1/auth/login', {
      method: 'POST',
      body: JSON.stringify({ password }),
    })
    if (res.ok) {
      setIsAuthenticated(true)
      return true
    }
    return false
  }

  return { isAuthenticated, isLoading, checkAuth, login }
}
