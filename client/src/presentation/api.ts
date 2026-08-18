import { z } from "zod";
import { safeFetch } from "../shared/api/fetch";
import {
  PRESENTATIONS,
  presentationById,
  studentsByPresentation,
  studentById,
} from "../shared/cfg/routes";
import { HTTP_METHOD_POST, HTTP_METHOD_PUT, HTTP_METHOD_DELETE } from "../shared/cfg/http";
import { PresentationSchema, PresentationPreviewSchema, StudentSchema } from "./types";

export function listPresentations() {
  return safeFetch(PRESENTATIONS, z.array(PresentationPreviewSchema));
}

export function getPresentation(id: string) {
  return safeFetch(presentationById(id), PresentationSchema);
}

export function createPresentation(id: string, title: string) {
  return safeFetch(PRESENTATIONS, PresentationSchema, {
    method: HTTP_METHOD_POST,
    body: JSON.stringify({ id, title }),
  });
}

export function updatePresentation(id: string, title: string, content: string) {
  return safeFetch(presentationById(id), PresentationSchema, {
    method: HTTP_METHOD_PUT,
    body: JSON.stringify({ title, content }),
  });
}

export function deletePresentation(id: string) {
  return safeFetch(presentationById(id), z.undefined(), { method: HTTP_METHOD_DELETE });
}

export function listStudents(presentationId: string) {
  return safeFetch(studentsByPresentation(presentationId), z.array(StudentSchema));
}

export function getStudent(presentationId: string, studentId: string) {
  return safeFetch(studentById(presentationId, studentId), StudentSchema);
}

export function createStudent(presentationId: string, id: string, name: string) {
  return safeFetch(studentsByPresentation(presentationId), StudentSchema, {
    method: HTTP_METHOD_POST,
    body: JSON.stringify({ id, name }),
  });
}

export function updateStudent(presentationId: string, studentId: string, name: string) {
  return safeFetch(studentById(presentationId, studentId), StudentSchema, {
    method: HTTP_METHOD_PUT,
    body: JSON.stringify({ name }),
  });
}

export function deleteStudent(presentationId: string, studentId: string) {
  return safeFetch(studentById(presentationId, studentId), z.undefined(), {
    method: HTTP_METHOD_DELETE,
  });
}
