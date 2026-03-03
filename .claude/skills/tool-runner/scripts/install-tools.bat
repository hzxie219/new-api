@echo off
REM Lint 工具自动安装脚本 - Windows 版本
REM 根据语言类型自动安装所需的 Lint 工具

setlocal EnableDelayedExpansion

set LANGUAGE=%1
set MODE=%2
if "%MODE%"=="" set MODE=fast

set SANGFOR_PYPI=http://mirrors.sangfor.org/pypi/simple
set SANGFOR_GO_PROXY=http://mirrors.sangfor.org/nexus/repository/go-proxy-group

REM 颜色代码（Windows 10+）
set "GREEN=[92m"
set "YELLOW=[93m"
set "RED=[91m"
set "NC=[0m"

REM 检查工具是否已安装
:check_tool
    where %1 >nul 2>&1
    if %errorlevel% equ 0 (
        echo %GREEN%[INFO]%NC% %1 已安装
        exit /b 0
    ) else (
        echo %YELLOW%[WARN]%NC% %1 未安装
        exit /b 1
    )

REM 安装 Go Lint 工具
:install_go_tools
    echo %GREEN%[INFO]%NC% 开始安装 Go Lint 工具...

    REM 配置 Go 代理
    set GO111MODULE=on
    set GOPROXY=%SANGFOR_GO_PROXY%
    set GOSUMDB=off

    echo %GREEN%[INFO]%NC% 已配置 Go Proxy: %GOPROXY%

    REM 检查 golangci-lint
    call :check_tool golangci-lint
    if %errorlevel% neq 0 (
        echo %GREEN%[INFO]%NC% 正在安装 golangci-lint...

        REM 检查是否有 go 命令
        where go >nul 2>&1
        if %errorlevel% equ 0 (
            REM 使用 go install 安装
            go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

            REM 添加 GOPATH\bin 到 PATH
            for /f "tokens=*" %%i in ('go env GOPATH') do set GOPATH=%%i
            set PATH=%PATH%;%GOPATH%\bin

            call :check_tool golangci-lint
            if !errorlevel! equ 0 (
                echo %GREEN%[INFO]%NC% ✅ golangci-lint 安装成功
                exit /b 0
            )
        )

        REM 备用方法：提示手动安装
        echo %YELLOW%[WARN]%NC% 自动安装失败，请手动安装 golangci-lint
        echo %YELLOW%[WARN]%NC% 安装方法1: choco install golangci-lint
        echo %YELLOW%[WARN]%NC% 安装方法2: scoop install golangci-lint
        echo %YELLOW%[WARN]%NC% 安装方法3: 下载二进制文件 - https://github.com/golangci/golangci-lint/releases
        exit /b 1
    )

    exit /b 0

REM 安装 Python Lint 工具
:install_python_tools
    echo %GREEN%[INFO]%NC% 开始安装 Python Lint 工具...

    set tools_to_install=
    set install_needed=0

    REM 检查 pylint
    call :check_tool pylint
    if %errorlevel% neq 0 (
        set tools_to_install=!tools_to_install! pylint
        set install_needed=1
    )

    REM 检查 flake8
    call :check_tool flake8
    if %errorlevel% neq 0 (
        set tools_to_install=!tools_to_install! flake8
        set install_needed=1
    )

    REM 检查 black
    call :check_tool black
    if %errorlevel% neq 0 (
        set tools_to_install=!tools_to_install! black
        set install_needed=1
    )

    REM 如果需要安装
    if !install_needed! equ 1 (
        echo %GREEN%[INFO]%NC% 需要安装的工具:!tools_to_install!
        echo %GREEN%[INFO]%NC% 正在使用深信服镜像源安装...

        REM 使用深信服 PyPI 镜像安装
        pip3 install -i %SANGFOR_PYPI% --trusted-host mirrors.sangfor.org!tools_to_install!

        if !errorlevel! equ 0 (
            echo %GREEN%[INFO]%NC% ✅ Python Lint 工具安装成功
            exit /b 0
        ) else (
            echo %RED%[ERROR]%NC% ❌ Python Lint 工具安装失败
            exit /b 1
        )
    ) else (
        echo %GREEN%[INFO]%NC% ✅ 所有 Python Lint 工具已安装
        exit /b 0
    )

