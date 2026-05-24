# Git Workflow

## Branch Workflow

Do not make task changes directly on `main`. Create a task branch from `main` before implementation work.

Branch prefixes should match the task type:

| Prefix | Use |
| --- | --- |
| `docs/` | Documentation-only changes |
| `chore/` | Project maintenance, configuration, and repository hygiene |
| `feat/` | New user-facing or system capability |
| `fix/` | Bug fixes |
| `ci/` | CI/CD changes |

Examples:

```text
chore/config-example-local
docs/update-local-development
feat/link-redirect
fix/cache-invalidation
ci/add-github-actions
```

## Commit Message Format

zero-link uses Conventional Commits:

```text
<type>(<scope>): <subject>
```

The scope is optional. The subject must be concise, written in English, and no longer than 72 characters.

## Allowed Types

| Type | Use |
| --- | --- |
| `feat` | New user-facing or system capability |
| `fix` | Bug fix |
| `docs` | Documentation-only change |
| `style` | Formatting or style-only change |
| `refactor` | Code restructuring without behavior change |
| `test` | Test changes |
| `chore` | Project maintenance |
| `ci` | CI/CD changes |
| `build` | Build system, dependencies, packaging |
| `perf` | Performance improvement |
| `revert` | Revert a previous commit |
| `security` | Security improvement or vulnerability fix |

## Valid Examples

```text
docs: add project documentation set
chore(project): initialize repository skeleton
build(nix): add go-zero tooling
feat(link-api): add redirect route
fix(cache): invalidate stale short link cache
test(link-rpc): add redirect resolution tests
security(auth): harden password hashing
```

## Invalid Examples

```text
add docs
Docs: add docs
feat add api
chore_: random cleanup
docs: this subject is intentionally written far beyond the maximum subject length and should fail
```

## Local Hook Installation

Install repository hooks with:

```bash
make install-hooks
```

This configures Git to read hooks from `.githooks/`.

## Bypassing The Hook

Avoid bypassing the hook. If an emergency requires it, Git supports `--no-verify`, but the commit message should still be corrected before review.
