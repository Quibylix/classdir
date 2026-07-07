import { useEffect, useRef, useState } from 'react'
import { useParams, Link } from 'react-router'
import {
  Box, Button, Center, Code, Group, Loader, NumberInput, Stack, Text, Title,
} from '@mantine/core'
import { CaretLeftIcon } from '@phosphor-icons/react/dist/csr/CaretLeft'
import { CaretRightIcon } from '@phosphor-icons/react/dist/csr/CaretRight'
import { useSafeWebSocket } from '../../shared/hooks/use-websocket'
import { WSOutputMessageSchema, type Slide } from '../types'
import { usePresentation } from '../hooks/use-presentation'
import { WS_V1, CLIENT_CONFIGURE } from '../../shared/cfg/routes'
import { WS_STATUS } from '../../shared/types'
import { POST_MSG_TYPE, WS_CMD_INIT_PRESENTATION, WS_CMD_NEXT_SLIDE, WS_CMD_PREV_SLIDE, WS_CMD_GO_TO_SLIDE, CDN_REVEAL_CSS, CDN_REVEAL_THEME_CSS, CDN_REVEAL_JS } from '../cfg'

function buildPresentHtml(slides: Slide[], initialSlide: number): string {
  const slidesHtml = slides.map(s => `<section>${s.content}</section>`).join('\n')
  return `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <link rel="stylesheet" href="${CDN_REVEAL_CSS}">
  <link rel="stylesheet" href="${CDN_REVEAL_THEME_CSS}">
  <style>
    * { margin: 0; padding: 0; box-sizing: border-box; }
    html, body { width: 100%; height: 100%; overflow: hidden; }
  </style>
</head>
<body>
  <div class="reveal" id="reveal">
    <div class="slides">${slidesHtml}</div>
  </div>
  <script src="${CDN_REVEAL_JS}"></script>
  <script>
    Reveal.initialize({ transition: 'slide', progress: false, controls: false }).then(function() {
      Reveal.slide(${initialSlide});
    });
    window.addEventListener('message', function(e) {
      if (e.data.type === '${POST_MSG_TYPE.Navigate}') Reveal.slide(e.data.index);
    });
  </script>
</body>
</html>`
}

export function ControlView() {
  const { id } = useParams<{ id: string }>()
  const { presentation, isLoading: presLoading } = usePresentation(id ?? '')
  const [cachedHtml, setCachedHtml] = useState<string>()
  const [slideCount, setSlideCount] = useState<number>(0)
  const [currentSlide, setCurrentSlide] = useState(0)
  const [loading, setLoading] = useState(true)
  const [fetchError, setFetchError] = useState<string | null>(null)
  const [goToValue, setGoToValue] = useState('')
  const [roomCode, setRoomCode] = useState<string>()
  const joinedRef = useRef(false)
  const iframeRef = useRef<HTMLIFrameElement>(null)

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

      if ('event' in msg) {
        setCurrentSlide(msg.data.current_slide)
        iframeRef.current?.contentWindow?.postMessage({ type: POST_MSG_TYPE.Navigate, index: msg.data.current_slide }, window.location.origin)
        return
      }

      setSlideCount(msg.data.slides.length)
      setCurrentSlide(msg.data.current_index)
      if (msg.data.room_code) {
        setRoomCode(msg.data.room_code)
      }
      setCachedHtml(buildPresentHtml(msg.data.slides, msg.data.current_index))
      setLoading(false)
      return;
    },
  })

  useEffect(() => {
    if (status === WS_STATUS.Connected && id && !joinedRef.current) {
      send({ command: WS_CMD_INIT_PRESENTATION, parameters: { presentation_id: id } })
      joinedRef.current = true
    }
    if (status === WS_STATUS.Disconnected) {
      joinedRef.current = false
    }
  }, [status, id, send])

  function handlePrev() {
    send({ command: WS_CMD_PREV_SLIDE, parameters: {} })
  }

  function handleNext() {
    send({ command: WS_CMD_NEXT_SLIDE, parameters: {} })
  }

  function handleGoTo() {
    const num = parseInt(goToValue, 10)
    if (isNaN(num) || num < 1 || num > slideCount) return
    send({ command: WS_CMD_GO_TO_SLIDE, parameters: { slide_number: num - 1 } })
    setGoToValue('')
  }

  if (presLoading || loading) {
    return <Center h="100vh" bg="dark.9"><Loader /></Center>
  }

  if (fetchError) {
    return (
      <Center h="100vh" bg="dark.9">
        <Stack align="center">
          <Text c="red">{fetchError}</Text>
          <Button component={Link} to={CLIENT_CONFIGURE}>Back</Button>
        </Stack>
      </Center>
    )
  }

  if (slideCount === 0) {
    return (
      <Center h="100vh" bg="dark.9">
        <Stack align="center">
          <Text c="dimmed">No slides in this presentation</Text>
          <Button component={Link} to={CLIENT_CONFIGURE}>Back</Button>
        </Stack>
      </Center>
    )
  }

  return (
    <Box p="md" bg="dark.9" style={{ height: '100vh', display: 'flex', flexDirection: 'column' }}>
      <Group mb="md">
        <Button
          component={Link}
          to={CLIENT_CONFIGURE}
          variant="outline"
          color="gray"
          size="sm"
          style={{ borderColor: 'var(--mantine-color-dark-6)' }}
        >
          &larr; Back
        </Button>
        <Title order={3} c="white">{presentation?.title ?? 'Control'}</Title>
        {roomCode && (
          <Code c="gray.3" fz={18} fw={700} bg="dark.7" style={{ letterSpacing: '0.1em' }}>
            {roomCode}
          </Code>
        )}
      </Group>

      <iframe
        ref={iframeRef}
        srcDoc={cachedHtml}
        title="Presentation"
        style={{ width: '100%', height: 0, border: 'none', display: 'block', flex: 1 }}
      />

      <Group justify="center" mb="md">
        <Button
          size="xl"
          variant="outline"
          color="gray"
          onClick={handlePrev}
          disabled={currentSlide <= 0}
          px="lg"
          style={{ borderColor: 'var(--mantine-color-dark-6)' }}
        >
          <CaretLeftIcon size={24} />
        </Button>
        <Text size="xl" c="gray.3" style={{ minWidth: 100, textAlign: 'center' }}>
          {currentSlide + 1} / {slideCount}
        </Text>
        <Button
          size="xl"
          variant="outline"
          color="gray"
          onClick={handleNext}
          disabled={currentSlide >= slideCount - 1}
          px="lg"
          style={{ borderColor: 'var(--mantine-color-dark-6)' }}
        >
          <CaretRightIcon size={24} />
        </Button>
      </Group>

      <form onSubmit={(e) => { e.preventDefault(); handleGoTo() }}>
        <Group justify="center">
          <NumberInput
            placeholder={`Go to slide (1-${slideCount})`}
            value={goToValue}
            min={1}
            max={slideCount}
            onChange={(value) => setGoToValue(value.toString())}
            miw={160}
          />
          <Button type="submit" disabled={!goToValue}>Go</Button>
        </Group>
      </form>
    </Box>
  )
}
