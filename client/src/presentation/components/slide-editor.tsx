import { useRef, useEffect, useMemo } from "react";
import { Box, Button, Group, Paper, Stack } from "@mantine/core";
import { buildPreviewHtml } from "../utils/reveal-html";
import { useCodeMirrorEditor } from "../hooks/use-codemirror-editor";

type SlideEditorProps = {
  content: string;
  onSave: (content: string) => void;
  isSaving: boolean;
};

export function SlideEditor({ content, onSave, isSaving }: SlideEditorProps) {
  const { editorRef, doc, setContent, setReadOnly } = useCodeMirrorEditor(content, isSaving);

  const lastPropContentRef = useRef(content);

  useEffect(() => {
    if (content === lastPropContentRef.current) return;
    lastPropContentRef.current = content;

    setContent(content);
  }, [content, setContent]);

  useEffect(() => {
    setReadOnly(isSaving);
  }, [isSaving, setReadOnly]);

  const handleSave = () => onSave(doc);

  const previewHtml = useMemo(() => {
    const slides = content.split(/^---+\s*\n/m).map((p) => p.replace(/^\n+/, ""));
    return buildPreviewHtml(slides, 0);
  }, [content]);

  return (
    <Stack h="100%">
      <Group justify="space-between">
        <Group>
          <Button onClick={handleSave} loading={isSaving} disabled={isSaving}>
            Save
          </Button>
        </Group>
      </Group>
      <Group align="stretch" h="100%" gap="md">
        <Paper
          ref={editorRef}
          flex={1}
          withBorder
          h="100%"
          style={{
            overflow: "hidden",
            opacity: isSaving ? 0.6 : 1,
            pointerEvents: isSaving ? "none" : "auto",
          }}
        />
        <Box flex={1} h="100%">
          <iframe
            srcDoc={previewHtml}
            title="Slide Preview"
            style={{
              width: "100%",
              height: "100%",
              border: "1px solid var(--mantine-color-default-border)",
              borderRadius: "var(--mantine-radius-md)",
            }}
          />
        </Box>
      </Group>
    </Stack>
  );
}
