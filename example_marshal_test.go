package jsonapi_test

import (
	"fmt"
	"github.com/pieoneers/jsonapi-go"
	"time"
)

type TestMarshalMeta struct {
	Count int `json:"count"`
}

type MarshalBook struct {
	ID              string    `json:"-"`
	AuthorID        string    `json:"-"`
	Title           string    `json:"title"`
	PublicationDate time.Time `json:"publication_date"`
}

func (b MarshalBook) GetID() string {
	return b.ID
}

func (b MarshalBook) GetType() string {
	return "books"
}

func (b MarshalBook) GetData() interface{} {
	return b
}

func (b MarshalBook) GetRelationships() map[string]interface{} {
	relationships := make(map[string]interface{})

	relationships["author"] = jsonapi.ResourceObjectIdentifier{
		ID:   b.AuthorID,
		Type: "authors",
	}

	return relationships
}

func (b MarshalBook) GetIncluded() []interface{} {
	var included []interface{}

	for _, author := range authors {
		if author.ID == b.AuthorID {
			included = append(included, author)
		}
	}

	return included
}

type MarshalBooks []MarshalBook

func (b MarshalBooks) GetData() interface{} {
	return b
}

func (b MarshalBooks) GetMeta() interface{} {
	return TestMarshalMeta{Count: len(books)}
}

