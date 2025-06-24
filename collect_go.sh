#!/usr/bin/env bash
# собирает все .go в project_overview.go
OUTPUT="project_overview.go"
echo "// Объединённый обзор кода проекта" > "$OUTPUT"
find . -type f -name '*.go' \
  ! -path './vendor/*' \
  ! -path './tests/*' \
  -print0 | sort -z | while IFS= read -r -d $'\0' file; do
    echo -e "\n// ===== File: ${file#./} =====\n" >> "$OUTPUT"
    sed '/^package /d' "$file" >> "$OUTPUT"
done
echo "Готово: см. $OUTPUT"
