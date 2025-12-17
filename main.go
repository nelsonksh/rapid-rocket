package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

// External API Models
type TransactionCounts struct {
	Total           int `json:"total"`
	MintAccessToken int `json:"mint_access_token"`
	CreateCourse    int `json:"create_course"`
}

type AnalyticsAPIResponse struct {
	Count TransactionCounts `json:"count"`
}

type TransactionAPIResponse struct {
	TxHash      string   `json:"tx_hash"`
	Types       []string `json:"types"`
	SubmittedAt string   `json:"submitted_at"`
}

// Data Models
type Analytics struct {
	TotalTransactions int
	ActiveAddresses   int
	TotalBlocks       int
	NetworkLoad       int
	AvgBlockTime      int
	TotalValue        string
	CourseCount       int
	ProjectCount      int
}

type Transaction struct {
	Hash      string
	Timestamp string
	Amount    string
	Types     []string
}

type Contribution struct {
	ID        string
	Title     string
	Timestamp string
	Author    string
}

type SearchResult struct {
	Type     string
	ID       string
	Title    string
	Subtitle string
	Details  string
	Link     string
}

func main() {
	// Router
	mux := http.NewServeMux()

	// Routes
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/docs", docsHandler)

	// HTMX Fragment API
	mux.HandleFunc("/api/analytics", analyticsHandler)
	mux.HandleFunc("/api/transactions", transactionsHandler)
	mux.HandleFunc("/api/contributions", contributionsHandler)
	mux.HandleFunc("/search", searchHandler)

	// Static files
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	// Server config
	addr := ":8080"
	log.Printf("🚀 Server starting on http://localhost%s", addr)

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

// Handlers

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	tmplPath := filepath.Join("views", "index.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, "Could not load template", http.StatusInternalServerError)
		log.Printf("Template error: %v", err)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		log.Printf("Template execution error: %v", err)
	}
}

func docsHandler(w http.ResponseWriter, r *http.Request) {
	tmplPath := filepath.Join("views", "docs.html")
	http.ServeFile(w, r, tmplPath)
}

func analyticsHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Fetch from Real API
	resp, err := http.Get("https://preprod.andamioscan.andamio.space/v2/transactions/count")
	if err != nil {
		log.Printf("API fetch error: %v", err)
		http.Error(w, "Failed to fetch data", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// 2. Decode JSON
	var apiData AnalyticsAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiData); err != nil {
		log.Printf("JSON decode error: %v", err)
		http.Error(w, "Failed to parse data", http.StatusInternalServerError)
		return
	}

	// 3. Map to View Model
	// We only map the fields provided by the API: Total, MintAccessToken (Users), CreateCourse (Courses)
	// Consistently maintain Mock or Placeholder values for missing data to avoid empty cards
	data := Analytics{
		TotalTransactions: apiData.Count.Total,
		ActiveAddresses:   apiData.Count.MintAccessToken, // Mapping MintAccessToken to Users
		TotalBlocks:       8945234,                       // Mock (Missing in API)
		NetworkLoad:       78,                            // Mock (Missing in API)
		AvgBlockTime:      20,                            // Mock (Missing in API)
		TotalValue:        "45.2B ADA",                   // Mock (Missing in API)
		CourseCount:       apiData.Count.CreateCourse,    // Mapping CreateCourse
		ProjectCount:      8,                             // Mock (Missing in API)
	}

	tmplPath := filepath.Join("views", "analytics.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Could not load template", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Template execution error: %v", err)
	}
}

func transactionsHandler(w http.ResponseWriter, r *http.Request) {
	// Mock Transactions
	txs := []Transaction{
		{Hash: "8a9b0c1d...", Timestamp: "5 minutes ago", Amount: "1,250.50 ADA", Types: []string{"Payment", "Fee"}},
		{Hash: "7z8a9b0c...", Timestamp: "8 minutes ago", Amount: "450.00 ADA", Types: []string{"Payment"}},
		{Hash: "6y7z8a9b...", Timestamp: "12 minutes ago", Amount: "2,100.25 ADA", Types: []string{"Stake"}},
	}

	tmplPath := filepath.Join("views", "transactions.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Could not load template", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, txs); err != nil {
		log.Printf("Template execution error: %v", err)
	}
}

func contributionsHandler(w http.ResponseWriter, r *http.Request) {
	// Mock Contributions
	contribs := []Contribution{
		{ID: "cnt_1", Title: "Completed Module 1", Timestamp: "2 minutes ago", Author: "Student #42"},
		{ID: "cnt_2", Title: "Submitted Assignment", Timestamp: "15 minutes ago", Author: "Student #88"},
		{ID: "cnt_3", Title: "Updated Project Info", Timestamp: "1 hour ago", Author: "Project #12"},
	}

	tmplPath := filepath.Join("views", "contributions.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Could not load template", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, contribs); err != nil {
		log.Printf("Template execution error: %v", err)
	}
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		w.Write([]byte(""))
		return
	}

	var results []SearchResult
	qLower := strings.ToLower(query)

	// 1. Transaction Search (Real API)
	if strings.HasPrefix(qLower, "tx_") || len(query) == 64 {
		url := fmt.Sprintf("https://preprod.andamioscan.andamio.space/v2/transactions/%s", query)
		resp, err := http.Get(url)
		if err == nil && resp.StatusCode == http.StatusOK {
			defer resp.Body.Close()
			var txRes []TransactionAPIResponse
			if err := json.NewDecoder(resp.Body).Decode(&txRes); err == nil && len(txRes) > 0 {
				tx := txRes[0]
				results = append(results, SearchResult{
					Type:     "transaction",
					ID:       tx.TxHash,
					Title:    "Transaction",
					Subtitle: tx.TxHash,
					Details:  fmt.Sprintf("%v • %s", tx.Types, tx.SubmittedAt),
					Link:     "#", // Link to details page when implemented
				})
			}
		}
	} else if strings.HasPrefix(qLower, "addr") || len(query) > 50 {
		results = append(results, SearchResult{
			Type:     "address",
			ID:       "addr_001",
			Title:    "Address",
			Subtitle: query,
			Details:  "Balance: 125,450.75 ADA • 342 transactions",
			Link:     "#",
		})
	} else {
		// Generic fallback or numeric check
		results = append(results, SearchResult{
			Type:     "block",
			ID:       "block_sample",
			Title:    "Block",
			Subtitle: "#8945234",
			Details:  "245 transactions • 64.5 KB",
			Link:     "#",
		})
	}

	data := struct {
		Count   int
		Query   string
		Results []SearchResult
	}{
		Count:   len(results),
		Query:   query,
		Results: results,
	}

	tmplPath := filepath.Join("views", "search.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Could not load template", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Template execution error: %v", err)
	}
}
