package main

import (
	"context"
	"crypto/ed25519"
	"encoding/json"
	"flag"
	"fmt"
	"ideal-core/pkg/crypto"
	"ideal-core/pkg/journal"
	"ideal-core/pkg/yggdrasil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

var (
	dataDir    = flag.String("data", "~/.ideal-core", "Directory for keys and data")
	bootstrap  = flag.String("bootstrap", "", "Comma-separated Yggdrasil peers to bootstrap with")
	genKey     = flag.Bool("genkey", false, "Generate new keypair and exit")
	port       = flag.String("port", "8080", "Port for local web server")
	bindAddr   = flag.String("bind", "127.0.0.1", "Address to bind web server")
	ollamaHost = flag.String("ollama", "http://localhost:11434", "Ollama API host")
	ollamaModel= flag.String("ollama-model", "bge-m3", "Ollama model for embeddings")
	useOllama  = flag.Bool("use-ollama", false, "Enable Ollama embeddings (requires running Ollama server)")
)

// Global instances
var (
	journalInstance *journal.Journal
	keyPair         *crypto.KeyPair
)

func main() {
	flag.Parse()

	dir := os.ExpandEnv(*dataDir)
	if err := os.MkdirAll(dir, 0700); err != nil {
		log.Fatalf("Failed to create data dir: %v", err)
	}

	// Key management
	keyPath := filepath.Join(dir, "private.key")
	pubKeyPath := filepath.Join(dir, "public.key")
	var err error

	if *genKey {
		keyPair, err = crypto.GenerateKeyPair()
		if err != nil {
			log.Fatalf("Key generation failed: %v", err)
		}
		if err := crypto.SavePrivateKey(keyPair.PrivateKey, keyPath); err != nil {
			log.Fatalf("Failed to save private key: %v", err)
		}
		if err := os.WriteFile(pubKeyPath, keyPair.PublicKey, 0644); err != nil {
			log.Fatalf("Failed to save public key: %v", err)
		}
		fmt.Printf("✅ New keypair generated:\n")
		fmt.Printf("   Public ID: %s\n", keyPair.ToHex())
		fmt.Printf("   Yggdrasil IP: %s\n", crypto.DeriveYggdrasilIP(keyPair.PublicKey))
		fmt.Printf("   ⚠️  %s", crypto.SecurityWarning())
		fmt.Printf("   Private key saved to: %s\n", keyPath)
		return
	}

	// Load or create key
	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		fmt.Println("🔑 No key found, generating new one...")
		keyPair, err = crypto.GenerateKeyPair()
		if err != nil {
			log.Fatalf("Key generation failed: %v", err)
		}
		if err := crypto.SavePrivateKey(keyPair.PrivateKey, keyPath); err != nil {
			log.Fatalf("Failed to save private key: %v", err)
		}
		if err := os.WriteFile(pubKeyPath, keyPair.PublicKey, 0644); err != nil {
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
		keyPair = &crypto.KeyPair{PrivateKey: priv, PublicKey: ed25519.PublicKey(pub)}
	}

	fmt.Printf("🗝️  Node ID: %s\n", keyPair.ToHex()[:16]+"...")
	fmt.Printf("🌐 App Yggdrasil IP: %s\n", crypto.DeriveYggdrasilIP(keyPair.PublicKey))

	// Initialize journal with Ollama option
	journalCfg := journal.JournalConfig{
		DataDir:        dir,
		OllamaHost:     *ollamaHost,
		OllamaModel:    *ollamaModel,
		UseOllamaEmbed: *useOllama,
		DefaultMode:    journal.EntryTypeCBT,
	}
	journalInstance, err = journal.NewJournal(journalCfg)
	if err != nil {
		log.Fatalf("Failed to initialize journal: %v", err)
	}
	fmt.Printf("📓 Journal initialized: %d entries loaded\n", len(journalInstance.GetEntries(journal.EntryFilters{})))

	// Yggdrasil client (optional)
	yggAvailable := checkYggdrasil("/usr/bin/yggdrasil")
	if !yggAvailable {
		log.Printf("⚠️  Yggdrasil not found, running in local-only mode")
	} else {
		fmt.Println("✅ Yggdrasil service detected")
	}
	ygg, err := yggdrasil.NewClient(keyPair.ToHex(), "/usr/bin/yggdrasil", yggAvailable)
	if err != nil {
		log.Printf("⚠️  Yggdrasil client init failed: %v", err)
	} else {
		defer ygg.Close()
	}

	// Start web server
	go startWebServer(*bindAddr, *port)

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

	// Yggdrasil message listener (optional)
	if ygg != nil {
		go func() {
			if err := ygg.Receive(ctx, func(msg []byte) error {
				fmt.Printf("📥 Received %d bytes (encrypted)\n", len(msg))
				// TODO: decrypt and process message
				return nil
			}); err != nil && err != context.Canceled {
				log.Printf("Receive error: %v", err)
			}
		}()
	}

	fmt.Printf("✅ Node + Server running at http://%s:%s\n", *bindAddr, *port)
	fmt.Println("🔐 Security: Your private key is stored encrypted at rest.")
	fmt.Println("📓 Journal: CBT + Gratitude modes with semantic search")
	if *useOllama {
		fmt.Printf("🤖 Ollama: %s/%s\n", *ollamaHost, *ollamaModel)
	}

	// Keep alive
	<-ctx.Done()
	fmt.Println("👋 Node stopped.")
}

// startWebServer запускает HTTP-сервер с API для дневника
func startWebServer(bindAddr, port string) {
	// Static files
	fs := http.FileServer(http.Dir("web"))
	http.Handle("/", fs)
	http.Handle("/web/", http.StripPrefix("/web/", fs))

	// Journal API endpoints
	http.HandleFunc("/api/journal/stats", handleJournalStats)
	http.HandleFunc("/api/journal/entries", handleJournalEntries)
	http.HandleFunc("/api/journal/search", handleJournalSearch)
	http.HandleFunc("/api/journal/export/md", handleJournalExportMD)

	// Health check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	addr := fmt.Sprintf("%s:%s", bindAddr, port)
	log.Printf("🌐 Web server listening on http://%s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Web server error: %v", err)
	}
}

