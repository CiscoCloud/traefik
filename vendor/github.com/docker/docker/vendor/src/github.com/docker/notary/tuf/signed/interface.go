package signed

import (
	"github.com/docker/notary/tuf/data"
)

// KeyService provides management of keys locally. It will never
// accept or provide private keys. Communication between the KeyService
// and a SigningService happen behind the Create function.
type KeyService interface {
	// Create issues a new key pair and is responsible for loading
	// the private key into the appropriate signing service.
	Create(role, gun, algorithm string) (data.PublicKey, error)

	// AddKey adds a private key to the specified role and gun
	AddKey(role, gun string, key data.PrivateKey) error

	// GetKey retrieves the public key if present, otherwise it returns nil
	GetKey(keyID string) data.PublicKey

	// GetPrivateKey retrieves the private key and role if present, otherwise
	// it returns nil
	GetPrivateKey(keyID string) (data.PrivateKey, string, error)

	// RemoveKey deletes the specified key
	RemoveKey(keyID string) error

	// ListKeys returns a list of key IDs for the role
	ListKeys(role string) []string

	// ListAllKeys returns a map of all available signing key IDs to role
	ListAllKeys() map[string]string
}

// CryptoService is deprecated and all instances of its use should be
// replaced with KeyService
type CryptoService interface {
	KeyService
}

// Verifier defines an interface for verfying signatures. An implementer
// of this interface should verify signatures for one and only one
// signing scheme.
type Verifier interface {
	Verify(key data.PublicKey, sig []byte, msg []byte) error
}
