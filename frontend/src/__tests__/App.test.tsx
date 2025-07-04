import { render, screen, waitFor } from '@testing-library/react';
import { vi, beforeAll, beforeEach, describe, it, expect } from 'vitest';
import { fireEvent } from '@testing-library/react';
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
  (global as any).window.matchMedia = vi.fn().mockReturnValue({
    matches: false,
    addListener: vi.fn(),
    removeListener: vi.fn(),
  });
  class WS {
    static instances: WS[] = [];
    onmessage: ((e: MessageEvent) => void) | null = null;
    constructor(url: string) {
      WS.instances.push(this);
    }
    send() {}
    close() {}
  }
  vi.stubGlobal('WebSocket', WS as any);
});

beforeEach(() => {
  (fetch as any).mockReset();
  (fetch as any).mockImplementation(async (url: string) => {
    if (url.includes('sessions')) return { ok: true, json: async () => [] };
    if (url.includes('models')) return { ok: true, json: async () => [] };
    if (url.includes('resource')) return { ok: true, json: async () => ({ cpu: 0, memory: 0 }) };
    if (url.includes('billing')) return { ok: true, json: async () => ({ tool: '', url: '', usage: { total_granted: 0, total_used: 0, total_available: 0 } }) };
    if (url.includes('theme')) return { ok: true, json: async () => ({ theme: 'light' }) };
    if (url.includes('config')) return { ok: true, json: async () => ({ concurrency: 1, workingDir: '', logLevel: 'info', logPath: '/tmp/log' }) };
    if (url.includes('history')) return { ok: true, json: async () => [] };
    return { ok: true, json: async () => ({}) };
  });
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

  it('animates streamed characters', async () => {
    vi.useFakeTimers();
    render(<App />);
    const input = screen.getAllByRole('textbox')[1];
    fireEvent.change(input, { target: { value: 'hi' } });
    (window.backend.RunPrompt as any).mockResolvedValueOnce('id1');
    fireEvent.click(screen.getByText('Send'));
    await waitFor(() => (global as any).WebSocket.instances.length > 0);
    const ws = (global as any).WebSocket.instances[0] as any;
    ws.onmessage?.({ data: 'h' });
    ws.onmessage?.({ data: 'i' });
    ws.onmessage?.({ data: '\n' });
    vi.advanceTimersByTime(100);
    await waitFor(() => {
      expect(screen.getByText(/hi/)).toBeInTheDocument();
    });
    vi.useRealTimers();
  });

  it('posts settings including log config', async () => {
    render(<App />);
    await waitFor(() => screen.getByText('Settings'));
    fireEvent.click(screen.getByText('Settings'));
    fireEvent.change(screen.getByLabelText('Log Level'), { target: { value: 'debug' } });
    fireEvent.change(screen.getByLabelText('Log Path'), { target: { value: '/tmp/out.log' } });
    (fetch as any).mockResolvedValueOnce({ ok: true, json: async () => ({}) });
    fireEvent.click(screen.getByText('Save'));
    await waitFor(() => expect((fetch as any).mock.calls.length).toBeGreaterThan(1));
    const call = (fetch as any).mock.calls.pop();
    expect(call[0]).toContain('/config');
    const opts = call[1];
    const body = JSON.parse(opts.body);
    expect(body).toEqual({ concurrency: 1, workingDir: '', logLevel: 'debug', logPath: '/tmp/out.log' });
  });
});
