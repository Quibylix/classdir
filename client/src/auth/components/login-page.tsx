import { useEffect, useState } from 'react'
import { useNavigate, Link } from 'react-router'
import { Box, Button, Center, Container, Group, Paper, PasswordInput, Stack, Text, Title } from '@mantine/core'
import { PresentationIcon } from '@phosphor-icons/react/dist/csr/Presentation'
import { useAuth } from '../hooks/use-auth'
import { CLIENT_LOGIN, CLIENT_CONFIGURE } from '../../shared/cfg/routes'
import { ERR_AUTH_INVALID_PASSWORD, ERR_AUTH_CONNECTION } from '../../shared/cfg/messages'
import { LOGIN_RESULT } from '../types'

export function LoginPage() {
  const { isAuthenticated, isLoading, login } = useAuth()
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [loggingIn, setLoggingIn] = useState(false)
  const navigate = useNavigate()

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
    const result = await login(password).finally(() => setLoggingIn(false))
    if (result === LOGIN_RESULT.Ok) {
      navigate(CLIENT_CONFIGURE)
    } else if (result === LOGIN_RESULT.Invalid) {
      setError(ERR_AUTH_INVALID_PASSWORD)
    } else {
      setError(ERR_AUTH_CONNECTION)
    }
  }

  return (
    <Box bg="dark.9" mih="100vh">
      <Box
        component="header"
        pos="fixed"
        top={0}
        left={0}
        right={0}
        bg="color-mix(in srgb, var(--mantine-color-dark-9) 80%, transparent)"
        style={{
          zIndex: 40,
          backdropFilter: 'blur(12px)',
          borderBottom: '1px solid var(--mantine-color-dark-8)',
        }}
      >
        <Container size="lg">
          <Group h={64} justify="space-between">
            <Group gap="xs">
              <Center
                w={36}
                h={36}
                bg="blue.6"
                style={{ borderRadius: 8, boxShadow: '0 4px 12px color-mix(in srgb, var(--mantine-color-blue-5) 20%, transparent)' }}
              >
                <PresentationIcon size={20} color="white" />
              </Center>
              <Text fw={700} size="xl" c="white" lts="-0.02em">
                Class<Text component="span" fw="inherit" c="blue.4">Dir</Text>
              </Text>
            </Group>

            <Button
              component={Link}
              to={isAuthenticated ? CLIENT_CONFIGURE : CLIENT_LOGIN}
              variant="outline"
              color="gray"
              size="sm"
              style={{ borderColor: 'var(--mantine-color-dark-6)' }}
            >
              {isAuthenticated ? 'Configure' : 'Admin Access'}
            </Button>
          </Group>
        </Container>
      </Box>

      <Center h="100vh">
        <Paper
          p="xl"
          radius="md"
          maw={400}
          w="100%"
          bg="dark.8"
          style={{ border: '1px solid var(--mantine-color-dark-7)' }}
        >
          <form onSubmit={handleLogin}>
            <Stack align="center" gap="md">
              <PresentationIcon size={48} color="var(--mantine-color-blue-4)" />
              <Title c="white">ClassDir</Title>
              <PasswordInput
                placeholder="Password"
                value={password}
                onChange={(e) => setPassword(e.currentTarget.value)}
                error={error}
                autoFocus
                w="100%"
              />
              <Button type="submit" loading={loggingIn} disabled={!password} fullWidth color="blue">
                Login
              </Button>
            </Stack>
          </form>
        </Paper>
      </Center>
    </Box>
  )
}
