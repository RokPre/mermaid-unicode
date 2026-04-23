# TODO

## Current TODOs

No current task selected.

## Active TODOs

- [ ] 1 Sequence expand-sequence-core
- [ ] 2 Shared extract-style-and-color-model
- [ ] 2 ER add-er-diagram-renderer
- [ ] 2 Class add-class-diagram-renderer
- [ ] 2 State add-state-diagram-renderer
- [ ] 3 Requirement add-requirement-diagram-renderer
- [ ] 3 Mindmap add-mindmap-tree-renderer
- [ ] 4 Timeline add-timeline-renderer
- [ ] 4 UserJourney add-user-journey-renderer
- [ ] 4 Charts add-pie-and-quadrant-renderers
- [ ] 5 Gantt add-gantt-timeline-renderer
- [ ] 5 Gitgraph add-gitgraph-lane-renderer
- [ ] 7 ZenUML defer-zenuml-subset
- [ ] 8 Docs document-supported-mermaid-subsets

## Done TODOs

- [x] 0 Architecture add-diagram-registry
- [x] 1 Flowchart audit-flowchart-parity
- [x] 1 Sequence add-sequence-actor-declarations
- [x] 1 Flowchart add-expanded-shape-syntax
- [x] 1 Flowchart add-reverse-flowchart-directions
- [x] 1 Flowchart add-open-solid-connectors
- [x] 1 Docs document-flowchart-support-matrix
- [x] 0 Graph existing-flowchart-renderer
- [x] 0 Sequence existing-sequence-renderer
- [x] 1 Graph add-unicode-graph-rendering
- [x] 1 Graph prioritize-styled-arrows
- [x] 1 Styling add-graph-colors

## Task Details

### add-diagram-registry

Priority: 0
Area: Architecture
Status: done
Depends on: none

Goal:
Replace the hard-coded diagram detection fallback with an explicit registry for supported diagram families and clear errors for unsupported Mermaid syntax.

Context:
`cmd/diagram.go` currently detects sequence diagrams, then flowchart/graph diagrams, then falls back to `GraphDiagram` for unknown input. That behavior will become wrong as ER, class, state, and other diagram families are added.

Expected changes:
- Add a small diagram descriptor or registry with detector, constructor, and type name fields.
- Register the current sequence and graph/flowchart handlers.
- Detect documented Mermaid keywords such as `classDiagram`, `stateDiagram`, `erDiagram`, `journey`, `gantt`, `pie`, `quadrantChart`, `requirementDiagram`, `gitGraph`, `mindmap`, `timeline`, and `zenuml`.
- Return a clear unsupported-diagram error for known-but-unimplemented diagram types.
- Add tests covering graph, sequence, unsupported known types, comments/blank lines, and unknown input.

Acceptance criteria:
- `go test ./...` passes.
- Existing graph and sequence examples still render.
- A known but unsupported type such as `classDiagram` does not fall through to graph parsing.
- Error text names the unsupported diagram type and the currently supported types.

Notes:
Keep the registry small and local unless adding a package-level abstraction clearly reduces duplication.

### audit-flowchart-parity

Priority: 1
Area: Flowchart
Status: done
Depends on: add-diagram-registry

Goal:
Audit Mermaid flowchart syntax against the current graph parser and implement the highest-value missing terminal-friendly features.

Context:
The renderer already supports graph/flowchart, Unicode boxes, rounded/stadium/double/database/circle/decision/hexagon/parallelogram approximations, styled edges, labels, subgraph parsing, class colors, and directions `TD`, `TB`, and `LR`. Mermaid docs also include `BT`, `RL`, expanded `@{ shape: ... }` syntax, many shape aliases, additional edge heads, and richer subgraph behavior.

Expected changes:
- Add `BT` and `RL` direction support if feasible with the current layout model.
- Add parser support for a focused subset of `A@{ shape: ... }` aliases that map to existing terminal shapes.
- Verify support for `-->`, `---`, `==>`, `-.->`, `-.-`, labels, chained links, and node class shorthand.
- Improve or explicitly document unsupported link heads such as circle and cross heads.
- Add tests and golden fixtures for every supported operator or shape added.

Acceptance criteria:
- A flowchart parity table exists in tests or docs showing supported and unsupported syntax.
- New supported flowchart syntax renders in Unicode and ASCII mode.
- Unsupported flowchart syntax fails clearly or is documented as ignored behavior.
- `go test ./...` passes.

