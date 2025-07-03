import { useState, useEffect, useRef } from "react";
import { EventsOn } from "@wailsio/runtime";
import { maskSecrets } from "./utils";

declare global {
  interface Window {
    backend: {
      RunPrompt(model: string, prompt: string): Promise<string>;
    };
  }
}
import "./App.css";
import SessionList from "./components/SessionList";
import ModelSelector from "./components/ModelSelector";
import ResourceMeter from "./components/ResourceMeter";
import ThemeToggle from "./components/ThemeToggle";
import HistoryViewer from "./components/HistoryViewer";
import BillingLink from "./components/BillingLink";

function App() {
  const [prompt, setPrompt] = useState("");
  const [response, setResponse] = useState("");
  const queueRef = useRef<string[]>([]);
  const [model, setModel] = useState("");
  const [alert, setAlert] = useState("");
  const [concurrency, setConcurrency] = useState(1);
  const [workingDir, setWorkingDir] = useState("");
  const [logLevel, setLogLevel] = useState("info");
  const [logPath, setLogPath] = useState("");
  const [showSettings, setShowSettings] = useState(false);

  useEffect(() => {
    const load = async () => {
      try {
        const res = await fetch("http://localhost:8080/config");
        if (res.ok) {
          const data = await res.json();
          setConcurrency(data.concurrency);
          setWorkingDir(data.workingDir);
          if (data.logLevel) setLogLevel(data.logLevel);
          if (data.logPath) setLogPath(data.logPath);
        }
      } catch (err) {
        console.error(err);
      }
    };
    load();
  }, []);

  useEffect(() => {
    const unsub = EventsOn("model:switched", (data: any) => {
      if (data && data.from && data.to) {
        setAlert(`\u26A0\uFE0F Model switched from ${data.from} to ${data.to}`);
        setTimeout(() => setAlert(""), 4000);
      }
    });
    return () => {
      unsub();
    };
  }, []);

  useEffect(() => {
    const reduce = window.matchMedia("(prefers-reduced-motion: reduce)").matches;
    const interval = setInterval(() => {
      if (queueRef.current.length === 0) {
        return;
      }
      if (reduce) {
        setResponse((p) => p + queueRef.current.join(""));
        queueRef.current = [];
        return;
      }
      const char = queueRef.current.shift();
      if (char) {
        setResponse((p) => p + char);
      }
    }, 30);
    return () => clearInterval(interval);
  }, []);

  const send = async () => {
    const id = await window.backend.RunPrompt(model, prompt);
    setResponse("");
    queueRef.current = [];
    const ws = new WebSocket(`ws://localhost:8080/streamchars?id=${id}`);
    ws.onmessage = (e) => {
      queueRef.current.push(...(maskSecrets(e.data)));
    };
  };

  const saveSettings = async () => {
    try {
      await fetch("http://localhost:8080/config", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          concurrency,
          workingDir,
          logLevel,
          logPath,
        }),
      });
      setShowSettings(false);
    } catch (err) {
      console.error(err);
    }
  };

  return (
    <div className="flex h-screen transition-colors duration-300 font-sans">
      <aside className="w-48 bg-gray-100 dark:bg-gray-800 overflow-y-auto">
        <SessionList />
        <HistoryViewer />
      </aside>
      <main className="flex-1 p-4 flex flex-col space-y-4">
        <div className="flex justify-between items-center">
          <ModelSelector selected={model} onChange={setModel} />
          <div className="flex items-center space-x-2">
            <ResourceMeter />
            <BillingLink />
            <ThemeToggle />
            <button
              className="border rounded px-2 py-1"
              onClick={() => setShowSettings(true)}
            >
              Settings
            </button>
          </div>
        </div>
        {alert && (
          <div className="text-yellow-600 transition-opacity duration-300" role="alert">
            {alert}
          </div>
        )}
        {showSettings && (
          <div className="border p-2 rounded space-y-2 transition-all duration-300" role="dialog">
            <div>
              <label className="block mb-1">Concurrent Agents</label>
              <input
                type="number"
                min={1}
                max={5}
                value={concurrency}
                onChange={(e) => setConcurrency(Number(e.target.value))}
                className="border p-1 rounded w-full"
              />
            </div>
            <div>
              <label className="block mb-1">Working Directory</label>
              <input
                type="text"
                value={workingDir}
                onChange={(e) => setWorkingDir(e.target.value)}
                className="border p-1 rounded w-full"
              />
            </div>
            <div>
              <label className="block mb-1">Log Level</label>
              <select
                value={logLevel}
                onChange={(e) => setLogLevel(e.target.value)}
                className="border p-1 rounded w-full"
              >
                <option value="error">error</option>
                <option value="info">info</option>
                <option value="debug">debug</option>
              </select>
            </div>
            <div>
              <label className="block mb-1">Log Path</label>
              <input
                type="text"
                value={logPath}
                onChange={(e) => setLogPath(e.target.value)}
                className="border p-1 rounded w-full"
              />
            </div>
            <div className="flex space-x-2">
              <button
                className="bg-primary-600 text-white px-3 py-1 rounded"
                onClick={saveSettings}
              >
                Save
              </button>
              <button
                className="border px-3 py-1 rounded"
                onClick={() => setShowSettings(false)}
              >
                Cancel
              </button>
            </div>
          </div>
        )}
        <textarea
          className="border rounded p-2 flex-1"
          value={prompt}
          onChange={(e) => setPrompt(e.target.value)}
        />
        <button
          className="bg-primary-600 text-white px-4 py-2 rounded"
          onClick={send}
        >
          Send
        </button>
        <pre
          className="whitespace-pre-wrap bg-gray-50 dark:bg-gray-900 text-gray-800 dark:text-gray-100 p-2 rounded flex-1 overflow-y-auto"
          aria-live="polite"
        >
          {response}
        </pre>
      </main>
    </div>
  );
}

export default App;
