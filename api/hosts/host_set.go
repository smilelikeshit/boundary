// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots at
// 2020-05-03 14:33:51.4625345 -0400 EDT m=+0.042707301
package hosts

import (
	"encoding/json"
	"time"

	"github.com/fatih/structs"
	"github.com/hashicorp/watchtower/api/internal/strutil"
)

type HostSet struct {
	defaultFields []string

	// Canonical path of the resource from the API's base URI
	// Output only.
	Path *string `json:"path,omitempty"`
	// The type of the resource, to help differentiate schemas
	Type *string `json:"type,omitempty"`
	// Friendly name, if set
	FriendlyName *string `json:"friendly_name,omitempty"`
	// The time this host was created
	// Output only.
	CreatedTime time.Time `json:"created_time,omitempty"`
	// The time this host was last updated
	// Output only.
	UpdatedTime time.Time `json:"updated_time,omitempty"`
	// Whether the host set is disabled
	Disabled *bool `json:"disabled,omitempty"`
	// The total count of hosts in this host set
	// Output only.
	Size *int64 `json:"size,omitempty"`
	// A list of hosts in this host set
	// TODO: Figure out if this should be in the basic HostSet view and what
	// view to use on the Hosts.
	// Output only.
	Hosts []*Host `json:"hosts,omitempty"`
}

func (s *HostSet) SetDefault(key string) {
	s.defaultFields = strutil.AppendIfMissing(s.defaultFields, key)
}

func (s *HostSet) UnsetDefault(key string) {
	s.defaultFields = strutil.StrListDelete(s.defaultFields, key)
}

func (s HostSet) MarshalJSON() ([]byte, error) {
	m := structs.Map(s)
	if m == nil {
		m = make(map[string]interface{})
	}
	for _, k := range s.defaultFields {
		m[k] = nil
	}
	return json.Marshal(m)
}
