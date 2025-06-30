import { render, screen, waitFor } from "@testing-library/react";
import { vi, describe, it, beforeEach, expect } from "vitest";
import ModelSelector from "../components/ModelSelector";

vi.stubGlobal("fetch", vi.fn());

describe("ModelSelector", () => {
  beforeEach(() => {
    (fetch as any).mockReset();
    (fetch as any).mockResolvedValue({ ok: true, json: async () => [] });
  });

  it("highlights on selection change", async () => {
    (fetch as any).mockResolvedValueOnce({ ok: true, json: async () => ["a", "b"] });
    const { rerender } = render(<ModelSelector selected="a" onChange={() => {}} />);
    await waitFor(() => screen.getByRole("combobox"));
    rerender(<ModelSelector selected="b" onChange={() => {}} />);
    const select = screen.getByRole("combobox");
    expect(select).toHaveStyle({
      borderColor: "rgb(37, 99, 235)",
      backgroundColor: "rgb(239, 246, 255)",
    });
    await new Promise((r) => setTimeout(r, 300));
    expect(select).not.toHaveStyle({ borderColor: "rgb(37, 99, 235)" });
  });
});
