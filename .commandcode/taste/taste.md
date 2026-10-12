# Taste

## Communication & documentation

- Wants project documentation in both English and Persian (Farsi) — include an RTL-laid-out Persian version alongside the English one, kept in parity (same structure/sections), with code blocks kept LTR. Confidence: 0.75
- Documents are expected to cover practical usage end to end: what a feature is, installation/per-OS instructions, field-by-field reference, client usage, and troubleshooting/FAQ. Confidence: 0.55

## Workflow & tooling

- Prefers fully automated update/release pipelines over manual human-review gates (chose "fully automatic release" over PR-gated or manual-dispatch options). Confidence: 0.5
- Expects verification of changes through the remote CI (GitHub Actions build runs) rather than local compilation alone — asks the agent to inspect the repo's Action builds. Confidence: 0.4
