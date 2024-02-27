package rdstation

import (
	"github.com/flan6/rdstation/entity"
	"github.com/flan6/rdstation/internal/client"
)

const (
	AcademyTagCancelled = entity.AcademyTagCancelled
	AcademyActive       = entity.AcademyActive
	EmailOptOut         = entity.EmailOptOut
)

type (
	Lead    = entity.Lead
	Token   = entity.Token
	Secret  = entity.Secret
	RDError = client.RDError
	Errors  = client.Errors
)
