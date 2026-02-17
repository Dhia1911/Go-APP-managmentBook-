package main

import (
	"encoding/json"
	"fmt"
	"net/http" // http client and sever implemenation
	"strconv"  // functions , convert strings to toher type s
	"strings"  // manulpation string
	"sync"
)

// Book management

type Book struct {
	// json syntax `json:"field name "`

	ID              int    `json:"id"`
	Title           string `json:"title"`
	Auther          string `json:"auther"`
	Publicationyear int    `json:"publication_year"`
}

var (
	books = make(map[int]Book)
	mutex = &sync.Mutex{}
)

func main() {

	http.HandleFunc("/books", getBooks)
	http.HandleFunc("/books/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getBook(w, r)
		case http.MethodPost:
			createBook(w, r)
		case http.MethodPut:
			updateBook(w, r)
		case http.MethodDelete:
			deleteBook(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)

}

// Functions of our CRUD

func getBooks(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock() // ensure function after it's done to unlock
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

// get book by id
func getBook(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	idStr := strings.TrimPrefix(r.URL.Path, "/books/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID ", http.StatusBadRequest)
		return
	}

	book, exists := books[id]
	if !exists {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)

}

// create a new book
func createBook(w http.ResponseWriter, r *http.Request) {

	mutex.Lock()
	defer mutex.Unlock()

	var book Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, "Invalid input ", http.StatusBadRequest)
		return
	}

	// user didnt add an id
	if book.ID == 0 {
		http.Error(w, "Book id is required  ", http.StatusBadRequest)
		return
	}

	books[book.ID] = book

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)

}

// update book

func updateBook(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	idStr := strings.TrimPrefix(r.URL.Path, "/books/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID ", http.StatusBadRequest)
		return
	}

	var updateBook Book
	if err := json.NewDecoder(r.Body).Decode(&updateBook); err != nil {
		http.Error(w, "Invalid input ", http.StatusBadRequest)
		return
	}

	if updateBook.ID != id {
		http.Error(w, "Book id mismatch", http.StatusBadRequest)
		return
	}

	books[id] = updateBook

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updateBook)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {

	mutex.Lock()
	defer mutex.Unlock()

	idStr := strings.TrimPrefix(r.URL.Path, "/books/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID ", http.StatusBadRequest)
		return
	}

	if _, exists := books[id]; !exists {
		http.Error(w, "Book not found ", http.StatusNotFound)
		return
	}

	delete(books, id)

	w.WriteHeader(http.StatusNoContent)

}