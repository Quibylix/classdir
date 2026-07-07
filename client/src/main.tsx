import '@mantine/core/styles.css'
import { MantineProvider } from '@mantine/core'
import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { BrowserRouter, Routes, Route } from 'react-router'
import { AuthProvider } from './auth/auth-context'
import { ProtectedRoute } from './auth/components/protected-route'
import { LandingPage } from './auth/components/landing-page'
import { LoginPage } from './auth/components/login-page'
import { PresentationListPage } from './presentation/components/presentation-list-page'
import { PresentationDetailPage } from './presentation/components/presentation-detail-page'
import { PresentView } from './presentation/components/present-view'
import { ControlView } from './presentation/components/control-view'
import { ROOT, CLIENT_LOGIN, CLIENT_CONFIGURE, clientConfigure, clientPresent, clientControl } from './shared/cfg/routes'

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <AuthProvider>
      <BrowserRouter>
        <MantineProvider defaultColorScheme="dark">
          <Routes>
            <Route path={ROOT} element={<LandingPage />} />
            <Route path={CLIENT_LOGIN} element={<LoginPage />} />
            <Route path={CLIENT_CONFIGURE} element={<ProtectedRoute><PresentationListPage /></ProtectedRoute>} />
            <Route path={clientConfigure(':id')} element={<ProtectedRoute><PresentationDetailPage /></ProtectedRoute>} />
            <Route path={clientPresent(':code')} element={<PresentView />} />
            <Route path={clientControl(':id')} element={<ProtectedRoute><ControlView /></ProtectedRoute>} />
          </Routes>
        </MantineProvider>
      </BrowserRouter>
    </AuthProvider>
  </StrictMode >,
)