Notes:
Completed as a terminal-friendly flowchart parity pass. Remaining work such as circle/cross link heads or exact Mermaid SVG shape parity should be tracked as separate future tasks if it becomes necessary.

### document-flowchart-support-matrix

Priority: 1
Area: Docs
Status: done
Depends on: audit-flowchart-parity

Goal:
Document the supported flowchart subset and the intentional Mermaid browser-renderer differences in README.

Context:
Flowchart support has expanded to include Unicode shapes, selected expanded shape metadata, reverse directions, styled connectors, colors, and subgraphs. The README still had stale checklist entries that said several of those features were unsupported.

Expected changes:
- Update the graph/flowchart styling section with supported shape aliases, connectors, directions, and unsupported boundaries.
- Update the supported diagram checklist so it matches current code behavior.
- Keep unsupported browser-only behavior and exact SVG parity clearly outside the current terminal scope.

Acceptance criteria:
- README lists the current supported flowchart directions.
- README lists current supported connector families and selected expanded shape metadata.
- README calls out unsupported circle/cross link heads and browser interactivity.

Notes:
This closes the current `audit-flowchart-parity` task at the documented terminal-renderer subset.

### add-sequence-actor-declarations

Priority: 1
Area: Sequence
Status: done
Depends on: add-diagram-registry

Goal:
Support Mermaid sequence `actor` declarations and aliases.

Context:
The sequence parser already supported `participant` declarations, aliases, and implicit participants from messages. Mermaid sequence diagrams also allow `actor` declarations. The current terminal renderer does not have a separate actor glyph, so actors are rendered through the same participant box layout while preserving ids and labels.

Expected changes:
- Parse `actor Alice`.
- Parse `actor A as Alice`.
- Parse quoted actor ids such as `actor "External User" as User`.
- Render actor labels in the existing sequence participant row.
- Add parser and render tests.

Acceptance criteria:
- Actor declarations create participants with the expected id and label.
- Actor aliases render the alias label.
- `go test ./...` passes.

Notes:
This is a first bounded slice of `expand-sequence-core`; visual actor-specific glyphs can be a later renderer enhancement if needed.

### add-expanded-shape-syntax

Priority: 1
Area: Flowchart
Status: done
Depends on: add-diagram-registry

Goal:
Support Mermaid's expanded flowchart node metadata syntax for shape aliases that map cleanly to existing terminal node shapes.

Context:
Mermaid v11.3.0+ supports `A@{ shape: rect }` style shape metadata. The renderer already had terminal approximations for square, rounded, stadium, double/subroutine, database, circle, decision, hexagon, and parallelogram shapes through older delimiter syntax.

Expected changes:
- Parse `@{ ... }` node metadata without treating it as part of the node id.
- Map supported `shape:` aliases onto existing terminal shapes.
- Preserve optional `label:` metadata and class shorthand such as `:::important`.
- Leave unsupported expanded shapes as normal square nodes with the clean node id.
- Add parser and render tests.

Acceptance criteria:
- `A@{ shape: rounded, label: "Research" }` renders with the rounded box glyphs and label `Research`.
- Supported aliases for existing shape families parse into the expected `graphNodeShape`.
- Unsupported expanded shapes do not corrupt the node id.
- `go test ./...` passes.

Notes:
This completes one slice of `audit-flowchart-parity`; `BT`/`RL`, extra link heads, and a support matrix remain active follow-ups.

### add-reverse-flowchart-directions

Priority: 1
Area: Flowchart
Status: done
Depends on: add-diagram-registry

Goal:
Support Mermaid `graph RL`, `flowchart RL`, `graph BT`, and `flowchart BT` directions in the terminal graph renderer.

Context:
The renderer previously accepted `LR`, `TD`, and `TB`, while Mermaid also documents `RL` and `BT`. The existing layout model already separates horizontal and vertical flow, so reverse directions can be represented by mirroring the mapped grid before pathing and drawing.

Expected changes:
- Parse `RL` and `BT` graph/flowchart headers.
- Keep horizontal and vertical layout decisions orientation-aware instead of checking only `LR` or `TD`.
- Mirror node grid coordinates for reversed directions before calculating paths, sizes, and subgraph bounds.
- Render `RL` arrows leftward and `BT` arrows upward.
- Add parser and render tests.

