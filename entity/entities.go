package entity

import (
	"time"

	"golang.org/x/oauth2"
)

const (
	// TAGS
	AcademyTagCancelled = "exac"
	AcademyActive       = "acativo"
	EmailOptOut         = "descadastrado"
)

type Lead struct {
	Uuid          string   `json:"uuid,omitempty"`
	Name          string   `json:"name,omitempty"`
	Email         string   `json:"email,omitempty"`
	JobTitle      string   `json:"job_title,omitempty"`
	Bio           string   `json:"bio,omitempty"`
	Website       string   `json:"website,omitempty"`
	PersonalPhone string   `json:"personal_phone,omitempty"`
	MobilePhone   string   `json:"mobile_phone,omitempty"`
	City          string   `json:"city,omitempty"`
	State         string   `json:"state,omitempty"`
	Country       string   `json:"country,omitempty"`
	Twitter       string   `json:"twitter,omitempty"`
	Facebook      string   `json:"facebook,omitempty"`
	Linkedin      string   `json:"linkedin,omitempty"`
	Tags          []string `json:"tags,omitempty"`
	ExtraEmails   []string `json:"extra_emails,omitempty"`
}

func (l *Lead) Empty() bool {
	return l == nil ||
		l.Uuid == "" &&
			l.Name == "" &&
			l.Email == "" &&
			l.JobTitle == "" &&
			l.Bio == "" &&
			l.Website == "" &&
			l.PersonalPhone == "" &&
			l.MobilePhone == "" &&
			l.City == "" &&
			l.State == "" &&
			l.Country == "" &&
			l.Twitter == "" &&
			l.Facebook == "" &&
			l.Linkedin == "" &&
			l.Tags == nil &&
			l.ExtraEmails == nil
}

func (l *Lead) HasTag(tag string) bool {
	for _, t := range l.Tags {
		if t == tag {
			return true
		}
	}

	return false
}

type Secret struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RefreshToken string `json:"refresh_token"`
}

type Token struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	CreationDate string `json:"creation_date"`
}

func (t *Token) Auth2Token() *oauth2.Token {
	return &oauth2.Token{
		AccessToken:  t.AccessToken,
		RefreshToken: t.RefreshToken,
		Expiry:       time.Now().Add(time.Second * time.Duration(t.ExpiresIn)),
	}
}
