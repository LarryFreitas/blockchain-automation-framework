package model

import (
	"fmt"
	"time"
)

type (
	// SwID represents a software identification tag
	SwID struct {
		// PrimaryTag identifies the software asset
		PrimaryTag string `json:"primary_tag"`
		// XML is the contents of the SwID document in xml format
		XML string `json:"xml"`
		// Asset is the ID of the associated license
		Asset string `json:"asset"`
		// License is the ID of the associated license
		License string `json:"license"`
		// LeaseExpiration is the date when the lease associated with this SwID expires
		LeaseExpiration time.Time `json:"lease_expiration"`
	}
)

const SwIDPrefix = "swid:"

// SwIDKey returns the key for a swid tag on the ledger.  SwIDs are stored with the format: "swid:<primary_tag>".
func SwIDKey(name string) string {
	return fmt.Sprintf("%s%s", SwIDPrefix, name)
}
