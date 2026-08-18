import { beforeEach, vi } from "vitest";
import { render } from "vitest-browser-react";
import { MantineProvider } from "@mantine/core";
import { test, expect } from "../../shared/test/fixture";
import { resetStudentApi, seedStudents } from "../../shared/test/handlers";
import { StudentsPanel } from "./students-panel";

const presentationId = "0192e5a0-7b7f-7b7f-8b7f-0192e5a07b7f";
const studentId = "0192e5a0-7b7f-7b7f-8b7f-0192e5a07b80";

function renderPanel() {
  return render(
    <MantineProvider defaultColorScheme="dark">
      <StudentsPanel presentationId={presentationId} />
    </MantineProvider>,
  );
}

beforeEach(() => {
  resetStudentApi();
});

test("toggles the student list open and closed", async () => {
  const screen = await renderPanel();

  await expect.element(screen.getByRole("button", { name: /show students \(0\)/i })).toBeVisible();
  await screen.getByRole("button", { name: /show students \(0\)/i }).click();

  await expect.element(screen.getByRole("button", { name: /hide students \(0\)/i })).toBeVisible();
  await expect.element(screen.getByText("No students registered yet")).toBeVisible();

  await screen.getByRole("button", { name: /hide students \(0\)/i }).click();
  await expect.element(screen.getByRole("button", { name: /show students \(0\)/i })).toBeVisible();
});

test("adds a student to the list", async () => {
  const screen = await renderPanel();

  await screen.getByRole("button", { name: /show students \(0\)/i }).click();
  const input = screen.getByPlaceholder("Student name");
  await input.fill("Alice");
  await screen.getByRole("button", { name: /^add$/i }).click();

  await expect.element(screen.getByText("Alice")).toBeVisible();
  await expect.element(screen.getByRole("button", { name: /students \(1\)/i })).toBeVisible();
});

test("shows a conflict error when adding a duplicate name", async () => {
  seedStudents(presentationId, [{ id: studentId, name: "Alice" }]);
  const screen = await renderPanel();

  await screen.getByRole("button", { name: /show students \(1\)/i }).click();
  await expect.element(screen.getByText("Alice")).toBeVisible();

  const input = screen.getByPlaceholder("Student name");
  await input.fill("Alice");
  await screen.getByRole("button", { name: /^add$/i }).click();

  await expect.element(screen.getByText("a student with this name already exists")).toBeVisible();
});

test("renames a student inline", async () => {
  seedStudents(presentationId, [{ id: studentId, name: "Alice" }]);
  const screen = await renderPanel();

  await screen.getByRole("button", { name: /show students \(1\)/i }).click();
  await expect.element(screen.getByText("Alice")).toBeVisible();

  await screen.getByRole("button", { name: "Edit student" }).click();
  const editInput = screen.getByPlaceholder("New name");
  await editInput.fill("Bob");
  await screen.getByRole("button", { name: /^save$/i }).click();

  await expect.element(screen.getByText("Bob")).toBeVisible();
  await vi.waitFor(() => {
    expect(screen.getByText("Alice").query()).toBeNull();
  });
});

test("deletes a student via the confirm modal", async () => {
  seedStudents(presentationId, [{ id: studentId, name: "Alice" }]);
  const screen = await renderPanel();

  await screen.getByRole("button", { name: /show students \(1\)/i }).click();
  await expect.element(screen.getByText("Alice")).toBeVisible();

  await screen.getByRole("button", { name: "Delete student" }).click();
  await expect.element(screen.getByRole("dialog")).toBeVisible();
  await expect.element(screen.getByText(/Are you sure you want to delete "Alice"/)).toBeVisible();

  await screen.getByRole("button", { name: /^delete$/i }).click();
  await expect.element(screen.getByText("No students registered yet")).toBeVisible();
});
