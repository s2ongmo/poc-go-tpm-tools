# PoC: Supply Chain Attack via Mutable Branch Reference in GitHub Actions Release Pipeline

## Overview

This repository demonstrates a supply chain vulnerability in `google/go-tpm-tools` where a third-party GitHub Action (`dtolnay/rust-toolchain@stable`) is referenced by a **mutable branch name** in the release workflow. An attacker who compromises the action owner's account can inject arbitrary code into the release pipeline without modifying the victim repository.

## Architecture

```
┌─────────────────────────────────┐     ┌──────────────────────────────────┐
│  s2ongmo/poc-go-tpm-tools       │     │  s2ongmo/poc-action              │
│  (simulates google/go-tpm-tools)│     │  (simulates dtolnay/rust-toolchain)
│                                 │     │                                  │
│  .github/workflows/             │     │  action.yml (on `stable` branch) │
│    releaser.yaml                │────▶│                                  │
│      uses: poc-action@stable    │     │  Phase 1: BENIGN (v1.0.0 build) │
│                                 │     │  Phase 2: COMPROMISED (v1.0.1)  │
│  Workflow NEVER changes.        │     │  Branch force-pushed by attacker │
└─────────────────────────────────┘     └──────────────────────────────────┘
```

## Evidence

### Phase 1: Benign Release (v1.0.0)

- **Workflow run**: https://github.com/s2ongmo/poc-go-tpm-tools/actions/runs/23230947906
- **Release**: https://github.com/s2ongmo/poc-go-tpm-tools/releases/tag/v1.0.0
- Binary output: `INTEGRITY_CHECK=CLEAN`
- `ACTION_VERSION=benign-1.0.0`
- SHA256: `5bb593696b50c44fd75b1ab863cff5ee10963a4ad425fbc25dd396be9d879eb0`

### Phase 2: Compromised Release (v1.0.1)

- **Workflow run**: https://github.com/s2ongmo/poc-go-tpm-tools/actions/runs/23230995530
- **Release**: https://github.com/s2ongmo/poc-go-tpm-tools/releases/tag/v1.0.1
- Binary output: `INTEGRITY_CHECK=COMPROMISED`
- `ACTION_VERSION=COMPROMISED-by-supply-chain-attack`
- SHA256: `3fa09332f420a72934ec9aa3e9deefe786a098d71e157dc5dc61852ca223d7a1`

### What changed between v1.0.0 and v1.0.1?

| Component | v1.0.0 → v1.0.1 |
|-----------|------------------|
| `poc-go-tpm-tools` source code | **NO CHANGE** |
| `poc-go-tpm-tools` workflow | **NO CHANGE** |
| `poc-action@stable` branch content | **CHANGED** (force-push by simulated attacker) |
| Released binary | **DIFFERENT** (compromised source injected at build time) |

The **only** change was a force-push to the `stable` branch of the third-party action repository. The victim repository (`poc-go-tpm-tools`) was never modified.

## Attack Chain Demonstrated

1. `poc-action` created with benign `action.yml` on `stable` branch
2. `poc-go-tpm-tools` references `s2ongmo/poc-action@stable` in release workflow (mirrors `dtolnay/rust-toolchain@stable`)
3. v1.0.0 tag pushed → benign release built and published
4. **Attacker force-pushes** compromised `action.yml` to `stable` branch of `poc-action`
5. v1.0.1 tag pushed → **compromised binary built and published**
6. The compromised action:
   - Accessed the runner filesystem (`ls -la *.go`)
   - Read environment variables (`GITHUB_REPOSITORY`, `GITHUB_REF`, etc.)
   - **Modified `main.go` source code** before the build step
   - Set `ACTION_VERSION=COMPROMISED-by-supply-chain-attack`

## How This Maps to google/go-tpm-tools

| PoC | Real Target |
|-----|-------------|
| `s2ongmo/poc-action@stable` | `dtolnay/rust-toolchain@stable` |
| `s2ongmo/poc-go-tpm-tools` | `google/go-tpm-tools` |
| `releaser.yaml` with `contents: write` | Identical workflow structure and permissions |
| Simulated binary `gotpm-poc` | Real binary `gotpm` (53K+ downloads across releases) |
| Force-push to `stable` branch | Same — no branch protection on `dtolnay/rust-toolchain` |

## Verification

1. Check that `poc-go-tpm-tools` workflow was **never modified** between v1.0.0 and v1.0.1
2. Compare the two releases — different binaries, different SHA256 hashes
3. Check `poc-action` commit history on `stable` branch — the compromised commit is visible
4. Review workflow logs for both runs to see benign vs compromised action output
