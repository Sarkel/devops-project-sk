import { describe, it, expect, vi } from 'vitest'
import { screen, fireEvent } from '@testing-library/react'
import { renderWithChakra } from '../test/testUtils'
import ChartFilters from './ChartFilters'

describe('ChartFilters', () => {
  const mockProps = {
    startDate: '2025-09-15T00:00',
    endDate: '2025-09-15T23:59',
    aggregation: '',
    types: ['api', 'local'],
    onStartDateChange: vi.fn(),
    onEndDateChange: vi.fn(),
    onAggregationChange: vi.fn(),
    onTypeToggle: vi.fn(),
    onResetFilters: vi.fn(),
  }

  it('renders filters heading', () => {
    renderWithChakra(<ChartFilters {...mockProps} />)
    expect(screen.getByText('Chart Filters')).toBeInTheDocument()
  })

  it('renders start date input', () => {
    renderWithChakra(<ChartFilters {...mockProps} />)
    expect(screen.getByText('Start Date:')).toBeInTheDocument()
    const startDateInput = screen.getAllByDisplayValue('2025-09-15T00:00')[0]
    expect(startDateInput).toBeInTheDocument()
  })

  it('renders end date input', () => {
    renderWithChakra(<ChartFilters {...mockProps} />)
    expect(screen.getByText('End Date:')).toBeInTheDocument()
    const endDateInput = screen.getAllByDisplayValue('2025-09-15T23:59')[0]
    expect(endDateInput).toBeInTheDocument()
  })

  it('renders aggregation select', () => {
    renderWithChakra(<ChartFilters {...mockProps} />)
    expect(screen.getByText('Aggregation:')).toBeInTheDocument()
    const select = screen.getByRole('combobox')
    expect(select).toBeInTheDocument()
  })

  it('renders sensor type checkboxes', () => {
    renderWithChakra(<ChartFilters {...mockProps} />)
    expect(screen.getByText('Sensor Types:')).toBeInTheDocument()
    const apiCheckbox = screen.getByRole('checkbox', { name: /api/i })
    const localCheckbox = screen.getByRole('checkbox', { name: /local/i })
    expect(apiCheckbox).toBeInTheDocument()
    expect(localCheckbox).toBeInTheDocument()
  })

  it('calls onStartDateChange when start date is changed', () => {
    renderWithChakra(<ChartFilters {...mockProps} />)
    const startDateInput = screen.getAllByDisplayValue('2025-09-15T00:00')[0]
    fireEvent.change(startDateInput, { target: { value: '2025-09-14T00:00' } })
    expect(mockProps.onStartDateChange).toHaveBeenCalledWith('2025-09-14T00:00')
  })

  it('calls onEndDateChange when end date is changed', () => {
    renderWithChakra(<ChartFilters {...mockProps} />)
    const endDateInput = screen.getAllByDisplayValue('2025-09-15T23:59')[0]
    fireEvent.change(endDateInput, { target: { value: '2025-09-16T23:59' } })
    expect(mockProps.onEndDateChange).toHaveBeenCalledWith('2025-09-16T23:59')
  })

  it('calls onAggregationChange when aggregation is changed', () => {
    renderWithChakra(<ChartFilters {...mockProps} />)
    const select = screen.getByRole('combobox')
    fireEvent.change(select, { target: { value: 'day' } })
    expect(mockProps.onAggregationChange).toHaveBeenCalledWith('day')
  })

  it('calls onTypeToggle when API checkbox is clicked', () => {
    renderWithChakra(<ChartFilters {...mockProps} />)
    const apiCheckbox = screen.getByRole('checkbox', { name: /api/i })
    fireEvent.click(apiCheckbox)
    expect(mockProps.onTypeToggle).toHaveBeenCalledWith('api')
  })

  it('calls onTypeToggle when Local checkbox is clicked', () => {
    renderWithChakra(<ChartFilters {...mockProps} />)
    const localCheckbox = screen.getByRole('checkbox', { name: /local/i })
    fireEvent.click(localCheckbox)
    expect(mockProps.onTypeToggle).toHaveBeenCalledWith('local')
  })

  it('shows correct checkbox states', () => {
    renderWithChakra(<ChartFilters {...mockProps} />)
    const apiCheckbox = screen.getByRole('checkbox', { name: /api/i })
    const localCheckbox = screen.getByRole('checkbox', { name: /local/i })
    expect(apiCheckbox).toBeChecked()
    expect(localCheckbox).toBeChecked()
  })

  it('renders reset filters button', () => {
    renderWithChakra(<ChartFilters {...mockProps} />)
    const resetButton = screen.getByLabelText('Reset filters')
    expect(resetButton).toBeInTheDocument()
  })

  it('calls onResetFilters when reset button is clicked', () => {
    renderWithChakra(<ChartFilters {...mockProps} />)
    const resetButton = screen.getByLabelText('Reset filters')
    fireEvent.click(resetButton)
    expect(mockProps.onResetFilters).toHaveBeenCalled()
  })

  it('shows aggregation options', () => {
    renderWithChakra(<ChartFilters {...mockProps} />)
    expect(screen.getByText('----')).toBeInTheDocument()
    expect(screen.getByText('Day')).toBeInTheDocument()
  })

  it('handles unchecked types', () => {
    const propsWithNoTypes = { ...mockProps, types: [] }
    renderWithChakra(<ChartFilters {...propsWithNoTypes} />)
    const apiCheckbox = screen.getByRole('checkbox', { name: /api/i })
    const localCheckbox = screen.getByRole('checkbox', { name: /local/i })
    expect(apiCheckbox).not.toBeChecked()
    expect(localCheckbox).not.toBeChecked()
  })
})
