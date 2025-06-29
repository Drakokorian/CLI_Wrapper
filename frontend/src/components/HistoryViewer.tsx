import { useEffect, useState, ChangeEvent } from "react";
import { maskSecrets } from "../utils";

interface Record {
  id: string;
  prompt: string;
  response: string;
  model: string;
  timestamp: string;
  success: boolean;
}

export default function HistoryViewer() {
  const [query, setQuery] = useState("");
  const [records, setRecords] = useState<Record[]>([]);

  const search = async () => {
    const res = await fetch(`/history/search?q=${encodeURIComponent(query)}`);
    if (res.ok) {
      setRecords(await res.json());
    }
  };

  const exportHist = async () => {
    const res = await fetch("/history/export");
    if (res.ok) {
      const blob = await res.blob();
      const url = URL.createObjectURL(blob);
      const a = document.createElement("a");
      a.href = url;
      a.download = "history.json";
      a.click();
      URL.revokeObjectURL(url);
    }
  };

  const importHist = async (e: ChangeEvent<HTMLInputElement>) => {
    if (!e.target.files?.length) return;
    const text = await e.target.files[0].text();
    await fetch("/history", { method: "POST", headers: { "Content-Type": "application/json" }, body: text });
    e.target.value = "";
    search();
  };

  useEffect(() => {
    search();
  }, []);

  return (
    <div className="p-2 border-t mt-2">
      <h2 className="text-lg font-semibold mb-2">History</h2>
      <div className="flex mb-2">
        <input className="border p-1 flex-1" value={query} onChange={(e) => setQuery(e.target.value)} />
        <button className="bg-blue-500 text-white px-2 ml-2" onClick={search}>Search</button>
      </div>
      <div className="flex items-center space-x-2 mb-2">
        <button className="bg-green-500 text-white px-2" onClick={exportHist}>Export</button>
        <label className="bg-gray-200 p-1 cursor-pointer">
          Import<input type="file" className="hidden" onChange={importHist} />
        </label>
      </div>
      <ul className="space-y-1 overflow-y-auto" style={{ maxHeight: "200px" }}>
        {records.map((r) => (
          <li key={r.id} className="text-sm">
            <div className="font-mono">{maskSecrets(r.prompt)}</div>
            <div className="text-gray-500 truncate">{maskSecrets(r.response)}</div>
          </li>
        ))}
      </ul>
    </div>
  );
}
