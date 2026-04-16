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
    setError(null);
    try {
      let rawImages: any[];

      // Check if we are running in Tauri
      if ((window as any).__TAURI_INTERNALS__) {
        const sidecar = Command.sidecar("bin/minicontainer", ["images", "--json"]);
        const output = await sidecar.execute();
        
        if (output.code !== 0) {
          throw new Error(output.stderr || "Failed to list images");
        }

        rawImages = JSON.parse(output.stdout);
      } else {
        // Fallback to API server
        const response = await fetch('http://localhost:8080/api/images');
        if (!response.ok) {
          throw new Error('Failed to fetch from API server');
        }
        rawImages = await response.json();
      }

      const mappedImages: Image[] = (rawImages || []).map((img: any) => {
        // Higher-level mapping that handles both cases and uses Names as fallback
        const names = img.Names || img.names || [];
        let repo = img.Repository || img.repository || "";
        let tag = img.Tag || img.tag || "";

        if ((!repo || repo === "<none>") && names.length > 0) {
          const parts = names[0].split(':');
          if (parts.length > 1) {
            repo = parts.slice(0, -1).join(':');
            tag = parts[parts.length - 1];
          } else {
            repo = names[0];
            tag = "latest";
          }
        }

        return {
          id: img.ID || img.Id || img.id,
          repository: repo || "unknown",
          tag: tag || "latest",
          size: img.Size || img.size || 0,
          names: names
        };
      });
      
      setImages(mappedImages);
      setError(null);
    } catch (err: any) {
      console.error("Images error:", err);
      setError(err.message);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    refreshAction();
  }, [refreshAction]);

  const removeImage = async (id: string) => {
    try {
      if ((window as any).__TAURI_INTERNALS__) {
        const sidecar = Command.sidecar("bin/minicontainer", ["rmi", "--force", id]);
        await sidecar.execute();
      } else {
        const response = await fetch(`http://localhost:8080/api/rmi?id=${id}`, { method: 'POST' });
        if (!response.ok) throw new Error('Failed to remove image');
      }
      await refreshAction();
    } catch (err: any) {
      console.error("Remove image error:", err);
      setError(err.message);
    }
  };

  const buildImage = async (tags: string[], context: string = ".", dockerfile?: string) => {
    try {
      setLoading(true);
      if ((window as any).__TAURI_INTERNALS__) {
        const args = ['build', ...tags.flatMap(t => ['-t', t])];
        if (dockerfile) args.push('-f', dockerfile);
        args.push(context);
        const sidecar = Command.sidecar("bin/minicontainer", args);
        await sidecar.execute();
      } else {
        const res = await fetch(`http://localhost:8080/api/build`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ tags, context, dockerfile })
        });
        if (!res.ok) throw new Error(await res.text());
      }
      await refreshAction();
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const clearError = () => setError(null);
  return { images, loading, error, refreshAction, removeImage, buildImage, clearError };
}
