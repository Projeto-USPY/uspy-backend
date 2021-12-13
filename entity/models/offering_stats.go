package models

// OfferingStats represents the evaluation stats for an offering.
//
// It is not a DTO and is only used in internal logic.
type OfferingStats struct {
	Approval    int
	Disapproval int
	Neutral     int
}
