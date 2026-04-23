# Sequence Actor Declaration Report

# Context

- Problem: Mermaid sequence diagrams allow `actor` declarations, but the parser only accepted `participant` declarations and implicit participants from messages.
- Constraints: The current terminal sequence renderer has one participant-box layout and no dedicated actor glyph. The smallest useful implementation was to parse actors and render them through that existing layout while preserving ids and aliases.

# Goals

- Primary success criteria: Support `actor Alice`, `actor A as Alice`, and quoted actor ids in sequence diagrams.
- Secondary success criteria: Keep existing participant parsing unchanged, document the support in README, and pass the full Go test suite.

# Approach

- Chosen approach: Added an `actorRegex` with the same id/alias capture behavior as `participantRegex`, then let `parseParticipant` accept either declaration kind.
- Rejected options: Did not add a separate actor renderer yet. That would require a design decision for terminal actor glyphs and is not necessary for parsing Mermaid sequence examples correctly.

# Implementation

- Architecture / flow: During sequence parsing, each non-message line is checked against participant declarations. The same code path now accepts both `participant` and `actor`, creates a `Participant`, and stores the id/label in the participant map.
- Key files or components: `pkg/sequence/parser.go` adds `actorRegex`. `pkg/sequence/sequence_test.go` adds parser and render coverage for actors. `README.md` lists actor declarations under supported sequence features. `TODO.md` records `add-sequence-actor-declarations`.
- Example: `sequenceDiagram\nactor A as Alice\nA->>A: Think` now renders a sequence participant labeled `Alice`.

# Results

- Outputs: Actor declarations and aliases now parse and render through the existing sequence layout.
- Metrics or observations: Existing participant aliases and implicit participants still pass focused tests.
- Verification: Ran `go test ./pkg/sequence -run 'TestActor|TestParticipantAlias|TestParse'`; it passed. Ran `go test ./...`; it passed.

# Decisions

- Tradeoffs made: Actors currently render like participants. This is a readable terminal approximation and keeps the parser compatible with Mermaid syntax without introducing an unproven visual convention.

# Limitations

- Known issues, uncertainties, or risks: The renderer does not visually distinguish actors from participants. Sequence notes, activation bars, and fragment blocks remain unsupported.

# Next steps

1. Continue `expand-sequence-core` with a small renderer-visible feature such as `Note left of`, `Note right of`, and `Note over`.
2. Decide later whether actors need distinct terminal glyphs after notes/fragments are in place.

# Reproducibility

1. Run `go test ./pkg/sequence -run 'TestActor|TestParticipantAlias|TestParse'`.
2. Run `go test ./...`.
