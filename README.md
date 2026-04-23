# Mermaid ASCII / Unicode

Render Mermaid diagrams in your terminal. Unicode box drawing is the default, including rounded boxes, double-line boxes, decision-node approximations, heavy connectors, and dashed connectors. Use `--ascii` or `-a` when you need plain ASCII output.

## Installation

You can download the binary from Github releases:

```bash
# Get the latest release
$ curl -s https://api.github.com/repos/AlexanderGrooff/mermaid-ascii/releases/latest | grep "browser_download_url.*mermaid-ascii" | grep "$(uname)_$(uname -m)" | cut -d: -f2,3 | tr -d \" | wget -qi -
# Unzip it
$ tar xvzf mermaid-ascii_*.tar.gz
$ ./mermaid-ascii --help
```

You can also build it yourself:

```bash
$ git clone
$ cd mermaid-ascii
$ go build
$ mermaid-ascii --help
```

Or using Nix:
```bash
$ git clone
$ cd mermaid-ascii
$ nix build
$ ./result/bin/mermaid-ascii --help
```

## Usage

You can render graphs directly from the command line or start a web interface to render them interactively.

```bash
$ cat test.mermaid
graph LR
A --> B & C
B --> C & D
D --> C
$ mermaid-ascii --file test.mermaid
┌───┐     ┌───┐     ┌───┐
│   │     │   │     │   │
│ A ├────►│ B ├────►│ D │
│   │     │   │     │   │
└─┬─┘     └─┬─┘     └─┬─┘
  │         │         │  
  │         │         │  
  │         │         │  
  │         │         │  
  │         ▼         │  
  │       ┌───┐       │  
  │       │   │       │  
  └──────►│ C │◄──────┘  
          │   │          
          └───┘          

# Increase horizontal spacing
$ mermaid-ascii --file test.mermaid -x 8
┌───┐        ┌───┐        ┌───┐
│   │        │   │        │   │
│ A ├───────►│ B ├───────►│ D │
│   │        │   │        │   │
└─┬─┘        └─┬─┘        └─┬─┘
  │            │            │  
  │            │            │  
  │            │            │  
  │            │            │  
  │            ▼            │  
  │          ┌───┐          │  
  │          │   │          │  
  └─────────►│ C │◄─────────┘  
             │   │             
             └───┘             

# Increase box padding
$ mermaid-ascii -f ./test.mermaid -p 3
┌───────┐     ┌───────┐     ┌───────┐
│       │     │       │     │       │
│       │     │       │     │       │
│       │     │       │     │       │
│   A   ├────►│   B   ├────►│   D   │
│       │     │       │     │       │
│       │     │       │     │       │
│       │     │       │     │       │
└───┬───┘     └───┬───┘     └───┬───┘
    │             │             │
    │             │             │
    │             │             │
    │             │             │
    │             ▼             │
    │         ┌───────┐         │
    │         │       │         │
    │         │       │         │
    │         │       │         │
    └────────►│   C   │◄────────┘
              │       │
              │       │
              │       │
              └───────┘

# Labeled edges
$ cat test.mermaid
graph LR
A --> B
A --> C
B --> C
B -->|example| D
D --> C
$ mermaid-ascii -f ./test.mermaid
┌───┐     ┌───┐         ┌───┐
│   │     │   │         │   │
│ A ├────►│ B ├─example►│ D │
│   │     │   │         │   │
└─┬─┘     └─┬─┘         └─┬─┘
  │         │             │  
  │         │             │  
  │         │             │  
  │         │             │  
  │         ▼             │  
  │       ┌───┐           │  
  │       │   │           │  
  └──────►│ C │◄──────────┘  
          │   │              
          └───┘              

# Top-down layout
$ cat test.mermaid
graph TD
A --> B
A --> C
B --> C
B -->|example| D
D --> C
$ mermaid-ascii -f ./test.mermaid
┌─────────┐          
│         │          
│    A    ├───────┐  
│         │       │  
└────┬────┘       │  
     │            │  
     │            │  
     │            │  
     │            │  
     ▼            ▼  
┌─────────┐     ┌───┐
│         │     │   │
│    B    ├────►│ C │
│         │     │   │
└────┬────┘     └───┘
     │            ▲  
     │            │  
  example         │  
     │            │  
     ▼            │  
┌─────────┐       │  
│         │       │  
│    D    ├───────┘  
│         │          
└─────────┘          

# Read from stdin
$ cat test.mermaid | mermaid-ascii
┌───┐     ┌───┐     ┌───┐
│   │     │   │     │   │
│ A ├────►│ B ├────►│ D │
│   │     │   │     │   │
└─┬─┘     └─┬─┘     └─┬─┘
  │         │         │  
  │         │         │  
  │         │         │  
  │         │         │  
  │         ▼         │  
  │       ┌───┐       │  
  │       │   │       │  
  └──────►│ C │◄──────┘  
          │   │          
          └───┘          

# Only ASCII
$ cat test.mermaid | mermaid-ascii --ascii
+---+     +---+     +---+
|   |     |   |     |   |
| A |---->| B |---->| D |
|   |     |   |     |   |
+---+     +---+     +---+
  |         |         |
  |         |         |
  |         |         |
  |         |         |
  |         v         |
  |       +---+       |
  |       |   |       |
  ------->| C |<-------
          |   |
          +---+

# Unicode node shapes
$ cat shapes.mermaid
graph LR
A(Rounded)
B[[Double]]
C{Decision}
$ mermaid-ascii -f shapes.mermaid
╭─────────╮
│         │
│ Rounded │
│         │
╰─────────╯
           
           
           
           
           
╔════════╗
║        ║
║ Double ║
║        ║
╚════════╝
           
           
           
           
           
◇────────────◇
│            │
│  Decision  │
│            │
◇────────────◇

# Unicode edge styles
$ cat edges.mermaid
graph LR
A ==> B
B -.-> C
C -.- D
$ mermaid-ascii -f edges.mermaid
┌───┐     ┌───┐     ┌───┐     ┌───┐
│   │     │   │     │   │     │   │
│ A ├━━━━►│ B ├┄┄┄┄►│ C ├┄┄┄┄┄│ D │
│   │     │   │     │   │     │   │
└───┘     └───┘     └───┘     └───┘

# Using Docker
$ docker build -t mermaid-ascii .
$ echo 'sequenceDiagram
Alice->>Bob: Hello
Bob-->>Alice: Hi' | docker run -i mermaid-ascii -f -
┌───────┐     ┌─────┐
│ Alice │     │ Bob │
└───┬───┘     └──┬──┘
    │            │
    │ Hello      │
    ├───────────►│
    │            │
    │ Hi         │
    │◄┈┈┈┈┈┈┈┈┈┈┈┤
    │            │

# Graph diagrams work too
$ echo 'graph LR
A-->B-->C' | docker run -i mermaid-ascii -f -
┌───┐     ┌───┐     ┌───┐
│   │     │   │     │   │
│ A ├────►│ B ├────►│ C │
│   │     │   │     │   │
└───┘     └───┘     └───┘

# Run web interface
$ docker run -p 3001:3001 mermaid-ascii web --port 3001
# Then visit http://localhost:3001
```

