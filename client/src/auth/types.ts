export const LOGIN_RESULT = {
  Ok: 'ok',
  Invalid: 'invalid',
  Error: 'error',
} as const

export type LoginResult = (typeof LOGIN_RESULT)[keyof typeof LOGIN_RESULT]
