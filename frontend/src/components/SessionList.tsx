import { useEffect, useState } from "react";

export default function SessionList() {
  const [sessions, setSessions] = useState<string[]>([]);

  useEffect(() => {
    const fetchSessions = async () => {
      try {
        const res = await fetch("http://localhost:8080/sessions");
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
    <div className="p-2">
      <h2 className="text-lg font-semibold mb-2">Sessions</h2>
      <ul className="space-y-1">
        {sessions.map((id) => (
          <li key={id} className="text-sm truncate">
            {id}
          </li>
        ))}
      </ul>
    </div>
  );
}
