#!/bin/bash
# Bulk CUID Migration Replacements
# Run from project root: bash CUID_BULK_REPLACEMENTS.sh

set -e

echo "Starting CUID Migration - Bulk Replacements"
echo "WARNING: This will modify many files. Ensure git is clean first!"
echo ""

# STEP 1: Update Service Interfaces and Implementations
echo "[1/5] Updating service interfaces (uint -> string)..."

# Replace in all service files
find ./internal/service -name "*.go" -type f -exec sed -i 's/, userID uint/, userID string/g' {} \;
find ./internal/service -name "*.go" -type f -exec sed -i 's/, classID uint/, classID string/g' {} \;
find ./internal/service -name "*.go" -type f -exec sed -i 's/, projectID uint/, projectID string/g' {} \;
find ./internal/service -name "*.go" -type f -exec sed -i 's/, categoryID uint/, categoryID string/g' {} \;
find ./internal/service -name "*.go" -type f -exec sed -i 's/, adminID uint/, adminID string/g' {} \;
find ./internal/service -name "*.go" -type f -exec sed -i 's/, certID uint/, certID string/g' {} \;
find ./internal/service -name "*.go" -type f -exec sed -i 's/, showcaseID uint/, showcaseID string/g' {} \;
find ./internal/service -name "*.go" -type f -exec sed -i 's/(ctx context.Context, userID uint)/(ctx context.Context, userID string)/g' {} \;
find ./internal/service -name "*.go" -type f -exec sed -i 's/uint(userID)/userID/g' {} \;
find ./internal/service -name "*.go" -type f -exec sed -i 's/uint(classID)/classID/g' {} \;

echo "✓ Service interfaces updated"

# STEP 2: Update Repository Interfaces  
echo "[2/5] Updating repository interfaces (uint -> string)..."

find ./internal/repository -name "*.go" -type f -exec sed -i 's/FindByID(.*uint/FindByID(ctx context.Context, id string/g' {} \;
find ./internal/repository -name "*.go" -type f -exec sed -i 's/, userID uint/, userID string/g' {} \;
find ./internal/repository -name "*.go" -type f -exec sed -i 's/, classID uint/, classID string/g' {} \;
find ./internal/repository -name "*.go" -type f -exec sed -i 's/, id uint/, id string/g' {} \;

echo "✓ Repository interfaces updated"

# STEP 3: Update Remaining Handler Files
echo "[3/5] Updating remaining handler files..."

find ./internal/handler -name "*.go" -type f -exec sed -i 's/c.Locals("userID").(uint)/c.Locals("userID").(string)/g' {} \;
find ./internal/handler -name "*.go" -type f -exec sed -i 's/c.ParamsInt(/c.Params(/g' {} \;
find ./internal/handler -name "*.go" -type f -exec sed -i 's/c.QueryInt(/c.Query(/g' {} \;
find ./internal/handler -name "*.go" -type f -exec sed -i 's/uint(/"/g' {} \;

echo "✓ Handler files updated"

# STEP 4: Update Seed Command
echo "[4/5] Updating seed command..."

sed -i 's/uint/string/g' ./cmd/seed/main.go
# Fix: revert specific lines that shouldn't be strings (like loop indices, counts)
sed -i 's/string(1)/1/g' ./cmd/seed/main.go

echo "✓ Seed command updated"

# STEP 5: Verify Build
echo "[5/5] Verifying build..."

if go build ./cmd/api > /dev/null 2>&1; then
    echo "✓ Build successful!"
else
    echo "⚠ Build has errors. Check output:"
    go build ./cmd/api
fi

echo ""
echo "✓ Migration complete! Review changes and commit to git"
