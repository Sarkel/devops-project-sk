import { describe, it, expect, vi } from 'vitest'
import { screen, fireEvent, waitFor } from '@testing-library/react'
import { act } from 'react'
import { renderWithChakra } from '../test/testUtils'
import Sensors from './Sensors'

// Mock hooks
vi.mock('../hooks/useLocations', () => ({
  useLocations: () => ({
    locations: [
      { sid: 'loc1', name: 'Location 1' },
      { sid: 'loc2', name: 'Location 2' },
    ],
    selectedLocation: 'loc1',
    setSelectedLocation: vi.fn(),
  }),
}))

vi.mock('../hooks/useSensorSummary', () => ({
  useSensorSummary: () => ({
    api: {
      temperature: 23.5,
      trend: 1.2,
      timestamp: '2025-09-15T12:00:00Z',
    },
    local: {
      temperature: 22.8,
      trend: -0.5,
      timestamp: '2025-09-15T12:00:00Z',
    },
  }),
}))

vi.mock('../hooks/useSensorData', () => ({
  useSensorData: () => ({
    labels: ['2025-09-15 10:00', '2025-09-15 11:00'],
    datasets: [
      {
        label: 'API Sensor',
        data: [22.5, 23.0],
        borderColor: 'rgb(59, 130, 246)',
        backgroundColor: 'rgba(59, 130, 246, 0.5)',
      },
    ],
  }),
}))

// Mock components
vi.mock('../components/LocationSelector', () => ({
  default: ({ locations, onRefresh }) => (
    <div data-testid="location-selector">
      <button onClick={onRefresh}>Refresh</button>
      {locations.map(loc => <div key={loc.sid}>{loc.name}</div>)}
    </div>
  ),
}))

vi.mock('../components/SensorCard', () => ({
  default: ({ title, data }) => (
    <div data-testid="sensor-card">
      {title}: {data ? data.temperature : 'No data'}
    </div>
  ),
}))

vi.mock('../components/ChartFilters', () => ({
  default: ({ onResetFilters, onTypeToggle }) => (
    <div data-testid="chart-filters">
      <button onClick={onResetFilters}>Reset</button>
      <button onClick={() => onTypeToggle('api')}>Toggle API</button>
    </div>
  ),
}))

vi.mock('../components/TemperatureChart', () => ({
  default: () => <div data-testid="temperature-chart">Chart</div>,
}))

describe('Sensors', () => {
  it('renders all components', () => {
    renderWithChakra(<Sensors />)
    expect(screen.getByTestId('location-selector')).toBeInTheDocument()
    expect(screen.getAllByTestId('sensor-card')).toHaveLength(2)
    expect(screen.getByTestId('chart-filters')).toBeInTheDocument()
    expect(screen.getByTestId('temperature-chart')).toBeInTheDocument()
  })

  it('displays API sensor data', () => {
    renderWithChakra(<Sensors />)
    expect(screen.getByText(/API Sensor: 23.5/)).toBeInTheDocument()
  })

  it('displays Local sensor data', () => {
    renderWithChakra(<Sensors />)
    expect(screen.getByText(/Local Sensor: 22.8/)).toBeInTheDocument()
  })

  it('handles refresh button click', async () => {
    renderWithChakra(<Sensors />)
    const refreshButton = screen.getByText('Refresh')
    await act(async () => {
      fireEvent.click(refreshButton)
    })
    // Component should re-render without errors
    await waitFor(() => {
      expect(screen.getByTestId('location-selector')).toBeInTheDocument()
    })
  })

  it('handles reset filters click', async () => {
    renderWithChakra(<Sensors />)
    const resetButton = screen.getByText('Reset')
    await act(async () => {
      fireEvent.click(resetButton)
    })
    // Component should handle reset without errors
    await waitFor(() => {
      expect(screen.getByTestId('chart-filters')).toBeInTheDocument()
    })
  })

  it('handles type toggle', async () => {
    renderWithChakra(<Sensors />)
    const toggleButton = screen.getByText('Toggle API')
    await act(async () => {
      fireEvent.click(toggleButton)
    })
    // Component should handle toggle without errors
    await waitFor(() => {
      expect(screen.getByTestId('chart-filters')).toBeInTheDocument()
    })
  })
})
