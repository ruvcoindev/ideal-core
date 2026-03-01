package main

import (
	"context"
	"crypto/ed25519"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"ideal-core/pkg/crypto"
	"ideal-core/pkg/yggdrasil-go"
)

var (
	dataDir    = flag.String("data", "~/.ideal-core", "Directory for keys and data")
	bootstrap  = flag.String("bootstrap", "", "Comma-separated Yggdrasil peers to bootstrap with")
	genKey     = flag.Bool("genkey", false, "Generate new keypair and exit")
	exportBackup = flag.String("export-backup", "", "Export encrypted backup to file (requires -password)")
	importBackup = flag.String("import-backup", "", "Import encrypted backup from file (requires -password)")
	password   = flag.String("password", "", "Password for backup encryption/decryption")
	yggPath    = flag.String("yggdrasil", "/usr/bin/yggdrasil", "Path to yggdrasil binary")
	port       = flag.String("port", "8080", "Port for local web server")
	bindAddr   = flag.String("bind", "127.0.0.1", "Address to bind web server (127.0.0.1 for local, 0.0.0.0 for LAN)")
)

func main() {
	flag.Parse()

	dir := os.ExpandEnv(*dataDir)
	if err := os.MkdirAll(dir, 0700); err != nil {
		log.Fatalf("Failed to create data dir: %v", err)
	}

	// Handle backup export/import
	if *exportBackup != "" {
		handleExportBackup(dir, *exportBackup, *password)
		return
	}
	if *importBackup != "" {
		handleImportBackup(*importBackup, *password)
		return
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
		fmt.Printf("   App Yggdrasil IP: %s\n", crypto.DeriveYggdrasilIP(kp.PublicKey))
		fmt.Printf("   ⚠️  %s", crypto.SecurityWarning())
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
	appYggIP := crypto.DeriveYggdrasilIP(kp.PublicKey)
	fmt.Printf("🌐 App Yggdrasil IP: %s\n", appYggIP)
	fmt.Printf("🌐 Web UI: http://%s:%s\n", *bindAddr, *port)

	// Check Yggdrasil availability
	yggAvailable := checkYggdrasil(*yggPath)
	if !yggAvailable {
		log.Printf("⚠️  Yggdrasil not found at %s, using fallback transport", *yggPath)
	} else {
		fmt.Printf("✅ Yggdrasil service detected\n")
	}

	// Connect to Yggdrasil mesh
	ygg, err := yggdrasil.NewClient(kp.ToHex(), *yggPath, yggAvailable)
	if err != nil {
		log.Fatalf("Yggdrasil init failed: %v", err)
	}
	defer ygg.Close()

	if *bootstrap != "" {
		peers := strings.Split(*bootstrap, ",")
		if err := ygg.Bootstrap(peers); err != nil {
			log.Printf("⚠️  Bootstrap warning: %v", err)
		}
	}

	// Start local web server (BIND TO LOCALHOST, NOT YGG IP)
	go startWebServer(*bindAddr, *port, kp, appYggIP)

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
			fmt.Printf("📥 Received %d bytes (encrypted)\n", len(msg))
			return nil
		}); err != nil && err != context.Canceled {
			log.Printf("Receive error: %v", err)
		}
	}()

	fmt.Println("✅ Node + Server running. Press Ctrl+C to stop.")
	fmt.Println("🔐 Security: Your private key is stored encrypted at rest.")

	// Keep alive
	<-ctx.Done()
	fmt.Println("👋 Node stopped.")
}

// startWebServer запускает локальный веб-сервер
func startWebServer(bindAddr, port string, kp *crypto.KeyPair, appYggIP string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "🗝️ IDEAL CORE\nNode ID: %s\nApp Yggdrasil IP: %s\n", kp.ToHex()[:16], appYggIP)
	})
	http.HandleFunc("/api/keys", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"public_key":"%s","app_yggdrasil_ip":"%s"}`, kp.ToHex(), appYggIP)
	})
	
	addr := fmt.Sprintf("%s:%s", bindAddr, port)
	log.Printf("🌐 Web server listening on http://%s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Printf("Web server error: %v", err)
	}
}

func handleExportBackup(dir, backupPath, password string) {
	if password == "" {
		log.Fatal("Password required for backup encryption. Use -password flag.")
	}
	keyPath := filepath.Join(dir, "private.key")
	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		log.Fatalf("No private key found at %s. Generate one first with -genkey", keyPath)
	}
	priv, err := crypto.LoadPrivateKey(keyPath)
	if err != nil {
		log.Fatalf("Failed to load private key: %v", err)
	}
	pub := priv.Public().(ed25519.PublicKey)
	kp := &crypto.KeyPair{PrivateKey: priv, PublicKey: pub}
	backup, err := kp.ExportEncryptedBackup(password)
	if err != nil {
		log.Fatalf("Failed to create encrypted backup: %v", err)
	}
	if err := crypto.SaveEncryptedBackup(backup, backupPath); err != nil {
		log.Fatalf("Failed to save backup: %v", err)
	}
	fmt.Printf("✅ Encrypted backup saved to: %s\n", backupPath)
	fmt.Printf("🔐 Keep this file safe. To restore: -import-backup %s -password YOUR_PASSWORD\n", backupPath)
}

func handleImportBackup(backupPath, password string) {
	if password == "" {
		log.Fatal("Password required for backup decryption. Use -password flag.")
	}
	backup, err := crypto.LoadEncryptedBackup(backupPath)
	if err != nil {
		log.Fatalf("Failed to load backup: %v", err)
	}
	kp, err := crypto.ImportEncryptedBackup(backup, password)
	if err != nil {
		log.Fatalf("Failed to decrypt backup: %v", err)
	}
	fmt.Printf("✅ Backup imported successfully!\n")
	fmt.Printf("   Public ID: %s\n", kp.ToHex())
	fmt.Printf("   App Yggdrasil IP: %s\n", crypto.DeriveYggdrasilIP(kp.PublicKey))
}

func checkYggdrasil(path string) bool {
	if path == "" {
		return false
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
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
	return false
}
