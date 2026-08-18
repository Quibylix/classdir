import { describe, it, expect, vi, beforeEach } from "vitest";
import { api } from "../shared/api/client";
import { studentsByPresentation, studentById } from "../shared/cfg/routes";
import { listStudents, getStudent, createStudent, updateStudent, deleteStudent } from "./api";

vi.mock("../shared/api/client", () => ({ api: vi.fn() }));

const presentationId = "0192e5a0-7b7f-7b7f-8b7f-0192e5a07b7f";
const studentId = "0192e5a0-7b7f-7b7f-8b7f-0192e5a07b80";

function mockJson(data: unknown, status = 200) {
  vi.mocked(api).mockResolvedValue(
    new Response(JSON.stringify({ data }), {
      status,
      headers: { "Content-Type": "application/json" },
    }),
  );
}

describe("student api", () => {
  beforeEach(() => {
    vi.mocked(api).mockReset();
  });

  it("listStudents requests the students collection", async () => {
    mockJson([]);
    const result = await listStudents(presentationId);
    expect(api).toHaveBeenCalledWith(studentsByPresentation(presentationId), undefined);
    expect(result.isOk()).toBe(true);
    expect(result._unsafeUnwrap()).toEqual([]);
  });

  it("listStudents parses the student array from the response", async () => {
    const students = [
      { id: studentId, name: "Alice" },
      { id: "0192e5a0-7b7f-7b7f-8b7f-0192e5a07b81", name: "Bob" },
    ];
    mockJson(students);
    const result = await listStudents(presentationId);
    expect(result.isOk()).toBe(true);
    expect(result._unsafeUnwrap()).toEqual(students);
  });

  it("getStudent requests the single student", async () => {
    mockJson({ id: studentId, name: "Alice" });
    const result = await getStudent(presentationId, studentId);
    expect(api).toHaveBeenCalledWith(studentById(presentationId, studentId), undefined);
    expect(result.isOk()).toBe(true);
  });

  it("createStudent posts id and name to the students collection", async () => {
    mockJson({ id: studentId, name: "Alice" }, 201);
    const result = await createStudent(presentationId, studentId, "Alice");
    expect(api).toHaveBeenCalledWith(studentsByPresentation(presentationId), {
      method: "POST",
      body: JSON.stringify({ id: studentId, name: "Alice" }),
    });
    expect(result.isOk()).toBe(true);
    expect(result._unsafeUnwrap()).toEqual({ id: studentId, name: "Alice" });
  });

  it("updateStudent puts the new name to the single student", async () => {
    mockJson({ id: studentId, name: "Bob" });
    const result = await updateStudent(presentationId, studentId, "Bob");
    expect(api).toHaveBeenCalledWith(studentById(presentationId, studentId), {
      method: "PUT",
      body: JSON.stringify({ name: "Bob" }),
    });
    expect(result.isOk()).toBe(true);
    expect(result._unsafeUnwrap()).toEqual({ id: studentId, name: "Bob" });
  });

  it("deleteStudent deletes the single student", async () => {
    vi.mocked(api).mockResolvedValue(new Response(null, { status: 204 }));
    const result = await deleteStudent(presentationId, studentId);
    expect(api).toHaveBeenCalledWith(studentById(presentationId, studentId), {
      method: "DELETE",
    });
    expect(result.isOk()).toBe(true);
  });
});
