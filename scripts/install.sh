#!/usr/bin/env bash
set -euo pipefail

# Minimum versions
REQUIRED_GO=1.24
REQUIRED_NODE=18

# Detect operating system
case "$(uname -s)" in
    Linux*) OS="linux" ;;
    Darwin*) OS="darwin" ;;
    CYGWIN*|MINGW*|MSYS*) OS="windows" ;;
    *) echo "Unsupported platform: $(uname -s)" >&2; exit 1 ;;
esac
ARCH=$(uname -m)

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

    local pkg url tmp
    pkg="go${REQUIRED_GO}.1"
    case "$OS" in
        linux)   pkg+=".linux-amd64.tar.gz" ;;
        darwin)  pkg+=".darwin-amd64.tar.gz" ;;
        windows) pkg+=".windows-amd64.zip" ;;
    esac
    url="https://go.dev/dl/${pkg}"
    tmp=$(mktemp -d)
    trap 'rm -rf "$tmp"' RETURN
    echo "Installing Go from $url"
    curl -fsSL "$url" -o "$tmp/go.tgz"
    case "$OS" in
        windows)
            unzip -q "$tmp/go.tgz" -d /usr/local
            ;;
        *)
            sudo tar -C /usr/local -xzf "$tmp/go.tgz"
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
            curl -fsSL https://nodejs.org/dist/latest-v18.x/node-v18.20.3-linux-x64.tar.gz -o node.tgz
            sudo tar -C /usr/local --strip-components=1 -xzf node.tgz
            rm node.tgz
            ;;
        darwin)
            if command -v brew >/dev/null 2>&1; then
                brew install node@18
                brew link --force --overwrite node@18
            else
                curl -fsSL https://nodejs.org/dist/latest-v18.x/node-v18.20.3-darwin-x64.tar.gz -o node.tgz
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
