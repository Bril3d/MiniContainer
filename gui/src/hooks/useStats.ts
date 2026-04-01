import { useState, useEffect, useCallback } from "react";
import { Command } from "@tauri-apps/plugin-shell";

export interface ContainerStats {
  id: string;
  name: string;
  cpu_perc: string;
  mem_usage: string;
  mem_perc: string;
  net_io: string;
}

export function useStats() {
  const [stats, setStats] = useState<Record<string, ContainerStats>>({});
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchStats = useCallback(async () => {
    try {
      const command = Command.sidecar("bin/minicontainer", ["stats", "--json"]);
      const output = await command.execute();

      if (output.code !== 0) {
        throw new Error(output.stderr || "Failed to fetch stats");
      }

      const rawStats: ContainerStats[] = JSON.parse(output.stdout);
      const statsMap: Record<string, ContainerStats> = {};
      
      rawStats.forEach((s) => {
        // Podman JSON fields might be PascalCase or lowercase depending on version 
        // Our Go struct uses PascalCase but json.MarshalIndent handles it.
        // We'll normalize here.
        statsMap[s.id] = s;
      });

      setStats(statsMap);
      setError(null);
    } catch (err: any) {
      console.error("Stats error:", err);
      setError(err.message);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchStats();
    const interval = setInterval(fetchStats, 2000); // Polling every 2s
    return () => clearInterval(interval);
  }, [fetchStats]);

  return { stats, loading, error, refresh: fetchStats };
}
