import { useEffect, useState } from "react";

export default function ResourceMeter() {
  const [cpu, setCpu] = useState(0);
  const [mem, setMem] = useState(0);

  useEffect(() => {
    const fetchUsage = async () => {
      try {
        const res = await fetch("http://localhost:8080/resource");
        if (res.ok) {
          const data = await res.json();
          setCpu(data.cpu);
          setMem(data.memory);
        }
      } catch (err) {
        console.error(err);
      }
    };
    fetchUsage();
    const id = setInterval(fetchUsage, 2000);
    return () => clearInterval(id);
  }, []);

  return (
    <div className="w-40">
      <div className="mb-2">
        <span className="text-sm">CPU {cpu.toFixed(1)}%</span>
        <div className="w-full bg-gray-200 h-2 rounded">
          <div
            className="bg-primary-600 h-2 rounded"
            style={{ width: `${cpu}%` }}
          />
        </div>
      </div>
      <div>
        <span className="text-sm">RAM {mem.toFixed(1)}%</span>
        <div className="w-full bg-gray-200 h-2 rounded">
          <div
            className="bg-success-600 h-2 rounded"
            style={{ width: `${mem}%` }}
          />
        </div>
      </div>
    </div>
  );
}
