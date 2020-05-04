package jsonapi_test

import (
	"fmt"
	"github.com/pieoneers/jsonapi-go"
	"time"
)

var bookJSON = []byte(`
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
`)

var authorJSON = []byte(`
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
`)

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

var authorsJSON = []byte(`
{
	"data": [
    {
  		"type": "authors",
  		"id": "1",
  		"attributes": {
  			"first_name": "Alan A. A.",
  			"last_name": "Donovan"
  		}
  	},
    {
  		"type": "authors",
  		"id": "2",
  		"attributes": {
  			"first_name": "Lex",
  			"last_name": "Sheehan"
  		}
  	},
    {
  		"type": "authors",
  		"id": "3",
  		"attributes": {
  			"first_name": "William",
  			"last_name": "Kennedy"
  		}
  	}
  ]
}
`)

type UnmarshalBook struct {
	ID              string    `json:"-"`
	Type            string    `json:"-"`
	AuthorID        string    `json:"-"`
	Title           string    `json:"title"`
	PublicationDate time.Time `json:"publication_date"`
}

func (b *UnmarshalBook) SetID(id string) error {
	b.ID = id
	return nil
}

func (b *UnmarshalBook) SetType(t string) error {
	b.Type = t
	return nil
}

func (b *UnmarshalBook) SetData(to func(target interface{}) error) error {
	return to(b)
}

func (b *UnmarshalBook) SetRelationships(relationships map[string]interface{}) error {
	if relationship, ok := relationships["author"]; ok {
		b.AuthorID = relationship.(*jsonapi.ResourceObjectIdentifier).ID
	}

	return nil
}

type UnmarshalBooks []UnmarshalBook

func (b *UnmarshalBooks) SetData(to func(target interface{}) error) error {
	return to(b)
}

type UnmarshalAuthor struct {
	ID        string `json:"-"`
	Type      string `json:"-"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (a *UnmarshalAuthor) SetID(id string) error {
	a.ID = id
	return nil
}

func (a *UnmarshalAuthor) SetType(t string) error {
	a.Type = t
	return nil
}

func (a *UnmarshalAuthor) SetData(to func(target interface{}) error) error {
	return to(a)
}

type UnmarshalAuthors []UnmarshalAuthor

func (a *UnmarshalAuthors) SetData(to func(target interface{}) error) error {
	return to(a)
}

func printBook(b UnmarshalBook) {
	fmt.Printf("ID:\t%v,\nType:\t%v\nAuthorID:\t%v\nTitle:\t%v\nPublicationDate:\t%v\n", b.ID, b.Type, b.AuthorID, b.Title, b.PublicationDate)
}

func printAuthor(a UnmarshalAuthor) {
	fmt.Printf("ID:\t%v,\nType:\t%v\nFirstName:\t%v\nLastName:\t%v\n", a.ID, a.Type, a.FirstName, a.LastName)
}

func ExampleUnmarshal() {
	book := UnmarshalBook{}
	books := UnmarshalBooks{}
	author := UnmarshalAuthor{}
	authors := UnmarshalAuthors{}

	fmt.Printf("Book\n")
	jsonapi.Unmarshal(bookJSON, &book)
	printBook(book)

	fmt.Printf("\nBooks\n")
	jsonapi.Unmarshal(booksJSON, &books)
	for _, b := range books {
		printBook(b)
	}

	fmt.Printf("\nAuthor\n")
	jsonapi.Unmarshal(authorJSON, &author)
	printAuthor(author)

	fmt.Printf("\nAuthors\n")
	jsonapi.Unmarshal(authorsJSON, &authors)
	for _, a := range authors {
		printAuthor(a)
	}
}
// Output:
// Book
// ID:	1,
// Type:	books
// AuthorID:	1
// Title:	Go Programming Language
// PublicationDate:	2015-01-01 00:00:00 +0000 UTC
//
// Books
// ID:	1,
// Type:	books
// AuthorID:	1
// Title:	Go Programming Language
// PublicationDate:	2015-01-01 00:00:00 +0000 UTC
// ID:	2,
// Type:	books
// AuthorID:	2
// Title:	Learning Functional Programming in Go
// PublicationDate:	2017-11-01 00:00:00 +0000 UTC
// ID:	3,
// Type:	books
// AuthorID:	3
// Title:	Go in Action
// PublicationDate:	2015-11-01 00:00:00 +0000 UTC
//
// Author
// ID:	1,
// Type:	authors
// FirstName:	Alan A. A.
// LastName:	Donovan
//
// Authors
// ID:	1,
// Type:	authors
// FirstName:	Alan A. A.
// LastName:	Donovan
// ID:	2,
// Type:	authors
// FirstName:	Lex
// LastName:	Sheehan
// ID:	3,
// Type:	authors
// FirstName:	William
// LastName:	Kennedy
