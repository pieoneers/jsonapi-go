package jsonapi_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/pieoneers/jsonapi-go"
)

type Author struct {
  ID   string `json:"-"`
  Name string `json:"name"`
}

func (a Author) GetID() string {
  return a.ID
}

func (a Author) GetType() string {
  return "authors"
}

type Reader struct {
  ID   string `json:"-"`
  Name string `json:"name"`
}

func (r Reader) GetID() string {
  return r.ID
}

func (r Reader) GetType() string {
  return "people"
}

type Book struct {
  ID    string `json:"-"`
  Title string `json:"title"`
  Year  string `json:"year"`
}

func (b Book) GetID() string {
  return b.ID
}

func (b Book) GetType() string {
  return "books"
}

func (b *Book) SetID(id string) error {
  b.ID = id
  return nil
}

type BookWithAuthor struct {
  Book
  Author Author `json:"-"`
}

func (b BookWithAuthor) GetRelationships() map[string]interface{} {
  return map[string]interface{}{
    "author": b.Author,
  }
}

func (b *BookWithAuthor) SetRelationships(relationships map[string]interface{}) error {
  resourceID := relationships["author"].(*ResourceObjectIdentifier)
  b.Author = Author{ ID: resourceID.ID }
  return nil
}

type BookWithAuthorIncluded struct {
  BookWithAuthor
}

func (b BookWithAuthorIncluded) GetIncluded() []interface{} {
  return []interface{}{ b.Author }
}

type BookWithReaders struct {
  Book
  Readers []Reader `json:"-"`
}

func (b BookWithReaders) GetRelationships() map[string]interface{} {
  return map[string]interface{}{
    "readers": b.Readers,
  }
}

func (b *BookWithReaders) SetRelationships(relationships map[string]interface{}) error {
  resourceIDs := relationships["readers"].([]*ResourceObjectIdentifier)

  for _, resourceID := range resourceIDs {
    b.Readers = append(b.Readers, Reader{ ID: resourceID.ID })
  }

  return nil
}

type BookWithReadersIncluded struct {
  BookWithReaders
}

func (b BookWithReadersIncluded) GetIncluded() []interface{} {
  var included []interface{}

  for _, reader := range b.Readers {
    included = append(included, reader)
  }

  return included
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

    It("marshals single resource object with one to many relationship included", func() {
      book := BookWithReadersIncluded{
        BookWithReaders: BookWithReaders{
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
          },
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
            }
          ]
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
      Ω(actual).Should(MatchJSON(expected))
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
      Ω(actual).Should(MatchJSON(expected))
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

      err := Unmarshal(payload, &actual)

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

      err := Unmarshal(payload, &actual)

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

      err := Unmarshal(payload, &actual)

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

      err := Unmarshal(payload, &actual)

      Ω(err).Should(BeNil())
      Ω(actual).Should(Equal(expected))
    })
  })
})
