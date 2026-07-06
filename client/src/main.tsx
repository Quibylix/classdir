import '@mantine/core/styles.css'
import { MantineProvider } from '@mantine/core'
import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { BrowserRouter, Routes, Route } from 'react-router'
import { AuthProvider } from './auth/auth-context'
import { ProtectedRoute } from './auth/components/protected-route'
import { LandingPage } from './auth/components/landing-page'
import { PresentationListPage } from './presentation/components/presentation-list-page'
import { PresentationDetailPage } from './presentation/components/presentation-detail-page'
import { PresentView } from './presentation/components/present-view'
import { ControlView } from './presentation/components/control-view'

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <AuthProvider>
      <BrowserRouter>
        <MantineProvider>
          <Routes>
            <Route path="/" element={<LandingPage />} />
            <Route path="/configure" element={<ProtectedRoute><PresentationListPage /></ProtectedRoute>} />
            <Route path="/configure/:id" element={<ProtectedRoute><PresentationDetailPage /></ProtectedRoute>} />
            <Route path="/present/:id" element={<ProtectedRoute><PresentView /></ProtectedRoute>} />
            <Route path="/control/:id" element={<ProtectedRoute><ControlView /></ProtectedRoute>} />
          </Routes>
        </MantineProvider>
      </BrowserRouter>
    </AuthProvider>
  </StrictMode >,
)