// handleJournalStats — GET /api/journal/stats
func handleJournalStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	stats := journalInstance.GetCombinedStats()
	json.NewEncoder(w).Encode(stats)
}

// handleJournalEntries — GET/POST /api/journal/entries
func handleJournalEntries(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	switch r.Method {
	case http.MethodGet:
		// Parse filters from query params
		filters := journal.EntryFilters{
			Type:     r.URL.Query().Get("type"),
			PersonID: r.URL.Query().Get("person"),
			Phase:    r.URL.Query().Get("phase"),
			Tag:      r.URL.Query().Get("tag"),
		}
		entries := journalInstance.GetEntries(filters)
		json.NewEncoder(w).Encode(entries)
		
	case http.MethodPost:
		var entry journal.ThoughtEntry
		if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		entry.Timestamp = time.Now()
		
		// Route to appropriate add method based on type
		var err error
		switch entry.Type {
		case journal.EntryTypeCBT:
			err = journalInstance.AddCBTEntry(
				entry.Situation,
				entry.AutomaticThought,
				entry.Emotions,
				entry.Intensity,
			)
		case journal.EntryTypeGratitude:
			err = journalInstance.AddGratitudeEntry(entry.GratitudeItems, entry.Notes)
		default:
			// Fallback to universal method
			err = journalInstance.AddEntry(entry)
		}
		
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(entry)
		
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleJournalSearch — GET /api/journal/search?q=...
func handleJournalSearch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	query := r.URL.Query().Get("q")
	limit := 10
	if l := r.URL.Query().Get("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}
	
	entries := journalInstance.SearchByMeaning(query, limit)
	json.NewEncoder(w).Encode(entries)
}

// handleJournalExportMD — GET /api/journal/export/md
func handleJournalExportMD(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	// Generate markdown to temp file
	tmpPath := filepath.Join(os.TempDir(), "journal_export.md")
	if err := journalInstance.ExportToMarkdown(tmpPath); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "text/markdown")
	w.Header().Set("Content-Disposition", "attachment; filename=journal.md")
	http.ServeFile(w, r, tmpPath)
	
	// Cleanup
	defer os.Remove(tmpPath)
}

// checkYggdrasil проверяет доступность сервиса Yggdrasil
func checkYggdrasil(path string) bool {
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
