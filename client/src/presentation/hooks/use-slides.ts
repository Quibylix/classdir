import { useState, useEffect, useCallback } from "react";
import { DEFAULT_SLIDE_CONTENT } from "../cfg";

export function useSlides(initialContent: string) {
  const [slides, setSlides] = useState<string[]>(splitContent(initialContent));

  useEffect(() => {
    setSlides(splitContent(initialContent));
  }, [initialContent]);

  const addSlide = useCallback(() => {
    setSlides((prev) => [...prev, DEFAULT_SLIDE_CONTENT]);
  }, []);

  const removeSlide = useCallback((index: number) => {
    setSlides((prev) => prev.filter((_, i) => i !== index));
  }, []);

  const updateSlideContent = useCallback((index: number, content: string) => {
    setSlides((prev) => prev.map((s, i) => (i === index ? content : s)));
  }, []);

  const joinSlides = useCallback(() => {
    return slides.map((slide) => slide.trim()).join("\n---\n");
  }, [slides]);

  return {
    slides,
    addSlide,
    removeSlide,
    updateSlideContent,
    joinSlides,
  };
}

function splitContent(content: string): string[] {
  if (!content) return [""];
  return content.split(/^---+\s*\n/m).map((p) => p.replace(/^\n+/, ""));
}
