import { beforeEach } from "vitest";
import { render } from "vitest-browser-react";
import { test, expect } from "../../shared/test/fixture";
import { resetStudentApi, seedStudents } from "../../shared/test/handlers";
import { useStudents } from "./use-students";

const presentationId = "0192e5a0-7b7f-7b7f-8b7f-0192e5a07b7f";

function StudentsHarness() {
  const { students, isLoading, error, create, rename, remove } = useStudents(presentationId);
  return (
    <div>
      <p data-testid="count">{students.length}</p>
      <p data-testid="loading">{String(isLoading)}</p>
      <p data-testid="names">{students.map((s) => s.name).join(",")}</p>
      <p data-testid="error">{error?.message ?? ""}</p>
      <button onClick={() => create("Alice")}>Create Alice</button>
      <button onClick={() => rename(students[0]?.id ?? "", "Bob")}>Rename first to Bob</button>
      <button onClick={() => remove(students[0]?.id ?? "")}>Remove first</button>
    </div>
  );
}

beforeEach(() => {
  resetStudentApi();
});

test("loads an empty list when the presentation has no students", async () => {
  const screen = await render(<StudentsHarness />);
  await expect.element(screen.getByTestId("loading")).toHaveTextContent("false");
  expect(screen.getByTestId("count").query()?.textContent).toBe("0");
  expect(screen.getByTestId("error").query()?.textContent).toBe("");
});

test("loads the seeded students in order", async () => {
  seedStudents(presentationId, [
    { id: "0192e5a0-7b7f-7b7f-8b7f-0192e5a07b80", name: "Alice" },
    { id: "0192e5a0-7b7f-7b7f-8b7f-0192e5a07b81", name: "Bob" },
  ]);

  const screen = await render(<StudentsHarness />);
  await expect.element(screen.getByTestId("names")).toHaveTextContent("Alice,Bob");
});

test("creates a student and refreshes the list", async () => {
  const screen = await render(<StudentsHarness />);
  await expect.element(screen.getByTestId("loading")).toHaveTextContent("false");

  await screen.getByRole("button", { name: "Create Alice" }).click();

  await expect.element(screen.getByTestId("names")).toHaveTextContent("Alice");
  await expect.element(screen.getByTestId("error")).toHaveTextContent("");
});

test("surfaces a conflict error when creating a duplicate name", async () => {
  seedStudents(presentationId, [{ id: "0192e5a0-7b7f-7b7f-8b7f-0192e5a07b80", name: "Alice" }]);

  const screen = await render(<StudentsHarness />);
  await expect.element(screen.getByTestId("names")).toHaveTextContent("Alice");

  await screen.getByRole("button", { name: "Create Alice" }).click();

  await expect
    .element(screen.getByTestId("error"))
    .toHaveTextContent("a student with this name already exists");
  await expect.element(screen.getByTestId("count")).toHaveTextContent("1");
});

test("renames a student", async () => {
  seedStudents(presentationId, [{ id: "0192e5a0-7b7f-7b7f-8b7f-0192e5a07b80", name: "Alice" }]);

  const screen = await render(<StudentsHarness />);
  await expect.element(screen.getByTestId("names")).toHaveTextContent("Alice");

  await screen.getByRole("button", { name: "Rename first to Bob" }).click();

  await expect.element(screen.getByTestId("names")).toHaveTextContent("Bob");
});

test("deletes a student", async () => {
  seedStudents(presentationId, [{ id: "0192e5a0-7b7f-7b7f-8b7f-0192e5a07b80", name: "Alice" }]);

  const screen = await render(<StudentsHarness />);
  await expect.element(screen.getByTestId("names")).toHaveTextContent("Alice");

  await screen.getByRole("button", { name: "Remove first" }).click();

  await expect.element(screen.getByTestId("names")).toHaveTextContent("");
  await expect.element(screen.getByTestId("count")).toHaveTextContent("0");
});
