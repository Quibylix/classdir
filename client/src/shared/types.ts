export const WS_STATUS = {
  Connecting: 'connecting',
  Connected: 'connected',
  Disconnected: 'disconnected',
} as const

export type WSStatus = (typeof WS_STATUS)[keyof typeof WS_STATUS]
