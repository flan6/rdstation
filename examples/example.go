package main

import (
	"fmt"

	"github.com/flan6/rdstation"
)

func main() {
	ClientID := ""
	ClientSecret := ""
	RefreshToken := ""

	rd := rdstation.NewRDStation(ClientID, ClientSecret, RefreshToken)
	lead, err := rd.GetLeadByEmail("b@qual.work")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(lead)
}