### Unicode Graph Styling

Unicode output is the default. Use `--ascii` when you need plain ASCII-only diagrams.

Supported graph node shape mappings:

| Mermaid syntax | Terminal approximation |
| --- | --- |
| `A[Text]` or `A["Text"]` | square box, `┌ ┐ └ ┘` |
| `A(Text)` | rounded box, `╭ ╮ ╰ ╯` |
| `A([Text])` | stadium-like rounded box with extra horizontal padding |
| `A[[Text]]` | double-line box, `╔ ╗ ╚ ╝` |
| `A[(Text)]` | database/cylinder approximation using a rounded box |
| `A((Text))` | circle-like rounded/stadium approximation |
| `A{Text}` | decision approximation with `◇` endpoints |
| `A{{Text}}` | hexagon-like approximation |
| `A[/Text/]` | parallelogram-like approximation |
| `A@{ shape: rounded, label: "Text" }` | expanded Mermaid shape metadata mapped to the same terminal shape set |

Supported expanded `shape:` aliases include the terminal-friendly Mermaid families for rectangles, rounded boxes, stadiums, subroutines, databases, circles, decisions, hexagons, and parallelograms. Unsupported expanded shapes are rendered as normal square nodes with the clean node id instead of treating the metadata as part of the label.

Supported graph edge mappings:

