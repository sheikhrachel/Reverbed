package handlers

const (
	StatusOK                  = 200
	StatusBadRequest          = 400
	StatusNotFound            = 404
	StatusInternalServerError = 500

	// core

	PathRoot   = "/"
	PathHealth = "/health"

	// editing

	PathEditTrack = "/edit_track/:edit_type"
)
