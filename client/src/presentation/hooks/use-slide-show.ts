import { useEffect, useRef, useState } from 'react'
import { useSafeWebSocket } from '../../shared/hooks/use-websocket'
import { WSOutputMessageSchema, type AnnotationOperation } from '../types'
import { WS_V1 } from '../../shared/cfg/routes'
import { WS_STATUS } from '../../shared/types'
import { WS_EVENT_ANNOTATION_ADDED, WS_EVENT_SLIDE_CHANGED, POST_MSG_TYPE } from '../cfg'
import { buildPresentHtml } from '../utils/reveal-html'

export function useSlideShow(initCommand: { command: string; parameters: Record<string, string> } | null) {
  const [cachedHtml, setCachedHtml] = useState<string>()
  const [slideCount, setSlideCount] = useState(0)
  const [currentSlide, setCurrentSlide] = useState(0)
  const [loading, setLoading] = useState(true)
  const [fetchError, setFetchError] = useState<string | null>(null)
  const [roomCode, setRoomCode] = useState<string | undefined>()
  const [operationsBySlide, setOperationsBySlide] = useState<Record<string, AnnotationOperation[]>>({})
  const joinedRef = useRef(false)
  const iframeRef = useRef<HTMLIFrameElement>(null)
  const canvasRef = useRef<HTMLCanvasElement>(null)

  const wsUrl = `${window.location.protocol === 'https:' ? 'wss:' : 'ws:'}//${window.location.host}${WS_V1}`

  const { status, send } = useSafeWebSocket({
    url: wsUrl,
    schema: WSOutputMessageSchema,
    onMessage(msg) {
      if ('error' in msg) {
        setFetchError(msg.error.message)
        setLoading(false)
        return
      }

      if (!('event' in msg)) {
        setCachedHtml(buildPresentHtml(msg.data.slides, msg.data.current_index))
        setSlideCount(msg.data.slides.length)
        setCurrentSlide(msg.data.current_index)
        if (msg.data.room_code) {
          setRoomCode(msg.data.room_code)
        }
        setLoading(false)
        return
      }

      if (msg.event === WS_EVENT_SLIDE_CHANGED) {
        setCurrentSlide(msg.data.current_slide)
        iframeRef.current?.contentWindow?.postMessage(
          { type: POST_MSG_TYPE.Navigate, index: msg.data.current_slide },
          window.location.origin,
        )
        return
      }

      if (msg.event === WS_EVENT_ANNOTATION_ADDED) {
        setOperationsBySlide(prev => {
          const slideKey = String(currentSlide)
          const existing = prev[slideKey] ?? []
          if (existing.some(op => op.id === msg.data.id)) return prev
          return { ...prev, [slideKey]: [...existing, msg.data] }
        })
        return
      }

      setOperationsBySlide(msg.data.operations_by_slide)
    },
  })

  useEffect(() => {
    if (status === WS_STATUS.Connected && initCommand && !joinedRef.current) {
      send(initCommand)
      joinedRef.current = true
    }
    if (status === WS_STATUS.Disconnected) {
      joinedRef.current = false
    }
  }, [status, initCommand, send])

  return {
    send,
    cachedHtml,
    slideCount,
    currentSlide,
    loading,
    fetchError,
    roomCode,
    operationsBySlide,
    setOperationsBySlide,
    iframeRef,
    canvasRef,
  }
}
