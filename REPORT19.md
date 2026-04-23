# User Journey Renderer Report

# Context

- Problem: Mermaid `journey` diagrams were unsupported, though the syntax maps cleanly to a terminal table.
- Constraints: User journeys do not need graph layout; readability comes from section grouping, score validation, and actor lists.

# Goals

- Primary success criteria: Detect `journey`, parse title/sections/tasks, validate scores from 1 to 5, render task rows, and pass focused plus full tests.
- Secondary success criteria: Provide ASCII and Unicode score bars.

# Approach

- Chosen approach: Added `pkg/journey` with a small parser and a sectioned task-table renderer.
- Rejected options: Did not use the graph renderer because journey diagrams are naturally tabular.

# Implementation

- Architecture / flow: `journey.Parse` reads `journey`, optional `title`, `section`, and task rows in the form `Task: score: actor, actor`. `journey.Render` emits section headers, scores, bars, and actors.
- Key files or components: `pkg/journey/parser.go`, `pkg/journey/renderer.go`, `pkg/journey/journey_test.go`, `cmd/diagram.go`, `cmd/diagram_test.go`, `README.md`, and `TODO.md`.
- Example: `Write code: 5: Me` renders a `5/5` score with a full bar and actor list.

# Results

- Outputs: `journey` is now supported as a sectioned terminal table.
- Metrics or observations: The priority-4 `add-user-journey-renderer` TODO is complete.
- Verification: `go test ./pkg/journey` passed. `go test ./cmd -run 'TestDiagramFactory'` passed. `go test ./...` passed.

# Decisions

- Tradeoffs made: Actor wrapping is not yet implemented; actor lists stay on the row for this first deterministic renderer.

# Limitations

- Known issues, uncertainties, or risks: Very long actor lists can produce wide rows.

# Next steps

1. Implement the priority-4 pie and quadrant chart renderers.
2. Add actor wrapping later if real journey diagrams need narrower output.

# Reproducibility

1. Run `go test ./pkg/journey`.
2. Run `go test ./cmd -run 'TestDiagramFactory'`.
3. Run `go test ./...`.
