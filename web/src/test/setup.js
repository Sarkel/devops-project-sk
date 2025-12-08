import { afterEach, vi, beforeAll } from 'vitest'
import { cleanup } from '@testing-library/react'
import '@testing-library/jest-dom/vitest'

// Mock HTMLCanvasElement
beforeAll(() => {
  HTMLCanvasElement.prototype.getContext = vi.fn(() => ({
    fillRect: vi.fn(),
    clearRect: vi.fn(),
    getImageData: vi.fn(),
    putImageData: vi.fn(),
    createImageData: vi.fn(),
    setTransform: vi.fn(),
    drawImage: vi.fn(),
    save: vi.fn(),
    fillText: vi.fn(),
    restore: vi.fn(),
    beginPath: vi.fn(),
    moveTo: vi.fn(),
    lineTo: vi.fn(),
    closePath: vi.fn(),
    stroke: vi.fn(),
    translate: vi.fn(),
    scale: vi.fn(),
    rotate: vi.fn(),
    arc: vi.fn(),
    fill: vi.fn(),
    measureText: vi.fn(() => ({ width: 0 })),
    transform: vi.fn(),
    rect: vi.fn(),
    clip: vi.fn(),
    canvas: {
      width: 500,
      height: 500,
    },
  }))

  HTMLCanvasElement.prototype.toDataURL = vi.fn(() => 'data:image/png;base64,')

  // Suppress chart.js console errors
  const originalError = console.error
  console.error = (...args) => {
    if (
      typeof args[0] === 'string' &&
      (args[0].includes('Failed to create chart') ||
       args[0].includes("can't acquire context"))
    ) {
      return
    }
    originalError.apply(console, args)
  }
})

// Mock fetch globally
global.fetch = vi.fn(() =>
  Promise.resolve({
    ok: true,
    json: () => Promise.resolve([]),
  })
)

// Cleanup after each test case
afterEach(() => {
  cleanup()
  vi.clearAllMocks()
})
