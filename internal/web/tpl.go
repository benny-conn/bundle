package web

import (
	"fmt"
	"strings"
	"text/template"
	"time"

	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/internal"
)

type profile struct {
	Id               string
	Username         string
	Email            string
	Tag              string
	StripeInfo       userStripeInfo
	PurchasedPlugins []string
}

type userStripeInfo struct {
	Id               string
	DetailsSubmitted bool
	ChargesEnabled   bool
}

type templateData struct {
	Profile         profile
	Plugin          *api.Plugin
	Plugins         []*api.Plugin
	PurchaseSession string
	Page            int
	Functions       functions
	Readme          string
	Info            string
	Error           errorData
	Referrer        string
}

type functions struct {
	Math     func(int, int, string) int
	Date     func(int64) string
	Contains func([]string, string) bool
	FmtMoney func(int) string
}

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("assets/templates/*.gohtml"))
}

func fillFunctions(data templateData) templateData {
	math := func(i int, p int, op string) int {
		result := 0
		switch op {
		case "add":
			result = i + p
		case "sub":
			result = i - p
		case "mul":
			result = i * p
		case "div":
			result = i / p
		}
		return result
	}
	date := func(i int64) string {
		d := time.Unix(i, 0)
		return d.Format("Mon Jan 2 15:04:05 MST")
	}
	contains := internal.Contains

	money := func(i int) string {
		new := float64(i) / 100
		s := fmt.Sprintf("$%f", new)
		return strings.TrimRight(s, "0")
	}

	data.Functions = functions{
		Math:     math,
		Date:     date,
		Contains: contains,
		FmtMoney: money,
	}
	return data

}

func userToProfile(u *api.User) profile {
	p := profile{
		Id:       u.Id,
		Username: u.Username,
		Email:    u.Email,
		Tag:      u.Tag,
		StripeInfo: userStripeInfo{
			Id: u.StripeId,
		},
	}
	if u.Purchases != nil {
		s := make([]string, len(u.Purchases))
		for i, v := range u.Purchases {
			s[i] = v.ObjectId
		}
		p.PurchasedPlugins = s
	}
	return p
}
