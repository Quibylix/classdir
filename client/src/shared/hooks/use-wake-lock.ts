import { useEffect } from 'react'

export function useWakeLock() {
  useEffect(() => {
    let sentinel: WakeLockSentinel | null = null

    async function request() {
      try {
        sentinel = await navigator.wakeLock.request('screen')
      } catch {
        // silently ignore — may fail without user gesture or on unsupported browsers
      }
    }

    request()

    function onVisibilityChange() {
      if (document.visibilityState === 'visible' && !sentinel) {
        request()
      }
    }

    document.addEventListener('visibilitychange', onVisibilityChange)

    return () => {
      document.removeEventListener('visibilitychange', onVisibilityChange)
      sentinel?.release()
      sentinel = null
    }
  }, [])

  return { supported: 'wakeLock' in navigator }
}
