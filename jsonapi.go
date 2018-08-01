package jsonapi

import (
  "bytes"
  "errors"
  "reflect"
  "encoding/json"
)

const ContentType = "application/vnd.api+json"

type MarshalResourceIdentifier interface {
  GetID()   string
  GetType() string
}

type UnmarshalResourceIdentifier interface {
  SetID(string) error
}

type MarshalRelationships interface {
  GetRelationships() map[string]interface{}
}

type UnmarshalRelationships interface {
  SetRelationships(map[string]interface{}) error
}

type MarshalIncluded interface {
  GetIncluded() []interface{}
}

type document struct {
  Data     *documentData     `json:"data,omitempty"`
  Errors   []*ErrorObject    `json:"errors,omitempty"`
  Included []*ResourceObject `json:"included,omitempty"`
}

type documentData struct {
  One    *ResourceObject
  Many []*ResourceObject
}

type relationship struct {
  Data *relationshipData `json:"data"`
}

type relationshipData struct {
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
  Relationships map[string]*relationship `json:"relationships,omitempty"`
}

type ErrorObject struct {
  Title  string             `json:"title,omitempty"`
  Source *ErrorObjectSource `json:"source,omitempty"`
}

type ErrorObjectSource struct {
  Pointer string `json:"pointer,omitempty"`
}

func(d *documentData) MarshalJSON() ([]byte, error) {
  if d.One != nil {
    return json.Marshal(d.One)
  }
  return json.Marshal(d.Many)
}

func(d *documentData) UnmarshalJSON(data []byte) error {
  if bytes.HasPrefix(data, []byte("{")) {
    return json.Unmarshal(data, &d.One)
  }

  if bytes.HasPrefix(data, []byte("[")) {
    return json.Unmarshal(data, &d.Many)
  }

  return nil
}

func(d *relationshipData) MarshalJSON() ([]byte, error) {
  if d.One != nil {
    return json.Marshal(d.One)
  }
  return json.Marshal(d.Many)
}

func(d *relationshipData) UnmarshalJSON(data []byte) error {
  if bytes.HasPrefix(data, []byte("{")) {
    return json.Unmarshal(data, &d.One)
  }

  if bytes.HasPrefix(data, []byte("[")) {
    return json.Unmarshal(data, &d.Many)
  }

  return nil
}

func Marshal(payload interface{}) ([]byte, error) {
  var (
    doc *document
    err error
  )

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

func marshalDocumentStruct(payload interface{}) (*document, error) {
  one, included, err := marshalResourceObject(payload.(MarshalResourceIdentifier))
  if err != nil {
    return nil, err
  }

  doc := &document{
    Data: &documentData{
      One: &one,
    },
    Included: included,
  }

  return doc, nil
}

func marshalDocumentSlice(payload interface{}) (*document, error) {
  var doc *document

  if errorObjects, ok := payload.([]*ErrorObject); ok {
    doc = &document{
      Errors: errorObjects,
    }
  } else {
    var included []*ResourceObject

    many := []*ResourceObject{}

    value := reflect.ValueOf(payload)

    for i := 0; i < value.Len(); i++ {
      one, inc, err := marshalResourceObject(value.Index(i).Interface().(MarshalResourceIdentifier))
      if err != nil {
        return nil, err
      }

      many = append(many, &one)

      for _, i := range inc {
        included = append(included, i)
      }
    }

    doc = &document{
      Data: &documentData{
        Many: many,
      },
      Included: included,
    }
  }

  return doc, nil
}

func marshalResourceObjectIdentifier(mri MarshalResourceIdentifier) ResourceObjectIdentifier {
  return ResourceObjectIdentifier{ ID: mri.GetID(), Type: mri.GetType() }
}

func marshalResourceObject(mri MarshalResourceIdentifier) (ResourceObject, []*ResourceObject, error) {
  var included []*ResourceObject

  one := ResourceObject{
    ResourceObjectIdentifier: marshalResourceObjectIdentifier(mri),
  }

  attributes, err := json.Marshal(mri)
  if err != nil {
    return one, included, err
  }

  one.Attributes = attributes

  if mr, ok := mri.(MarshalRelationships); ok {
    one.Relationships = marshalRelationships(mr)

    if mi, ok := mri.(MarshalIncluded); ok {
      included, err = marshalIncluded(mi)
      if err != nil {
        return one, included, err
      }
    }
  }

  return one, included, nil
}

func marshalRelationships(mr MarshalRelationships) map[string]*relationship {
  relationships := map[string]*relationship{}

  for key, value := range mr.GetRelationships() {
    relationships[key] = marshalRelationship(value)
  }

  return relationships
}

func marshalRelationship(payload interface{}) *relationship {
  var rel *relationship

  switch reflect.TypeOf(payload).Kind() {
  case reflect.Struct:
    rel = marshalRelationshipStruct(payload)
  case reflect.Slice:
    rel = marshalRelationshipSlice(payload)
  }

  return rel
}

func marshalRelationshipStruct(payload interface{}) *relationship {
  one := marshalResourceObjectIdentifier(payload.(MarshalResourceIdentifier))

  return &relationship{
    Data: &relationshipData{
      One: &one,
    },
  }
}

func marshalRelationshipSlice(payload interface{}) *relationship {
  many := []*ResourceObjectIdentifier{}

  value := reflect.ValueOf(payload)

  for i := 0; i < value.Len(); i++ {
    one := marshalResourceObjectIdentifier(value.Index(i).Interface().(MarshalResourceIdentifier))

    many = append(many, &one)
  }

  return &relationship{
    Data: &relationshipData{
      Many: many,
    },
  }
}

func marshalIncluded(mi MarshalIncluded) ([]*ResourceObject, error) {
  var included []*ResourceObject

  for _, value := range mi.GetIncluded() {
    inc, _, err := marshalResourceObject(value.(MarshalResourceIdentifier))
    if err != nil {
      return included, err
    }

    included = append(included, &inc)
  }

  return included, nil
}

func Unmarshal(data []byte, target interface{}) error {
  var err error

  doc := &document{}

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
  asserted := target.(UnmarshalResourceIdentifier)

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
    asserted := new.Interface().(UnmarshalResourceIdentifier)

    err := unmarshalResourceObject(one, asserted)
    if err != nil {
      return nil
    }

    val = reflect.Append(val, new)
  }

  ptr.Elem().Set(val)

  return nil
}

func unmarshalResourceObject(res *ResourceObject, ui UnmarshalResourceIdentifier) error {
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
    err = unmarshalRelationships(res, ur)
    if err != nil {
      return err
    }
  }

  return nil
}

func unmarshalRelationships(res *ResourceObject, ur UnmarshalRelationships) error {
  relationships := map[string]interface{}{}

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

  err := ur.SetRelationships(relationships)
  if err != nil {
    return err
  }

  return nil
}
