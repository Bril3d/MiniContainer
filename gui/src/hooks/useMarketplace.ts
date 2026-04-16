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
        const response = await fetch('http://localhost:8080/api/presets');
        if (!response.ok) throw new Error('Failed to fetch presets');
        const data = await response.json();
        
        // Convert to Record if necessary, or just set data
        setPresets(data);
      } catch (err: any) {
        console.error("Fetch presets error:", err);
        setError("Could not load presets from backend. Is the API running?");
      } finally {
        setLoading(false);
      }
    };

    fetchAction();
  }, []);

  const deployAction = async (presetId: string) => {
    try {
      if ((window as any).__TAURI_INTERNALS__) {
        const sidecar = Command.sidecar("bin/minicontainer", ["run", presetId]);
        const output = await sidecar.execute();
        return output.code === 0;
      } else {
        // Fallback to API server
        const response = await fetch(`http://localhost:8080/api/deploy?id=${presetId}`, { method: 'POST' });
        if (!response.ok) throw new Error('Failed to deploy preset');
        return true;
      }
    } catch (err: any) {
      console.error("Deploy error:", err);
      setError(err.message);
      return false;
    }
  };

  const clearError = () => setError(null);
  
  return { presets, loading, error, deployAction, clearError };
}
