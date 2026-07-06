import { Center, Loader } from '@mantine/core'
import { useNavigate } from 'react-router'
import { useAuth } from '../hooks/use-auth'
import { ROOT } from '../../shared/cfg/routes'
import { useEffect } from 'react';

export function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const { isAuthenticated, isLoading } = useAuth()
  const navigate = useNavigate();

  useEffect(() => {
    if (!isLoading && !isAuthenticated) {
      navigate(ROOT, { replace: true });
    }
  }, [isAuthenticated, isLoading, navigate]);

  if (isLoading) return <Center h="100vh"><Loader /></Center>

  return isAuthenticated ? children : null;
}
