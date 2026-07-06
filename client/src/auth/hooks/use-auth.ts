import { useContext } from 'react'
import { AuthContext } from '../auth-context'
import { ERR_AUTH_NO_PROVIDER } from '../cfg'
import type { AuthState } from '../auth-context'

export function useAuth(): AuthState {
  const ctx = useContext(AuthContext)
  if (!ctx) throw new Error(ERR_AUTH_NO_PROVIDER)
  return ctx
}
