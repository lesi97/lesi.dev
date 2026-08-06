Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

$repoRoot = Resolve-Path -LiteralPath (Join-Path $PSScriptRoot "..")
$apiDir = Join-Path $repoRoot "apps\api"
$tmpDir = Join-Path $repoRoot "tmp"
$goTmpDir = Join-Path $tmpDir "go-build"

New-Item -ItemType Directory -Force -Path $tmpDir | Out-Null
New-Item -ItemType Directory -Force -Path $goTmpDir | Out-Null

$exeName = "api-dev"
if ($IsWindows -or $env:OS -eq "Windows_NT") {
    $exeName = "api-dev.exe"
}

$exePath = Join-Path $tmpDir $exeName
$env:GOTMPDIR = $goTmpDir

Push-Location $apiDir
try {
    go build -o $exePath ./cmd/main.go
    & $exePath
} finally {
    Pop-Location
}
