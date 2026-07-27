import { useRef, useEffect, useState, useCallback } from "react";
import type { EditorView } from "codemirror";

export function useCodeMirrorEditor(initialContent: string, initialReadOnly: boolean) {
  const editorRef = useRef<HTMLDivElement>(null);
  const viewRef = useRef<EditorView | null>(null);
  const [doc, setDoc] = useState(initialContent);

  const readOnlyRef = useRef(initialReadOnly);
  const pendingDocRef = useRef(doc);

  const setReadOnlyRef = useRef<((v: boolean) => void) | null>(null);

  useEffect(() => {
    let isMounted = true;

    (async () => {
      if (!editorRef.current) return;

      const { EditorView, basicSetup } = await import("codemirror");
      const { Compartment, EditorState } = await import("@codemirror/state");
      const { html } = await import("@codemirror/lang-html");

      if (!isMounted || !editorRef.current) return;

      const readOnlyCompartment = new Compartment();
      editorRef.current.innerHTML = "";

      viewRef.current = new EditorView({
        doc: pendingDocRef.current,
        extensions: [
          basicSetup,
          html(),
          readOnlyCompartment.of(EditorState.readOnly.of(readOnlyRef.current)),
          EditorView.theme(
            { "&": { height: "100%" }, ".cm-scroller": { overflow: "auto" } },
            { dark: true },
          ),
          EditorView.updateListener.of((update) => {
            if (update.docChanged) setDoc(update.state.doc.toString());
          }),
        ],
        parent: editorRef.current,
      });

      setReadOnlyRef.current = (readOnly: boolean) => {
        viewRef.current?.dispatch({
          effects: readOnlyCompartment.reconfigure(EditorState.readOnly.of(readOnly)),
        });
      };
      setReadOnlyRef.current(readOnlyRef.current);
    })();

    return () => {
      isMounted = false;
      viewRef.current?.destroy();
      viewRef.current = null;
    };
  }, []);

  const setContent = useCallback((content: string) => {
    setDoc(content);
    pendingDocRef.current = content;
    if (viewRef.current && viewRef.current.state.doc.toString() !== content) {
      viewRef.current.dispatch({
        changes: { from: 0, to: viewRef.current.state.doc.length, insert: content },
      });
    }
  }, []);

  const setReadOnly = useCallback((readOnly: boolean) => {
    readOnlyRef.current = readOnly;
    setReadOnlyRef.current?.(readOnly);
  }, []);

  return { editorRef, doc, setContent, setReadOnly };
}
