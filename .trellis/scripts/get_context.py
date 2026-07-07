from __future__ import annotations

import argparse
import json
import os
from pathlib import Path


ROOT = Path(__file__).resolve().parents[2]
TRELLIS = ROOT / ".trellis"
TASKS = TRELLIS / "tasks"


def current_task() -> dict[str, str] | None:
    current = TASKS / "current"
    if not current.exists():
        return None
    target = current.read_text(encoding="utf-8").strip()
    if not target:
        return None
    task_dir = TASKS / target
    status = "unknown"
    status_file = task_dir / "status"
    if status_file.exists():
        status = status_file.read_text(encoding="utf-8").strip()
    return {"slug": target, "status": status, "path": str(task_dir.relative_to(ROOT))}


def mode_default() -> None:
    print(f"Repository: {ROOT.name}")
    print(f"Trellis root: {TRELLIS.relative_to(ROOT)}")
    task = current_task()
    if task:
        print(f"Active task: {task['slug']}")
        print(f"Status: {task['status']}")
        print(f"Task path: {task['path']}")
    else:
        print("Active task: none")
    print("Journal: .trellis/journal.md")


def mode_phase(step: str | None) -> None:
    if step == "2.1":
        print("Phase 2.1: Before development")
        print("- Read active task artifacts.")
        print("- Read relevant .trellis/spec indexes.")
        print("- Run implementation and verification commands from the task.")
        return
    print("Phase Index")
    print("1. Planning: prd.md, optional design.md and implement.md.")
    print("2. Development: read specs before editing.")
    print("3. Verification: run tests, browser checks, and deployment checks when relevant.")
    print("4. Completion: commit intentional files and report residual risks.")


def mode_packages() -> None:
    packages = {
        "app/backend": ["main.go", "route*.go", "db/", "smzdm/", "push/", "file/"],
        "app/frontend": ["template/html/index.html"],
        "app/deploy": ["Dockerfile", "render.yaml", "readme.md"],
    }
    print(json.dumps(packages, ensure_ascii=False, indent=2))


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("--mode", choices=["phase", "packages"], default="default")
    parser.add_argument("--step")
    parser.add_argument("--platform")
    args = parser.parse_args()

    if args.mode == "phase":
        mode_phase(args.step)
    elif args.mode == "packages":
        mode_packages()
    else:
        mode_default()


if __name__ == "__main__":
    main()
