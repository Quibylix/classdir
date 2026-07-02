import '@mantine/core/styles.css'
import { MantineProvider } from '@mantine/core'
import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { BrowserRouter, Routes, Route } from 'react-router'
import { LandingPage } from './auth/components/LandingPage'
import { PresentationListPage } from './presentation/components/PresentationListPage'
import { PresentationDetailPage } from './presentation/components/PresentationDetailPage'

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <BrowserRouter>
      <MantineProvider>
        <Routes>
          <Route path="/" element={<LandingPage />} />
          <Route path="/configure" element={<PresentationListPage />} />
          <Route path="/configure/:id" element={<PresentationDetailPage />} />
        </Routes>
      </MantineProvider>
    </BrowserRouter>
  </StrictMode>,
)
