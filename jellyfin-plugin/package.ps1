# Simple packaging script
$version = "1.0.1"
$binPath = "Jellyfin.Plugin.TorrServer\bin\Release\net8.0"
$packageDir = "packages\TorrServer"
$zipPath = "packages\Jellyfin.Plugin.TorrServer_$version.zip"

# Create package directory
if (Test-Path "packages") { Remove-Item "packages" -Recurse -Force }
New-Item -ItemType Directory -Path $packageDir -Force | Out-Null

# Copy DLL
Copy-Item "$binPath\Jellyfin.Plugin.TorrServer.dll" $packageDir

# Copy HTML files if they exist
if (Test-Path "$binPath\Configuration") {
    Copy-Item "$binPath\Configuration" $packageDir -Recurse
}

# Create meta.json
@"
{
  "guid": "a8b12c3d-4e5f-6a7b-8c9d-0e1f2a3b4c5d",
  "name": "TorrServer",
  "description": "Manage torrents directly from Jellyfin with TorrServer integration",
  "overview": "TorrServer integration for Jellyfin",
  "owner": "Ernous",
  "category": "General",
  "version": "$version",
  "targetAbi": "10.9.0.0"
}
"@ | Out-File "$packageDir\meta.json" -Encoding UTF8

# Create ZIP
Compress-Archive -Path "$packageDir\*" -DestinationPath $zipPath -Force

# Calculate checksum
$hash = (Get-FileHash -Path $zipPath -Algorithm MD5).Hash.ToLower()

Write-Host "Package created: $zipPath"
Write-Host "MD5 Checksum: $hash"
Write-Host ""
Write-Host "Update manifest.json with:"
Write-Host "  version: $version"
Write-Host "  checksum: $hash"
