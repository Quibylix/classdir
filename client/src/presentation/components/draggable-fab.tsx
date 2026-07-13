import { useRef, useState, type ReactNode } from 'react'
import { Box, ActionIcon } from '@mantine/core'

type DraggableFabProps = {
  icon: ReactNode
  onClick: () => void
  label: string
  initialTop: number
  initialRight: number
}

const BUTTON_SIZE = 48;

export function DraggableFab({ icon, onClick, label, initialTop, initialRight }: DraggableFabProps) {
  const [pos, setPos] = useState({ top: initialTop, right: initialRight })
  const dragRef = useRef<{ startX: number; startY: number; startTop: number; startRight: number } | null>(null)

  const handlePointerDown = (e: React.PointerEvent) => {
    e.currentTarget.setPointerCapture(e.pointerId)
    dragRef.current = {
      startX: e.clientX,
      startY: e.clientY,
      startTop: pos.top,
      startRight: pos.right,
    }
  }

  const handlePointerMove = (e: React.PointerEvent) => {
    const d = dragRef.current
    if (!d) return
    setPos({
      top: Math.min(Math.max(0, d.startTop + (e.clientY - d.startY)), window.innerHeight - BUTTON_SIZE),
      right: Math.min(Math.max(0, d.startRight - (e.clientX - d.startX)), window.innerWidth - BUTTON_SIZE),
    })
  }

  const handlePointerUp = (e: React.PointerEvent) => {
    const d = dragRef.current
    if (!d) return
    const moved = Math.abs(e.clientX - d.startX) + Math.abs(e.clientY - d.startY)
    dragRef.current = null
    if (moved < 5) {
      onClick()
    }
  }

  return (
    <Box
      pos="fixed"
      top={pos.top}
      right={pos.right}
      style={{ zIndex: 1000, cursor: 'grab', touchAction: 'none' }}
      onPointerDown={handlePointerDown}
      onPointerMove={handlePointerMove}
      onPointerUp={handlePointerUp}
      onPointerCancel={() => { dragRef.current = null }}
    >
      <ActionIcon
        title={label}
        w={BUTTON_SIZE}
        h={BUTTON_SIZE}
        size="lg"
        variant="filled"
        color="dark"
        aria-label={label}
      >
        {icon}
      </ActionIcon>
    </Box>
  )
}
