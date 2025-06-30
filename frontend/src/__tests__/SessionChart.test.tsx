import { render, screen, waitFor } from "@testing-library/react";
import { vi } from "vitest";
import SessionChart from "../components/SessionChart";

vi.stubGlobal("fetch", vi.fn());

it("renders metric lines", async () => {
  (fetch as any).mockResolvedValueOnce({
    ok: true,
    json: async () => [{ id: "1", timestamp: "t", cpu: 10, memory: 20 }],
  });
  render(<SessionChart id="1" />);
  await waitFor(() => {
    expect(screen.getByRole("img")).toBeInTheDocument();
  });
});
