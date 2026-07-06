import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router'
import { Box, Button, Center, Paper, PasswordInput, Stack, Title } from '@mantine/core'
import { PresentationIcon } from '@phosphor-icons/react/dist/csr/Presentation'
import { useAuth } from '../hooks/use-auth'
import { CLIENT_CONFIGURE } from '../../shared/cfg/routes'
import { ERR_AUTH_INVALID_PASSWORD, ERR_AUTH_CONNECTION } from '../../shared/cfg/messages'
import { LOGIN_RESULT } from '../types'

export function LandingPage() {
  const { isAuthenticated, isLoading, login } = useAuth()
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [loggingIn, setLoggingIn] = useState(false)
  const navigate = useNavigate();

  useEffect(() => {
    if (!isLoading && isAuthenticated) {
      navigate(CLIENT_CONFIGURE)
    }
  }, [isLoading, isAuthenticated, navigate])

  async function handleLogin(e: React.SubmitEvent) {
    e.preventDefault()
    if (!password) return
    setLoggingIn(true)
    setError('')
    const result = await login(password)
    if (result === LOGIN_RESULT.Invalid) setError(ERR_AUTH_INVALID_PASSWORD)
    else if (result === LOGIN_RESULT.Error) setError(ERR_AUTH_CONNECTION)
    setLoggingIn(false)
  }

  return isAuthenticated ? null : (
    <Box h="100vh" bg="gray.0">
      <Center h="100vh">
        <Paper shadow="md" p="xl" radius="md" maw={400} w="100%">
          <form onSubmit={handleLogin}>
            <Stack align="center" gap="md">
              <PresentationIcon size="48" />
              <Title>ClassDir</Title>
              <PasswordInput
                placeholder="Password"
                value={password}
                onChange={(e) => setPassword(e.currentTarget.value)}
                error={error}
                autoFocus
                w="100%"
              />
              <Button
                type="submit"
                loading={loggingIn}
                disabled={!password}
                fullWidth
              >
                Login
              </Button>
            </Stack>
          </form>
        </Paper>
      </Center>
    </Box>
  )
}
