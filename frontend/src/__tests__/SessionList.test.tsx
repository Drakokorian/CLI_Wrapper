import { render, screen, waitFor } from '@testing-library/react';
import SessionList from '../components/SessionList';
import { vi } from 'vitest';

vi.stubGlobal('fetch', vi.fn());

describe('SessionList', () => {
  beforeEach(() => {
    ;(fetch as any).mockReset();
  });

  it('displays fetched session IDs', async () => {
    (fetch as any).mockResolvedValueOnce({
      ok: true,
      json: async () => ['abc123', 'xyz789'],
    });

    render(<SessionList />);

    await waitFor(() => {
      expect(screen.getByText('abc123')).toBeInTheDocument();
      expect(screen.getByText('xyz789')).toBeInTheDocument();
    });
  });
});
