# GoTunnel Windows Installer
# The open-source ngrok killer

param(
    [switch]$Silent,
    [switch]$Uninstall
)

$ErrorActionPreference = "Stop"

# Configuration
$AppName = "GoTunnel"
$AppVersion = "1.0.0"
$InstallDir = "$env:PROGRAMFILES\GoTunnel"
$StartMenuDir = "$env:APPDATA\Microsoft\Windows\Start Menu\Programs\GoTunnel"
$DesktopShortcut = "$env:USERPROFILE\Desktop\GoTunnel Dashboard.lnk"

# Colors for output
$Green = "Green"
$Yellow = "Yellow"
$Red = "Red"
$Cyan = "Cyan"

function Write-ColorOutput {
    param([string]$Message, [string]$Color = "White")
    Write-Host $Message -ForegroundColor $Color
}

function Show-Banner {
    Write-ColorOutput @"
╔══════════════════════════════════════════════════════════════╗
║                    🚀 GoTunnel Installer                     ║
║              The open-source ngrok killer                    ║
║                                                              ║
║  Self-hosted, secure, and blazingly fast tunneling          ║
╚══════════════════════════════════════════════════════════════╝
"@ $Cyan
}

function Test-Administrator {
    $currentUser = [Security.Principal.WindowsIdentity]::GetCurrent()
    $principal = New-Object Security.Principal.WindowsPrincipal($currentUser)
    return $principal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
}

function Install-GoTunnel {
    Write-ColorOutput "Installing GoTunnel..." $Green
    
    # Create installation directory
    if (!(Test-Path $InstallDir)) {
        New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
        Write-ColorOutput "Created installation directory: $InstallDir" $Green
    }
    
    # Download GoTunnel binaries
    Write-ColorOutput "Downloading GoTunnel binaries..." $Yellow
    
    $releasesUrl = "https://api.github.com/repos/ogrok/gotunnel/releases/latest"
    try {
        $release = Invoke-RestMethod -Uri $releasesUrl
        $assets = $release.assets
        
        # Find Windows assets
        $serverAsset = $assets | Where-Object { $_.name -like "*gotunnel-server*" -and $_.name -like "*windows*" }
        $clientAsset = $assets | Where-Object { $_.name -like "*gotunnel-client*" -and $_.name -like "*windows*" }
        $dashboardAsset = $assets | Where-Object { $_.name -like "*dashboard*" -and $_.name -like "*windows*" }
        
        if ($serverAsset -and $clientAsset) {
            # Download server
            Write-ColorOutput "Downloading server binary..." $Yellow
            Invoke-WebRequest -Uri $serverAsset.browser_download_url -OutFile "$InstallDir\gotunnel-server.exe"
            
            # Download client
            Write-ColorOutput "Downloading client binary..." $Yellow
            Invoke-WebRequest -Uri $clientAsset.browser_download_url -OutFile "$InstallDir\gotunnel-client.exe"
            
            # Download dashboard if available
            if ($dashboardAsset) {
                Write-ColorOutput "Downloading dashboard..." $Yellow
                Invoke-WebRequest -Uri $dashboardAsset.browser_download_url -OutFile "$InstallDir\gotunnel-dashboard.exe"
            }
            
            Write-ColorOutput "Binaries downloaded successfully!" $Green
        } else {
            Write-ColorOutput "Could not find Windows binaries. Building from source..." $Yellow
            Build-FromSource
        }
    } catch {
        Write-ColorOutput "Failed to download binaries. Building from source..." $Yellow
        Build-FromSource
    }
    
    # Create start menu directory
    if (!(Test-Path $StartMenuDir)) {
        New-Item -ItemType Directory -Path $StartMenuDir -Force | Out-Null
    }
    
    # Create shortcuts
    Create-Shortcuts
    
    # Add to PATH
    Add-ToPath
    
    # Create configuration
    Create-Configuration
    
    Write-ColorOutput "GoTunnel installed successfully!" $Green
    Write-ColorOutput "You can now run 'gotunnel-server' and 'gotunnel-client' from anywhere." $Cyan
}

