package oidc

import "strings"

type Scope string

const (
	ScopeOpenID  Scope = "openid"
	ScopeOffline Scope = "offline_access"
)

type Scopes []Scope

func (s Scopes) String() string {
	strs := make([]string, len(s))
	for i, scope := range s {
		strs[i] = string(scope)
	}
	return strings.Join(strs, " ")
}

type ScopeString string

func (s ScopeString) Array() Scopes {
	strs := strings.FieldsFunc(string(s), func(r rune) bool {
		switch r {
		case ' ', '\t', '\n', '\r':
			return true
		}
		return false
	})
	scopes := make([]Scope, len(strs))
	for i, str := range strs {
		scopes[i] = Scope(str)
	}
	return scopes
}
