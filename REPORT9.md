# Sequence Activation Directive Report

# Context

- Problem: Mermaid sequence diagrams support activation spans, but the renderer did not parse `activate` or `deactivate` directives.
- Constraints: The renderer is line-oriented and already has participant lifelines. The smallest useful slice was standalone directives, not Mermaid's message suffix shorthand such as `A->>+B`.

# Goals

- Primary success criteria: Parse `activate A` and `deactivate A` as ordered sequence events.
- Secondary success criteria: Render visible activation bars, keep messages and notes working, document the feature, and pass the full Go test suite.

# Approach

- Chosen approach: Added activation events to the existing ordered `SequenceItem` model introduced for notes. Rendering tracks active participant indexes while walking items and swaps the normal lifeline glyph for an activation glyph while active.
- Rejected options: Did not implement `+` and `-` message suffix shorthand in this slice because it requires expanding message arrow parsing and activation state changes around messages.

# Implementation

- Architecture / flow: `Parse` recognizes `activate` and `deactivate` before notes and messages. `Render` maintains an active participant map and passes it to lifeline, message, self-message, and note rendering helpers.
- Key files or components: `pkg/sequence/parser.go` defines `Activation` and `activationRegex`. `pkg/sequence/renderer.go` renders activation bars. `pkg/sequence/charset.go` adds Unicode and ASCII activation glyphs. `pkg/sequence/sequence_test.go` covers parsing and output. `README.md` and `TODO.md` record the new feature.
- Example: `activate B` marks Bob's lifeline with `┃` in Unicode mode until `deactivate B`.

# Results

- Outputs: Standalone activation directives now render visibly in sequence diagrams.
- Metrics or observations: Existing sequence golden tests still pass because inactive diagrams render exactly as before.
- Verification: Ran `go test ./pkg/sequence -run 'TestSequenceActivations|TestActivationRegex|TestSequenceNotes|TestParse'`; it passed. Ran `go test ./...`; it passed.

# Decisions

- Tradeoffs made: ASCII activation bars use `#`; Unicode activation bars use `┃`. This keeps activation visibly distinct without adding a wider visual block that would disturb terminal alignment.

# Limitations

- Known issues, uncertainties, or risks: Mermaid's activation shorthand on message arrows, such as `->>+` and `-->>-`, is still unsupported. Nested activation depth is represented as active/inactive rather than stacked bars.

# Next steps

1. Continue `expand-sequence-core` with activation shorthand on message arrows.
2. After shorthand, consider loop/alt/opt/par frames as the next larger sequence feature.

# Reproducibility

1. Run `go test ./pkg/sequence -run 'TestSequenceActivations|TestActivationRegex|TestSequenceNotes|TestParse'`.
2. Run `go test ./...`.
