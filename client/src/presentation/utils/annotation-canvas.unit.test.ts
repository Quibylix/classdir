import { describe, it, expect } from "vitest";
import { visibleStrokes, toPercent } from "./annotation-canvas";
import type { AnnotationOperation } from "../types";

function stroke(id: string, color = "#ff0000"): AnnotationOperation {
  return {
    type: "stroke",
    id,
    payload: { points: [{ x: 1, y: 2 }], color, thickness: 3 },
  };
}

describe("visibleStrokes", () => {
  it("returns an empty list for no operations", () => {
    expect(visibleStrokes([])).toEqual([]);
  });

  it("accumulates stroke payloads in order", () => {
    expect(visibleStrokes([stroke("a"), stroke("b")])).toEqual([
      { points: [{ x: 1, y: 2 }], color: "#ff0000", thickness: 3 },
      { points: [{ x: 1, y: 2 }], color: "#ff0000", thickness: 3 },
    ]);
  });

  it("clears accumulated strokes on a clear op", () => {
    expect(visibleStrokes([stroke("a"), { type: "clear", id: "c" }, stroke("b")])).toEqual([
      { points: [{ x: 1, y: 2 }], color: "#ff0000", thickness: 3 },
    ]);
  });

  it("the last clear wipes everything; only strokes after the last clear survive", () => {
    const ops: AnnotationOperation[] = [
      stroke("s1", "#00ff00"),
      { type: "clear", id: "c1" },
      stroke("s2", "#0000ff"),
      stroke("s3", "#ffff00"),
      { type: "clear", id: "c2" },
      stroke("s4", "#ff00ff"),
    ];
    expect(visibleStrokes(ops)).toEqual([
      { points: [{ x: 1, y: 2 }], color: "#ff00ff", thickness: 3 },
    ]);
  });

  it("ignores unknown operation types", () => {
    const ops = [
      stroke("a"),
      { type: "bogus" as const, id: "x" },
      stroke("b"),
    ] as unknown as AnnotationOperation[];
    expect(visibleStrokes(ops)).toHaveLength(2);
  });
});

describe("toPercent", () => {
  it("converts client coords to percentages relative to the rect", () => {
    const rect = { left: 10, top: 20, width: 100, height: 50 } as DOMRect;
    expect(toPercent(60, 45, rect)).toEqual({ x: 50, y: 50 });
  });

  it("handles leading edge (top-left corner)", () => {
    const rect = { left: 0, top: 0, width: 200, height: 100 } as DOMRect;
    expect(toPercent(0, 0, rect)).toEqual({ x: 0, y: 0 });
  });

  it("handles trailing edge (bottom-right corner)", () => {
    const rect = { left: 0, top: 0, width: 200, height: 100 } as DOMRect;
    expect(toPercent(200, 100, rect)).toEqual({ x: 100, y: 100 });
  });
});
