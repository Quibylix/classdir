export const API_PREFIX = '/api/v1'

export const AUTH_LOGIN = `${API_PREFIX}/auth/login`
export const AUTH_LOGOUT = `${API_PREFIX}/auth/logout`
export const AUTH_CHECK = `${API_PREFIX}/auth/check`

export const CLIENT_CONFIGURE = '/configure'

export const PRESENTATIONS = `${API_PREFIX}/presentation`
export const presentationById = (id: string) => `${API_PREFIX}/presentation/${id}`
