# Flowchart Support Matrix Report

# Context

- Problem: Flowchart support had advanced through several implementation slices, but README still described stale limitations such as missing subgraph support, missing non-rectangle shapes, and only `LR`/`TD` directions.
- Constraints: This was a documentation closure task for the current terminal-renderer subset, not a new parser feature. The docs needed to be specific about supported syntax and explicit about unsupported Mermaid browser/SVG behavior.

# Goals

- Primary success criteria: Update README so users can see the supported flowchart directions, shape syntax, connector syntax, styling features, and unsupported boundaries.
- Secondary success criteria: Mark the current `audit-flowchart-parity` TODO as done and keep verification green.

# Approach

- Chosen approach: Added a compact support matrix under the existing Unicode graph styling section and refreshed the older supported-diagram checklist.
- Rejected options: Did not add separate full documentation files because the user-facing quick reference belongs in README and the repo already has TODO/report files for internal handoff.

# Implementation

- Architecture / flow: README now documents expanded `@{ shape: ... }` metadata, `LR`/`RL`/`TD`/`TB`/`BT` directions, supported connectors, edge labels, colors, subgraphs, and unsupported browser-only or SVG-exact behavior.
- Key files or components: `README.md` contains the support matrix. `TODO.md` moves `audit-flowchart-parity` to done and records `document-flowchart-support-matrix`.
- Example: The README now lists `A --- B`, `A ---|label| B`, and selected expanded shape aliases next to the older Mermaid node and edge syntax.

# Results

- Outputs: Flowchart docs now match the implemented terminal subset after the direction, shape, and connector work.
- Metrics or observations: No code behavior changed in this task.
- Verification: Ran `go test ./...`; it passed.

# Decisions

- Tradeoffs made: The docs explicitly state that circle/cross link heads and exact Mermaid SVG parity remain unsupported. This is clearer than implying the terminal renderer is a full Mermaid replacement.

# Limitations

- Known issues, uncertainties, or risks: README now documents the flowchart subset, but it does not yet provide a complete support matrix for every other Mermaid diagram family researched in `TODO.md`.

# Next steps

1. Start `expand-sequence-core` with a small sequence feature such as actor declarations or notes.
2. Keep broader diagram-family support matrices for later, after more renderers exist.

# Reproducibility

1. Run `go test ./...`.
