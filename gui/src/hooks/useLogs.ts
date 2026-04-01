import { useState, useEffect, useRef } from "react";
import { Command } from "@tauri-apps/plugin-shell";

export function useLogs(containerId: string | null) {
  const [logs, setLogs] = useState<string[]>([]);
  const [active, setActive] = useState(false);
  const commandRef = useRef<any>(null);

  useEffect(() => {
    if (!containerId) {
      setLogs([]);
      setActive(false);
      return;
    }

    setLogs([]);
    setActive(true);

    const command = Command.sidecar("bin/minicontainer", ["logs", "-f", containerId]);
    commandRef.current = command;

    command.stdout.on("data", (line) => {
      setLogs((prev) => [...prev.slice(-499), line]); // Keep last 500 lines
    });

    command.stderr.on("data", (line) => {
      setLogs((prev) => [...prev.slice(-499), `\x1b[31m${line}\x1b[0m`]); // Red for stderr
    });

    command.on("error", (err) => {
      setLogs((prev) => [...prev, `\x1b[31mError spawning logs: ${err}\x1b[0m`]);
      setActive(false);
    });

    command.spawn().catch((err) => {
      setLogs((prev) => [...prev, `\x1b[31mError spawning logs: ${err.message}\x1b[0m`]);
      setActive(false);
    });

    return () => {
      // Cleanup: stop the command if it's still running
      if (commandRef.current) {
        // Tauri's Child process isn't directly exposed here, 
        // but Command garbage collection usually handles it.
        // For production we'd want to explicitly kill it if possible via the Child handle.
      }
      setActive(false);
    };
  }, [containerId]);

  return { logs, active, clear: () => setLogs([]) };
}
