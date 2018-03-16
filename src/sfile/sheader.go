package sfile

import (
	"strings"
)

// HeaderFormat is a simple interface to enforce map of attributes.
type HeaderFormat interface {
	Attributes() map[string]interface{}
	SetAttributes(att map[string]interface{})
	Get(key string) interface{}
	Set(key string, value interface{})
}

// SimpleHeader is an object that just holds a map object of the attributes
type SimpleHeader struct {
	attributes map[string]interface{}
}

// Attributes returns the internal map of attributes
func (sh *SimpleHeader) Attributes() map[string]interface{} {
	return sh.attributes
}

// SetAttributes sets the internal map of attributes with normalization.
func (sh *SimpleHeader) SetAttributes(att map[string]interface{}) {
	m := map[string]interface{}{}
	for key, value := range att {
		m[strings.ToLower(key)] = value
	}
	sh.attributes = m
}

// Get the value for a given key
func (sh *SimpleHeader) Get(key string) interface{} {
	return sh.attributes[strings.ToLower(key)]
}

// Set creates or updates a key and value pair in the map
func (sh *SimpleHeader) Set(key string, value interface{}) {
	sh.attributes[strings.ToLower(key)] = value
}
