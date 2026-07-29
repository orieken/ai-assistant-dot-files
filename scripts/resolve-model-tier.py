#!/usr/bin/env python3
"""
Resolves portable `model_tier` declarations in agent frontmatter against
target platform mappings in shared/model-defaults.yaml (or user overrides).

Emits WARNING logs for platforms that do not honor per-agent model selection (null mapping).
"""

import sys
import os
import re
import argparse

PLATFORM_DOCS = {
    "cursor": "https://docs.cursor.com/context/model-select",
    "copilot": "https://docs.github.com/en/copilot/using-github-copilot",
    "gemini_antigravity": "https://cloud.google.com/gemini/docs",
    "roo_code": "https://github.com/RooVetGit/Roo-Code",
    "cline": "https://github.com/cline/cline"
}

def parse_frontmatter(content):
    if not content.startswith("---"):
        return {}
    parts = content.split("---", 2)
    if len(parts) < 3:
        return {}
    yaml_lines = parts[1].strip().split("\n")
    data = {}
    for line in yaml_lines:
        line = line.strip()
        if not line or line.startswith("#"):
            continue
        if ":" in line:
            k, v = line.split(":", 1)
            data[k.strip()] = v.strip()
    return data

def parse_simple_yaml(path):
    if not os.path.isfile(path):
        return {}
    res = {}
    current_section = None
    with open(path, "r", encoding="utf-8") as f:
        for line in f:
            line_str = line.strip()
            if not line_str or line_str.startswith("#"):
                continue
            if line_str.endswith(":") and not line_str.startswith("-"):
                current_section = line_str[:-1].strip()
                res[current_section] = {}
            elif ":" in line_str and current_section:
                k, v = line_str.split(":", 1)
                val = v.strip()
                if val == "null":
                    val = None
                res[current_section][k.strip()] = val
    return res

def main():
    parser = argparse.ArgumentParser(description="Resolve agent model_tier for target platform")
    parser.add_argument("--platform", required=True, help="Target platform name")
    parser.add_argument("--defaults", required=True, help="Path to shared/model-defaults.yaml")
    parser.add_argument("--overrides", help="Path to user model-overrides.yaml")
    parser.add_argument("--agents-dir", required=True, help="Path to agents directory")
    args = parser.parse_args()

    platform_key = args.platform.replace("-", "_")
    defaults = parse_simple_yaml(args.defaults)
    overrides = parse_simple_yaml(args.overrides) if args.overrides else {}

    platform_mapping = overrides.get(platform_key) or defaults.get(platform_key) or {}

    if not os.path.isdir(args.agents_dir):
        sys.exit(0)

    agent_files = [
        f for f in os.listdir(args.agents_dir)
        if f.endswith(".md") and f != "CHANGELOG.md"
    ]

    warn_count = 0
    for file_name in sorted(agent_files):
        path = os.path.join(args.agents_dir, file_name)
        with open(path, "r", encoding="utf-8") as f:
            content = f.read()
        
        fm = parse_frontmatter(content)
        agent_name = fm.get("name", file_name.replace(".md", ""))
        tier = fm.get("model_tier")
        model_override = fm.get("model")

        if not tier:
            continue

        resolved = platform_mapping.get(tier)

        if resolved is None:
            doc_link = PLATFORM_DOCS.get(platform_key, "platform documentation")
            print(f"WARN: agent '{agent_name}' requested model_tier: '{tier}'")
            print(f"      platform '{args.platform}' does not honor per-agent model selection.")
            print(f"      Set your global model preference via {doc_link}.")
            warn_count += 1

    sys.exit(0)

if __name__ == "__main__":
    main()
