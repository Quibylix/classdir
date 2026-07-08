export const WS_CMD_INIT_PRESENTATION = 'init_presentation'
export const WS_CMD_JOIN_ROOM = 'join_room'
export const WS_CMD_NEXT_SLIDE = 'next_slide'
export const WS_CMD_PREV_SLIDE = 'prev_slide'
export const WS_CMD_GO_TO_SLIDE = 'go_to_slide'
export const WS_CMD_ANNOTATION = 'annotation'

export const WS_EVENT_SLIDE_CHANGED = 'slide_changed'
export const WS_EVENT_ANNOTATION_ADDED = 'annotation_added'
export const WS_EVENT_ANNOTATIONS_BATCH = 'annotations_batch'

export const WS_ANNOTATION_TYPE_STROKE = 'stroke'
export const WS_ANNOTATION_TYPE_CLEAR = 'clear'

export const POST_MSG_TYPE = {
  Navigate: 'navigate',
} as const

export const CDN_REVEAL_VERSION = 'reveal.js@6'
export const CDN_REVEAL_CSS = `https://cdn.jsdelivr.net/npm/${CDN_REVEAL_VERSION}/dist/reveal.css`
export const CDN_REVEAL_THEME_CSS = `https://cdn.jsdelivr.net/npm/${CDN_REVEAL_VERSION}/dist/theme/black.css`
export const CDN_REVEAL_JS = `https://cdn.jsdelivr.net/npm/${CDN_REVEAL_VERSION}/dist/reveal.js`

export const DEFAULT_SLIDE_CONTENT = '<h1>New Slide</h1>'

export const ANNOTATION_COLORS = ['#ff0000', '#00ff00', '#0000ff', '#ffff00', '#ff00ff', '#00ffff', '#ffffff']
export const ANNOTATION_DEFAULT_COLOR = '#ff0000'
export const ANNOTATION_DEFAULT_THICKNESS = 3
export const ANNOTATION_MIN_THICKNESS = 1
export const ANNOTATION_MAX_THICKNESS = 10
