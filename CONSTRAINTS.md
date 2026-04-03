# Constraints

> Architectural rules for this project. Validated by `/stack-audit`.
> Component-level constraints (if any) are in each component's CONSTRAINTS.md and checked automatically.

## Constraints

### No Defensive Error Handling
**Rule**: Do not add try/catch blocks to log-and-continue, and do not null-check mandatory objects. Let exceptions propagate naturally. Use preconditions for invariants.
**Why**: We are in experimenting/revising phase. Defensive error handling covers up failure modes and makes root causes harder to find. Errors should be loud.
**Verify**: `grep -rn 'catch\s*(' --include='*.ts' --include='*.go' | grep -v node_modules | grep -v gen/`
**Scope**: project-wide

### No Workarounds Without Root Cause
**Rule**: Always find the root cause of an issue before proposing a fix. Never create workarounds without asking.
**Why**: Workarounds accumulate and hide the real problem. Past sessions wasted time patching symptoms.
**Verify**: manual
**Scope**: project-wide

### Use copyUnit() for Unit Copies
**Rule**: When creating unit copies (e.g., for history recording), always use `copyUnit()` in services/moves.go. Never manually construct a `&v1.Unit{}` copy.
**Why**: Manual copies miss fields when new proto fields are added to Unit (e.g., Shortcut was forgotten).
**Verify**: `grep -rn 'v1\.Unit{' --include='*.go' services/ | grep -v 'copyUnit\|func copyUnit\|test'`
**Scope**: services/

### Auth Through OneAuth Only
**Rule**: All authentication must go through oneauth middleware. No direct JWT parsing, no custom auth headers, no manual token verification.
**Why**: Direct auth handling skips token rotation, session invalidation, and rate limiting provided by the stack.
**Verify**: `grep -rn 'jwt\.Parse\|jwt\.Verify\|jwt\.New' --include='*.go' | grep -v oneauth | grep -v vendor`
**Scope**: project-wide
<!-- Candidate for promotion to oneauth/CONSTRAINTS.md if seen in other projects -->

### Rate Limiting Through GoAppLib Only
**Rule**: Use goapplib rate limiting middleware. No hand-rolled rate limiters.
**Why**: Stack provides consistent rate limiting with auth-aware tiers. Custom implementations drift.
**Verify**: `grep -rn 'rate.*limit\|RateLimit\|rateLim' --include='*.go' | grep -v goapplib | grep -v vendor | grep -v _test.go`
**Scope**: project-wide
<!-- Candidate for promotion to goapplib/CONSTRAINTS.md if seen in other projects -->

### No Manual Builds
**Rule**: Do not run `npm build`, `npm run build`, or `buf generate` manually. The web module and proto files auto-rebuild on change. Do not rebuild the server — devloop runs it continuously.
**Why**: Manual builds conflict with the file-watching build pipeline and cause confusing stale state.
**Verify**: manual
**Scope**: project-wide

### Lazy Top-Up Pattern for Units
**Rule**: Units must not have their movement points reset at turn start. Use `topUpUnitIfNeeded()` on-demand when a unit is accessed for actions or options.
**Why**: Eager reset at turn start causes state inconsistencies when units are accessed between turns or during replays. Lazy pattern ensures consistency.
**Verify**: `grep -rn 'DistanceLeft\s*=' --include='*.go' services/ | grep -v topUpUnit | grep -v _test.go | grep -v copyUnit`
**Scope**: services/
