import '@mantine/core/styles.css'
import { MantineProvider } from '@mantine/core'
import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { BrowserRouter, Routes, Route } from 'react-router'
import { AuthProvider } from './auth/auth-context'
import { ProtectedRoute } from './auth/components/ProtectedRoute'
import { LandingPage } from './auth/components/LandingPage'
import { PresentationListPage } from './presentation/components/PresentationListPage'
import { PresentationDetailPage } from './presentation/components/PresentationDetailPage'

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <AuthProvider>
      <BrowserRouter>
        <MantineProvider>
          <Routes>
            <Route path="/" element={<LandingPage />} />
            <Route path="/configure" element={<ProtectedRoute><PresentationListPage /></ProtectedRoute>} />
            <Route path="/configure/:id" element={<ProtectedRoute><PresentationDetailPage /></ProtectedRoute>} />
          </Routes>
        </MantineProvider>
      </BrowserRouter>
    </AuthProvider>
  </StrictMode >,
)
