import { z } from 'zod'
import { WS_EVENT_SLIDE_CHANGED, WS_EVENT_ANNOTATION_ADDED, WS_EVENT_ANNOTATIONS_BATCH, WS_ANNOTATION_TYPE_CLEAR, WS_ANNOTATION_TYPE_STROKE } from './cfg'

export const SlideSchema = z.object({
  id: z.string(),
  content: z.string(),
})
export type Slide = z.infer<typeof SlideSchema>

export const PresentationSchema = z.object({
  id: z.string(),
  title: z.string(),
  slide_order: z.array(z.string()),
  slides: z.array(SlideSchema),
})
export type Presentation = z.infer<typeof PresentationSchema>

export const PresentationPreviewSchema = z.object({
  id: z.string(),
  title: z.string(),
})
export type PresentationPreview = z.infer<typeof PresentationPreviewSchema>

export const AnnotationPointSchema = z.object({
  x: z.number(),
  y: z.number(),
})
export type AnnotationPoint = z.infer<typeof AnnotationPointSchema>

export const AnnotationPayloadSchema = z.object({
  points: z.array(AnnotationPointSchema),
  color: z.string(),
  thickness: z.number(),
})
export type AnnotationPayload = z.infer<typeof AnnotationPayloadSchema>

export const AnnotationOperationSchema = z.object({
  type: z.literal(WS_ANNOTATION_TYPE_CLEAR),
  id: z.string(),
}).or(z.object({
  type: z.literal(WS_ANNOTATION_TYPE_STROKE),
  id: z.string(),
  payload: AnnotationPayloadSchema,
}))
export type AnnotationOperation = z.infer<typeof AnnotationOperationSchema>

export const WSOutputMessageSchema = z.object({
  event: z.literal(WS_EVENT_SLIDE_CHANGED),
  data: z.object({
    current_slide: z.number(),
  }),
}).or(z.object({
  event: z.literal(WS_EVENT_ANNOTATION_ADDED),
  data: AnnotationOperationSchema,
})).or(z.object({
  event: z.literal(WS_EVENT_ANNOTATIONS_BATCH),
  data: z.object({
    operations_by_slide: z.record(z.string(), z.array(AnnotationOperationSchema)),
  }),
})).or(z.object({
  data: z.object({
    slides: z.array(SlideSchema),
    current_index: z.number(),
    room_code: z.string().optional(),
  })
})).or(z.object({
  error: z.object({
    code: z.string(),
    message: z.string()
  }),
}))

export type WSOutputMessage = z.infer<typeof WSOutputMessageSchema>
