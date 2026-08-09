import { describe, it, expect } from "vitest";
import { SLIDE_SEPARATOR, splitSlides, buildPresentHtml, buildPreviewHtml } from "./reveal-html";
import { CDN_REVEAL_CSS, CDN_REVEAL_THEME_CSS, CDN_REVEAL_JS } from "../cfg";

describe("SLIDE_SEPARATOR", () => {
  it("matches a line of dashes", () => {
    expect(SLIDE_SEPARATOR.test("---")).toBe(true);
    expect(SLIDE_SEPARATOR.test("----")).toBe(true);
  });

  it("does not match plain text", () => {
    expect(SLIDE_SEPARATOR.test("hello")).toBe(false);
  });
});

describe("splitSlides", () => {
  it("returns a single slide when there is no separator", () => {
    expect(splitSlides("Hello world")).toEqual(["Hello world"]);
  });

  it("splits on --- ; trailing newline of the previous slide is preserved, leading newline of next slide is stripped", () => {
    expect(splitSlides("First\n---\nSecond")).toEqual(["First\n", "Second"]);
  });

  it("preserves intra-slide newlines but strips the leading newline of each piece", () => {
    expect(splitSlides("First\n\n---\n\nSecond")).toEqual(["First\n\n", "Second"]);
  });

  it("handles a trailing separator by producing an empty trailing slide", () => {
    expect(splitSlides("First\n---\n")).toEqual(["First\n", ""]);
  });

  it("handles multiple consecutive separators", () => {
    expect(splitSlides("---\n---\n---")).toEqual(["", "", "", ""]);
  });

  it("does not split on dashes mid-line (separator must be the whole line)", () => {
    expect(splitSlides("a-b-c\n---\nd")).toEqual(["a-b-c\n", "d"]);
  });
});

describe("buildPresentHtml", () => {
  it("wraps each slide in a <section> inside the reveal container", () => {
    const html = buildPresentHtml(["<h1>A</h1>", "<h1>B</h1>"], 1);
    expect(html).toContain("<section><h1>A</h1></section>");
    expect(html).toContain("<section><h1>B</h1></section>");
  });

  it("references the reveal CDN assets", () => {
    const html = buildPresentHtml(["x"], 0);
    expect(html).toContain(`href="${CDN_REVEAL_CSS}"`);
    expect(html).toContain(`href="${CDN_REVEAL_THEME_CSS}"`);
    expect(html).toContain(`src="${CDN_REVEAL_JS}"`);
  });

  it("seeds Reveal.slide with the initialSlide argument", () => {
    expect(buildPresentHtml(["a", "b"], 7)).toContain("Reveal.slide(7)");
  });

  it("registers a navigate postMessage handler", () => {
    expect(buildPresentHtml(["a"], 0)).toContain("e.data.type === 'navigate'");
  });
});

describe("buildPreviewHtml", () => {
  it("wraps each slide in a <section>", () => {
    const html = buildPreviewHtml(["<p>1</p>", "<p>2</p>"], 0);
    expect(html).toContain("<section><p>1</p></section>");
    expect(html).toContain("<section><p>2</p></section>");
  });

  it("references the reveal CDN assets", () => {
    const html = buildPreviewHtml(["a"], 0);
    expect(html).toContain(`href="${CDN_REVEAL_CSS}"`);
    expect(html).toContain(`src="${CDN_REVEAL_JS}"`);
  });

  it("seeds Reveal.slide with the targetIndex argument", () => {
    expect(buildPreviewHtml(["a", "b"], 3)).toContain("Reveal.slide(3)");
  });
});
