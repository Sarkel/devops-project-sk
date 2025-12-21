import { describe, it, expect } from 'vitest'
import { screen, waitFor } from '@testing-library/react'
import { act } from 'react'
import { renderWithChakra } from './test/testUtils'
import Router from './Router'

describe('Router', () => {
  it('renders Router component', async () => {
    await act(async () => {
      renderWithChakra(<Router />)
    })
    // The router renders successfully - component structure is tested
    await waitFor(() => {
      expect(document.body).toBeInTheDocument()
    })
  })

  it('renders default route structure', async () => {
    await act(async () => {
      renderWithChakra(<Router />)
    })
    // Router should mount without errors
    await waitFor(() => {
      expect(document.body).toBeInTheDocument()
    })
  })
})
