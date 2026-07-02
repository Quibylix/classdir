import { useRef, useEffect, useState, useCallback } from 'react'
import { EditorView, basicSetup } from 'codemirror'
import { html } from '@codemirror/lang-html'
import { Box, Button, Group, Paper, Stack } from '@mantine/core'
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
  <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/reveal.js@5/dist/reveal.css">
  <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/reveal.js@5/dist/theme/black.css">
</head>
<body>
  <div class="reveal">
    <div class="slides">
${slidesHtml}
    </div>
  </div>
  <script src="https://cdn.jsdelivr.net/npm/reveal.js@5/dist/reveal.js"><\/script>
  <script>
    Reveal.initialize({ embedded: true }).then(() => {
      Reveal.slide(${targetIndex});
    });
  <\/script>
</body>
</html>`
}

export function SlideEditor({ slides, currentIndex, onSave, isSaving }: SlideEditorProps) {
  const editorRef = useRef<HTMLDivElement>(null)
  const viewRef = useRef<EditorView | null>(null)
  const previewRef = useRef<HTMLIFrameElement>(null)
  const [content, setContent] = useState(() => slides[currentIndex]?.content ?? '')
  const [previewKey, setPreviewKey] = useState(0)

  useEffect(() => {
    if (!editorRef.current) return

    const view = new EditorView({
      doc: content,
      extensions: [
        basicSetup,
        html(),
        EditorView.updateListener.of((update) => {
          if (update.docChanged) {
            setContent(update.state.doc.toString())
          }
        }),
      ],
      parent: editorRef.current,
    })
    viewRef.current = view

    return () => view.destroy()
  }, [])

  useEffect(() => {
    const view = viewRef.current
    if (!view) return
    const newContent = slides[currentIndex]?.content ?? ''
    if (view.state.doc.toString() !== newContent) {
      view.dispatch({
        changes: { from: 0, to: view.state.doc.length, insert: newContent },
      })
    }
  }, [currentIndex, slides])

  const handleSave = useCallback(() => {
    onSave(currentIndex, content)
  }, [onSave, currentIndex, content])

  const handleRefresh = useCallback(() => {
    setPreviewKey((k) => k + 1)
  }, [])

  const previewHtml = buildPreviewHtml(slides, currentIndex)

  return (
    <Stack h="100%">
      <Group justify="space-between">
        <Group>
          <Button onClick={handleSave} loading={isSaving} disabled={isSaving}>
            Save
          </Button>
          <Button variant="subtle" onClick={handleRefresh}>
            Refresh Preview
          </Button>
        </Group>
      </Group>
      <Group align="stretch" h="100%" gap="md">
        <Paper ref={editorRef} flex={1} withBorder style={{ overflow: 'auto', minHeight: 400 }} />
        <Box flex={1} style={{ minHeight: 400 }}>
          <iframe
            key={previewKey}
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
