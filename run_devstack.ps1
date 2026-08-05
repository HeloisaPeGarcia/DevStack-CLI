param (
    [string]$Command = "bootstrap",
    [string]$Stack = "Go Backend + React Frontend",
    [string]$ProjectName = "my-awesome-app",
    [switch]$DryRun
)

$GoPath = Get-Command go -ErrorAction SilentlyContinue

if (-not $GoPath) {
    Write-Host "ℹ Procurando a instalação do Go no sistema..." -ForegroundColor Yellow
    $CommonPaths = @(
        "C:\Program Files\Go\bin\go.exe",
        "C:\Go\bin\go.exe",
        "$env:LocalAppData\Programs\Go\bin\go.exe"
    )
    foreach ($path in $CommonPaths) {
        if (Test-Path $path) {
            $GoExe = $path
            break
        }
    }
} else {
    $GoExe = "go"
}

if ($GoExe) {
    Write-Host "✔ Executando DevStack via Go ($GoExe)..." -ForegroundColor Green
    $cmdArgs = @("run", "./cmd/devstack", $Command, "--stack", "$Stack", "--project-name", "$ProjectName")
    if ($DryRun) { $cmdArgs += "--dry-run" }
    
    & $GoExe $cmdArgs
} else {
    Write-Host "✖ Go SDK não encontrado no PATH. Instale o Go via Winget:" -ForegroundColor Red
    Write-Host "  winget install --id GoLang.Go" -ForegroundColor Cyan
}
