#!/bin/bash
# Lint 工具自动安装脚本 - Linux/macOS 版本
# 根据语言类型自动安装所需的 Lint 工具

set -e

LANGUAGE=$1
MODE=${2:-"fast"}  # fast 或 deep
SANGFOR_PYPI="http://mirrors.sangfor.org/pypi/simple"
SANGFOR_GO_PROXY="http://mirrors.sangfor.org/nexus/repository/go-proxy-group"

# 颜色输出
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查工具是否已安装
check_tool() {
    local tool=$1
    if command -v "$tool" &> /dev/null; then
        log_info "$tool 已安装"
        return 0
    else
        log_warn "$tool 未安装"
        return 1
    fi
}

# 安装 Go Lint 工具
install_go_tools() {
    log_info "开始安装 Go Lint 工具..."

    # 配置 Go 代理
    export GO111MODULE=on
    export GOPROXY=$SANGFOR_GO_PROXY
    export GOSUMDB=off

    log_info "已配置 Go Proxy: $GOPROXY"

    # 安装 golangci-lint
    if ! check_tool golangci-lint; then
        log_info "正在安装 golangci-lint..."

        # 方法1: 使用 go install（推荐）
        if command -v go &> /dev/null; then
            go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

            # 确保 GOPATH/bin 在 PATH 中
            export PATH=$PATH:$(go env GOPATH)/bin

            if check_tool golangci-lint; then
                log_info "✅ golangci-lint 安装成功"
                return 0
            fi
        fi

        # 方法2: 使用官方安装脚本（备用）
        log_warn "使用备用安装方法..."
        curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

        if check_tool golangci-lint; then
            log_info "✅ golangci-lint 安装成功（备用方法）"
            return 0
        else
            log_error "❌ golangci-lint 安装失败"
            return 1
        fi
    fi

    return 0
}

# 安装 Python Lint 工具
install_python_tools() {
    log_info "开始安装 Python Lint 工具..."

    local tools_to_install=()

    # 检查并记录需要安装的工具
    if ! check_tool pylint; then
        tools_to_install+=("pylint")
    fi

    if ! check_tool flake8; then
        tools_to_install+=("flake8")
    fi

    if ! check_tool black; then
        tools_to_install+=("black")
    fi

    # 如果有需要安装的工具
    if [ ${#tools_to_install[@]} -gt 0 ]; then
        log_info "需要安装的工具: ${tools_to_install[*]}"
        log_info "正在使用深信服镜像源安装..."

        # 使用深信服 PyPI 镜像一次性安装所有工具
        if pip3 install -i "$SANGFOR_PYPI" --trusted-host mirrors.sangfor.org "${tools_to_install[@]}"; then
            log_info "✅ Python Lint 工具安装成功"

            # 验证安装
            local failed=0
            for tool in "${tools_to_install[@]}"; do
                if ! check_tool "$tool"; then
                    log_error "❌ $tool 验证失败"
                    failed=1
                fi
            done

            return $failed
        else
            log_error "❌ Python Lint 工具安装失败"
            return 1
        fi
    else
        log_info "✅ 所有 Python Lint 工具已安装"
        return 0
    fi
}

# 安装 Java Lint 工具
install_java_tools() {
    log_info "开始检查 Java Lint 工具..."

    # checkstyle 和 spotbugs 通常需要手动下载或通过包管理器安装
    # 这里提供安装指导

    local missing_tools=()

    if ! check_tool checkstyle; then
        missing_tools+=("checkstyle")
        log_warn "checkstyle 未安装"
        log_info "安装方法1: brew install checkstyle (macOS)"
        log_info "安装方法2: apt-get install checkstyle (Ubuntu/Debian)"
        log_info "安装方法3: 下载 jar 文件 - https://github.com/checkstyle/checkstyle/releases"
    fi

    if ! command -v spotbugs &> /dev/null && [ ! -f "/usr/local/bin/spotbugs" ]; then
        missing_tools+=("spotbugs")
        log_warn "spotbugs 未安装"
        log_info "安装方法1: brew install spotbugs (macOS)"
        log_info "安装方法2: 下载并解压 - https://github.com/spotbugs/spotbugs/releases"
    fi

    if [ ${#missing_tools[@]} -gt 0 ]; then
        log_warn "⚠️  Java Lint 工具需要手动安装: ${missing_tools[*]}"
        log_warn "工具检查将在运行时跳过未安装的工具"
        return 1
    else
        log_info "✅ Java Lint 工具已安装"
        return 0
    fi
}

# 主函数
main() {
    if [ -z "$LANGUAGE" ]; then
        log_error "用法: $0 <language> [mode]"
        log_error "language: go, python, java"
        log_error "mode: fast, deep (默认: fast)"
        exit 1
    fi

    log_info "========================================="
    log_info "Lint 工具自动安装脚本"
    log_info "语言: $LANGUAGE"
    log_info "模式: $MODE"
    log_info "========================================="
    echo

    case "$LANGUAGE" in
        go|golang)
            install_go_tools
            exit_code=$?
            ;;
        python|py)
            install_python_tools
            exit_code=$?
            ;;
        java)
            install_java_tools
            exit_code=$?
            ;;
        *)
            log_error "不支持的语言: $LANGUAGE"
            log_error "支持的语言: go, python, java"
            exit 1
            ;;
    esac

    echo
    log_info "========================================="
    if [ $exit_code -eq 0 ]; then
        log_info "✅ 安装检查完成，所有工具就绪"
    else
        log_warn "⚠️  部分工具安装失败或需要手动安装"
        log_warn "Lint 检查将跳过未安装的工具"
    fi
    log_info "========================================="

    exit $exit_code
}

main "$@"
