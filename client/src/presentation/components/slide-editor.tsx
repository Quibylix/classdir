import { useRef, useEffect, useState, useCallback } from 'react'
import type { EditorView } from 'codemirror'
import { Box, Button, Group, Paper, Stack } from '@mantine/core'
import { CDN_REVEAL_CSS, CDN_REVEAL_THEME_CSS, CDN_REVEAL_JS } from '../cfg'
import type { Slide } from '../types'

type SlideEditorProps = {
  slides: Slide[]
  currentIndex: number
  onSave: (index: number, content: string) => void
  isSaving: boolean
}

function buildPreviewHtml(allSlides: Slide[], targetIndex: number): string {
  const slidesHtml = allSlides.map((s) => `<section>${s.content}</section>`).join('\n')
  return `<!DOCTYPE html>
<html>
<head>
  <link rel="stylesheet" href="${CDN_REVEAL_CSS}">
  <link rel="stylesheet" href="${CDN_REVEAL_THEME_CSS}">
</head>
<body>
  <div class="reveal" id="reveal" style="opacity: 0;">
    <div class="slides">
${slidesHtml}
    </div>
  </div>
  <script src="${CDN_REVEAL_JS}"></script>
  <script>
    Reveal.initialize({transition: 'none', progress: false}).then(() => {
      Reveal.slide(${targetIndex});
      Reveal.configure({ transition: 'slide'})
      document.getElementById('reveal').style.opacity = '1';
    });
  </script>
</body>
</html>`
}

export function SlideEditor({ slides, currentIndex, onSave, isSaving }: SlideEditorProps) {
  const editorRef = useRef<HTMLDivElement>(null)
  const viewRef = useRef<EditorView | null>(null)
  const previewRef = useRef<HTMLIFrameElement>(null)
  const [content, setContent] = useState(() => slides[currentIndex]?.content ?? '')

  const lastIndexRef = useRef(currentIndex)
  const currentPropContent = slides[currentIndex]?.content ?? ''
  const lastPropContentRef = useRef(currentPropContent)

  useEffect(() => {
    let isMounted = true;

    async function loadCodeMirror() {
      if (!editorRef.current) return;

      const { EditorView, basicSetup } = await import('codemirror');
      const { html } = await import('@codemirror/lang-html');

      if (!isMounted) return;

      if (editorRef.current) {
        editorRef.current.innerHTML = '';
      }

      const view = new EditorView({
        doc: content,
        extensions: [
          basicSetup,
          html(),
          EditorView.theme({}, { dark: true }),
          EditorView.updateListener.of((update) => {
            if (update.docChanged) {
              setContent(update.state.doc.toString());
            }
          }),
        ],
        parent: editorRef.current,
      });

      viewRef.current = view;
    }

    loadCodeMirror();

    return () => {
      isMounted = false;
      viewRef.current?.destroy();
      viewRef.current = null;
    };
  }, []);

  useEffect(() => {
    const view = viewRef.current
    if (!view) return

    const hasSlideChanged = currentIndex !== lastIndexRef.current
    const hasPropContentChanged = currentPropContent !== lastPropContentRef.current

    if (hasSlideChanged || hasPropContentChanged) {
      lastIndexRef.current = currentIndex
      lastPropContentRef.current = currentPropContent

      setContent(currentPropContent)

      if (view.state.doc.toString() !== currentPropContent) {
        view.dispatch({
          changes: { from: 0, to: view.state.doc.length, insert: currentPropContent },
        })
      }
    }
  }, [currentIndex, currentPropContent])

  const handleSave = useCallback(() => {
    onSave(currentIndex, content)
  }, [onSave, currentIndex, content])

  const previewHtml = buildPreviewHtml(slides, currentIndex)

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
        <Paper ref={editorRef} flex={1} withBorder mih={400} style={{ overflow: 'auto' }} />
        <Box flex={1} mih={400}>
          <iframe
            ref={previewRef}
            srcDoc={previewHtml}
            title="Slide Preview"
            style={{ width: '100%', height: '100%', border: '1px solid var(--mantine-color-default-border)', borderRadius: 'var(--mantine-radius-md)' }}
          />
        </Box>
      </Group>
    </Stack>
  )
}
