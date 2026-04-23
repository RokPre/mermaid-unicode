# Expanded Flowchart Shape Syntax Report

# Context

- Problem: Mermaid v11.3.0+ supports expanded flowchart shape metadata such as `A@{ shape: rounded }`, but the terminal renderer only understood older delimiter syntax such as `A(Rounded)` and `A{Decision}`.
- Constraints: The renderer already has a finite set of terminal shape approximations. The work needed to map Mermaid aliases onto those existing shapes without attempting full SVG shape parity.

# Goals

- Primary success criteria: Parse expanded flowchart node metadata and render supported `shape:` aliases using the existing Unicode and ASCII node shape machinery.
- Secondary success criteria: Preserve optional `label:` metadata, preserve class shorthand such as `:::important`, keep unsupported shape names from corrupting node IDs, and pass the full Go test suite.

# Approach

- Chosen approach: Added a metadata parser inside the existing graph node parser. The parser recognizes `@{ ... }`, extracts properties, maps supported `shape:` aliases to `graphNodeShape`, and leaves unknown shape values unset.
- Rejected options: Did not implement all Mermaid expanded shapes. Many documented Mermaid shapes have no stable terminal equivalent yet, so mapping only aliases that match existing renderer shapes keeps behavior predictable.

# Implementation

- Architecture / flow: `parseNode` now checks for expanded metadata before delimiter-based shape parsing. `parseExpandedNodeProperties` extracts metadata key/value pairs. `expandedNodeShape` maps aliases such as `rect`, `rounded`, `stadium`, `subroutine`, `db`, `circle`, `decision`, `hexagon`, and `lean-r` into the existing terminal shape enum.
- Key files or components: `cmd/parse.go` contains the parser and alias mapping. `cmd/parse_test.go` covers parser behavior, labels, class shorthand, unsupported shapes, and shape persistence across bare references. `cmd/render_graph_test.go` covers rendered Unicode glyphs for expanded shape syntax. `TODO.md` records `add-expanded-shape-syntax` as a completed flowchart subtask.
- Example: `A@{ shape: rounded, label: "Research" }` now parses as node `A`, renders label `Research`, and uses rounded corners in Unicode mode.

# Results

- Outputs: Expanded shape syntax works for the renderer's existing shape families: square, rounded, stadium, double/subroutine, database, circle, decision, hexagon, and parallelogram.
- Metrics or observations: Unsupported expanded shapes such as `cloud` no longer become part of the node name; they render as an unshaped node with the clean node id.
- Verification: Ran `go test ./cmd -run 'TestParseNodeExpandedShape|TestMermaidFileToMapKeepsExpandedNodeShape|TestRenderGraphUsesExpandedShapeSyntax'`; it passed. Ran `go test ./...`; it passed.

# Decisions

- Tradeoffs made: Unsupported expanded shapes are ignored rather than hard errors. That matches the current graph parser style, avoids breaking diagrams that contain future Mermaid shapes, and keeps the node id usable in subsequent edges.

# Limitations

- Known issues, uncertainties, or risks: This does not add exact terminal renderings for Mermaid shapes such as cloud, document, hourglass, braces, flags, or storage variants. It also does not add `BT`/`RL` layout support or extra link heads.

# Next steps

1. Continue `audit-flowchart-parity` with `BT` and `RL` direction support if the current layout model can be extended cleanly.
2. Add a flowchart support matrix to README after the next parity slice so users can see which expanded shapes map to terminal approximations.

# Reproducibility

1. Run `go test ./cmd -run 'TestParseNodeExpandedShape|TestMermaidFileToMapKeepsExpandedNodeShape|TestRenderGraphUsesExpandedShapeSyntax'`.
2. Run `go test ./...`.
