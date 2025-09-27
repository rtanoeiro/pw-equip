#!/bin/bash

# Build script for PW Equipment Changer
# Handles different platforms and build modes

echo "Building PW Equipment Changer..."

# Detect OS
OS=$(uname -s)

case $OS in
    "Darwin")
        echo "Building for macOS..."
        # On macOS, build normally
        go build -o pw-equip-changer
        echo "✅ macOS build complete: pw-equip-gui"
        ;;
    "Linux")
        echo "Building for Linux..."
        # On Linux, build normally
        go build -o pw-equip-changer
        echo "✅ Linux build complete: pw-equip-gui"
        ;;
    "MINGW"*|"MSYS"*|"CYGWIN"*)
        echo "Building for Windows..."
        # On Windows, use windowsgui flag to prevent console window
        go build -ldflags="-H windowsgui" -o pw-equip-changer.exe
        echo "✅ Windows build complete: pw-equip-gui.exe"
        ;;
    *)
        echo "Unknown OS: $OS"
        echo "Building with default settings..."
        go build -o pw-equip-changer
        ;;
esac

echo "Build completed!"
