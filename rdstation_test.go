package rdstation

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/flan6/rdstation/entity"
	"github.com/flan6/rdstation/test/mocks"
)

func TestNewRdStation(t *testing.T) {
	rd := NewRDStation("", "", "")
	require.NotNil(t, rd)
}

func TestRdStation_GetLeadByEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockClient(ctrl)

	rd := &rdStation{client: client}
	require.NotNil(t, rd)

	t.Run("success", func(t *testing.T) {
		lead := entity.Lead{
			Name:  "nome",
			Email: "email",
		}
		ret, err := json.Marshal(lead)
		require.NoError(t, err)

		client.EXPECT().Request(fmt.Sprintf("%s%semail:%s", RDURL, RDLeadPath, lead.Email), http.MethodGet, nil).
			Return(ret, nil)

		res, err := rd.GetLeadByEmail(lead.Email)

		require.NoError(t, err)
		require.Equal(t, lead.Name, res.Name)
	})

	t.Run("error", func(t *testing.T) {
		lead := entity.Lead{Email: "email"}

		client.EXPECT().Request(fmt.Sprintf("%s%semail:%s", RDURL, RDLeadPath, lead.Email), http.MethodGet, nil).
			Return(nil, errors.New("err"))

		res, err := rd.GetLeadByEmail(lead.Email)

		require.Error(t, err)
		require.Empty(t, res)
	})
}

func TestRdStation_CreateLead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockClient(ctrl)

	rd := &rdStation{client: client}
	require.NotNil(t, rd)

	lead := entity.Lead{
		Name:  "nome",
		Email: "email",
	}
	data, err := json.Marshal(lead)
	require.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		expected := entity.Lead{
			Uuid:  "aksjdnasd",
			Email: "email",
			Name:  "nome",
		}
		response, err := json.Marshal(
			entity.Lead{
				Uuid:  "aksjdnasd",
				Email: "email",
				Name:  "nome",
			})

		require.NoError(t, err)

		client.EXPECT().Request(fmt.Sprintf("%s%s", RDURL, RDLeadPath), http.MethodPost, data).
			Return(response, nil)

		got, err := rd.CreateLead(&lead)

		require.NoError(t, err)
		require.Equal(t, expected, *got)
	})

	t.Run("error", func(t *testing.T) {
		target := errors.New("batata")
		client.EXPECT().Request(fmt.Sprintf("%s%s", RDURL, RDLeadPath), http.MethodPost, data).
			Return(nil, target)

		got, err := rd.CreateLead(&lead)
		require.Equal(t, target, err)
		require.Error(t, err)
		require.Nil(t, got)
	})
}

func TestRdStation_DeleteLeadByEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockClient(ctrl)

	rd := &rdStation{client: client}
	require.NotNil(t, rd)

	t.Run("success", func(t *testing.T) {
		lead := entity.Lead{Email: "email"}

		client.EXPECT().Request(fmt.Sprintf("%s%semail:%s", RDURL, RDLeadPath, lead.Email), http.MethodDelete, nil).
			Return(nil, nil)

		err := rd.DeleteLeadByEmail(lead.Email)
		require.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		lead := entity.Lead{Email: "email"}

		client.EXPECT().Request(fmt.Sprintf("%s%semail:%s", RDURL, RDLeadPath, lead.Email), http.MethodDelete, nil).
			Return(nil, errors.New("err"))

		err := rd.DeleteLeadByEmail(lead.Email)
		require.Error(t, err)
	})
}

func TestRdStation_UpdateLead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockClient(ctrl)

	rd := &rdStation{client: client}
	require.NotNil(t, rd)

	lead := entity.Lead{
		Name:  "nome",
		Email: "email",
	}
	data, err := json.Marshal(lead)
	require.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		client.EXPECT().Request(fmt.Sprintf("%s%semail:%s", RDURL, RDLeadPath, lead.Email), http.MethodPatch, data).
			Return(nil, nil)

		err = rd.UpdateLead(&lead)
		require.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		client.EXPECT().Request(fmt.Sprintf("%s%semail:%s", RDURL, RDLeadPath, lead.Email), http.MethodPatch, data).
			Return(nil, errors.New("ah sei la"))

		err = rd.UpdateLead(&lead)
		require.Error(t, err)
	})
}

