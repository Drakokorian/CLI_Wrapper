import { useState } from "react";
import "./App.css";
import SessionList from "./components/SessionList";
import ModelSelector from "./components/ModelSelector";
import ResourceMeter from "./components/ResourceMeter";

function App() {
  const [prompt, setPrompt] = useState("");
  const [response, setResponse] = useState("");
  const [model, setModel] = useState("");

  const send = async () => {
    const res = await window.backend.RunPrompt(model, prompt);
    setResponse(res);
  };

  return (
    <div className="flex h-screen">
      <aside className="w-48 bg-gray-100 dark:bg-gray-800 overflow-y-auto">
        <SessionList />
      </aside>
      <main className="flex-1 p-4 flex flex-col space-y-4">
        <div className="flex justify-between items-center">
          <ModelSelector selected={model} onChange={setModel} />
          <ResourceMeter />
        </div>
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
