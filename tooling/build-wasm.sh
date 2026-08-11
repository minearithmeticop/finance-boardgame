#!/usr/bin/env bash
# Build Game Engine (Go) เป็น WebAssembly สำหรับรันใน browser
# Output: apps/web/public/wasm/{engine.wasm, wasm_exec.js}
#
# รันจาก root: pnpm build:wasm   (หรือ bash tooling/build-wasm.sh)
set -euo pipefail

# Resolve repo root = parent dir ของ tooling/
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
OUT="$ROOT/apps/web/public/wasm"
mkdir -p "$OUT"

echo "→ Compiling engine.wasm (GOOS=js GOARCH=wasm)…"
GOOS=js GOARCH=wasm go build -trimpath \
  -o "$OUT/engine.wasm" \
  "$ROOT/packages/engine/cmd/wasm"

echo "→ Locating wasm_exec.js (Go runtime support)…"
GOROOT="$(go env GOROOT)"
SRC=""
for candidate in "$GOROOT/lib/wasm/wasm_exec.js" "$GOROOT/misc/wasm/wasm_exec.js"; do
  if [ -f "$candidate" ]; then
    SRC="$candidate"
    break
  fi
done
if [ -z "$SRC" ]; then
  echo "✗ wasm_exec.js ไม่พบใน GOROOT ($GOROOT)" >&2
  exit 1
fi
cp "$SRC" "$OUT/wasm_exec.js"

echo "✓ Done"
echo "  engine.wasm : $(du -h "$OUT/engine.wasm" | cut -f1)"
echo "  wasm_exec.js: $(du -h "$OUT/wasm_exec.js" | cut -f1)"
