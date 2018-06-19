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

type MarshalRelationships interface {
  GetRelationships() map[string]interface{}
}

type UnmarshalRelationships interface {
  SetRelationships(map[string]interface{}) error
}

type Document struct {
  Data   *DocumentData  `json:"data,omitempty"`
  Errors []*ErrorObject `json:"errors,omitempty"`
}

type DocumentData struct {
  One    *ResourceObject
  Many []*ResourceObject
}

type Relationship struct {
  Data *RelationshipData `json:"data"`
}

type RelationshipData struct {
  One    *ResourceObjectIdentifier
  Many []*ResourceObjectIdentifier
}

type ResourceObjectIdentifier struct {
  Type string `json:"type"`
  ID   string `json:"id,omitempty"`
}

type ResourceObject struct {
  ResourceObjectIdentifier
  Attributes    json.RawMessage          `json:"attributes"`
  Relationships map[string]*Relationship `json:"relationships,omitempty"`
}

type ErrorObject struct {
  Title  string             `json:"title,omitempty"`
  Source *ErrorObjectSource `json:"source,omitempty"`
}

type ErrorObjectSource struct {
  Pointer string `json:"pointer,omitempty"`
}

func(d *DocumentData) MarshalJSON() ([]byte, error) {
  if d.One != nil {
    return json.Marshal(d.One)
  }
  return json.Marshal(d.Many)
}

func(d *DocumentData) UnmarshalJSON(data []byte) error {
  if bytes.HasPrefix(data, []byte("{")) {
    return json.Unmarshal(data, &d.One)
  }

  if bytes.HasPrefix(data, []byte("[")) {
    return json.Unmarshal(data, &d.Many)
  }

  return nil
}

func(d *RelationshipData) MarshalJSON() ([]byte, error) {
  if d.One != nil {
    return json.Marshal(d.One)
  }
  return json.Marshal(d.Many)
}

func(d *RelationshipData) UnmarshalJSON(data []byte) error {
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
    doc, err = marshalDocumentStruct(payload)
  case reflect.Slice:
    doc, err = marshalDocumentSlice(payload)
  }
  if err != nil {
    return nil, err
  }

  return json.Marshal(doc)
}

func marshalDocumentStruct(payload interface{}) (*Document, error) {
  one := &ResourceObject{}

  err := marshalResourceObject(payload.(MarshalIdentifier), one)
  if err != nil {
    return nil, err
  }

  doc := &Document{
    Data: &DocumentData{
      One: one,
    },
  }

  return doc, nil
}

func marshalDocumentSlice(payload interface{}) (*Document, error) {
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
      Data: &DocumentData{
        Many: many,
      },
    }
  }

  return doc, nil
}

func marshalResourceObjectIdentifier(i MarshalIdentifier, r *ResourceObjectIdentifier) {
  r.ID = i.GetID()
  r.Type = i.GetType()
}

func marshalResourceObject(i MarshalIdentifier, r *ResourceObject) error {
  attrs, err := json.Marshal(i)
  if err != nil {
    return err
  }

  r.ID = i.GetID()
  r.Type = i.GetType()
  r.Attributes = attrs

  if asserted, ok := i.(MarshalRelationships); ok {
    r.Relationships = make(map[string]*Relationship)

    for key, value := range asserted.GetRelationships() {
      r.Relationships[key] = marshalRelationship(value)
    }
  }

  return nil
}

func marshalRelationship(payload interface{}) *Relationship {
  var rel *Relationship

  switch reflect.TypeOf(payload).Kind() {
  case reflect.Struct:
    rel = marshalRelationshipStruct(payload)
  case reflect.Slice:
    rel = marshalRelationshipSlice(payload)
  }

  return rel
}

func marshalRelationshipStruct(payload interface{}) *Relationship {
  one := &ResourceObjectIdentifier{}

  marshalResourceObjectIdentifier(payload.(MarshalIdentifier), one)

  return &Relationship{
    Data: &RelationshipData{
      One: one,
    },
  }
}

func marshalRelationshipSlice(payload interface{}) *Relationship {
  value := reflect.ValueOf(payload)

  many := []*ResourceObjectIdentifier{}

  for i := 0; i < value.Len(); i++ {
    one := &ResourceObjectIdentifier{}

    marshalResourceObjectIdentifier(value.Index(i).Interface().(MarshalIdentifier), one)

    many = append(many, one)
  }

  return &Relationship{
    Data: &RelationshipData{
      Many: many,
    },
  }
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

func unmarshalResourceObject(res *ResourceObject, ui UnmarshalIdentifier) error {
  var err error

  err = json.Unmarshal(res.Attributes, ui)
  if err != nil {
    return err
  }

  err = ui.SetID(res.ID)
  if err != nil {
    return err
  }

  if ur, ok := ui.(UnmarshalRelationships); ok {
    relationships := make(map[string]interface{})

    for k, v := range res.Relationships {
      one := v.Data.One
      if one != nil {
        relationships[k] = one
      }

      many := v.Data.Many
      if many != nil {
        relationships[k] = many
      }
    }

    err = ur.SetRelationships(relationships)
    if err != nil {
      return err
    }
  }

  return nil
}
