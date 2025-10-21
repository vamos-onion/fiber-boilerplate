package util

import (
	"encoding/json"
	"reflect"
)

type structBlock struct{}

// ReflectBlock :
type ReflectBlock struct {
	TagName string
	Type    reflect.Type
	Value   reflect.Value
}

// FieldBlock :
type FieldBlock struct {
	Tag   string
	Value reflect.Value
}

func (s structBlock) Reflect(src interface{}, tagName string) (refl *ReflectBlock) {
	refl = &ReflectBlock{
		TagName: tagName,
	}

	refl.Value = reflect.ValueOf(src)

	refl.Type = reflect.TypeOf(src)
	if refl.Type.Kind() == reflect.Ptr {
		refl.Type = refl.Type.Elem()

		if refl.Value.IsNil() {
			refl.Value = reflect.New(refl.Type)
		}

	}

	refl.Value = reflect.Indirect(refl.Value)

	return
}

func (s structBlock) Field(refl *ReflectBlock, idx int) *FieldBlock {
	return &FieldBlock{
		Tag:   refl.Type.Field(idx).Tag.Get(refl.TagName),
		Value: refl.Value.Field(idx),
	}
}

func (s structBlock) ToMap(src interface{}, tagName string) (objMap map[string]interface{}) {
	refl := s.Reflect(src, tagName)

	objMap = make(map[string]interface{})
	for i := 0; i < refl.Type.NumField(); i++ {
		field := s.Field(refl, i)
		if field.Tag != "" && field.Tag != "-" {
			objMap[field.Tag] = field.Value.Interface()
		}
	}

	return
}

func (s structBlock) ToStringMap(src interface{}, tagName string) (strMap map[string]string) {
	refl := s.Reflect(src, tagName)

	strMap = make(map[string]string)
	for i := 0; i < refl.Type.NumField(); i++ {
		field := s.Field(refl, i)
		if field.Tag != "" && field.Tag != "-" {
			strMap[field.Tag] = field.Value.String()
		}
	}

	return
}

func (s structBlock) Convert(src interface{}, dst interface{}) (err error) {
	bytes, err := json.Marshal(src)
	if err != nil {
		return
	}

	err = json.Unmarshal(bytes, dst)
	return
}
