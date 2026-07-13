import { useEffect, useCallback, useState } from 'react'
import { useParams, Link } from 'react-router'
import {
  Box, Button, Center, Code, Group, Loader, NumberInput, Slider, Stack, Text, Title, ActionIcon, Tooltip,
} from '@mantine/core'
import { CaretLeftIcon } from '@phosphor-icons/react/dist/csr/CaretLeft'
import { CaretRightIcon } from '@phosphor-icons/react/dist/csr/CaretRight'
import { PencilIcon } from '@phosphor-icons/react/dist/csr/Pencil'
import { EyeIcon } from '@phosphor-icons/react/dist/csr/Eye'
import { TrashIcon } from '@phosphor-icons/react/dist/csr/Trash'
import { useSlideShow } from '../hooks/use-slide-show'
import { usePresentation } from '../hooks/use-presentation'
import { CLIENT_CONFIGURE } from '../../shared/cfg/routes'
import { WS_CMD_INIT_PRESENTATION, WS_CMD_NEXT_SLIDE, WS_CMD_PREV_SLIDE, WS_CMD_GO_TO_SLIDE, WS_CMD_ANNOTATION, ANNOTATION_COLORS, ANNOTATION_DEFAULT_COLOR, ANNOTATION_DEFAULT_THICKNESS, ANNOTATION_MIN_THICKNESS, ANNOTATION_MAX_THICKNESS } from '../cfg'
import { visibleStrokes, drawStrokes, toPercent } from '../utils/annotation-canvas'
import type { AnnotationPoint, AnnotationOperation } from '../types'
import { uuidv7 } from '../../shared/util/uuid'

