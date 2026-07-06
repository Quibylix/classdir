import { z } from 'zod'
import { safeFetch } from '../shared/api/fetch'
import { PRESENTATIONS, presentationById, presentationSlides, slideById } from '../shared/cfg/routes'
import { HTTP_METHOD_POST, HTTP_METHOD_PUT, HTTP_METHOD_DELETE } from '../shared/cfg/http'
import { PresentationSchema, PresentationPreviewSchema, SlideSchema } from './types'

export function listPresentations() {
  return safeFetch(PRESENTATIONS, z.array(PresentationPreviewSchema))
}

export function getPresentation(id: string) {
  return safeFetch(presentationById(id), PresentationSchema)
}

export function createPresentation(id: string, title: string) {
  return safeFetch(PRESENTATIONS, PresentationSchema, {
    method: HTTP_METHOD_POST,
    body: JSON.stringify({ id, title }),
  })
}

export function updatePresentationTitle(id: string, title: string) {
  return safeFetch(presentationById(id), PresentationSchema, {
    method: HTTP_METHOD_PUT,
    body: JSON.stringify({ title }),
  })
}

export function deletePresentation(id: string) {
  return safeFetch(presentationById(id), z.undefined(), { method: HTTP_METHOD_DELETE })
}

export function createSlide(presId: string, id: string, content: string) {
  return safeFetch(presentationSlides(presId), SlideSchema, {
    method: HTTP_METHOD_POST,
    body: JSON.stringify({ id, content }),
  })
}

export function getSlide(presId: string, slideId: string) {
  return safeFetch(slideById(presId, slideId), SlideSchema)
}

export function updateSlide(presId: string, slideId: string, content: string) {
  return safeFetch(slideById(presId, slideId), SlideSchema, {
    method: HTTP_METHOD_PUT,
    body: JSON.stringify({ content }),
  })
}

export function deleteSlide(presId: string, slideId: string) {
  return safeFetch(slideById(presId, slideId), z.undefined(), { method: HTTP_METHOD_DELETE })
}