| Mermaid syntax | Terminal connector |
| --- | --- |
| `A --> B` | light line, `────►` |
| `A --- B` | light line without arrowhead, `────` |
| `A ---|label| B` | labeled light line without arrowhead |
| `A ==> B` | heavy line, `━━━━►` |
| `A -.-> B` | dashed line with arrowhead, `┄┄┄┄►` |
| `A -.- B` | dashed line without arrowhead |

Supported graph directions are `LR`, `RL`, `TD`, `TB`, and `BT` for both `graph` and `flowchart`.

Supported flowchart features:

| Feature | Status |
| --- | --- |
| Node labels, including quoted labels and multiline label breaks | Supported |
| Node shape syntax and selected expanded `@{ shape: ... }` aliases | Supported |
| `-->`, `---`, `==>`, `-.->`, and `-.-` connectors | Supported |
| Edge labels using `-->|label|`, `---|label|`, `==>|label|`, and `-.->|label|` | Supported |
| Chained links and `A & B` shorthand | Supported |
| `classDef`, `class`, and `linkStyle` color styling | Supported |
| `subgraph ... end` frames, including nested and labeled subgraphs | Supported |
| Mermaid click callbacks, browser interactivity, circle/cross link heads, and exact SVG shape parity | Not supported |

You can also choose default styles for bare nodes and standard arrows:

```bash
$ mermaid-ascii -f graph.mermaid --box-style rounded --edge-style heavy
```

`--box-style` accepts `square`, `rounded`, `double`, or `heavy`.
`--edge-style` accepts `light`, `heavy`, or `dashed`.
Explicit Mermaid syntax still wins over these defaults, so `A[[Text]]` stays double-line and `A ==> B` stays heavy even when different defaults are configured.

### Sequence Diagrams

Sequence diagrams are also fully supported! They visualize message flows between participants over time.
The renderer supports participant and actor declarations, notes, activation directives and shorthand, and Mermaid fragment frames such as `loop`, `alt`, `else`, `opt`, and `par`.

```bash
# Simple sequence diagram
$ cat sequence.mermaid
sequenceDiagram
Alice->>Bob: Hello Bob!
Bob-->>Alice: Hi Alice!
$ mermaid-ascii -f sequence.mermaid
┌───────┐     ┌─────┐
│ Alice │     │ Bob │
└───┬───┘     └──┬──┘
    │            │
    │ Hello Bob! │
    ├───────────►│
    │            │
    │ Hi Alice!  │
    │◄┈┈┈┈┈┈┈┈┈┈┈┤
    │            │

# Solid arrows (->>) and dotted arrows (-->>)
$ cat sequence.mermaid
sequenceDiagram
Client->>Server: Request
Server-->>Client: Response
$ mermaid-ascii -f sequence.mermaid
┌────────┐     ┌────────┐
│ Client │     │ Server │
└───┬────┘     └───┬────┘
    │              │
    │   Request    │
    ├─────────────►│
    │              │
    │   Response   │
    │◄┈┈┈┈┈┈┈┈┈┈┈┈┈┤
    │              │

# Multiple participants
$ cat sequence.mermaid
sequenceDiagram
Alice->>Bob: Hello!
Bob->>Charlie: Forward message
Charlie-->>Alice: Got it!
$ mermaid-ascii -f sequence.mermaid
┌───────┐     ┌─────┐     ┌─────────┐
│ Alice │     │ Bob │     │ Charlie │
└───┬───┘     └──┬──┘     └────┬────┘
    │            │              │
    │   Hello!   │              │
    ├───────────►│              │
    │            │              │
    │            │ Forward message
    │            ├─────────────►│
    │            │              │
    │         Got it!           │
    │◄┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┤
    │            │              │

# Self-messages
$ cat sequence.mermaid
sequenceDiagram
Alice->>Alice: Think
Alice->>Bob: Hello
$ mermaid-ascii -f sequence.mermaid
┌───────┐     ┌─────┐
│ Alice │     │ Bob │
└───┬───┘     └──┬──┘
    │            │
    │ Think      │
    ├──┐         │
    │  │         │
    │◄─┘         │
    │            │
    │ Hello      │
    ├───────────►│
    │            │

# Explicit participant declarations with aliases
$ cat sequence.mermaid
sequenceDiagram
participant A as Alice
participant B as Bob
A->>B: Message from Alice
B-->>A: Reply to Alice
$ mermaid-ascii -f sequence.mermaid
┌───────┐     ┌─────┐
│ Alice │     │ Bob │
└───┬───┘     └──┬──┘
    │            │
    │ Message from Alice
    ├───────────►│
    │            │
    │ Reply to Alice
    │◄┈┈┈┈┈┈┈┈┈┈┈┤
    │            │

# ASCII mode for sequence diagrams
$ cat sequence.mermaid | mermaid-ascii --ascii
+-------+     +-----+
| Alice |     | Bob |
+---+---+     +--+--+
    |            |
    | Hello Bob! |
    +----------->|
    |            |
    | Hi Alice!  |
    |<...........+
    |            |

```

