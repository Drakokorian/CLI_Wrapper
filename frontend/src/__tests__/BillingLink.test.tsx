import { render, screen, waitFor } from "@testing-library/react";
import { vi, beforeEach, it, describe, expect } from "vitest";
import BillingLink from "../components/BillingLink";

vi.stubGlobal("fetch", vi.fn());
const openURL = vi.fn();
vi.mock("@wailsio/runtime", () => ({
  OpenURL: (...args: any[]) => openURL(...args),
}));

beforeEach(() => {
  (fetch as any).mockReset();
  (fetch as any).mockResolvedValue({ ok: true, json: async () => ({}) });
  openURL.mockReset();
});

describe("BillingLink", () => {
  it("fetches info and opens url", async () => {
    (fetch as any).mockResolvedValueOnce({
      ok: true,
      json: async () => ({
        tool: "openai",
        url: "https://billing",
        usage: { total_granted: 10, total_used: 3, total_available: 7 },
      }),
    });

    render(<BillingLink />);
    await waitFor(() => screen.getByText("View Billing"));
    expect(screen.getByLabelText("usage")).toHaveTextContent("Used 3.00 / 10.00");
    screen.getByText("View Billing").click();
    expect(openURL).toHaveBeenCalledWith("https://billing");
  });
});
