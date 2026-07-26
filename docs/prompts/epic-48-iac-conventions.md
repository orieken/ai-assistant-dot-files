# Epic 48 — Infrastructure-as-Code Conventions

Source: `docs/audits/framework-gap-audit-2026-07-25.md` § Dimension 3.

## Target repo

`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files`. Do NOT push.

## Prior context

`shared/rules/*.md` covers 5 languages (typescript, go, python, java, csharp) plus cross-cutting rules (design-principles, architecture-guardrails, testing-conventions, approval-gates). No IaC guardrails today. `devops-engineer` and `sre-engineer` agents produce IaC-adjacent artifacts but have no rule doc to check against.

## Scope

**One commit.** Create `shared/rules/iac-conventions.md` matching the shape of existing convention files. Sections:

- **Terraform / OpenTofu** — module layout, state locking, no-hardcoded-secrets, `terraform fmt` + `tflint`, explicit provider versions
- **Dockerfile** — multi-stage builds, non-root USER, minimal base images, no `apt-get` without version pinning, no secrets in ENV
- **Kubernetes / Helm** — resource limits mandatory, no `latest` image tags, ServiceAccount + RBAC per app, network policies, secrets from ExternalSecrets/SealedSecrets not plain-text
- **GitHub Actions** — pinned action versions (SHA, not tags), `permissions:` block explicit (least-privilege), no secrets in `run:` echo, `pull_request_target` used only with explicit checkout of PR SHA

Cross-reference:
- `shared/rules/architecture-guardrails.md` #3 (no hardcoded secrets)
- `shared/agents/devops-engineer.md` (should reference this new rule file)

Update `shared/agents/devops-engineer.md`'s "Before beginning any task, read..." block to include `shared/rules/iac-conventions.md`.

Commit: `docs(rules): add iac-conventions.md + wire into devops-engineer (Epic 48)`.

## Discipline

Standard — match other prompts in `docs/prompts/`.

## Escalation

- If existing conventions in another rule doc contradict what you're about to write here — halt, describe.
- If devops-engineer's existing behavior conflicts with a rule you're adding — halt (may need to update the agent).

## Report (under 100 words)

```
Commit: <sha>
Sections landed: Terraform, Docker, K8s, GHA (or subset)
devops-engineer.md updated: yes | no
Health-check green: yes
```

Go.
