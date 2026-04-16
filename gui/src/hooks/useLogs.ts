import { useState, useEffect, useRef } from "react";
import { Command } from "@tauri-apps/plugin-shell";

export function useLogs(containerId: string | null) {
  const [logs, setLogs] = useState<string[]>([]);
  const [active, setActive] = useState(false);
  const commandRef = useRef<any>(null);
  const childRef = useRef<any>(null);

  useEffect(() => {
    if (!containerId) {
      setLogs([]);
      setActive(false);
      return;
    }

    setLogs([]);
    if (!(window as any).__TAURI_INTERNALS__) {
      setActive(false);
      return;
    }
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

    command.spawn().then((child) => {
      childRef.current = child;
    }).catch((err) => {
      setLogs((prev) => [...prev, `\x1b[31mError spawning logs: ${err.message}\x1b[0m`]);
      setActive(false);
    });

    return () => {
      if (childRef.current) {
        childRef.current.kill();
      }
      setActive(false);
    };
  }, [containerId]);

  return { logs, active, clear: () => setLogs([]) };
}
