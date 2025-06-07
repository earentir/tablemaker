# ASCII Table Generator

A Go application that reads JSON configuration files and generates beautiful ASCII tables or PNG images with customizable fonts and styling.

## Features

- 📊 Generate ASCII tables from JSON configuration
- 🎨 Export to PNG with custom fonts and sizes
- 🖥️ **Automatic system font detection** - no need to download fonts!
- 🔧 Configurable table styles (extensible architecture)
- 📝 Support for markdown-style bold text in cells
- 🎯 **Cell alignment options** - left, center, right alignment per column
- 📏 **Smart column sizing** - automatically calculates optimal widths
- 🖼️ Transparent PNG backgrounds
- 💾 Output to file or stdout
- 🎯 Generic table rendering system for easy style additions

## Installation

### Prerequisites

- Go 1.19 or higher
- System fonts (automatically detected) OR custom TTF fonts

### Build from source

```bash
# Clone or create the project directory
mkdir ascii-table-generator && cd ascii-table-generator

# Initialize and install dependencies
go mod tidy

# Build the application
make build
# or
go build -o ascii-table-generator .
```

### Font Setup

The application automatically detects and uses system fonts! No manual font setup required.

**Automatic System Font Detection:**
- **Windows**: Searches Windows/Fonts directory
- **macOS**: Searches system and user font directories
- **Linux**: Searches /usr/share/fonts and user directories

**Supported Font Names** (searched automatically):
- Sans-serif: "DejaVu Sans", "Arial", "Helvetica", "Liberation Sans", "FreeSans", "Noto Sans"
- Bold: "DejaVu Sans Bold", "Arial Bold", "Helvetica Bold", etc.
- Monospace: "DejaVu Sans Mono", "Courier New", "Liberation Mono", "FreeMono", "Noto Sans Mono", "Consolas"

**Custom Fonts** (optional):
```bash
mkdir fonts
# Place your custom TTF files in the fonts/ directory
```

## Usage

### Basic Usage

```bash
# Generate ASCII table to stdout
./ascii-table-generator -input example.json

# Save ASCII table to file
./ascii-table-generator -input example.json -out output.txt

# Generate PNG image
./ascii-table-generator -input example.json -png -out table.png
```

### Command Line Options

- `-input <file>`: Input JSON configuration file (required)
- `-out <file>`: Output file path (optional, defaults to stdout for text)
- `-png`: Generate PNG output instead of text

## JSON Configuration Format

```json
{
  "type": "single-line-full",
  "name": "Table Name",
  "headers": ["Column 1", "Column 2", "Column 3"],
  "alignment": ["left", "center", "right"],
  "rows": [
    ["**Bold Text**", "Normal text", "More content"],
    ["Row 2 Col 1", "Row 2 Col 2", "Row 2 Col 3"]
  ],
  "png": {
    "title_font": {
      "path": "./fonts/DejaVuSans-Bold.ttf",
      "size": 14
    },
    "content_font": {
      "path": "./fonts/DejaVuSans.ttf",
      "size": 12
    },
    "ascii_font": {
      "path": "./fonts/DejaVuSansMono.ttf",
      "size": 12
    }
  }
}
```

### Configuration Fields

- **type**: Table style ("single-line-full" or "double-line-full")
- **name**: Table name/title (for reference)
- **headers**: Array of column headers
- **rows**: Array of row data (each row is an array of cell values)
- **alignment**: Array of alignment options for each column (optional)
  - Options: "left", "center"/"centre", "right"
  - If not specified, defaults to "left" for all columns
  - Can specify fewer alignments than columns (remaining default to "left")
- **png**: PNG generation configuration (optional for PNG output)
  - **title_font**: Font configuration for bold text
  - **content_font**: Font configuration for regular text
  - **ascii_font**: Font configuration for table borders

### Font Configuration

You can specify fonts in multiple ways:

**1. System Font Names** (recommended - automatically detected):
```json
"title_font": {
  "path": "Arial Bold",
  "size": 14
}
```

**2. Relative/Absolute File Paths**:
```json
"title_font": {
  "path": "./fonts/MyFont-Bold.ttf",
  "size": 14
}
```

**3. Auto-Detection** (leave path empty):
```json
"png": {
  "title_font": {"path": "", "size": 14},
  "content_font": {"path": "", "size": 12},
  "ascii_font": {"path": "", "size": 12}
}
```

### Text Formatting

- Use `**text**` for bold formatting in cells
- Bold text will be rendered with the title font in PNG output
- Regular text uses the content font
- ASCII table borders use the ASCII font

### Column Alignment

Control how text is aligned within each column:

- **"left"**: Text aligned to the left (default)
- **"center"** or **"centre"**: Text centered in the column
- **"right"**: Text aligned to the right

Example:
```json
{
  "headers": ["Name", "Status", "Score", "Notes"],
  "alignment": ["left", "center", "right", "left"],
  "rows": [
    ["John Doe", "Active", "95.5", "Excellent performance"],
    ["Jane Smith", "Pending", "87.2", "Good work"]
  ]
}
```

This creates:
```
┌───────────┬─────────┬───────┬─────────────────────┐
│ Name      │ Status  │ Score │ Notes               │
├───────────┼─────────┼───────┼─────────────────────┤
│ John Doe  │ Active  │  95.5 │ Excellent performance │
│ Jane Smith│ Pending │  87.2 │ Good work           │
└───────────┴─────────┴───────┴─────────────────────┘
```

## Table Styles

Currently supported table styles:

### single-line-full
Creates tables with single-line borders using Unicode box-drawing characters:

