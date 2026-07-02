import { useState } from 'react'
import { useParams, useNavigate, Link } from 'react-router'
import { Box, Button, Card, Center, Group, Loader, Stack, Text, TextInput, Title } from '@mantine/core'
import { usePresentation } from '../hooks/usePresentation'
import { deletePresentation } from '../api'
import { DeleteModal } from './DeleteModal'

export function PresentationDetailPage() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()

  const { presentation, isLoading, isSaving, error, updateTitle } = usePresentation(id ?? "")
  const [editTitle, setEditTitle] = useState('')
  const [editing, setEditing] = useState(false)
  const [deleteModalOpen, setDeleteModalOpen] = useState(false)
  const [isDeleting, setIsDeleting] = useState(false)

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
          <Button component={Link} to="/configure">Back to Presentations</Button>
        </Stack>
      </Center>
    )
  }

  if (!presentation) {
    return (
      <Center h="100vh">
        <Stack align="center">
          <Text>Presentation not found</Text>
          <Button component={Link} to="/configure">Back to Presentations</Button>
        </Stack>
      </Center>
    )
  }

  function handleEditStart() {
    if (!presentation) return
    setEditTitle(presentation.title)
    setEditing(true)
  }

  function handleSave() {
    if (!editTitle.trim()) return
    updateTitle(editTitle.trim())
    setEditing(false)
  }

  function handleDeleteConfirm() {
    if (!presentation) return
    setIsDeleting(true)
    deletePresentation(presentation.id).match(
      () => navigate('/configure'),
      () => { setIsDeleting(false); setDeleteModalOpen(false) },
    )
  }

  return (
    <Box p="xl" maw={800} mx="auto">
      <Group justify="space-between" mb="lg">
        <Button component={Link} to="/configure" variant="subtle">
          &larr; Back
        </Button>
        <Button color="red" onClick={() => setDeleteModalOpen(true)}>Delete Presentation</Button>
      </Group>

      {editing ? (
        <form onSubmit={(e) => { e.preventDefault(); handleSave() }}>
          <Group mb="lg">
            <TextInput
              value={editTitle}
              onChange={(e) => setEditTitle(e.currentTarget.value)}
              autoFocus
              style={{ flex: 1 }}
            />
            <Button type="submit" loading={isSaving} disabled={!editTitle.trim() || isSaving}>Save</Button>
            <Button variant="subtle" onClick={() => setEditing(false)}>Cancel</Button>
          </Group>
        </form>
      ) : (
        <Group mb="lg" justify="space-between">
          <Title>{presentation.title}</Title>
          <Button variant="subtle" onClick={handleEditStart}>Edit</Button>
        </Group>
      )}

      <DeleteModal
        opened={deleteModalOpen}
        onClose={() => setDeleteModalOpen(false)}
        onConfirm={handleDeleteConfirm}
        title={presentation.title}
        isLoading={isDeleting}
      />

      <Title order={3} mb="md">Slides</Title>
      {presentation.slides.length === 0 ? (
        <Text c="dimmed">No slides yet</Text>
      ) : (
        <Stack>
          {presentation.slides.map((slide) => (
            <Card key={slide.id} shadow="sm" padding="md" radius="md">
              <Text fw={500} size="sm" c="dimmed" mb="xs">
                Slide {slide.slide_number}
              </Text>
              <Text>{slide.content}</Text>
            </Card>
          ))}
        </Stack>
      )}
    </Box>
  )
}