Acceptance criteria:
- `graph RL` renders the child left of the parent with a left arrowhead.
- `graph BT` renders the child above the parent with an up arrowhead.
- Existing `LR`, `TD`, and `TB` behavior remains covered by the full test suite.
- `go test ./...` passes.

Notes:
This completes the reverse-direction slice of `audit-flowchart-parity`. Complex subgraph layouts in reversed directions may still need golden fixtures later.

### add-open-solid-connectors

Priority: 1
Area: Flowchart
Status: done
Depends on: add-diagram-registry

Goal:
Support Mermaid open solid flowchart connectors using `---` and labeled `---|label|` syntax.

Context:
The renderer already supported arrow connectors `-->`, heavy arrows `==>`, dashed arrows `-.->`, and open dashed connectors `-.-`. Mermaid flowcharts also use `---` for links without arrowheads.

Expected changes:
- Parse `A --- B` as a light connector without an arrowhead.
- Parse `A ---|label| B` as a labeled light connector without an arrowhead.
- Preserve existing arrow and dashed connector behavior.
- Add parser and render tests proving open solid connectors do not draw arrowheads.

Acceptance criteria:
- `A --- B` renders a solid line between nodes with no arrowhead.
- `A ---|open| B` preserves the label and still omits the arrowhead.
- `go test ./...` passes.

Notes:
This completes one more connector slice under `audit-flowchart-parity`; bidirectional, circle-head, and cross-head operators remain future work.

### expand-sequence-core

Priority: 1
Area: Sequence
Status: pending
Depends on: add-diagram-registry

Goal:
Expand the existing sequence renderer toward the Mermaid sequence core: participants, actors, aliases, message variants, notes, activations, and common fragments.

Context:
Mermaid sequence docs include participants/actors, aliases, messages, notes, activation/deactivation, loops, `alt`/`else`, `opt`, `par`, and related fragments. The repo already has `pkg/sequence`, so sequence work should extend that package instead of creating a new renderer.

Expected changes:
- Audit current `pkg/sequence` parser and renderer behavior against the official syntax.
- Add missing participant/actor/alias handling.
- Add or verify message arrow variants and labels.
- Add notes over, left of, and right of participants if missing.
- Add activation/deactivation bars if missing.
- Add loop, alt/else, opt, and par frames as Unicode/ASCII lane-spanning blocks.
- Add example tests and focused parser tests.

Acceptance criteria:
- Mermaid examples for the supported subset render without parse errors.
- Unsupported sequence features fail clearly or are documented as unsupported.
- Unicode and ASCII output remain deterministic.
- `go test ./...` passes.

Notes:
Browser-specific links/actions should remain out of scope.

### extract-style-and-color-model

Priority: 2
Area: Shared
Status: pending
Depends on: add-diagram-registry

Goal:
Make Mermaid style/class/color parsing reusable across graph-like and non-graph renderers.

Context:
Graph color support currently lives in graph-specific parsing and rendering paths. ER, class, state, requirement, quadrant, timeline, gitgraph, and user journey all have style or color concepts in the docs.

Expected changes:
- Identify style parsing code that can move out of graph-specific types without destabilizing graph rendering.
- Define a shared style model for fill, stroke, text color, and emphasis where terminal output can represent them.
- Keep ANSI behavior optional and compatible with ASCII output.
- Add tests for style precedence: direct style, class style, default style, and unstyled output.

Acceptance criteria:
- Existing graph color tests still pass.
- At least one shared style helper is usable outside `cmd/parse.go`.
- No renderer is forced to support CSS properties it cannot render.
- `go test ./...` passes.

Notes:
Prefer small extraction over a large theming framework.

### add-er-diagram-renderer

Priority: 2
Area: ER
Status: pending
Depends on: add-diagram-registry, extract-style-and-color-model

Goal:
Add a terminal renderer for a practical Mermaid `erDiagram` subset.

Context:
ER diagrams are a strong fit for terminal boxes: entities become boxes, attributes become rows, and relationships become labeled lines with cardinality markers.

Expected changes:
- Add `erDiagram` detection through the registry.
- Parse entity declarations, relationships, labels, attributes, aliases, PK/FK/UK markers, and direction.
- Render entity boxes with attribute compartments.
- Render identifying and non-identifying relationships with solid and dashed connectors.
- Approximate crow's foot cardinality markers in ASCII and Unicode.

Acceptance criteria:
- Official-doc style ER examples for entities, attributes, and relationships render.
- Invalid relationship syntax produces actionable parse errors.
- Unicode and ASCII fixtures cover cardinality markers.
- `go test ./...` passes.

