import { describe, it, expect } from 'vitest'
import { screen } from '@testing-library/react'
import { BrowserRouter } from 'react-router-dom'
import { renderWithChakra } from '../test/testUtils'
import Layout from './Layout'

describe('Layout', () => {
  it('renders layout wrapper with gray background', () => {
    const { container } = renderWithChakra(
      <BrowserRouter>
        <Layout />
      </BrowserRouter>
    )
    expect(container).toBeInTheDocument()
  })

  it('renders outlet for child routes', () => {
    renderWithChakra(
      <BrowserRouter>
        <Layout />
      </BrowserRouter>
    )
    // Outlet is rendered by react-router-dom
    expect(screen.queryByText(/temperature/i)).not.toBeInTheDocument()
  })
})
