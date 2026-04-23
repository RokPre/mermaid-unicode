# State Diagram Renderer Report

# Context

- Problem: `stateDiagram` and `stateDiagram-v2` were still known unsupported Mermaid types, so simple state machines with start/end markers and labeled transitions could not render.
- Constraints: State diagrams include deep features such as nested composites and concurrency. The first slice needed to be deterministic and useful without claiming full browser parity.

# Goals

- Primary success criteria: Detect state diagrams, parse simple states and transitions, render start/end markers and labels, and pass focused plus full tests.
- Secondary success criteria: Include first-pass support for notes, choice/fork/join markers, and simple composite frames while keeping other diagram renderers stable.

# Approach

- Chosen approach: Added a dedicated `pkg/statediagram` package and registered it through the existing diagram factory. The renderer uses pairwise transition rows for simple states and additional note/composite frame rows for supporting syntax.
- Rejected options: Did not force state diagrams through the flowchart renderer. State-specific start/end markers, descriptions, and composite frames are clearer in a state-specific parser.

# Implementation

- Architecture / flow: `statediagram.Parse` reads `stateDiagram` and `stateDiagram-v2`, optional `direction`, state aliases, descriptions, transitions, notes, choice/fork/join declarations, and simple `state Parent { ... }` composite blocks. `statediagram.Render` creates Unicode or ASCII state boxes, compact start/end markers, transition connectors, note boxes, and composite frames.
- Key files or components: `pkg/statediagram/parser.go` handles syntax and validation. `pkg/statediagram/renderer.go` handles terminal output. `cmd/diagram.go` registers state diagrams. `pkg/statediagram/statediagram_test.go` and `cmd/diagram_test.go` cover parser, renderer, and factory behavior. `README.md` and `TODO.md` document support.
- Example: `[*] --> Still` renders a compact start marker connected to the `Still` state. `Still --> [*] : done` renders an end marker with the `done` label.

# Results

- Outputs: `stateDiagram` and `stateDiagram-v2` are now supported for a practical terminal subset.
- Metrics or observations: The priority-2 `add-state-diagram-renderer` TODO is complete.
- Verification: `go test ./pkg/statediagram` passed. `go test ./cmd -run 'TestDiagramFactory'` passed. `go test ./...` passed.

# Decisions

- Tradeoffs made: Composite states render as simple frames listing child state ids rather than nested full graph layouts. This keeps the first pass stable and easy to extend.

# Limitations

- Known issues, uncertainties, or risks: Deep nested composite layout, concurrency regions, and block-form multiline notes are not fully implemented. Transition layout repeats state boxes per transition.

# Next steps

1. Implement the priority-3 requirement diagram renderer.
2. Revisit global graph-style layout for ER/class/state diagrams if repeated boxes become too verbose.

# Reproducibility

1. Run `go test ./pkg/statediagram`.
2. Run `go test ./cmd -run 'TestDiagramFactory'`.
3. Run `go test ./...`.
