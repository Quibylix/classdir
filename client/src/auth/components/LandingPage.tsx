import { useEffect, useState } from 'react'
import { Button, Center, Loader, PasswordInput, Stack, Title } from '@mantine/core'
import { useAuth } from '../hooks/useAuth'

export function LandingPage() {
  const { isAuthenticated, isLoading, checkAuth, login } = useAuth()
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [loggingIn, setLoggingIn] = useState(false)

  useEffect(() => { checkAuth() }, [])

  async function handleLogin() {
    setLoggingIn(true)
    setError('')
    const ok = await login(password)
    if (!ok) setError('Invalid password')
    setLoggingIn(false)
  }

  if (isLoading) {
    return (
      <Center h="100vh">
        <Loader />
      </Center>
    )
  }

  if (isAuthenticated) {
    return (
      <Center h="100vh">
        <Stack align="center">
          <Title>ClassDir</Title>
          <Button component="a" href="/configure" size="lg">
            Go to Dashboard
          </Button>
        </Stack>
      </Center>
    )
  }

  return (
    <Center h="100vh">
      <Stack align="center">
        <Title>ClassDir</Title>
        <PasswordInput
          placeholder="Password"
          value={password}
          onChange={(e) => setPassword(e.currentTarget.value)}
          onKeyDown={(e) => e.key === 'Enter' && handleLogin()}
          error={error}
        />
        <Button onClick={handleLogin} loading={loggingIn}>
          Login
        </Button>
      </Stack>
    </Center>
  )
}
