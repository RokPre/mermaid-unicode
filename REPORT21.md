# Gantt Renderer Report

# Context

- Problem: Mermaid `gantt` diagrams were unsupported, preventing basic project timeline input from rendering.
- Constraints: Full Gantt support requires date math, dependencies, exclusions, and milestone semantics. The first pass needed a conservative, deterministic terminal representation.

# Goals

- Primary success criteria: Detect `gantt`, parse title, date format, sections, and task rows, render terminal task bars, and pass focused plus full tests.
- Secondary success criteria: Preserve raw task specs so unsupported advanced syntax remains visible instead of being silently dropped.

# Approach

- Chosen approach: Added `pkg/gantt` with source-order task parsing and fixed-width terminal bars.
- Rejected options: Did not implement scaled date layout or full Day.js date parsing in this slice.

# Implementation

- Architecture / flow: `gantt.Parse` reads `gantt`, `title`, `dateFormat`, `section`, and task rows split on the first colon. `gantt.Render` emits title/date format lines, section headers, task names, fixed bars, and the original task spec.
- Key files or components: `pkg/gantt/gantt.go`, `pkg/gantt/gantt_test.go`, `cmd/diagram.go`, `cmd/diagram_test.go`, `README.md`, and `TODO.md`.
- Example: `Task one :a1, 2024-01-01, 7d` renders `Task one` with a terminal bar and the raw schedule spec.

# Results

- Outputs: `gantt` is now supported for a basic source-order timeline table.
- Metrics or observations: The priority-5 `add-gantt-timeline-renderer` TODO is complete at the conservative first-pass scope.
- Verification: `go test ./pkg/gantt` passed. `go test ./cmd -run 'TestDiagramFactory'` passed. `go test ./...` passed.

# Decisions

- Tradeoffs made: Task rows preserve raw specs instead of partially parsing dates and dependencies. This avoids misleading date calculations before a real date model exists.

# Limitations

- Known issues, uncertainties, or risks: Date scaling, milestones, excludes, weekends, and dependency layout are not implemented.

# Next steps

1. Implement the priority-5 gitgraph lane renderer.
2. Add date scaling later with explicit date-format support and focused tests.

# Reproducibility

1. Run `go test ./pkg/gantt`.
2. Run `go test ./cmd -run 'TestDiagramFactory'`.
3. Run `go test ./...`.
