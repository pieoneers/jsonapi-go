// Copyright (c) 2020 Pieoneers Software Incorporated. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package jsonapi

import (
	"bytes"
	"encoding/json"
	"reflect"
)

// ContentType describes data content type.
const ContentType = "application/vnd.api+json"

// MarshalResourceIdentifier interface should be implemented to be able marshal Go struct into JSON API document.
//
// GetID example:
//
//    func(s SomeStruct) GetID() string {
//      return s.ID
//    }
//
//    func(s SomeStruct) GetType() string {
//      return d.Type
//    }
//
// or
//
//    func(s SomeStruct) GetType() string {
//      return "some-resource-type"
//    }
//
type MarshalResourceIdentifier interface {
	GetID() string
	GetType() string
}

// UnmarshalResourceIdentifier interface should be implemented to be able unmarshal JSON API document into Go struct.
//
// SetID, SetType examples:
//
//    func(s *SomeStruct) SetID(id string) error {
//      s.ID = id
//      return nil
//    }
//
//    func(s *SomeStruct) SetType(t string) error {
//      s.Type = t
//      return nil
//    }
//
// or
//
//    func(s *SomeStruct) SetType(string) error {
//      return nil
//    }
//
type UnmarshalResourceIdentifier interface {
	SetID(string) error
	SetType(string) error
}

// MarshalRelationships interface should be implemented to be able marshal JSON API document relationships.
//
// GetRelationships example:
//
//    func(s SomeStruct) GetRelationships() map[string]interface{} {
//      relationships := make(map[string]interface{})
//
//      relationships["relation"] = jsonapi.ResourceObjectIdentifier{
//        ID: s.RelationID,
//        Type: "relation-type",
//      }
//
//      return relationships
//    }
//
type MarshalRelationships interface {
	GetRelationships() map[string]interface{}
}

// UnmarshalRelationships interface should be implemented to be able unmarshal JSON API document relationships into Go struct.
//
// SetRelationships example:
//
//    func (s *SomeStruct) SetRelationships(relationships map[string]interface{}) error {
//    	if relationship, ok := relationships["relation"]; ok {
//    		s.RealtionID = relationship.(*jsonapi.ResourceObjectIdentifier).ID
//    	}
//
//    	return nil
//    }
//
type UnmarshalRelationships interface {
	SetRelationships(map[string]interface{}) error
}

// MarshalData interface should be implemented to be able get data from Go struct and marshal it.
//
// GetData example:
//
//    func(s SomeStruct) GetData() interface{} {
//      return s
//    }
//
type MarshalData interface {
	GetData() interface{}
}

// UnmarshalData interface should be implemented to be able unmarshal data from JSON API document into Go struct.
//
// SetData example:
//
//    func(s *SomeStruct) SetData(to func(target interface{}) error) error {
//      return to(s)
//    }
//
// NOTE: If you are using SomeStruct collections, you should implement additional data type, e.g.:
// type SomeStructs []SomeStruct
//
// Then you should implement SetData method for SomeStructs:
//
//    func(s *SomeStructs) SetData(to func(target interface{}) error) error {
//      return to(s)
//    }
//
type UnmarshalData interface {
	SetData(func(interface{}) error) error
}

// MarshalErrors interface should be implemented to be able marshal errors into JSON API document.
//
// GetErrors example:
//
//    type SomeErrorType struct {
//      Code string
//      Title string
//      Pointer string
//    }
//
//    type SomeErrorTypes []SomeErrorType
//
//    func(e SomeErrors) GetErrors() []*jsonapi.ErrorObject {
//      var errs []*jsonapi.ErrorObject
//
//      for _, err := range e {
//        errs = append(errs, &jsonapi.ErrorObject{
//          Title: err.Title,
//          Code: err.Code,
//          Source: jsonapi.ErrorObjectSource{
//            Pointer: err.Pointer,
//          },
//        })
//      }
//
//      return errs
//    }
//
type MarshalErrors interface {
	GetErrors() []*ErrorObject
}

// UnmarshalErrors interface should be implemented to be able unmarshal errors from JSON API document.
//
// SetErrors example:
//
//   type SomeError struct {
//      Code string
//      Title string
//      Pointer string
//    }
//
//   type SomeErrors struct {
//     Errors []SomeError
//   }
//
//   func(v SomeErrors) SetErrors(errs []*jsonapi.ErrorObject) error {
//     var someErrors []SomeError
//
//     for _, err := range errs {
//       someErrors = append(someErrors, SomeError{
//         Title: err.Title,
//         Code: err.Code,
//         Pointer: err.Source.Pointer,
//       })
//     }
//
//     v.Errors = someErrors
//
//     return nil
//   }
//
type UnmarshalErrors interface {
	SetErrors(errors []*ErrorObject) error
}

// MarshalIncluded interface should be implemented to be able marshal JSON API document included.
//
// GetIncluded example:
//
//    func(v SomeStruct) GetIncluded() []interface{} {
//      var included []interface{}
//
//      /*
//        Get some additional data here and put it into `items` variables
//      `items` data type should implement MarshalResourceIdentifier and MarshalData interface.
//      */
//      for _, item := range items {
//        included = append(included, item)
//      }
//
//      return included
//    }
//
type MarshalIncluded interface {
	GetIncluded() []interface{}
}

