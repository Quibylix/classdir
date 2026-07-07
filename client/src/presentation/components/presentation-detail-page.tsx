import { useState } from 'react'
import { useParams, useNavigate, Link } from 'react-router'
import { Box, Button, Center, Group, Loader, Stack, Text, TextInput, Title } from '@mantine/core'
import { CaretLeftIcon } from '@phosphor-icons/react/dist/csr/CaretLeft'
import { CaretRightIcon } from '@phosphor-icons/react/dist/csr/CaretRight'
import { PlusIcon } from '@phosphor-icons/react/dist/csr/Plus'
import { TrashIcon } from '@phosphor-icons/react/dist/csr/Trash'
import { usePresentation } from '../hooks/use-presentation'
import { useSlides } from '../hooks/use-slides'
import { deletePresentation } from '../api'
import { DeleteModal } from './delete-modal'
import { SlideEditor } from './slide-editor'
import { CLIENT_CONFIGURE, clientPresent, clientControl } from '../../shared/cfg/routes'

export function PresentationDetailPage() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()

  const { presentation, isLoading, isSaving, error, updateTitle } = usePresentation(id ?? "")
  const {
    slides, currentIndex,
    isAdding, isSaving: isSlideSaving, isDeleting: isSlideDeleting,
    error: slidesError,
    addSlide, saveSlide, removeSlide, goToSlide,
  } = useSlides(id ?? "", presentation?.slides)

  const currentSlide = slides[currentIndex];

  const [editTitle, setEditTitle] = useState('')
  const [editing, setEditing] = useState(false)
  const [deleteModalOpen, setDeleteModalOpen] = useState(false)
  const [isDeleting, setIsDeleting] = useState(false)

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
          <Button component={Link} to={CLIENT_CONFIGURE}>Back to Presentations</Button>
        </Stack>
      </Center>
    )
  }

  if (!presentation) {
    return (
      <Center h="100vh" bg="dark.9">
        <Stack align="center">
          <Text>Presentation not found</Text>
          <Button component={Link} to={CLIENT_CONFIGURE}>Back to Presentations</Button>
        </Stack>
      </Center>
    )
  }

  function handleEditStart() {
    if (!presentation) return
    setEditTitle(presentation.title)
    setEditing(true)
  }

  function handleSaveTitle() {
    if (!editTitle.trim()) return
    updateTitle(editTitle.trim())?.finally(() => setEditing(false))
  }

  function handleDeleteConfirm() {
    if (!presentation) return
    setIsDeleting(true)
    deletePresentation(presentation.id).match(
      () => navigate('/configure'),
      () => { setIsDeleting(false); setDeleteModalOpen(false) },
    )
  }

  if (!presentation) return null

  return (
    <Stack h="100vh" p="md" gap="sm" bg="dark.9">
      <Group justify="space-between">
        <Group>
          <Button
            component={Link}
            to={CLIENT_CONFIGURE}
            variant="outline"
            color="gray"
            size="sm"
            style={{ borderColor: 'var(--mantine-color-dark-6)' }}
          >
            &larr; Back
          </Button>
          <Button component={Link} to={clientPresent(id!)} variant="light" size="sm">Present</Button>
          <Button component={Link} to={clientControl(id!)} variant="light" size="sm">Control</Button>
        </Group>
        <Group>
          <Button color="red" onClick={() => setDeleteModalOpen(true)}>Delete Presentation</Button>
        </Group>
      </Group>

      {editing ? (
        <form onSubmit={(e) => { e.preventDefault(); handleSaveTitle() }}>
          <Group>
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
        <Group justify="space-between">
          <Group>
            <Title c="white">{presentation.title}</Title>
            <Button variant="outline" color="gray" size="sm" style={{ borderColor: 'var(--mantine-color-dark-6)' }} onClick={handleEditStart}>Edit</Button>
          </Group>
        </Group>
      )}

      {slidesError && (
        <Text c="red" size="sm">{slidesError.message}</Text>
      )}

      <Group justify="space-between">
        <Group>
          <Button
            variant="outline"
            color="gray"
            size="sm"
            onClick={() => goToSlide(currentIndex - 1)}
            disabled={currentIndex <= 0}
            px="xs"
            style={{ borderColor: 'var(--mantine-color-dark-6)' }}
          >
            <CaretLeftIcon size={16} />
          </Button>
          <Text size="sm" c="dimmed">
            {slides.length > 0 ? `${currentIndex + 1} / ${slides.length}` : '0 / 0'}
          </Text>
          <Button
            variant="outline"
            color="gray"
            size="sm"
            onClick={() => goToSlide(currentIndex + 1)}
            disabled={currentIndex >= slides.length - 1}
            px="xs"
            style={{ borderColor: 'var(--mantine-color-dark-6)' }}
          >
            <CaretRightIcon size={16} />
          </Button>
        </Group>
        <Group>
          <Button
            leftSection={<PlusIcon size={16} />}
            onClick={addSlide}
            loading={isAdding}
            disabled={isAdding}
          >
            Add Slide
          </Button>
          <Button
            leftSection={<TrashIcon size={16} />}
            color="red"
            variant="outline"
            onClick={() => removeSlide(currentIndex)}
            loading={isSlideDeleting}
            disabled={!currentSlide || isSlideDeleting}
          >
            Delete Slide
          </Button>
        </Group>
      </Group>

      <Box style={{ flex: 1, minHeight: 0 }}>
        {currentSlide ? (
          <SlideEditor
            slides={slides}
            currentIndex={currentIndex}
            onSave={saveSlide}
            isSaving={isSlideSaving}
          />
        ) : (
          <Center h="100%">
            <Stack align="center" gap="md">
              <Text c="dimmed">No slides yet</Text>
              <Button onClick={addSlide} loading={isAdding} disabled={isAdding}>
                Add your first slide
              </Button>
            </Stack>
          </Center>
        )}
      </Box>

      <DeleteModal
        opened={deleteModalOpen}
        onClose={() => setDeleteModalOpen(false)}
        onConfirm={handleDeleteConfirm}
        title={presentation.title}
        isLoading={isDeleting}
      />
    </Stack>
  )
}