Notes:
Exact crow's foot geometry is not required; stable readable markers are.

### add-class-diagram-renderer

Priority: 2
Area: Class
Status: pending
Depends on: add-diagram-registry, extract-style-and-color-model

Goal:
Add a terminal renderer for a practical Mermaid `classDiagram` subset.

Context:
Class diagrams map well to compartment boxes and graph relationships. The docs include class declarations, members, visibility, relationships, labels, two-way relations, lollipop interfaces, namespaces, cardinality, annotations, and styling.

Expected changes:
- Add `classDiagram` detection through the registry.
- Parse class declarations, labels, member blocks, colon member declarations, visibility markers, and method detection via `()`.
- Render classes as boxes with name, attributes, and operations compartments.
- Parse and render core relationship operators: `<|--`, `*--`, `o--`, `-->`, `--`, `..>`, `..|>`, and `..`.
- Support relationship labels and cardinality strings.

Acceptance criteria:
- Basic class, members, relationships, and labels render in Unicode and ASCII.
- Unsupported namespace/lollipop syntax either renders in a documented approximation or returns a clear error.
- `go test ./...` passes.

Notes:
Do not let class parser changes alter flowchart edge parsing.

### add-state-diagram-renderer

Priority: 2
Area: State
Status: pending
Depends on: add-diagram-registry, extract-style-and-color-model

Goal:
Add a terminal renderer for a practical Mermaid `stateDiagram`/`stateDiagram-v2` subset.

Context:
State diagrams are graph-like but have unique start/end markers, composite states, choice/fork/join markers, notes, and concurrency syntax.

Expected changes:
- Add state diagram detection through the registry.
- Parse state declarations, descriptions, transitions, transition labels, comments, and direction.
- Render start/end states using compact symbols.
- Render simple states and transitions via the graph drawing infrastructure where practical.
- Add first-pass support for choice, fork, join, and notes.
- Add composite state frames after simple states are stable.

Acceptance criteria:
- Simple state diagrams from the docs render.
- Start/end markers and transition labels are covered by tests.
- Composite syntax is either supported in a focused subset or rejected clearly.
- `go test ./...` passes.

Notes:
Nested composite layout can be iterative; avoid promising full depth parity immediately.

### add-requirement-diagram-renderer

Priority: 3
Area: Requirement
Status: pending
Depends on: add-diagram-registry, extract-style-and-color-model

Goal:
Add a terminal renderer for Mermaid `requirementDiagram`.

Context:
Requirement diagrams define structured requirement and element blocks plus labeled relationships. They fit terminal boxes and graph connectors well.

Expected changes:
- Parse requirement blocks with type, name, id, text, risk, and verification method.
- Parse element blocks with name, type, and docref.
- Parse relationships with contains, copies, derives, satisfies, verifies, refines, and traces.
- Render requirement and element boxes with readable field rows.
- Render relationships with labels.

Acceptance criteria:
- The larger official-doc style example renders in a deterministic terminal layout.
- Missing required fields produce useful parse errors.
- Styling/classes work if shared style extraction is available.
- `go test ./...` passes.

Notes:
Keep SysML semantics descriptive; the renderer should not validate requirement quality.

### add-mindmap-tree-renderer

Priority: 3
Area: Mindmap
Status: pending
Depends on: add-diagram-registry

Goal:
Add a terminal tree renderer for Mermaid `mindmap`.

Context:
Mindmap syntax is indentation-based and experimental. A terminal tree is more realistic than a radial mindmap and still preserves the hierarchy.

Expected changes:
- Add `mindmap` detection through the registry.
- Parse indentation levels into a tree.
- Render with Unicode tree glyphs and ASCII fallback.
- Support a focused set of node shape syntax where it maps cleanly to current graph node shapes.
- Ignore or clearly reject icons in the first pass.

Acceptance criteria:
- Nested mindmap examples render as stable trees.
- Bad indentation returns clear parse errors.
- Unicode and ASCII modes are both covered.
- `go test ./...` passes.

Notes:
Radial layout and icon integration are out of scope for the first implementation.

### add-timeline-renderer

Priority: 4
Area: Timeline
Status: pending
Depends on: add-diagram-registry, extract-style-and-color-model

Goal:
Add a terminal renderer for Mermaid `timeline`.