```bash
$ mermaid-ascii --help
Generate ASCII diagrams from mermaid code.

Usage:
  mermaid-ascii [flags]
  mermaid-ascii [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  web         HTTP server for rendering mermaid diagrams.

Flags:
  -a, --ascii               Don't use extended character set
  -p, --borderPadding int   Padding between text and border (default 1)
      --box-style string    Default Unicode box style for nodes without explicit shape (square, rounded, double, heavy) (default "square")
  -c, --coords              Show coordinates
      --edge-style string   Default Unicode edge style for standard arrows (light, heavy, dashed) (default "light")
  -f, --file string         Mermaid file to parse (use '-' for stdin)
  -h, --help                help for mermaid-ascii
  -x, --paddingX int        Horizontal space between nodes (default 5)
  -y, --paddingY int        Vertical space between nodes (default 5)
  -v, --verbose             Verbose output

Use "mermaid-ascii [command] --help" for more information about a command.

# And some ridiculous example
$ mermaid-ascii -f complex.mermaid
┌───┐     ┌───┐     ┌───┐     ┌───┐     ┌───┐     ┌───┐
│   │     │   │     │   │     │   │     │   │     │   │
│ A ├────►│ B ├──┬─►│ E ├──┬─►│ M ├──┬─►│ U ├──┬─►│ W │
│   │     │   │  │  │   │  │  │   │  │  │   │  │  │   │
└─┬─┘     └─┬─┘  │  └─┬─┘  │  └─┬─┘  │  └─┬─┘  │  └─┬─┘
  │         │    │    │    │    │    │    ▲    │    │  
  │         │    │    │    │    │    │    │    │    │  
  │    ┌────┘    ├────┘    │    ├────┘    ├────┼────┘  
  │    │         │         │    │         │    │       
  │    │         │         │    ▼         ▼    │       
  │    │  ┌───┐  │  ┌───┐  │  ┌─┴─┐     ┌───┐  │  ┌───┐
  │    │  │   │  │  │   │  │  │   │     │   │  │  │   │
  ├────┼─►│ C ├──┼─►│ F │  ├─►│ Q ├────►│ Y │◄─┼─►│ V │
  │    │  │   │  │  │   │  │  │   │     │   │  │  │   │
  │    │  └─┬─┘  │  └─┬─┘  │  └───┘     └───┘  │  └─┬─┘
  │    │    │    │    │    │                   │    ▲  
  │    │    │    │    │    │                   │    │  
  │    └────┼────┤    └────┤                   │    │  
  │         │    │         │                   │    │  
  │         ▼    │         │                   │    │  
  │       ┌─┴─┐  │  ┌───┐  │  ┌───┐     ┌───┐  │    │  
  │       │   │  │  │   │  │  │   │     │   │  │    │  
  └──────►│ D │  ├─►│ G │  ├─►│ L ├──┬─►│ T ├──┼────┤  
          │   │  │  │   │  │  │   │  │  │   │  │    │  
          └─┬─┘  │  └─┬─┘  │  └─┬─┘  │  └─┬─┘  │    │  
            │    │    │    │    │    │    ▲    │    │  
            │    │    │    │    │    │    │    │    │  
            │    │    ├────┼────┘    │    ├────┤    │  
            │    │    │    │         │    │    │    │  
            │    │    ▼    │         │    ▼    │    │  
            │    │  ┌─┴─┐  │  ┌───┐  │  ┌───┐  │    │  
            │    │  │   │  │  │   │  │  │   │  │    │  
            │    ├─►│ H │  ├─►│ J │  ├─►│ X │◄─┼────┤  
            │    │  │   │  │  │   │  │  │   │  │    │  
            │    │  └─┬─┘  │  └─┬─┘  │  └───┘  │    │  
            │    │    │    │    │    │         │    │  
            │    │    │    │    │    │         │    │  
            │    └────┼────┤    └────┤    ┌────┤    │  
            │         │    │         │    │    │    │  
            │         ▼    │         │    │    │    │  
            │       ┌─┴─┐  │  ┌───┐  │  ┌─┴─┐  │    │  
            │       │   │  │  │   │  │  │   │  │    │  
            └──────►│ I │  ├─►│ K │  ├─►│ R ├──┼────┘  
                    │   │  │  │   │  │  │   │  │       
                    └───┘  │  └─┬─┘  │  └───┘  │       
                           │    │    │         │       
                           │    │    │         │       
                           │    ├────┼────┬────┤       
                           │    │    │    │    │       
                           │    ▼    │    │    │       
                           │  ┌─┴─┐  │  ┌─┴─┐  │       
                           │  │   │  │  │   │  │       
                           ├─►│ N │  ├─►│ O │  │       
                           │  │   │  │  │   │  │       
                           │  └───┘  │  └─┬─┘  │       
                           │         │    │    │       
                           │         │    │    │       
                           ├────┬────┤    ├────┘       
                           │    │    │    │            
                           │    ▼    │    ▼            
                           │  ┌─┴─┐  │  ┌─┴─┐          
                           │  │   │  │  │   │          
                           └─►│ P │  └─►│ S │          
                              │   │     │   │          
                              └───┘     └───┘          

```

