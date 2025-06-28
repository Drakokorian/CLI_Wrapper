import { useEffect, useState } from "react";

interface Props {
  selected: string;
  onChange: (model: string) => void;
}

export default function ModelSelector({ selected, onChange }: Props) {
  const [models, setModels] = useState<string[]>([]);

  useEffect(() => {
    const fetchModels = async () => {
      try {
        const res = await fetch("http://localhost:8080/models");
        if (res.ok) {
          const data = await res.json();
          setModels(data);
        }
      } catch (err) {
        console.error(err);
      }
    };
    fetchModels();
  }, []);

  return (
    <select
      className="border rounded px-2 py-1"
      value={selected}
      onChange={(e) => onChange(e.target.value)}
    >
      {models.map((m) => (
        <option key={m} value={m}>
          {m}
        </option>
      ))}
    </select>
  );
}