Context:
Timeline syntax is simple: optional title, time periods, events, sections, optional direction, and color schemes. The docs mark timeline as experimental.

Expected changes:
- Parse `timeline`, optional `LR`/`TD`, title, sections, periods, and one or more events per period.
- Render `LR` as compact columns or rows and `TD` as a vertical chronology.
- Wrap long event text deterministically.
- Reuse shared color/style support where practical.

Acceptance criteria:
- Sectioned and unsectioned timeline examples render.
- Multiple events under one period render in source order.
- Direction handling is tested.
- `go test ./...` passes.

Notes:
Do not implement browser theme variables beyond simple terminal color mapping.

### add-user-journey-renderer

Priority: 4
Area: UserJourney
Status: pending
Depends on: add-diagram-registry

Goal:
Add a terminal renderer for Mermaid `journey`.

Context:
User journey syntax is small: sections and task lines containing a task name, score from 1 to 5, and comma-separated actors.

Expected changes:
- Add `journey` detection through the registry.
- Parse sections and task rows.
- Validate score range 1..5.
- Render a sectioned table with task, score, and actors.
- Optionally show score with proportional bars or colored markers.

Acceptance criteria:
- Official-doc style journey examples render.
- Invalid scores produce useful errors.
- Long actor lists wrap without breaking columns.
- `go test ./...` passes.

Notes:
This renderer does not need graph layout.

### add-pie-and-quadrant-renderers

Priority: 4
Area: Charts
Status: pending
Depends on: add-diagram-registry, extract-style-and-color-model

Goal:
Add practical terminal renderers for Mermaid `pie` and `quadrantChart`.

Context:
Pie charts can be represented better as tables and proportional bars in terminal cells. Quadrant charts can be represented as a fixed grid with x/y points in the 0..1 range.

Expected changes:
- Add `pie` and `quadrantChart` detection through the registry.
- Parse pie title, optional showData, labels, and numeric values.
- Render pie as sorted or source-order labels with values, percentages, and bars.
- Parse quadrant title, axes, quadrant labels, point labels, x/y values, and basic point styling.
- Render quadrant as a deterministic grid with labels.

Acceptance criteria:
- Pie values total and percentages are correct.
- Quadrant points reject values outside 0..1.
- Unicode and ASCII grid/bar output is tested.
- `go test ./...` passes.

Notes:
Do not try to draw circular pies in the terminal.

### add-gantt-timeline-renderer

Priority: 5
Area: Gantt
Status: pending
Depends on: add-diagram-registry

Goal:
Add a simplified Mermaid `gantt` renderer for project timelines.

Context:
Gantt syntax is useful but broader than a simple graph parser. The docs include dates, durations, dependencies, statuses, milestones, excludes, weekend rules, and vertical markers.

Expected changes:
- Parse title, dateFormat for a focused subset, sections, task rows, status tags, task IDs, `after` dependencies, durations, milestones, and vertical markers.
- Render a scaled terminal timeline with section rows and status markers.
- Support `YYYY-MM-DD` first before broadening dateFormat support.
- Add clear errors for unsupported date formats.

Acceptance criteria:
- Sequential tasks, explicit dates, durations, dependencies, and milestones render.
- Date calculations are covered by unit tests.
- Unsupported date format behavior is documented.
- `go test ./...` passes.

Notes:
Full Day.js date parsing and exclude/weekend parity should wait until the core renderer is stable.

### add-gitgraph-lane-renderer

Priority: 5
Area: Gitgraph
Status: pending
Depends on: add-diagram-registry, extract-style-and-color-model

Goal:
Add a useful terminal lane renderer for Mermaid `gitGraph`.

Context:
Gitgraph maps well to terminal lanes but has many details: commits, IDs, tags, types, branches, checkout, merge, cherry-pick, branch ordering, orientations, parallel commits, labels, and colors.

Expected changes:
- Add `gitGraph` detection through the registry.
- Parse commits, IDs, tags, branch, checkout, merge, and a narrow cherry-pick subset.
- Render branches as lanes with commit markers and merge connectors.
- Support `LR` first, then consider `TB` and `BT`.
- Map branch colors through shared style/color helpers.

Acceptance criteria:
- A branch/checkout/merge example renders with correct lane continuity.
- Commit IDs and tags appear in output when present.
- Unsupported configuration flags produce clear errors or documented ignored behavior.
- `go test ./...` passes.

