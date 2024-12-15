# Makefile for the Valmet QueryForge project

setup:
	ollama serve &
	ollama pull llama3.2:1b
	ollama pull llama3.2:1b
	ollama pull all-minilm:33m
	ollama pull all-minilm:125m
