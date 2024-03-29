// Code generated by go generate; DO NOT EDIT.
package network

import (
	"net/url"

	"github.com/containers/podman/v5/pkg/bindings/internal/util"
)

// Changed returns true if named field has been set
func (o *ExtraCreateOptions) Changed(fieldName string) bool {
	return util.Changed(o, fieldName)
}

// ToParams formats struct fields to be passed to API service
func (o *ExtraCreateOptions) ToParams() (url.Values, error) {
	return util.ToParams(o)
}

// WithIgnoreIfExists set field IgnoreIfExists to given value
func (o *ExtraCreateOptions) WithIgnoreIfExists(value bool) *ExtraCreateOptions {
	o.IgnoreIfExists = &value
	return o
}

// GetIgnoreIfExists returns value of field IgnoreIfExists
func (o *ExtraCreateOptions) GetIgnoreIfExists() bool {
	if o.IgnoreIfExists == nil {
		var z bool
		return z
	}
	return *o.IgnoreIfExists
}
