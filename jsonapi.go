package jsonapi

import (
  "sort"
  "bytes"
  "strings"
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

type Document struct {
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

func(roi ResourceObjectIdentifier) GetID() string {
  return roi.ID
}

func(roi ResourceObjectIdentifier) GetType() string {
  return roi.Type
}

type ResourceObject struct {
  ResourceObjectIdentifier
  Attributes    json.RawMessage          `json:"attributes"`
  Relationships map[string]*relationship `json:"relationships,omitempty"`
}

type ErrorObject struct {
  Title  string            `json:"title,omitempty"`
  Source ErrorObjectSource `json:"source,omitempty"`
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

func(d *documentData) UnmarshalJSON(payload []byte) error {
  if bytes.HasPrefix(payload, []byte("{")) {
    return json.Unmarshal(payload, &d.One)
  }

  if bytes.HasPrefix(payload, []byte("[")) {
    return json.Unmarshal(payload, &d.Many)
  }

  return nil
}

func(d *relationshipData) MarshalJSON() ([]byte, error) {
  if d.One != nil {
    return json.Marshal(d.One)
  }
  return json.Marshal(d.Many)
}

func(d *relationshipData) UnmarshalJSON(payload []byte) error {
  if bytes.HasPrefix(payload, []byte("{")) {
    return json.Unmarshal(payload, &d.One)
  }

  if bytes.HasPrefix(payload, []byte("[")) {
    return json.Unmarshal(payload, &d.Many)
  }

  return nil
}

func Marshal(payload interface{}) ([]byte, error) {
  var (
    doc *Document
    err error
  )

  doc, err = marshalDocument(payload)
  if err != nil {
    return nil, err
  }

  return json.Marshal(doc)
}

func marshalDocument(payload interface{}) (*Document, error) {
  var (
    doc *Document
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

  sort.Slice(doc.Included, func(i, j int) bool {
    return strings.Compare(doc.Included[i].Type, doc.Included[j].Type) == -1
  })

  sort.Slice(doc.Included, func(i, j int) bool {
    return strings.Compare(doc.Included[i].ID, doc.Included[j].ID) == -1
  })

  return doc, nil
}

func marshalDocumentStruct(payload interface{}) (*Document, error) {
  doc := &Document{
    Data: &documentData{},
  }

  one, included, err := marshalResourceObject(payload.(MarshalResourceIdentifier))
  if err != nil {
    return nil, err
  }

  doc.Data.One = &one

  for _, inc := range included {
    for _, one := range inc {
      doc.Included = append(doc.Included, one)
    }
  }

  return doc, nil
}

func marshalDocumentSlice(payload interface{}) (*Document, error) {
  var doc *Document

  if errorObjects, ok := payload.([]*ErrorObject); ok {
    doc = &Document{
      Errors: errorObjects,
    }
  } else {
    doc = &Document{
      Data: &documentData{
        Many: []*ResourceObject{},
      },
    }

    included := map[string]map[string]*ResourceObject{}

    value := reflect.ValueOf(payload)

    for i := 0; i < value.Len(); i++ {
      one, inc, err := marshalResourceObject(value.Index(i).Interface().(MarshalResourceIdentifier))
      if err != nil {
        return nil, err
      }

      doc.Data.Many = append(doc.Data.Many, &one)

      for typ, value := range inc {
        if _, ok := included[typ]; !ok {
          included[typ] = make(map[string]*ResourceObject)
        }

        for id, one := range value {
          if _, ok := included[typ][id]; !ok {
            included[typ][id] = one
          }
        }
      }
    }

    for _, inc := range included {
      for _, one := range inc {
        doc.Included = append(doc.Included, one)
      }
    }
  }

  return doc, nil
}

func marshalResourceObjectIdentifier(mri MarshalResourceIdentifier) ResourceObjectIdentifier {
  return ResourceObjectIdentifier{ ID: mri.GetID(), Type: mri.GetType() }
}

func marshalResourceObject(mri MarshalResourceIdentifier) (ResourceObject, map[string]map[string]*ResourceObject, error) {
  var included map[string]map[string]*ResourceObject

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
  var relationship *relationship

  switch reflect.TypeOf(payload).Kind() {
  case reflect.Struct:
    relationship = marshalRelationshipStruct(payload)
  case reflect.Slice:
    relationship = marshalRelationshipSlice(payload)
  }

  return relationship
}

func marshalRelationshipStruct(payload interface{}) *relationship {
  relationship := &relationship{
    Data: &relationshipData{},
  }

  one := marshalResourceObjectIdentifier(payload.(MarshalResourceIdentifier))
  relationship.Data.One = &one

  return relationship
}

func marshalRelationshipSlice(payload interface{}) *relationship {
  relationship := &relationship{
    Data: &relationshipData{},
  }

  value := reflect.ValueOf(payload)

  for i := 0; i < value.Len(); i++ {
    one := marshalResourceObjectIdentifier(value.Index(i).Interface().(MarshalResourceIdentifier))
    relationship.Data.Many = append(relationship.Data.Many, &one)
  }

  sort.Slice(relationship.Data.Many, func(i, j int) bool {
    return strings.Compare(relationship.Data.Many[i].ID, relationship.Data.Many[j].ID) == -1
  })

  return relationship
}

func marshalIncluded(mi MarshalIncluded) (map[string]map[string]*ResourceObject, error) {
  included := make(map[string]map[string]*ResourceObject)

  for _, value := range mi.GetIncluded() {
    inc, _, err := marshalResourceObject(value.(MarshalResourceIdentifier))
    if err != nil {
      return included, err
    }

    typ, id := inc.Type, inc.ID

    if _, ok := included[typ]; !ok {
      included[typ] = make(map[string]*ResourceObject)
    }

    if _, ok := included[typ][id]; !ok {
      included[typ][id] = &inc
    }
  }

  return included, nil
}

func Unmarshal(data []byte, target interface{}) (*Document, error) {
  var err error

  doc := &Document{}

  err = json.Unmarshal(data, doc)
  if err != nil {
    return doc, err
  }

  errs := doc.Errors
  if errs != nil {
    return doc, err
  }

  one := doc.Data.One
  if one != nil {
    err = unmarshalOne(one, target)
  }

  many := doc.Data.Many
  if many != nil {
    err = unmarshalMany(many, target)
  }

  return doc, err
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

func unmarshalResourceObject(ro *ResourceObject, ui UnmarshalResourceIdentifier) error {
  var err error

  err = json.Unmarshal(ro.Attributes, ui)
  if err != nil {
    return err
  }

  err = ui.SetID(ro.ID)
  if err != nil {
    return err
  }

  if ur, ok := ui.(UnmarshalRelationships); ok {
    err = unmarshalRelationships(ro, ur)
    if err != nil {
      return err
    }
  }

  return nil
}

func unmarshalRelationships(ro *ResourceObject, ur UnmarshalRelationships) error {
  relationships := map[string]interface{}{}

  for k, v := range ro.Relationships {
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
