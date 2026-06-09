#!/bin/bash
# CUID Migration Helper Script
# This script documents all the replacements needed to complete the CUID migration
# Run after reading CUID_MIGRATION_GUIDE.md

echo "=== CUID Migration - Pattern Examples ==="
echo ""
echo "1. HANDLER FILES - Replace uint with string"
echo "   Pattern: c.Locals(\"userID\").(uint)"
echo "   Replace: c.Locals(\"userID\").(string)"
echo ""
echo "   Pattern: c.ParamsInt(\"id\")"
echo "   Replace: c.Params(\"id\")"
echo ""
echo "   Pattern: c.QueryInt(\"user_id\")"
echo "   Replace: c.Query(\"user_id\")"
echo ""
echo "2. SERVICE FILES - Update method signatures"
echo "   Pattern: func (s *XService) Method(ctx context.Context, userID uint, itemID uint) error"
echo "   Replace: func (s *XService) Method(ctx context.Context, userID string, itemID string) error"
echo ""
echo "3. REPOSITORY FILES - Update method signatures"
echo "   Pattern: func (r *XRepository) FindByID(id uint) (*domain.X, error)"
echo "   Replace: func (r *XRepository) FindByID(id string) (*domain.X, error)"
echo ""
echo "4. INTERFACES - Update interface method signatures"
echo "   Pattern: interface IXService{ Method(ctx context.Context, id uint) error }"
echo "   Replace: interface IXService{ Method(ctx context.Context, id string) error }"
echo ""
echo "=== FILES TO UPDATE ==="
echo ""
echo "HANDLERS (8 files):"
find ./internal/handler -name "*.go" -type f | sed 's/^/  - /'
echo ""
echo "SERVICES (7 files):"
find ./internal/service -name "*.go" -type f ! -name "service.go" | sed 's/^/  - /'
echo ""
echo "REPOSITORIES (many files):"
find ./internal/repository -name "*.go" -type f ! -name "init.go" | sed 's/^/  - /'
echo ""
echo "SEEDS:"
echo "  - cmd/seed/main.go"
