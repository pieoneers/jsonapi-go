package jsonapi_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "jsonapi-go"
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

func (a Reader) GetID() string {
  return a.ID
}

func (a Reader) GetType() string {
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
  ID     string `json:"-"`
  Title  string `json:"title"`
  Year   string `json:"year"`
  Author Author `json:"-"`
}

func (b BookWithAuthor) GetID() string {
  return b.ID
}

func (b BookWithAuthor) GetType() string {
  return "books"
}

func (b BookWithAuthor) GetRelationships() map[string]interface{} {
  return map[string]interface{}{
    "author": b.Author,
  }
}

func (b *BookWithAuthor) SetID(id string) error {
  b.ID = id
  return nil
}

func (b *BookWithAuthor) SetRelationships(relationships map[string]interface{}) error {
  resourceID := relationships["author"].(*ResourceObjectIdentifier)
  b.Author = Author{ ID: resourceID.ID }
  return nil
}

type BookWithReaders struct {
  ID      string   `json:"-"`
  Title   string   `json:"title"`
  Year    string   `json:"year"`
  Readers []Reader `json:"-"`
}

func (b BookWithReaders) GetID() string {
  return b.ID
}

func (b BookWithReaders) GetType() string {
  return "books"
}

func (b BookWithReaders) GetRelationships() map[string]interface{} {
  return map[string]interface{}{
    "readers": b.Readers,
  }
}

func (b *BookWithReaders) SetID(id string) error {
  b.ID = id
  return nil
}

func (b *BookWithReaders) SetRelationships(relationships map[string]interface{}) error {
  resourceIDs := relationships["readers"].([]*ResourceObjectIdentifier)

  for _, resourceID := range resourceIDs {
    b.Readers = append(b.Readers, Reader{ ID: resourceID.ID })
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

    It("marshals single resource object with one to one relationship", func() {
      book := BookWithAuthor{
        ID:    "1",
        Title: "An Introduction to Programming in Go",
        Year:  "2012",
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

    It("marshals single resource object with one to many relationship", func() {
      book := BookWithReaders{
        ID:    "1",
        Title: "An Introduction to Programming in Go",
        Year:  "2012",
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
        Title: "An Introduction to Programming in Go",
        Year:  "2012",
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
          ID:    "1",
          Title: "An Introduction to Programming in Go",
          Year:  "2012",
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
