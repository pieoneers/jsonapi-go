package jsonapi

import (
  "bytes"
  "errors"
  "reflect"
  "encoding/json"
)

const ContentType = "application/vnd.api+json"

type MarshalIdentifier interface {
  GetID()   string
  GetType() string
}

type UnmarshalIdentifier interface {
  SetID(string) error
}

type Document struct {
  Data     *Data        `json:"data,omitempty"`
  Errors []*ErrorObject `json:"errors,omitempty"`
}

type Data struct {
  One    *ResourceObject
  Many []*ResourceObject
}

type ResourceObject struct {
  Type       string          `json:"type"`
  ID         string          `json:"id,omitempty"`
  Attributes json.RawMessage `json:"attributes,omitempty"`
}

type ErrorObject struct {
  Title  string             `json:"title,omitempty"`
  Source *ErrorObjectSource `json:"source,omitempty"`
}

type ErrorObjectSource struct {
  Pointer string `json:"pointer,omitempty"`
}

func(d *Data) MarshalJSON() ([]byte, error) {
  if d.One != nil {
    return json.Marshal(d.One)
  }
  return json.Marshal(d.Many)
}

func(d *Data) UnmarshalJSON(data []byte) error {
  if bytes.HasPrefix(data, []byte("{")) {
    return json.Unmarshal(data, &d.One)
  }

  if bytes.HasPrefix(data, []byte("[")) {
    return json.Unmarshal(data, &d.Many)
  }

  return nil
}

func Marshal(payload interface{}) ([]byte, error) {
  var doc *Document
  var err error

  switch reflect.TypeOf(payload).Kind() {
  case reflect.Struct:
    doc, err = marshalStruct(payload)
  case reflect.Slice:
    doc, err = marshalSlice(payload)
  }
  if err != nil {
    return nil, err
  }

  return json.Marshal(doc)
}

func marshalStruct(payload interface{}) (*Document, error) {
  one := &ResourceObject{}

  err := marshalResourceObject(payload.(MarshalIdentifier), one)
  if err != nil {
    return nil, err
  }

  doc := &Document{
    Data: &Data{
      One: one,
    },
  }

  return doc, nil
}

func marshalSlice(payload interface{}) (*Document, error) {
  var doc *Document

  errorObjects, ok := payload.([]*ErrorObject)
  if ok {
    doc = &Document{
      Errors: errorObjects,
    }
  } else {
    value := reflect.ValueOf(payload)

    many := []*ResourceObject{}

    for i := 0; i < value.Len(); i++ {
      one := &ResourceObject{}

      err := marshalResourceObject(value.Index(i).Interface().(MarshalIdentifier), one)
      if err != nil {
        return nil, err
      }

      many = append(many, one)
    }

    doc = &Document{
      Data: &Data{
        Many: many,
      },
    }
  }

  return doc, nil
}

func marshalResourceObject(i MarshalIdentifier, r *ResourceObject) error {
  attrs, err := json.Marshal(i)
  if err != nil {
    return err
  }

  r.ID = i.GetID()
  r.Type = i.GetType()
  r.Attributes = attrs

  return nil
}

func Unmarshal(data []byte, target interface{}) error {
  var err error

  doc := &Document{}

  err = json.Unmarshal(data, doc)
  if err != nil {
    return err
  }

  if doc.Data == nil {
    return errors.New("The root object must have the data key")
  }

  one := doc.Data.One
  if one != nil {
    err = unmarshalOne(one, target)
    if err != nil {
      return err
    }
  }

  many := doc.Data.Many
  if many != nil {
    err = unmarshalMany(many, target)
    if err != nil {
      return err
    }
  }

  return nil
}

func unmarshalOne(one *ResourceObject, target interface{}) error {
  asserted := target.(UnmarshalIdentifier)

  err := unmarshalResourceObject(one, asserted)
  if err != nil {
    return err
  }

  return nil
}

func unmarshalMany(many []*ResourceObject, target interface{}) error {
  typ := reflect.TypeOf(target).Elem().Elem().Elem()
  ptr := reflect.ValueOf(target)
  val := ptr.Elem()

  for _, one := range many {
    new := reflect.New(typ)
    asserted := new.Interface().(UnmarshalIdentifier)

    err := unmarshalResourceObject(one, asserted)
    if err != nil {
      return nil
    }

    val = reflect.Append(val, new)
  }

  ptr.Elem().Set(val)

  return nil
}

func unmarshalResourceObject(r *ResourceObject, i UnmarshalIdentifier) error {
  var err error

  err = json.Unmarshal(r.Attributes, i)
  if err != nil {
    return err
  }

  err = i.SetID(r.ID)
  if err != nil {
    return err
  }

  return nil
}
