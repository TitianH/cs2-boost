@echo off
echo Building cs2-boost...
echo.

go mod tidy
if %errorlevel% neq 0 (
    echo Failed to tidy modules
    exit /b 1
)

go build -ldflags="-s -w" -o cs2-boost.exe .
if %errorlevel% neq 0 (
    echo Build failed!
    exit /b 1
)

echo.
echo Build successful: cs2-boost.exe
echo.
echo Usage:
echo   .\cs2-boost.exe install
echo   .\cs2-boost.exe uninstall
echo   .\cs2-boost.exe status
