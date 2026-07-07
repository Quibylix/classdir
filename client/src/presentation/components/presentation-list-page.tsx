import { useState } from 'react'
import { Link } from 'react-router'
import { Box, Button, Card, Center, Container, Group, Loader, Modal, SimpleGrid, Stack, Text, TextInput, Title } from '@mantine/core'
import { TrashIcon } from '@phosphor-icons/react/dist/csr/Trash'
import { usePresentationList } from '../hooks/use-presentation-list'
import { useAuth } from '../../auth/hooks/use-auth'
import { DeleteModal } from './delete-modal'
import { clientConfigure, clientControl, CLIENT_LOGIN } from '../../shared/cfg/routes'
import styles from './presentation-list-page.module.css'

export function PresentationListPage() {
  const { presentations, isLoading, isCreating, isDeleting, error, refresh, create, remove } = usePresentationList()
  const { logout } = useAuth()
  const [createModalOpen, setCreateModalOpen] = useState(false)
  const [createTitle, setCreateTitle] = useState('')
  const [deletingId, setDeletingId] = useState<string | null>(null)

  function handleCreate() {
    if (!createTitle.trim()) return
    create(createTitle.trim())
    setCreateTitle('')
    setCreateModalOpen(false)
  }

  const deletingTitle = deletingId ? presentations.find(p => p.id === deletingId)?.title ?? '' : ''

  if (isLoading) {
    return (
      <Center h="100vh" bg="dark.9">
        <Loader />
      </Center>
    )
  }

  if (error) {
    return (
      <Center h="100vh" bg="dark.9">
        <Stack align="center">
          <Text c="red">{error.message}</Text>
          <Button onClick={refresh}>Retry</Button>
        </Stack>
      </Center>
    )
  }

  return (
    <Box bg="dark.9" mih="100vh">
      <Box
        component="header"
        pos="fixed"
        top={0}
        left={0}
        right={0}
        style={{
          zIndex: 40,
          background: 'color-mix(in srgb, var(--mantine-color-dark-9) 80%, transparent)',
          backdropFilter: 'blur(12px)',
          borderBottom: '1px solid var(--mantine-color-dark-8)',
        }}
      >
        <Container size="lg">
          <Group h={64} justify="space-between">
            <Group gap={6}>
              <Center
                w={32}
                h={32}
                bg="blue.6"
                style={{ borderRadius: 6, boxShadow: '0 4px 12px color-mix(in srgb, var(--mantine-color-blue-5) 20%, transparent)' }}
              >
                <Text c="white" fw={700} size="sm">C</Text>
              </Center>
              <Text fw={700} size="lg" c="white" lts="-0.02em">ClassDir</Text>
            </Group>
            <Group gap="md">
              <Button
                variant="outline"
                color="gray"
                size="sm"
                style={{ borderColor: 'var(--mantine-color-dark-6)' }}
                onClick={() => { logout(); window.location.href = CLIENT_LOGIN }}
              >
                Logout
              </Button>
            </Group>
          </Group>
        </Container>
      </Box>

      <Container size="lg" pt={80} pb="xl">
        <Group justify="space-between" mb="lg">
          <Title order={2} c="white">Presentations</Title>
          <Button onClick={() => setCreateModalOpen(true)}>New Presentation</Button>
        </Group>

        {presentations.length === 0 ? (
          <Center py="xl">
            <Stack align="center">
              <Text c="dimmed">No presentations yet</Text>
              <Button onClick={() => setCreateModalOpen(true)}>Create your first presentation</Button>
            </Stack>
          </Center>
        ) : (
          <SimpleGrid cols={{ base: 1, sm: 2, md: 3 }}>
            {presentations.map((p) => (
              <Card
                className={styles.card}
                key={p.id}
                bg="dark.8"
                padding="lg"
                radius="md"
              >
                <Stack gap="xs">
                  <Text className={styles.titleLink} component={Link} to={clientConfigure(p.id)} fw={500}>{p.title}</Text>
                  <Group gap="xs">
                    <Button component={Link} to={clientControl(p.id)} size="xs" variant="light">Control</Button>
                    <Button
                      variant="light"
                      color="red"
                      size="xs"
                      ml="auto"
                      leftSection={<TrashIcon size={12} />}
                      onClick={(e) => { e.stopPropagation(); setDeletingId(p.id) }}
                    >
                      Delete
                    </Button>
                  </Group>
                </Stack>
              </Card>
            ))}
          </SimpleGrid>
        )}
      </Container>

      <DeleteModal
        opened={deletingId !== null}
        onClose={() => setDeletingId(null)}
        onConfirm={() => { if (deletingId) remove(deletingId); setDeletingId(null) }}
        title={deletingTitle}
        isLoading={isDeleting}
      />

      <Modal opened={createModalOpen} onClose={() => setCreateModalOpen(false)} title="New Presentation">
        <form onSubmit={(e) => { e.preventDefault(); handleCreate() }}>
          <Stack>
            <TextInput
              placeholder="Presentation title"
              value={createTitle}
              onChange={(e) => setCreateTitle(e.currentTarget.value)}
              autoFocus
            />
            <Button type="submit" loading={isCreating} disabled={!createTitle.trim() || isCreating}>Create</Button>
          </Stack>
        </form>
      </Modal>
    </Box>
  )
}
