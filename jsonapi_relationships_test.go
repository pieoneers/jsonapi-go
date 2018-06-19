package jsonapi

import (
  "testing"
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

type BookWithAuthor struct {
  ID     string         `json:"-"`
  Title  string         `json:"title"`
  Year   string         `json:"year"`
  Author Author         `json:"-"`
  Errors []*ErrorObject `json:"errors,omitempty"`
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

type BookWithReaders struct {
  ID      string         `json:"-"`
  Title   string         `json:"title"`
  Year    string         `json:"year"`
  Readers []Reader       `json:"-"`
  Errors  []*ErrorObject `json:"errors,omitempty"`
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

func TestMarshalWithRelationships(t *testing.T) {

  t.Run("Resource object with one to one relationship", func(t *testing.T) {
    book := BookWithAuthor{
      ID:    "1",
      Title: "An Introduction to Programming in Go",
      Year:  "2012",
      Author: Author{
        ID:   "1",
        Name: "Caleb Doxsey",
      },
    }

    bytes, _ := Marshal(book)
    actual   := string(bytes)
    expected := `{"data":{"type":"books","id":"1","attributes":{"title":"An Introduction to Programming in Go","year":"2012"},"relationships":{"author":{"data":{"type":"authors","id":"1"}}}}}`

    if actual != expected {
      t.Log("\nActual:\t\t", actual, "\nExpected:\t", expected)
      t.Fail()
    }
  })

  t.Run("Resource object with one to many relationship", func(t *testing.T) {
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

    bytes, _ := Marshal(book)
    actual   := string(bytes)
    expected := `{"data":{"type":"books","id":"1","attributes":{"title":"An Introduction to Programming in Go","year":"2012"},"relationships":{"readers":{"data":[{"type":"people","id":"1"},{"type":"people","id":"2"}]}}}}`

    if actual != expected {
      t.Log("\nActual:\t\t", actual, "\nExpected:\t", expected)
      t.Fail()
    }
  })
}

func TestUnmarshalWithRelationships(t *testing.T) {

  t.Run("Resource object with one to one relationship", func(t *testing.T) {
    payload := []byte(`{"data":{"type":"books","attributes":{"title":"An Introduction to Programming in Go","year":"2012"},"relationships":{"author":{"data":{"type":"authors","id":"1"}}}}}`)

    actual   := BookWithAuthor{}
    expected := BookWithAuthor{
      Title: "An Introduction to Programming in Go",
      Year:  "2012",
      Author: Author{
        ID: "1",
      },
    }

    Unmarshal(payload, &actual)

    if actual.Title != expected.Title || actual.Year != expected.Year || actual.Author.ID != expected.Author.ID {
      t.Log("\nActual:\t\t", actual, "\nExpected:\t", expected)
      t.Fail()
    }
  })

  t.Run("Resource object with one to many relationship", func(t *testing.T) {
    payload := []byte(`{"data":[{"type":"books","id":"1","attributes":{"title":"An Introduction to Programming in Go","year":"2012"},"relationships":{"readers":{"data":[{"type":"people","id":"1"},{"type":"people","id":"2"}]}}}]}`)

    actual   := []*BookWithReaders{}
    expected := []*BookWithReaders{
      {
        ID:    "1",
        Title: "An Introduction to Programming in Go",
        Year:  "2012",
        Readers: []Reader{
          {
            ID: "1",
          },
          {
            ID: "2",
          },
        },
      },
    }

    Unmarshal(payload, &actual)

    for i, _ := range actual {
      if actual[i].ID != expected[i].ID || actual[i].Title != expected[i].Title || actual[i].Year != expected[i].Year {
        t.Log("\nActual:\t\t", *actual[i], "\nExpected:\t", expected[i])
        t.Fail()
      }
    }
  })
}
