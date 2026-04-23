# Open Solid Flowchart Connector Report

# Context

- Problem: Mermaid flowcharts support open solid links with `---`, but the renderer only supported arrow links, heavy arrow links, dashed arrow links, and open dashed links.
- Constraints: The change needed to fit the existing graph parser style, keep current arrow behavior unchanged, and preserve the arrow-style priority model.

# Goals

- Primary success criteria: Parse and render `A --- B` as a solid connector without an arrowhead.
- Secondary success criteria: Support labeled open solid connectors with `A ---|label| B`, add focused tests, and pass the full Go test suite.

# Approach

- Chosen approach: Added two parser patterns beside the existing edge patterns. Both call the existing `setEdgeWithLabel` helper with light line style and `hasArrowHead=false`.
- Rejected options: Did not introduce a new edge type because the existing `textEdge` fields already model line style, label, and arrowhead presence.

# Implementation

- Architecture / flow: `graphProperties.parseString` now recognizes `---` and `---|label|` before standard `-->` arrows. Parsed edges use `graphEdgeLineStyleLight` and no arrowhead. Rendering reuses the current line drawing path.
- Key files or components: `cmd/parse.go` contains the new parse patterns. `cmd/parse_test.go` verifies line style, labels, and `hasArrowHead=false`. `cmd/render_graph_test.go` verifies rendered solid lines and no arrowhead. `TODO.md` records `add-open-solid-connectors` as complete.
- Example: `graph LR\nA ---|open| B` now renders a light connector labeled `open` without drawing `►`.

# Results

- Outputs: Open solid Mermaid connectors are now part of the supported flowchart subset.
- Metrics or observations: Existing styled edge behavior remains unchanged, including heavy, dashed arrow, and open dashed output.
- Verification: Ran `go test ./cmd -run 'TestMermaidFileToMapParsesEdgeLineStyles|TestRenderGraphUsesStyledEdgeGlyphs'`; it passed. Ran `go test ./...`; it passed.

# Decisions

- Tradeoffs made: `---` uses the default light edge style path with `lineStyleSet=false`, so configured default edge styles can still affect standard unstyled connectors consistently with existing `-->` behavior.

# Limitations

- Known issues, uncertainties, or risks: This does not add Mermaid circle-head, cross-head, or bidirectional operators. It also keeps the parser's existing requirement for spaces around link operators.

# Next steps

1. Continue `audit-flowchart-parity` by documenting the supported flowchart subset in README, including supported connectors and unsupported link heads.
2. Decide whether to add circle-head and cross-head connectors as terminal approximations or explicitly reject them with clear errors.

# Reproducibility

1. Run `go test ./cmd -run 'TestMermaidFileToMapParsesEdgeLineStyles|TestRenderGraphUsesStyledEdgeGlyphs'`.
2. Run `go test ./...`.
