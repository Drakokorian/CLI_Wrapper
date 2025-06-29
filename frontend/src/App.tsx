import { useState, useEffect } from "react";
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

function App() {
  const [prompt, setPrompt] = useState("");
  const [response, setResponse] = useState("");
  const [model, setModel] = useState("");
  const [alert, setAlert] = useState("");

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

  const send = async () => {
    const res = await window.backend.RunPrompt(model, prompt);
    setResponse(maskSecrets(res));
  };

  return (
    <div className="flex h-screen">
      <aside className="w-48 bg-gray-100 dark:bg-gray-800 overflow-y-auto">
        <SessionList />
        <HistoryViewer />
      </aside>
      <main className="flex-1 p-4 flex flex-col space-y-4">
        <div className="flex justify-between items-center">
          <ModelSelector selected={model} onChange={setModel} />
          <div className="flex items-center space-x-2">
            <ResourceMeter />
            <ThemeToggle />
          </div>
        </div>
        {alert && (
          <div className="text-yellow-600" role="alert">
            {alert}
          </div>
        )}
        <textarea
          className="border rounded p-2 flex-1"
          value={prompt}
          onChange={(e) => setPrompt(e.target.value)}
        />
        <button
          className="bg-blue-500 text-white px-4 py-2 rounded"
          onClick={send}
        >
          Send
        </button>
        <pre className="whitespace-pre-wrap bg-gray-50 p-2 rounded flex-1 overflow-y-auto">
          {response}
        </pre>
      </main>
    </div>
  );
}

export default App;
