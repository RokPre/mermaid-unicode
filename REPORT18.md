# Timeline Renderer Report

# Context

- Problem: Mermaid `timeline` diagrams were unsupported, despite having a compact syntax that maps well to terminal rows.
- Constraints: Timeline syntax is experimental and includes browser theme options. The first terminal slice should preserve source order and avoid theme complexity.

# Goals

- Primary success criteria: Detect `timeline`, parse title, sections, periods, events, and direction, render TD/LR output, and pass focused plus full tests.
- Secondary success criteria: Keep multiple events per period readable and deterministic.

# Approach

- Chosen approach: Added `pkg/timeline` with a source-order parser and two render modes: vertical rows for TD/default and compact horizontal rows for LR.
- Rejected options: Did not implement browser theme variables or color schemes in this slice.

# Implementation

- Architecture / flow: `timeline.Parse` reads `timeline`, optional inline direction, `title`, `section`, `direction`, and item rows of the form `period : event : event`. `timeline.Render` emits section headers and period rows, preserving event order.
- Key files or components: `pkg/timeline/parser.go`, `pkg/timeline/renderer.go`, `pkg/timeline/timeline_test.go`, `cmd/diagram.go`, `cmd/diagram_test.go`, `README.md`, and `TODO.md`.
- Example: `2024 : Idea : Research` renders one period row with `Idea` followed by an indented `Research` event.

# Results

- Outputs: `timeline` is now supported for title, sections, periods, multiple events, and TD/LR rendering.
- Metrics or observations: The priority-4 `add-timeline-renderer` TODO is complete.
- Verification: `go test ./pkg/timeline` passed. `go test ./cmd -run 'TestDiagramFactory'` passed. `go test ./...` passed.

# Decisions

- Tradeoffs made: LR output is a compact row representation instead of a proportional date scale, because Mermaid timeline periods are labels rather than guaranteed machine-readable dates.

# Limitations

- Known issues, uncertainties, or risks: Theme/color schemes and text wrapping are not implemented beyond deterministic row output.

# Next steps

1. Implement the priority-4 user journey renderer.
2. Add wrapping only if real timeline examples exceed practical terminal widths.

# Reproducibility

1. Run `go test ./pkg/timeline`.
2. Run `go test ./cmd -run 'TestDiagramFactory'`.
3. Run `go test ./...`.
