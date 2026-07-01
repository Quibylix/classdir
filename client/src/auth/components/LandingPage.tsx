import { useEffect, useState } from 'react'
import { Box, Button, Center, Loader, Paper, PasswordInput, Stack, Title } from '@mantine/core'
import { PresentationIcon } from '@phosphor-icons/react/dist/csr/Presentation'
import { useAuth } from '../hooks/useAuth'

export function LandingPage() {
  const { isAuthenticated, isLoading, checkAuth, login, logout } = useAuth()
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [loggingIn, setLoggingIn] = useState(false)

  useEffect(() => { checkAuth() }, [])

  async function handleLogin() {
    if (!password) return
    setLoggingIn(true)
    setError('')
    const result = await login(password)
    if (result === 'invalid') setError('Invalid password')
    else if (result === 'error') setError('Could not connect to server')
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
      <Box h="100vh" bg="gray.0">
        <Center h="100vh">
          <Stack align="center" gap="md">
            <PresentationIcon size="48" />
            <Title>ClassDir</Title>
            <Button component="a" href="/configure" size="lg">
              Go to Dashboard
            </Button>
            <Button variant="subtle" onClick={logout}>
              Logout
            </Button>
          </Stack>
        </Center>
      </Box>
    )
  }

  return (
    <Box h="100vh" bg="gray.0">
      <Center h="100vh">
        <Paper shadow="md" p="xl" radius="md" maw={400} w="100%">
          <Stack align="center" gap="md">
            <PresentationIcon size="48" />
            <Title>ClassDir</Title>
            <PasswordInput
              placeholder="Password"
              value={password}
              onChange={(e) => setPassword(e.currentTarget.value)}
              onKeyDown={(e) => e.key === 'Enter' && handleLogin()}
              error={error}
              autoFocus
              w="100%"
            />
            <Button
              onClick={handleLogin}
              loading={loggingIn}
              disabled={!password}
              fullWidth
            >
              Login
            </Button>
          </Stack>
        </Paper>
      </Center>
    </Box>
  )
}
