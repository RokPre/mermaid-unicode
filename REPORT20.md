# Pie And Quadrant Chart Renderer Report

# Context

- Problem: Mermaid `pie` and `quadrantChart` diagrams were unsupported. Both can be represented in the terminal more effectively as bars/tables or grids than as browser-style graphics.
- Constraints: Terminal cells are not suitable for circular pie geometry, and quadrant points need deterministic placement with clear validation.

# Goals

- Primary success criteria: Detect and render `pie` and `quadrantChart`, validate numeric values, and pass focused plus full tests.
- Secondary success criteria: Provide ASCII and Unicode output modes.

# Approach

- Chosen approach: Added `pkg/charts` with separate pie and quadrant parsers/renderers. Pie charts render source-order percentage bars; quadrant charts render an 11x11 coordinate grid plus point labels.
- Rejected options: Did not attempt circular terminal pies.

# Implementation

- Architecture / flow: `charts.ParsePie` reads title, `showData`, labels, and values. `charts.RenderPie` computes totals and percentages. `charts.ParseQuadrant` reads title, axes, quadrant labels, and `[x, y]` points while rejecting values outside 0..1. `charts.RenderQuadrant` plots points on a fixed grid.
- Key files or components: `pkg/charts/pie.go`, `pkg/charts/quadrant.go`, `pkg/charts/charts_test.go`, `cmd/diagram.go`, `cmd/diagram_test.go`, `README.md`, and `TODO.md`.
- Example: `Dogs : 75` and `Cats : 25` render as 75% and 25% bars.

# Results

- Outputs: `pie` and `quadrantChart` are now supported.
- Metrics or observations: The priority-4 `add-pie-and-quadrant-renderers` TODO is complete.
- Verification: `go test ./pkg/charts` passed. `go test ./cmd -run 'TestDiagramFactory'` passed. `go test ./...` passed.

# Decisions

- Tradeoffs made: Pie slices preserve source order instead of sorting, matching Mermaid input order and keeping output predictable.

# Limitations

- Known issues, uncertainties, or risks: Quadrant labels and axes are parsed but only the title, grid, and point list are rendered in this first pass.

# Next steps

1. Implement the priority-5 Gantt renderer.
2. Add richer quadrant axis labeling later if terminal examples need it.

# Reproducibility

1. Run `go test ./pkg/charts`.
2. Run `go test ./cmd -run 'TestDiagramFactory'`.
3. Run `go test ./...`.
