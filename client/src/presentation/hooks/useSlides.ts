import { useState, useEffect, useCallback } from 'react'
import { createSlide, updateSlide, deleteSlide } from '../api'
import { uuidv7 } from '../../shared/util/uuid'
import type { Slide } from '../types'
import type { FetchError } from '../../shared/api/fetch'

export function useSlides(presId: string, initialSlides?: Slide[]) {
  const [slides, setSlides] = useState<Slide[]>(initialSlides ?? [])
  const [currentIndex, setCurrentIndex] = useState(0)
  const [isAdding, setIsAdding] = useState(false)
  const [isSaving, setIsSaving] = useState(false)
  const [isDeleting, setIsDeleting] = useState(false)
  const [error, setError] = useState<FetchError | null>(null)

  useEffect(() => {
    setSlides(initialSlides ?? [])
  }, [initialSlides])

  useEffect(() => {
    if (currentIndex >= slides.length) {
      setCurrentIndex(Math.max(0, slides.length - 1))
    }
  }, [slides, currentIndex])

  const addSlide = useCallback(() => {
    setIsAdding(true)
    setError(null)
    const id = uuidv7()
    createSlide(presId, id, '<h1>New Slide</h1>').match(
      (slide) => {
        setSlides(prev => [...prev, slide])
        setCurrentIndex(slides.length)
        setIsAdding(false)
      },
      (e) => {
        setError(e)
        setIsAdding(false)
      },
    )
  }, [presId, slides.length])

  const saveSlide = useCallback((index: number, content: string) => {
    const slide = slides[index]
    if (!slide) return
    setIsSaving(true)
    setError(null)
    updateSlide(presId, slide.id, content).match(
      (updated) => {
        setSlides(prev => prev.map((s, i) => i === index ? updated : s))
        setIsSaving(false)
      },
      (e) => {
        setError(e)
        setIsSaving(false)
      },
    )
  }, [presId, slides])

  const removeSlide = useCallback((index: number) => {
    const slide = slides[index]
    if (!slide) return
    setIsDeleting(true)
    setError(null)
    deleteSlide(presId, slide.id).match(
      () => {
        setSlides(prev => prev.filter((_, i) => i !== index))
        setIsDeleting(false)
      },
      (e) => {
        setError(e)
        setIsDeleting(false)
      },
    )
  }, [presId, slides])

  const goToSlide = useCallback((index: number) => {
    if (index >= 0 && index < slides.length) {
      setCurrentIndex(index)
    }
  }, [slides.length])

  return {
    slides,
    currentSlide: slides[currentIndex] ?? null,
    currentIndex,
    isAdding,
    isSaving,
    isDeleting,
    error,
    addSlide,
    saveSlide,
    removeSlide,
    goToSlide,
  }
}
