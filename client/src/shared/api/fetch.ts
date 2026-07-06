import { ResultAsync } from 'neverthrow'
import { z } from 'zod'
import { api } from './client'
import { HTTP_NO_CONTENT, ERR_CODE_UNKNOWN } from '../cfg/http'
import { toApiError, isApiError } from './errors'
import type { ApiError } from './errors'

export type FetchError = ApiError | z.ZodError<unknown>

export function safeFetch<T>(path: string, schema: z.ZodType<T>, options?: RequestInit): ResultAsync<T, FetchError> {
  return ResultAsync.fromPromise(
    api(path, options).then(async (res) => {
      if (!res.ok) throw await toApiError(res)
      if (res.status === HTTP_NO_CONTENT) return schema.parse(undefined)
      const body = await res.json()
      return schema.parse(body.data)
    }),
    (e) => {
      if (isApiError(e)) return e
      if (e instanceof z.ZodError) return e
      return { code: ERR_CODE_UNKNOWN, message: String(e), status: 0 }
    },
  )
}
