import { describe, it, expect, vi, beforeEach } from "vitest";
import { isApiError, toApiError } from "./errors";
import { ERR_CODE_UNKNOWN } from "../cfg/http";

describe("isApiError", () => {
  it("returns true for a well-formed ApiError", () => {
    expect(isApiError({ code: "E", message: "m", status: 400 })).toBe(true);
  });

  it("returns false for a missing field", () => {
    expect(isApiError({ code: "E", message: "m" })).toBe(false);
  });

  it("returns false for wrong-typed status", () => {
    expect(isApiError({ code: "E", message: "m", status: "400" })).toBe(false);
  });

  it("returns false for arbitrary objects", () => {
    expect(isApiError(null)).toBe(false);
    expect(isApiError(undefined)).toBe(false);
    expect(isApiError({})).toBe(false);
    expect(isApiError(42)).toBe(false);
  });
});

describe("toApiError", () => {
  beforeEach(() => {
    vi.useFakeTimers();
  });

  function fakeResponse(body: unknown, status: number, statusText = ""): Response {
    return {
      status,
      statusText,
      json: async () => body,
    } as unknown as Response;
  }

  it("maps {error:{code,message}} to an ApiError with status attached", async () => {
    const res = fakeResponse({ error: { code: "BAD", message: "nope" } }, 400, "Bad Request");
    expect(await toApiError(res)).toEqual({
      code: "BAD",
      message: "nope",
      status: 400,
    });
  });

  it("falls back to UNKNOWN + statusText when body isn't the error envelope", async () => {
    const res = fakeResponse({ unrelated: "field" }, 404, "Not Found");
    expect(await toApiError(res)).toEqual({
      code: ERR_CODE_UNKNOWN,
      message: "Not Found",
      status: 404,
    });
  });

  it("falls back when json() rejects (malformed body)", async () => {
    const res = {
      status: 500,
      statusText: "Server Error",
      json: async () => {
        throw new SyntaxError("bad json");
      },
    } as unknown as Response;
    expect(await toApiError(res)).toEqual({
      code: ERR_CODE_UNKNOWN,
      message: "Server Error",
      status: 500,
    });
  });
});
