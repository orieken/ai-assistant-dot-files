#!/usr/bin/env python3
"""Validate YAML frontmatter of markdown files against a JSON Schema.

Used by ``scripts/health-check.sh`` for the "Frontmatter JSON Schema
Validation" step. Kept as a separate file so the shell script stays
readable and this can be re-run standalone.

Requires ``jsonschema`` and ``PyYAML`` (both widely available via pip).
If either import fails, exit code 3 signals "validator unavailable" so
the calling shell can skip the check gracefully instead of failing.

Usage:
    validate-frontmatter.py <schema.json> <file1.md> [file2.md ...]

Exit codes:
    0 - every file passed validation
    1 - one or more files failed validation
    2 - bad invocation (missing args)
    3 - required Python library not installed (jsonschema or PyYAML)
"""

from __future__ import annotations

import datetime as _dt
import json
import re
import sys
from pathlib import Path

try:
    import jsonschema
    import yaml
except ImportError as exc:  # pragma: no cover — reported via exit code 3
    print(f"validate-frontmatter: missing dependency: {exc}", file=sys.stderr)
    sys.exit(3)


FRONTMATTER_RE = re.compile(r"^---\s*\n(.*?)\n---\s*\n", re.DOTALL)


def _stringify_dates(value):
    """Coerce YAML date/datetime nodes to ISO strings for JSON Schema
    string validation. Matches how yaml-language-server (the Red Hat VS
    Code extension) treats YAML nodes when validating against JSON
    Schemas — so the IDE experience and the shell check agree.
    """
    if isinstance(value, dict):
        return {k: _stringify_dates(v) for k, v in value.items()}
    if isinstance(value, list):
        return [_stringify_dates(v) for v in value]
    if isinstance(value, (_dt.date, _dt.datetime)):
        return value.isoformat()
    return value


def _extract_frontmatter(md_path: Path):
    text = md_path.read_text()
    match = FRONTMATTER_RE.match(text)
    if match is None:
        return None
    return _stringify_dates(yaml.safe_load(match.group(1)))


def _validate_one(validator, md_path: Path):
    fm = _extract_frontmatter(md_path)
    if fm is None:
        return "skip", []
    errors = sorted(validator.iter_errors(fm), key=lambda e: list(e.absolute_path))
    return ("fail" if errors else "pass"), errors


def main(argv):
    if len(argv) < 3:
        print(__doc__, file=sys.stderr)
        return 2

    schema = json.loads(Path(argv[1]).read_text())
    validator = jsonschema.Draft7Validator(schema)

    failed = 0
    for md_arg in argv[2:]:
        md_path = Path(md_arg)
        status, errors = _validate_one(validator, md_path)
        if status == "skip":
            print(f"  SKIP  {md_path} — no frontmatter block")
            continue
        if status == "fail":
            failed += 1
            print(f"  FAIL  {md_path}")
            for err in errors:
                loc = "/".join(str(p) for p in err.absolute_path) or "<root>"
                print(f"          - {loc}: {err.message}")
        else:
            print(f"  PASS  {md_path}")

    return 1 if failed else 0


if __name__ == "__main__":
    sys.exit(main(sys.argv))
