import { describe, it, expect } from "vitest";
import { maskSecrets } from "../utils";

describe("maskSecrets", () => {
  const lengths = [48, 50, 60];
  it.each(lengths)("masks OpenAI keys of length %i", (len) => {
    const key = "sk-" + "a".repeat(len);
    const result = maskSecrets(`prefix ${key} suffix`);
    expect(result).toBe(`prefix ${"*".repeat(key.length)} suffix`);
  });
});
