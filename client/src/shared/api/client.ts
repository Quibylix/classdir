import { CREDENTIALS_INCLUDE, HEADER_CONTENT_TYPE, CONTENT_TYPE_JSON } from '../cfg/http'

const BASE = import.meta.env.VITE_API_URL

export async function api(path: string, options?: RequestInit): Promise<Response> {
  return fetch(`${BASE}${path}`, {
    credentials: CREDENTIALS_INCLUDE,
    headers: { [HEADER_CONTENT_TYPE]: CONTENT_TYPE_JSON },
    ...options,
  })
}
