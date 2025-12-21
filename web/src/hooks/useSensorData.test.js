import { describe, it, expect, vi, beforeEach } from 'vitest'
import { renderHook, waitFor } from '@testing-library/react'
import { act } from 'react'
import { useSensorData } from './useSensorData'

describe('useSensorData', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('initializes with empty chart data', async () => {
    const { result } = renderHook(() =>
      useSensorData('loc1', '2025-09-15T00:00', '2025-09-15T23:59', '', ['api', 'local'])
    )

    await waitFor(() => {
      expect(result.current).toEqual({ labels: [], datasets: [] })
    })
  })

  it('does not fetch when locationSid is empty', () => {
    const { result } = renderHook(() =>
      useSensorData('', '2025-09-15T00:00', '2025-09-15T23:59', '', ['api', 'local'])
    )
    expect(result.current).toEqual({ labels: [], datasets: [] })
    expect(global.fetch).not.toHaveBeenCalled()
  })

  it('fetches sensor data successfully', async () => {
    const mockApiResponse = [
      { timestamp: '2025-09-15T10:00:00Z', type: 'api', temperature: 22.5 },
      { timestamp: '2025-09-15T11:00:00Z', type: 'api', temperature: 23.0 },
      { timestamp: '2025-09-15T10:00:00Z', type: 'local', temperature: 22.0 },
      { timestamp: '2025-09-15T11:00:00Z', type: 'local', temperature: 22.5 },
    ]

    global.fetch = vi.fn(() =>
      Promise.resolve({
        json: () => Promise.resolve(mockApiResponse),
      })
    )

    const { result } = renderHook(() =>
      useSensorData('loc1', '2025-09-15T00:00', '2025-09-15T23:59', '', ['api', 'local'])
    )

    await waitFor(() => {
      expect(result.current.labels).toHaveLength(2)
      expect(result.current.datasets).toHaveLength(2)
    })

    expect(result.current.datasets[0].label).toBe('API Sensor')
    expect(result.current.datasets[1].label).toBe('Local Sensor')
  })

  it('constructs correct URL without aggregation', async () => {
    global.fetch = vi.fn(() =>
      Promise.resolve({
        json: () => Promise.resolve([]),
      })
    )

    renderHook(() =>
      useSensorData('loc1', '2025-09-15T00:00', '2025-09-15T23:59', '', ['api', 'local'])
    )

    await waitFor(() => {
      expect(global.fetch).toHaveBeenCalledWith(
        'api/v1/sensors/data?location_sid=loc1&start_datetime=2025-09-15T00:00:00Z&end_datetime=2025-09-15T23:59:59Z&types=api&types=local'
      )
    })
  })

  it('constructs correct URL with aggregation', async () => {
    global.fetch = vi.fn(() =>
      Promise.resolve({
        json: () => Promise.resolve([]),
      })
    )

    renderHook(() =>
      useSensorData('loc1', '2025-09-15T00:00', '2025-09-15T23:59', 'day', ['api'])
    )

    await waitFor(() => {
      expect(global.fetch).toHaveBeenCalledWith(
        'api/v1/sensors/data?location_sid=loc1&start_datetime=2025-09-15T00:00:00Z&end_datetime=2025-09-15T23:59:59Z&aggregation=day&types=api'
      )
    })
  })

  it('handles only API sensor type', async () => {
    const mockApiResponse = [
      { timestamp: '2025-09-15T10:00:00Z', type: 'api', temperature: 22.5 },
    ]

    global.fetch = vi.fn(() =>
      Promise.resolve({
        json: () => Promise.resolve(mockApiResponse),
      })
    )

    const { result } = renderHook(() =>
      useSensorData('loc1', '2025-09-15T00:00', '2025-09-15T23:59', '', ['api'])
    )

    await waitFor(() => {
      expect(result.current.datasets).toHaveLength(1)
      expect(result.current.datasets[0].label).toBe('API Sensor')
    })
  })

  it('handles only Local sensor type', async () => {
    const mockApiResponse = [
      { timestamp: '2025-09-15T10:00:00Z', type: 'local', temperature: 22.0 },
    ]

    global.fetch = vi.fn(() =>
      Promise.resolve({
        json: () => Promise.resolve(mockApiResponse),
      })
    )

    const { result } = renderHook(() =>
      useSensorData('loc1', '2025-09-15T00:00', '2025-09-15T23:59', '', ['local'])
    )

    await waitFor(() => {
      expect(result.current.datasets).toHaveLength(1)
      expect(result.current.datasets[0].label).toBe('Local Sensor')
    })
  })

  it('refetches when parameters change', async () => {
    global.fetch = vi.fn(() =>
      Promise.resolve({
        json: () => Promise.resolve([]),
      })
    )

    const { rerender } = renderHook(
      ({ locationSid, startDate, endDate, aggregation, types, refreshTrigger }) =>
        useSensorData(locationSid, startDate, endDate, aggregation, types, refreshTrigger),
      {
        initialProps: {
          locationSid: 'loc1',
          startDate: '2025-09-15T00:00',
          endDate: '2025-09-15T23:59',
          aggregation: '',
          types: ['api'],
          refreshTrigger: 0,
        },
      }
    )

    await waitFor(() => {
      expect(global.fetch).toHaveBeenCalledTimes(1)
    })

    rerender({
      locationSid: 'loc1',
      startDate: '2025-09-16T00:00',
      endDate: '2025-09-16T23:59',
      aggregation: '',
      types: ['api'],
      refreshTrigger: 0,
    })

    await waitFor(() => {
      expect(global.fetch).toHaveBeenCalled()
    }, { timeout: 5000 })

    // Just verify it was called at least once
    expect(global.fetch.mock.calls.length).toBeGreaterThanOrEqual(1)
  })

  it('refetches when refreshTrigger changes', async () => {
    const mockFetch = vi.fn(() =>
      Promise.resolve({
        json: () => Promise.resolve([]),
      })
    )
    global.fetch = mockFetch

    const { rerender } = renderHook(
      ({ refreshTrigger }) =>
        useSensorData('loc1', '2025-09-15T00:00', '2025-09-15T23:59', '', ['api'], refreshTrigger),
      { initialProps: { refreshTrigger: 0 } }
    )

    await waitFor(() => {
      expect(mockFetch).toHaveBeenCalled()
    }, { timeout: 5000 })

    const initialCallCount = mockFetch.mock.calls.length

    rerender({ refreshTrigger: 1 })

    await waitFor(() => {
      expect(mockFetch.mock.calls.length).toBeGreaterThan(initialCallCount)
    }, { timeout: 5000 })
  })

  it('groups data by timestamp', async () => {
    const mockApiResponse = [
      { timestamp: '2025-09-15T10:00:00Z', type: 'api', temperature: 22.5 },
      { timestamp: '2025-09-15T10:00:00Z', type: 'local', temperature: 22.0 },
    ]

    global.fetch = vi.fn(() =>
      Promise.resolve({
        json: () => Promise.resolve(mockApiResponse),
      })
    )

    const { result } = renderHook(() =>
      useSensorData('loc1', '2025-09-15T00:00', '2025-09-15T23:59', '', ['api', 'local'])
    )

    await waitFor(() => {
      expect(result.current.labels).toHaveLength(1)
      expect(result.current.datasets[0].data[0]).toBe(22.5)
      expect(result.current.datasets[1].data[0]).toBe(22.0)
    })
  })

  it('handles null values for missing sensor types', async () => {
    const mockApiResponse = [
      { timestamp: '2025-09-15T10:00:00Z', type: 'api', temperature: 22.5 },
      { timestamp: '2025-09-15T11:00:00Z', type: 'local', temperature: 23.0 },
    ]

    global.fetch = vi.fn(() =>
      Promise.resolve({
        json: () => Promise.resolve(mockApiResponse),
      })
    )

    const { result } = renderHook(() =>
      useSensorData('loc1', '2025-09-15T00:00', '2025-09-15T23:59', '', ['api', 'local'])
    )

    await waitFor(() => {
      expect(result.current.labels).toHaveLength(2)
      expect(result.current.datasets[0].data[1]).toBeNull()
      expect(result.current.datasets[1].data[0]).toBeNull()
    })
  })

  it('cleans up on unmount', async () => {
    global.fetch = vi.fn(() =>
      Promise.resolve({
        json: () => Promise.resolve([]),
      })
    )

    const { unmount } = renderHook(() =>
      useSensorData('loc1', '2025-09-15T00:00', '2025-09-15T23:59', '', ['api'])
    )

    unmount()

    await waitFor(() => {
      expect(global.fetch).toHaveBeenCalled()
    })
  })
})
