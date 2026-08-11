import { test, expect } from "../shared/test/fixture";
import { AUTH_CHECK } from "../shared/cfg/routes";
import { api } from "../shared/api/client";
import { listPresentations } from "./api";

test("MSW worker is running and intercepts AUTH_CHECK", async ({ worker }) => {
  const res = await api(AUTH_CHECK);
  expect(res.status).toBe(204);
  expect(res.ok).toBe(true);
  expect(worker.listHandlers().length).toBeGreaterThan(0);
});

test("MSW serves an empty presentation list through the API client", async () => {
  const list = await listPresentations().unwrapOr([{ id: "n/a", title: "n/a" }]);
  expect(list).toEqual([]);
});
