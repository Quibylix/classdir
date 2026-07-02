import { z } from 'zod'
import { safeFetch } from '../shared/api/fetch'
import { PRESENTATIONS, presentationById } from '../shared/cfg/routes'
import { PresentationSchema, PresentationPreviewSchema } from './types'

export function listPresentations() {
  return safeFetch(PRESENTATIONS, z.array(PresentationPreviewSchema))
}

export function getPresentation(id: string) {
  return safeFetch(presentationById(id), PresentationSchema)
}

export function createPresentation(id: string, title: string) {
  return safeFetch(PRESENTATIONS, PresentationSchema, {
    method: 'POST',
    body: JSON.stringify({ id, title }),
  })
}

export function updatePresentationTitle(id: string, title: string) {
  return safeFetch(presentationById(id), PresentationSchema, {
    method: 'PUT',
    body: JSON.stringify({ title }),
  })
}

export function deletePresentation(id: string) {
  return safeFetch(presentationById(id), z.undefined(), { method: 'DELETE' })
}
