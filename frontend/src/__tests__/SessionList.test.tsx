import { render, screen, waitFor } from "@testing-library/react";
import SessionList from "../components/SessionList";
import { vi } from "vitest";

vi.stubGlobal("fetch", vi.fn());
class WS {
  static instances: WS[] = [];
  onmessage: ((e: MessageEvent) => void) | null = null;
  constructor(url: string) {
    WS.instances.push(this);
  }
  close() {}
}
vi.stubGlobal("WebSocket", WS as any);

describe("SessionList", () => {
  beforeEach(() => {
    (fetch as any).mockReset();
    (fetch as any).mockResolvedValue({ ok: true, json: async () => [] });
    (WebSocket as any).instances.length = 0;
  });

  it("displays fetched session IDs", async () => {
    (fetch as any).mockResolvedValueOnce({
      ok: true,
      json: async () => ["abc123", "xyz789"],
    });

    render(<SessionList />);

    await waitFor(() => {
      expect(screen.getByText("abc123")).toBeInTheDocument();
      expect(screen.getByText("xyz789")).toBeInTheDocument();
    });
  });

  it("appends output from websocket", async () => {
    vi.useFakeTimers();
    (fetch as any).mockResolvedValue({
      ok: true,
      json: async () => ["id1"],
    });
    render(<SessionList />);
    vi.advanceTimersByTime(0);
    await waitFor(() => {
      expect(screen.getByText("id1")).toBeInTheDocument();
    });
    await waitFor(() => (WebSocket as any).instances.length > 0);
    const ws = (WebSocket as any).instances[0] as any;
    ws.onmessage?.({ data: "hello" });
    await waitFor(() => {
      expect(screen.getByText("hello")).toBeInTheDocument();
    });
    vi.useRealTimers();
  });
});
