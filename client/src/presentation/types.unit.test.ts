import { describe, it, expect } from "vitest";
import {
  PresentationSchema,
  PresentationPreviewSchema,
  AnnotationOperationSchema,
  WSOutputMessageSchema,
} from "./types";

describe("PresentationSchema", () => {
  it("parses a full presentation shape", () => {
    expect(PresentationSchema.safeParse({ id: "x", title: "T", content: "c" }).success).toBe(true);
  });

  it("rejects missing content", () => {
    expect(PresentationSchema.safeParse({ id: "x", title: "T" }).success).toBe(false);
  });
});

describe("PresentationPreviewSchema", () => {
  it("parses id + title only and ignores extras", () => {
    expect(PresentationPreviewSchema.safeParse({ id: "x", title: "T", extra: 1 }).success).toBe(
      true,
    );
  });
});

describe("AnnotationOperationSchema", () => {
  it("accepts a clear op", () => {
    expect(AnnotationOperationSchema.safeParse({ type: "clear", id: "c" }).success).toBe(true);
  });

  it("accepts a stroke op with payload", () => {
    expect(
      AnnotationOperationSchema.safeParse({
        type: "stroke",
        id: "s",
        payload: {
          points: [{ x: 0, y: 1 }],
          color: "#fff",
          thickness: 2,
        },
      }).success,
    ).toBe(true);
  });

  it("rejects a stroke op without payload", () => {
    expect(AnnotationOperationSchema.safeParse({ type: "stroke", id: "s" }).success).toBe(false);
  });

  it("rejects an unknown type", () => {
    expect(AnnotationOperationSchema.safeParse({ type: "noise", id: "n" }).success).toBe(false);
  });
});

describe("WSOutputMessageSchema", () => {
  it("accepts slide_changed", () => {
    expect(
      WSOutputMessageSchema.safeParse({
        event: "slide_changed",
        data: { current_slide: 3 },
      }).success,
    ).toBe(true);
  });

  it("accepts annotation_added", () => {
    expect(
      WSOutputMessageSchema.safeParse({
        event: "annotation_added",
        data: { type: "clear", id: "c" },
      }).success,
    ).toBe(true);
  });

  it("accepts annotations_batch with per-slide operations", () => {
    expect(
      WSOutputMessageSchema.safeParse({
        event: "annotations_batch",
        data: {
          operations_by_slide: {
            "0": [{ type: "clear", id: "c1" }],
            "1": [
              {
                type: "stroke",
                id: "s1",
                payload: {
                  points: [{ x: 0, y: 0 }],
                  color: "#fff",
                  thickness: 1,
                },
              },
            ],
          },
        },
      }).success,
    ).toBe(true);
  });

  it("accepts the slides payload event", () => {
    expect(
      WSOutputMessageSchema.safeParse({
        data: { slides: ["a", "b"], current_index: 0, room_code: "XYZ" },
      }).success,
    ).toBe(true);
  });

  it("accepts the slides payload event without room_code", () => {
    expect(
      WSOutputMessageSchema.safeParse({
        data: { slides: ["a"], current_index: 0 },
      }).success,
    ).toBe(true);
  });

  it("accepts an error envelope", () => {
    expect(
      WSOutputMessageSchema.safeParse({
        error: { code: "E", message: "m" },
      }).success,
    ).toBe(true);
  });

  it("rejects a slide_changed with string current_slide", () => {
    expect(
      WSOutputMessageSchema.safeParse({
        event: "slide_changed",
        data: { current_slide: "x" },
      }).success,
    ).toBe(false);
  });

  it("rejects a message with no recognised shape", () => {
    expect(WSOutputMessageSchema.safeParse({ foo: "bar" }).success).toBe(false);
  });
});
