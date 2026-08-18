import { useState } from "react";
import {
  ActionIcon,
  Button,
  Collapse,
  Group,
  Loader,
  Paper,
  ScrollArea,
  Stack,
  Text,
  TextInput,
} from "@mantine/core";
import { PencilIcon } from "@phosphor-icons/react/dist/csr/Pencil";
import { TrashIcon } from "@phosphor-icons/react/dist/csr/Trash";
import { useStudents } from "../hooks/use-students";
import { DeleteModal } from "./delete-modal";

type StudentsPanelProps = {
  presentationId: string;
};

export function StudentsPanel({ presentationId }: StudentsPanelProps) {
  const { students, isLoading, isCreating, isUpdating, isDeleting, error, create, rename, remove } =
    useStudents(presentationId);

  const [open, setOpen] = useState(false);
  const [newName, setNewName] = useState("");
  const [editingId, setEditingId] = useState<string | null>(null);
  const [editName, setEditName] = useState("");
  const [deletingId, setDeletingId] = useState<string | null>(null);

  const deletingName = deletingId ? (students.find((s) => s.id === deletingId)?.name ?? "") : "";

  function handleAdd() {
    const name = newName.trim();
    if (!name || isCreating) return;
    create(name);
    setNewName("");
  }

  function handleRenameSave() {
    if (!editingId || isUpdating) return;
    const name = editName.trim();
    if (!name) return;
    rename(editingId, name)?.finally(() => setEditingId(null));
  }

  function handleDeleteConfirm() {
    if (!deletingId || isDeleting) return;
    remove(deletingId)?.finally(() => setDeletingId(null));
  }

  return (
    <Paper withBorder bg="dark.8" p="sm">
      <Stack gap="sm">
        <Button variant="subtle" color="gray" onClick={() => setOpen((v) => !v)}>
          {open ? "Hide students" : "Show students"} ({students.length})
        </Button>

        <Collapse expanded={open}>
          <Stack gap="sm">
            {isLoading ? (
              <Loader size="sm" />
            ) : (
              <>
                <Group>
                  <form
                    onSubmit={(e) => {
                      e.preventDefault();
                      handleAdd();
                    }}
                  >
                    <TextInput
                      placeholder="Student name"
                      value={newName}
                      onChange={(e) => setNewName(e.currentTarget.value)}
                      flex={1}
                    />

                    <Button
                      type="submit"
                      disabled={!newName.trim() || isCreating}
                      loading={isCreating}
                    >
                      Add
                    </Button>
                  </form>
                </Group>

                {students.length === 0 ? (
                  <Text c="dimmed" size="sm">
                    No students registered yet
                  </Text>
                ) : (
                  <ScrollArea.Autosize mah={280}>
                    <Stack gap={4}>
                      {students.map((s) => (
                        <Group key={s.id} justify="space-between" wrap="nowrap" gap="xs">
                          {editingId === s.id ? (
                            <form
                              onSubmit={(e) => {
                                e.preventDefault();
                                handleRenameSave();
                              }}
                            >
                              <TextInput
                                placeholder="New name"
                                value={editName}
                                onChange={(e) => setEditName(e.currentTarget.value)}
                                style={{ flex: 1 }}
                              />
                              <Button
                                type="submit"
                                size="xs"
                                loading={isUpdating}
                                disabled={!editName.trim() || isUpdating}
                              >
                                Save
                              </Button>
                              <Button onClick={() => setEditingId(null)} size="xs" variant="subtle">
                                Cancel
                              </Button>
                            </form>
                          ) : (
                            <>
                              <Text size="sm" style={{ flex: 1 }}>
                                {s.name}
                              </Text>
                              <Group gap={4}>
                                <ActionIcon
                                  size="sm"
                                  variant="outline"
                                  color="gray"
                                  aria-label="Edit student"
                                  onClick={() => {
                                    setEditingId(s.id);
                                    setEditName(s.name);
                                  }}
                                >
                                  <PencilIcon size={14} />
                                </ActionIcon>
                                <ActionIcon
                                  size="sm"
                                  variant="outline"
                                  color="red"
                                  aria-label="Delete student"
                                  onClick={() => setDeletingId(s.id)}
                                >
                                  <TrashIcon size={14} />
                                </ActionIcon>
                              </Group>
                            </>
                          )}
                        </Group>
                      ))}
                    </Stack>
                  </ScrollArea.Autosize>
                )}

                {error && (
                  <Text c="red" size="sm">
                    {error.message}
                  </Text>
                )}
              </>
            )}
          </Stack>
        </Collapse>
      </Stack>

      <DeleteModal
        opened={deletingId !== null}
        onClose={() => setDeletingId(null)}
        onConfirm={handleDeleteConfirm}
        title={deletingName}
        heading="Delete student"
        isLoading={isDeleting}
      />
    </Paper>
  );
}
