import { useState } from "react";
import { useParams, useNavigate, Link } from "react-router";
import { Box, Button, Center, Group, Loader, Stack, Text, TextInput, Title } from "@mantine/core";
import { usePresentation } from "../hooks/use-presentation";
import { deletePresentation } from "../api";
import { DeleteModal } from "./delete-modal";
import { SlideEditor } from "./slide-editor";
import { CLIENT_CONFIGURE, clientControl } from "../../shared/cfg/routes";

export function PresentationDetailPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();

  const { presentation, isLoading, isSaving, error, saveContent } = usePresentation(id ?? "");

  const [editTitle, setEditTitle] = useState("");
  const [editing, setEditing] = useState(false);
  const [deleteModalOpen, setDeleteModalOpen] = useState(false);
  const [isDeleting, setIsDeleting] = useState(false);

  if (isLoading) {
    return (
      <Center h="100vh" bg="dark.9">
        <Loader />
      </Center>
    );
  }

  if (error) {
    return (
      <Center h="100vh" bg="dark.9">
        <Stack align="center">
          <Text c="red">{error.message}</Text>
          <Button component={Link} to={CLIENT_CONFIGURE}>
            Back to Presentations
          </Button>
        </Stack>
      </Center>
    );
  }

  if (!presentation) {
    return (
      <Center h="100vh" bg="dark.9">
        <Stack align="center">
          <Text>Presentation not found</Text>
          <Button component={Link} to={CLIENT_CONFIGURE}>
            Back to Presentations
          </Button>
        </Stack>
      </Center>
    );
  }

  function handleEditStart() {
    if (!presentation) return;
    setEditTitle(presentation.title);
    setEditing(true);
  }

  function handleSaveTitle() {
    if (!editTitle.trim() || !presentation) return;
    saveContent(editTitle.trim(), presentation.content)?.finally(() => setEditing(false));
  }

  function handleDeleteConfirm() {
    if (!presentation) return;
    setIsDeleting(true);
    deletePresentation(presentation.id).match(
      () => navigate("/configure"),
      () => {
        setIsDeleting(false);
        setDeleteModalOpen(false);
      },
    );
  }

  function handleSaveContent(content: string) {
    if (!presentation) return;
    saveContent(presentation.title, content);
  }

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
            style={{ borderColor: "var(--mantine-color-dark-6)" }}
          >
            &larr; Back
          </Button>
          <Button component={Link} to={clientControl(id!)} variant="light" size="sm">
            Control
          </Button>
        </Group>
        <Group>
          <Button color="red" onClick={() => setDeleteModalOpen(true)}>
            Delete Presentation
          </Button>
        </Group>
      </Group>

      {editing ? (
        <form
          onSubmit={(e) => {
            e.preventDefault();
            handleSaveTitle();
          }}
        >
          <Group>
            <TextInput
              value={editTitle}
              onChange={(e) => setEditTitle(e.currentTarget.value)}
              autoFocus
              style={{ flex: 1 }}
            />
            <Button type="submit" loading={isSaving} disabled={!editTitle.trim() || isSaving}>
              Save
            </Button>
            <Button variant="subtle" onClick={() => setEditing(false)}>
              Cancel
            </Button>
          </Group>
        </form>
      ) : (
        <Group justify="space-between">
          <Group>
            <Title c="white">{presentation.title}</Title>
            <Button
              variant="outline"
              color="gray"
              size="sm"
              style={{ borderColor: "var(--mantine-color-dark-6)" }}
              onClick={handleEditStart}
            >
              Edit
            </Button>
          </Group>
        </Group>
      )}

      <Box style={{ flex: 1, minHeight: 0 }}>
        <SlideEditor
          content={presentation.content}
          onSave={handleSaveContent}
          isSaving={isSaving}
        />
      </Box>

      <DeleteModal
        opened={deleteModalOpen}
        onClose={() => setDeleteModalOpen(false)}
        onConfirm={handleDeleteConfirm}
        title={presentation.title}
        isLoading={isDeleting}
      />
    </Stack>
  );
}
