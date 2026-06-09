# ============================================================
# reset_db.ps1 - Drop & recreate database schema for dev
# Usage: .\scripts\reset_db.ps1
# ============================================================

# Load .env values
$envFile = Join-Path $PSScriptRoot ".." ".env"
$envVars = @{}
Get-Content $envFile | Where-Object { $_ -match '^\s*[^#].*=' } | ForEach-Object {
    $parts = $_ -split '=', 2
    $envVars[$parts[0].Trim()] = $parts[1].Trim()
}

$DB_HOST     = if ($envVars["DB_HOST"])     { $envVars["DB_HOST"] }     else { "localhost" }
$DB_PORT     = if ($envVars["DB_PORT"])     { $envVars["DB_PORT"] }     else { "5432" }
$DB_USER     = if ($envVars["DB_USER"])     { $envVars["DB_USER"] }     else { "postgres" }
$DB_NAME     = if ($envVars["DB_NAME"])     { $envVars["DB_NAME"] }     else { "jvalleyverse" }
$DB_PASSWORD = if ($envVars["DB_PASSWORD"]) { $envVars["DB_PASSWORD"] } else { "" }

Write-Host ""
Write-Host "=================================================" -ForegroundColor Cyan
Write-Host "  JValleyverse DB Reset (Development Only)" -ForegroundColor Cyan
Write-Host "=================================================" -ForegroundColor Cyan
Write-Host "  Host : $DB_HOST`:$DB_PORT"
Write-Host "  DB   : $DB_NAME"
Write-Host "  User : $DB_USER"
Write-Host ""
Write-Host "⚠️  This will DROP all tables and data!" -ForegroundColor Yellow
$confirm = Read-Host "Type 'yes' to continue"
if ($confirm -ne "yes") {
    Write-Host "Aborted." -ForegroundColor Red
    exit 1
}

# Set PGPASSWORD so psql doesn't prompt
$env:PGPASSWORD = $DB_PASSWORD

Write-Host ""
Write-Host "→ Dropping and recreating public schema..." -ForegroundColor Yellow

$sql = @"
DROP SCHEMA public CASCADE;
CREATE SCHEMA public;
GRANT ALL ON SCHEMA public TO $DB_USER;
GRANT ALL ON SCHEMA public TO public;
"@

$result = $sql | psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME 2>&1
if ($LASTEXITCODE -ne 0) {
    Write-Host "✗ psql failed. Make sure PostgreSQL is running and psql is in PATH." -ForegroundColor Red
    Write-Host $result
    exit 1
}

Write-Host "✓ Schema dropped and recreated." -ForegroundColor Green

# Drop the env var after use
Remove-Item Env:PGPASSWORD -ErrorAction SilentlyContinue

Write-Host ""
Write-Host "→ Running AutoMigrate + Seed via Go..." -ForegroundColor Yellow

# Go back to project root
$projectRoot = Join-Path $PSScriptRoot ".."
Set-Location $projectRoot

go run ./cmd/seed/main.go
if ($LASTEXITCODE -ne 0) {
    Write-Host "✗ Seed failed. Check errors above." -ForegroundColor Red
    exit 1
}

Write-Host ""
Write-Host "=================================================" -ForegroundColor Green
Write-Host "  ✅ Reset complete! Database is fresh." -ForegroundColor Green
Write-Host "  Run: go run ./cmd/api/main.go" -ForegroundColor Green
Write-Host "=================================================" -ForegroundColor Green
