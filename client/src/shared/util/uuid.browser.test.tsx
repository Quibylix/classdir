import { describe, it, expect } from "../test/fixture";
import { render } from "vitest-browser-react";
import { uuidv7 } from "./uuid";

describe("uuidv7 (browser)", () => {
  it("renders a generated uuid into the DOM", async () => {
    const id = uuidv7();
    const screen = await render(<p data-testid="uid">{id}</p>);
    await expect
      .element(screen.getByTestId("uid"))
      .toHaveTextContent(/^[0-9a-f]{8}-[0-9a-f]{4}-7[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/);
  });
});
