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
    try {
      setLoading(true);
      setError(null);

      let rawData: any[];

      // Check if we are running in Tauri
      if ((window as any).__TAURI_INTERNALS__) {
        const sidecar = Command.sidecar("bin/minicontainer", ["ps", "--json"]);
        const output = await sidecar.execute();
        
        if (output.code !== 0) {
          throw new Error(output.stderr || "Failed to list containers");
        }

        rawData = JSON.parse(output.stdout);
      } else {
        // Fallback to API server
        const response = await fetch('http://localhost:8080/api/ps');
        if (!response.ok) {
          throw new Error('Failed to fetch from API server');
        }
        rawData = await response.json();
      }

      // Map Go runtime Container struct (PascalCase JSON tags) to Frontend Container interface
      const mappedData: Container[] = rawData.map((c: any) => {
        const id = (c.ID || "").substring(0, 12);
        const name = c.Names && c.Names.length > 0 ? c.Names[0].replace(/^\//, '') : "unnamed";
        
        let portStr = "";
        if (c.Ports) {
          portStr = c.Ports.map((p: any) => 
            p.host_port > 0 ? `${p.host_port}->${p.container_port}/${p.protocol}` : `${p.container_port}/${p.protocol}`
          ).join(", ");
        }

        return {
          id: id,
          image: c.Image || "unknown",
          status: c.Status || "Unknown",
          ports: portStr,
          names: name
        };
      });

      setContainers(mappedData);
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
    if ((window as any).__TAURI_INTERNALS__) {
      const sidecar = Command.sidecar("bin/minicontainer", ["start", name]);
      await sidecar.execute();
    } else {
      const response = await fetch(`http://localhost:8080/api/start?name=${name}`, { method: 'POST' });
      if (!response.ok) throw new Error('Failed to start container');
    }
    await refreshAction();
  };

  const stopContainer = async (name: string) => {
    if ((window as any).__TAURI_INTERNALS__) {
      const sidecar = Command.sidecar("bin/minicontainer", ["stop", name]);
      await sidecar.execute();
    } else {
      const response = await fetch(`http://localhost:8080/api/stop?name=${name}`, { method: 'POST' });
      if (!response.ok) throw new Error('Failed to stop container');
    }
    await refreshAction();
  };

  const removeContainer = async (name: string) => {
    if ((window as any).__TAURI_INTERNALS__) {
      const sidecar = Command.sidecar("bin/minicontainer", ["rm", "--force", name]);
      await sidecar.execute();
    } else {
      const response = await fetch(`http://localhost:8080/api/remove?name=${name}`, { method: 'POST' });
      if (!response.ok) throw new Error('Failed to remove container');
    }
    await refreshAction();
  };

  const pauseContainer = async (id: string) => {
    try {
      setLoading(true);
      if ((window as any).__TAURI_INTERNALS__) {
        const { invoke } = await import('@tauri-apps/api/core');
        await invoke('exec_podman', { args: ['pause', id] });
      } else {
        const res = await fetch(`http://localhost:8080/api/pause?id=${id}`, { method: 'POST' });
        if (!res.ok) throw new Error(await res.text());
      }
      await refreshAction();
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const unpauseContainer = async (id: string) => {
    try {
      setLoading(true);
      if ((window as any).__TAURI_INTERNALS__) {
        const { invoke } = await import('@tauri-apps/api/core');
        await invoke('exec_podman', { args: ['unpause', id] });
      } else {
        const res = await fetch(`http://localhost:8080/api/unpause?id=${id}`, { method: 'POST' });
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
  
  return { containers, loading, error, refreshAction, startContainer, stopContainer, removeContainer, pauseContainer, unpauseContainer, clearError };
}
