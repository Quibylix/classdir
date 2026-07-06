export const API_PREFIX = '/api/v1'

export const AUTH_LOGIN = `${API_PREFIX}/auth/login`
export const AUTH_LOGOUT = `${API_PREFIX}/auth/logout`
export const AUTH_CHECK = `${API_PREFIX}/auth/check`

export const CLIENT_CONFIGURE = '/configure'

export const PRESENTATIONS = `${API_PREFIX}/presentation`
export const presentationById = (id: string) => `${API_PREFIX}/presentation/${id}`
export const presentationSlides = (presId: string) => `${API_PREFIX}/presentation/${presId}/slides`
export const slideById = (presId: string, slideId: string) => `${API_PREFIX}/presentation/${presId}/slides/${slideId}`

export const WS_V1 = '/ws/v1'

export const CLIENT_PRESENT = '/present'
export const CLIENT_CONTROL = '/control'
