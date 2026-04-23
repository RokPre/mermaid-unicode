# Shared Style Model Extraction Report

# Context

- Problem: Graph color parsing and rendering used graph-local style structures and helper functions, which would force future ER, class, state, and other diagram renderers to duplicate style parsing.
- Constraints: The change needed to preserve existing graph output and avoid creating a broad theming framework before additional renderers exist.

# Goals

- Primary success criteria: Provide a reusable style model outside `cmd/parse.go`, keep existing graph color behavior passing, and add focused tests for style parsing and precedence.
- Secondary success criteria: Keep ANSI/HTML wrapping behavior compatible with the existing `StyleType` setting and avoid forcing unsupported CSS behavior onto renderers.

# Approach

- Chosen approach: Added a small shared style API in `pkg/diagram` and adapted graph code to use it through a local alias. The graph parser and renderer still own graph-specific semantics, but the style map, class, normalization, precedence, and text wrapping helpers are reusable.
- Rejected options: Did not introduce a full theme engine or renderer-wide style registry. Current diagram support only needs a small stable base for upcoming renderer slices.

# Implementation

- Architecture / flow: `diagram.StyleMap` represents parsed CSS-like key/value styles. `diagram.StyleClass` names a style map. `diagram.ParseStyleMap`, `diagram.ResolveStyle`, `diagram.NormalizeStyleColor`, and `diagram.WrapTextInStyle` provide shared parsing, precedence, normalization, and output wrapping.
- Key files or components: `pkg/diagram/style.go` contains the shared model. `cmd/parse.go` delegates `classDef` and `linkStyle` parsing to the shared parser. `cmd/color.go` reads exported `StyleClass.Styles`. `cmd/draw.go` delegates color wrapping and normalization to `pkg/diagram`.
- Example: `classDef warning stroke:#ff0000,color:#00ff00,fill:#111111` is now parsed through `diagram.ParseStyleMap` and still renders the same graph ANSI/HTML styling as before.

# Results

- Outputs: A shared style package is available to future diagram renderers without importing graph parser internals.
- Metrics or observations: The priority-2 `extract-style-and-color-model` TODO is complete with no graph color behavior change expected.
- Verification: `go test ./pkg/diagram -run 'TestParseStyleMap|TestResolveStyle|TestWrapTextInStyleHTML|TestNormalizeStyleColor'` passed. `go test ./cmd -run 'TestMermaidFileToMapParsesClassAndLinkStyles|TestRenderGraphAppliesNodeAndEdgeColors'` passed. `go test ./...` passed.

# Decisions

- Tradeoffs made: `ResolveStyle` treats `none` and `transparent` as absent values rather than overriding lower-precedence styles. That matches the terminal renderer's current capability: it can apply colors, but it does not model full CSS cascade removal semantics.

# Limitations

- Known issues, uncertainties, or risks: Existing graph nodes still only receive class-based styles, and link styles still apply by graph edge index. Direct per-node `style` statements are not implemented in this slice.

# Next steps

1. Implement the priority-2 ER renderer using the shared style model where Mermaid styles overlap with graph styling.
2. Add direct style statement support later if upcoming diagram renderers need it.

# Reproducibility

1. Run `go test ./pkg/diagram -run 'TestParseStyleMap|TestResolveStyle|TestWrapTextInStyleHTML|TestNormalizeStyleColor'`.
2. Run `go test ./cmd -run 'TestMermaidFileToMapParsesClassAndLinkStyles|TestRenderGraphAppliesNodeAndEdgeColors'`.
3. Run `go test ./...`.
