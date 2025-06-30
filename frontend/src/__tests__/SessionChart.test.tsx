import { render, screen, waitFor } from "@testing-library/react";
import { vi } from "vitest";
import SessionChart from "../components/SessionChart";

vi.stubGlobal("fetch", vi.fn());

beforeEach(() => {
  (fetch as any).mockReset();
  (fetch as any).mockResolvedValue({ ok: true, json: async () => [] });
});

it("renders metric lines", async () => {
  (fetch as any).mockResolvedValueOnce({
    ok: true,
    json: async () => [{ id: "1", timestamp: "t", cpu: 10, memory: 20 }],
  });
  render(<SessionChart id="1" />);
  await waitFor(() => {
    const img = screen.getByRole("img");
    expect(img).toBeInTheDocument();
    const lines = img.querySelectorAll('polyline');
    expect(lines[0]).toHaveAttribute('stroke', 'rgb(29, 78, 216)');
    expect(lines[1]).toHaveAttribute('stroke', 'rgb(4, 120, 87)');
  });
});
