@echo off
:: Ensure script runs with administrative privileges
net session >nul 2>&1
if %errorlevel% neq 0 (
    echo Please run this script as an administrator.
    pause
    exit /b
)

:: Install Ollama using winget
echo Installing Ollama...
winget install --id=Ollama.Ollama -e

if %errorlevel% neq 0 (
    echo Failed to install Ollama. Exiting script.
    pause
    exit /b
)

:: Verify Ollama is installed
echo Verifying Ollama installation...
where ollama >nul 2>&1
if %errorlevel% neq 0 (
    echo Ollama command-line tool not found. Please check the installation.
    pause
    exit /b
)

:: Pull the specified models
echo Pulling models...
ollama pull qwen2.5:0.5b
ollama pull llama3.2:1b
ollama pull llama3.2:3b
ollama pull phi3:3.8b

if %errorlevel% neq 0 (
    echo Failed to pull one or more models. Check the error messages above.
    pause
    exit /b
)

echo Ollama and models installed successfully!
pause
exit /b
