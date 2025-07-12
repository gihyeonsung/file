import { useEffect, useState } from "react";

import type { File } from "./types/file";
import {
  deleteFilesId,
  getFiles,
  getFilesId,
  postFiles,
  postFilesId,
} from "./api/file";

const File = ({
  file,
  handleDelete,
  uploadingProgress,
}: {
  file: File;
  handleDelete: () => void;
  uploadingProgress: number | null;
}) => {
  const [isHovering, setIsHovering] = useState(false);

  const sizeInBytes = file.size ?? 0;
  const sizeInKb = Math.round(sizeInBytes / 1024);
  const sizeInMb = Math.round(sizeInBytes / 1024 / 1024);
  const sizeInGb = Math.round(sizeInBytes / 1024 / 1024 / 1024);
  const sizeString =
    sizeInBytes < 1024
      ? `${sizeInBytes}B`
      : sizeInKb < 1024
      ? `${sizeInKb}KiB`
      : sizeInMb < 1024
      ? `${sizeInMb}MiB`
      : `${sizeInGb}GiB`;

  const handleClick = async (e: React.MouseEvent<HTMLDivElement>) => {
    e.preventDefault();

    try {
      const blob = await getFilesId(file.id);

      const url = window.URL.createObjectURL(blob);
      const link = document.createElement("a");
      link.href = url;

      const fileName = file.path.split("/").pop() || "download";
      link.download = fileName;

      document.body.appendChild(link);
      link.click();

      document.body.removeChild(link);
      window.URL.revokeObjectURL(url);
    } catch {
      alert("download failed");
    }
  };

  const handleMouseOver = (e: React.MouseEvent<HTMLDivElement>) => {
    e.preventDefault();
    setIsHovering(true);
  };

  const handleMouseLeave = (e: React.MouseEvent<HTMLDivElement>) => {
    e.preventDefault();
    setIsHovering(false);
  };

  const handleClickDelete = async (e: React.MouseEvent<HTMLDivElement>) => {
    e.preventDefault();
    e.stopPropagation();

    handleDelete();
  };

  return (
    <div
      className="flex flex-row gap-2 cursor-pointer hover:bg-neutral-100"
      onClick={handleClick}
      onMouseOver={handleMouseOver}
      onMouseLeave={handleMouseLeave}
    >
      <div className="text-sm text-neutral-500">{file.path}</div>
      <div className="text-sm text-neutral-500">{sizeString}</div>
      {uploadingProgress !== null && (
        <div className="text-sm text-neutral-500">
          {uploadingProgress.toFixed(2)}%
        </div>
      )}

      {isHovering && (
        <div className="text-sm text-neutral-500" onClick={handleClickDelete}>
          delete
        </div>
      )}
    </div>
  );
};

function App() {
  const [path] = useState<string>("/");
  const [files, setFiles] = useState<File[]>([]);
  const [uploadingFileId, setUploadingFileId] = useState<string | null>(null);
  const [uploadingProgress, setUploadingProgress] = useState<number | null>(null);

  useEffect(() => {
    (async () => {
      const files = await getFiles(null, null, [path]);
      setFiles(files);
    })();
  }, [path]);

  const handleDragOver = (e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault();
  };

  const handleDrop = async (e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault();

    for (const f of e.dataTransfer.files) {
      if (f.type === "" && f.size === 0) {
        continue;
      }
      
      const p = path + f.name;
      await postFiles(p);

      const fileUploaded = (await getFiles(null, [p], null))[0];
      setFiles((files) => [...files, fileUploaded].sort((a, b) => a.path.localeCompare(b.path)));

      const formData = new FormData();
      formData.append("file", f);
      setUploadingFileId(fileUploaded.id);
      setUploadingProgress(0);
      await postFilesId(fileUploaded.id, formData, (p) => {
        setUploadingProgress(p);
      });
      setUploadingFileId(null);
      setUploadingProgress(null);

      const files = await getFiles(null, null, [path]);
      setFiles(files);
    }
  };

  const handleDelete = async (id: string) => {
    await deleteFilesId(id);
    const files = await getFiles(null, null, [path]);
    setFiles(files);
  };

  return (
    <div
      className="flex flex-col items-center h-screen"
      onDragOver={handleDragOver}
      onDrop={handleDrop}
    >
      <div className="w-full h-full p-4">
        <h1 className="text-sm font-bold text-neutral-700">{path}</h1>

        <div className="flex flex-col">
          {files.map((f) => {
            return (
              <File
                key={f.id}
                file={f}
                handleDelete={() => handleDelete(f.id)}
                uploadingProgress={
                  uploadingFileId === f.id ? uploadingProgress : null
                }
              />
            );
          })}
        </div>
      </div>
    </div>
  );
}

export default App;
