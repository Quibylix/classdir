import { z } from 'zod'

export const SlideMetadataSchema = z.object({
  title: z.string(),
  author: z.string(),
})
export type SlideMetadata = z.infer<typeof SlideMetadataSchema>

export const SlideSchema = z.object({
  id: z.string(),
  slide_number: z.number(),
  content: z.string(),
  metadata: SlideMetadataSchema,
})
export type Slide = z.infer<typeof SlideSchema>

export const PresentationSchema = z.object({
  id: z.string(),
  title: z.string(),
  slides: z.array(SlideSchema),
})
export type Presentation = z.infer<typeof PresentationSchema>

export const PresentationPreviewSchema = z.object({
  id: z.string(),
  title: z.string(),
})
export type PresentationPreview = z.infer<typeof PresentationPreviewSchema>
