package entity

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHasTag(t *testing.T) {
	lead := Lead{
		Tags: []string{
			"doce",
			"solidao",
		},
	}

	require.True(t, lead.HasTag("doce"))
	require.False(t, lead.HasTag("salgado"))
}

func TestEmpty(t *testing.T) {
	emptyLead := Lead{}

	require.True(t, emptyLead.Empty())

	lead := Lead{
		Name:    "biro",
		Website: "https://website.com",
	}

	require.False(t, lead.Empty())
}
