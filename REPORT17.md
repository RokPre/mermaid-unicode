# Mindmap Tree Renderer Report

# Context

- Problem: Mermaid `mindmap` input was still unsupported, even though a terminal tree is a practical approximation of the hierarchy.
- Constraints: Mermaid mindmaps can include radial layout and icons, but those do not translate cleanly to terminal output.

# Goals

- Primary success criteria: Detect `mindmap`, parse indentation into a hierarchy, render Unicode and ASCII trees, and pass focused plus full tests.
- Secondary success criteria: Normalize simple node shape wrappers without adding graph-shape rendering complexity.

# Approach

- Chosen approach: Added a small `pkg/mindmap` parser that builds a tree from indentation and a renderer that emits standard tree glyphs.
- Rejected options: Did not attempt radial layout or icon support; both are outside a first terminal renderer pass.

# Implementation

- Architecture / flow: `mindmap.Parse` reads the root and children by indentation level, rejects multiple roots and inconsistent outdents, and strips simple Mermaid shape delimiters from node text. `mindmap.Render` walks the tree recursively and emits Unicode `├──`/`└──` or ASCII `|--`/``--` branches.
- Key files or components: `pkg/mindmap/parser.go`, `pkg/mindmap/renderer.go`, `pkg/mindmap/mindmap_test.go`, `cmd/diagram.go`, `cmd/diagram_test.go`, `README.md`, and `TODO.md`.
- Example: A root with children `Idea` and `Plan` renders as a stable terminal tree with nested branches.

# Results

- Outputs: `mindmap` is now supported as a terminal tree renderer.
- Metrics or observations: The priority-3 `add-mindmap-tree-renderer` TODO is complete.
- Verification: `go test ./pkg/mindmap` passed. `go test ./cmd -run 'TestDiagramFactory'` passed. `go test ./...` passed.

# Decisions

- Tradeoffs made: Shape syntax is normalized to text rather than rendered as boxes, keeping the mindmap output dense and tree-like.

# Limitations

- Known issues, uncertainties, or risks: Radial mindmap layout, icons, and rich per-node shapes are not implemented.

# Next steps

1. Implement the priority-4 timeline renderer.
2. Add icon handling only if a terminal-safe representation is defined.

# Reproducibility

1. Run `go test ./pkg/mindmap`.
2. Run `go test ./cmd -run 'TestDiagramFactory'`.
3. Run `go test ./...`.
