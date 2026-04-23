# Sequence Activation Shorthand Report

# Context

- Problem: Standalone sequence activation directives were supported, but Mermaid's common message-arrow activation shorthand such as `A->>+B` and `B-->>-A` still failed to parse.
- Constraints: The renderer already tracks activation state while walking ordered sequence items. The shorthand needed to integrate with that state without changing inactive message rendering.

# Goals

- Primary success criteria: Support `->>+`, `-->>+`, `->>-`, and `-->>-` on currently supported sequence message arrows.
- Secondary success criteria: Keep standalone activation directives, notes, and normal messages working, document the shorthand, and pass the full Go test suite.

# Approach

- Chosen approach: Extended message parsing with an optional `+` or `-` suffix immediately after the arrow token. `+` creates a post-message activation for the receiver. `-` creates a post-message deactivation for the sender.
- Rejected options: Did not create separate synthetic activation sequence items around shorthand messages. Keeping activation changes on the message preserves source order while avoiding extra rendered lifeline rows.

# Implementation

- Architecture / flow: `messageRegex` now captures an optional activation suffix. `Message` has post-message activation/deactivation fields. Rendering applies those activation state changes immediately after drawing each message.
- Key files or components: `pkg/sequence/parser.go` parses the suffix and stores activation records. `pkg/sequence/renderer.go` applies message activation changes. `pkg/sequence/sequence_test.go` covers parsing and rendered activation output. `README.md` and `TODO.md` record support.
- Example: `A->>+B: Start` activates `B` after the message. `B-->>-A: Done` renders while `B` is active, then deactivates `B`.

# Results

- Outputs: Message activation shorthand now works for solid and dotted message arrows.
- Metrics or observations: Existing sequence tests still pass, including standalone activation and note coverage.
- Verification: Ran `go test ./pkg/sequence -run 'TestSequenceActivationShorthand|TestMessageRegex|TestSequenceActivations|TestParse'`; it passed. Ran `go test ./...`; it passed.

# Decisions

- Tradeoffs made: The `-` suffix deactivates the sender after the message, matching the common Mermaid response pattern `B-->>-A`.

# Limitations

- Known issues, uncertainties, or risks: Only the currently supported message arrows, `->>` and `-->>`, have shorthand support. Other Mermaid sequence arrow families remain unsupported.

# Next steps

1. Continue `expand-sequence-core` with loop/alt/opt/par fragments if a larger sequence renderer slice is acceptable.
2. Add a README sequence support matrix once fragment support is clearer.

# Reproducibility

1. Run `go test ./pkg/sequence -run 'TestSequenceActivationShorthand|TestMessageRegex|TestSequenceActivations|TestParse'`.
2. Run `go test ./...`.
