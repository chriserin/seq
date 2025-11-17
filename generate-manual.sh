#!/bin/bash

# Script to generate PDF manual from README and documentation files

set -e

# Function to replace documentation links with internal anchors
replace_doc_links() {
    sed -E 's/\[([^]]+)\]\(docs\/beat-creation\.md\)/[\1](#beat-creation)/g;
        s/\[([^]]+)\]\(docs\/note-alteration\.md\)/[\1](#note-alteration)/g;
        s/\[([^]]+)\]\(docs\/key-cycles\.md\)/[\1](#key-cycles)/g;
        s/\[([^]]+)\]\(docs\/arrangement\.md\)/[\1](#arrangement)/g;
        s/\[([^]]+)\]\(docs\/overlay-key\.md\)/[\1](#overlay-key)/g;
        s/\[([^]]+)\]\(docs\/actions\.md\)/[\1](#actions-reference)/g;
        s/\[([^]]+)\]\(docs\/key-mappings\.md\)/[\1](#key-mappings)/g;
        s/\[([^]]+)\]\(docs\/core-concepts\.md\)/[\1](#core-concepts)/g'
}

# Check if pandoc is installed
if ! command -v pandoc &>/dev/null; then
    echo "Error: pandoc is not installed."
    echo "Please install pandoc:"
    echo "  macOS:   brew install pandoc"
    echo "  Linux:   sudo apt-get install pandoc"
    echo "  Windows: choco install pandoc"
    exit 1
fi

# Output file
OUTPUT="sq-manual.pdf"

# Create temporary file to hold combined markdown
TEMP_MD=$(mktemp)

# Title page content
cat >"$TEMP_MD" <<EOF
---
title: "sq - MIDI Sequencer Manual"
author: "sq Project"
date: "Generated on $(date +%Y-%m-%d)"
toc: true
toc-depth: 2
---

\newpage

EOF

# Add README (filter out Documentation and Development sections)
echo "Adding README.md..."
echo -e "\n# Introduction\n" >>"$TEMP_MD"
awk '
BEGIN { skip = 0 }
/^## Documentation/ { skip = 1; next }
/^## Development/ { skip = 1; next }
/^## License/ { skip = 1; next }
/^## Contributing/ { skip = 1; next }
/^## Support/ { skip = 1; next }
/^## Install from Source/ { skip = 1; next }
/^## / && skip == 1 { skip = 0 }
skip == 0 && NR > 1 { print }
' README.md |
    replace_doc_links |
    awk '{print} /^```/ {print ""}' \
        >>"$TEMP_MD"
echo -e "\n\n" >>"$TEMP_MD"

# Add documentation files in logical order
DOCS=(
    "docs/core-concepts.md"
    "docs/key-mappings.md"
    "docs/beat-creation.md"
    "docs/note-alteration.md"
    "docs/actions.md"
    "docs/arrangement.md"
    "docs/key-cycles.md"
    "docs/overlay-key.md"
)

for doc in "${DOCS[@]}"; do
    if [ -f "$doc" ]; then
        echo "Adding $doc..."

        # Special handling for key-mappings.md to control table column widths
        if [ "$doc" = "docs/key-mappings.md" ]; then
            # Limit table separator dashes to control column widths
            # This gives pandoc hints about relative column widths
            awk '
            /^\| *-+ *\| *-+ *\| *-+ *\|/ {
                # Column 1: 27 chars, Column 2: 15 chars, Column 3: 85 chars (total ~120)
                print "| --------------------------- | --------------- | ------------------------------------------------------------------------------------- |"
                next
            }
            { print }
            ' "$doc" |
                replace_doc_links |
                awk '{print} /^```/ {print ""}' \
                    >>"$TEMP_MD"
        elif [ "$doc" = "docs/actions.md" ]; then
            # Add custom anchor to Actions to avoid conflict with README subsection
            cat "$doc" |
                sed -E 's/^# Actions$/# Actions {#actions-reference}/' |
                replace_doc_links |
                awk '{print} /^```/ {print ""}' \
                    >>"$TEMP_MD"
        elif [ "$doc" = "docs/arrangement.md" ]; then
            # Special handling for arrangement.md to preserve code blocks with unicode
            cat "$doc" |
                replace_doc_links |
                awk '{print} /^```/ {print ""}' \
                    >>"$TEMP_MD"
        else
            cat "$doc" |
                replace_doc_links |
                awk '{print} /^```/ {print ""}' \
                    >>"$TEMP_MD"
        fi

        echo -e "\n\n" >>"$TEMP_MD"
    else
        echo "Warning: $doc not found, skipping..."
    fi
done

# Generate PDF with pandoc
echo "Generating PDF..."
pandoc "$TEMP_MD" \
    -o "$OUTPUT" \
    --pdf-engine=xelatex \
    --toc \
    --toc-depth=2 \
    --number-sections \
    -V geometry:"margin=0.75in" \
    -V fontsize=10pt \
    -V documentclass=article \
    -V colorlinks=true \
    -V linkcolor=blue \
    -V urlcolor=blue \
    -V tables=true \
    -V mainfont="Helvetica Neue" \
    -V monofont="FreeMono" \
    -V header-includes='\usepackage{pdflscape}' \
    -V header-includes='\usepackage{fancyvrb}' \
    -V header-includes='\usepackage{upquote}' \
    -V header-includes='\usepackage{booktabs}' \
    -V header-includes='\renewcommand{\arraystretch}{1.5}' \
    -V header-includes='\DefineVerbatimEnvironment{Highlighting}{Verbatim}{commandchars=\\\{\},fontsize=\small}' \
    --syntax-highlighting=tango

# Clean up
rm "$TEMP_MD"

echo "✓ Manual generated: $OUTPUT"
