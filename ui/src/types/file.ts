export interface File {
  id: string;
  createdAt: string;
  updatedAt: string;
  path: string;
  pathRemote: string | null;
  size: number | null;
}