import { useEffect, useState, useRef } from "react";

interface Props {
  selected: string;
  onChange: (model: string) => void;
}

export default function ModelSelector({ selected, onChange }: Props) {
  const [models, setModels] = useState<string[]>([]);
  const prev = useRef<string>(selected);
  const [flash, setFlash] = useState(false);

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

  useEffect(() => {
    if (prev.current !== selected) {
      setFlash(true);
      const id = setTimeout(() => setFlash(false), 300);
      prev.current = selected;
      return () => clearTimeout(id);
    }
  }, [selected]);

  return (
    <select
      className={`border rounded px-2 py-1 transition-colors ${flash ? 'border-blue-500 bg-blue-50' : ''}`}
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
