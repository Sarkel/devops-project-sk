import { describe, it, expect, vi, beforeEach } from 'vitest'
import { renderHook, waitFor } from '@testing-library/react'
import { act } from 'react'
import { useSensorSummary } from './useSensorSummary'

describe('useSensorSummary', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('initializes with null summary', async () => {
    const { result } = renderHook(() => useSensorSummary('loc1'))

    await waitFor(() => {
      expect(result.current).toBeNull()
    })
  })

  it('does not fetch when locationSid is empty', () => {
    const { result } = renderHook(() => useSensorSummary(''))
    expect(result.current).toBeNull()
    expect(global.fetch).not.toHaveBeenCalled()
  })

  it('fetches summary successfully', async () => {
    const mockSummary = {
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
    }

    global.fetch = vi.fn(() =>
      Promise.resolve({
        json: () => Promise.resolve(mockSummary),
      })
    )

    const { result } = renderHook(() => useSensorSummary('loc1'))

    await waitFor(() => {
      expect(result.current).toEqual(mockSummary)
    })

    expect(global.fetch).toHaveBeenCalledWith('api/v1/sensors/summary?location_sid=loc1')
  })

  it('handles fetch error', async () => {
    const consoleErrorSpy = vi.spyOn(console, 'error').mockImplementation(() => {})

    global.fetch = vi.fn(() => Promise.reject(new Error('Network error')))

    const { result } = renderHook(() => useSensorSummary('loc1'))

    await waitFor(() => {
      expect(result.current).toBeNull()
    })

    expect(consoleErrorSpy).toHaveBeenCalledWith(
      'Error fetching summary:',
      expect.any(Error)
    )

    consoleErrorSpy.mockRestore()
  })

  it('refetches when locationSid changes', async () => {
    const mockSummary1 = {
      api: { temperature: 23.5, trend: 1.2, timestamp: '2025-09-15T12:00:00Z' },
      local: { temperature: 22.8, trend: -0.5, timestamp: '2025-09-15T12:00:00Z' },
    }

    const mockSummary2 = {
      api: { temperature: 24.0, trend: 0.5, timestamp: '2025-09-15T13:00:00Z' },
      local: { temperature: 23.2, trend: 0.3, timestamp: '2025-09-15T13:00:00Z' },
    }

    global.fetch = vi
      .fn()
      .mockResolvedValueOnce({ json: () => Promise.resolve(mockSummary1) })
      .mockResolvedValueOnce({ json: () => Promise.resolve(mockSummary2) })

    const { result, rerender } = renderHook(
      ({ locationSid }) => useSensorSummary(locationSid),
      { initialProps: { locationSid: 'loc1' } }
    )

    await waitFor(() => {
      expect(result.current).toEqual(mockSummary1)
    })

    rerender({ locationSid: 'loc2' })

    await waitFor(() => {
      expect(result.current).toEqual(mockSummary2)
    })

    expect(global.fetch).toHaveBeenCalledTimes(2)
  })

  it('refetches when refreshTrigger changes', async () => {
    const mockSummary1 = {
      api: { temperature: 23.5, trend: 1.2, timestamp: '2025-09-15T12:00:00Z' },
      local: { temperature: 22.8, trend: -0.5, timestamp: '2025-09-15T12:00:00Z' },
    }

    const mockSummary2 = {
      api: { temperature: 24.0, trend: 0.5, timestamp: '2025-09-15T13:00:00Z' },
      local: { temperature: 23.2, trend: 0.3, timestamp: '2025-09-15T13:00:00Z' },
    }

    global.fetch = vi
      .fn()
      .mockResolvedValueOnce({ json: () => Promise.resolve(mockSummary1) })
      .mockResolvedValueOnce({ json: () => Promise.resolve(mockSummary2) })

    const { result, rerender } = renderHook(
      ({ refreshTrigger }) => useSensorSummary('loc1', refreshTrigger),
      { initialProps: { refreshTrigger: 0 } }
    )

    await waitFor(() => {
      expect(result.current).toEqual(mockSummary1)
    })

    rerender({ refreshTrigger: 1 })

    await waitFor(() => {
      expect(result.current).toEqual(mockSummary2)
    })

    expect(global.fetch).toHaveBeenCalledTimes(2)
  })

  it('cleans up on unmount', async () => {
    global.fetch = vi.fn(() =>
      Promise.resolve({
        json: () => Promise.resolve({}),
      })
    )

    const { unmount } = renderHook(() => useSensorSummary('loc1'))
    unmount()

    // Should not cause any errors after unmount
    await waitFor(() => {
      expect(global.fetch).toHaveBeenCalled()
    })
  })
})
