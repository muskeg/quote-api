// Package main implements a simple RESTful API for managing quotes.
// It provides endpoints to retrieve all quotes, get a random quote,
// fetch a quote by ID, and add new quotes. Quotes are persisted in a JSON file.
package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// quote represents a single quote entity with a unique ID and quote text.
type quote struct {
	ID    string `json:"id"`
	Quote string `json:"quote"`
}

// quotes is the in-memory storage for all quotes loaded from the JSON file.
var quotes = []quote{
}

// quotesMutex protects the quotes slice from concurrent access during reloads
var quotesMutex sync.RWMutex

func main() {
	const (
		ColorYellow = "\033[33m" // Yellow foreground for warning messages
		ColorReset  = "\033[0m"  // Reset to default color
	)

	viper.SetConfigName("config")
	viper.AddConfigPath("./data")
	err := viper.ReadInConfig()

	// Handle errors
	if err != nil {
	panic(fmt.Errorf("fatal error config file: %w", err))
}
	// Load the quotes from the JSON file
	loadedQuotes, err := loadFromJSON("./data/quotes.json")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Loaded quotes: %+v\n\n", loadedQuotes)
	quotes = loadedQuotes
	
	// Display a warning if no quotes were found in the file
	if len(quotes) == 0 {
		warningMessage := "No quotes found in quotes.json."
		fmt.Printf("%sWARNING: %s%s\n", ColorYellow, warningMessage, ColorReset)
		fmt.Println("add quotes using the POST /quotes endpoint.")
		fmt.Println("Example:")
		fmt.Println(`curl -X POST http://localhost:8080/quotes \`)
		fmt.Println(`  -H "Content-Type: application/json" \`)
		fmt.Println(`  -d '{"quote": "If it builds, it ships."}'`) 
	}

	// Set up Gin router with API endpoints
	trustedProxies := viper.GetStringSlice("trustedProxies")
	if len(trustedProxies) == 0 {
		trustedProxies = nil // Disable trusted proxies if none are configured
	}
	router := gin.Default()
	router.SetTrustedProxies(trustedProxies)

	// Check for read-only mode
	readOnly := viper.GetBool("readOnly")
	if readOnly {
		fmt.Println("API running in read-only mode - POST requests will be disabled")
	}

	// Middleware to block POST requests in read-only mode
	router.Use(func(c *gin.Context) {
		if readOnly && c.Request.Method == "POST" {
			c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "API is in read-only mode"})
			c.Abort()
			return
		}
		c.Next()
	})

	router.GET("/quote", getRandomQuote)
	router.GET("/quotes", getQuotes)
	router.GET("/quote/:id", getQuoteByID)
	router.GET("/quote/next", getNextQuote)
	router.POST("/quotes", addQuote)
	router.POST("/reload", func(c *gin.Context) {
		if err := reloadQuotes(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reload quotes"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Quotes reloaded successfully"})
	})
	router.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	
	// Start periodic reload in a goroutine
	var wg sync.WaitGroup
	wg.Add(1)
	go periodicReload(&wg)
	
	router.Run()
}

// getQuotes handles GET /quotes requests by returning all available quotes.
// Returns a 200 OK with the JSON array of all quotes.
func getQuotes(c *gin.Context) {
	quotesMutex.RLock()
	defer quotesMutex.RUnlock()
	c.IndentedJSON(http.StatusOK, quotes)
}

// getRandomQuote handles GET /quote requests by returning a random quote.
// Returns a 200 OK with a randomly selected quote in JSON format.
func getRandomQuote(c *gin.Context) {
	quotesMutex.RLock()
	defer quotesMutex.RUnlock()
	
	if len(quotes) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "no quotes available"})
		return
	}
	
	randIndex := rand.Intn(len(quotes))
	c.IndentedJSON(http.StatusOK, quotes[randIndex])
}

// getQuoteByID handles GET /quote/:id requests by looking up a specific quote by its ID.
// Path parameter:
//   - id: The unique identifier of the quote
//
// Returns:
//   - 200 OK with the quote in JSON format if found
//   - 404 Not Found if no quote with the given ID exists
func getQuoteByID(c *gin.Context) {
	id := c.Param("id")
	
	quotesMutex.RLock()
	defer quotesMutex.RUnlock()
	
	for _, q := range quotes {
		if q.ID == id {
			c.IndentedJSON(http.StatusOK, q)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "quote not found"})
}

// addQuote handles POST /quotes requests to add a new quote.
// It automatically generates an ID based on the current number of quotes.
// The request body should contain a JSON object with a "quote" field.
//
// Request body:
//   - {"quote": "The quote text"}
//
// Returns:
//   - 201 Created with the newly created quote (including its ID) on success
//   - 400 Bad Request if the request body is invalid
func addQuote(c *gin.Context) {
	var quoteRequest struct {
		Quote string `json:"quote"`
	}

	if err := c.BindJSON(&quoteRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}
	
	quotesMutex.Lock()
	// Generate ID based on the slice length
	newID := fmt.Sprintf("%d", len(quotes)+1)

	// Create the new quote object
	newQuote := quote{
		ID:    newID,
		Quote: quoteRequest.Quote,
	}

	// Add to the in-memory quotes slice
	quotes = append(quotes, newQuote)
	quotesMutex.Unlock()
	
	c.IndentedJSON(http.StatusCreated, newQuote)

	// Persist the updated quotes to the JSON file
	if err := saveToJSON("./data/quotes.json", quotes); err != nil {
		panic(err)
	}
	fmt.Println("Saved quotes to quotes.json")
}

func getNextQuote(c *gin.Context) {
	quotesMutex.RLock()
	quoteLen := len(quotes)
	quotesMutex.RUnlock()
	
	if quoteLen == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "no quotes available"})
		return
	}

	// Read the quoteIndex cookie
	idx := 0
	if cookie, err := c.Cookie("quoteIndex"); err == nil {
		var parsed int
		if _, err := fmt.Sscanf(cookie, "%d", &parsed); err == nil && parsed >= 0 {
			idx = parsed
		}
	}

	// Serve the quote at idx (wrap around)
	quoteIdx := idx % quoteLen
	
	// Set the cookie for the next index BEFORE sending the response
	nextIdx := (quoteIdx + 1) % quoteLen
	// maxAge: 30 days, path: /, domain: (leave empty for current domain), 
	// secure: false (allow HTTP), httpOnly: false (allow curl to store)
	c.SetCookie("quoteIndex", fmt.Sprintf("%d", nextIdx), 3600*24*30, "/", "", false, false)
	
	// Return flattened structure with nextIndex, id, and quote
	quotesMutex.RLock()
	response := gin.H{
		"nextIndex": nextIdx,
		"id": quotes[quoteIdx].ID,
		"quote": quotes[quoteIdx].Quote,
	}
	quotesMutex.RUnlock()
	
	c.IndentedJSON(http.StatusOK, response)
}

// saveToJSON persists the quotes slice to a JSON file on disk.
// Parameters:
//   - filename: The name of the file to save to
//   - users: The slice of quotes to save
//
// Returns:
//   - error: Any error encountered during the save operation
func saveToJSON(filename string, users []quote) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(users)
}

// loadFromJSON loads quotes from a JSON file on disk.
// If the file doesn't exist, it creates a new file with an empty array.
// Parameters:
//   - filename: The name of the file to load from
//
// Returns:
//   - []quote: The loaded quotes or an empty slice if the file is empty or new
//   - error: Any error encountered during the load operation
func loadFromJSON(filename string) ([]quote, error) {
	var loadedQuotes []quote

	// Check if the file exists
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		// File doesn't exist, create it with empty array
		file, err := os.Create(filename)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		// Write empty JSON array
		_, err = file.WriteString("[]")
		if err != nil {
			return nil, err
		}

		return loadedQuotes, nil
	} else if err != nil {
		// Some other error occurred
		return nil, err
	}

	// File exists, open and decode
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&loadedQuotes)
	return loadedQuotes, err
}

// periodicReload periodically reloads the quotes from the JSON file at the specified interval.
// It runs in a separate goroutine and updates the in-memory quotes slice.
// The reload interval is configured via the "reloadInterval" setting (in seconds).
func periodicReload(wg *sync.WaitGroup) {
	defer wg.Done()

	// Get the reload interval from the config (default: 60 seconds)
	interval := viper.GetInt("reloadInterval")
	if interval <= 0 {
		interval = 60
	}

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		// Reload the quotes from the JSON file
		loadedQuotes, err := loadFromJSON("./data/quotes.json")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reloading quotes: %v\n", err)
			continue
		}

		// Update the in-memory quotes slice
		quotesMutex.Lock()
		quotes = loadedQuotes
		quotesMutex.Unlock()

		fmt.Printf("Reloaded %d quotes from quotes.json\n", len(loadedQuotes))
	}
}

// reloadQuotes reloads the quotes from the JSON file and updates the in-memory slice.
// This function is thread-safe and can be called from any goroutine.
func reloadQuotes() error {
	loadedQuotes, err := loadFromJSON("./data/quotes.json")
	if err != nil {
		return err
	}

	quotesMutex.Lock()
	quotes = loadedQuotes
	quotesMutex.Unlock()

	fmt.Printf("Manually reloaded %d quotes from quotes.json\n", len(loadedQuotes))
	return nil
}
