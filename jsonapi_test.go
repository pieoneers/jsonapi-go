package jsonapi

import (
  "testing"
)

type Book struct {
  ID     string         `json:"-"`
  Title  string         `json:"title"`
  Year   string         `json:"year"`
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

  t.Run("One resource object", func(t *testing.T) {
    book := Book{
      ID:    "1",
      Title: "An Introduction to Programming in Go",
      Year:  "2012",
    }

    bytes, _ := Marshal(book)
    actual   := string(bytes)
    expected := `{"data":{"type":"books","id":"1","attributes":{"title":"An Introduction to Programming in Go","year":"2012"}}}`

    if actual != expected {
      t.Log("\nActual:\t\t", actual, "\nExpected:\t", expected)
      t.Fail()
    }
  })

  t.Run("Many resource objects", func(t *testing.T) {
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

    result, _ := Marshal(books)
    actual    := string(result)
    expected  := `{"data":[{"type":"books","id":"1","attributes":{"title":"An Introduction to Programming in Go","year":"2012"}},{"type":"books","id":"2","attributes":{"title":"Introducing Go","year":"2016"}}]}`

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

func TestUnmarshal(t *testing.T) {

  t.Run("One resource object", func(t *testing.T) {
    payload := []byte(`{"data":{"type":"books","attributes":{"title":"An Introduction to Programming in Go","year":"2012"}}}`)

    actual   := Book{}
    expected := Book{
      Title: "An Introduction to Programming in Go",
      Year:  "2012",
    }

    Unmarshal(payload, &actual)

    if actual.Title != expected.Title || actual.Year != expected.Year {
      t.Log("\nActual:\t\t", actual, "\nExpected:\t", expected)
      t.Fail()
    }
  })

  t.Run("Many resource objects", func(t *testing.T) {
    payload := []byte(`{"data":[{"type":"books","id":"1","attributes":{"title":"An Introduction to Programming in Go","year":"2012"}},{"type":"books","id":"2","attributes":{"title":"Introducing Go","year":"2016"}}]}`)

    actual   := []*Book{}
    expected := []Book{
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

    Unmarshal(payload, &actual)

    for i, _ := range actual {
      if actual[i].ID != expected[i].ID || actual[i].Title != expected[i].Title || actual[i].Year != expected[i].Year {
        t.Log("\nActual:\t\t", *actual[i], "\nExpected:\t", expected[i])
        t.Fail()
      }
    }
  })
}
