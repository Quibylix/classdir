import { useEffect, useCallback } from 'react'
import { useParams } from 'react-router'
import { Box, Center, Container, Loader, Paper, Text } from '@mantine/core'
import { useSlideShow } from '../hooks/use-slide-show'
import { WS_CMD_JOIN_ROOM } from '../cfg'
import { visibleStrokes, drawStrokes } from '../utils/annotation-canvas'

export function PresentView() {
  const { code } = useParams<{ code: string }>()
  const { cachedHtml, slideCount, currentSlide, loading, fetchError, iframeRef, canvasRef, operationsBySlide } =
    useSlideShow(code ? { command: WS_CMD_JOIN_ROOM, parameters: { room_code: code } } : null)

  const triggerRedraw = useCallback(() => {
    const canvas = canvasRef.current
    if (!canvas) return
    const ctx = canvas.getContext('2d')
    if (!ctx) return

    const ops = operationsBySlide[String(currentSlide)] ?? []
    const strokes = visibleStrokes(ops)
    drawStrokes(ctx, strokes, undefined, canvas.width, canvas.height)
  }, [canvasRef, operationsBySlide, currentSlide])

  useEffect(() => {
    if (loading) return

    const canvas = canvasRef.current
    if (!canvas) return
    const parent = canvas.parentElement
    if (!parent) return

    const observer = new ResizeObserver(() => {
      const dpr = window.devicePixelRatio || 1
      canvas.width = parent.offsetWidth * dpr
      canvas.height = parent.offsetHeight * dpr

      triggerRedraw()
    })

    observer.observe(parent)
    return () => observer.disconnect()
  }, [loading, canvasRef, triggerRedraw])

  useEffect(() => {
    if (loading) return
    triggerRedraw()
  }, [loading, triggerRedraw])

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
