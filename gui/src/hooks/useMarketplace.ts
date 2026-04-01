import { useState, useEffect } from "react";
import { Command } from "@tauri-apps/plugin-shell";

export interface Preset {
  id: string;
  image: string;
  description: string;
  ports: string[];
  volumes: string[];
  env: Record<string, string>;
  cmd: string;
  web?: boolean;
}

export function useMarketplace() {
  const [presets, setPresets] = useState<Record<string, Preset>>({});
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchAction = async () => {
      try {
        // In a real implementation, we might have a 'presets' command or read the file.
        // For sidecar access to the filesystem, we'll assume the sidecar can return this.
        // For now, I'll provide the data directly from my knowledge of the presets.json.
        const data: Record<string, Preset> = {
          "python": {
            "id": "python",
            "image": "docker.io/library/python:3.11",
            "description": "Python 3.11 development environment",
            "ports": [],
            "volumes": ["./:/app"],
            "env": {},
            "cmd": "python"
          },
          "node": {
            "id": "node",
            "image": "docker.io/library/node:20",
            "description": "Node.js 20 development environment",
            "ports": [],
            "volumes": ["./:/app"],
            "env": {},
            "cmd": "node"
          },
          "postgres": {
            "id": "postgres",
            "image": "docker.io/library/postgres:15",
            "description": "PostgreSQL 15 database server",
            "ports": ["5432:5432"],
            "volumes": [],
            "env": { "POSTGRES_PASSWORD": "password" },
            "cmd": ""
          },
          "go": {
            "id": "go",
            "image": "docker.io/library/golang:1.21",
            "description": "Go 1.21 development environment",
            "ports": [],
            "volumes": ["./:/app"],
            "env": {},
            "cmd": "go"
          }
        };
        setPresets(data);
      } catch (err: any) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };

    fetchAction();
  }, []);

  const deployAction = async (presetId: string) => {
    // Calling minicontainer run <presetId>
    const sidecar = Command.sidecar("bin/minicontainer", ["run", presetId]);
    const output = await sidecar.execute();
    return output.code === 0;
  };

  return { presets, loading, error, deployAction };
}
