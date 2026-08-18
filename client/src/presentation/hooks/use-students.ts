import { useState, useEffect, useCallback } from "react";
import { listStudents, createStudent, updateStudent, deleteStudent } from "../api";
import { uuidv7 } from "../../shared/util/uuid";
import type { Student } from "../types";
import type { FetchError } from "../../shared/api/fetch";

export function useStudents(presentationId: string) {
  const [students, setStudents] = useState<Student[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [isCreating, setIsCreating] = useState(false);
  const [isUpdating, setIsUpdating] = useState(false);
  const [isDeleting, setIsDeleting] = useState(false);
  const [error, setError] = useState<FetchError | null>(null);

  const refresh = useCallback(() => {
    if (!presentationId) return;

    setIsLoading(true);
    setError(null);
    listStudents(presentationId).match(
      (data) => {
        setStudents(data);
        setIsLoading(false);
      },
      (e) => {
        setError(e);
        setIsLoading(false);
      },
    );
  }, [presentationId]);

  useEffect(() => {
    refresh();
  }, [refresh]);

  function create(name: string) {
    if (isCreating || !presentationId) return;

    const id = uuidv7();
    setIsCreating(true);
    setError(null);
    return createStudent(presentationId, id, name).match(
      () => {
        setIsCreating(false);
        refresh();
      },
      (e) => {
        setIsCreating(false);
        setError(e);
      },
    );
  }

  function rename(studentId: string, name: string) {
    if (isUpdating || !presentationId) return;

    setIsUpdating(true);
    setError(null);
    return updateStudent(presentationId, studentId, name).match(
      () => {
        setIsUpdating(false);
        refresh();
      },
      (e) => {
        setIsUpdating(false);
        setError(e);
      },
    );
  }

  function remove(studentId: string) {
    if (isDeleting || !presentationId) return;

    setIsDeleting(true);
    setError(null);
    return deleteStudent(presentationId, studentId).match(
      () => {
        setIsDeleting(false);
        refresh();
      },
      (e) => {
        setIsDeleting(false);
        setError(e);
      },
    );
  }

  return {
    students,
    isLoading,
    isCreating,
    isUpdating,
    isDeleting,
    error,
    refresh,
    create,
    rename,
    remove,
  };
}
