import { useState, useRef, useEffect } from "react";
import { useContainers } from "./hooks/useContainers";
import { useImages } from "./hooks/useImages";
import { useMarketplace } from "./hooks/useMarketplace";
import { useStats } from "./hooks/useStats";
import { useLogs } from "./hooks/useLogs";
import { 
  Terminal, 
  Activity, 
  Layers, 
  Archive, 
  Database, 
  Globe, 
  Cpu, 
  Zap, 
  Package, 
  Binary, 
  Coffee, 
  Code, 
  Gem, 
  Layout, 
  Box, 
  Search,
  Filter
} from "lucide-react";

function App() {
  const [activeTab, setActiveTab] = useState("dashboard");
  const [selectedLogContainer, setSelectedLogContainer] = useState<{id: string, name: string} | null>(null);

  // Global stats polling
  const { stats } = useStats();

  return (
    <div className="flex h-screen bg-background text-text-main font-sans overflow-hidden">
      {/* Sidebar - Nocturnal Darker */}
      <aside className="w-64 bg-surface border-r border-border-subtle flex flex-col p-6 z-10">
        <div className="flex items-center gap-3 mb-10">
          <div className="w-8 h-8 rounded-lg bg-primary flex items-center justify-center shadow-[0_0_15px_rgba(0,245,255,0.3)]">
            <span className="text-background font-bold text-lg">M</span>
          </div>
          <h1 className="text-xl font-bold tracking-tight">MiniContainer</h1>
        </div>

        <nav className="flex-1 space-y-2">
          <NavItem 
            label="Dashboard" 
            active={activeTab === "dashboard"} 
            onClick={() => setActiveTab("dashboard")} 
          />
          <NavItem 
            label="Marketplace" 
            active={activeTab === "marketplace"} 
            onClick={() => setActiveTab("marketplace")} 
          />
          <NavItem 
            label="Image Library" 
            active={activeTab === "images"} 
            onClick={() => setActiveTab("images")} 
          />
        </nav>

        <div className="mt-auto pt-6 border-t border-border-subtle">
          <div className="flex items-center gap-3 text-sm text-text-dim">
            <div className="status-pulse"></div>
            <span>Runtime Active</span>
          </div>
        </div>
      </aside>

      {/* Main Content Area */}
      <main className="flex-1 overflow-y-auto relative p-8">
        {/* Background glow effects */}
        <div className="absolute top-[-10%] right-[-10%] w-[40%] h-[40%] bg-primary/5 blur-[120px] rounded-full pointer-events-none"></div>
        <div className="absolute bottom-[-5%] left-[-5%] w-[30%] h-[30%] bg-accent/5 blur-[100px] rounded-full pointer-events-none"></div>

        <header className="mb-10">
          <h2 className="text-3xl font-bold mb-2 lowercase opacity-90 tracking-tight">
            {activeTab === "dashboard" ? "Active Envs" : activeTab}
          </h2>
          <p className="text-text-dim max-w-lg">
            {activeTab === "dashboard" 
              ? "Your current fast-path development containers. Clean, isolated, and instant." 
              : "Expansive repositories for pre-configured development environments."}
          </p>
        </header>

        {/* Content Panel */}
        <div className="glass-panel p-6 min-h-[500px] relative">
          {activeTab === "dashboard" && <Dashboard stats={stats} onShowLogs={setSelectedLogContainer} />}
          {activeTab === "marketplace" && <Marketplace />}
          {activeTab === "images" && <ImageLibrary />}
        </div>

        {/* Log Viewer Modal */}
        {selectedLogContainer && (
          <LogViewer 
            container={selectedLogContainer} 
            onClose={() => setSelectedLogContainer(null)} 
          />
        )}
      </main>
    </div>
  );
}

function ErrorToast({ message, onClose }: { message: string, onClose: () => void }) {
  useEffect(() => {
    const timer = setTimeout(onClose, 8000);
    return () => clearTimeout(timer);
  }, [onClose]);

  return (
    <div className="fixed bottom-8 right-8 z-[100] max-w-sm animate-in fade-in slide-in-from-bottom-4 duration-300">
      <div className="bg-surface border-l-4 border-accent p-4 shadow-2xl rounded-r-md flex items-start gap-4">
        <div className="flex-1">
          <h5 className="text-accent font-bold text-xs uppercase tracking-widest mb-1">System Alert</h5>
          <p className="text-sm text-text-main leading-relaxed">{message}</p>
        </div>
        <button onClick={onClose} className="text-text-dim hover:text-white transition-colors">&times;</button>
      </div>
    </div>
  );
}

