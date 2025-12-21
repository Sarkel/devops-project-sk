import { describe, it, expect, vi, beforeEach } from 'vitest'
import { renderHook, waitFor } from '@testing-library/react'
import { act } from 'react'
import { useLocations } from './useLocations'

describe('useLocations', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('initializes with empty locations', async () => {
    const { result } = renderHook(() => useLocations())

    await waitFor(() => {
      expect(result.current.locations).toEqual([])
      expect(result.current.selectedLocation).toBe('')
    })
  })

  it('fetches locations successfully', async () => {
    const mockLocations = [
      { sid: 'loc1', name: 'Location 1' },
      { sid: 'loc2', name: 'Location 2' },
    ]

    global.fetch = vi.fn(() =>
      Promise.resolve({
        json: () => Promise.resolve(mockLocations),
      })
    )

    const { result } = renderHook(() => useLocations())

    await waitFor(() => {
      expect(result.current.locations).toEqual(mockLocations)
    })

    expect(result.current.selectedLocation).toBe('loc1')
    expect(global.fetch).toHaveBeenCalledWith('/api/v1/locations')
  })

  it('handles fetch error', async () => {
    const consoleErrorSpy = vi.spyOn(console, 'error').mockImplementation(() => {})

    global.fetch = vi.fn(() => Promise.reject(new Error('Network error')))

    const { result } = renderHook(() => useLocations())

    await waitFor(() => {
      expect(result.current.locations).toEqual([])
    })

    expect(consoleErrorSpy).toHaveBeenCalledWith(
      'Error fetching locations:',
      expect.any(Error)
    )

    consoleErrorSpy.mockRestore()
  })

  it('handles empty response', async () => {
    global.fetch = vi.fn(() =>
      Promise.resolve({
        json: () => Promise.resolve([]),
      })
    )

    const { result } = renderHook(() => useLocations())

    await waitFor(() => {
      expect(result.current.locations).toEqual([])
    })

    expect(result.current.selectedLocation).toBe('')
  })

  it('handles null response', async () => {
    global.fetch = vi.fn(() =>
      Promise.resolve({
        json: () => Promise.resolve(null),
      })
    )

    const { result } = renderHook(() => useLocations())

    await waitFor(() => {
      expect(result.current.locations).toEqual([])
    })
  })

  it('cleans up on unmount', async () => {
    const mockLocations = [{ sid: 'loc1', name: 'Location 1' }]

    global.fetch = vi.fn(() =>
      Promise.resolve({
        json: () => Promise.resolve(mockLocations),
      })
    )

    const { unmount } = renderHook(() => useLocations())
    unmount()

    // Should not cause any errors after unmount
    await waitFor(() => {
      expect(global.fetch).toHaveBeenCalled()
    })
  })

  it('allows setting selected location', async () => {
    const mockLocations = [
      { sid: 'loc1', name: 'Location 1' },
      { sid: 'loc2', name: 'Location 2' },
    ]

    global.fetch = vi.fn(() =>
      Promise.resolve({
        json: () => Promise.resolve(mockLocations),
      })
    )

    const { result } = renderHook(() => useLocations())

    await waitFor(() => {
      expect(result.current.selectedLocation).toBe('loc1')
    })

    await act(async () => {
      result.current.setSelectedLocation('loc2')
    })

    await waitFor(() => {
      expect(result.current.selectedLocation).toBe('loc2')
    })
  })
})
