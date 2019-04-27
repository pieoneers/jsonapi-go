package jsonapi_test

import (
  "sort"
  . "github.com/pieoneers/jsonapi-go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type Book struct {
  ID    string   `json:"-"`
  Title string   `json:"title"`
  Year  string   `json:"year"`
}

func(b Book) GetID() string {
  return b.ID
}

func(b *Book) SetID(id string) error {
  b.ID = id
  return nil
}

func(b Book) GetType() string {
  return "books"
}

type BookWithMeta struct {
  Book
  Meta BookMeta `json:"-"`
}

func(b BookWithMeta) GetMeta() interface{} {
  return b.Meta
}

type BookMeta struct {
  Sold int `json:"sold,omitempty"`
}

type BookView struct {
  Book Book `json:"-"`
}

func(v BookView) GetData() interface{} {
  return v.Book
}

func(v *BookView) SetData(to func(target interface{}) error) error {
  return to(&v.Book)
}

type BookWithMetaView struct {
  Book BookWithMeta `json:"-"`
}

func(d BookWithMetaView) GetData() interface{} {
  return d.Book
}

type BookWithAuthorView struct {
  Book BookWithAuthor `json:"-"`
}

func(v BookWithAuthorView) GetData() interface{} {
  return v.Book
}

func(v *BookWithAuthorView) SetData(to func(target interface{}) error) error {
  return to(&v.Book)
}

type BookWithAuthorIncludedView struct {
  BookWithAuthorView
}

func(v BookWithAuthorIncludedView) GetIncluded() []interface{} {
  return []interface{}{ v.Book.Author }
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
  if author, ok := relationships["author"]; ok {
    b.Author = Author{ ID: author.(*ResourceObjectIdentifier).ID }
  }

  return nil
}

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

type BookWithReadersView struct {
  Book BookWithReaders `json:"-"`
}

func(v BookWithReadersView) GetData() interface{} {
  return v.Book
}

func(v *BookWithReadersView) SetData(to func(target interface{}) error) error {
  return to(&v.Book)
}

type BookWithReadersIncludedView struct {
  BookWithReadersView
}

func(v BookWithReadersIncludedView) GetIncluded() []interface{} {
  var included []interface{}

  relationships := v.Book.GetRelationships()

  for _, reader := range relationships["readers"].(Readers) {
    included = append(included, reader)
  }

  return included
}

type BookWithReaders struct {
  Book
  Readers Readers `json:"-"`
}

func(b BookWithReaders) GetRelationships() map[string]interface{} {
  return map[string]interface{}{
    "readers": b.Readers,
  }
}

