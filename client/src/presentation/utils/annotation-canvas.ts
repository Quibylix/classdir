import { WS_ANNOTATION_TYPE_CLEAR, WS_ANNOTATION_TYPE_STROKE } from '../cfg'
import type { AnnotationOperation, AnnotationPayload, AnnotationPoint } from '../types'

export function visibleStrokes(operations: AnnotationOperation[]): AnnotationPayload[] {
  const strokes: AnnotationPayload[] = []
  for (const op of operations) {
    if (op.type === WS_ANNOTATION_TYPE_CLEAR) {
      strokes.length = 0
    } else if (op.type === WS_ANNOTATION_TYPE_STROKE) {
      strokes.push(op.payload)
    }
  }
  return strokes
}

export function drawStrokes(
  ctx: CanvasRenderingContext2D,
  strokes: AnnotationPayload[],
  currentPoints: AnnotationPoint[] | undefined,
  width: number,
  height: number,
): void {
  ctx.clearRect(0, 0, width, height)

  for (const stroke of strokes) {
    drawSingleStroke(ctx, stroke.points, stroke.color, stroke.thickness, width, height)
  }

  if (currentPoints && currentPoints.length > 1) {
    drawSingleStroke(ctx, currentPoints, '#ffffff', 2, width, height)
  }
}

function drawSingleStroke(
  ctx: CanvasRenderingContext2D,
  points: AnnotationPoint[],
  color: string,
  thickness: number,
  width: number,
  height: number,
): void {
  if (points.length < 2) return

  ctx.beginPath()
  ctx.strokeStyle = color
  ctx.lineWidth = thickness
  ctx.lineCap = 'round'
  ctx.lineJoin = 'round'

  const first = points[0]
  ctx.moveTo((first.x / 100) * width, (first.y / 100) * height)

  for (let i = 1; i < points.length; i++) {
    const p = points[i]
    ctx.lineTo((p.x / 100) * width, (p.y / 100) * height)
  }

  ctx.stroke()
}

export function toPercent(clientX: number, clientY: number, rect: DOMRect): AnnotationPoint {
  return {
    x: ((clientX - rect.left) / rect.width) * 100,
    y: ((clientY - rect.top) / rect.height) * 100,
  }
}