func (b MarshalBooks) GetIncluded() []interface{} {
	var included []interface{}

	authorsMap := make(map[string]MarshalAuthor)

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

type MarshalAuthor struct {
	ID        string `json:"-"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (a MarshalAuthor) GetID() string {
	return a.ID
}

func (a MarshalAuthor) GetType() string {
	return "authors"
}

func (a MarshalAuthor) GetData() interface{} {
	return a
}

func (a MarshalAuthor) GetIncluded() []interface{} {
	var included []interface{}

	for _, book := range books {
		if book.AuthorID == a.ID {
			included = append(included, book)
		}
	}

	return included
}

type MarshalAuthors []MarshalAuthor

func (a MarshalAuthors) GetMeta() interface{} {
	return TestMarshalMeta{Count: len(authors)}
}

func (a MarshalAuthors) GetData() interface{} {
	return a
}

func (a MarshalAuthors) GetIncluded() []interface{} {
	var included []interface{}

	booksMap := make(map[string]MarshalBook)

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
	authors MarshalAuthors
	books   MarshalBooks
)

func ExampleMarshal() {
	var publicationDate time.Time

	alan := MarshalAuthor{
		ID:        "1",
		FirstName: "Alan A. A.",
		LastName:  "Donovan",
	}
	authors = append(authors, alan)

	lex := MarshalAuthor{
		ID:        "2",
		FirstName: "Lex",
		LastName:  "Sheehan",
	}
	authors = append(authors, lex)

	william := MarshalAuthor{
		ID:        "3",
		FirstName: "William",
		LastName:  "Kennedy",
	}
	authors = append(authors, william)

	publicationDate, _ = time.Parse(time.RFC3339, "2015-01-01T00:00:00Z")

	book1 := MarshalBook{
		ID:              "1",
		Title:           "Go Programming Language",
		AuthorID:        alan.ID,
		PublicationDate: publicationDate,
	}
	books = append(books, book1)

	publicationDate, _ = time.Parse(time.RFC3339, "2017-11-01T00:00:00Z")

	book2 := MarshalBook{
		ID:              "2",
		Title:           "Learning Functional Programming in Go",
		AuthorID:        lex.ID,
		PublicationDate: publicationDate,
	}
	books = append(books, book2)

	publicationDate, _ = time.Parse(time.RFC3339, "2015-11-01T00:00:00Z")

	book3 := MarshalBook{
		ID:              "3",
		Title:           "Go in Action",
		AuthorID:        william.ID,
		PublicationDate: publicationDate,
	}
	books = append(books, book3)

	bookJSON, _ := jsonapi.Marshal(book1)
	fmt.Printf("book JSON:\n%v\n", string(bookJSON))
	booksJSON, _ := jsonapi.Marshal(books)
	fmt.Printf("books JSON:\n%v\n", string(booksJSON))
	authorJSON, _ := jsonapi.Marshal(alan)
	fmt.Printf("author JSON:\n%v\n", string(authorJSON))
	authorsJSON, _ := jsonapi.Marshal(authors)
	fmt.Printf("authors JSON:\n%v\n", string(authorsJSON))
}
// Output:
// book JSON: {
// 	"data": {
// 		"type": "books",
// 		"id": "1",
// 		"attributes": {
// 			"title": "Go Programming Language",
// 			"publication_date": "2015-01-01T00:00:00Z"
// 		},
// 		"relationships": {
// 			"author": {
// 				"data": {
// 					"type": "authors",
// 					"id": "1"
// 				}
// 			}
// 		}
// 	},
// 	"included": [{
// 		"type": "authors",
// 		"id": "1",
// 		"attributes": {
// 			"first_name": "Alan A. A.",
// 			"last_name": "Donovan"
// 		}
// 	}]
// }
//
// books JSON: {
// 	"data": [{
// 		"type": "books",
// 		"id": "1",
// 		"attributes": {
// 			"title": "Go Programming Language",
// 			"publication_date": "2015-01-01T00:00:00Z"
// 		},
// 		"relationships": {
// 			"author": {
// 				"data": {
// 					"type": "authors",
// 					"id": "1"
// 				}
// 			}
// 		}
// 	}, {
// 		"type": "books",
// 		"id": "2",
// 		"attributes": {
// 			"title": "Learning Functional Programming in Go",
// 			"publication_date": "2017-11-01T00:00:00Z"
// 		},
// 		"relationships": {
// 			"author": {
// 				"data": {
// 					"type": "authors",
// 					"id": "2"
// 				}
// 			}
// 		}
// 	}, {
// 		"type": "books",
// 		"id": "3",
// 		"attributes": {
// 			"title": "Go in Action",
// 			"publication_date": "2015-11-01T00:00:00Z"
// 		},
// 		"relationships": {
// 			"author": {
// 				"data": {
// 					"type": "authors",
// 					"id": "3"
// 				}
// 			}
// 		}
// 	}],
// 	"included": [{
// 		"type": "authors",
// 		"id": "1",
// 		"attributes": {
// 			"first_name": "Alan A. A.",
// 			"last_name": "Donovan"
// 		}
// 	}, {
// 		"type": "authors",
// 		"id": "2",
// 		"attributes": {
// 			"first_name": "Lex",
// 			"last_name": "Sheehan"
// 		}
// 	}, {
// 		"type": "authors",
// 		"id": "3",
// 		"attributes": {
// 			"first_name": "William",
// 			"last_name": "Kennedy"
// 		}
// 	}],
// 	"meta": {
// 		"count": 3
// 	}
// }
//
// author JSON: {
// 	"data": {
// 		"type": "authors",
// 		"id": "1",
// 		"attributes": {
// 			"first_name": "Alan A. A.",
// 			"last_name": "Donovan"
// 		}
// 	},
// 	"included": [{
// 		"type": "books",
// 		"id": "1",
// 		"attributes": {
// 			"title": "Go Programming Language",
// 			"publication_date": "2015-01-01T00:00:00Z"
// 		},
// 		"relationships": {
// 			"author": {
// 				"data": {
// 					"type": "authors",
// 					"id": "1"
// 				}
// 			}
// 		}
// 	}]
// }
//
// authors JSON: {
// 	"data": [{
// 		"type": "authors",
// 		"id": "1",
// 		"attributes": {
// 			"first_name": "Alan A. A.",
// 			"last_name": "Donovan"
// 		}
// 	}, {
// 		"type": "authors",
// 		"id": "2",
// 		"attributes": {
// 			"first_name": "Lex",
// 			"last_name": "Sheehan"
// 		}
// 	}, {
// 		"type": "authors",
// 		"id": "3",
// 		"attributes": {
// 			"first_name": "William",
// 			"last_name": "Kennedy"
// 		}
// 	}],
// 	"included": [{
// 		"type": "books",
// 		"id": "1",
// 		"attributes": {
// 			"title": "Go Programming Language",
// 			"publication_date": "2015-01-01T00:00:00Z"
// 		},
// 		"relationships": {
// 			"author": {
// 				"data": {
// 					"type": "authors",
// 					"id": "1"
// 				}
// 			}
// 		}
// 	}, {
// 		"type": "books",
// 		"id": "2",
// 		"attributes": {
// 			"title": "Learning Functional Programming in Go",
// 			"publication_date": "2017-11-01T00:00:00Z"
// 		},
// 		"relationships": {
// 			"author": {
// 				"data": {
// 					"type": "authors",
// 					"id": "2"
// 				}
// 			}
// 		}
// 	}, {
// 		"type": "books",
// 		"id": "3",
// 		"attributes": {
// 			"title": "Go in Action",
// 			"publication_date": "2015-11-01T00:00:00Z"
// 		},
// 		"relationships": {
// 			"author": {
// 				"data": {
// 					"type": "authors",
// 					"id": "3"
// 				}
// 			}
// 		}
// 	}],
// 	"meta": {
// 		"count": 3
// 	}
// }
