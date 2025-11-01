#!/usr/bin/env pwsh
# Build and package Jellyfin TorrServer Plugin

param(
    [string]$Version = "1.0.1",
    [string]$Configuration = "Release"
)

$ErrorActionPreference = "Stop"

Write-Host "Building Jellyfin TorrServer Plugin v$Version..." -ForegroundColor Green

# Clean previous builds
$projectPath = Join-Path $PSScriptRoot "Jellyfin.Plugin.TorrServer"
$outputPath = Join-Path $PSScriptRoot "bin" $Configuration "net8.0"
$packagePath = Join-Path $PSScriptRoot "packages"

if (Test-Path $packagePath) {
    Remove-Item $packagePath -Recurse -Force
}
New-Item -ItemType Directory -Path $packagePath -Force | Out-Null

# Build the project
Write-Host "Building project..." -ForegroundColor Yellow
Set-Location $projectPath
dotnet build -c $Configuration

if ($LASTEXITCODE -ne 0) {
    Write-Error "Build failed!"
    exit 1
}

# Create package directory structure
$pluginDir = Join-Path $packagePath "TorrServer"
New-Item -ItemType Directory -Path $pluginDir -Force | Out-Null

# Copy DLL and dependencies
Write-Host "Copying files..." -ForegroundColor Yellow
$binPath = Join-Path $projectPath "bin" $Configuration "net8.0"
Copy-Item (Join-Path $binPath "Jellyfin.Plugin.TorrServer.dll") $pluginDir
Copy-Item (Join-Path $binPath "Jellyfin.Plugin.TorrServer.pdb") $pluginDir -ErrorAction SilentlyContinue

# Copy HTML files if they exist in output
$configPath = Join-Path $binPath "Configuration"
if (Test-Path $configPath) {
    $targetConfigPath = Join-Path $pluginDir "Configuration"
    New-Item -ItemType Directory -Path $targetConfigPath -Force | Out-Null
    Copy-Item (Join-Path $configPath "*.html") $targetConfigPath
}

# Create meta.json
$metaJson = @{
    guid = "a8b12c3d-4e5f-6a7b-8c9d-0e1f2a3b4c5d"
    name = "TorrServer"
    description = "Manage torrents directly from Jellyfin with TorrServer integration"
    overview = "TorrServer integration for Jellyfin"
    owner = "Ernous"
    category = "General"
    version = $Version
    targetAbi = "10.9.0.0"
    changelog = "Bug fixes and improvements"
    timestamp = (Get-Date -Format "yyyy-MM-ddTHH:mm:ssZ")
} | ConvertTo-Json -Depth 10

$metaJson | Out-File (Join-Path $pluginDir "meta.json") -Encoding UTF8

# Create ZIP archive
Write-Host "Creating ZIP archive..." -ForegroundColor Yellow
$zipPath = Join-Path $packagePath "Jellyfin.Plugin.TorrServer_$Version.zip"
Compress-Archive -Path "$pluginDir\*" -DestinationPath $zipPath -Force

# Calculate MD5 checksum
Write-Host "Calculating checksum..." -ForegroundColor Yellow
$md5 = Get-FileHash -Path $zipPath -Algorithm MD5
$checksum = $md5.Hash.ToLower()

Write-Host "`nBuild completed successfully!" -ForegroundColor Green
Write-Host "Package: $zipPath" -ForegroundColor Cyan
Write-Host "Checksum (MD5): $checksum" -ForegroundColor Cyan
Write-Host "`nUpdate manifest.json with this checksum!" -ForegroundColor Yellow

# Update manifest.json
$manifestPath = Join-Path $PSScriptRoot "manifest.json"
if (Test-Path $manifestPath) {
    Write-Host "`nUpdating manifest.json..." -ForegroundColor Yellow
    $manifest = Get-Content $manifestPath -Raw | ConvertFrom-Json
    $manifest[0].versions[0].version = $Version
    $manifest[0].versions[0].checksum = $checksum
    $manifest[0].versions[0].timestamp = (Get-Date -Format "yyyy-MM-ddTHH:mm:ssZ")
    $manifest | ConvertTo-Json -Depth 10 | Out-File $manifestPath -Encoding UTF8
    Write-Host "Manifest updated!" -ForegroundColor Green
}

Set-Location $PSScriptRoot
