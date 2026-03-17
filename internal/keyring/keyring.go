package keyring

import (
	"crypto/md5"
	"fmt"

	"github.com/zalando/go-keyring"
)

const baseServiceName = "gearsec"

func getServiceName(profile, key string) string {
	if profile == "" {
		profile = "default"
	}
	hash := md5.Sum([]byte(fmt.Sprintf("%s-%s", profile, key)))
	return fmt.Sprintf("%s-%x", baseServiceName, hash)
}

// Set stores a credential in the OS keyring.
func Set(profile, key, value string) error {
	return keyring.Set(getServiceName(profile, key), key, value)
}

// Get retrieves a credential from the OS keyring.
func Get(profile, key string) (string, error) {
	return keyring.Get(getServiceName(profile, key), key)
}

// Delete removes a credential from the OS keyring.
func Delete(profile, key string) error {
	return keyring.Delete(getServiceName(profile, key), key)
}
