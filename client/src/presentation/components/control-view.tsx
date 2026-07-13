import { useEffect, useState } from 'react'
import { useParams, Link } from 'react-router'
import {
  Box, Button, Center, Code, Group, Loader, NumberInput, Slider, Stack, Text, ActionIcon, Paper, Container, Modal,
} from '@mantine/core'
import { PencilIcon } from '@phosphor-icons/react/dist/csr/Pencil'
import { EyeIcon } from '@phosphor-icons/react/dist/csr/Eye'
import { CodeSimpleIcon } from '@phosphor-icons/react/dist/csr/CodeSimple'
import { TrashIcon } from '@phosphor-icons/react/dist/csr/Trash'
import { useSlideShow } from '../hooks/use-slide-show'
import { useAnnotation } from '../hooks/use-annotation'
import { useSlideGestures } from '../hooks/use-slide-gestures'
import { DraggableFab } from './draggable-fab'
import { CLIENT_CONFIGURE } from '../../shared/cfg/routes'
import { WS_CMD_INIT_PRESENTATION, ANNOTATION_COLORS, ANNOTATION_MIN_THICKNESS, ANNOTATION_MAX_THICKNESS } from '../cfg'
import { visibleStrokes, drawStrokes } from '../utils/annotation-canvas'

export function ControlView() {
  const { id } = useParams<{ id: string }>()
  const { send, cachedHtml, slideCount, currentSlide, loading, fetchError, roomCode, iframeRef, canvasRef, operationsBySlide, setOperationsBySlide } =
    useSlideShow(id ? { command: WS_CMD_INIT_PRESENTATION, parameters: { presentation_id: id } } : null)
  const { drawMode, setDrawMode, annotationColor, setAnnotationColor, annotationThickness, setAnnotationThickness, currentPoints, handlePointerDown, handlePointerMove, handlePointerUp, handleClear } =
    useAnnotation({ send, currentSlide, setOperationsBySlide, canvasRef })
  const { goToValue, setGoToValue, goToModalOpen, setGoToModalOpen, handleZoneClick, handleGoToSubmit } =
    useSlideGestures({ send, slideCount, currentSlide })

  const [drawMenuOpen, setDrawMenuOpen] = useState(false)
  const [showRoomCode, setShowRoomCode] = useState(false)

  useEffect(() => {
    const canvas = canvasRef.current
    const container = canvas?.parentElement
    if (!container || !canvas) return

    const observer = new ResizeObserver(() => {
      canvas.width = container.offsetWidth
      canvas.height = container.offsetHeight
    })
    observer.observe(container)
    return () => observer.disconnect()
  }, [])

  useEffect(() => {
    const canvas = canvasRef.current
    if (!canvas) return
    const ctx = canvas.getContext('2d')
    if (!ctx) return

    const ops = operationsBySlide[String(currentSlide)] ?? []
    const strokes = visibleStrokes(ops)
    drawStrokes(ctx, strokes, currentPoints, canvas.width, canvas.height)
  }, [operationsBySlide, currentSlide, currentPoints])

  if (loading) {
    return <Center h="100vh" bg="dark.9"><Loader /></Center>
  }

  if (fetchError) {
    return (
      <Center h="100vh" bg="dark.9">
        <Stack align="center">
          <Text c="red">{fetchError}</Text>
          <Button component={Link} to={CLIENT_CONFIGURE}>Back</Button>
        </Stack>
      </Center>
    )
  }

  if (slideCount === 0) {
    return (
      <Center h="100vh" bg="dark.9">
        <Stack align="center">
          <Text c="dimmed">No slides in this presentation</Text>
          <Button component={Link} to={CLIENT_CONFIGURE}>Back</Button>
        </Stack>
      </Center>
    )
  }

  return (
    <Container fluid m={0} p={0} h="100vh" bg="dark.9" pos="relative">
      <Box m="auto" inset={0} mah="100%" maw="100%" pos="absolute" bg="#000" style={{ aspectRatio: '48/35' }}>
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
            pointerEvents: drawMode ? 'auto' : 'none',
            cursor: drawMode ? 'crosshair' : 'default',
            touchAction: drawMode ? 'none' : 'auto',
          }}
          onPointerDown={drawMode ? handlePointerDown : undefined}
          onPointerMove={drawMode ? handlePointerMove : undefined}
          onPointerUp={drawMode ? handlePointerUp : undefined}
        />
        {!drawMode && (
          <Box pos="absolute" inset={0} style={{ display: 'flex', cursor: 'pointer' }}>
            <Box
              flex="0 0 35%"
              onClick={() => handleZoneClick('prev')}
            />
            <Box
              flex="0 0 30%"
              onClick={() => handleZoneClick('goto')}
            />
            <Box
              flex="0 0 35%"
              onClick={() => handleZoneClick('next')}
            />
          </Box>
        )}
      </Box>

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

      {drawMode && drawMenuOpen && (
        <Paper pos="fixed" top={80} right={24} p="sm" withBorder bg="dark.8" style={{ zIndex: 1001 }}>
          <Stack gap="xs">
            <Group gap={4}>
              {ANNOTATION_COLORS.map(c => (
                <ActionIcon
                  key={c}
                  size="sm"
                  style={{
                    backgroundColor: c,
                    border: c === annotationColor ? '2px solid white' : '2px solid transparent',
                    borderRadius: '50%',
                  }}
                  onClick={() => setAnnotationColor(c)}
                />
              ))}
            </Group>
            <Slider
              value={annotationThickness}
              onChange={setAnnotationThickness}
              min={ANNOTATION_MIN_THICKNESS}
              max={ANNOTATION_MAX_THICKNESS}
              style={{ width: 120 }}
            />
            <Group>
              <ActionIcon variant="outline" color="gray" onClick={handleClear}>
                <TrashIcon size={16} />
              </ActionIcon>
              <Button size="xs" variant="subtle" color="gray" onClick={() => { setDrawMode(false); setDrawMenuOpen(false) }}>
                Exit draw
              </Button>
            </Group>
          </Stack>
        </Paper>
      )}

      <DraggableFab
        icon={drawMode ? <EyeIcon size={20} /> : <PencilIcon size={20} />}
        onClick={() => drawMode ? setDrawMenuOpen(v => !v) : setDrawMode(true)}
        label={drawMode ? 'Open draw menu' : 'Enter draw mode'}
        initialTop={24}
        initialRight={24}
      />

      <DraggableFab
        icon={<CodeSimpleIcon size={20} />}
        onClick={() => setShowRoomCode(v => !v)}
        label="Toggle room code"
        initialTop={80}
        initialRight={24}
      />

      {showRoomCode && roomCode && (
        <Box pos="fixed" top={140} right={24} style={{ zIndex: 1000 }}>
          <Code c="gray.3" fz={18} fw={700} bg="dark.7" lts="0.1em">
            {roomCode}
          </Code>
        </Box>
      )}

      <Modal
        opened={goToModalOpen}
        onClose={() => setGoToModalOpen(false)}
        title="Go to slide"
        centered
      >
        <NumberInput
          placeholder={`Slide (1-${slideCount})`}
          value={goToValue}
          min={1}
          max={slideCount}
          onChange={(value) => setGoToValue(value.toString())}
        />
        <Button fullWidth mt="md" onClick={handleGoToSubmit} disabled={!goToValue}>
          Go
        </Button>
      </Modal>
    </Container>
  )
}
