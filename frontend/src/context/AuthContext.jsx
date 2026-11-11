import { useCallback, useEffect, useMemo, useState } from 'react'
import { login as requestLogin } from '../services/authService.js'
import { setAccessToken } from '../services/apiClient.js'
import AuthContext from './authContext.js'

export function AuthProvider({ children }) {
  const [session, setSession] = useState(null)

  const clearSession = useCallback(() => {
    setAccessToken(null)
    setSession(null)
  }, [])

  const login = useCallback(async (email, password) => {
    const result = await requestLogin(email, password)
    const expiresAt = Date.parse(result.expires_at)
    if (expiresAt <= Date.now()) {
      throw new Error('The server returned an expired access token.')
    }

    setAccessToken(result.access_token)
    setSession({
      token: result.access_token,
      expiresAt,
      user: result.user,
    })
    return result.user
  }, [])

  useEffect(() => {
    const onUnauthorized = () => clearSession()
    window.addEventListener('auth:unauthorized', onUnauthorized)
    return () => window.removeEventListener('auth:unauthorized', onUnauthorized)
  }, [clearSession])

  useEffect(() => {
    if (!session) return undefined

    const remaining = Math.max(session.expiresAt - Date.now(), 0)
    const timeout = window.setTimeout(clearSession, remaining)
    return () => window.clearTimeout(timeout)
  }, [session, clearSession])

  const value = useMemo(
    () => ({
      user: session?.user ?? null,
      isAuthenticated: Boolean(session?.token),
      login,
      logout: clearSession,
    }),
    [session, login, clearSession],
  )

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}
