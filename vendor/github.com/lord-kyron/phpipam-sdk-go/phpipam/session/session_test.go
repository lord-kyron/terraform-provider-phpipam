package session

import (
	"reflect"
	"testing"

	"github.com/pavel-z1/phpipam-sdk-go/phpipam"
)

func phpipamConfig() phpipam.Config {
	return phpipam.Config{
		AppID:    "0123456789abcdefgh",
		Endpoint: "http://localhost/api",
		Password: "changeit",
		Username: "nobody",
	}
}

func fullSessionConfig() *Session {
	return &Session{
		Config: phpipamConfig(),
		Token: Token{
			String: "foobarbazboop",
		},
	}
}

func TestNewSession(t *testing.T) {
	cfg := phpipamConfig()

	expected := &Session{
		Config: phpipamConfig(),
	}

	actual := NewSession(cfg)

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected session to be %#v, got %#v", expected, actual)
	}
}
