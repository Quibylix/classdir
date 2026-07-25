import { useState, useEffect } from "react";
import { getPresentation, updatePresentation } from "../api";
import type { Presentation } from "../types";
import type { FetchError } from "../../shared/api/fetch";

export function usePresentation(id: string) {
  const [presentation, setPresentation] = useState<Presentation | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isSaving, setIsSaving] = useState(false);
  const [error, setError] = useState<FetchError | null>(null);

  useEffect(() => {
    setIsLoading(true);
    setError(null);

    getPresentation(id)
      .match(
        (data) => setPresentation(data),
        (e) => setError(e),
      )
      .finally(() => setIsLoading(false));
  }, [id]);

  function saveContent(title: string, content: string) {
    if (isSaving) return;

    setIsSaving(true);
    setError(null);

    return updatePresentation(id, title, content)
      .match(
        (data) => setPresentation(data),
        (e) => setError(e),
      )
      .finally(() => setIsSaving(false));
  }

  return { presentation, isLoading, isSaving, error, saveContent };
}
