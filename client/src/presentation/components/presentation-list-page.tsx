import { useState } from 'react'
import { Link } from 'react-router'
import { Box, Button, Card, Center, Group, Loader, Modal, SimpleGrid, Stack, Text, TextInput, Title } from '@mantine/core'
import { usePresentationList } from '../hooks/use-presentation-list'
import { useAuth } from '../../auth/hooks/use-auth'
import { DeleteModal } from './delete-modal'
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
      <Center h="100vh">
        <Loader />
      </Center>
    )
  }

  if (error) {
    return (
      <Center h="100vh">
        <Stack align="center">
          <Text c="red">{error.message}</Text>
          <Button onClick={refresh}>Retry</Button>
        </Stack>
      </Center>
    )
  }

  return (
    <Box p="xl">
      <Group justify="space-between" mb="lg">
        <Title>ClassDir</Title>
        <Button variant="subtle" onClick={logout}>Logout</Button>
      </Group>

      <Group justify="space-between" mb="md">
        <Title order={2}>Presentations</Title>
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
              shadow="sm"
              padding="lg"
              radius="md"
            >
              <Stack gap="xs">
                <Text className={styles.titleLink} component={Link} to={`/configure/${p.id}`} fw={500}>{p.title}</Text>
                <Group gap="xs">
                  <Button component={Link} to={`/present/${p.id}`} size="xs" variant="light">Present</Button>
                  <Button component={Link} to={`/control/${p.id}`} size="xs" variant="light">Control</Button>
                  <Button
                    variant="subtle"
                    color="red"
                    size="xs"
                    ml="auto"
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
