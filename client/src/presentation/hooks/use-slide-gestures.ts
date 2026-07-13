import { useState } from 'react'
import { WS_CMD_PREV_SLIDE, WS_CMD_NEXT_SLIDE, WS_CMD_GO_TO_SLIDE } from '../cfg'

interface UseSlideGesturesOptions {
  send: (data: unknown) => void
  slideCount: number
  currentSlide: number
}

export function useSlideGestures({ send, slideCount, currentSlide }: UseSlideGesturesOptions) {
  const [goToValue, setGoToValue] = useState('')
  const [goToModalOpen, setGoToModalOpen] = useState(false)

  const handleZoneClick = (zone: 'prev' | 'next' | 'goto') => {
    if (zone === 'prev' && currentSlide > 0) {
      send({ command: WS_CMD_PREV_SLIDE, parameters: {} })
      return
    }
    if (zone === 'next' && currentSlide < slideCount - 1) {
      send({ command: WS_CMD_NEXT_SLIDE, parameters: {} })
      return
    }
    if (zone === 'goto') {
      setGoToValue('')
      setGoToModalOpen(true)
    }
  }

  const handleGoToSubmit = () => {
    const num = parseInt(goToValue, 10)
    if (isNaN(num) || num < 1 || num > slideCount) return
    send({ command: WS_CMD_GO_TO_SLIDE, parameters: { slide_number: num - 1 } })
    setGoToModalOpen(false)
    setGoToValue('')
  }

  return {
    goToValue,
    setGoToValue,
    goToModalOpen,
    setGoToModalOpen,
    handleZoneClick,
    handleGoToSubmit,
  }
}
