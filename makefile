# Makefile for the Valmet QueryForge project

default:
	@echo "Installing the Valmet QueryForge project, and all its dependencies."
	@echo "Please wait... This may take a while."
	go install fyne.io/fyne/v2/cmd/fyne@latest
	go mod download
	ollama serve &
	ollama pull qwen2.5:0.5b
	ollama pull llama3.2:1b
	ollama pull llama3.2:3b
	ollama pull phi3:3.8b

macos:
	@echo "Building the Valmet QueryForge project for MacOS."
	@echo "Please wait... This may take a while."
	fyne package -os darwin --exe ./build -icon ./build_assets/valmet_logo_noname.png --src ./src --name "Valmet QueryForge"

windows:
	@echo "Building the Valmet QueryForge project for Windows."
	@echo "Please wait... This may take a while."
	fyne package -os windows --exe ./build -icon ./build_assets/valmet_logo_noname.png --src ./src --name "Valmet QueryForge"

linux:
	@echo "Building the Valmet QueryForge project for Linux."
	@echo "Please wait... This may take a while."
	fyne package -os linux --exe ./build -icon ./build_assets/valmet_logo_noname.png --src ./src --name "Valmet QueryForge"