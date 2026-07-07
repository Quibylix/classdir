import { useState } from 'react'
import { useNavigate } from 'react-router'
import { Button, Center, Container, Paper, Stack, Text, TextInput, Title } from '@mantine/core'
import { CLIENT_PRESENT } from '../../shared/cfg/routes'

export function PresentCodeEntry() {
  const [code, setCode] = useState('')
  const navigate = useNavigate()

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    if (code.trim().length === 8) {
      navigate(`${CLIENT_PRESENT}/${code.trim()}`)
    }
  }

  return (
    <Center h="100vh" bg="dark.9">
      <Container size="xs">
        <Paper bg="dark.8" p="xl" radius="md" bd="1px solid dark.6">
          <form onSubmit={handleSubmit}>
            <Stack align="center" gap="lg">
              <Title order={2} c="white" ta="center">Enter Presentation Code</Title>
              <Text c="dimmed" ta="center" size="sm">
                Enter the 8-digit code shared by your teacher
              </Text>
              <TextInput
                placeholder="00000000"
                value={code}
                onChange={(e) => setCode(e.currentTarget.value.replace(/\D/g, '').slice(0, 8))}
                size="xl"
                styles={{ input: { textAlign: 'center', letterSpacing: '0.3em', fontSize: '1.5rem', fontFamily: 'monospace' } }}
                autoFocus
              />
              <Button type="submit" size="lg" fullWidth disabled={code.length !== 8}>
                Join Presentation
              </Button>
            </Stack>
          </form>
        </Paper>
      </Container>
    </Center>
  )
}
