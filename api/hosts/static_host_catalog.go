// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots at
// 2020-05-03 14:33:51.4627143 -0400 EDT m=+0.042887101
package hosts

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
)

type StaticHostCatalog struct {
	*HostCatalog
}

func (s HostCatalog) AsStaticHostCatalog() (*StaticHostCatalog, error) {
	out := &StaticHostCatalog{
		HostCatalog: &s,
	}
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:  out,
		TagName: "json",
	})
	if err != nil {
		return nil, fmt.Errorf("error creating map decoder: %w", err)
	}

	if err := decoder.Decode(s.Attributes); err != nil {
		return nil, fmt.Errorf("error decoding attributes map: %w", err)
	}

	return out, nil
}
