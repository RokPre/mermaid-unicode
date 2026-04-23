# Requirement Diagram Renderer Report

# Context

- Problem: `requirementDiagram` was still a known unsupported Mermaid type, so requirement and element blocks could not be rendered in the terminal.
- Constraints: Requirement diagrams contain SysML-like semantics. The renderer needed to display the structure without trying to validate requirement quality or enforce domain rules.

# Goals

- Primary success criteria: Detect requirement diagrams, parse requirement and element blocks, parse labeled relationships, render deterministic terminal boxes, and pass the full test suite.
- Secondary success criteria: Keep field ordering stable and document the supported subset.

# Approach

- Chosen approach: Added a dedicated `pkg/requirementdiagram` parser and renderer with boxed items and pairwise relationship rows. Fields are sorted during rendering so output diffs remain stable.
- Rejected options: Did not build semantic validation for fields such as `risk` or `verifymethod`; that is outside terminal rendering scope.

# Implementation

- Architecture / flow: `requirementdiagram.Parse` reads `requirementDiagram`, requirement-like blocks, `element` blocks, key/value fields, and relationships of the form `source - satisfies -> target`. `requirementdiagram.Render` draws requirement and element boxes with field rows and connector labels.
- Key files or components: `pkg/requirementdiagram/parser.go`, `pkg/requirementdiagram/renderer.go`, `pkg/requirementdiagram/requirementdiagram_test.go`, `cmd/diagram.go`, `cmd/diagram_test.go`, `README.md`, and `TODO.md`.
- Example: `tester - verifies -> req` renders the `tester` element and `req` requirement connected by a `verifies` relationship.

# Results

- Outputs: `requirementDiagram` is now supported for requirement/element blocks and relationships.
- Metrics or observations: The priority-3 `add-requirement-diagram-renderer` TODO is complete.
- Verification: `go test ./pkg/requirementdiagram` passed. `go test ./cmd -run 'TestDiagramFactory'` passed. `go test ./...` passed.

# Decisions

- Tradeoffs made: Relationship rendering repeats boxes per relationship instead of computing a global graph layout, matching the current bounded renderer approach used by ER, class, and state diagrams.

# Limitations

- Known issues, uncertainties, or risks: The renderer does not enforce required fields or accepted enum values. Styling is not yet applied to requirement fields.

# Next steps

1. Implement the priority-3 mindmap tree renderer.
2. Add stricter requirement field validation only if real workflows need it.

# Reproducibility

1. Run `go test ./pkg/requirementdiagram`.
2. Run `go test ./cmd -run 'TestDiagramFactory'`.
3. Run `go test ./...`.
