package ciba

import (
	"net/url"
	"time"

	"github.com/daedaluz/goauth2/oidc"
)

type Option interface {
	Apply(*AuthSession)
}

type optionFunc func(*AuthSession)

func (f optionFunc) Apply(s *AuthSession) {
	f(s)
}

func WithLoginHint(hint string) Option {
	return optionFunc(func(s *AuthSession) {
		s.loginHint = hint
	})
}

func WithLoginHintToken(token string) Option {
	return optionFunc(func(s *AuthSession) {
		s.loginHintToken = token
	})
}

func WithIDTokenHint(token string) Option {
	return optionFunc(func(s *AuthSession) {
		s.idTokenHint = token
	})
}

func WithBindingMessage(msg string) Option {
	return optionFunc(func(s *AuthSession) {
		s.bindingMessage = msg
	})
}

func WithRequestedExpiry(d time.Duration) Option {
	return optionFunc(func(s *AuthSession) {
		s.requestedExpiry = d
	})
}

func WithPollInterval(d time.Duration) Option {
	return optionFunc(func(s *AuthSession) {
		s.pollInterval = d
	})
}

func WithACRValues(acr ...oidc.ACR) Option {
	return optionFunc(func(s *AuthSession) {
		for _, a := range acr {
			s.acrValues = append(s.acrValues, string(a))
		}
	})
}

func WithScope(scopes oidc.Scopes) Option {
	return optionFunc(func(s *AuthSession) {
		for _, sc := range scopes {
			s.scope = append(s.scope, string(sc))
		}
	})
}

func WithValues(v url.Values) Option {
	return optionFunc(func(s *AuthSession) {
		s.values = v
	})
}
