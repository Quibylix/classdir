import { useState, useEffect, useCallback } from 'react'
import { listPresentations, createPresentation, deletePresentation } from '../api'
import { uuidv7 } from '../../shared/util/uuid'
import type { PresentationPreview } from '../types'
import type { FetchError } from '../../shared/api/fetch'

export function usePresentationList() {
  const [presentations, setPresentations] = useState<PresentationPreview[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [isCreating, setIsCreating] = useState(false)
  const [isDeleting, setIsDeleting] = useState(false)
  const [error, setError] = useState<FetchError | null>(null)

  const refresh = useCallback(() => {
    setIsLoading(true)
    setError(null)
    listPresentations().match(
      (data) => { setPresentations(data); setIsLoading(false) },
      (e) => { setError(e); setIsLoading(false) },
    )
  }, [])

  useEffect(() => { refresh() }, [refresh])

  function create(title: string) {
    if (isCreating) return

    const id = uuidv7()
    setIsCreating(true)
    setError(null)
    return createPresentation(id, title).match(
      () => { setIsCreating(false); refresh() },
      (e) => { setIsCreating(false); setError(e) },
    )
  }

  function remove(id: string) {
    if (isDeleting) return

    setIsDeleting(true)
    setError(null)
    return deletePresentation(id).match(
      () => { setIsDeleting(false); refresh() },
      (e) => { setIsDeleting(false); setError(e) },
    )
  }

  return { presentations, isLoading, isCreating, isDeleting, error, refresh, create, remove }
}