Notes:
Rotated labels and exact temporal placement are not terminal goals.

### defer-zenuml-subset

Priority: 7
Area: ZenUML
Status: pending
Depends on: expand-sequence-core

Goal:
Decide whether ZenUML support should exist at all, and if yes, define a small subset that maps to the sequence model.

Context:
ZenUML uses a different syntax than Mermaid's normal sequence diagram. It includes participants, annotators, aliases, sync/async/create/reply messages, nested blocks, comments, loops, alt, opt, and parallel blocks.

Expected changes:
- Write a short design note or TODO update defining supported and unsupported ZenUML syntax.
- Prefer mapping only participants, aliases, messages, and simple blocks to `pkg/sequence`.
- Keep full nested parser parity out of near-term scope.

Acceptance criteria:
- The repo has a clear decision: unsupported with documented reason, or supported subset with examples.
- If implemented, ZenUML examples do not regress regular sequence rendering.
- `go test ./...` passes if code changes are made.

Notes:
This is intentionally lower priority than regular sequence support.

### document-supported-mermaid-subsets

Priority: 8
Area: Docs
Status: pending
Depends on: audit-flowchart-parity, expand-sequence-core

Goal:
Document exactly which Mermaid syntax subsets this terminal renderer supports and where it intentionally differs from Mermaid's browser renderer.

Context:
The official Mermaid syntax surface is large. Users need clear expectations for terminal rendering, Unicode/ASCII mode, color behavior, and unsupported browser-only features.

Expected changes:
- Update `README.md` with a support matrix by diagram type.
- List supported graph node shapes, arrow types, colors, and sequence features.
- Document unsupported features such as click callbacks, exact SVG layout, radial mindmaps, true circular pies, and full ZenUML.
- Add example commands for each supported diagram family as they are implemented.

Acceptance criteria:
- README support matrix matches code behavior.
- The docs explain Unicode default and ASCII fallback.
- Unsupported known diagram types are documented clearly until implemented.

Notes:
Keep docs synced with the registry so supported type names do not drift.

### existing-flowchart-renderer

Priority: 0
Area: Graph
Status: done
Depends on: none

Goal:
Maintain the existing flowchart/graph renderer.

Context:
The repo already renders `graph` and `flowchart` inputs through `GraphDiagram`.

Expected changes:
- None for this done item.

Acceptance criteria:
- Existing graph tests pass.

Notes:
Future graph work should preserve existing Mermaid graph behavior.

### existing-sequence-renderer

Priority: 0
Area: Sequence
Status: done
Depends on: none

Goal:
Maintain the existing sequence renderer.

Context:
The repo already has `pkg/sequence` and `SequenceDiagram`.

Expected changes:
- None for this done item.

Acceptance criteria:
- Existing sequence tests pass.

Notes:
Future sequence work should extend `pkg/sequence`.

### add-unicode-graph-rendering

Priority: 1
Area: Graph
Status: done
Depends on: existing-flowchart-renderer

Goal:
Use Unicode box drawing by default with ASCII fallback.

Context:
The README and code now describe Unicode as the default and `--ascii`/`-a` as fallback. Graph shapes include square, rounded, stadium, double, database, circle, decision, hexagon, and parallelogram approximations.

Expected changes:
- None for this done item.

Acceptance criteria:
- Unicode output is default.
- ASCII output is available with `--ascii` or `-a`.

Notes:
Keep all new renderers compatible with ASCII mode.

### prioritize-styled-arrows

Priority: 1
Area: Graph
Status: done
Depends on: add-unicode-graph-rendering

Goal:
Ensure higher-priority arrow styles win when drawn cells overlap.

Context:
`REPORT1.md` documents explicit arrow style priority: dashed < light < heavy. Heavy `==>` edges now win over default `-->` edges in overlapping cells.

Expected changes:
- None for this done item.

Acceptance criteria:
- Existing arrow priority tests pass.

Notes:
Future renderers with overlapping line styles should use the same priority principle.

### add-graph-colors

Priority: 1
Area: Styling
Status: done
Depends on: add-unicode-graph-rendering

Goal:
Support colored graph nodes, arrows, and text through Mermaid-style class/style data.

Context:
Recent work added color support for graph output so diagrams can differentiate nodes and edges more clearly.

Expected changes:
- None for this done item.

Acceptance criteria:
- Existing graph color tests pass.

Notes:
Color support should be extracted before many new diagram renderers duplicate it.
