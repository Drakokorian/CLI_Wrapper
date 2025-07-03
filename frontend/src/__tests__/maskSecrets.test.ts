import { describe, it, expect } from "vitest";
import { maskSecrets } from "../utils";

describe("maskSecrets", () => {
  it("masks 48-char OpenAI keys", () => {
    const key = "sk-" + "a".repeat(48);
    const result = maskSecrets(`prefix ${key} suffix`);
    expect(result).toBe(`prefix ${"*".repeat(key.length)} suffix`);
  });

  it("masks longer OpenAI keys", () => {
    const key = "sk-" + "b".repeat(60);
    const result = maskSecrets(`test ${key}`);
    expect(result).toBe(`test ${"*".repeat(key.length)}`);
  });
});
