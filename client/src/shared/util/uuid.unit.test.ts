import { describe, it, expect } from "vitest";
import { uuidv7 } from "./uuid";

const UUID_RE = /^[0-9a-f]{8}-[0-9a-f]{4}-7[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/;

describe("uuidv7", () => {
  it("matches the canonical UUID v7 layout (8-4-4-4-12) with version 7 and RFC 4122 variant", () => {
    for (let i = 0; i < 16; i++) {
      expect(uuidv7()).toMatch(UUID_RE);
    }
  });

  it("embeds the millisecond timestamp in the first 40 bits", () => {
    const ts = Date.now();
    const id = uuidv7();
    const msHex = id.slice(0, 13).replace("-", "");
    const msFromId = Number.parseInt(msHex, 16);
    expect(msFromId).toBeGreaterThanOrEqual(ts - 5);
    expect(msFromId).toBeLessThanOrEqual(ts + 5);
  });

  it("produces monotonically non-decreasing timestamps", () => {
    let prev = 0;
    for (let i = 0; i < 32; i++) {
      const msHex = uuidv7().slice(0, 13).replace("-", "");
      const ms = Number.parseInt(msHex, 16);
      expect(ms).toBeGreaterThanOrEqual(prev);
      prev = ms;
    }
  });

  it("returns distinct values for successive calls", () => {
    const seen = new Set<string>();
    for (let i = 0; i < 256; i++) seen.add(uuidv7());
    expect(seen.size).toBe(256);
  });

  it("has version nibble 7 at position 14 and variant in 8/9/a/b at position 19", () => {
    const id = uuidv7();
    expect(id[14]).toBe("7");
    expect(["8", "9", "a", "b"]).toContain(id[19]);
  });
});
