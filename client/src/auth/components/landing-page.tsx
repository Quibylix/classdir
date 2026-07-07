import { useEffect } from 'react'
import { useNavigate, Link } from 'react-router'
import { Box, Button, Center, Container, Grid, Group, Paper, Stack, Text, Title } from '@mantine/core'
import { PresentationIcon } from '@phosphor-icons/react/dist/csr/Presentation'
import { FolderOpenIcon } from '@phosphor-icons/react/dist/csr/FolderOpen'
import { CodeBlockIcon } from '@phosphor-icons/react/dist/csr/CodeBlock'
import { WifiHighIcon } from '@phosphor-icons/react/dist/csr/WifiHigh'
import { LockKeyIcon } from '@phosphor-icons/react/dist/csr/LockKey'
import { useAuth } from '../hooks/use-auth'
import classes from './landing-page.module.css'
import { CLIENT_LOGIN, CLIENT_CONFIGURE, CLIENT_PRESENT } from '../../shared/cfg/routes'

const FEATURES = [
  {
    icon: FolderOpenIcon,
    color: 'var(--mantine-color-blue-4)',
    title: 'Organize Your Lessons',
    desc: 'Keep all your presentations in one place. Add, rename, or remove them in a click.',
  },
  {
    icon: CodeBlockIcon,
    color: 'var(--mantine-color-indigo-4)',
    title: 'Built-in Slide Editor',
    desc: 'Write your slide content and see a live preview as you type. Simple and fast.',
  },
  {
    icon: WifiHighIcon,
    color: 'var(--mantine-color-teal-4)',
    title: 'Live Sync',
    desc: 'Advance a slide on your device and every student follows instantly. No delays, no confusion.',
  },
  {
    icon: LockKeyIcon,
    color: 'var(--mantine-color-yellow-4)',
    title: 'One Password, Done',
    desc: 'Just set a password and start. No student accounts, no email verification, no hassle.',
  },
]