function NavItem({ label, active, onClick }: { label: string; active: boolean; onClick: () => void }) {
  return (
    <button
      onClick={onClick}
      className={`w-full text-left px-4 py-3 rounded-md transition-all duration-200 flex items-center group ${
        active 
          ? "bg-primary/10 text-primary font-medium" 
          : "text-text-dim hover:text-text-main hover:bg-surface/50"
      }`}
    >
      <span className={`w-1.5 h-1.5 rounded-full mr-3 transition-all duration-300 ${
        active ? "bg-primary scale-100 shadow-[0_0_8px_#00F5FF]" : "bg-transparent scale-0 group-hover:scale-100 group-hover:bg-text-dim/30"
      }`}></span>
      {label}
    </button>
  );
}

function Dashboard({ stats, onShowLogs }: { stats: any, onShowLogs: (c: {id: string, name: string}) => void }) {
  const { containers, loading, error, refreshAction, startContainer, stopContainer, removeContainer, pauseContainer, unpauseContainer, restartContainer, execContainer, clearError } = useContainers();

  if (loading && containers.length === 0) return <div className="flex items-center justify-center h-64 text-text-dim">Loading node state...</div>;

  return (
    <div className="space-y-6">
      {error && <ErrorToast message={error} onClose={clearError} />}
      <div className="flex items-center justify-between">
        <h3 className="text-lg font-semibold tracking-tight text-text-main/90">Container Instance List</h3>
        <button onClick={refreshAction} className="text-text-dim hover:text-primary transition-colors text-sm font-mono uppercase tracking-widest">
          [ Sync ]
        </button>
      </div>

      <div className="overflow-x-auto">
        <table className="w-full text-left border-collapse">
          <thead>
            <tr className="border-b border-border-subtle">
              <th className="pb-4 font-mono text-[10px] uppercase tracking-widest text-text-dim">ID</th>
              <th className="pb-4 font-mono text-[10px] uppercase tracking-widest text-text-dim">NAME</th>
              <th className="pb-4 font-mono text-[10px] uppercase tracking-widest text-text-dim">CPU / MEM</th>
              <th className="pb-4 font-mono text-[10px] uppercase tracking-widest text-text-dim">STATUS</th>
              <th className="pb-4 font-mono text-[10px] uppercase tracking-widest text-text-dim text-right">ACTIONS</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-border-subtle/30">
            {containers.map((container) => {
              const s = stats[container.id.slice(0, 12)];
              const statusLower = container.status.toLowerCase();
              const isRunning = (statusLower.includes("run") || statusLower.includes("up")) && !statusLower.includes("paused");
              const isPaused = statusLower.includes("paused");

              return (
                <tr key={container.id} className="group hover:bg-white/[0.02] transition-colors">
                  <td className="py-4 font-mono text-xs text-text-dim">{container.id.slice(0, 8)}</td>
                  <td className="py-4 font-medium">{container.names}</td>
                  <td className="py-4">
                    {s ? (
                      <div className="flex gap-2">
                        <span className="text-[10px] font-mono bg-primary/5 text-primary border border-primary/20 px-1.5 rounded">
                          CPU {s.cpu_perc}
                        </span>
                        <span className="text-[10px] font-mono bg-accent/5 text-accent border border-accent/20 px-1.5 rounded">
                          RAM {s.mem_perc}
                        </span>
                      </div>
                    ) : (
                      <span className="text-[10px] font-mono text-text-dim">--</span>
                    )}
                  </td>
                  <td className="py-4 text-sm">
                    <span className={`flex items-center gap-2 ${
                      isRunning ? "text-primary" : (isPaused ? "text-yellow-500" : "text-text-dim")
                    }`}>
                      <span className={`w-1.5 h-1.5 rounded-full ${
                        isRunning ? "bg-primary animate-pulse shadow-[0_0_8px_#00F5FF]" : (isPaused ? "bg-yellow-500 shadow-[0_0_8px_#EAB308]" : "bg-text-dim/40")
                      }`}></span>
                      {container.status}
                    </span>
                  </td>
                  <td className="py-4 text-right">
                    <div className="flex items-center justify-end gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
                      <ActionBtn label="Logs" onClick={() => onShowLogs({ id: container.id, name: container.names })} />
                      
                      {isRunning ? (
                        <>
                          <ActionBtn label="Terminal" onClick={() => execContainer(container.id)} />
                          <ActionBtn label="Restart" onClick={() => restartContainer(container.id)} />
                          <ActionBtn label="Pause" onClick={() => pauseContainer(container.id)} />
                          <ActionBtn label="Stop" onClick={() => stopContainer(container.names)} />
                        </>
                      ) : isPaused ? (
                        <>
                          <ActionBtn label="Unpause" onClick={() => unpauseContainer(container.id)} />
                          <ActionBtn label="Stop" onClick={() => stopContainer(container.names)} />
                        </>
                      ) : (
                        <ActionBtn label="Start" onClick={() => startContainer(container.names)} />
                      )}
                      
                      <ActionBtn label="Remove" danger onClick={() => removeContainer(container.names)} />
                    </div>
                  </td>
                </tr>
              );
            })}
          </tbody>
        </table>
      </div>

      {containers.length === 0 && (
        <div className="text-center py-20 border-2 border-dashed border-border-subtle rounded-lg">
          <p className="text-text-dim mb-4 italic">No active containers found in cluster.</p>
          <button className="btn-primary text-sm">Initialize New Environment</button>
        </div>
      )}
    </div>
  );
}

