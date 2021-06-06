// Package secret contains implementation to interact with varios secret providers and injectors.
package secret

// Provider defines methods to interact with secret provider
type Provider interface {
	Get (key string) (string, error)
}