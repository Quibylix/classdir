import { useEffect, useRef, useState, useCallback } from 'react'
import type z from 'zod'

export type WSStatus = 'connecting' | 'connected' | 'disconnected'

type UseWebSocketOptions<T> = {
  url: string
  schema: z.ZodSchema<T>
  onMessage?: (msg: T) => void
}

export function useSafeWebSocket<T>({ url, onMessage, schema }: UseWebSocketOptions<T>) {
  const [status, setStatus] = useState<WSStatus>('disconnected')
  const wsRef = useRef<WebSocket | null>(null)
  const reconnectTimerRef = useRef<ReturnType<typeof setTimeout> | undefined>(undefined)
  const onMessageRef = useRef(onMessage)
  onMessageRef.current = onMessage

  const connect = useCallback(() => {
    if (wsRef.current?.readyState === WebSocket.OPEN) return

    setStatus('connecting')
    const ws = new WebSocket(url)
    wsRef.current = ws

    ws.onopen = () => setStatus('connected')
    ws.onclose = () => {
      setStatus('disconnected')
      reconnectTimerRef.current = setTimeout(connect, 3000)
    }
    ws.onmessage = (e) => {
      try {
        const parsed = JSON.parse(e.data)
        const result = schema.parse(parsed)
        onMessageRef.current?.(result)
      } catch {
        // ignore malformed
      }
    }
    ws.onerror = () => { }
  }, [url])

  useEffect(() => {
    connect()
    return () => {
      clearTimeout(reconnectTimerRef.current)
      wsRef.current?.close()
      wsRef.current = null
    }
  }, [connect])

  const send = useCallback((data: unknown) => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify(data))
    }
  }, [])

  return { status, send }
}
