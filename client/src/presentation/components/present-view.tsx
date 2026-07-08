import { useEffect, useRef, useState } from 'react'
import { useParams } from 'react-router'
import { Box, Center, Container, Loader, Paper, Text } from '@mantine/core'
import { useSafeWebSocket } from '../../shared/hooks/use-websocket'
import { WSOutputMessageSchema, type Slide, type AnnotationOperation } from '../types'
import { WS_V1 } from '../../shared/cfg/routes'
import { WS_STATUS } from '../../shared/types'
import { POST_MSG_TYPE, WS_CMD_JOIN_ROOM, WS_EVENT_ANNOTATION_ADDED, CDN_REVEAL_CSS, CDN_REVEAL_THEME_CSS, CDN_REVEAL_JS, WS_EVENT_SLIDE_CHANGED } from '../cfg'
import { visibleStrokes, drawStrokes } from '../utils/annotation-canvas'

function buildPresentHtml(slides: Slide[], initialSlide: number): string {
  const slidesHtml = slides.map(s => `<section>${s.content}</section>`).join('\n')
  return `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <link rel="stylesheet" href="${CDN_REVEAL_CSS}">
  <link rel="stylesheet" href="${CDN_REVEAL_THEME_CSS}">
  <style>
    * { margin: 0; padding: 0; box-sizing: border-box; }
    html, body { width: 100%; height: 100%; overflow: hidden; }
  </style>
</head>
<body>
  <div class="reveal" id="reveal">
    <div class="slides">${slidesHtml}</div>
  </div>
  <script src="${CDN_REVEAL_JS}"></script>
  <script>
    Reveal.initialize({ transition: 'slide', progress: false, controls: false, touch: false, scrollActivationWidth: null, }).then(function() {
      Reveal.slide(${initialSlide});
    });
    window.addEventListener('message', function(e) {
      if (e.data.type === '${POST_MSG_TYPE.Navigate}') Reveal.slide(e.data.index);
    });
  </script>
</body>
</html>`
}

export function PresentView() {
  const { code } = useParams<{ code: string }>()
  const iframeRef = useRef<HTMLIFrameElement>(null)
  const canvasRef = useRef<HTMLCanvasElement>(null)
  const [cachedHtml, setCachedHtml] = useState<string>()
  const [slideCount, setSlideCount] = useState<number>(0)
  const [currentSlide, setCurrentSlide] = useState(0)
  const [loading, setLoading] = useState(true)
  const [fetchError, setFetchError] = useState<string | null>(null)
  const [operationsBySlide, setOperationsBySlide] = useState<Record<string, AnnotationOperation[]>>({})
  const joinedRef = useRef(false)

  const wsUrl = `${window.location.protocol === 'https:' ? 'wss:' : 'ws:'}//${window.location.host}${WS_V1}`

  const { status, send } = useSafeWebSocket({
    url: wsUrl,
    schema: WSOutputMessageSchema,
    onMessage(msg) {
      if ('error' in msg) {
        setFetchError(msg.error.message)
        setLoading(false)
        return
      }

      if (!('event' in msg)) {
        setCachedHtml(buildPresentHtml(msg.data.slides, msg.data.current_index))
        setSlideCount(msg.data.slides.length)
        setCurrentSlide(msg.data.current_index)
        setLoading(false)
        return;
      }

      if (msg.event === WS_EVENT_SLIDE_CHANGED) {
        setCurrentSlide(msg.data.current_slide)
        iframeRef.current?.contentWindow?.postMessage({ type: POST_MSG_TYPE.Navigate, index: msg.data.current_slide }, window.location.origin)
        return;
      }

      if (msg.event === WS_EVENT_ANNOTATION_ADDED) {
        setOperationsBySlide(prev => {
          const slideKey = String(currentSlide)
          const existing = prev[slideKey] ?? []
          if (existing.some(op => op.id === msg.data.id)) return prev
          return { ...prev, [slideKey]: [...existing, msg.data] }
        })
        return;
      }

      setOperationsBySlide(msg.data.operations_by_slide)
    }
  })

  useEffect(() => {
    if (status === WS_STATUS.Connected && code && !joinedRef.current) {
      send({ command: WS_CMD_JOIN_ROOM, parameters: { room_code: code } })
      joinedRef.current = true
    }
    if (status === WS_STATUS.Disconnected) {
      joinedRef.current = false
    }
  }, [status, code, send])

  useEffect(() => {
    const canvas = canvasRef.current
    if (!canvas) return

    const parent = canvas.parentElement
    if (!parent) return

    const observer = new ResizeObserver(() => {
      canvas.width = parent.offsetWidth
      canvas.height = parent.offsetHeight
    })
    observer.observe(parent)
    return () => observer.disconnect()
  }, [])

  useEffect(() => {
    const canvas = canvasRef.current
    if (!canvas) return
    const ctx = canvas.getContext('2d')
    if (!ctx) return

    const ops = operationsBySlide[String(currentSlide)] ?? []
    const strokes = visibleStrokes(ops)
    drawStrokes(ctx, strokes, undefined, canvas.width, canvas.height)
  }, [operationsBySlide, currentSlide])

  if (loading) {
    return <Center h="100vh" bg="dark.9"><Loader /></Center>
  }

  if (fetchError) {
    return (
      <Center h="100vh" bg="dark.9">
        <Text c="red">{fetchError}</Text>
      </Center>
    )
  }

  if (slideCount === 0) {
    return (
      <Center h="100vh" bg="dark.9">
        <Text c="dimmed">No slides in this presentation</Text>
      </Center>
    )
  }

  return (
    <Container fluid m={0} p={0} h="100vh" bg="dark.9" pos="relative">
      <Box m="auto" inset={0} p={0} mah="100%" maw="100%" pos="absolute" bg="#000" style={{ aspectRatio: '48/35' }}>
        <iframe
          ref={iframeRef}
          srcDoc={cachedHtml}
          title="Presentation"
          style={{ width: '100%', height: '100%', border: 'none', display: 'block' }}
        />
        <canvas
          ref={canvasRef}
          style={{
            position: 'absolute',
            top: 0,
            left: 0,
            width: '100%',
            height: '100%',
            pointerEvents: 'none',
          }}
        />
        <Paper
          pos="fixed"
          bottom={16}
          right={16}
          p="4px 12px"
          bg="dark.8"
          bd="1px solid dark.6"
          style={{ borderRadius: 4 }}
        >
          <Text c="gray.3" fz={14} ff="monospace">
            {currentSlide + 1} / {slideCount}
          </Text>
        </Paper>
      </Box>
    </Container>
  )
}
