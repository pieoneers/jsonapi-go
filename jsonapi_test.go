package jsonapi

import (
  "testing"
)

type Book struct {
  ID     string         `json:"-"`
  Title  string         `json:"title"`
  Author string         `json:"author"`
  Errors []*ErrorObject `json:"errors,omitempty"`
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

func TestMarshal(t *testing.T) {

  t.Run("One resource object", func (t *testing.T) {
    book := Book{
      ID:     "1",
      Title:  "Introducing Go",
      Author: "Caleb Doxsey",
    }

    bytes, _ := Marshal(book)
    actual   := string(bytes)
    expected := `{"data":{"type":"books","id":"1","attributes":{"title":"Introducing Go","author":"Caleb Doxsey"}}}`

    if actual != expected {
      t.Log("\nActual:\t\t", actual, "\nExpected:\t", expected)
      t.Fail()
    }
  })

  t.Run("Many resource objects", func (t *testing.T) {
    books := []Book{
      {
        ID:     "1",
        Title:  "Introducing Go",
        Author: "Caleb Doxsey",
      },
      {
        ID:     "2",
        Title:  "An Introduction to Programming in Go",
        Author: "Caleb Doxsey",
      },
    }

    result, _ := Marshal(books)
    actual    := string(result)
    expected  := `{"data":[{"type":"books","id":"1","attributes":{"title":"Introducing Go","author":"Caleb Doxsey"}},{"type":"books","id":"2","attributes":{"title":"An Introduction to Programming in Go","author":"Caleb Doxsey"}}]}`

    if actual != expected {
      t.Log("\nActual:\t\t", actual, "\nExpected:\t", expected)
      t.Fail()
    }
  })

  t.Run("Error objects", func(t *testing.T) {
    book := Book{
      Errors: []*ErrorObject{
        {
          Title: "is required",
          Source: &ErrorObjectSource {
            Pointer: "/data/attributes/title",
          },
        },
      },
    }

    bytes, _ := Marshal(book.Errors)
    actual   := string(bytes)
    expected := `{"errors":[{"title":"is required","source":{"pointer":"/data/attributes/title"}}]}`

    if actual != expected {
      t.Log("\nActual:\t\t", actual, "\nExpected:\t", expected)
      t.Fail()
    }
  })
}
