import { Button, Group, Modal, Stack, Text } from '@mantine/core'

type DeleteModalProps = {
  opened: boolean
  onClose: () => void
  onConfirm: () => void
  title: string
  isLoading?: boolean
}

export function DeleteModal({ opened, onClose, onConfirm, title, isLoading }: DeleteModalProps) {
  return (
    <Modal opened={opened} onClose={onClose} title="Delete presentation">
      <Stack>
        <Text>
          Are you sure you want to delete "{title}"? This action cannot be undone.
        </Text>
        <Group justify="flex-end">
          <Button variant="subtle" onClick={onClose}>Cancel</Button>
          <Button color="red" onClick={onConfirm} loading={isLoading}>Delete</Button>
        </Group>
      </Stack>
    </Modal>
  )
}
