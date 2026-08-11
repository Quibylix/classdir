import { http, HttpResponse } from "msw";
import { AUTH_CHECK, PRESENTATIONS } from "../cfg/routes";

export const handlers = [
  http.get(AUTH_CHECK, () => new HttpResponse(null, { status: 204 })),
  http.get(PRESENTATIONS, () => HttpResponse.json({ data: [] })),
];
