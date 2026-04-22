# Edge Priority Rendering Report

# Context

- Problem: Styled graph edges had distinct glyphs, but overlapping or merged line cells could lose heavy styling when combined with default light arrows.
- Constraints: The change needed to keep existing pathfinding and layout behavior intact while making the final drawn output respect arrow type priority.

# Goals

- Primary success criteria: Give each graph arrow line style a priority and render higher-priority styles over lower-priority ones.
- Secondary success criteria: Keep heavy `==>` arrows above default `-->` arrows, preserve readable junctions, and pass the full Go suite.

# Approach

- Chosen approach: Defined explicit line priorities in the graph charset layer: dashed < light < heavy. Edge rendering is now sorted by effective priority, and junction merging picks the highest-priority glyph family when line connections overlap.
- Rejected options: Reworking pathfinding was unnecessary because this issue is about draw order and cell merging, not route selection.

# Implementation

- Architecture / flow: `graph.draw` now sorts edges by `edgeDrawPriority` before producing line, corner, arrowhead, box-start, and label drawings. `mergeJunctions` now considers glyph priority when combining compatible line cells, so heavy glyphs survive overlap with light defaults.
- Key files or components: `cmd/graph.go`, `cmd/graph_charset.go`, `cmd/graph_charset_test.go`, `cmd/testdata/extended-chars/styled_edges_lr.txt`, and `cmd/testdata/extended-chars/styled_edges_td.txt`.
- Example: A heavy edge leaving a node now uses a heavy junction such as `┣` or `┳`, and a heavy horizontal overlap remains `━` instead of degrading to `─`.

# Results

- Outputs: Edge style priority is explicit and covered by focused tests plus golden fixture updates.
- Metrics or observations: Heavy line glyphs now win over light and dashed glyphs when the connection topology is the same or when a mixed junction is needed.
- Verification: Ran `go test ./cmd -run 'TestGraphCharset|TestEdgeLineStylePriority|TestMergeDrawingsKeepsHigherPriorityLineGlyph'` and `go test ./...`; both passed.

# Decisions

- Tradeoffs made: Heavy mixed junctions use the heavy junction family (`╋`, `┳`, `┣`, etc.). Dashed remains lower priority than light because dashed connectors are less visually dominant than default solid connectors.

# Limitations

- Known issues, uncertainties, or risks: Priority affects the final glyph chosen for an overlapping cell; it does not alter pathfinding to avoid or prefer routes based on style.

# Next steps

1. Add a user-facing example if users need documentation for priority behavior.
2. Consider style-aware pathfinding only if diagrams with many overlapping connectors still produce ambiguous results.

# Reproducibility

1. Run `go test ./cmd -run 'TestGraphCharset|TestEdgeLineStylePriority|TestMergeDrawingsKeepsHigherPriorityLineGlyph'`.
2. Run `go test ./...`.
