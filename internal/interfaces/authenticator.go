// Package interfaces
package interfaces

import (
	"net/http"
)

type Authenticator interface {
	GetClient() (http.Client, error)
}
