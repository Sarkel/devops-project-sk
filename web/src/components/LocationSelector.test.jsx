import { describe, it, expect, vi } from 'vitest'
import { screen, fireEvent } from '@testing-library/react'
import { renderWithChakra } from '../test/testUtils'
import LocationSelector from './LocationSelector'

describe('LocationSelector', () => {
  const mockLocations = [
    { sid: 'loc1', name: 'Location 1' },
    { sid: 'loc2', name: 'Location 2' },
    { sid: 'loc3', name: 'Location 3' },
  ]

  const mockProps = {
    locations: mockLocations,
    selectedLocation: 'loc1',
    onLocationChange: vi.fn(),
    onRefresh: vi.fn(),
  }

  it('renders heading', () => {
    renderWithChakra(<LocationSelector {...mockProps} />)
    expect(screen.getByText('Temperature Checker')).toBeInTheDocument()
  })

  it('renders all location options', () => {
    renderWithChakra(<LocationSelector {...mockProps} />)
    mockLocations.forEach(loc => {
      expect(screen.getByText(loc.name)).toBeInTheDocument()
    })
  })

  it('shows selected location', () => {
    renderWithChakra(<LocationSelector {...mockProps} />)
    const select = screen.getByRole('combobox')
    expect(select).toHaveValue('loc1')
  })

  it('calls onLocationChange when location is changed', () => {
    renderWithChakra(<LocationSelector {...mockProps} />)
    const select = screen.getByRole('combobox')
    fireEvent.change(select, { target: { value: 'loc2' } })
    expect(mockProps.onLocationChange).toHaveBeenCalledWith('loc2')
  })

  it('renders refresh button', () => {
    renderWithChakra(<LocationSelector {...mockProps} />)
    const refreshButton = screen.getByLabelText('Refresh')
    expect(refreshButton).toBeInTheDocument()
  })

  it('calls onRefresh when refresh button is clicked', () => {
    renderWithChakra(<LocationSelector {...mockProps} />)
    const refreshButton = screen.getByLabelText('Refresh')
    fireEvent.click(refreshButton)
    expect(mockProps.onRefresh).toHaveBeenCalled()
  })

  it('renders with empty locations array', () => {
    const emptyProps = { ...mockProps, locations: [] }
    renderWithChakra(<LocationSelector {...emptyProps} />)
    expect(screen.getByText('Temperature Checker')).toBeInTheDocument()
  })
})
