import { z } from 'zod'

const ApiErrorSchema = z.object({
  error: z.object({
    code: z.string(),
    message: z.string(),
  }),
})

const ApiErrorObjSchema = z.object({
  code: z.string(),
  message: z.string(),
  status: z.number(),
})
export type ApiError = z.infer<typeof ApiErrorObjSchema>

export function isApiError(e: unknown): e is ApiError {
  return ApiErrorObjSchema.safeParse(e).success
}

export async function toApiError(res: Response): Promise<ApiError> {
  const body = await res.json().catch(() => ({}))
  const parsed = ApiErrorSchema.safeParse(body)
  if (parsed.success) {
    return { ...parsed.data.error, status: res.status }
  }
  return { code: 'UNKNOWN', message: res.statusText, status: res.status }
}
