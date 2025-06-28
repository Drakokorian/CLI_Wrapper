import { useState } from 'react';
import './App.css';

function App() {
  const [prompt, setPrompt] = useState('');
  const [response, setResponse] = useState('');
  const [tool, setTool] = useState('openai');

  const send = async () => {
    const res = await window.backend.RunPrompt(tool, prompt);
    setResponse(res);
  };

  return (
    <div>
      <select value={tool} onChange={e => setTool(e.target.value)}>
        <option value="openai">OpenAI</option>
        <option value="gemini">Gemini</option>
      </select>
      <textarea value={prompt} onChange={e => setPrompt(e.target.value)} />
      <button onClick={send}>Send</button>
      <pre>{response}</pre>
    </div>
  );
}

export default App;
