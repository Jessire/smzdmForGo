from __future__ import annotations

import argparse
from datetime import datetime
from pathlib import Path


ROOT = Path(__file__).resolve().parents[2]
TASKS = ROOT / ".trellis" / "tasks"


def create(title: str, slug: str) -> Path:
    prefix = datetime.now().strftime("%m-%d")
    name = f"{prefix}-{slug}"
    task_dir = TASKS / name
    task_dir.mkdir(parents=True, exist_ok=False)
    (task_dir / "status").write_text("planning\n", encoding="utf-8")
    (task_dir / "prd.md").write_text(
        f"# {title}\n\n"
        "## Goal\n\n"
        "- State the concrete user-facing outcome before implementation.\n\n"
        "## Confirmed Facts\n\n"
        "- Record repository-backed facts with file paths before implementation.\n\n"
        "## Requirements\n\n"
        "- Write each requirement so it can be verified.\n\n"
        "## Acceptance Criteria\n\n"
        "- [ ] Verification commands are recorded and pass.\n\n"
        "## Out Of Scope\n\n"
        "- Unrelated product behavior changes.\n",
        encoding="utf-8",
    )
    (TASKS / "current").write_text(name + "\n", encoding="utf-8")
    return task_dir


def start() -> None:
    current = TASKS / "current"
    if not current.exists():
        raise SystemExit("No active task")
    task_dir = TASKS / current.read_text(encoding="utf-8").strip()
    (task_dir / "status").write_text("in_progress\n", encoding="utf-8")


def complete() -> None:
    current = TASKS / "current"
    if not current.exists():
        raise SystemExit("No active task")
    task_dir = TASKS / current.read_text(encoding="utf-8").strip()
    (task_dir / "status").write_text("complete\n", encoding="utf-8")


def main() -> None:
    parser = argparse.ArgumentParser()
    sub = parser.add_subparsers(dest="cmd", required=True)
    create_parser = sub.add_parser("create")
    create_parser.add_argument("title")
    create_parser.add_argument("--slug", required=True)
    sub.add_parser("start")
    sub.add_parser("complete")
    args = parser.parse_args()

    if args.cmd == "create":
        print(create(args.title, args.slug))
    elif args.cmd == "start":
        start()
    elif args.cmd == "complete":
        complete()


if __name__ == "__main__":
    main()
