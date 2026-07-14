import { useState, useCallback } from 'react'
import { toPercent } from '../utils/annotation-canvas'
import type { AnnotationPoint } from '../types'
import { WS_CMD_ANNOTATION, ANNOTATION_DEFAULT_COLOR, ANNOTATION_DEFAULT_THICKNESS, WS_ANNOTATION_TYPE_CLEAR, WS_ANNOTATION_TYPE_STROKE } from '../cfg'
import { uuidv7 } from '../../shared/util/uuid'

interface UseAnnotationOptions {
  send: (data: unknown) => void
  canvasRef: React.RefObject<HTMLCanvasElement | null>
}

export function useAnnotation({ send, canvasRef }: UseAnnotationOptions) {
  const [drawMode, setDrawMode] = useState(false)
  const [annotationColor, setAnnotationColor] = useState(ANNOTATION_DEFAULT_COLOR)
  const [annotationThickness, setAnnotationThickness] = useState(ANNOTATION_DEFAULT_THICKNESS)
  const [isDrawing, setIsDrawing] = useState(false)
  const [currentPoints, setCurrentPoints] = useState<AnnotationPoint[]>([])

  const handlePointerDown = useCallback((e: React.PointerEvent<HTMLCanvasElement>) => {
    const canvas = canvasRef.current
    if (!canvas) return
    setIsDrawing(true)
    const rect = canvas.getBoundingClientRect()
    setCurrentPoints([toPercent(e.clientX, e.clientY, rect)])
    canvas.setPointerCapture(e.pointerId)
  }, [canvasRef])

  const handlePointerMove = useCallback((e: React.PointerEvent<HTMLCanvasElement>) => {
    if (!isDrawing) return
    const canvas = canvasRef.current
    if (!canvas) return
    const rect = canvas.getBoundingClientRect()
    setCurrentPoints(prev => [...prev, toPercent(e.clientX, e.clientY, rect)])
  }, [isDrawing, canvasRef])

  const handlePointerUp = useCallback(() => {
    if (!isDrawing) return
    setIsDrawing(false)

    const points = currentPoints
    setCurrentPoints([])

    if (points.length < 2) return

    const id = uuidv7()
    send({
      command: WS_CMD_ANNOTATION,
      parameters: {
        type: WS_ANNOTATION_TYPE_STROKE,
        id,
        payload: { points, color: annotationColor, thickness: annotationThickness },
      },
    })
  }, [isDrawing, currentPoints, annotationColor, annotationThickness, send])

  const handleClear = useCallback(() => {
    const id = uuidv7()
    send({
      command: WS_CMD_ANNOTATION,
      parameters: { type: WS_ANNOTATION_TYPE_CLEAR, id },
    })
  }, [send])

  return {
    drawMode,
    setDrawMode,
    annotationColor,
    setAnnotationColor,
    annotationThickness,
    setAnnotationThickness,
    currentPoints,
    handlePointerDown,
    handlePointerMove,
    handlePointerUp,
    handleClear,
  }
}
