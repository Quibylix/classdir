import { useRef, useEffect, useCallback, useMemo } from "react";
import { Box, Button, Group, Paper, Stack } from "@mantine/core";
import { buildPreviewHtml, splitSlides } from "../utils/reveal-html";
import { POST_MSG_TYPE, SLIDE_PREVIEW_DEBOUNCE_MS } from "../cfg";
import { useCodeMirrorEditor } from "../hooks/use-codemirror-editor";

type SlideEditorProps = {
  content: string;
  onSave: (content: string) => void;
  isSaving: boolean;
};

export function SlideEditor({ content, onSave, isSaving }: SlideEditorProps) {
  const { editorRef, doc, cursorSlide, setContent, setReadOnly } = useCodeMirrorEditor(
    content,
    isSaving,
  );

  const iframeRef = useRef<HTMLIFrameElement>(null);

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

  const previewHtml = useMemo(() => buildPreviewHtml(splitSlides(content), 0), [content]);

  useEffect(() => {
    const t = setTimeout(() => {
      iframeRef.current?.contentWindow?.postMessage(
        { type: POST_MSG_TYPE.UpdateSlides, slides: splitSlides(doc), index: cursorSlide },
        window.location.origin,
      );
    }, SLIDE_PREVIEW_DEBOUNCE_MS);
    return () => clearTimeout(t);
  }, [doc, cursorSlide]);

  const handleIframeLoad = useCallback(() => {
    iframeRef.current?.contentWindow?.postMessage(
      { type: POST_MSG_TYPE.UpdateSlides, slides: splitSlides(doc), index: cursorSlide },
      window.location.origin,
    );
  }, [doc, cursorSlide]);

  return (
    <Stack h="100%">
      <Group justify="space-between">
        <Group>
          <Button onClick={handleSave} loading={isSaving} disabled={isSaving}>
            Save
          </Button>
        </Group>
      </Group>
      <Group align="stretch" mih={0} gap="md">
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
            ref={iframeRef}
            srcDoc={previewHtml}
            title="Slide Preview"
            onLoad={handleIframeLoad}
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
