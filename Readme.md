# RailGuard Pro

![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)
![Fyne Version](https://img.shields.io/badge/Fyne-v2.4-blue?style=for-the-badge&logo=fyne)
![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)
![Platform](https://img.shields.io/badge/Platform-Android%20%7C%20Desktop-orange?style=for-the-badge)

> **RailGuard Pro** is a cross-platform safety system engineered for railway operators to manage train compositions, calculate precise braking parameters, and ensure compliance with dangerous goods regulations. Built with performance and accuracy in mind.

---

## 📸 Interface Visuals
<table>
  <tr>
    <td align="center" width="100%">
      <h3>📱 Mobile Interface</h3>
      <img src="assets/screenshots/image_2026-02-12_22-29-25.png" alt="Mobile UI" width="300">
      <br>
      <em>Android Optimized View</em>
    </td>
  </tr>
</table>
---

## ✨ Core Features

### 🚂 Smart UIC Decoding
Automatically identifies wagon types and validates axle counts instantly from the input UIC number. Eliminates manual entry errors by cross-referencing industry standards.

### ⚠️ Dangerous Goods Compatibility Matrix
A robust, built-in safety validation engine that checks for dangerous goods compatibility in real-time. It enforces separation rules and alerts operators to potential conflicts within the train composition.

### 🧮 Physics-Based Brake Calculations
Advanced braking physics engine that computes:
- **Max Safe Speed**: Based on current track slope and train weight.
- **Braking Percentage**: Real-time calculation of braking efficiency.
- **Brake Weight**: Aggregated brake weight calculations for the entire composition.

### 💾 SQLite Persistence
Reliable local data storage using SQLite. RailGuard Pro maintains a comprehensive history of train setups, allowing operators to save, load, and modify compositions without data loss.

### 📄 PDF License Generation
One-click export of official Brake Licenses. Generates industry-standard PDF documents ready for printing or digital transmission.

---

## 🏗 Technical Architecture

RailGuard Pro is built using a **Modular Clean Architecture** in Go, ensuring high maintainability and testability.

- **Language**: Golang 1.21+
- **GUI Framework**: [Fyne v2](https://fyne.io/) (Material Design) - providing a native experience across platforms.
- **Database**: SQLite (via `modernc.org/sqlite` or `mattn/go-sqlite3` w/ CGO) for robust, embedded SQL support.
- **PDF Engine**: `gofpdf` for programmatic document generation.

---

## 🚀 Installation & Build

### Prerequisites
- **Go**: Version 1.21 or higher.
- **C Compiler**: Required for Fyne and SQLite (e.g., GCC or MinGW).
- **Fyne Cross**: (Optional) For cross-compiling to Android.

### 🖥 Desktop Run
To run the application locally on your machine:

```bash
# Clone the repository
git clone https://github.com/raminfathi/RailGuard.git
cd RailGuard

# Install dependencies
go mod tidy

# Run the application
go run ./cmd/app
```

### 📱 Android Build
We use `fyne-cross` to build optimized APKs for Android.

1.  **Install fyne-cross**:
    ```bash
    go install github.com/fyne-io/fyne-cross/v2/cmd/fyne-cross@latest
    ```

2.  **Build APK**:
    ```bash
    # Build for ARM64 architecture
    fyne-cross android -arch arm64 -app-id com.raminfathi.railguard -icon Icon.png ./cmd/app
    ```
    The output APK will be located in the `fyne-cross/bin/android-arm64` directory.

---

## 📦 Release Guide

Follow these steps to publish a new version:

1.  **Update Version**: Bump the version in `cmd/app/main.go` and `FyneApp.toml`.
2.  **Tag release**:
    ```bash
    git tag -a v1.0.0 -m "Release v1.0.0"
    git push origin v1.0.0
    ```
3.  **Generate Artifacts**: Run the `fyne-cross` build command (see above).
4.  **Upload**: Go to GitHub Releases, create a new release pointing to the tag, and upload the generated `.apk` and desktop binaries.

---

## 🔮 Future Roadmap

- [ ] **Network Sync**: Cloud synchronization for sharing train compositions across team devices.
- [ ] **Multi-Language Support**: i18n implementation for global railway standards.
- [ ] **Digital Signatures**: Cryptographic signing of PDF brake licenses.

---

<p align="center">
  <sub>Developed by Ramin Fathi | RailGuard Pro © 2024</sub>
</p>