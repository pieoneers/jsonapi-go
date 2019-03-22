package jsonapi_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
  . "github.com/Benjamintf1/unmarshalledmatchers"

	. "github.com/pieoneers/jsonapi-go"
)

type Author struct {
  ID   string `json:"-"`
  Name string `json:"name"`
}

func(a Author) GetID() string {
  return a.ID
}

func(a Author) GetType() string {
  return "authors"
}

type Reader struct {
  ID   string `json:"-"`
  Name string `json:"name"`
}

func(r Reader) GetID() string {
  return r.ID
}

func(r Reader) GetType() string {
  return "people"
}

type Book struct {
  ID    string `json:"-"`
  Title string `json:"title"`
  Year  string `json:"year"`
}

func(b Book) GetID() string {
  return b.ID
}

func(b Book) GetType() string {
  return "books"
}

func(b *Book) SetID(id string) error {
  b.ID = id
  return nil
}

type BookWithMeta struct {
	Book
}

type BooksWithMeta []BookWithMeta

type BookMeta struct {
	NumberOfAuthors int `json:"number_of_authors"`
	NumberOfReaders int `json:"number_of_readers"`
	TotalRead 	 		int `json:"total_read"`
}

func(b BookWithMeta) GetMeta() interface{} {
	return BookMeta{
		NumberOfAuthors: 1,
		NumberOfReaders: 2,
		TotalRead: 	 		 3,
	}
}

func(b BooksWithMeta) GetMeta() interface{} {
	return BookMeta{
		NumberOfAuthors: 2,
		NumberOfReaders: 3,
		TotalRead: 	 		 4,
	}
}


type BookWithAuthor struct {
  Book
  Author Author `json:"-"`
}

func(b BookWithAuthor) GetRelationships() map[string]interface{} {
  return map[string]interface{}{
    "author": b.Author,
  }
}

func(b *BookWithAuthor) SetRelationships(relationships map[string]interface{}) error {
  resourceID := relationships["author"].(*ResourceObjectIdentifier)
  b.Author = Author{ ID: resourceID.ID }
  return nil
}

type BookWithAuthorIncluded struct {
  BookWithAuthor
}

func(b BookWithAuthorIncluded) GetIncluded() []interface{} {
  return []interface{}{ b.Author }
}

type BooksWithAuthorIncluded []BookWithAuthor

func(books BooksWithAuthorIncluded) GetIncluded() []interface{} {
  var included []interface{}

  for _, book := range books {
    included = append(included, book.Author)
  }

  return included
}

type BookWithReaders struct {
  Book
  Readers []Reader `json:"-"`
}

func(b BookWithReaders) GetRelationships() map[string]interface{} {
  return map[string]interface{}{
    "readers": b.Readers,
  }
}

func(b *BookWithReaders) SetRelationships(relationships map[string]interface{}) error {
  resourceIDs := relationships["readers"].([]*ResourceObjectIdentifier)

  for _, resourceID := range resourceIDs {
    b.Readers = append(b.Readers, Reader{ ID: resourceID.ID })
  }

  return nil
}

type BookWithReadersIncluded struct {
  BookWithReaders
}

func(b BookWithReadersIncluded) GetIncluded() []interface{} {
  var included []interface{}

  for _, reader := range b.Readers {
    included = append(included, reader)
  }

  return included
}

type BooksWithReadersIncluded []BookWithReaders

func(books BooksWithReadersIncluded) GetIncluded() []interface{} {
  var included []interface{}

  for _, book := range books {
    for _, reader := range book.Readers {
      included = append(included, reader)
    }
  }

  return included
}

type Read struct {
  ID     string `json:"-"`
  Reader Reader `json:"-"`
  Book   Book   `json:"-"`
}

func(r Read) GetType() string {
  return "reads"
}

func(r Read) GetID() string {
  return r.ID
}

func(r Read) GetRelationships() map[string]interface{} {
  return map[string]interface{}{
    "reader": r.Reader,
    "book": r.Book,
  }
}

func(r *Read) SetID(id string) error {
  return nil
}

