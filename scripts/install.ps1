$ErrorActionPreference = 'Stop'
Set-StrictMode -Version Latest

if ($PSVersionTable.PSVersion.Major -lt 5) {
    Write-Error 'PowerShell 5.1 or newer is required'
    exit 1
}

$LogPath = $env:SENTINEL_LOG_PATH
if (-not $LogPath) {
    if ($IsWindows) {
        $LogPath = Join-Path $env:SystemDrive 'Temp\sentinel.log'
    } else {
        $LogPath = '/var/log/sentinel.log'
    }
}

function Rotate-Log {
    $maxSize = 20MB
    $maxBackups = 5
    $maxAge = [TimeSpan]::FromDays(30)
    if (Test-Path $LogPath) {
        $info = Get-Item $LogPath
        if ($info.Length -ge $maxSize) {
            $ts = (Get-Date -Format 'yyyyMMddTHHmmssZ' -Date (Get-Date).ToUniversalTime())
            $backup = "$LogPath-$ts"
            Move-Item $LogPath $backup
            $cutoff = (Get-Date).ToUniversalTime().Add(-$maxAge)
            Get-ChildItem "$LogPath-*" | Where-Object { $_.LastWriteTimeUtc -lt $cutoff } | Remove-Item -ErrorAction SilentlyContinue
            $files = Get-ChildItem "$LogPath-*" | Sort-Object LastWriteTimeUtc -Descending
            if ($files.Count -gt $maxBackups) {
                $files | Select-Object -Skip $maxBackups | Remove-Item -ErrorAction SilentlyContinue
            }
        }
    }
}

function Log {
    param([string]$Message, [string]$Level = 'INFO')
    Rotate-Log
    $ts = (Get-Date -Date (Get-Date).ToUniversalTime() -Format o)
    $entry = @{timestamp=$ts; level=$Level; message=$Message} | ConvertTo-Json -Compress
    Add-Content -Path $LogPath -Value $entry
}

trap {
    Log $_.Exception.Message 'ERROR'
    exit 1
}

$REQUIRED_GO = '1.24'
$REQUIRED_NODE = '18'
$NODE_PATCH = '18.20.8'

function Version-GE {
    param([string]$Current, [string]$Required)
    try {
        [version]$c = $Current
        [version]$r = $Required
        return $c -ge $r
    } catch {
        return $false
    }
}

function Install-Go {
    if (Get-Command go -ErrorAction SilentlyContinue) {
        $ver = (& go version) -replace '.*go([0-9.]+).*', '$1'
        if (Version-GE $ver $REQUIRED_GO) {
            Log "Go $ver already installed"
            return
        }
        Log "Go $ver found but outdated"
    }
    if (Get-Command choco -ErrorAction SilentlyContinue) {
        Log 'Installing Go via Chocolatey'
        choco install -y golang --version $REQUIRED_GO | Out-Null
    } else {
        $pkg = "go${REQUIRED_GO}.1.windows-amd64.zip"
        $url = "https://go.dev/dl/$pkg"
        $tmp = New-Item -ItemType Directory -Path ([IO.Path]::GetTempPath()) -Name ([guid]::NewGuid())
        $zip = Join-Path $tmp $pkg
        Log "Downloading Go from $url"
        Invoke-WebRequest -Uri $url -OutFile $zip -UseBasicParsing
        Expand-Archive $zip -DestinationPath $tmp -Force
        Move-Item -Path (Join-Path $tmp 'go') -Destination 'C:\\Go' -Force
        Remove-Item $tmp -Recurse -Force
        $env:Path = 'C:\\Go\\bin;' + $env:Path
    }
}

function Install-Node {
    if (Get-Command node -ErrorAction SilentlyContinue) {
        $ver = (& node -v) -replace 'v', ''
        if (Version-GE $ver $REQUIRED_NODE) {
            Log "Node $ver already installed"
            return
        }
        Log "Node $ver found but outdated"
    }
    if (Get-Command choco -ErrorAction SilentlyContinue) {
        Log 'Installing Node via Chocolatey'
        choco install -y nodejs-lts | Out-Null
    } else {
        $url = "https://nodejs.org/dist/v$NODE_PATCH/node-v$NODE_PATCH-x64.msi"
        $msi = Join-Path ([IO.Path]::GetTempPath()) 'node.msi'
        Log "Downloading Node from $url"
        Invoke-WebRequest -Uri $url -OutFile $msi -UseBasicParsing
        Start-Process msiexec.exe -Wait -ArgumentList "/i `"$msi`" /qn"
        Remove-Item $msi -Force
    }
}

function Install-Wails {
    $wailsCmd = Get-Command wails -ErrorAction SilentlyContinue
    if ($wailsCmd) {
        try {
            $out = (& wails version) -join ' '
            if ($out -match 'v([0-9.]+)') {
                Log "Wails $($Matches[1]) already installed"
                return
            }
        } catch {
            # continue install if version check fails
        }
    }
    Log 'Installing Wails CLI'
    & go install github.com/wailsapp/wails/v2/cmd/wails@latest | Out-Null
}

Log 'starting install.ps1'

if (-not $env:INSTALL_NOP) {
    Install-Go
    Install-Node
    Install-Wails
}

Log 'install.ps1 completed'

