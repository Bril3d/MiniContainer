import { useState, useEffect, useCallback } from "react";
import { Command } from "@tauri-apps/plugin-shell";

export interface Container {
  id: string;
  image: string;
  status: string;
  ports: string;
  names: string;
}

export function useContainers() {
  const [containers, setContainers] = useState<Container[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const refreshAction = useCallback(async () => {
    setLoading(true);
    try {
      // In a real implementation with JSON output from Go:
      // const sidecar = Command.sidecar("bin/minicontainer", ["list", "--json"]);
      
      // For now, parsing the table output (or assuming a simplified JSON-like response)
      // Since our CLI 'list' currently returns a table, we'll need either JSON support in CLI 
      // or a clever parser. 
      // ACTION: I'll assume we might want to add --json to our Go CLI later.
      // For now, I'll mock the data to get the UI right, then we'll fix the bridge.
      
      const sidecar = Command.sidecar("bin/minicontainer", ["list"]);
      const output = await sidecar.execute();
      
      if (output.code !== 0) {
        throw new Error(output.stderr || "Failed to list containers");
      }

      // Mocking for Phase 3 visual fidelity while bridge is refined
      // Real implementation would parse 'output.stdout'
      const mockData: Container[] = [
        { id: "a1b2c3d4", image: "python:3.9", status: "Running", ports: "8080->80", names: "python-env" },
        { id: "e5f6g7h8", image: "node:18", status: "Stopped", ports: "3000->3000", names: "node-app" },
      ];
      setContainers(mockData);
      setError(null);
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    refreshAction();
  }, [refreshAction]);

  const startContainer = async (name: string) => {
    const sidecar = Command.sidecar("bin/minicontainer", ["run", name]);
    await sidecar.execute();
    await refreshAction();
  };

  const stopContainer = async (name: string) => {
    // Assuming 'stop' command exists or will be added
    const sidecar = Command.sidecar("bin/minicontainer", ["stop", name]);
    await sidecar.execute();
    await refreshAction();
  };

  const removeContainer = async (name: string) => {
    const sidecar = Command.sidecar("bin/minicontainer", ["remove", name]);
    await sidecar.execute();
    await refreshAction();
  };

  return { containers, loading, error, refreshAction, startContainer, stopContainer, removeContainer };
}