Colored output is also supported in terminals that support ANSI colors and in the web renderer.

Use `classDef` with `:::className` or `class nodeId className` to color node text, box borders, and node fill:

```bash
graph LR
classDef idea stroke:#2f80ed,color:#2f80ed,fill:#eaf3ff
classDef done stroke:#27ae60,color:#145a32,fill:#eafff2

A[Idea]:::idea --> B[Done]:::done
```

Supported node color keys:

| Style key | Effect |
| --- | --- |
| `stroke` | Node box border color |
| `color` | Node label text color |
| `fill` | Node interior background color |

Use `linkStyle` to color arrows and edge labels by edge index:

```bash
graph LR
A[Idea] -->|research| B[Research]
B ==> C[Plan]

linkStyle 0 stroke:#f2994a,color:#9b5100
linkStyle 1 stroke:#eb5757,color:#eb5757
```

Supported edge color keys:

| Style key | Effect |
| --- | --- |
| `stroke` | Arrow or connector color |
| `color` | Edge label text color. If omitted, the label uses `stroke`. |

This results in the following graph:

![](docs/colored_graph.png)

## How it works

We parse a mermaid file into basic components in order to render a grid. The grid is used for mapping purposes, which is eventually converted to a drawing.
The grid looks a bit like this:

```
There are three grid-points per node, and one in-between nodes.
These coords don't have to be the same size, as long as they
can be used for pathing purposes where we convert them to drawing
coordinates.
This allows us to navigate edges between nodes, like the arrow in this
drawing taking the path [(2,1), (3,1), (3,5), (4,5)].
    0      1      2  3  4      5      6
    |      |      |  |  |      |      |
    v      v      v  v  v      v      v
                                       
0-> +-------------+     +-------------+
    |             |     |             |
1-> |  Some text  |---  |  Some text  |
    |             |  |  |             |
2-> +-------------+  |  +-------------+
                     |                 
3->                  |                 
                     |                 
4-> +-------------+  |  +-------------+
    |             |  |  |             |
5-> |  Some text  |  -->|  Some text  |
    |             |     |             |
6-> +-------------+     +-------------+
```

You can show these coords in your graph by enabling the `--coords` flag:

