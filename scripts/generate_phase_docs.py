import os
import re
import json
from datetime import datetime, timezone

ROADMAP_FILE = os.path.join('docs', 'ROADMAP.md')
PHASE_DIR = os.path.join('docs', 'phases')

pattern = re.compile(r"\*\*Phase\s+(\d+):\s+([^*]+)\*\*")
phase_titles = {}
with open(ROADMAP_FILE, 'r', encoding='utf-8') as f:
    for line in f:
        m = pattern.search(line)
        if m:
            num = int(m.group(1))
            title = m.group(2).strip()
            phase_titles[num] = title

START = 9
END = 18

os.makedirs(PHASE_DIR, exist_ok=True)

for num in range(START, END+1):
    title = phase_titles.get(num, f"Phase {num}")
    meta = {
        'Phase': f"{num:03d}",
        'Title': title,
        'Created': datetime.now(timezone.utc).strftime('%Y-%m-%d'),
        'CompatibleWith': [
            'Windows 10+ (PowerShell 5.1+)',
            'Ubuntu 18.04+ (Bash 3.2+)',
            'macOS 10.15+ (zsh/bash)'
        ],
        'SprintWindow': 'Sprint 1 (Day 1–10)',
        'Author': 'AutoDocSystem'
    }

    tasks = []
    for t in range(1, 51):
        doc_cmd = f"echo {title.lower().replace(' ', '_')}_task_{t}"
        task = {
            'TaskID': f"PH{num:03d}-T{t:02d}",
            'TaskTitle': f"Task {t} for {title}",
            'TaskDetails': f"Detailed instructions for task {t} in {title}.",
            'Platform Compatibility Notes': 'PowerShell 5.1+, Bash 3.2+, zsh compatible.',
            'Security Considerations': 'Validate inputs and restrict permissions.',
            'Estimated Sprint Day': f"Day {((t-1)//5)+1}",
            'Documentation Block': f"```bash\n{doc_cmd}\n```"
        }
        tasks.append(task)

    base_name = f"Phase_{num:03d}_{title.replace(' ', '_').replace('/', '_')}"
    md_path = os.path.join(PHASE_DIR, base_name + '.md')
    json_path = os.path.join(PHASE_DIR, base_name + '.json')
    txt_path = os.path.join(PHASE_DIR, base_name + '.txt')

    with open(md_path, 'w', encoding='utf-8') as md:
        md.write('---\n')
        for k, v in meta.items():
            if k == 'CompatibleWith':
                md.write(f"{k}:\n")
                for item in v:
                    md.write(f"  - {item}\n")
            else:
                md.write(f"{k}: {v}\n")
        md.write('---\n\n')
        md.write(f"# Phase {num:03d}: {title}\n\n")
        md.write(f"## Objective\nThis phase focuses on {title.lower()} across all supported platforms.\n\n")
        md.write("## Tasks\n")
        for task in tasks:
            md.write(f"- **{task['TaskID']}**: {task['TaskTitle']}\n")
            md.write(f"  - TaskDetails: {task['TaskDetails']}\n")
            md.write(f"  - Platform Compatibility Notes: {task['Platform Compatibility Notes']}\n")
            md.write(f"  - Security Considerations: {task['Security Considerations']}\n")
            md.write(f"  - Estimated Sprint Day: {task['Estimated Sprint Day']}\n")
            md.write(f"  - Documentation Block:\n    {task['Documentation Block']}\n\n")

    with open(json_path, 'w', encoding='utf-8') as js:
        json.dump({'metadata': meta, 'tasks': tasks}, js, indent=2)

    with open(txt_path, 'w', encoding='utf-8') as tx:
        tx.write(f"Phase {num:03d}: {title}\n")
        tx.write(f"Created: {meta['Created']}\n")
        tx.write("Tasks:\n")
        for task in tasks:
            tx.write(f"{task['TaskID']}: {task['TaskTitle']} - {task['Estimated Sprint Day']}\n")
