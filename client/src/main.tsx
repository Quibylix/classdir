import '@mantine/core/styles.css'
import { MantineProvider } from '@mantine/core'
import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { LandingPage } from './auth/components/LandingPage'

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <MantineProvider>
      <LandingPage />
    </MantineProvider>
  </StrictMode>,
)