```bash
$ mermaid-ascii -f ./test.mermaid --coords
   01  23    45  67  89       0
   0123456789012345678901234567
0 0+---+     +---+   +--------+
  1|   |     |   |   |        |
1 2| A |-123>| B |-->|   D    |
  3|   |     |   |   |        |
2 4+---+     +---+   +--------+
  5  |         |          |
3 6  |         2          |
  7  |         v       123456
4 8  |       +---+        |
  9  |       |   |        |
510  ------->| C |<--------
 11          |   |
612          +---+
```

Note that with `--coords` enabled, the grid-coords shown show the starting location of the coord, not the center of the coord. This is why `(1,0)` is next to `(0,0)` instead of in the center of the `A` node.

## Supported Diagram Types

### Graphs / Flowcharts ✅
- [x] Graph directions (`graph LR`, `graph RL`, `graph TD`, `graph TB`, and `graph BT`)
- [x] Labelled edges (like `A -->|label| B`)
- [x] Multiple arrows on one line (like `A --> B --> C`)
- [x] `A & B` syntax
- [x] `classDef` and `class` for colored output
- [x] Prevent arrows overlapping nodes
- [x] `subgraph` support
- [x] Shapes other than rectangles
- [ ] Diagonal arrows

### Sequence Diagrams ✅
- [x] Basic message syntax (`A->>B: message`)
- [x] Solid arrows (`->>`) and dotted arrows (`-->>`)
- [x] Self-messages (`A->>A: think`)
- [x] Participant declarations (`participant Alice`)
- [x] Participant aliases (`participant A as Alice`)
- [x] Actor declarations and aliases (`actor A as Alice`)
- [x] Notes (`Note left of Alice`, `Note right of Alice`, `Note over Alice,Bob`)
- [x] Activation directives (`activate Alice`, `deactivate Alice`)
- [x] Activation shorthand (`A->>+B`, `B-->>-A`)
- [x] Fragment frames (`loop`, `alt`, `else`, `opt`, `par`)
- [x] Unicode support (emojis, CJK characters, etc.)
- [x] Both ASCII and Unicode rendering modes
- [ ] Activation boxes
- [x] Notes (`Note left of Alice: text`)
- [x] Loops, alt, opt, and par blocks

### ER Diagrams ✅
- [x] `erDiagram` detection and rendering
- [x] Entity declarations and aliases (`CUSTOMER[Customer]`)
- [x] Entity attribute blocks with key markers (`PK`, `FK`, `UK`)
- [x] Identifying (`--`) and non-identifying (`..`) relationships
- [x] Cardinality markers such as `||`, `|o`, `o{`, `}o`, and `}|`
- [x] ASCII and Unicode rendering modes

### Class Diagrams ✅
- [x] `classDiagram` detection and rendering
- [x] Class declarations and labels (`class BankAccount["Bank Account"]`)
- [x] Class member blocks and colon member declarations
- [x] Attribute and operation compartments
- [x] Core relationships (`<|--`, `*--`, `o--`, `-->`, `--`, `..>`, `..|>`, `..`)
- [x] Relationship labels and quoted cardinality strings
- [x] ASCII and Unicode rendering modes

## TODOs

The baseline components for Mermaid work, but there are a lot of things that are not supported yet. Here's a list of things that are not yet supported:

### Syntax support

- [x] Labelled edges (like `A -->|label| B`)
- [x] Graph directions like `graph LR`, `graph RL`, `graph TD`, `graph TB`, and `graph BT`
- [x] `classDef` and `class`
- [x] `A & B`
- [x] Multiple arrows on one line (like `A --> B --> C`)
- [x] `subgraph`
- [x] Shapes other than rectangles
- [x] Whitespacing and comments

### Rendering

- [x] Prevent arrows overlapping nodes
- [ ] Diagonal arrows
- [ ] Place nodes in a more compact way
- [ ] Prevent rendering more than X characters wide (like default 80 for terminal width)

### Sequence Diagram Improvements

- [x] Activation boxes (activate/deactivate)
- [x] Notes (`Note left of Alice: text`)
- [x] Loops, alt, opt, and par blocks

### General

- [ ] Support for more diagram types (class diagrams, state diagrams, etc.)
