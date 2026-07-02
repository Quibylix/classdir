import { useState, useEffect } from 'react'
import { getPresentation, updatePresentationTitle } from '../api'
import type { Presentation } from '../types'
import type { FetchError } from '../../shared/api/fetch'

export function usePresentation(id: string) {
  const [presentation, setPresentation] = useState<Presentation | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [isSaving, setIsSaving] = useState(false)
  const [error, setError] = useState<FetchError | null>(null)

  useEffect(() => {
    setIsLoading(true)
    setError(null)
    getPresentation(id).match(
      (data) => { setPresentation(data); setIsLoading(false) },
      (e) => { setError(e); setIsLoading(false) },
    )
  }, [id])

  function updateTitle(title: string) {
    if (isSaving) return

    setIsSaving(true)
    setError(null)
    return updatePresentationTitle(id, title).match(
      (data) => { setPresentation(data); setIsSaving(false) },
      (e) => { setError(e); setIsSaving(false) },
    )
  }

  return { presentation, isLoading, isSaving, error, updateTitle }
}
