import { Outlet } from 'react-router-dom'
import { Box, Container } from '@chakra-ui/react'

export default function Layout() {
  return (
    <Box minH="100vh" bg="gray.50" p={6}>
      <Container maxW="7xl">
        <Outlet />
      </Container>
    </Box>
  )
}
