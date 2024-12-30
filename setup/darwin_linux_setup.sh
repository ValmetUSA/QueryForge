#!/bin/bash

# Ensure the script is run with administrative privileges
if [ "$EUID" -ne 0 ]; then
    echo "Please run this script as root or with sudo."
    exit 1
fi

# Install Homebrew if not installed
if ! command -v brew &>/dev/null; then
    echo "Homebrew is not installed. Installing Homebrew..."
    /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
    if [ $? -ne 0 ]; then
        echo "Failed to install Homebrew. Exiting script."
        exit 1
    fi
else
    echo "Homebrew is already installed."
fi

# Install Ollama using Homebrew
echo "Installing Ollama..."
brew install ollama

if [ $? -ne 0 ]; then
    echo "Failed to install Ollama. Exiting script."
    exit 1
fi

# Verify Ollama installation
if ! command -v ollama &>/dev/null; then
    echo "Ollama command-line tool not found. Please check the installation."
    exit 1
fi

# Pull the specified models
echo "Pulling models..."
ollama pull qwen2.5:0.5b
ollama pull llama3.2:1b
ollama pull llama3.2:3b
ollama pull phi3:3.8b

if [ $? -ne 0 ]; then
    echo "Failed to pull one or more models. Check the error messages above."
    exit 1
fi

echo "Ollama and models installed successfully!"
exit 0
