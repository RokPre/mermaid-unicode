# Sequence Fragment Frame Report

# Context

- Problem: The sequence renderer supported ordered messages, notes, activation directives, and activation shorthand, but Mermaid fragment syntax such as `loop`, `alt`, `else`, `opt`, and `par` still failed to parse.
- Constraints: The existing sequence renderer is row-oriented and terminal-focused. Fragment support needed to preserve source order without attempting browser-style nested SVG regions around arbitrary row spans.

# Goals

- Primary success criteria: Parse common Mermaid sequence fragments, render readable lane-spanning markers, and reject unmatched or unterminated fragment markers with clear errors.
- Secondary success criteria: Keep existing sequence messages, notes, activations, Unicode output, ASCII mode, and full test coverage passing.

# Approach

- Chosen approach: Added fragments as another ordered `SequenceItem` variant. The parser tracks fragment depth while appending start, branch, and end markers; the renderer draws each marker as a full-width separator over the participant lifelines.
- Rejected options: Did not draw multi-row enclosing boxes around fragment contents. That would require a larger layout model that knows fragment spans before rendering and would be brittle for terminal wrapping.

# Implementation

- Architecture / flow: `Parse` now recognizes fragment starts (`loop`, `alt`, `opt`, `par`, plus `critical` and `break`), branch separators (`else` and `and`), and `end`. It validates branches and ends against a depth counter, then preserves fragments in `SequenceDiagram.Fragments` and `SequenceDiagram.Items`.
- Key files or components: `pkg/sequence/parser.go` owns parsing and validation. `pkg/sequence/renderer.go` renders fragment separators using the active ASCII or Unicode character set. `pkg/sequence/sequence_test.go` covers valid rendering and invalid fragment structure. `README.md` and `TODO.md` document the new support.
- Example: `loop Retry` renders as a top separator labeled `loop Retry`; `else Failure` renders as a branch separator; `end` renders as a closing separator.

# Results

- Outputs: Sequence diagrams can now include ordered fragment markers interleaved with messages, notes, and activation state.
- Metrics or observations: The fragment task closes the priority-1 `expand-sequence-core` TODO slice.
- Verification: `go test ./pkg/sequence -run 'TestSequenceFragments|TestSequenceFragmentValidation|TestSequenceActivationShorthand|TestParse'` passed. `go test ./...` passed.

# Decisions

- Tradeoffs made: Fragment separators are terminal-stable and deterministic, but they are approximations. The renderer prioritizes readability and source-order fidelity over exact Mermaid browser layout parity.

# Limitations

- Known issues, uncertainties, or risks: Fragment contents are not visually enclosed by vertical side borders across their full height. Other Mermaid sequence features such as additional arrow families and browser links/actions remain outside this slice.

# Next steps

1. Implement `extract-style-and-color-model` so upcoming diagram renderers can reuse color/style parsing instead of duplicating graph-specific logic.
2. Start a bounded renderer slice for one priority-2 diagram type after the shared style model exists.

# Reproducibility

1. Run `go test ./pkg/sequence -run 'TestSequenceFragments|TestSequenceFragmentValidation|TestSequenceActivationShorthand|TestParse'`.
2. Run `go test ./...`.
