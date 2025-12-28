package types

// NationalDefaultRegionID is used when a record is country-wide.
// Some older code referenced this constant but it was missing.
const NationalDefaultRegionID = "UA"

// DefaultParams returns module default parameters.
func DefaultParams() Params {
	// Keep permissive defaults so chain can boot.
	return Params{}
}

// Validate validates the set of params.
func (p Params) Validate() error {
	// Keep permissive for now; add checks when params fields are finalized.
	return nil
}