function Build-FromSource {
    Write-ColorOutput "Building GoTunnel from source..." $Yellow
    
    # Check if Go is installed
    try {
        $goVersion = go version
        Write-ColorOutput "Found Go: $goVersion" $Green
    } catch {
        Write-ColorOutput "Go not found. Installing Go..." $Yellow
        Install-Go
    }
    
    # Clone repository
    $tempDir = "$env:TEMP\gotunnel-build"
    if (Test-Path $tempDir) {
        Remove-Item $tempDir -Recurse -Force
    }
    
    git clone https://github.com/ogrok/gotunnel.git $tempDir
    Set-Location $tempDir
    
    # Build binaries
    Write-ColorOutput "Building server..." $Yellow
    go build -o "$InstallDir\gotunnel-server.exe" ./cmd/server
    
    Write-ColorOutput "Building client..." $Yellow
    go build -o "$InstallDir\gotunnel-client.exe" ./cmd/client
    
    # Build dashboard
    Write-ColorOutput "Building dashboard..." $Yellow
    Set-Location dashboard
    npm install
    npm run build
    
    # Copy dashboard files
    Copy-Item -Path "dist\*" -Destination $InstallDir -Recurse -Force
    
    # Cleanup
    Set-Location $env:TEMP
    Remove-Item $tempDir -Recurse -Force
    
    Write-ColorOutput "Build completed successfully!" $Green
}

function Install-Go {
    Write-ColorOutput "Installing Go..." $Yellow
    
    # Download Go installer
    $goUrl = "https://go.dev/dl/go1.24.5.windows-amd64.msi"
    $goInstaller = "$env:TEMP\go-installer.msi"
    
    Invoke-WebRequest -Uri $goUrl -OutFile $goInstaller
    
    # Install Go
    Start-Process -FilePath "msiexec.exe" -ArgumentList "/i", $goInstaller, "/quiet", "/norestart" -Wait
    
    # Refresh environment
    $env:PATH = [System.Environment]::GetEnvironmentVariable("PATH","Machine") + ";" + [System.Environment]::GetEnvironmentVariable("PATH","User")
    
    # Cleanup
    Remove-Item $goInstaller -Force
    
    Write-ColorOutput "Go installed successfully!" $Green
}

function Create-Shortcuts {
    Write-ColorOutput "Creating shortcuts..." $Yellow
    
    # Create start menu shortcuts
    $WshShell = New-Object -comObject WScript.Shell
    
    # Server shortcut
    $serverShortcut = $WshShell.CreateShortcut("$StartMenuDir\GoTunnel Server.lnk")
    $serverShortcut.TargetPath = "$InstallDir\gotunnel-server.exe"
    $serverShortcut.WorkingDirectory = $InstallDir
    $serverShortcut.Description = "Start GoTunnel Server"
    $serverShortcut.Save()
    
    # Client shortcut
    $clientShortcut = $WshShell.CreateShortcut("$StartMenuDir\GoTunnel Client.lnk")
    $clientShortcut.TargetPath = "$InstallDir\gotunnel-client.exe"
    $clientShortcut.WorkingDirectory = $InstallDir
    $clientShortcut.Description = "Start GoTunnel Client"
    $clientShortcut.Save()
    
    # Dashboard shortcut
    if (Test-Path "$InstallDir\gotunnel-dashboard.exe") {
        $dashboardShortcut = $WshShell.CreateShortcut("$StartMenuDir\GoTunnel Dashboard.lnk")
        $dashboardShortcut.TargetPath = "$InstallDir\gotunnel-dashboard.exe"
        $dashboardShortcut.WorkingDirectory = $InstallDir
        $dashboardShortcut.Description = "Open GoTunnel Dashboard"
        $dashboardShortcut.Save()
        
        # Desktop shortcut
        $desktopShortcut = $WshShell.CreateShortcut($DesktopShortcut)
        $desktopShortcut.TargetPath = "$InstallDir\gotunnel-dashboard.exe"
        $desktopShortcut.WorkingDirectory = $InstallDir
        $desktopShortcut.Description = "GoTunnel Dashboard - The open-source ngrok killer"
        $desktopShortcut.Save()
    }
    
    Write-ColorOutput "Shortcuts created successfully!" $Green
}