export function ControlView() {
  const { id } = useParams<{ id: string }>()
  const { presentation, isLoading: presLoading } = usePresentation(id ?? '')
  const { send, cachedHtml, slideCount, currentSlide, loading, fetchError, roomCode, iframeRef, canvasRef, operationsBySlide, setOperationsBySlide } =
    useSlideShow(id ? { command: WS_CMD_INIT_PRESENTATION, parameters: { presentation_id: id } } : null)

  const [goToValue, setGoToValue] = useState('')
  const [drawMode, setDrawMode] = useState(false)
  const [annotationColor, setAnnotationColor] = useState(ANNOTATION_DEFAULT_COLOR)
  const [annotationThickness, setAnnotationThickness] = useState(ANNOTATION_DEFAULT_THICKNESS)
  const [isDrawing, setIsDrawing] = useState(false)
  const [currentPoints, setCurrentPoints] = useState<AnnotationPoint[]>([])

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

  function handlePrev() {
    send({ command: WS_CMD_PREV_SLIDE, parameters: {} })
  }

  function handleNext() {
    send({ command: WS_CMD_NEXT_SLIDE, parameters: {} })
  }

  function handleGoTo() {
    const num = parseInt(goToValue, 10)
    if (isNaN(num) || num < 1 || num > slideCount) return
    send({ command: WS_CMD_GO_TO_SLIDE, parameters: { slide_number: num - 1 } })
    setGoToValue('')
  }

  const handlePointerDown = useCallback((e: React.PointerEvent<HTMLCanvasElement>) => {
    const canvas = canvasRef.current
    if (!canvas) return
    setIsDrawing(true)
    const rect = canvas.getBoundingClientRect()
    setCurrentPoints([toPercent(e.clientX, e.clientY, rect)])
    canvas.setPointerCapture(e.pointerId)
  }, [])

  const handlePointerMove = useCallback((e: React.PointerEvent<HTMLCanvasElement>) => {
    if (!isDrawing) return
    const canvas = canvasRef.current
    if (!canvas) return
    const rect = canvas.getBoundingClientRect()
    setCurrentPoints(prev => [...prev, toPercent(e.clientX, e.clientY, rect)])
  }, [isDrawing])

  const handlePointerUp = useCallback(() => {
    if (!isDrawing) return
    setIsDrawing(false)

    const points = currentPoints
    setCurrentPoints([])

    if (points.length < 2) return

    const id = uuidv7()
    const op: AnnotationOperation = {
      type: 'stroke',
      id,
      payload: { points, color: annotationColor, thickness: annotationThickness },
    }

    setOperationsBySlide(prev => {
      const slideKey = String(currentSlide)
      const existing = prev[slideKey] ?? []
      return { ...prev, [slideKey]: [...existing, op] }
    })

    send({
      command: WS_CMD_ANNOTATION,
      parameters: {
        type: 'stroke',
        id,
        payload: { points, color: annotationColor, thickness: annotationThickness },
      },
    })
  }, [isDrawing, currentPoints, annotationColor, annotationThickness, currentSlide, send])

  const handleClear = useCallback(() => {
    const id = uuidv7()
    send({
      command: WS_CMD_ANNOTATION,
      parameters: { type: 'clear', id },
    })
  }, [send])

  if (presLoading || loading) {
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
    <Box p="md" bg="dark.9" style={{ height: '100vh', display: 'flex', flexDirection: 'column' }}>
      <Group mb="md">
        <Button
          component={Link}
          to={CLIENT_CONFIGURE}
          variant="outline"
          color="gray"
          size="sm"
          style={{ borderColor: 'var(--mantine-color-dark-6)' }}
        >
          &larr; Back
        </Button>
        <Title order={3} c="white">{presentation?.title ?? 'Control'}</Title>
        {roomCode && (
          <Code c="gray.3" fz={18} fw={700} bg="dark.7" style={{ letterSpacing: '0.1em' }}>
            {roomCode}
          </Code>
        )}
        <Tooltip label={drawMode ? 'View mode' : 'Draw mode'}>
          <ActionIcon
            variant={drawMode ? 'filled' : 'outline'}
            color={drawMode ? 'red' : 'gray'}
            onClick={() => setDrawMode(v => !v)}
            size="lg"
          >
            {drawMode ? <EyeIcon size={20} /> : <PencilIcon size={20} />}
          </ActionIcon>
        </Tooltip>
      </Group>

      {drawMode && (
        <Group mb="md" gap="xs" wrap="nowrap">
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
            style={{ width: 80 }}
          />
          <ActionIcon variant="outline" color="gray" onClick={handleClear}>
            <TrashIcon size={16} />
          </ActionIcon>
        </Group>
      )}

      <Box pos="relative" flex={1}>
        <Box
          inset={0}
          m="auto"
          pos="absolute"
          maw="100%"
          mah="100%"
          style={{ aspectRatio: '48/35' }}
        >
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
        </Box>
      </Box>
      <Group justify="center" mb="md">
        <Button
          size="xl"
          variant="outline"
          color="gray"
          onClick={handlePrev}
          disabled={currentSlide <= 0}
          px="lg"
          style={{ borderColor: 'var(--mantine-color-dark-6)' }}
        >
          <CaretLeftIcon size={24} />
        </Button>
        <Text size="xl" c="gray.3" style={{ minWidth: 100, textAlign: 'center' }}>
          {currentSlide + 1} / {slideCount}
        </Text>
        <Button
          size="xl"
          variant="outline"
          color="gray"
          onClick={handleNext}
          disabled={currentSlide >= slideCount - 1}
          px="lg"
          style={{ borderColor: 'var(--mantine-color-dark-6)' }}
        >
          <CaretRightIcon size={24} />
        </Button>
      </Group>

      <form onSubmit={(e) => { e.preventDefault(); handleGoTo() }}>
        <Group justify="center">
          <NumberInput
            placeholder={`Go to slide (1-${slideCount})`}
            value={goToValue}
            min={1}
            max={slideCount}
            onChange={(value) => setGoToValue(value.toString())}
            miw={160}
          />
          <Button type="submit" disabled={!goToValue}>Go</Button>
        </Group>
      </form>
    </Box>
  )
}