export function LandingPage() {
  const { isAuthenticated, isLoading } = useAuth()
  const navigate = useNavigate()

  useEffect(() => {
    if (!isLoading && isAuthenticated) {
      navigate(CLIENT_CONFIGURE)
    }
  }, [isLoading, isAuthenticated, navigate])

  if (isLoading) return null

  return (
    <Box bg="dark.9" c="gray.3" mih="100vh">
      <Box
        component="header"
        pos="fixed"
        top={0}
        left={0}
        right={0}
        bg="color-mix(in srgb, var(--mantine-color-dark-9) 80%, transparent)"
        style={{
          zIndex: 40,
          backdropFilter: 'blur(12px)',
          borderBottom: '1px solid var(--mantine-color-dark-8)',
        }}
      >
        <Container size="lg">
          <Group h={64} justify="space-between">
            <Group gap="xs">
              <Center w={36} h={36} bg="blue.6" style={{ borderRadius: 8, boxShadow: '0 4px 12px color-mix(in srgb, var(--mantine-color-blue-5) 20%, transparent)' }}>
                <PresentationIcon size={20} color="white" />
              </Center>
              <Text fw={700} size="xl" c="white" lts="-0.02em">
                Class<Text fw="inherit" component="span" c="blue.4">Dir</Text>
              </Text>
            </Group>

            <Group gap="lg" visibleFrom="sm">
              <Box
                component="a"
                href="#features"
                fz={14}
                fw={500}
                c="gray.5"
                td="none"
                style={{ transition: 'color 0.15s', cursor: 'pointer' }}
              >
                Features
              </Box>
              <Button
                component={Link}
                to={isAuthenticated ? CLIENT_CONFIGURE : CLIENT_LOGIN}
                variant="outline"
                color="gray"
                size="sm"
                style={{ borderColor: 'var(--mantine-color-dark-6)' }}
              >
                {isAuthenticated ? 'Configure' : 'Admin Access'}
              </Button>
            </Group>
          </Group>
        </Container>
      </Box>

      <Box component="section" pos="relative" pt={128} pb={96} style={{ overflow: 'hidden' }}>
        <Box
          pos="absolute"
          top="-10%"
          left="20%"
          w={500}
          h={500}
          bg="blue.5"
          opacity={0.08}
          style={{ borderRadius: '50%', filter: 'blur(120px)', pointerEvents: 'none' }}
        />
        <Box
          pos="absolute"
          top="20%"
          right="10%"
          w={400}
          h={400}
          bg="indigo.5"
          opacity={0.08}
          style={{ borderRadius: '50%', filter: 'blur(100px)', pointerEvents: 'none' }}
        />

        <Container size="lg">
          <Stack align="center" gap="lg" ta="center">
            <Group
              display="inline-flex"
              px={12}
              py={4}
              fz={12}
              fw={500}
              bg="color-mix(in srgb, var(--mantine-color-blue-5) 10%, transparent)"
              c="var(--mantine-color-blue-4)"
              bd="1px solid color-mix(in srgb, var(--mantine-color-blue-5) 20%, transparent)"
              gap={6}
              style={{ borderRadius: 999 }}
            >
              <Box
                w={6}
                h={6}
                bg="teal.4"
                className={classes.dot}
              />
              No sign-up required
            </Group>

            <Title
              order={1}
              c="white"
              maw={800}
              fz="clamp(2rem, 6vw, 3.75rem)"
              fw={800}
              lts="-0.03em"
              lh={1.15}
            >
              Your classroom presentations,{' '}
              <span
                style={{
                  background: 'linear-gradient(to right, var(--mantine-color-blue-4), var(--mantine-color-blue-5))',
                  WebkitBackgroundClip: 'text',
                  WebkitTextFillColor: 'transparent',
                }}
              >
                synced in real time
              </span>
            </Title>

            <Text size="lg" c="dimmed" maw={560} lh={1.6}>
              Create beautiful lesson slides and guide your students through them in real time — from any device, no accounts needed.
            </Text>

            <Group mt="md" gap="md" justify="center">
              <Button
                component={Link}
                to={isAuthenticated ? CLIENT_CONFIGURE : CLIENT_LOGIN}
                size="lg"
                color="blue"
                fw={600}
                style={{ boxShadow: '0 4px 16px color-mix(in srgb, var(--mantine-color-blue-6) 20%, transparent)' }}
              >
                Get Started
              </Button>
              <Button
                component={Link}
                to={CLIENT_PRESENT}
                size="lg"
                variant="outline"
                color="gray"
                c="gray.5"
                style={{ borderColor: 'var(--mantine-color-dark-6)' }}
              >
                Join a Presentation
              </Button>
              <Button
                component="a"
                href="#features"
                variant="subtle"
                size="lg"
                c="gray.5"
                style={{ borderColor: 'var(--mantine-color-dark-6)' }}
              >
                Learn More
              </Button>
            </Group>
          </Stack>
        </Container>
      </Box>

      <Box component="section" id="features" py={80}>
        <Container size="lg">
          <Stack align="center" gap="xs" mb={64}>
            <Title order={2} c="white" ta="center" fz="clamp(1.5rem, 4vw, 2.25rem)" fw={700} lts="-0.02em">
              Everything you need to teach better
            </Title>
            <Text c="dimmed" maw={480} ta="center">
              ClassDir helps you prepare and deliver your classes without fighting with complicated tools.
            </Text>
          </Stack>

          <Grid>
            {FEATURES.map((feature) => (
              <Grid.Col key={feature.title} span={{ base: 12, sm: 6, lg: 3 }}>
                <Paper p="lg" radius="md" bg="dark.8" className={classes.card} h="100%">
                  <Stack gap="md">
                    <Box className={classes.iconBox} style={{ color: feature.color }}>
                      <feature.icon size={22} />
                    </Box>
                    <Title order={4} c="white" fz={18} fw={600}>
                      {feature.title}
                    </Title>
                    <Text size="sm" c="dimmed" lh={1.6}>
                      {feature.desc}
                    </Text>
                  </Stack>
                </Paper>
              </Grid.Col>
            ))}
          </Grid>
        </Container>
      </Box>

      <Box component="footer" py="lg" style={{ borderTop: '1px solid var(--mantine-color-dark-8)' }}>
        <Container size="lg">
          <Group justify="space-between" fz={13} c="gray.6">
            <Group gap="xs">
              <Text fw={600} c="gray.5">ClassDir</Text>
              <Text component="span" c="dimmed">—</Text>
              <Text c="dimmed">Teaching made simpler.</Text>
            </Group>
            <Text c="dimmed">{new Date().getFullYear()} ClassDir. Built for digital education.</Text>
          </Group>
        </Container>
      </Box>
    </Box>
  )
}