function Add-ToPath {
    Write-ColorOutput "Adding GoTunnel to PATH..." $Yellow
    
    $currentPath = [Environment]::GetEnvironmentVariable("PATH", "Machine")
    if ($currentPath -notlike "*$InstallDir*") {
        $newPath = $currentPath + ";" + $InstallDir
        [Environment]::SetEnvironmentVariable("PATH", $newPath, "Machine")
        Write-ColorOutput "Added to system PATH" $Green
    } else {
        Write-ColorOutput "Already in PATH" $Green
    }
}

function Create-Configuration {
    Write-ColorOutput "Creating default configuration..." $Yellow
    
    $configDir = "$env:USERPROFILE\.gotunnel"
    if (!(Test-Path $configDir)) {
        New-Item -ItemType Directory -Path $configDir -Force | Out-Null
    }
    
    $configContent = @"
# GoTunnel Configuration
# The open-source ngrok killer

server:
  port: 8080
  tls: false
  allowed_tokens:
    - "your-secret-token"

client:
  server: "localhost:8080"
  token: "your-secret-token"
  tls: false
  log_level: "info"

dashboard:
  enabled: true
  port: 3000
"@
    
    $configPath = "$configDir\config.yaml"
    $configContent | Out-File -FilePath $configPath -Encoding UTF8
    
    Write-ColorOutput "Configuration created: $configPath" $Green
}

function Uninstall-GoTunnel {
    Write-ColorOutput "Uninstalling GoTunnel..." $Yellow
    
    # Remove from PATH
    $currentPath = [Environment]::GetEnvironmentVariable("PATH", "Machine")
    $newPath = ($currentPath -split ';' | Where-Object { $_ -ne $InstallDir }) -join ';'
    [Environment]::SetEnvironmentVariable("PATH", $newPath, "Machine")
    
    # Remove shortcuts
    if (Test-Path $StartMenuDir) {
        Remove-Item $StartMenuDir -Recurse -Force
    }
    
    if (Test-Path $DesktopShortcut) {
        Remove-Item $DesktopShortcut -Force
    }
    
    # Remove installation directory
    if (Test-Path $InstallDir) {
        Remove-Item $InstallDir -Recurse -Force
    }
    
    # Remove configuration
    $configDir = "$env:USERPROFILE\.gotunnel"
    if (Test-Path $configDir) {
        Remove-Item $configDir -Recurse -Force
    }
    
    Write-ColorOutput "GoTunnel uninstalled successfully!" $Green
}

# Main execution
Show-Banner

if (!(Test-Administrator)) {
    Write-ColorOutput "This installer requires administrator privileges." $Red
    Write-ColorOutput "Please run PowerShell as Administrator and try again." $Red
    exit 1
}

if ($Uninstall) {
    Uninstall-GoTunnel
} else {
    Install-GoTunnel
    
    Write-ColorOutput @"

🎉 Installation Complete!

Next steps:
1. Open GoTunnel Dashboard from your desktop or start menu
2. Or run from command line:
   - gotunnel-server --port 8080 --allowed-tokens "your-secret-token"
   - gotunnel-client --server localhost:8080 --subdomain myapp --local-port 3000 --token "your-secret-token"

For more information, visit: https://github.com/ogrok/gotunnel

Happy tunneling! 🚀
"@ $Green
} 