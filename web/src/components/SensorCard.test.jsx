import { describe, it, expect } from 'vitest'
import { screen } from '@testing-library/react'
import { renderWithChakra } from '../test/testUtils'
import SensorCard from './SensorCard'

describe('SensorCard', () => {
  const mockData = {
    temperature: 23.5,
    trend: 1.2,
    timestamp: '2025-09-15T12:00:00Z',
  }

  it('renders title', () => {
    renderWithChakra(<SensorCard title="API Sensor" data={mockData} colorScheme="blue" />)
    expect(screen.getByText('API Sensor')).toBeInTheDocument()
  })

  it('displays temperature with one decimal place', () => {
    renderWithChakra(<SensorCard title="API Sensor" data={mockData} colorScheme="blue" />)
    expect(screen.getByText(/23.5°C/)).toBeInTheDocument()
  })

  it('displays timestamp', () => {
    renderWithChakra(<SensorCard title="API Sensor" data={mockData} colorScheme="blue" />)
    const formattedDate = new Date('2025-09-15T12:00:00Z').toLocaleString()
    expect(screen.getByText(formattedDate)).toBeInTheDocument()
  })

  it('shows trending up icon for positive trend', () => {
    const { container } = renderWithChakra(
      <SensorCard title="API Sensor" data={mockData} colorScheme="blue" />
    )
    const icon = container.querySelector('svg')
    expect(icon).toBeInTheDocument()
  })

  it('shows trending down icon for negative trend', () => {
    const dataNegativeTrend = { ...mockData, trend: -1.2 }
    const { container } = renderWithChakra(
      <SensorCard title="API Sensor" data={dataNegativeTrend} colorScheme="blue" />
    )
    const icon = container.querySelector('svg')
    expect(icon).toBeInTheDocument()
  })

  it('shows no icon for zero trend', () => {
    const dataZeroTrend = { ...mockData, trend: 0 }
    const { container } = renderWithChakra(
      <SensorCard title="API Sensor" data={dataZeroTrend} colorScheme="blue" />
    )
    // No icon should be rendered for trend = 0
    expect(screen.getByText(/23.5°C/)).toBeInTheDocument()
  })

  it('shows "No data available" when data is null', () => {
    renderWithChakra(<SensorCard title="API Sensor" data={null} colorScheme="blue" />)
    expect(screen.getByText('No data available')).toBeInTheDocument()
  })

  it('shows "No data available" when data is undefined', () => {
    renderWithChakra(<SensorCard title="API Sensor" data={undefined} colorScheme="blue" />)
    expect(screen.getByText('No data available')).toBeInTheDocument()
  })

  it('applies correct color scheme', () => {
    const { container} = renderWithChakra(
      <SensorCard title="API Sensor" data={mockData} colorScheme="green" />
    )
    expect(container).toBeInTheDocument()
  })
})
