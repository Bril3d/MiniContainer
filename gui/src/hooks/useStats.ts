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
      let rawStats: any[];

      // Check if we are running in Tauri
      if ((window as any).__TAURI_INTERNALS__) {
        const command = Command.sidecar("bin/minicontainer", ["stats", "--json"]);
        const output = await command.execute();

        if (output.code !== 0) {
          throw new Error(output.stderr || "Failed to fetch stats");
        }

        rawStats = JSON.parse(output.stdout);
      } else {
        // Fallback to API server
        const response = await fetch('http://localhost:8080/api/stats');
        if (!response.ok) {
          throw new Error('Failed to fetch stats from API server');
        }
        rawStats = await response.json();
      }

      const statsMap: Record<string, ContainerStats> = {};
      
      (rawStats || []).forEach((s) => {
        // Handle both PascalCase and lowercase from different runtime outputs
        const id = s.ID || s.id || "";
        const name = s.Name || s.name || "unknown";
        const cpu = s.CPUPerc || s.cpu_perc || "0%";
        const mem = s.MemUsage || s.mem_usage || "0MB";
        const memPerc = s.MemPerc || s.mem_perc || "0%";
        const net = s.NetIO || s.net_io || "0B / 0B";

        const shortId = id.substring(0, 12);
        statsMap[shortId] = {
          id: shortId,
          name: name,
          cpu_perc: cpu,
          mem_usage: mem,
          mem_perc: memPerc,
          net_io: net
        };
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