```
┌───────────────┬─────────┬────────────┬────────────────────────────────┐
│ Region        │ Sectors │   Typical  │ Purpose                        │
├───────────────┼─────────┼────────────┼────────────────────────────────┤
│ **Reserved**  │ 1–n     │ 1 (FAT12)  │ Boot sector + optional extras  │
├───────────────┼─────────┼────────────┼────────────────────────────────┤
│ **FAT #1**    │ varies  │            │ Allocation map                 │
└───────────────┴─────────┴────────────┴────────────────────────────────┘
```

### double-line-full
Creates tables with double-line borders using Unicode box-drawing characters:

```
╔═══════════════╦═════════╦════════════╦════════════════════════════════╗
║ Database      ║ Status  ║ Connections║ Response Time                  ║
╠═══════════════╬═════════╬════════════╬════════════════════════════════╣
║ **Primary DB**║ Online  ║ 45/100     ║ 12ms                           ║
╠═══════════════╬═════════╬════════════╬════════════════════════════════╣
║ **Replica DB**║ Online  ║ 23/50      ║ 8ms                            ║
╚═══════════════╩═════════╩════════════╩════════════════════════════════╝
```

## Examples

### Quick Start

```bash
make example          # Generate text output
make example-png      # Generate PNG output
```

### Custom Configuration

Create your own JSON file with alignment options:

```json
{
  "type": "double-line-full",
  "name": "Financial Report",
  "headers": ["Item", "Q1", "Q2", "Q3", "Q4", "Total"],
  "alignment": ["left", "right", "right", "right", "right", "right"],
  "rows": [
    ["**Revenue**", "$125,000", "$138,500", "$142,000", "$155,000", "$560,500"],
    ["**Expenses**", "$98,000", "$102,000", "$105,000", "$108,000", "$413,000"],
    ["**Profit**", "$27,000", "$36,500", "$37,000", "$47,000", "$147,500"],
    ["**Margin**", "21.6%", "26.4%", "26.1%", "30.3%", "26.3%"]
  ],
  "png": {
    "title_font": {
      "path": "Arial Bold",
      "size": 16
    },
    "content_font": {
      "path": "Arial",
      "size": 14
    },
    "ascii_font": {
      "path": "Consolas",
      "size": 14
    }
  }
}
```

## Development

### Project Structure

```
ascii-table-generator/
├── main.go                          # Main application entry point
├── tables/                          # ASCII table generation package
│   └── ascii_tables.go              # Table rendering logic and styles
├── output/                          # Output generation package
│   ├── png.go                       # PNG generation and text parsing
│   └── fonts.go                     # System font detection
├── go.mod                           # Go module definition
├── example.json                     # Example configuration
├── example-double.json              # Double-line style example
├── example-alignment.json           # Alignment demonstration
├── Makefile                         # Build automation
└── README.md                        # This file
```

### Package Architecture

The application is organized into clean, modular packages:

**Main Application (`main.go`)**
- CLI argument parsing
- JSON configuration loading
- Coordination between packages

**Tables Package (`tables/`)**
- ASCII table generation and rendering
- Table style definitions and management
- Column width calculation and text alignment
- Extensible architecture for new table styles

**Output Package (`output/`)**
- PNG image generation with font rendering
- System font detection and resolution
- Text parsing and categorization
- Cross-platform font path management

### Adding New Table Styles

The modular architecture makes adding new table styles simple:

1. **Define your style** in the `tables/ascii_tables.go` file:

```go
// Add to the tableStyles map
"my-custom-style": {
    TopLeft:     "╭",
    TopRight:    "╮",
    BottomLeft:  "╰",
    BottomRight: "╯",
    Horizontal:  "─",
    Vertical:    "│",
    TopJoin:     "┬",
    BottomJoin:  "┴",
    LeftJoin:    "├",
    RightJoin:   "┤",
    Cross:       "┼",
},
```

2. **Or register dynamically** in your code:
```go
import "ascii-table-generator/tables"

customStyle := tables.TableStyle{
    TopLeft: "╭", TopRight: "╮",
    // ... define all characters
}
tables.RegisterTableStyle("my-style", customStyle)
```

3. **Use it in your JSON**:
```json
{
  "type": "my-custom-style",
  ...
}
```

The generic renderer automatically handles the layout using your character set!

### Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes following the package architecture:
   - **Tables package**: For new table styles or rendering features
   - **Output package**: For new output formats or font improvements
   - **Main application**: For CLI or configuration changes
4. Add tests if applicable
5. Run `make fmt` and `make lint`
6. Submit a pull request

### Package Development

**Adding Table Features:**
- Modify `tables/ascii_tables.go` for rendering logic
- Add new styles to the `tableStyles` map
- Update alignment or formatting functions

**Adding Output Features:**
- Modify `output/png.go` for PNG generation
- Update `output/fonts.go` for font handling
- Add new output formats as separate files

**Testing Individual Packages:**
```bash
# Test tables package
go test ./tables/

# Test output package
go test ./output/

# Test everything
go test ./...
```

## Troubleshooting

### Common Issues

**Font loading errors**:
- The app will automatically find system fonts
- If custom fonts fail, check the file path and ensure TTF format
- Use system font names like "Arial" instead of file paths when possible

**PNG generation fails**:
- System fonts are detected automatically
- For custom fonts, ensure TTF files exist at specified paths
- Leave font paths empty to use automatic detection

**Unicode display issues**:
- Ensure your terminal supports Unicode box-drawing characters
- Try a different terminal or font if characters appear as boxes

## License

MIT License - feel free to use this project for any purpose.

## Future Enhancements

- Additional table styles (rounded corners, thick borders, etc.)
- Color support for PNG output
- CSV input support
- More font format support (OTF, WOFF)
- Table themes and presets
- Multi-line cell content support
- Custom border thickness
- Column width controls
- Cell padding customization
