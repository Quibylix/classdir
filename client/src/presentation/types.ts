import { z } from 'zod'

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

export const WSOutputMessageSchema = z.object({
  event: z.literal('slide_changed'),
  data: z.object({
    current_slide: z.number(),
  }),
}).or(z.object({
  data: z.object({
    slides: z.array(SlideSchema),
    current_index: z.number(),
  })
})).or(z.object({
  error: z.object({
    code: z.string(),
    message: z.string()
  }),
}))

export type WSOutputMessage = z.infer<typeof WSOutputMessageSchema>