// MarshalMeta interface should be implemented to be able marshal JSON API document meta.
//
// GetMeta example:
//
//    type Meta struct {
//      Count int `json:"count"`
//    }
//
//    func(v SomeStruct) GetMeta() interface{} {
//      return Meta{ Count: 42 }
//    }
//
type MarshalMeta interface {
	GetMeta() interface{}
}

// Document describes Go representation of JSON API document.
type Document struct {
	// Document data
	Data *documentData `json:"data,omitempty"`
	// Document errors
	Errors []*ErrorObject `json:"errors,omitempty"`
	// Document included
	Included []*ResourceObject `json:"included,omitempty"`
	// Document meta
	Meta json.RawMessage `json:"meta,omitempty"`
}

type documentData struct {
	One  *ResourceObject
	Many []*ResourceObject
}

type relationship struct {
	Data *relationshipData `json:"data"`
}

type relationshipData struct {
	One  *ResourceObjectIdentifier
	Many []*ResourceObjectIdentifier
}

// ResourceObjectIdentifier JSON API resource object.
type ResourceObjectIdentifier struct {
	Type string `json:"type"`
	ID   string `json:"id,omitempty"`
}

// GetID method returns ResourceObjectIdentifier ID.
func (roi ResourceObjectIdentifier) GetID() string {
	return roi.ID
}

// GetType method returns ResourceObjectIdentifier Type.
func (roi ResourceObjectIdentifier) GetType() string {
	return roi.Type
}

// ResourceObject extends ResourceObjectIdentifier with JSON API document Attributes, Meta and Relationships.
type ResourceObject struct {
	ResourceObjectIdentifier
	// Attributes JSON API document attributes raw data.
	Attributes json.RawMessage `json:"attributes,omitempty"`
	// Meta JSON API document meta raw data.
	Meta json.RawMessage `json:"meta,omitempty"`
	// Relationships JSON API document relationships raw data.
	Relationships map[string]*relationship `json:"relationships,omitempty"`
}

// ErrorObject JSON API error object https://jsonapi.org/format/#error-objects
type ErrorObject struct {
	// Title a short, human-readable summary of the problem.
	Title string `json:"title,omitempty"`
	// Code application specified value to identify the error.
	Code string `json:"code,omitempty"`
	// Source an object containing references to the source of the error.
	Source ErrorObjectSource `json:"source,omitempty"`
}

// ErrorObjectSource includes pointer ErrorObject.Source
type ErrorObjectSource struct {
	// Pointer a JSON Pointer [RFC6901] to the associated entity in the request document [e.g. "/data" for a primary data object, or "/data/attributes/title" for a specific attribute].
	Pointer string `json:"pointer,omitempty"`
}

func (d *documentData) MarshalJSON() ([]byte, error) {
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

func (d *documentData) UnmarshalJSON(payload []byte) error {
	if bytes.HasPrefix(payload, []byte("{")) {
		return json.Unmarshal(payload, &d.One)
	}

	if bytes.HasPrefix(payload, []byte("[")) {
		return json.Unmarshal(payload, &d.Many)
	}

	return nil
}

func (d *relationshipData) MarshalJSON() ([]byte, error) {
	if d.One != nil && len(d.One.ID) > 0 {
		return json.Marshal(d.One)
	}
	return json.Marshal(d.Many)
}

func (d *relationshipData) UnmarshalJSON(payload []byte) error {
	if bytes.HasPrefix(payload, []byte("{")) {
		return json.Unmarshal(payload, &d.One)
	}

	if bytes.HasPrefix(payload, []byte("[")) {
		return json.Unmarshal(payload, &d.Many)
	}

	return nil
}

// Marshal serialize Go struct into []byte JSON API document
// If the corresponding interfaces are implemented the output will contain, relationships, included, meta and errors.
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
	return ResourceObjectIdentifier{ID: mri.GetID(), Type: mri.GetType()}
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

// Unmarshal deserialize JSON API document into Gu sturct
// If the corresponding interfaces are implemented target will contain data from JSON API document relationships and errors.
func Unmarshal(data []byte, target interface{}) (*Document, error) {
	doc := &Document{}

	if err := json.Unmarshal(data, doc); err != nil {
		return doc, err
	}

	if asserted, ok := target.(UnmarshalData); ok && doc.Data != nil {

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
	}

	if asserted, ok := target.(UnmarshalErrors); ok && doc.Errors != nil {
		asserted.SetErrors(doc.Errors)
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

	if err := ui.SetType(ro.ResourceObjectIdentifier.Type); err != nil {
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
		data := v.Data

		if data != nil {
			if one := data.One; one != nil {
				relationships[k] = one
			}

			if many := data.Many; many != nil {
				relationships[k] = many
			}
		}
	}

	if err := ur.SetRelationships(relationships); err != nil {
		return err
	}

	return nil
}
