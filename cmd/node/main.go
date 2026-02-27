package main

import (
	"context"
	"crypto/ed25519"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"ideal-core/pkg/crypto"
	"ideal-core/pkg/rl"
	"ideal-core/pkg/yggdrasil"
)

var (
	dataDir   = flag.String("data", "~/.ideal-core", "Directory for keys and data")
	bootstrap = flag.String("bootstrap", "", "Comma-separated Yggdrasil peers to bootstrap with")
	genKey    = flag.Bool("genkey", false, "Generate new keypair and exit")
	yggPath   = flag.String("yggdrasil", "/usr/bin/yggdrasil", "Path to yggdrasil binary")
)

func main() {
	flag.Parse()

	// Resolve data directory
	dir := os.ExpandEnv(*dataDir)
	if err := os.MkdirAll(dir, 0700); err != nil {
		log.Fatalf("Failed to create data dir: %v", err)
	}

	// Key management
	keyPath := filepath.Join(dir, "private.key")
	pubKeyPath := filepath.Join(dir, "public.key")
	var kp *crypto.KeyPair

	if *genKey {
		var err error
		kp, err = crypto.GenerateKeyPair()
		if err != nil {
			log.Fatalf("Key generation failed: %v", err)
		}
		if err := crypto.SavePrivateKey(kp.PrivateKey, keyPath); err != nil {
			log.Fatalf("Failed to save private key: %v", err)
		}
		if err := os.WriteFile(pubKeyPath, kp.PublicKey, 0644); err != nil {
			log.Fatalf("Failed to save public key: %v", err)
		}
		fmt.Printf("✅ New keypair generated:\n")
		fmt.Printf("   Public ID: %s\n", kp.ToHex())
		fmt.Printf("   Yggdrasil IP: %s\n", crypto.DeriveYggdrasilIP(kp.PublicKey))
		fmt.Printf("   ⚠️  WARNING: Keep your private key secret! Never share it.\n")
		fmt.Printf("   Private key saved to: %s\n", keyPath)
		fmt.Printf("   Public key saved to: %s\n", pubKeyPath)
		return
	}

	// Load or create key
	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		fmt.Println("🔑 No key found, generating new one...")
		var err error
		kp, err = crypto.GenerateKeyPair()
		if err != nil {
			log.Fatalf("Key generation failed: %v", err)
		}
		if err := crypto.SavePrivateKey(kp.PrivateKey, keyPath); err != nil {
			log.Fatalf("Failed to save private key: %v", err)
		}
		if err := os.WriteFile(pubKeyPath, kp.PublicKey, 0644); err != nil {
			log.Fatalf("Failed to save public key: %v", err)
		}
	} else {
		priv, err := crypto.LoadPrivateKey(keyPath)
		if err != nil {
			log.Fatalf("Failed to load private key: %v", err)
		}
		pub, err := os.ReadFile(pubKeyPath)
		if err != nil {
			log.Fatalf("Failed to load public key: %v", err)
		}
		kp = &crypto.KeyPair{PrivateKey: priv, PublicKey: ed25519.PublicKey(pub)}
	}

	fmt.Printf("🗝️  Node ID: %s\n", kp.ToHex()[:16]+"...")
	fmt.Printf("🌐 Yggdrasil IP: %s\n", crypto.DeriveYggdrasilIP(kp.PublicKey))

	// Initialize RL agent
	agent := rl.NewAgent(0.1, 0.9, 0.1) // lr, discount, exploration
	_ = agent // Used for future extensions

	// Check Yggdrasil availability
	yggAvailable := checkYggdrasil(*yggPath)
	if !yggAvailable {
		log.Printf("⚠️  Yggdrasil not found at %s, using fallback transport", *yggPath)
	}

	// Connect to Yggdrasil mesh
	ygg, err := yggdrasil.NewClient(kp.ToHex(), *yggPath, yggAvailable)
	if err != nil {
		log.Fatalf("Yggdrasil init failed: %v", err)
	}

	if *bootstrap != "" {
		peers := strings.Split(*bootstrap, ",")
		if err := ygg.Bootstrap(peers); err != nil {
			log.Printf("⚠️  Bootstrap warning: %v", err)
		}
	}

	// Context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Signal handler
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\n🛑 Shutting down...")
		cancel()
	}()

	// Start message listener
	go func() {
		if err := ygg.Receive(ctx, func(msg []byte) error {
			// Decrypt and process message
			fmt.Printf("📥 Received %d bytes (encrypted)\n", len(msg))
			return nil
		}); err != nil && err != context.Canceled {
			log.Printf("Receive error: %v", err)
		}
	}()

	fmt.Println("✅ Node running. Press Ctrl+C to stop.")
	fmt.Println("🔐 Security reminder: Your private key is stored in:", keyPath)
	fmt.Println("   Never share it. Never upload it. Never transmit it.")

	// Demo: send a test intention after 3 seconds
	time.Sleep(3 * time.Second)
	testPayload := []byte("encrypted:intention:MuLaDhArA:stability")
	targetIP := crypto.DeriveYggdrasilIP(kp.PublicKey)
	if err := ygg.Send(targetIP, testPayload); err != nil {
		log.Printf("Send error: %v", err)
	}

	// Keep alive
	<-ctx.Done()
	fmt.Println("👋 Node stopped.")
}

// checkYggdrasil проверяет доступность сервиса Yggdrasil
func checkYggdrasil(path string) bool {
	// 1. Проверка бинарника
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	// 2. Проверка сокета (обычно /var/run/yggdrasil.sock или аналог)
	socketPaths := []string{
		"/var/run/yggdrasil.sock",
		"/run/yggdrasil.sock",
		filepath.Join(os.Getenv("HOME"), ".yggdrasil", "yggdrasil.sock"),
	}
	for _, sock := range socketPaths {
		if _, err := os.Stat(sock); err == nil {
			return true
		}
	}

	// 3. Проверка через yggdrasilctl
	cmd := fmt.Sprintf("%sctl getself", strings.TrimSuffix(path, "yggdrasil"))
	if _, err := execCommand(cmd); err == nil {
		return true
	}

	return false
}

// execCommand выполняет команду и возвращает вывод
func execCommand(cmd string) (string, error) {
	// Упрощённая реализация для MVP
	// В продакшене использовать os/exec
	return "", fmt.Errorf("not implemented")
}
