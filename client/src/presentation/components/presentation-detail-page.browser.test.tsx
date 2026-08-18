import { render } from "vitest-browser-react";
import { MemoryRouter, Route, Routes } from "react-router";
import { MantineProvider } from "@mantine/core";
import { test, expect } from "../../shared/test/fixture";
import { resetStudentApi, seedPresentation, seedStudents } from "../../shared/test/handlers";
import { clientConfigure } from "../../shared/cfg/routes";
import PresentationDetailPage from "./presentation-detail-page";

const presentationId = "0192e5a0-7b7f-7b7f-8b7f-0192e5a07b7f";

function renderDetailPage() {
  return render(
    <MantineProvider defaultColorScheme="dark">
      <MemoryRouter initialEntries={[`/configure/${presentationId}`]}>
        <Routes>
          <Route path={clientConfigure(":id")} element={<PresentationDetailPage />} />
        </Routes>
      </MemoryRouter>
    </MantineProvider>,
  );
}

test("renders a presentation with the students panel", async () => {
  resetStudentApi();
  seedPresentation({
    id: presentationId,
    title: "Intro to Math",
    content: "<h1>Hello</h1>",
  });
  seedStudents(presentationId, [{ id: "0192e5a0-7b7f-7b7f-8b7f-0192e5a07b80", name: "Alice" }]);

  const screen = await renderDetailPage();

  await expect.element(screen.getByText("Intro to Math")).toBeVisible();
  await expect.element(screen.getByRole("button", { name: /show students \(1\)/i })).toBeVisible();
});

test("shows an error message for a missing presentation", async () => {
  resetStudentApi();

  const screen = await renderDetailPage();

  await expect.element(screen.getByText(/presentation not found/i)).toBeVisible();
  await expect.element(screen.getByText("Back to Presentations")).toBeVisible();
});
