// noinspection ES6UnusedImports (Required to waitFor async events in tests)
import MutationObserver from "mutationobserver-shim";
import '@testing-library/jest-dom';
import { vi, expect } from 'vitest';
import axios from 'axios';

// Make vi available globally (like jest)
global.vi = vi;

// Mock axios
vi.mock('axios', () => {
  return {
    default: {
      get: vi.fn(),
      post: vi.fn(),
      put: vi.fn(),
      delete: vi.fn(),
    },
  };
});



// Mock window.location
Object.defineProperty(window, 'location', {
    value: {
        origin: "https://mock-your-api-calls"
    },
    writable: true,
});

// Setup global mocks
global.ResizeObserver = vi.fn().mockImplementation(() => ({
    observe: vi.fn(),
    unobserve: vi.fn(),
    disconnect: vi.fn(),
}));

// Mock matchMedia
Object.defineProperty(window, 'matchMedia', {
    writable: true,
    value: vi.fn().mockImplementation(query => ({
        matches: false,
        media: query,
        onchange: null,
        addListener: vi.fn(),
        removeListener: vi.fn(),
        addEventListener: vi.fn(),
        removeEventListener: vi.fn(),
        dispatchEvent: vi.fn(),
    })),
});