func(r *Read) SetRelationships(relationships map[string]interface{}) error {
  if resource, ok := relationships["reader"].(*ResourceObjectIdentifier); ok {
    r.Reader = Reader{ ID: resource.ID }
  }

  if resource, ok := relationships["book"].(*ResourceObjectIdentifier); ok {
    r.Book = Book{ ID: resource.ID }
  }

  return nil
}

var _ = Describe("JSONAPI", func() {

  Describe("Marshal", func() {
    It("marshals single resource object", func() {
      book := Book{
        ID:    "1",
        Title: "An Introduction to Programming in Go",
        Year:  "2012",
      }

      bytes, err := Marshal(book)

      actual   := string(bytes)
      expected := `
        {
          "data": {
            "type": "books",
            "id": "1",
            "attributes": {
              "title": "An Introduction to Programming in Go",
              "year": "2012"
            }
          }
        }
      `
      Ω(err).Should(BeNil())
      Ω(actual).Should(MatchJSON(expected))
    })

    It("marshals single resource object given as pointer", func() {
      book := &Book{
        ID:    "1",
        Title: "An Introduction to Programming in Go",
        Year:  "2012",
      }

      bytes, err := Marshal(book)

      actual   := string(bytes)
      expected := `
        {
          "data": {
            "type": "books",
            "id": "1",
            "attributes": {
              "title": "An Introduction to Programming in Go",
              "year": "2012"
            }
          }
        }
      `
      Ω(err).Should(BeNil())
      Ω(actual).Should(MatchJSON(expected))
    })

		It("marshals single resource object with meta information", func() {
      book := BookWithMeta{
				Book: Book{
					ID:    "1",
					Title: "An Introduction to Programming in Go",
					Year:  "2012",
				},
      }

      bytes, err := Marshal(book)

      actual   := string(bytes)
      expected := `
        {
          "data": {
            "type": "books",
            "id": "1",
            "attributes": {
              "title": "An Introduction to Programming in Go",
              "year": "2012"
            }
          },
					"meta": {
						"number_of_authors": 1,
						"number_of_readers": 2,
						"total_read": 3
					}
        }
      `
      Ω(err).Should(BeNil())
      Ω(actual).Should(MatchJSON(expected))
    })

    It("marshals single resource object with one to one relationship", func() {
      book := BookWithAuthor{
        Book: Book{
          ID:    "1",
          Title: "An Introduction to Programming in Go",
          Year:  "2012",
        },
        Author: Author{
          ID:   "1",
          Name: "Caleb Doxsey",
        },
      }

      bytes, err := Marshal(book)

      actual   := string(bytes)
      expected := `
        {
          "data": {
            "type": "books",
            "id": "1",
            "attributes": {
              "title": "An Introduction to Programming in Go",
              "year": "2012"
            },
            "relationships": {
              "author": {
                "data": { "type": "authors", "id": "1" }
              }
            }
          }
        }
      `
      Ω(err).Should(BeNil())
      Ω(actual).Should(MatchJSON(expected))
    })

    It("marshals single resource object with one to one relationship and no attributes", func() {
      book := Read{
        ID: "1",
        Book: Book{
          ID:    "1",
          Title: "An Introduction to Programming in Go",
          Year:  "2012",
        },
        Reader: Reader{
          ID:   "1",
          Name: "Fedor Khardikov",
        },
      }

      bytes, err := Marshal(book)

      actual   := string(bytes)
      expected := `
        {
          "data": {
            "type": "reads",
            "id": "1",
            "relationships": {
              "book": {
                "data": { "type": "books", "id": "1" }
              },
              "reader": {
                "data": { "type": "people", "id": "1" }
              }
            }
          }
        }
      `
      Ω(err).Should(BeNil())
      Ω(actual).Should(MatchJSON(expected))
    })

    It("marshals single resource object with one to one relationship included", func() {
      book := BookWithAuthorIncluded{
        BookWithAuthor: BookWithAuthor{
          Book: Book{
            ID:    "1",
            Title: "An Introduction to Programming in Go",
            Year:  "2012",
          },
          Author: Author{
            ID:   "1",
            Name: "Caleb Doxsey",
          },
        },
      }

      bytes, err := Marshal(book)

      actual   := string(bytes)
      expected := `
        {
          "data": {
            "type": "books",
            "id": "1",
            "attributes": {
              "title": "An Introduction to Programming in Go",
              "year": "2012"
            },
            "relationships": {
              "author": {
                "data": { "type": "authors", "id": "1" }
              }
            }
          },
          "included": [
            {
              "type": "authors",
              "id": "1",
              "attributes": {
                "name": "Caleb Doxsey"
              }
            }
          ]
        }
      `
      Ω(err).Should(BeNil())
      Ω(actual).Should(MatchJSON(expected))
    })

    It("marshals single resource object with one to many relationship", func() {
      book := BookWithReaders{
        Book: Book{
          ID:    "1",
          Title: "An Introduction to Programming in Go",
          Year:  "2012",
        },
        Readers: []Reader{
          {
            ID:   "1",
            Name: "Fedor Khardikov",
          },
          {
            ID:   "2",
            Name: "Andrew Manshin",
          },
        },
      }

      bytes, err := Marshal(book)

      actual   := string(bytes)
      expected := `
        {
          "data": {
            "type": "books",
            "id": "1",
            "attributes": {
              "title": "An Introduction to Programming in Go",
              "year": "2012"
            },
            "relationships": {
              "readers": {
                "data": [
                  { "type": "people", "id": "1" },
                  { "type": "people", "id": "2" }
                ]
              }
            }
          }
        }
      `

      Ω(err).Should(BeNil())
      Ω(actual).Should(MatchJSON(expected))
    })

    It("marshals single resource object with a to-many relationship", func() {
      book := BookWithReaders{
        Book: Book{
          ID:    "1",
          Title: "An Introduction to Programming in Go",
          Year:  "2012",
        },
        Readers: []Reader{
          {
            ID:   "1",
            Name: "Fedor Khardikov",
          },
          {
            ID:   "2",
            Name: "Andrew Manshin",
          },
        },
      }

      bytes, err := Marshal(book)

      actual   := string(bytes)
      expected := `
        {
          "data": {
            "type": "books",
            "id": "1",
            "attributes": {
              "title": "An Introduction to Programming in Go",
              "year": "2012"
            },
            "relationships": {
              "readers": {
                "data": [
                  { "type": "people", "id": "1" },
                  { "type": "people", "id": "2" }
                ]
              }
            }
          }
        }
      `

      Ω(err).Should(BeNil())
      Ω(actual).Should(MatchJSON(expected))
    })

    It("marshals single resource object with an empty to-many relationship included", func() {
      book := BookWithReadersIncluded{
        BookWithReaders: BookWithReaders{
          Book: Book{
            ID:    "1",
            Title: "An Introduction to Programming in Go",
            Year:  "2012",
          },
        },
      }

      bytes, err := Marshal(book)

      actual   := string(bytes)
      expected := `
        {
          "data": {
            "type": "books",
            "id": "1",
            "attributes": {
              "title": "An Introduction to Programming in Go",
              "year": "2012"
            },
            "relationships": {
              "readers": {
                "data": []
              }
            }
          }
        }
      `

      Ω(err).Should(BeNil())
      Ω(actual).Should(MatchJSON(expected))
    })

    It("marshals multiple resource objects", func() {
      books := []Book{
        {
          ID:    "1",
          Title: "An Introduction to Programming in Go",
          Year:  "2012",
        },
        {
          ID:    "2",
          Title: "Introducing Go",
          Year:  "2016",
        },
      }

      bytes, err := Marshal(books)

      actual   := string(bytes)
      expected := `
        {
          "data": [
            {
              "type": "books",
              "id": "1",
              "attributes": {
                "title": "An Introduction to Programming in Go",
                "year": "2012"
              }
            },
            {
              "type": "books",
              "id": "2",
              "attributes": {
                "title": "Introducing Go",
                "year": "2016"
              }
            }
          ]
        }
      `
      Ω(err).Should(BeNil())
      Ω(actual).Should(MatchJSON(expected))
    })

    It("marshals multiple resource objects given as pointers", func() {
      books := &[]*Book{
        {
          ID:    "1",
          Title: "An Introduction to Programming in Go",
          Year:  "2012",
        },
        {
          ID:    "2",
          Title: "Introducing Go",
          Year:  "2016",
        },
      }

      bytes, err := Marshal(books)

      actual   := string(bytes)
      expected := `
        {
          "data": [
            {
              "type": "books",
              "id": "1",
              "attributes": {
                "title": "An Introduction to Programming in Go",
                "year": "2012"
              }
            },
            {
              "type": "books",
              "id": "2",
              "attributes": {
                "title": "Introducing Go",
                "year": "2016"
              }
            }
          ]
        }
      `
      Ω(err).Should(BeNil())
      Ω(actual).Should(MatchJSON(expected))
    })

		It("marshals multiple resource objects with meta", func() {
      books := BooksWithMeta{
        {
					Book: Book{
						ID:    "1",
						Title: "An Introduction to Programming in Go",
						Year:  "2012",
					},
				},
				{
					Book: Book{
					  ID:    "2",
	          Title: "Introducing Go",
	          Year:  "2016",
					},
        },
      }

      bytes, err := Marshal(books)

      actual   := string(bytes)
      expected := `
        {
          "data": [
            {
              "type": "books",
              "id": "1",
              "attributes": {
                "title": "An Introduction to Programming in Go",
                "year": "2012"
              }
            },
            {
              "type": "books",
              "id": "2",
              "attributes": {
                "title": "Introducing Go",
                "year": "2016"
              }
            }
          ],
					"meta": {
						"number_of_authors": 2,
						"number_of_readers": 3,
						"total_read": 4
					}
        }
      `
      Ω(err).Should(BeNil())
      Ω(actual).Should(MatchJSON(expected))
    })


    It("marshals multiple resource objects with one to one relationships", func() {
      books := []BookWithAuthor{
        {
          Book: Book{
            ID:    "1",
            Title: "An Introduction to Programming in Go",
            Year:  "2012",
          },
          Author: Author{
            ID:   "1",
            Name: "Caleb Doxsey",
          },
        },
        {
          Book: Book{
            ID:    "2",
            Title: "Introducing Go",
            Year:  "2016",
          },
          Author: Author{
            ID:   "1",
            Name: "Caleb Doxsey",
          },
        },
      }

      bytes, err := Marshal(books)

      actual   := string(bytes)
      expected := `
        {
          "data": [
            {
              "type": "books",
              "id": "1",
              "attributes": {
                "title": "An Introduction to Programming in Go",
                "year": "2012"
              },
              "relationships": {
                "author": {
                  "data": { "type": "authors", "id": "1" }
                }
              }
            },
            {
              "type": "books",
              "id": "2",
              "attributes": {
                "title": "Introducing Go",
                "year": "2016"
              },
              "relationships": {
                "author": {
                  "data": { "type": "authors", "id": "1" }
                }
              }
            }
          ]
        }
      `
      Ω(err).Should(BeNil())
      Ω(actual).Should(MatchJSON(expected))
    })

    It("marshals multiple resource objects with one to one relationships included", func() {
      books := []BookWithAuthorIncluded{
        {
          BookWithAuthor: BookWithAuthor{
            Book: Book{
              ID:    "1",
              Title: "An Introduction to Programming in Go",
              Year:  "2012",
            },
            Author: Author{
              ID:   "1",
              Name: "Caleb Doxsey",
            },
          },
        },
        {
          BookWithAuthor: BookWithAuthor{
            Book: Book{
              ID:    "2",
              Title: "Introducing Go",
              Year:  "2016",
            },
            Author: Author{
              ID:   "1",
              Name: "Caleb Doxsey",
            },
          },
        },
      }

      bytes, err := Marshal(books)

      actual   := string(bytes)
      expected := `
        {
          "data": [
            {
              "type": "books",
              "id": "1",
              "attributes": {
                "title": "An Introduction to Programming in Go",
                "year": "2012"
              },
              "relationships": {
                "author": {
                  "data": { "type": "authors", "id": "1" }
                }
              }
            },
            {
              "type": "books",
              "id": "2",
              "attributes": {
                "title": "Introducing Go",
                "year": "2016"
              },
              "relationships": {
                "author": {
                  "data": { "type": "authors", "id": "1" }
                }
              }
            }
          ],
          "included": [
            {
              "type": "authors",
              "id": "1",
              "attributes": {
                "name": "Caleb Doxsey"
              }
            }
          ]
        }
      `
      Ω(err).Should(BeNil())
      Ω(actual).Should(MatchJSON(expected))
    })

    It("marshals resource objects collection with one to one relationships included", func() {
      books := BooksWithAuthorIncluded{
        {
          Book: Book{
            ID:    "1",
            Title: "An Introduction to Programming in Go",
            Year:  "2012",
          },
          Author: Author{
            ID:   "1",
            Name: "Caleb Doxsey",
          },
        },
        {
          Book: Book{
            ID:    "2",
            Title: "Introducing Go",
            Year:  "2016",
          },
          Author: Author{
            ID:   "1",
            Name: "Caleb Doxsey",
          },
        },
      }

      bytes, err := Marshal(books)

      actual   := string(bytes)
      expected := `
        {
          "data": [
            {
              "type": "books",
              "id": "1",
              "attributes": {
                "title": "An Introduction to Programming in Go",
                "year": "2012"
              },
              "relationships": {
                "author": {
                  "data": { "type": "authors", "id": "1" }
                }
              }
            },
            {
              "type": "books",
              "id": "2",
              "attributes": {
                "title": "Introducing Go",
                "year": "2016"
              },
              "relationships": {
                "author": {
                  "data": { "type": "authors", "id": "1" }
                }
              }
            }
          ],
          "included": [
            {
              "type": "authors",
              "id": "1",
              "attributes": {
                "name": "Caleb Doxsey"
              }
            }
          ]
        }
      `
      Ω(err).Should(BeNil())
      Ω(actual).Should(MatchJSON(expected))
    })

    It("marshals multiple resource objects with one to many relationships", func() {
      books := []BookWithReaders{
        {
          Book: Book{
            ID:    "1",
            Title: "An Introduction to Programming in Go",
            Year:  "2012",
          },
          Readers: []Reader{
            {
              ID:   "1",
              Name: "Fedor Khardikov",
            },
            {
              ID:   "2",
              Name: "Andrew Manshin",
            },
          },
        },
        {
          Book: Book{
            ID:    "2",
            Title: "Introducing Go",
            Year:  "2016",
          },
          Readers: []Reader{
            {
              ID:   "2",
              Name: "Andrew Manshin",
            },
            {
              ID:   "1",
              Name: "Fedor Khardikov",
            },
          },
        },
      }

      bytes, err := Marshal(books)

      actual   := string(bytes)
      expected := `
        {
          "data": [
            {
              "type": "books",
              "id": "1",
              "attributes": {
                "title": "An Introduction to Programming in Go",
                "year": "2012"
              },
              "relationships": {
                "readers": {
                  "data": [
                    { "type": "people", "id": "1" },
                    { "type": "people", "id": "2" }
                  ]
                }
              }
            },
            {
              "type": "books",
              "id": "2",
              "attributes": {
                "title": "Introducing Go",
                "year": "2016"
              },
              "relationships": {
                "readers": {
                  "data": [
                    { "type": "people", "id": "1" },
                    { "type": "people", "id": "2" }
                  ]
                }
              }
            }
          ]
        }
      `

      Ω(err).Should(BeNil())
      Ω(actual).Should(MatchUnorderedJSON(expected))
    })

    It("marshals multiple resource objects with one to many relationships included", func() {
      books := []BookWithReadersIncluded{
        {
          BookWithReaders {
            Book: Book{
              ID:    "1",
              Title: "An Introduction to Programming in Go",
              Year:  "2012",
            },
            Readers: []Reader{
              {
                ID:   "1",
                Name: "Fedor Khardikov",
              },
              {
                ID:   "2",
                Name: "Andrew Manshin",
              },
            },
          },
        },
        {
          BookWithReaders {
            Book: Book{
              ID:    "2",
              Title: "Introducing Go",
              Year:  "2016",
            },
            Readers: []Reader{
              {
                ID:   "3",
                Name: "Shane McCallum",
              },
              {
                ID:   "1",
                Name: "Fedor Khardikov",
              },
            },
          },
        },
      }

      bytes, err := Marshal(books)

      actual   := string(bytes)
      expected := `
        {
          "data": [
            {
              "type": "books",
              "id": "1",
              "attributes": {
                "title": "An Introduction to Programming in Go",
                "year": "2012"
              },
              "relationships": {
                "readers": {
                  "data": [
                    { "type": "people", "id": "1" },
                    { "type": "people", "id": "2" }
                  ]
                }
              }
            },
            {
              "type": "books",
              "id": "2",
              "attributes": {
                "title": "Introducing Go",
                "year": "2016"
              },
              "relationships": {
                "readers": {
                  "data": [
                    { "type": "people", "id": "1" },
                    { "type": "people", "id": "3" }
                  ]
                }
              }
            }
          ],
          "included": [
            {
              "type": "people",
              "id": "1",
              "attributes": {
                "name": "Fedor Khardikov"
              }
            },
            {
              "type": "people",
              "id": "2",
              "attributes": {
                "name": "Andrew Manshin"
              }
            },
            {
              "type": "people",
              "id": "3",
              "attributes": {
                "name": "Shane McCallum"
              }
            }
          ]
        }
      `

      Ω(err).Should(BeNil())
      Ω(actual).Should(MatchUnorderedJSON(expected))
    })

    It("marshals resource objects collection with one to many relationships included", func() {
      books := BooksWithReadersIncluded{
        {
          Book: Book{
            ID:    "1",
            Title: "An Introduction to Programming in Go",
            Year:  "2012",
          },
          Readers: []Reader{
            {
              ID:   "1",
              Name: "Fedor Khardikov",
            },
            {
              ID:   "2",
              Name: "Andrew Manshin",
            },
          },
        },
        {
          Book: Book{
            ID:    "2",
            Title: "Introducing Go",
            Year:  "2016",
          },
          Readers: []Reader{
            {
              ID:   "3",
              Name: "Shane McCallum",
            },
            {
              ID:   "1",
              Name: "Fedor Khardikov",
            },
          },
        },
      }

      bytes, err := Marshal(books)

      actual   := string(bytes)
      expected := `
        {
          "data": [
            {
              "type": "books",
              "id": "1",
              "attributes": {
                "title": "An Introduction to Programming in Go",
                "year": "2012"
              },
              "relationships": {
                "readers": {
                  "data": [
                    { "type": "people", "id": "1" },
                    { "type": "people", "id": "2" }
                  ]
                }
              }
            },
            {
              "type": "books",
              "id": "2",
              "attributes": {
                "title": "Introducing Go",
                "year": "2016"
              },
              "relationships": {
                "readers": {
                  "data": [
                    { "type": "people", "id": "1" },
                    { "type": "people", "id": "3" }
                  ]
                }
              }
            }
          ],
          "included": [
            {
              "type": "people",
              "id": "1",
              "attributes": {
                "name": "Fedor Khardikov"
              }
            },
            {
              "type": "people",
              "id": "2",
              "attributes": {
                "name": "Andrew Manshin"
              }
            },
            {
              "type": "people",
              "id": "3",
              "attributes": {
                "name": "Shane McCallum"
              }
            }
          ]
        }
      `

      Ω(err).Should(BeNil())
      Ω(actual).Should(MatchUnorderedJSON(expected))
    })

    It("marshals empty collection into empty array", func() {
      books := []Book{}

      bytes, err := Marshal(books)

      actual   := string(bytes)
      expected := `
        {
          "data": []
        }
      `
      Ω(err).Should(BeNil())
      Ω(actual).Should(MatchJSON(expected))
    })
  })

  Describe("Unmarshal", func() {
    It("unmarshals single resource object", func() {
      payload := []byte(`
        {
          "data": {
            "type": "books",
            "attributes": {
              "title": "An Introduction to Programming in Go",
              "year": "2012"
            }
          }
        }
      `)

      actual   := Book{}
      expected := Book{
        Title: "An Introduction to Programming in Go",
        Year:  "2012",
      }

      _, err := Unmarshal(payload, &actual)

      Ω(err).Should(BeNil())
      Ω(actual).Should(Equal(expected))
    })

    It("unmarshals single resource object with one to one relationship", func() {
      payload := []byte(`
        {
          "data": {
            "type": "books",
            "attributes": {
              "title": "An Introduction to Programming in Go",
              "year": "2012"
            },
            "relationships": {
              "author": {
                "data": { "type": "authors", "id": "1" }
              }
            }
          }
        }
      `)

      actual   := BookWithAuthor{}
      expected := BookWithAuthor{
        Book: Book{
          Title: "An Introduction to Programming in Go",
          Year:  "2012",
        },
        Author: Author{ ID: "1", },
      }

      _, err := Unmarshal(payload, &actual)

      Ω(err).Should(BeNil())
      Ω(actual).Should(Equal(expected))
    })

    It("unmarshals single resource object with one to one relationship and no attributes", func() {
      payload := []byte(`
        {
          "data": {
            "type": "reads",
            "relationships": {
              "reader": {
                "data": { "type": "people", "id": "1" }
              },
              "book": {
                "data": { "type": "books", "id": "1" }
              }
            }
          }
        }
      `)

      actual   := Read{}
      expected := Read{
        Reader: Reader{
          ID: "1",
        },
        Book: Book{
          ID: "1",
        },
      }

      _, err := Unmarshal(payload, &actual)

      Ω(err).Should(BeNil())
      Ω(actual).Should(Equal(expected))
    })

    It("unmarshals single resource object with one to many relationship", func() {
      payload := []byte(`
        {
          "data": [
            {
              "type": "books",
              "id": "1",
              "attributes": {
                "title": "An Introduction to Programming in Go",
                "year": "2012"
              },
              "relationships": {
                "readers": {
                  "data": [
                    { "type": "people", "id": "1" },
                    { "type": "people", "id": "2" }
                  ]
                }
              }
            }
          ]
        }
      `)

      actual   := []*BookWithReaders{}
      expected := []*BookWithReaders{
        {
          Book: Book{
            ID:    "1",
            Title: "An Introduction to Programming in Go",
            Year:  "2012",
          },
          Readers: []Reader{
            { ID: "1" },
            { ID: "2" },
          },
        },
      }

      _, err := Unmarshal(payload, &actual)

      Ω(err).Should(BeNil())
      Ω(actual).Should(Equal(expected))
    })

    It("unmarshals multiple resource objects", func() {
      payload := []byte(`
        { "data":
          [
            {
              "type": "books",
              "id": "1",
              "attributes": {
                "title": "An Introduction to Programming in Go",
                "year": "2012"
              }
            },
            { "type": "books",
              "id": "2",
              "attributes": {
                "title": "Introducing Go",
                "year": "2016"
              }
            }
          ]
        }
      `)

      actual   := []*Book{}
      expected := []*Book{
        {
          ID:    "1",
          Title: "An Introduction to Programming in Go",
          Year:  "2012",
        },
        {
          ID:    "2",
          Title: "Introducing Go",
          Year:  "2016",
        },
      }

      _, err := Unmarshal(payload, &actual)

      Ω(err).Should(BeNil())
      Ω(actual).Should(Equal(expected))
    })

    It("unmarshals error objects", func() {
      payload := []byte(`
        {
          "errors": [
            {
              "title": "is required",
              "source": {
                "pointer": "/data/attributes/title"
              }
            },
            {
              "title": "is required",
              "source": {
                "pointer": "/data/attributes/year"
              }
            }
          ]
        }
      `)

      doc, err := Unmarshal(payload, &Book{})

      actual := doc.Errors
      expected := []*ErrorObject{
        {
          Title: "is required",
          Source: ErrorObjectSource{
            Pointer: "/data/attributes/title",
          },
        },
        {
          Title: "is required",
          Source: ErrorObjectSource{
            Pointer: "/data/attributes/year",
          },
        },
      }

      Ω(err).ShouldNot(HaveOccurred())
      Ω(actual).Should(Equal(expected))
    })
  })
})
