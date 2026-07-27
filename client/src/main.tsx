import "@mantine/core/styles.css";
import { Center, Container, Loader, MantineProvider } from "@mantine/core";
import { lazy, StrictMode, Suspense } from "react";
import { createRoot } from "react-dom/client";
import { BrowserRouter, Routes, Route } from "react-router";
import { AuthProvider } from "./auth/auth-provider";
import { ProtectedRoute } from "./auth/components/protected-route";
import {
  ROOT,
  CLIENT_LOGIN,
  CLIENT_CONFIGURE,
  CLIENT_PRESENT,
  clientConfigure,
  clientPresent,
  clientControl,
} from "./shared/cfg/routes";

export const LandingPageLazy = lazy(() => import("./auth/components/landing-page"));
export const LoginPageLazy = lazy(() => import("./auth/components/login-page"));
export const PresentationListPageLazy = lazy(
  () => import("./presentation/components/presentation-list-page"),
);
export const PresentationDetailPageLazy = lazy(
  () => import("./presentation/components/presentation-detail-page"),
);
export const PresentViewLazy = lazy(() => import("./presentation/components/present-view"));
export const PresentCodeEntryLazy = lazy(
  () => import("./presentation/components/present-code-entry"),
);
export const ControlViewLazy = lazy(() => import("./presentation/components/control-view"));

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <AuthProvider>
      <BrowserRouter>
        <MantineProvider defaultColorScheme="dark">
          <Suspense
            fallback={
              <Container fluid h="100dvh">
                <Center w="100%" h="100%">
                  <Loader />
                </Center>
              </Container>
            }
          >
            <Routes>
              <Route path={ROOT} element={<LandingPageLazy />} />
              <Route path={CLIENT_LOGIN} element={<LoginPageLazy />} />
              <Route
                path={CLIENT_CONFIGURE}
                element={
                  <ProtectedRoute>
                    <PresentationListPageLazy />
                  </ProtectedRoute>
                }
              />
              <Route
                path={clientConfigure(":id")}
                element={
                  <ProtectedRoute>
                    <PresentationDetailPageLazy />
                  </ProtectedRoute>
                }
              />
              <Route path={clientPresent(":code")} element={<PresentViewLazy />} />
              <Route path={CLIENT_PRESENT} element={<PresentCodeEntryLazy />} />
              <Route
                path={clientControl(":id")}
                element={
                  <ProtectedRoute>
                    <ControlViewLazy />
                  </ProtectedRoute>
                }
              />
            </Routes>
          </Suspense>
        </MantineProvider>
      </BrowserRouter>
    </AuthProvider>
  </StrictMode>,
);
