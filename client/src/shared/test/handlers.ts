import { http, HttpResponse } from "msw";
import {
  AUTH_CHECK,
  PRESENTATIONS,
  presentationById,
  studentsByPresentation,
  studentById,
} from "../cfg/routes";
import type { Student } from "../../presentation/types";

type StoredPresentation = { id: string; title: string; content: string };

const studentStore = new Map<string, Student[]>();
const presentationStore = new Map<string, StoredPresentation>();

export function seedPresentation(presentation: StoredPresentation) {
  presentationStore.set(presentation.id, { ...presentation });
}

export function seedStudents(presentationId: string, students: Student[]) {
  studentStore.set(
    presentationId,
    students.map((s) => ({ ...s })),
  );
}

export function resetStudentApi() {
  studentStore.clear();
  presentationStore.clear();
}

export const handlers = [
  http.get(AUTH_CHECK, () => new HttpResponse(null, { status: 204 })),
  http.get(PRESENTATIONS, () => HttpResponse.json({ data: [] })),
  http.get(presentationById(":id"), ({ params }) => {
    const presentation = presentationStore.get(String(params.id));
    if (!presentation) {
      return HttpResponse.json(
        { error: { code: "NOT_FOUND", message: "presentation not found" } },
        { status: 404 },
      );
    }
    return HttpResponse.json({ data: presentation });
  }),
  http.get(studentsByPresentation(":presentationId"), ({ params }) => {
    const students = studentStore.get(String(params.presentationId)) ?? [];
    return HttpResponse.json({ data: students });
  }),
  http.post(studentsByPresentation(":presentationId"), async ({ request, params }) => {
    const presentationId = String(params.presentationId);
    const body = (await request.json()) as { id: string; name: string };
    const students = studentStore.get(presentationId) ?? [];
    if (students.some((s) => s.id === body.id)) {
      return HttpResponse.json(
        { error: { code: "CONFLICT", message: "a student with this id already exists" } },
        { status: 409 },
      );
    }
    if (students.some((s) => s.name === body.name)) {
      return HttpResponse.json(
        { error: { code: "CONFLICT", message: "a student with this name already exists" } },
        { status: 409 },
      );
    }
    const student = { id: body.id, name: body.name };
    students.push(student);
    studentStore.set(presentationId, students);
    return HttpResponse.json({ data: student }, { status: 201 });
  }),
  http.put(studentById(":presentationId", ":studentId"), async ({ request, params }) => {
    const presentationId = String(params.presentationId);
    const studentId = String(params.studentId);
    const body = (await request.json()) as { name: string };
    const students = studentStore.get(presentationId) ?? [];
    const student = students.find((s) => s.id === studentId);
    if (!student) {
      return HttpResponse.json(
        { error: { code: "NOT_FOUND", message: "student not found" } },
        { status: 404 },
      );
    }
    if (students.some((s) => s.id !== studentId && s.name === body.name)) {
      return HttpResponse.json(
        { error: { code: "CONFLICT", message: "a student with this name already exists" } },
        { status: 409 },
      );
    }
    student.name = body.name;
    return HttpResponse.json({ data: { ...student } });
  }),
  http.delete(studentById(":presentationId", ":studentId"), ({ params }) => {
    const presentationId = String(params.presentationId);
    const studentId = String(params.studentId);
    const students = studentStore.get(presentationId) ?? [];
    const next = students.filter((s) => s.id !== studentId);
    if (next.length === students.length) {
      return HttpResponse.json(
        { error: { code: "NOT_FOUND", message: "student not found" } },
        { status: 404 },
      );
    }
    studentStore.set(presentationId, next);
    return new HttpResponse(null, { status: 204 });
  }),
];