REM 安装 Java Lint 工具
:install_java_tools
    echo %GREEN%[INFO]%NC% 开始检查 Java Lint 工具...

    set missing_tools=

    REM 检查 checkstyle
    call :check_tool checkstyle
    if %errorlevel% neq 0 (
        set missing_tools=!missing_tools! checkstyle
        echo %YELLOW%[WARN]%NC% checkstyle 未安装
        echo %YELLOW%[WARN]%NC% 安装方法1: choco install checkstyle
        echo %YELLOW%[WARN]%NC% 安装方法2: 下载 jar 文件 - https://github.com/checkstyle/checkstyle/releases
    )

    REM 检查 spotbugs
    where spotbugs >nul 2>&1
    if %errorlevel% neq 0 (
        if not exist "C:\Program Files\spotbugs\bin\spotbugs.bat" (
            set missing_tools=!missing_tools! spotbugs
            echo %YELLOW%[WARN]%NC% spotbugs 未安装
            echo %YELLOW%[WARN]%NC% 安装方法1: choco install spotbugs
            echo %YELLOW%[WARN]%NC% 安装方法2: 下载并解压 - https://github.com/spotbugs/spotbugs/releases
        )
    )

    if not "!missing_tools!"=="" (
        echo %YELLOW%[WARN]%NC% ⚠️  Java Lint 工具需要手动安装:!missing_tools!
        echo %YELLOW%[WARN]%NC% 工具检查将在运行时跳过未安装的工具
        exit /b 1
    ) else (
        echo %GREEN%[INFO]%NC% ✅ Java Lint 工具已安装
        exit /b 0
    )

REM 主函数
:main
    if "%LANGUAGE%"=="" (
        echo %RED%[ERROR]%NC% 用法: %~nx0 ^<language^> [mode]
        echo %RED%[ERROR]%NC% language: go, python, java
        echo %RED%[ERROR]%NC% mode: fast, deep (默认: fast)
        exit /b 1
    )

    echo %GREEN%[INFO]%NC% =========================================
    echo %GREEN%[INFO]%NC% Lint 工具自动安装脚本 - Windows
    echo %GREEN%[INFO]%NC% 语言: %LANGUAGE%
    echo %GREEN%[INFO]%NC% 模式: %MODE%
    echo %GREEN%[INFO]%NC% =========================================
    echo.

    if /i "%LANGUAGE%"=="go" goto :run_go
    if /i "%LANGUAGE%"=="golang" goto :run_go
    if /i "%LANGUAGE%"=="python" goto :run_python
    if /i "%LANGUAGE%"=="py" goto :run_python
    if /i "%LANGUAGE%"=="java" goto :run_java

    echo %RED%[ERROR]%NC% 不支持的语言: %LANGUAGE%
    echo %RED%[ERROR]%NC% 支持的语言: go, python, java
    exit /b 1

:run_go
    call :install_go_tools
    set exit_code=%errorlevel%
    goto :finish

:run_python
    call :install_python_tools
    set exit_code=%errorlevel%
    goto :finish

:run_java
    call :install_java_tools
    set exit_code=%errorlevel%
    goto :finish

:finish
    echo.
    echo %GREEN%[INFO]%NC% =========================================
    if !exit_code! equ 0 (
        echo %GREEN%[INFO]%NC% ✅ 安装检查完成，所有工具就绪
    ) else (
        echo %YELLOW%[WARN]%NC% ⚠️  部分工具安装失败或需要手动安装
        echo %YELLOW%[WARN]%NC% Lint 检查将跳过未安装的工具
    )
    echo %GREEN%[INFO]%NC% =========================================

    exit /b !exit_code!

call :main
