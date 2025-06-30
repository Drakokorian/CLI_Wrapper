import { render, screen, waitFor } from "@testing-library/react";
import { vi, describe, it, beforeEach, expect } from "vitest";
import ModelSelector from "../components/ModelSelector";

vi.stubGlobal("fetch", vi.fn());

describe("ModelSelector", () => {
  beforeEach(() => {
    (fetch as any).mockReset();
  });

  it("highlights on selection change", async () => {
    vi.useFakeTimers();
    (fetch as any).mockResolvedValueOnce({ ok: true, json: async () => ["a", "b"] });
    const { rerender } = render(<ModelSelector selected="a" onChange={() => {}} />);
    await waitFor(() => screen.getByRole("combobox"));
    rerender(<ModelSelector selected="b" onChange={() => {}} />);
    const select = screen.getByRole("combobox");
    expect(select.className).toMatch(/border-blue-500/);
    vi.advanceTimersByTime(300);
    expect(select.className).not.toMatch(/border-blue-500/);
    vi.useRealTimers();
  });
});
