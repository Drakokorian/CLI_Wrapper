import { useEffect, useState } from "react";

function SessionItem({ id }: { id: string }) {
  const [lines, setLines] = useState<string[]>([]);
  useEffect(() => {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const ws = new WebSocket(`${protocol}//${window.location.host}/stream?id=${id}`);
    ws.onmessage = (e) => setLines((p) => [...p, e.data as string]);
    return () => ws.close();
  }, [id]);
  return (
    <li className="text-sm">
      <div className="truncate font-mono">{id}</div>
      <pre className="text-xs whitespace-pre-wrap">{lines.join("\n")}</pre>
    </li>
  );
}

export default function SessionList() {
  const [sessions, setSessions] = useState<string[]>([]);

  useEffect(() => {
    const fetchSessions = async () => {
      try {
        const res = await fetch("/sessions");
        if (res.ok) {
          const data = await res.json();
          setSessions(data);
        }
      } catch (err) {
        console.error(err);
      }
    };
    fetchSessions();
    const id = setInterval(fetchSessions, 1000);
    return () => clearInterval(id);
  }, []);

  return (
    <div className="p-2 font-sans">
      <h2 className="text-lg font-semibold mb-2">Sessions</h2>
      <ul className="space-y-1">
        {sessions.map((id) => (
          <SessionItem key={id} id={id} />
        ))}
      </ul>
    </div>
  );
}
