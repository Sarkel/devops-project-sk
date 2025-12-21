import { describe, it, expect } from 'vitest'
import { screen } from '@testing-library/react'
import { renderWithChakra } from '../test/testUtils'
import TemperatureChart from './TemperatureChart'

describe('TemperatureChart', () => {
  const mockData = {
    labels: ['2025-09-15 10:00', '2025-09-15 11:00', '2025-09-15 12:00'],
    datasets: [
      {
        label: 'API Sensor',
        data: [22.5, 23.0, 23.5],
        borderColor: 'rgb(59, 130, 246)',
        backgroundColor: 'rgba(59, 130, 246, 0.5)',
      },
      {
        label: 'Local Sensor',
        data: [22.0, 22.5, 23.0],
        borderColor: 'rgb(34, 197, 94)',
        backgroundColor: 'rgba(34, 197, 94, 0.5)',
      },
    ],
  }

  it('renders chart heading', () => {
    renderWithChakra(<TemperatureChart data={mockData} />)
    expect(screen.getByText('Temperature History')).toBeInTheDocument()
  })

  it('renders chart with data', () => {
    const { container } = renderWithChakra(<TemperatureChart data={mockData} />)
    const canvas = container.querySelector('canvas')
    expect(canvas).toBeInTheDocument()
  })

  it('renders chart with empty data', () => {
    const emptyData = { labels: [], datasets: [] }
    const { container } = renderWithChakra(<TemperatureChart data={emptyData} />)
    const canvas = container.querySelector('canvas')
    expect(canvas).toBeInTheDocument()
  })

  it('renders white background card', () => {
    const { container } = renderWithChakra(<TemperatureChart data={mockData} />)
    expect(container).toBeInTheDocument()
  })
})
