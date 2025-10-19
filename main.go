package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type quote struct {
	ID    string `json:"id"`
	Quote string `json:"quote"`
}

var quotes = []quote{
	// {ID: "1", Quote: "I'm not a fan of using CI as a way for Software Development"},
	// {ID: "2", Quote: "これを読めるのはちょっとムズ"},
	// {ID: "3", Quote: "The quick brown fox jumps over the lazy dog."},
}

func main() {
	const (
		ColorYellow = "\033[33m" // Yellow foreground
		ColorReset  = "\033[0m"  // Reset to default color
	)

	// Load the slice from the file
	loadedQuotes, err := loadFromJSON("quotes.json")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Loaded quotes: %+v\n", loadedQuotes)
	quotes = loadedQuotes
	
	if len(quotes) == 0 {
		warningMessage := "No quotes found in quotes.json."
		fmt.Printf("%sWARNING: %s%s\n", ColorYellow, warningMessage, ColorReset)
		fmt.Println("add quotes using the POST /quotes endpoint.")
		fmt.Println("Example curl command:")
		fmt.Println(`curl -X POST http://localhost:8080/quotes \`)
		fmt.Println(`  -H "Content-Type: application/json" \`)
		fmt.Println(`  -d '{"quote": "The only way to do great work is to love what you do."}'`) 
	}

	// Set up Gin router
	router := gin.Default()
	router.GET("/quote", getRandomQuote)
	router.GET("/quotes", getQuotes)
	router.GET("/quote/:id", getQuoteByID)
	router.POST("/quotes", addQuote)
	router.Run(":8080")
}

func getQuotes(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, quotes)
}

func getRandomQuote(c *gin.Context) {
	randIndex := rand.Intn(len(quotes))
	c.IndentedJSON(http.StatusOK, quotes[randIndex])
}

func getQuoteByID(c *gin.Context) {
	id := c.Param("id")
	for _, q := range quotes {
		if q.ID == id {
			c.IndentedJSON(http.StatusOK, q)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "quote not found"})
}

func addQuote(c *gin.Context) {
	var quoteRequest struct {
		Quote string `json:"quote"`
	}

	if err := c.BindJSON(&quoteRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}
	newID := fmt.Sprintf("%d", len(quotes)+1)

	newQuote := quote{
		ID:    newID,
		Quote: quoteRequest.Quote,
	}

	quotes = append(quotes, newQuote)
	c.IndentedJSON(http.StatusCreated, newQuote)

	if err := saveToJSON("quotes.json", quotes); err != nil {
		panic(err)
	}
	fmt.Println("Saved quotes to quotes.json")
}

func saveToJSON(filename string, users []quote) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(users)
}

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