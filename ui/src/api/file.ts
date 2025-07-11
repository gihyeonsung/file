import type { File } from "../types/file";

export const getFiles = async (
  ids: string[] | null,
  paths: string[] | null,
  pathsLike: string[] | null
): Promise<File[]> => {
  const url = new URL("/api/v1/files", window.location.origin);
  if (ids) {
    url.searchParams.set("ids", ids.join(","));
  }
  if (paths) {
    url.searchParams.set("paths", paths.join(","));
  }
  if (pathsLike) {
    url.searchParams.set("paths-like", pathsLike.join(","));
  }

  const response = await fetch(url);
  const responseJson = await response.json();
  return responseJson.data;
};

export const getFilesId = async (id: string): Promise<Blob> => {
  const response = await fetch(`/api/v1/files/${id}`, {
    method: "GET",
  });
  return response.blob();
};

export const postFiles = async (path: string): Promise<File> => {
  const response = await fetch(`/api/v1/files`, {
    method: "POST",
    body: JSON.stringify({ path }),
  });
  return response.json();
};

export const postFilesId = async (
  id: string,
  data: FormData
): Promise<null> => {
  const response = await fetch(`/api/v1/files/${id}`, {
    method: "POST",
    body: data,
  });
  await response.json();

  return null;
};

export const deleteFilesId = async (id: string): Promise<File> => {
  const response = await fetch(`/api/v1/files/${id}`, {
    method: "DELETE",
  });
  return response.json();
};
