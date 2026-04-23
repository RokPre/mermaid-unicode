# Reverse Flowchart Direction Support Report

# Context

- Problem: Mermaid flowcharts document `RL` and `BT` directions, but the renderer only accepted `LR`, `TD`, and `TB`.
- Constraints: The graph layout code was built around two orientation families: horizontal and vertical. The change needed to add reverse directions without rewriting pathfinding or disturbing existing `LR`/`TD` output.

# Goals

- Primary success criteria: Support `graph RL`, `flowchart RL`, `graph BT`, and `flowchart BT` through parse, layout, pathing, and rendering.
- Secondary success criteria: Preserve existing graph behavior, render visible left/up arrowheads for reverse directions, and pass the full Go test suite.

# Approach

- Chosen approach: Kept the existing layout algorithm, then mirrored mapped node coordinates on the primary axis for reversed directions before sizing and path calculation. Direction checks now ask whether the graph is horizontal or vertical instead of comparing only against `LR` and `TD`.
- Rejected options: Did not rewrite the layout engine around full directional constraints. That would be a broad change and unnecessary for the current Mermaid parity gap.

# Implementation

- Architecture / flow: `mermaidFileToMap` now accepts `RL` and `BT` headers. `graph.isHorizontalLayout`, `graph.isVerticalLayout`, and `graph.isReversedLayout` centralize direction checks. `mirrorMappingForReversedDirection` mirrors node grid coordinates and rebuilds the occupied grid before edge paths are determined.
- Key files or components: `cmd/parse.go` handles the new headers. `cmd/graph.go`, `cmd/direction.go`, `cmd/mapping_node.go`, and `cmd/mapping_edge.go` contain orientation-aware layout/path changes. `cmd/parse_test.go` and `cmd/render_graph_test.go` cover parsing and rendered output.
- Example: `graph RL\nA --> B` now renders `B` to the left of `A` with a left arrowhead. `graph BT\nA --> B` now renders `B` above `A` with an up arrowhead.

# Results

- Outputs: Reverse directions are now part of the supported flowchart subset for graph and flowchart aliases.
- Metrics or observations: The existing graph tests still pass after replacing direct `LR`/`TD` comparisons with orientation helpers.
- Verification: Ran `go test ./cmd -run 'TestMermaidFileToMapParsesReverseDirections|TestRenderGraphSupportsRightToLeftDirection|TestRenderGraphSupportsBottomToTopDirection'`; it passed. Ran `go test ./...`; it passed.

# Decisions

- Tradeoffs made: Mirroring happens after node mapping and before path calculation. This keeps the mature placement algorithm intact while still giving reverse directions correct visual order and arrow orientation.

# Limitations

- Known issues, uncertainties, or risks: Complex reversed-direction subgraph layouts do not yet have dedicated golden fixtures. The full test suite covers existing subgraph behavior, but future work should add explicit `RL`/`BT` subgraph cases if users rely on them heavily.

# Next steps

1. Continue `audit-flowchart-parity` by documenting the supported flowchart subset in README, including `LR`, `RL`, `TD`, `TB`, `BT`, expanded shape aliases, and current link operators.
2. Add tests or docs for unsupported flowchart link heads such as Mermaid circle and cross heads so users see the boundary clearly.

# Reproducibility

1. Run `go test ./cmd -run 'TestMermaidFileToMapParsesReverseDirections|TestRenderGraphSupportsRightToLeftDirection|TestRenderGraphSupportsBottomToTopDirection'`.
2. Run `go test ./...`.
