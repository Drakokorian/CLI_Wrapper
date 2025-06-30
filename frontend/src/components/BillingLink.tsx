import { useEffect, useState } from "react";
import { OpenURL } from "@wailsio/runtime";

interface UsageTotals {
  total_granted: number;
  total_used: number;
  total_available: number;
}

interface BillingResponse {
  tool: string;
  url: string;
  usage?: UsageTotals;
}

export default function BillingLink() {
  const [info, setInfo] = useState<BillingResponse | null>(null);

  useEffect(() => {
    const load = async () => {
      try {
        const res = await fetch("/billing");
        if (res.ok) {
          const data: BillingResponse = await res.json();
          setInfo(data);
        }
      } catch (err) {
        console.error(err);
      }
    };
    load();
  }, []);

  if (!info) return null;

  const handleOpen = () => {
    OpenURL(info.url);
  };

  return (
    <div className="flex items-center space-x-2 font-sans">
      {info.usage && (
        <span className="text-sm" aria-label="usage">
          Used {info.usage.total_used.toFixed(2)} / {info.usage.total_granted.toFixed(2)}
        </span>
      )}
      <button className="border rounded px-2 py-1 bg-primary-600 text-white" onClick={handleOpen}>
        View Billing
      </button>
    </div>
  );
}
