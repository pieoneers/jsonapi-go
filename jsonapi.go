package jsonapi

import (
  "bytes"
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

type MarshalData interface {
  GetData() interface{}
}

type UnmarshalData interface {
  SetData(func(interface{}) error) error
}

type MarshalErrors interface {
  GetErrors() []*ErrorObject
}

type UnmarshalErrors interface {
  SetErrors(errors []*ErrorObject) error
}

type MarshalIncluded interface {
  GetIncluded() []interface{}
}

type MarshalMeta interface {
  GetMeta() interface{}
}

type Document struct {
  Data     *documentData     `json:"data,omitempty"`
  Errors   []*ErrorObject    `json:"errors,omitempty"`
  Included []*ResourceObject `json:"included,omitempty"`
  Meta     json.RawMessage   `json:"meta,omitempty"`
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
  Attributes    json.RawMessage          `json:"attributes,omitempty"`
  Meta          json.RawMessage          `json:"meta,omitempty"`
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
  var err error

  buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)

  if d.One != nil {
    err = enc.Encode(d.One)

    return buf.Bytes(), err
  }

  err = enc.Encode(d.Many)

  return buf.Bytes(), err
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
  if d.One != nil && len(d.One.ID) > 0 {
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

  val := reflect.ValueOf(payload)
  i := val.Interface()

  if val.Kind() == reflect.Ptr {
    val = val.Elem()
    i = val.Interface()
  }

  doc, err = marshalDocument(i)
  if err != nil {
    return nil, err
  }

  buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)

  err = enc.Encode(doc)

  return buf.Bytes(), err
}

func marshalDocument(payload interface{}) (*Document, error) {
  doc := &Document{}

  switch asserted := payload.(type) {
	case MarshalData:
    doc.Data = &documentData{}

    data := asserted.GetData()

    switch reflect.TypeOf(data).Kind() {
    case reflect.Struct:
      if one, err := marshalResourceObject(data.(MarshalResourceIdentifier)); err == nil {
        doc.Data.One = &one
      } else {
        return nil, err
      }
    case reflect.Slice:
      if many, err := marshalResourceObjects(data); err == nil {
        doc.Data.Many = many
      } else {
        return nil, err
      }
    }
	case MarshalErrors:
    doc.Errors = asserted.GetErrors()
	}

  if mi, ok := payload.(MarshalIncluded); ok {
    if included, err := marshalIncluded(mi); err == nil {
      doc.Included = included
    } else {
      return nil, err
    }
  }

  if mm, ok := payload.(MarshalMeta); ok {
    if meta, err := marshalMeta(mm); err == nil {
      if !bytes.Equal(meta, []byte("{}\n")) {
        doc.Meta = meta
      }
    } else {
      return nil, err
    }
  }

  return doc, nil
}

func marshalResourceObjectIdentifier(mri MarshalResourceIdentifier) ResourceObjectIdentifier {
  return ResourceObjectIdentifier{ ID: mri.GetID(), Type: mri.GetType() }
}

func marshalResourceObject(mri MarshalResourceIdentifier) (ResourceObject, error) {
  one := ResourceObject{
    ResourceObjectIdentifier: marshalResourceObjectIdentifier(mri),
  }

  buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)

  err := enc.Encode(mri)
  if err != nil {
    return one, err
  }

  attributes := buf.Bytes()

  if !bytes.Equal(attributes, []byte("{}\n")) {
    one.Attributes = attributes
  }

  if mm, ok := mri.(MarshalMeta); ok {
    if meta, err := marshalMeta(mm); err == nil {
      if !bytes.Equal(meta, []byte("{}\n")) {
        one.Meta = meta
      }
    } else {
      return one, err
    }
  }

  if mr, ok := mri.(MarshalRelationships); ok {
    one.Relationships = marshalRelationships(mr)
  }

  return one, nil
}

