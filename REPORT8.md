# Sequence Note Support Report

# Context

- Problem: Mermaid sequence diagrams support notes, but the renderer only modeled participant declarations and messages.
- Constraints: Notes need to preserve source order with messages. The renderer already has a line-oriented sequence layout and shared ASCII/Unicode charset, so notes should reuse that machinery instead of introducing a separate layout engine.

# Goals

- Primary success criteria: Support `Note left of A: text`, `Note right of A: text`, and `Note over A,B: text`.
- Secondary success criteria: Render notes as terminal boxes, keep messages and notes in source order, document the feature, and pass the full Go test suite.

# Approach

- Chosen approach: Added ordered sequence items that can hold either a message or a note. Messages still populate the existing `Messages` slice for compatibility, while rendering now iterates over ordered items.
- Rejected options: Did not store notes separately and render them after messages, because that would lose Mermaid source order and make notes less useful.

# Implementation

- Architecture / flow: `Parse` now checks note syntax before messages. Parsed notes reference participants through the existing participant map, which means notes can introduce implicit participants just like messages. `Render` uses `SequenceDiagram.Items` to render messages and notes in order.
- Key files or components: `pkg/sequence/parser.go` defines `Note`, `NotePosition`, `SequenceItem`, `noteRegex`, and note parsing helpers. `pkg/sequence/renderer.go` adds note box placement. `pkg/sequence/sequence_test.go` covers note parsing, ordering, and rendered output. `README.md` and `TODO.md` record the new support.
- Example: `Note over A,B: Shared context` now renders a boxed note spanning the two participant lanes.

# Results

- Outputs: Sequence notes render in Unicode and ASCII modes through the existing charset.
- Metrics or observations: A Unicode indexing bug in note overlay placement was caught by the focused note test and fixed by indexing overlay text as runes instead of bytes.
- Verification: Ran `go test ./pkg/sequence -run 'TestSequenceNotes|TestNoteRegex|TestParse'`; it passed after the rune-indexing fix. Ran `go test ./...`; it passed.

# Decisions

- Tradeoffs made: Notes are single-line boxes. That matches the first useful Mermaid subset while leaving multiline note wrapping for later.

# Limitations

- Known issues, uncertainties, or risks: Notes do not yet support multiline text or markdown formatting. Actor-specific visuals, activation bars, and loop/alt/opt/par fragments remain unsupported.

# Next steps

1. Continue `expand-sequence-core` with activation/deactivation bars or loop/alt/opt/par fragments.
2. Add golden fixtures for note rendering if the visual format starts changing frequently.

# Reproducibility

1. Run `go test ./pkg/sequence -run 'TestSequenceNotes|TestNoteRegex|TestParse'`.
2. Run `go test ./...`.
