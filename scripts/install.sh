#!/usr/bin/env bash
set -euo pipefail

# Minimum versions
REQUIRED_GO=1.24
REQUIRED_NODE=18

# Detect operating system
if command -v cmd.exe >/dev/null 2>&1; then
    OS="windows"
else
    case "$(uname -s)" in
        Linux*) OS="linux" ;;
        Darwin*) OS="darwin" ;;
        *) echo "Unsupported platform: $(uname -s)" >&2; exit 1 ;;
    esac
fi
ARCH=$(uname -m)

# Extract a zip archive to the given destination using whatever tool is
# available. This supports environments as old as Windows Server 2008 where
# `unzip` may not be installed by default.
extract_zip() {
    local archive=$1 dest=$2
    if command -v unzip >/dev/null 2>&1; then
        unzip -q "$archive" -d "$dest"
    elif command -v 7z >/dev/null 2>&1; then
        7z x -y "$archive" -o"$dest" >/dev/null
    elif command -v powershell.exe >/dev/null 2>&1; then
        powershell.exe -NoProfile -Command "Expand-Archive -Path '$archive' -DestinationPath '$dest' -Force" >/dev/null
    else
        echo "No unzip utility found. Please install unzip or 7zip." >&2
        return 1
    fi
}

version_ge() {
    # returns 0 if $1 >= $2
    [ "$(printf '%s\n' "$2" "$1" | sort -V | head -n1)" = "$2" ]
}

install_go() {
    if command -v go >/dev/null 2>&1; then
        local current
        current=$(go version | awk '{print $3}' | sed 's/go//')
        if version_ge "$current" "$REQUIRED_GO"; then
            echo "Go $current already installed"
            return
        fi
    fi

    local pkg url tmp archive
    pkg="go${REQUIRED_GO}.1"
    case "$OS" in
        linux)
            pkg+=".linux-amd64.tar.gz"
            archive="go.tgz"
            ;;
        darwin)
            pkg+=".darwin-amd64.tar.gz"
            archive="go.tgz"
            ;;
        windows)
            pkg+=".windows-amd64.zip"
            archive="go.zip"
            ;;
    esac
    url="https://go.dev/dl/${pkg}"
    tmp=$(mktemp -d)
    trap 'rm -rf "$tmp"' RETURN
    echo "Installing Go from $url"
    curl -fsSL "$url" -o "$tmp/$archive"
    case "$OS" in
        windows)
            extract_zip "$tmp/$archive" /usr/local
            ;;
        *)
            sudo tar -C /usr/local -xzf "$tmp/$archive"
            ;;
    esac
    export PATH="/usr/local/go/bin:$PATH"
}

install_node() {
    if command -v node >/dev/null 2>&1; then
        local current
        current=$(node -v | sed 's/v//')
        if version_ge "$current" "$REQUIRED_NODE"; then
            echo "Node $current already installed"
            return
        fi
    fi

    echo "Installing Node.js"
    case "$OS" in
        linux)
            curl -fsSL https://nodejs.org/dist/v18.20.8/node-v18.20.8-linux-x64.tar.gz -o node.tgz
            sudo tar -C /usr/local --strip-components=1 -xzf node.tgz
            rm node.tgz
            ;;
        darwin)
            if command -v brew >/dev/null 2>&1; then
                brew install node@18
                brew link --force --overwrite node@18
            else
                curl -fsSL https://nodejs.org/dist/v18.20.8/node-v18.20.8-darwin-x64.tar.gz -o node.tgz
                sudo tar -C /usr/local --strip-components=1 -xzf node.tgz
                rm node.tgz
            fi
            ;;
        windows)
            if command -v choco >/dev/null 2>&1; then
                choco install -y nodejs-lts
            else
                echo "Please install Node.js LTS manually." >&2
                exit 1
            fi
            ;;
    esac
}

install_wails() {
    echo "Installing Wails CLI"
    go install github.com/wailsapp/wails/v2/cmd/wails@latest
}

install_go
install_node
install_wails
