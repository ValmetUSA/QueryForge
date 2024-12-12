<div align="center">
  
# <img src="assets/img/valmet_logo_nobg.png" alt="Valmet" width="500" height="200"/>
# QueryForge
Produced by [VII @ Valmet](https://www.valmet.com/automation/industrial-internet/)

### A simple, yet robust, private local AI RAG chat application for Ollama - designed for conversations involving sensitive data.
___
</div>

## 📖 Table of Contents:
* [Project Background](#-project-background)
* [How to Run](#-how-to-run)
* [Building From Source](#-building-from-source-advanced-users-only)

---

## 🤓 Project Background
QueryForge is a lightweight, secure, and efficient local AI-driven RAG (Retrieval Augmented Generation) chat application designed specifically for conversations involving sensitive data. Built entirely in the Go programming language, QueryForge offers a robust, easy-to-use interface for querying a local AI model (on Ollama) without relying on cloud-based services, ensuring privacy and security for users dealing with confidential information.

This project was developed by [Valmet of North America](https://www.valmet.com). Valmet is an international automation and services company that specializes in providing solutions to industries such as pulp, paper, and energy. With a strong focus on advancing sustainable practices, Valmet offers a range of technologies and services aimed at improving operational efficiency and environmental performance. Their expertise spans from process automation to machine vision systems, and they work with clients globally to implement cutting-edge solutions that drive innovation.

Valmet’s commitment to integrating advanced technologies, such as artificial intelligence, with traditional industries allows them to support customers in transforming their operations. The development of QueryForge is part of Valmet's ongoing efforts to explore new frontiers in automation and AI applications, particularly in sectors that require high levels of data privacy and security. By focusing on local, private AI-driven solutions, Valmet aims to provide businesses with powerful tools for managing sensitive data while maintaining complete control over the information flow.

Valmet of North America is based in Atlanta, Georgia, and is part of the broader Valmet global network, which itself spans over 30 countries worldwide.
___

### Features:
- Private and Local: All operations occur entirely on the local machine, ensuring sensitive data remains private.
- Customizable AI Model: Select from different base conversational and embedding models to fine-tune the AI's responses.
- Easy Folder Selection: Choose a directory for running the RAG search, streamlining the process of retrieving relevant documents for AI-based responses.
- Progress Feedback: A progress bar indicates the status of queries, giving users visibility into processing times.
- Clipboard Integration: Copy and paste functionality is available directly from the toolbar, enhancing usability.
- Simple Interface: Designed with an intuitive Fyne-based GUI for seamless interaction.

### Underlying Code
The core application is written in Go and utilizes the Fyne GUI framework for creating a cross-platform desktop application. Here's an overview of some key components:
- Fyne Framework: Used for creating the UI elements - it also ensures cross-platform compatibility.
- Effortless AI Model Integration: The application queries an AI model in Ollama, which processes the input text and retrieves relevant responses. This keeps things simple for the end user, since Ollama handles model management - no AI interfacing is done manually by the user.
- Folder Selection: Users can select a folder to run the RAG search, our chunking process for document handling gets you an asnwer quickly.
- Settings and Model Selection: The app allows users to select base conversational models and embedding models to customize the AI's behavior.
___

## 🦙 How To Run
TBD
____

## 👷‍♂️ Building From Source (Advanced Users Only!)
> [!IMPORTANT]
> Before you start - make sure you have the [GoLang compiler downloaded and installed for your operating system](https://go.dev/doc/install)!

1. Clone this git repo, either run the command below or download the repo as a zip:
   ```
   git clone https://github.com/ValmetUSA/QueryForge.git
   ```
2. Use the `cd` command to naviagate to the directory containing QueryForge (likely your home directory), and download all build dependencies by running the command below:
   ```
   go mod download
   ```
3. Use the `cd` command again to naviagate to the `./src` folder.

4. Run the following command to build an executable for your operating system:
  ```
  go run .
  ```
5. Run the following command to run the program directly from source:
  ```
  go build .
  ```
___

To Do (Desciption, % of entire project):
- [x] Get UI and Framework done - 30%
- [x] Add RAG - 30%
- [x] Finish front end (as of recently, mostly working with a few bugs) - 20%
- [ ] Finish documentation on GitHub - 5%
- [ ] Finish file picker with RAG, along with adjustments - 10%
- [ ] Write script to setup Ollama - 5%
