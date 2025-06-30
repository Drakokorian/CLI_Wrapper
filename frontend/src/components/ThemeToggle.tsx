import { useEffect, useState } from "react";

export default function ThemeToggle() {
  const [theme, setTheme] = useState<"light" | "dark">("light");

  useEffect(() => {
    const load = async () => {
      try {
        const res = await fetch("http://localhost:8080/theme");
        if (res.ok) {
          const data = await res.json();
          if (data.theme === "light" || data.theme === "dark") {
            setTheme(data.theme);
          }
        }
      } catch (err) {
        console.error(err);
      }
    };
    load();
  }, []);

  useEffect(() => {
    document.documentElement.classList.toggle("dark", theme === "dark");
  }, [theme]);

  const toggle = async () => {
    const next = theme === "light" ? "dark" : "light";
    setTheme(next);
    try {
      await fetch("http://localhost:8080/theme", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ theme: next }),
      });
    } catch (err) {
      console.error(err);
    }
  };

  return (
    <button className="border rounded px-2 py-1 font-sans transition-colors" onClick={toggle} aria-label="Toggle theme">
      {theme === "light" ? "Dark Mode" : "Light Mode"}
    </button>
  );
}
