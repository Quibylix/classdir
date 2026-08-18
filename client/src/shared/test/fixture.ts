import { describe, expect, it, test as base } from "vitest";
import type { SetupWorker } from "msw/browser";
import { worker } from "./browser";

let started = false;

export const test = base.extend<{ worker: SetupWorker }>({
  worker: [
    // oxlint-disable-next-line no-empty-pattern
    async ({}, use) => {
      if (!started) {
        await worker.start();
        started = true;
      }
      await use(worker);
      worker.resetHandlers();
    },
    { auto: true },
  ],
});

export { describe, it, expect };