func(b *BookWithReaders) SetRelationships(relationships map[string]interface{}) error {
  relationship := relationships["readers"].([]*ResourceObjectIdentifier)

  for _, reader := range relationship {
    b.Readers = append(b.Readers, Reader{ ID: reader.ID })
  }

  return nil
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

type Readers []Reader

type BooksView struct {
  Books Books     `json:"-"`
  Meta  BooksMeta `json:"-"`
}

func(v BooksView) GetData() interface{} {
  return v.Books
}

func(v *BooksView) SetData(to func(target interface{}) error) error {
  return to(&v.Books)
}

type Books []Book

type BooksWithAuthorsView struct {
  Books []BookWithAuthor `json:"-"`
}

func(v BooksWithAuthorsView) GetData() interface{} {
  return v.Books
}

type BooksWithAuthorsIncludedView struct {
  BooksWithAuthorsView
}

func(v BooksWithAuthorsIncludedView) GetIncluded() []interface{} {
  var included []interface{}

  for _, book := range v.Books {
    relationships := book.GetRelationships()

    included = append(included, relationships["author"])
  }

  return included
}

type BooksWithReadersView struct {
  Books []BookWithReaders `json:"-"`
}

func(v BooksWithReadersView) GetData() interface{} {
  return v.Books
}

type BooksWithReadersIncludedView struct {
  BooksWithReadersView
}

func(v BooksWithReadersIncludedView) GetIncluded() []interface{} {
  var included []interface{}

  filter := make(map[string]Reader)

  for _, book := range v.Books {
    relationships := book.GetRelationships()

    for _, reader := range relationships["readers"].(Readers) {
      filter[reader.ID] = reader
    }
  }

  var readers Readers

  for _, reader := range filter {
    readers = append(readers, reader)
  }

  sort.Slice(readers, func(i, j int) bool {
		return readers[i].ID < readers[j].ID
	})

  for _, reader := range readers {
    included = append(included, reader)
  }

  return included
}

type BooksViewWithMeta struct {
  BooksView
  Meta  BooksMeta `json:"-"`
}

func(v BooksViewWithMeta) GetMeta() interface{} {
  return v.Meta
}

type BooksMeta struct {
  Count int `json:"count"`
}

type OrderView struct {
  Order Order `json:"-"`
}

func(v OrderView) GetData() interface{} {
  return v.Order
}

func(v *OrderView) SetData(to func(target interface{}) error) error {
  return to(&v.Order)
}

type Order struct {
  ID     string `json:"-"`
  Book   Book   `json:"-"`
  Reader Reader `json:"-"`
}

func(o Order) GetType() string {
  return "orders"
}

func(o Order) GetID() string {
  return o.ID
}

func(o *Order) SetID(id string) error {
  o.ID = id
  return nil
}

func(o Order) GetRelationships() map[string]interface{} {
  return map[string]interface{}{
    "book": o.Book,
    "reader": o.Reader,
  }
}

func(o *Order) SetRelationships(relationships map[string]interface{}) error {

  if book, ok := relationships["book"].(*ResourceObjectIdentifier); ok {
    o.Book = Book{ ID: book.ID }
  }

  if reader, ok := relationships["reader"].(*ResourceObjectIdentifier); ok {
    o.Reader = Reader{ ID: reader.ID }
  }

  return nil
}

type ErrorsView struct {
  ValidationErrors []*ErrorObject `json:"-"`
}

func(v ErrorsView) GetErrors() []*ErrorObject {
  return v.ValidationErrors
}

func(v *ErrorsView) SetData(to func(target interface{}) error) error {
  return nil
}

func(v *ErrorsView) SetErrors(errors []*ErrorObject) error {
  v.ValidationErrors = errors
  return nil
}

var _ = Describe("JSONAPI", func() {

  Describe("Marshal", func() {

    It("marshals single resource object", func() {
      view := BookView{
        Book: Book{
          ID:    "1",
          Title: "An Introduction to Programming in Go",
          Year:  "2012",
        },
      }

      result, err := Marshal(view)

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

      Ω(result).Should(MatchJSON(expected))
      Ω(err).ShouldNot(HaveOccurred())
    })

    It("marshals single resource object given as pointer", func() {
      view := &BookView{
        Book: Book{
          ID:    "1",
          Title: "An Introduction to Programming in Go",
          Year:  "2012",
        },
      }

      result, err := Marshal(view)

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

      Ω(result).Should(MatchJSON(expected))
      Ω(err).ShouldNot(HaveOccurred())
    })

		It("marshals single resource object with meta", func() {
      view := BookWithMetaView{
        Book: BookWithMeta{
          Book: Book{
            ID:    "1",
            Title: "An Introduction to Programming in Go",
            Year:  "2012",
          },
          Meta: BookMeta{
            Sold: 10,
          },
        },
      }

      result, err := Marshal(view)

      expected := `
        {
          "data": {
            "type": "books",
            "id": "1",
            "attributes": {
              "title": "An Introduction to Programming in Go",
              "year": "2012"
            },
            "meta": {
              "sold": 10
            }
          }
        }
      `

      Ω(result).Should(MatchJSON(expected))
      Ω(err).ShouldNot(HaveOccurred())
    })

    It("marshals single resource object with to-one relationship", func() {
      view := BookWithAuthorView{
        Book: BookWithAuthor{
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

      result, err := Marshal(view)

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

      Ω(result).Should(MatchJSON(expected))
      Ω(err).ShouldNot(HaveOccurred())
    })

    It("marshals single resource object with empty to-one relationship", func() {
      view := BookWithAuthorView{
        Book: BookWithAuthor{
          Book: Book{
            ID:    "1",
            Title: "An Introduction to Programming in Go",
            Year:  "2012",
          },
        },
      }

      result, err := Marshal(view)

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
                "data": null
              }
            }
          }
        }
      `

      Ω(result).Should(MatchJSON(expected))
      Ω(err).ShouldNot(HaveOccurred())
    })

    It("marshals single resource object with to-one relationship and no attributes", func() {
      view := OrderView{
        Order: Order{
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
        },
      }

      result, err := Marshal(view)

      expected := `
        {
          "data": {
            "type": "orders",
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

      Ω(result).Should(MatchJSON(expected))
      Ω(err).ShouldNot(HaveOccurred())
    })

    It("marshals single resource object with to-one relationship included", func() {
      view := BookWithAuthorIncludedView{
        BookWithAuthorView: BookWithAuthorView{
          Book: BookWithAuthor{
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
      }

      result, err := Marshal(view)

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

      Ω(result).Should(MatchJSON(expected))
      Ω(err).ShouldNot(HaveOccurred())
    })

    It("marshals single resource object with to-many relationship", func() {
      view := BookWithReadersView{
        Book: BookWithReaders{
          Book: Book{
            ID:    "1",
            Title: "An Introduction to Programming in Go",
            Year:  "2012",
          },
          Readers: Readers{
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

      result, err := Marshal(view)

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

      Ω(result).Should(MatchJSON(expected))
      Ω(err).ShouldNot(HaveOccurred())
    })

    It("marshals single resource object with to-many relationship included", func() {
      view := BookWithReadersIncludedView{
        BookWithReadersView: BookWithReadersView{
          Book: BookWithReaders{
            Book: Book{
              ID:    "1",
              Title: "An Introduction to Programming in Go",
              Year:  "2012",
            },
            Readers: Readers{
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
      }

      result, err := Marshal(view)

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

      Ω(result).Should(MatchJSON(expected))
      Ω(err).ShouldNot(HaveOccurred())
    })

    It("marshals single resource object with empty to-many relationship", func() {
      view := BookWithReadersView{
        Book: BookWithReaders{
          Book: Book{
            ID:    "1",
            Title: "An Introduction to Programming in Go",
            Year:  "2012",
          },
        },
      }

      result, err := Marshal(view)

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

      Ω(result).Should(MatchJSON(expected))
      Ω(err).ShouldNot(HaveOccurred())
    })

    It("marshals resource objects collection", func() {
      view := BooksView{
        Books: Books{
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
        },
      }

      result, err := Marshal(view)

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

      Ω(result).Should(MatchJSON(expected))
      Ω(err).ShouldNot(HaveOccurred())
    })

    It("marshals resource objects collection given as pointer", func() {
      view := &BooksView{
        Books: Books{
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
        },
      }

      result, err := Marshal(view)

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

      Ω(result).Should(MatchJSON(expected))
      Ω(err).ShouldNot(HaveOccurred())
    })

		It("marshals resource objects collection with meta", func() {
      view := BooksViewWithMeta{
        BooksView: BooksView{
          Books: Books{
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
          },
        },
        Meta: BooksMeta{
          Count: 2,
        },
      }

      result, err := Marshal(view)

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
            "count": 2
          }
        }
      `

      Ω(result).Should(MatchJSON(expected))
      Ω(err).ShouldNot(HaveOccurred())
    })

    It("marshals resource objects collection with to-one relationships", func() {
      view := BooksWithAuthorsView{
        Books: []BookWithAuthor{
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
        },
      }

      result, err := Marshal(view)

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

      Ω(result).Should(MatchJSON(expected))
      Ω(err).ShouldNot(HaveOccurred())
    })

    It("marshals resource objects collection with to-one relationships included", func() {
      view := BooksWithAuthorsIncludedView{
        BooksWithAuthorsView: BooksWithAuthorsView{
          Books: []BookWithAuthor{
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
          },
        },
      }

      result, err := Marshal(view)

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
            },
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

      Ω(result).Should(MatchJSON(expected))
      Ω(err).ShouldNot(HaveOccurred())
    })

    It("marshals resource objects collection with to-many relationships", func() {
      view := BooksWithReadersView{
        Books: []BookWithReaders{
          {
            Book: Book{
              ID:    "1",
              Title: "An Introduction to Programming in Go",
              Year:  "2012",
            },
            Readers: Readers{
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
            Readers: Readers{
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
        },
      }

      result, err := Marshal(view)

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
                    { "type": "people", "id": "2" },
                    { "type": "people", "id": "1" }
                  ]
                }
              }
            }
          ]
        }
      `

      Ω(result).Should(MatchJSON(expected))
      Ω(err).ShouldNot(HaveOccurred())
    })

    It("marshals resource objects collection with to-many relationships included", func() {
      view := BooksWithReadersIncludedView{
        BooksWithReadersView: BooksWithReadersView{
          Books: []BookWithReaders{
            {
              Book: Book{
                ID:    "1",
                Title: "An Introduction to Programming in Go",
                Year:  "2012",
              },
              Readers: Readers{
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
              Readers: Readers{
                {
                  ID:   "3",
                  Name: "Sasha Petrulevich",
                },
                {
                  ID:   "1",
                  Name: "Fedor Khardikov",
                },
              },
            },
          },
        },
      }

      result, err := Marshal(view)

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
                    { "type": "people", "id": "3" },
                    { "type": "people", "id": "1" }
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
                "name": "Sasha Petrulevich"
              }
            }
          ]
        }
      `

      Ω(result).Should(MatchJSON(expected))
      Ω(err).ShouldNot(HaveOccurred())
    })

    It("marshals empty resource objects collection into empty array", func() {
      view := BooksView{}

      result, err := Marshal(view)

      expected := `
        {
          "data": []
        }
      `

      Ω(result).Should(MatchJSON(expected))
      Ω(err).ShouldNot(HaveOccurred())
    })

    It("marshals error objects collection", func() {
      view := ErrorsView{
        ValidationErrors: []*ErrorObject{
          {
            Title: "is required",
            Source: ErrorObjectSource{
              Pointer: "/data/attributes/title",
            },
          },
          {
            Title: "is required to be in the past",
            Source: ErrorObjectSource{
              Pointer: "/data/attributes/year",
            },
          },
        },
      }

      result, err := Marshal(view)

      expected := `
        {
          "errors": [
            {
              "title": "is required",
              "source": {
                "pointer": "/data/attributes/title"
              }
            },
            {
              "title": "is required to be in the past",
              "source": {
                "pointer": "/data/attributes/year"
              }
            }
          ]
        }
      `

      Ω(result).Should(MatchJSON(expected))
      Ω(err).ShouldNot(HaveOccurred())
    })
  })

  Describe("Unmarshal", func() {

    It("unmarshals single resource object", func() {
      payload := []byte(`
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
      `)

      result   := BookView{}
      expected := BookView{
        Book: Book{
          ID: "1",
          Title: "An Introduction to Programming in Go",
          Year:  "2012",
        },
      }

      _, err := Unmarshal(payload, &result)

      Ω(result).Should(Equal(expected))
      Ω(err).ShouldNot(HaveOccurred())
    })

    It("unmarshals resource object with to-one relationship", func() {
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

      result   := BookWithAuthorView{}
      expected := BookWithAuthorView{
        Book: BookWithAuthor{
          Book: Book{
            Title: "An Introduction to Programming in Go",
            Year:  "2012",
          },
          Author: Author{ ID: "1" },
        },
      }

      _, err := Unmarshal(payload, &result)

      Ω(result).Should(Equal(expected))
      Ω(err).ShouldNot(HaveOccurred())
    })

    It("unmarshals resource object with to-one relationship and no attributes", func() {
      payload := []byte(`
        {
          "data": {
            "type": "orders",
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
      `)

      result   := OrderView{}
      expected := OrderView{
        Order: Order{
          ID: "1",
          Book: Book{ ID: "1" },
          Reader: Reader{ ID: "1" },
        },
      }

      _, err := Unmarshal(payload, &result)


      Ω(result).Should(Equal(expected))
      Ω(err).ShouldNot(HaveOccurred())
    })

    It("unmarshals resource object with to-many relationship", func() {
      payload := []byte(`
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
      `)

      result   := BookWithReadersView{}
      expected := BookWithReadersView{
        Book: BookWithReaders{
          Book: Book{
            ID:    "1",
            Title: "An Introduction to Programming in Go",
            Year:  "2012",
          },
          Readers: Readers{
            { ID: "1" },
            { ID: "2" },
          },
        },
      }

      _, err := Unmarshal(payload, &result)

      Ω(result).Should(Equal(expected))
      Ω(err).ShouldNot(HaveOccurred())
    })

    It("unmarshals resource objects collection", func() {
      payload := []byte(`
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
      `)

      result   := BooksView{}
      expected := BooksView{
        Books: Books{
          {
            ID: "1",
            Title: "An Introduction to Programming in Go",
            Year:  "2012",
          },
          {
            ID: "2",
            Title: "Introducing Go",
            Year:  "2016",
          },
        },
      }

      _, err := Unmarshal(payload, &result)

      Ω(err).ShouldNot(HaveOccurred())
      Ω(result).Should(Equal(expected))
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
              "title": "is required to be in the past",
              "source": {
                "pointer": "/data/attributes/year"
              }
            }
          ]
        }
      `)

      result   := ErrorsView{}
      expected := ErrorsView{
        ValidationErrors: []*ErrorObject{
          {
            Title: "is required",
            Source: ErrorObjectSource{
              Pointer: "/data/attributes/title",
            },
          },
          {
            Title: "is required to be in the past",
            Source: ErrorObjectSource{
              Pointer: "/data/attributes/year",
            },
          },
        },
      }

      _, err := Unmarshal(payload, &result)

      Ω(result).Should(Equal(expected))
      Ω(err).ShouldNot(HaveOccurred())
    })
  })
})