function ImageLibrary() {
  const { images, loading, error, refreshAction, removeImage, clearError } = useImages();

  if (loading && images.length === 0) return <div className="flex items-center justify-center h-64 text-text-dim">Scanning local registry...</div>;

  return (
    <div className="space-y-6">
      {error && <ErrorToast message={error} onClose={clearError} />}
      <div className="flex items-center justify-between">
        <h3 className="text-lg font-semibold tracking-tight text-text-main/90">Local Image Cache</h3>
        <button onClick={refreshAction} className="text-text-dim hover:text-primary transition-colors text-sm font-mono uppercase tracking-widest">
          [ Re-index ]
        </button>
      </div>

      <div className="overflow-x-auto">
        <table className="w-full text-left border-collapse">
          <thead>
            <tr className="border-b border-border-subtle">
              <th className="pb-4 font-mono text-[10px] uppercase tracking-widest text-text-dim">REPOSITORY</th>
              <th className="pb-4 font-mono text-[10px] uppercase tracking-widest text-text-dim">TAG</th>
              <th className="pb-4 font-mono text-[10px] uppercase tracking-widest text-text-dim">ID</th>
              <th className="pb-4 font-mono text-[10px] uppercase tracking-widest text-text-dim">SIZE</th>
              <th className="pb-4 font-mono text-[10px] uppercase tracking-widest text-text-dim text-right">ACTIONS</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-border-subtle/30">
            {images.map((img) => (
              <tr key={img.id} className="group hover:bg-white/[0.02] transition-colors">
                <td className="py-4 font-medium">{img.repository}</td>
                <td className="py-4 text-sm text-text-dim">{img.tag}</td>
                <td className="py-4 font-mono text-xs text-text-dim">{img.id.slice(0, 12)}</td>
                <td className="py-4 text-sm text-text-dim">{(img.size / (1024 * 1024)).toFixed(1)} MB</td>
                <td className="py-4 text-right">
                  <div className="flex items-center justify-end gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
                    <ActionBtn label="Prune" danger onClick={() => removeImage(img.id)} />
                  </div>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {images.length === 0 && (
        <div className="text-center py-20 border-2 border-dashed border-border-subtle rounded-lg">
          <p className="text-text-dim">No local images detected in podman storage.</p>
        </div>
      )}
    </div>
  );
}

function Marketplace() {
  const { presets, loading, error, deployAction, clearError } = useMarketplace();
  const [deploying, setDeploying] = useState<string | null>(null);
  const [selectedCategory, setSelectedCategory] = useState<string>("All");

  const categories = ["All", "Develop", "Database", "Service", "Tools"];

  if (loading && Object.keys(presets).length === 0) return <div className="flex items-center justify-center h-64 text-text-dim">Loading index...</div>;

  const handleDeploy = async (id: string) => {
    setDeploying(id);
    try {
      await deployAction(id);
    } finally {
      setDeploying(null);
    }
  };

  const filteredPresets = Object.entries(presets).filter(([_, p]) => 
    selectedCategory === "All" || p.category === selectedCategory
  );

  const IconMap: Record<string, any> = {
    python: Code,
    package: Package,
    binary: Binary,
    cpu: Cpu,
    database: Database,
    box: Box,
    zap: Zap,
    globe: Globe,
    coffee: Coffee,
    code: Code,
    gem: Gem,
    layout: Layout,
  };

  return (
    <div className="space-y-8">
      <div className="flex items-center justify-between">
        <div className="flex bg-surface/50 p-1 rounded-lg border border-border-subtle">
          {categories.map(cat => (
            <button
              key={cat}
              onClick={() => setSelectedCategory(cat)}
              className={`px-4 py-1.5 rounded-md text-xs font-medium transition-all ${
                selectedCategory === cat 
                  ? "bg-primary text-background shadow-[0_0_10px_rgba(0,245,255,0.4)]" 
                  : "text-text-dim hover:text-text-main hover:bg-white/5"
              }`}
            >
              {cat}
            </button>
          ))}
        </div>
        <div className="relative group">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-text-dim group-focus-within:text-primary transition-colors" />
          <input 
            type="text" 
            placeholder="Search templates..." 
            className="bg-surface/50 border border-border-subtle rounded-lg pl-10 pr-4 py-1.5 text-xs focus:border-primary/50 outline-none transition-all w-64"
          />
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 relative">
        {error && <div className="col-span-full"><ErrorToast message={error} onClose={clearError} /></div>}
        {filteredPresets.map(([key, preset]) => {
          const IconObj = IconMap[preset.icon || "box"] || Box;
          return (
            <div key={key} className="bg-surface/50 border border-border-subtle p-6 rounded-lg hover:border-primary/30 transition-all duration-300 group flex flex-col h-full">
              <div className="flex items-center justify-between mb-4">
                <div className="flex items-center gap-3">
                  <div className="w-10 h-10 rounded-lg bg-white/5 flex items-center justify-center group-hover:bg-primary/10 transition-colors">
                    <IconObj className="w-5 h-5 text-text-dim group-hover:text-primary transition-colors" />
                  </div>
                  <h4 className="text-lg font-bold group-hover:text-primary transition-colors capitalize">{key}</h4>
                </div>
                <span className="font-mono text-[10px] text-text-dim bg-white/5 px-2 py-0.5 rounded uppercase">{preset.category || "Misc"}</span>
              </div>
              
              <p className="text-sm text-text-dim mb-6 italic flex-grow">"{preset.description}"</p>
              
              <div className="space-y-3 mb-8">
                <div className="flex justify-between text-[11px] font-mono">
                  <span className="text-text-dim/60">SOURCE</span>
                  <span className="truncate max-w-[150px]">{preset.image.split('/').pop()}</span>
                </div>
                <div className="flex justify-between text-[11px] font-mono">
                  <span className="text-text-dim/60">NETWORK</span>
                  <span>{preset.ports && preset.ports.length > 0 ? preset.ports.join(', ') : "Internal Only"}</span>
                </div>
              </div>

              <button 
                disabled={!!deploying}
                onClick={() => handleDeploy(key)}
                className="w-full py-3 bg-white/5 border border-white/10 rounded font-bold text-xs uppercase tracking-widest hover:bg-primary hover:text-background transition-all duration-200 active:scale-95 disabled:opacity-50"
              >
                {deploying === key ? "Initializing..." : `Launch ${key}`}
              </button>
            </div>
          );
        })}
      </div>
    </div>
  );
}

function LogViewer({ container, onClose }: { container: {id: string, name: string}, onClose: () => void }) {
  const { logs, active } = useLogs(container.id);
  const logEndRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    logEndRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [logs]);

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-background/80 backdrop-blur-sm">
      <div className="w-full max-w-4xl h-[80vh] bg-surface border border-border-subtle rounded-xl shadow-2xl flex flex-col overflow-hidden relative">
        {/* Header */}
        <div className="p-4 border-b border-border-subtle flex items-center justify-between bg-white/[0.02]">
          <div className="flex items-center gap-3">
            <div className="w-2 h-2 rounded-full bg-primary animate-pulse"></div>
            <h4 className="font-bold text-sm tracking-tight capitalize">
              Terminal: <span className="text-primary">{container.name}</span>
            </h4>
            {!active && <span className="text-[10px] uppercase font-mono text-accent bg-accent/10 px-2 rounded">Disconnected</span>}
          </div>
          <button 
            onClick={onClose}
            className="text-text-dim hover:text-white transition-colors text-xl leading-none"
          >
            &times;
          </button>
        </div>

        {/* Content */}
        <div className="flex-1 bg-black/40 p-6 overflow-y-auto font-mono text-sm leading-relaxed scrollbar-thin scrollbar-thumb-white/10">
          {logs.length === 0 ? (
            <div className="text-text-dim italic opacity-50">Waiting for output line...</div>
          ) : (
            logs.map((log, i) => (
              <div key={i} className="mb-0.5 whitespace-pre-wrap break-all opacity-90 hover:opacity-100">
                <span className="text-[10px] text-text-dim/30 mr-4 select-none">{(i + 1).toString().padStart(3, '0')}</span>
                {log}
              </div>
            ))
          )}
          <div ref={logEndRef} />
        </div>

        {/* Footer */}
        <div className="p-3 bg-white/[0.01] border-t border-border-subtle flex justify-between items-center px-6">
          <span className="text-[10px] font-mono text-text-dim/60 uppercase">Streaming active session</span>
          <span className="text-[10px] font-mono text-text-dim/60">{logs.length} lines buffered</span>
        </div>
      </div>
    </div>
  );
}



function ActionBtn({ label, onClick, danger }: { label: string; onClick: () => void; danger?: boolean }) {
  return (
    <button
      onClick={onClick}
      className={`px-3 py-1 rounded text-[10px] font-mono uppercase tracking-tighter transition-all duration-200 border ${
        danger 
          ? "border-accent/40 text-accent/80 hover:bg-accent hover:text-white" 
          : "border-primary/40 text-primary/80 hover:bg-primary hover:text-background"
      }`}
    >
      {label}
    </button>
  );
}

export default App;
