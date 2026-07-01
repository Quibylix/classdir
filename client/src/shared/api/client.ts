const BASE = import.meta.env.VITE_API_URL

export async function api(path: string, options?: RequestInit): Promise<Response> {
  return fetch(`${BASE}${path}`, {
    credentials: 'include',
    headers: { 'Content-Type': 'application/json' },
    ...options,
  })
}
