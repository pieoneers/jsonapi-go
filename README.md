# jsonapi-go
# Go jsonapi client

[![Go Report Card](https://goreportcard.com/badge/github.com/pieoneers/jsonapi-go)](https://goreportcard.com/report/github.com/pieoneers/jsonapi-go)
[![GoDoc](https://godoc.org/github.com/pieoneers/jsonapi-go?status.svg)](https://godoc.org/github.com/pieoneers/jsonapi-go)

Lightweight [JSON API](https://jsonapi.org/) implementation for Go.

### Installing

``` go get -u "github.com/pieoneers/jsonapi-go" ```

### Running the tests
Go to jsonapi-go package directory and run:

```go test```

### Usage
#### Marshal
For instance we have Go program what implements a simple library, we have `Book` and `Author` structs:
```go
type Book struct {
	ID              string   
	AuthorID        string   
	Title           string   
	PublicationDate time.Time
}

type Author struct {
	ID        string
	FirstName string
	LastName  string
}
```

We want to produce the JSON representation of the data. Let's modify the struct by adding json tags to it and implement functions `GetID` and `GetType` as required for [MarshalResourceIdentifier interface](https://godoc.org/github.com/pieoneers/jsonapi-go#MarshalResourceIdentifier), one more function `GetData` required for [MarshalData interface](https://godoc.org/github.com/pieoneers/jsonapi-go#MarshalData).

```go
type Book struct {
	ID              string    `json:"-"`
	AuthorID        string    `json:"-"`
	Title           string    `json:"title"`
	PublicationDate time.Time `json:"publication_date"`
}

func (b Book) GetID() string {
	return b.ID
}

func (b Book) GetType() string {
	return "books"
}

func (b Book) GetData() interface{} {
	return b
}

type Author struct {
	ID        string `json:"-"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (a Author) GetID() string {
	return a.ID
}

func (a Author) GetType() string {
	return "authors"
}

func (a Author) GetData() interface{} {
	return a
}
```

By running [Marshal](https://godoc.org/github.com/pieoneers/jsonapi-go#Marshal) function for `Book` and `Author` the output will be a `[]byte` of json data.

Initial data:
```go
alan := Author{
	ID:        "1",
	FirstName: "Alan A. A.",
	LastName:  "Donovan",
}

publicationDate, _ := time.Parse(time.RFC3339, "2015-01-01T00:00:00Z")

book := Book{
	ID:              "1",
	Title:           "Go Programming Language",
	AuthorID:        alan.ID,
	PublicationDate: publicationDate,
}
```
Running `Marshal`:
```go
bookJSON, _ := jsonapi.Marshal(book1)
authorJSON, _ := jsonapi.Marshal(alan)
```

Output:
```json
{
	"data": {
		"type": "books",
		"id": "1",
		"attributes": {
			"title": "Go Programming Language",
			"publication_date": "2015-01-01T00:00:00Z"
		}
	}
}
```

and
```json
{
	"data": {
		"type": "authors",
		"id": "1",
		"attributes": {
			"first_name": "Alan A. A.",
			"last_name": "Donovan"
		}
	}
}
```

##### Relationships
Add `relationships` to the resource is easy to do by implementing [MarshalRelationships interface](https://godoc.org/github.com/pieoneers/jsonapi-go#MarshalRelationships)
e.g. for `Book` we will add `GetRelationships` function.
```go
func (b Book) GetRelationships() map[string]interface{} {
	relationships := make(map[string]interface{})

	relationships["author"] = jsonapi.ResourceObjectIdentifier{
		ID:   b.AuthorID,
		Type: "authors",
	}

	return relationships
}
```

The `Marshal` output will be:

```json
{
	"data": {
		"type": "books",
		"id": "1",
		"attributes": {
			"title": "Go Programming Language",
			"publication_date": "2015-01-01T00:00:00Z"
		},
		"relationships": {
			"author": {
				"data": {
					"type": "authors",
					"id": "1"
				}
			}
		}
	}
}
```

##### Included
Adding to JSON document `included` is easy by adding function `GetIncluded` it will implement [MarshalIncluded interface](https://godoc.org/github.com/pieoneers/jsonapi-go#MarshalIncluded).
```go
func (b Book) GetIncluded() []interface{} {
	var included []interface{}

  //`authors` a global array with all authors, but it can be a DB or something else
	for _, author := range authors {
		if author.ID == b.AuthorID {
			included = append(included, author)
		}
	}

	return included
}
```

When JSON document is a collection of resources better way to implement additional data type and `Books` method `GetIncluded`, otherwise all same included resources may be included several times.
```go
type Books []Book

func (b Books) GetIncluded() []interface{} {
  //do something
}
```

`Book` with `included` will be looks like:
```json
{
	"data": {
		"type": "books",
		"id": "1",
		"attributes": {
			"title": "Go Programming Language",
			"publication_date": "2015-01-01T00:00:00Z"
		},
		"relationships": {
			"author": {
				"data": {
					"type": "authors",
					"id": "1"
				}
			}
		}
	},
	"included": [{
		"type": "authors",
		"id": "1",
		"attributes": {
			"first_name": "Alan A. A.",
			"last_name": "Donovan"
		}
	}]
}
```

##### Meta

`Book`'s `GetMeta` will implement [MarshalMeta interface](https://godoc.org/github.com/pieoneers/jsonapi-go#MarshalMeta).
`meta` section will be added to json document.

```go
func (b Book) GetMeta() interface{} {
	return Meta{ ReadCount: 42 }
}
```

the output:
```json
{
	"data": {
		"type": "books",
		"id": "1",
		"attributes": {
			"title": "Go Programming Language",
			"publication_date": "2015-01-01T00:00:00Z"
		},
		"meta": {
			"read_count": 42
		},
		"relationships": {
			"author": {
				"data": {
					"type": "authors",
					"id": "1"
				}
			}
		}
	},
	"included": [{
		"type": "authors",
		"id": "1",
		"attributes": {
			"first_name": "Alan A. A.",
			"last_name": "Donovan"
		}
	}],
	"meta": {
		"read_count": 42
	}
}
```

Here we can find that `meta` includes in resource object and into JSON document.
It happened, because resource and document may have their own `meta`, to avoid such behavior we should implement `Book` document data type and implement `GetMeta` for the `Book`'s document.

```go
type BookDocumentMeta struct {
	TotalCount int `json:"total_count"`
}

type BookDocument struct {
  Data Book
  Meta BookDocumentMeta
}

func (b BookDocument) GetData() interface{} {
	return b.Data
}

func (b BookDocument) GetMeta() interface{} {
	return b.Meta
}

...

bookJSON, _ := jsonapi.Marshal(BookDocument{
  Data: book, //book from previous examples
  Meta: BookDocumentMeta{ 17 }, //let's imagine that we have 17 books in our library
})
```

Output:

```json
{
	"data": {
		"type": "books",
		"id": "1",
		"attributes": {
			"title": "Go Programming Language",
			"publication_date": "2015-01-01T00:00:00Z"
		},
		"meta": {
			"read_count": 42
		},
		"relationships": {
			"author": {
				"data": {
					"type": "authors",
					"id": "1"
				}
			}
		}
	},
	"meta": {
		"total_count": 17
	}
}
```

#### Unmarshal

Now we will implement the examples, how to put values from JSON document into Go struct.
We will use the same data types `Book` and `Author`

But the set of functions will be different.
`Book`:
```go
type Book struct {
	ID              string    `json:"-"`
	Type            string    `json:"-"`
	AuthorID        string    `json:"-"`
	Title           string    `json:"title"`
	PublicationDate time.Time `json:"publication_date"`
}

func (b *Book) SetID(id string) error {
	b.ID = id
	return nil
}

func (b *Book) SetType(t string) error {
	b.Type = t
	return nil
}

func (b *Book) SetData(to func(target interface{}) error) error {
	return to(b)
}

```

Here is JSON data:
```json
{
  	"data": {
  		"type": "books",
  		"id": "1",
  		"attributes": {
  			"title": "Go Programming Language",
  			"publication_date": "2015-01-01T00:00:00Z"
  		},
      "relationships": {
        "author": {
          "data": {
            "type": "authors",
            "id": "1"
          }
        }
      }
  	}
  }
```
Let's call [Unmarshal](https://godoc.org/github.com/pieoneers/jsonapi-go#Unmarshal) function.


```go
  book := Book{}
	jsonapi.Unmarshal(bookJSON, &book) //bookJSON is []byte of JSON data.
```

At the end we will have Go struct:
```go
_ = Book{
  ID:	      "1",
  Type:	    "books",
  AuthorID:	"",
  Title:	   "Go Programming Language"
  PublicationDate:	2015-01-01 00:00:00 +0000 UTC,
}
```
But the AuthorID is empty, to set this relationship we should implement [UnmarshalRelationships](https://godoc.org/github.com/pieoneers/jsonapi-go#UnmarshalRelationships) interface, by creating function SetRelationships:
```go
func (b *Book) SetRelationships(relationships map[string]interface{}) error {
	if relationship, ok := relationships["author"]; ok {
		b.AuthorID = relationship.(*jsonapi.ResourceObjectIdentifier).ID
	}

	return nil
}
```
call Unmarshal again and look at the result.

##### Collection unmarshal

When you need to unmarshal collections, you should implement SetData function for collection datatype.
```go
type Books []Book

func (b *Books) SetData(to func(target interface{}) error) error {
	return to(b)
}
```
If resource `Book` data type have SetRelationships function, in the collection all relationships will be filled.
Example:

```go
var booksJSON = []byte(`
  {
  	"data":
    [
      {
    		"type": "books",
    		"id": "1",
    		"attributes": {
    			"title": "Go Programming Language",
    			"publication_date": "2015-01-01T00:00:00Z"
    		},
    		"relationships": {
    			"author": {
    				"data": {
    					"type": "authors",
    					"id": "1"
    				}
    			}
    		}
    	},
      {
    		"type": "books",
    		"id": "2",
    		"attributes": {
    			"title": "Learning Functional Programming in Go",
    			"publication_date": "2017-11-01T00:00:00Z"
    		},
    		"relationships": {
    			"author": {
    				"data": {
    					"type": "authors",
    					"id": "2"
    				}
    			}
    		}
    	},
      {
    		"type": "books",
    		"id": "3",
    		"attributes": {
    			"title": "Go in Action",
    			"publication_date": "2015-11-01T00:00:00Z"
    		},
    		"relationships": {
    			"author": {
    				"data": {
    					"type": "authors",
    					"id": "3"
    				}
    			}
    		}
    	}
    ]
  }
`)

books := Books{}
jsonapi.Unmarshal(booksJSON, &books)
```

The output:
```go
_ = Books{
  {  
    ID:	"1",
    Type:	"books",
    AuthorID:	"1",
    Title:	"Go Programming Language",
    PublicationDate:	2015-01-01 00:00:00 +0000 UTC,
  },
  {
    ID:	"2",
    Type:	"books",
    AuthorID:	"2",
    Title:	"Learning Functional Programming in Go",
    PublicationDate:	2017-11-01 00:00:00 +0000 UTC,
  },
  {
    ID:	"3",
    Type:	"books",
    AuthorID:	"3",
    Title:	"Go in Action",
    PublicationDate:	2015-11-01 00:00:00 +0000 UTC,
  },
}
```


### Example Library client and server applications.

server.go
```go
package main

import (
	"github.com/pieoneers/jsonapi-go"
	"log"
	"net/http"
	"time"
)

type Meta struct {
	Count int `json:"count"`
}

type Book struct {
	ID              string    `json:"-"`
	AuthorID        string    `json:"-"`
	Title           string    `json:"title"`
	PublicationDate time.Time `json:"publication_date"`
}

func (b Book) GetID() string {
	return b.ID
}

func (b Book) GetType() string {
	return "books"
}

func (b Book) GetData() interface{} {
	return b
}

func (b Book) GetRelationships() map[string]interface{} {
	relationships := make(map[string]interface{})

	relationships["author"] = jsonapi.ResourceObjectIdentifier{
		ID:   b.AuthorID,
		Type: "authors",
	}

	return relationships
}

func (b Book) GetIncluded() []interface{} {
	var included []interface{}

	for _, author := range authors {
		if author.ID == b.AuthorID {
			included = append(included, author)
		}
	}

	return included
}

type Books []Book

func (b Books) GetData() interface{} {
	return b
}

func (b Books) GetMeta() interface{} {
	return Meta{Count: len(books)}
}

func (b Books) GetIncluded() []interface{} {
	var included []interface{}

	authorsMap := make(map[string]Author)

	for _, book := range b {
		for _, author := range authors {
			if book.AuthorID == author.ID {
				authorsMap[author.ID] = author
			}
		}
	}

	for _, author := range authorsMap {
		included = append(included, author)
	}

	return included
}

type Author struct {
	ID        string `json:"-"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (a Author) GetID() string {
	return a.ID
}

func (a Author) GetType() string {
	return "authors"
}

func (a Author) GetData() interface{} {
	return a
}

func (a Author) GetIncluded() []interface{} {
	var included []interface{}

	for _, book := range books {
		if book.AuthorID == a.ID {
			included = append(included, book)
		}
	}

	return included
}

type Authors []Author

func (a Authors) GetMeta() interface{} {
	return Meta{Count: len(authors)}
}

func (a Authors) GetData() interface{} {
	return a
}

func (a Authors) GetIncluded() []interface{} {
	var included []interface{}

	booksMap := make(map[string]Book)

	for _, author := range a {
		for _, book := range books {
			if book.AuthorID == author.ID {
				booksMap[book.ID] = book
			}
		}
	}

	for _, book := range booksMap {
		included = append(included, book)
	}

	return included
}

var (
	authors Authors
	books   Books
)

func bookHandler(w http.ResponseWriter, req *http.Request) {

	id := req.URL.Path[len("/books/"):]

	for _, book := range books {
		if book.ID == id {
			bookData, err := jsonapi.Marshal(book)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/vnd.api+json")
			w.Write(bookData)
			w.WriteHeader(http.StatusOK)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}

func booksHandler(w http.ResponseWriter, req *http.Request) {

	booksData, err := jsonapi.Marshal(books)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/vnd.api+json")
	w.Write(booksData)
	w.WriteHeader(http.StatusOK)
}

func authorHandler(w http.ResponseWriter, req *http.Request) {

	id := req.URL.Path[len("/authors/"):]

	for _, author := range authors {
		if author.ID == id {
			authorData, err := jsonapi.Marshal(author)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/vnd.api+json")
			w.Write(authorData)
			w.WriteHeader(http.StatusOK)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}

func authorsHandler(w http.ResponseWriter, req *http.Request) {

	authorsData, err := jsonapi.Marshal(authors)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/vnd.api+json")
	w.Write(authorsData)
	w.WriteHeader(http.StatusOK)
}

func main() {
	var publicationDate time.Time

	alan := Author{ID: "1", FirstName: "Alan A. A.", LastName: "Donovan"}
	authors = append(authors, alan)

	lex := Author{ID: "2", FirstName: "Lex", LastName: "Sheehan"}
	authors = append(authors, lex)

	william := Author{ID: "3", FirstName: "William", LastName: "Kennedy"}
	authors = append(authors, william)

	publicationDate, _ = time.Parse(time.RFC3339, "2015-01-01T00:00:00Z")

	book1 := Book{
		ID:              "1",
		Title:           "Go Programming Language",
		AuthorID:        alan.ID,
		PublicationDate: publicationDate,
	}
	books = append(books, book1)

	publicationDate, _ = time.Parse(time.RFC3339, "2017-11-01T00:00:00Z")

	book2 := Book{
		ID:              "2",
		Title:           "Learning Functional Programming in Go",
		AuthorID:        lex.ID,
		PublicationDate: publicationDate,
	}
	books = append(books, book2)

	publicationDate, _ = time.Parse(time.RFC3339, "2015-11-01T00:00:00Z")

	book3 := Book{
		ID:              "3",
		Title:           "Go in Action",
		AuthorID:        william.ID,
		PublicationDate: publicationDate,
	}
	books = append(books, book3)

	http.HandleFunc("/books/", bookHandler)
	http.HandleFunc("/books", booksHandler)
	http.HandleFunc("/authors/", authorHandler)
	http.HandleFunc("/authors", authorsHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

client.go
```go
package main

import  (
  "fmt"
	"io/ioutil"
	"log"
	"net/http"
  "time"
  "github.com/pieoneers/jsonapi-go"
)

type Book struct {
  ID       string `json:"-"`
  Type     string `json:"-"`
  AuthorID string `json:"-"`
  Title    string `json:"title"`
  PublicationDate time.Time `json:"publication_date"`
}

func(b *Book) SetID(id string) error {
  b.ID = id
  return nil
}

func(b *Book) SetType(t string) error {
  b.Type = t
  return nil
}

func(b *Book) SetData(to func(target interface{}) error) error {
  return to(b)
}

func(b *Book) SetRelationships(relationships map[string]interface{}) error {
  if relationship, ok := relationships["author"]; ok {
    b.AuthorID = relationship.(*jsonapi.ResourceObjectIdentifier).ID
  }

  return nil
}

type Books []Book

func(b *Books) SetData(to func(target interface{}) error) error {
  return to(b)
}

type Author struct {
  ID string `json:"-"`
  Type string `json:"-"`
  FirstName string `json:"first_name"`
  LastName string `json:"last_name"`
}

func(a *Author) SetID(id string) error {
  a.ID = id
  return nil
}

func(a *Author) SetType(t string) error {
  a.Type = t
  return nil
}

func(a *Author) SetData(to func(target interface{}) error) error {
  return to(a)
}

type Authors []Author

func(a *Authors) SetData(to func(target interface{}) error) error {
  return to(a)
}

func printBook(b Book) {
  fmt.Printf("ID:\t%v,\nType:\t%v\nAuthorID:\t%v\nTitle:\t%v\nPublicationDate:\t%v\n", b.ID, b.Type, b.AuthorID, b.Title, b.PublicationDate)
}

func printAuthor(a Author) {
  fmt.Printf("ID:\t%v,\nType:\t%v\nFirstName:\t%v\nLastName:\t%v\n", a.ID, a.Type, a.FirstName, a.LastName)
}

func GetBooks() (books Books){

    res, err := http.Get("http://localhost:8080/books")

  	if err != nil {
  		log.Fatal(err)
  	}

    if res.StatusCode != http.StatusOK {
      return
    }

  	booksJSON, err := ioutil.ReadAll(res.Body)
  	res.Body.Close()

    if err != nil {
  		log.Fatal(err)
  	}
    _, jsonapiErr := jsonapi.Unmarshal(booksJSON, &books)

    if jsonapiErr != nil {
      log.Fatal(jsonapiErr)
    }

    return books
}

func GetAuthors() (authors Authors){

    res, err := http.Get("http://localhost:8080/authors")

  	if err != nil {
  		log.Fatal(err)
  	}

    if res.StatusCode != http.StatusOK {
      return
    }

  	authorsJSON, err := ioutil.ReadAll(res.Body)
  	res.Body.Close()

    if err != nil {
  		log.Fatal(err)
  	}
    _, jsonapiErr := jsonapi.Unmarshal(authorsJSON, &authors)

    if jsonapiErr != nil {
      log.Fatal(jsonapiErr)
    }

    return authors
}

func main() {
  books := GetBooks()

  for _, book := range books {
    printBook(book)
  }

  authors := GetAuthors()

  for _, author := range authors {
    printAuthor(author)
  }
}
```

### See also
* [jsonapi-client-go](https://github.com/pieoneers/jsonapi-client-go)
