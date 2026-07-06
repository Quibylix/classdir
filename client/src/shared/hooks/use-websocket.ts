import { useEffect, useRef, useState, useCallback } from 'react'
import type z from 'zod'
import { WS_STATUS } from '../types'
import { WS_RECONNECT_TIMEOUT_MS } from '../cfg/ws'
import type { WSStatus } from '../types'

type UseWebSocketOptions<T> = {
  url: string
  schema: z.ZodSchema<T>
  onMessage?: (msg: T) => void
}

export function useSafeWebSocket<T>({ url, onMessage, schema }: UseWebSocketOptions<T>) {
  const [status, setStatus] = useState<WSStatus>(WS_STATUS.Disconnected)
  const wsRef = useRef<WebSocket | null>(null)
  const reconnectTimerRef = useRef<ReturnType<typeof setTimeout> | undefined>(undefined)
  const onMessageRef = useRef(onMessage)
  onMessageRef.current = onMessage

  const connect = useCallback(() => {
    if (wsRef.current?.readyState === WebSocket.OPEN) return

    setStatus(WS_STATUS.Connecting)
    const ws = new WebSocket(url)
    wsRef.current = ws

    ws.onopen = () => setStatus(WS_STATUS.Connected)
    ws.onclose = () => {
      setStatus(WS_STATUS.Disconnected)
      reconnectTimerRef.current = setTimeout(connect, WS_RECONNECT_TIMEOUT_MS)
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
