# Planning Notes

## Immediate Goals
- Review existing repository structure, tooling, and conventions to tailor guidelines.
- Extract key build/test workflows from Makefile, README, and Go modules.
- Capture coding style and testing expectations based on current project artifacts.
- Draft AGENTS.md aligning with prompt scope and repository practices.

## Open Questions / Assumptions
- Assume Go is primary language (per `go.mod` and `cmd/internal` layout).
- Need confirmation if additional languages/tools (e.g., frontend) exist beyond inspected dirs.
- Clarify whether CI/CD conventions exist; will infer from Makefile and repo structure if not documented.

## Next Steps
1. Inspect Makefile, README, and `internal` / `cmd` directories for relevant info.
2. Identify testing frameworks and existing conventions from go files and docs.
3. Draft concise repository guidelines (200-400 words) with required sections.
4. Review draft for tone, accuracy, and word count before finalizing AGENTS.md.
