package yggdrasil

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

// CheckAndInstall проверяет наличие Yggdrasil и предлагает установку
func CheckAndInstall() (installed bool, path string, err error) {
	// Проверяем стандартные пути
	paths := []string{
		"/usr/bin/yggdrasil",
		"/usr/local/bin/yggdrasil",
		"/opt/yggdrasil/bin/yggdrasil",
	}

	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return true, p, nil
		}
	}

	// Проверяем через PATH
	if path, err := exec.LookPath("yggdrasil"); err == nil {
		return true, path, nil
	}

	// Не найдено — предлагаем установку
	return false, "", fmt.Errorf("yggdrasil not found")
}

// GetInstallCommand возвращает команду для установки в зависимости от ОС
func GetInstallCommand() string {
	switch runtime.GOOS {
	case "linux":
		// Определяем дистрибутив
		if _, err := os.Stat("/etc/debian_version"); err == nil {
			return "sudo apt update && sudo apt install -y yggdrasil"
		}
		if _, err := os.Stat("/etc/redhat-release"); err == nil {
			return "sudo dnf install -y yggdrasil"
		}
		return "curl -s https://yggdrasil-network.github.io/install.sh | sudo bash"
	case "darwin":
		return "brew install yggdrasil"
	case "windows":
		return "choco install yggdrasil"
	default:
		return "See https://yggdrasil-network.github.io/install.html"
	}
}

// PrintInstallInstructions выводит инструкции по установке
func PrintInstallInstructions() {
	fmt.Println(`
🌐 Yggdrasil not found. To enable mesh networking:

1. Install Yggdrasil:
   ` + GetInstallCommand() + `

2. Start the service:
   sudo systemctl enable --now yggdrasil

3. Verify installation:
   yggdrasilctl getself

4. Run ideal-core again.

Alternatively, run in fallback mode with:
   go run cmd/node/main.go -yggdrasil ""

⚠️  Note: Fallback mode is for local testing only.
   For production mesh networking, Yggdrasil is required.
`)
}
