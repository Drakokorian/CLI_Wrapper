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

  const highlightStyle = flash
    ? { borderColor: '#2563eb', backgroundColor: '#eff6ff' }
    : {};

  return (
    <select
      className="border rounded px-2 py-1 font-sans transition-colors"
      style={highlightStyle}
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
