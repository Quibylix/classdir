export const WS_CMD_INIT_PRESENTATION = 'init_presentation'
export const WS_CMD_JOIN_ROOM = 'join_room'
export const WS_CMD_NEXT_SLIDE = 'next_slide'
export const WS_CMD_PREV_SLIDE = 'prev_slide'
export const WS_CMD_GO_TO_SLIDE = 'go_to_slide'

export const WS_EVENT_SLIDE_CHANGED = 'slide_changed'

export const POST_MSG_TYPE = {
  Navigate: 'navigate',
} as const

export const CDN_REVEAL_VERSION = 'reveal.js@6'
export const CDN_REVEAL_CSS = `https://cdn.jsdelivr.net/npm/${CDN_REVEAL_VERSION}/dist/reveal.css`
export const CDN_REVEAL_THEME_CSS = `https://cdn.jsdelivr.net/npm/${CDN_REVEAL_VERSION}/dist/theme/black.css`
export const CDN_REVEAL_JS = `https://cdn.jsdelivr.net/npm/${CDN_REVEAL_VERSION}/dist/reveal.js`

export const DEFAULT_SLIDE_CONTENT = '<h1>New Slide</h1>'
