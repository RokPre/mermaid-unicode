# Diagram Registry Implementation Report

# Context

- Problem: `DiagramFactory` defaulted unknown input to `GraphDiagram`, so a known Mermaid diagram such as `classDiagram` could be misdetected as graph and fail later with a graph parse error.
- Constraints: The change needed to preserve the current supported renderers, sequence and graph/flowchart, while making future diagram families explicit. Existing graph files can start with `paddingX=` or `paddingY=` directives before the graph header, so detection still needed to skip those directives.

# Goals

- Primary success criteria: Add a small diagram registry that detects supported diagrams and rejects known but unimplemented Mermaid diagram types with clear errors.
- Secondary success criteria: Cover graph, flowchart, sequence, comments/blank lines, graph padding directives, known unsupported types, unknown input, and empty/comment-only input with tests.

# Approach

- Chosen approach: Kept the registry local to `cmd/diagram.go` because only the command factory needs it today. The registry stores type name, detector, and constructor for supported diagrams, plus a separate unsupported registry for known Mermaid syntax families.
- Rejected options: Did not introduce a package-level plugin interface yet. That would be premature before any third renderer exists and would add churn outside the current task.

# Implementation

- Architecture / flow: `DiagramFactory` trims input, finds the first real diagram line, checks the supported registry, checks the known unsupported registry, and then returns an unknown-type error if nothing matches.
- Key files or components: `cmd/diagram.go` now contains `supportedDiagramRegistry`, `unsupportedDiagramRegistry`, keyword helpers, and supported-type error formatting. `cmd/diagram_test.go` adds focused factory tests. `TODO.md` marks `add-diagram-registry` done.
- Example: `classDiagram` now returns an unsupported diagram error naming `classDiagram` and supported types `sequence, graph` instead of returning `GraphDiagram`.

# Results

- Outputs: Known future Mermaid types are explicitly detected: `classDiagram`, `stateDiagram`, `stateDiagram-v2`, `erDiagram`, `journey`, `gantt`, `pie`, `quadrantChart`, `requirementDiagram`, `gitGraph`, `mindmap`, `timeline`, and `zenuml`.
- Metrics or observations: Existing sequence and graph rendering paths still pass the full test suite. Graph detection still supports comments, blank lines, `graph`, `flowchart`, and pre-header padding directives.
- Verification: Ran `go test ./cmd -run 'TestDiagramFactory|TestSequenceDiagramIntegration_ErrorHandling'`; it passed. Ran `go test ./...`; it passed.

# Decisions

- Tradeoffs made: The supported type list reports `sequence, graph` because those are the actual implemented renderer type names. `flowchart` is treated as graph because it is rendered by the existing `GraphDiagram`.

# Limitations

- Known issues, uncertainties, or risks: The registry only detects additional Mermaid families; it does not implement rendering for them. Unknown free-form text now fails during detection instead of falling through to graph parsing, which is intentional but changes where the error is raised.

# Next steps

1. Start `audit-flowchart-parity` to close high-value flowchart gaps such as `BT`/`RL` directions and selected expanded shape aliases.
2. Start `expand-sequence-core` to add the next Mermaid sequence features on top of `pkg/sequence`.

# Reproducibility

1. Run `go test ./cmd -run 'TestDiagramFactory|TestSequenceDiagramIntegration_ErrorHandling'`.
2. Run `go test ./...`.
