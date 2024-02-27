package rdstation

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/flan6/rdstation/entity"
	"github.com/flan6/rdstation/internal/client"

	"golang.org/x/oauth2"
)

const (
	RDURL           = "https://api.rd.services/"
	RDLeadPath      = "platform/contacts/"
	RefreshTokenURL = "auth/token/"
)

type RDStation interface {
	GetLeadByEmail(email string) (*entity.Lead, error)
	DeleteLeadByEmail(email string) error
	UpdateLead(leads *entity.Lead) error
	AddTags(lead *entity.Lead, tags []string) error
	RemoveTags(lead *entity.Lead, tags []string) error
	CreateLead(lead *entity.Lead) (*entity.Lead, error)
}

type rdStation struct {
	client client.Client
}

func NewRDStation(clientID, clientSecret, refreshToken string) RDStation {
	secret := entity.Secret{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RefreshToken: refreshToken,
	}

	endpoint := oauth2.Endpoint{
		AuthURL:   "https://api.rd.services/auth/",
		TokenURL:  fmt.Sprintf("%s%s", RDURL, RefreshTokenURL),
		AuthStyle: oauth2.AuthStyleInParams,
	}

	cl, err := client.NewClient(secret, endpoint)
	if err != nil {
		return nil
	}

	return &rdStation{
		client: cl,
	}
}

func (rd rdStation) GetLeadByEmail(email string) (*entity.Lead, error) {
	ret, err := rd.client.Request(fmt.Sprintf("%s%semail:%s", RDURL, RDLeadPath, email), http.MethodGet, nil)
	if err != nil {
		return nil, err
	}

	var lead entity.Lead
	err = json.Unmarshal(ret, &lead)
	if err != nil {
		return nil, err
	}

	return &lead, nil
}

func (rd rdStation) CreateLead(lead *entity.Lead) (*entity.Lead, error) {
	data, err := json.Marshal(lead)
	if err != nil {
		return nil, err
	}

	data, err = rd.client.Request(fmt.Sprintf("%s%s", RDURL, RDLeadPath), http.MethodPost, data)
	if err != nil {
		return nil, err
	}
	var newLead entity.Lead
	err = json.Unmarshal(data, &newLead)
	if err != nil {
		return nil, err
	}

	return &newLead, nil
}

func (rd rdStation) DeleteLeadByEmail(email string) error {
	_, err := rd.client.Request(fmt.Sprintf("%s%semail:%s", RDURL, RDLeadPath, email), http.MethodDelete, nil)
	return err
}

func (rd rdStation) UpdateLead(lead *entity.Lead) error {
	data, err := json.Marshal(lead)
	if err != nil {
		return err
	}

	_, err = rd.client.Request(fmt.Sprintf("%s%semail:%s", RDURL, RDLeadPath, lead.Email), http.MethodPatch, data)

	return err
}

func (rd rdStation) AddTags(lead *entity.Lead, tags []string) error {
	for _, tag := range tags {
		if lead.HasTag(tag) {
			lead.Tags = removeFromList(lead.Tags, tag)
		}
	}

	if lead.HasTag(entity.AcademyTagCancelled) && contains(tags, entity.AcademyActive) {
		lead.Tags = removeFromList(lead.Tags, entity.AcademyTagCancelled)
	}
	if lead.HasTag(entity.AcademyActive) && contains(tags, entity.AcademyTagCancelled) {
		lead.Tags = removeFromList(lead.Tags, entity.AcademyActive)
	}

	lead.Tags = append(lead.Tags, tags...)

	data, err := json.Marshal(map[string][]string{"tags": lead.Tags})
	if err != nil {
		return err
	}

	_, err = rd.client.Request(fmt.Sprintf("%s%semail:%s", RDURL, RDLeadPath, lead.Email), http.MethodPatch, data)

	return err
}

func (rd rdStation) RemoveTags(lead *entity.Lead, tags []string) error {
	for _, tag := range tags {
		if lead.HasTag(tag) {
			lead.Tags = removeFromList(lead.Tags, tag)
		}
	}

	var (
		data []byte
		err  error
	)
	if lead.Tags != nil {
		data, err = json.Marshal(map[string][]string{"tags": lead.Tags})
		if err != nil {
			return err
		}
	} else {
		empty := []string{}
		data, err = json.Marshal(map[string][]string{"tags": empty})
		if err != nil {
			return err
		}
	}

	_, err = rd.client.Request(fmt.Sprintf("%s%semail:%s", RDURL, RDLeadPath, lead.Email), http.MethodPatch, data)

	return err
}

func contains(tags []string, tag string) bool {
	for _, t := range tags {
		if t == tag {
			return true
		}
	}

	return false
}

func removeFromList(tags []string, tag string) []string {
	j := 0
	for i, v := range tags {
		if v != tag {
			tags[j] = tags[i]
			j++
		}
	}

	return tags[:j]
}