func TestRdStation_AddTags(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockClient(ctrl)

	rd := &rdStation{client: client}
	require.NotNil(t, rd)

	t.Run("add tags", func(t *testing.T) {
		tests := map[string]struct {
			lead entity.Lead
			tags []string
			want []string
		}{
			"no tag":                      {lead: entity.Lead{Email: "test@test.com"}, tags: []string{"exac"}, want: []string{"exac"}},
			"insert academy active":       {lead: entity.Lead{Email: "test@test.com"}, tags: []string{"acativo"}, want: []string{"acativo"}},
			"unusual tag":                 {lead: entity.Lead{Email: "test@test.com", Tags: []string{"exac"}}, tags: []string{"biro loco"}, want: []string{"exac", "biro loco"}},
			"duplicated academy canceled": {lead: entity.Lead{Email: "test@test.com", Tags: []string{"exac", "exac", "potato"}}, tags: []string{"exac"}, want: []string{"potato", "exac"}},
			"has academy canceled, insert academy active": {lead: entity.Lead{Email: "test@test.com", Tags: []string{"exac"}}, tags: []string{"acativo"}, want: []string{"acativo"}},
			"has academy active, insert academy canceled": {lead: entity.Lead{Email: "test@test.com", Tags: []string{"acativo"}}, tags: []string{"exac"}, want: []string{"exac"}},
		}

		for _, test := range tests {
			data, err := json.Marshal(map[string][]string{"tags": test.want})
			require.NoError(t, err)

			client.EXPECT().Request(fmt.Sprintf("%s%semail:%s", RDURL, RDLeadPath, test.lead.Email), http.MethodPatch, data).
				Return(nil, nil)

			err = rd.AddTags(&test.lead, test.tags)
			require.NoError(t, err)
			require.Equal(t, test.want, test.lead.Tags)
		}
	})
}

func TestRdStation_RemoveTags(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockClient(ctrl)

	rd := &rdStation{client: client}
	require.NotNil(t, rd)

	t.Run("remove tags", func(t *testing.T) {
		tests := map[string]struct {
			lead entity.Lead
			tags []string
			want []string
		}{
			"has tag":                      {lead: entity.Lead{Email: "test@test.com", Tags: []string{"exac"}}, tags: []string{"exac"}, want: []string{}},
			"remove both academy canceled": {lead: entity.Lead{Email: "test@test.com", Tags: []string{"exac", "exac", "potato"}}, tags: []string{"exac"}, want: []string{"potato"}},
			"no change":                    {lead: entity.Lead{Email: "test@test.com", Tags: []string{"exac"}}, tags: []string{"acativo"}, want: []string{"exac"}},
		}

		for _, test := range tests {
			data, err := json.Marshal(map[string][]string{"tags": test.want})
			require.NoError(t, err)

			client.EXPECT().Request(fmt.Sprintf("%s%semail:%s", RDURL, RDLeadPath, test.lead.Email), http.MethodPatch, data).
				Return(nil, nil)

			err = rd.RemoveTags(&test.lead, test.tags)
			require.NoError(t, err)
			require.Equal(t, test.want, test.lead.Tags)
		}
	})

	t.Run("has nil tags", func(t *testing.T) {
		tests := map[string]struct {
			lead entity.Lead
			tags []string
			want []string
		}{
			"lead has nil tags": {lead: entity.Lead{Email: "test@test.com"}, tags: []string{"exac"}, want: []string(nil)},
		}

		for _, test := range tests {
			empty := []string{}
			data, err := json.Marshal(map[string][]string{"tags": empty})
			require.NoError(t, err)

			client.EXPECT().Request(fmt.Sprintf("%s%semail:%s", RDURL, RDLeadPath, test.lead.Email), http.MethodPatch, data).
				Return(nil, nil)

			err = rd.RemoveTags(&test.lead, test.tags)
			require.NoError(t, err)
			require.Equal(t, test.want, test.lead.Tags)
		}
	})
}
