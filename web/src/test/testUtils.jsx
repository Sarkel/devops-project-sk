import { render } from '@testing-library/react'
import { ChakraProvider, defaultSystem } from '@chakra-ui/react'

export function renderWithChakra(ui, options) {
  return render(
    <ChakraProvider value={defaultSystem}>
      {ui}
    </ChakraProvider>,
    options
  )
}

export * from '@testing-library/react'
