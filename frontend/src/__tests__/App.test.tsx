import { render, screen, waitFor } from '@testing-library/react';
import { vi } from 'vitest';
import App from '../App';

const events: Record<string, (data: any) => void> = {};
vi.mock('@wailsio/runtime', () => ({
  EventsOn: (event: string, handler: (data: any) => void) => {
    events[event] = handler;
    return () => delete events[event];
  },
}));

beforeAll(() => {
  vi.stubGlobal('fetch', vi.fn());
  (global as any).window.backend = { RunPrompt: vi.fn() };
});

beforeEach(() => {
  (fetch as any).mockReset();
  (window.backend.RunPrompt as any).mockReset();
});

describe('App logging', () => {
  it('shows alert on model switch event', async () => {
    render(<App />);

    events['model:switched']?.({ from: 'old', to: 'new' });

    await waitFor(() => {
      expect(screen.getByRole('alert')).toHaveTextContent('Model switched from old to new');
    });
  });
});
