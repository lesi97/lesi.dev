param(
  [string]$Tag = "latest",
  [string]$Platform = "linux/amd64"
)

$EnvFile = Join-Path $PSScriptRoot "..\.env.local"

if (Test-Path $EnvFile) {
  Get-Content $EnvFile | ForEach-Object {
    if ($_ -match "^\s*([^#=]+)\s*=\s*(.+)\s*$") {
      $name  = $matches[1]
      $value = $matches[2]
      Set-Item -Path "env:$name" -Value $value
    }
  }
}

if (-not $env:GHCR_USER) { throw "Set GHCR_USER" }
if (-not $env:GHCR_TOKEN) { throw "Set GHCR_TOKEN" }
if (-not $env:GHCR_IMAGE) { throw "Set GHCR_IMAGE like ghcr.io/owner/repo" }

$env:GHCR_TOKEN | docker login ghcr.io -u $env:GHCR_USER --password-stdin | Out-Null

docker buildx create --use | Out-Null
docker buildx inspect --bootstrap | Out-Null

docker buildx build --push --platform $Platform -t "$env:GHCR_IMAGE`:$Tag" .
