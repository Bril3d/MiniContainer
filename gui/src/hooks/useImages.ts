import { useState, useEffect, useCallback } from "react";
import { Command } from "@tauri-apps/plugin-shell";

export interface Image {
  id: string;
  repository: string;
  tag: string;
  size: number;
  names: string[];
}

export function useImages() {
  const [images, setImages] = useState<Image[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const refreshAction = useCallback(async () => {
    setLoading(true);
    try {
      const sidecar = Command.sidecar("bin/minicontainer", ["images", "--json"]);
      const output = await sidecar.execute();
      
      if (output.code !== 0) {
        throw new Error(output.stderr || "Failed to list images");
      }

      const parsed: Image[] = JSON.parse(output.stdout);
      setImages(parsed);
      setError(null);
    } catch (err: any) {
      setError(err.message);
      // Fallback/Mock for demonstration if podman fails
      if (images.length === 0) {
        setImages([
          { id: "sha256:123", repository: "python", tag: "3.11", size: 120 * 1024 * 1024, names: ["python:3.11"] },
          { id: "sha256:456", repository: "node", tag: "20", size: 350 * 1024 * 1024, names: ["node:20"] },
        ]);
      }
    } finally {
      setLoading(false);
    }
  }, [images.length]);

  useEffect(() => {
    refreshAction();
  }, [refreshAction]);

  const removeImage = async (id: string) => {
    const sidecar = Command.sidecar("bin/minicontainer", ["rmi", id]);
    await sidecar.execute();
    await refreshAction();
  };

  return { images, loading, error, refreshAction, removeImage };
}
