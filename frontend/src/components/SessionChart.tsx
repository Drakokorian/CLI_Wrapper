import { useEffect, useState } from "react";

interface Metric {
  id: string;
  timestamp: string;
  cpu: number;
  memory: number;
}

export default function SessionChart({ id }: { id: string }) {
  const [metrics, setMetrics] = useState<Metric[]>([]);

  useEffect(() => {
    const fetchMetrics = async () => {
      const res = await fetch(`/resource/session/${encodeURIComponent(id)}`);
      if (res.ok) {
        setMetrics(await res.json());
      }
    };
    fetchMetrics();
  }, [id]);

  const width = 200;
  const height = 80;
  const points = (key: "cpu" | "memory") =>
    metrics
      .map((m, i) => {
        const x = (i / Math.max(metrics.length - 1, 1)) * width;
        const y = height - ((m[key] as number) / 100) * height;
        return `${x},${y}`;
      })
      .join(" ");

  return (
    <svg width={width} height={height} className="border" role="img">
      {metrics.length > 0 && (
        <>
          <polyline
            fill="none"
            stroke="blue"
            strokeWidth="2"
            points={points("cpu")}
          />
          <polyline
            fill="none"
            stroke="green"
            strokeWidth="2"
            points={points("memory")}
          />
        </>
      )}
    </svg>
  );
}