func marshalResourceObjects(payload interface{}) ([]*ResourceObject, error) {
  many := []*ResourceObject{}

  value := reflect.ValueOf(payload)

  for i := 0; i < value.Len(); i++ {
    one, err := marshalResourceObject(value.Index(i).Interface().(MarshalResourceIdentifier))
    if err != nil {
      return many, err
    }

    many = append(many, &one)
  }

  return many, nil
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
    Data: &relationshipData{
      Many: make([]*ResourceObjectIdentifier, 0),
    },
  }

  value := reflect.ValueOf(payload)

  for i := 0; i < value.Len(); i++ {
    one := marshalResourceObjectIdentifier(value.Index(i).Interface().(MarshalResourceIdentifier))
    relationship.Data.Many = append(relationship.Data.Many, &one)
  }

  return relationship
}

func marshalIncluded(mi MarshalIncluded) ([]*ResourceObject, error) {
  var included []*ResourceObject

  for _, value := range mi.GetIncluded() {
    ro, err := marshalResourceObject(value.(MarshalResourceIdentifier))
    if err != nil {
      return included, err
    }

    included = append(included, &ro)
  }

  return included, nil
}

func marshalMeta(mm MarshalMeta) (json.RawMessage, error) {
  buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)

  meta := mm.GetMeta()

  err := enc.Encode(meta)

  return buf.Bytes(), err
}

func Unmarshal(data []byte, target interface{}) (*Document, error) {
  doc := &Document{}

  if err := json.Unmarshal(data, doc); err != nil {
    return doc, err
  }

  switch asserted := target.(type) {
  case UnmarshalData:
    if one := doc.Data.One; one != nil {
      if err := asserted.SetData(func(target interface{}) error {
        return unmarshalOne(one, target)
      }); err != nil {
        return doc, err
      }
    }

    if many := doc.Data.Many; many != nil {
      if err := asserted.SetData(func(target interface{}) error {
        return unmarshalMany(many, target)
      }); err != nil {
        return doc, err
      }
    }
  case UnmarshalErrors:
    if errors := doc.Errors; errors != nil {
      asserted.SetErrors(errors)
    }
  }

  return doc, nil
}

func unmarshalOne(one *ResourceObject, target interface{}) error {
  return unmarshalResourceObject(one, target.(UnmarshalResourceIdentifier))
}

func unmarshalMany(many []*ResourceObject, target interface{}) error {
  ptr := reflect.ValueOf(target)
  val := ptr.Elem()

  typ := reflect.TypeOf(target).Elem().Elem()
  knd := typ.Kind()

  if knd == reflect.Ptr {
    typ = typ.Elem()
  }

  for _, one := range many {
    new := reflect.New(typ)

    if err := unmarshalResourceObject(one, new.Interface().(UnmarshalResourceIdentifier)); err != nil {
      return err
    }

    if knd == reflect.Struct {
      new = new.Elem()
    }

    val = reflect.Append(val, new)
  }

  ptr.Elem().Set(val)

  return nil
}

func unmarshalResourceObject(ro *ResourceObject, ui UnmarshalResourceIdentifier) error {
  if len(ro.Attributes) > 0 {
    if err := json.Unmarshal(ro.Attributes, ui); err != nil {
      return err
    }
  }

  if err := ui.SetID(ro.ID); err != nil {
    return err
  }

  if ur, ok := ui.(UnmarshalRelationships); ok {
    if err := unmarshalRelationships(ro, ur); err != nil {
      return err
    }
  }

  return nil
}

func unmarshalRelationships(ro *ResourceObject, ur UnmarshalRelationships) error {
  relationships := map[string]interface{}{}

  for k, v := range ro.Relationships {
    if one := v.Data.One; one != nil {
      relationships[k] = one
    }

    if many := v.Data.Many; many != nil {
      relationships[k] = many
    }
  }

  if err := ur.SetRelationships(relationships); err != nil {
    return err
  }

  return nil
}
