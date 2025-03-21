# CT_HW1
Computer Technology Course - Homework 1

## [Internet Download Manager in Terminal](https://github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW1)
A powerful concurrent download manager written in Go that supports multi-part downloads, bandwidth throttling, queue management, and state persistence. The user interface is implemented using [BubbleTea](https://github.com/charmbracelet/bubbletea) library.

The whole project was dissected into three parts: handling downloading, handling app state, and implementing tui, which were respectively implemented by [Parisa Sadat Mousavi](https://github.com/parisam83), [Mohammad Barekatain](https://github.com/mammedbrk), and [Mehdi Taheri](https://github.com/Mefi22).


![GoDownloadManager](https://github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW1/DownloadsTab.png)

## Features

- **Parallel Downloads**: Splits files into multiple parts for faster concurrent downloading
- **Pause & Resume**: Pause downloads and resume them later from where they left off
- **Download Queue**: Organize downloads in queues with prioritization
- **Bandwidth Control**: Limit download speeds to prevent network saturation
- **Download Progress**: Real-time monitoring of download progress and speed
- **Persistence**: Save download state and configuration across sessions
- **Error Recovery**: Automatically handles connection issues and retries
- **File Integrity**: Ensures downloaded files are complete and correctly merged


## Installation

### Prerequisites

- Go 1.16 or higher
- Git

### From Source

```bash
# Clone the repository
git clone https://github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW1.git
cd godownloadmanager

# Build the application
go build -o gdm

# Install globally (optional)
go install
```

### Using Go Modules

```bash
go get github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW1
```

---

Built with ❤️ by [Kafsh e Mardane Varzeshi Hypo Test Team](https://github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team)
