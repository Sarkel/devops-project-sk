import { describe, it, expect } from 'vitest'
import { render, waitFor } from '@testing-library/react'
import { ChakraProvider, defaultSystem } from '@chakra-ui/react'
import App from './App'

describe('App', () => {
  it('renders without crashing', async () => {
    const { container } = render(
      <ChakraProvider value={defaultSystem}>
        <App />
      </ChakraProvider>
    )

    await waitFor(() => {
      expect(container).toBeTruthy()
    })
  })
})
