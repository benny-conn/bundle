package web

import (
	"net/http"
	"os"

	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/internal/gate"
	stripe "github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/account"
	"github.com/stripe/stripe-go/v72/accountlink"
)

func stripeAuthHandlerFunc(w http.ResponseWriter, r *http.Request) {

	pro, err := getProfFromCookie(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}

	gs := gate.NewGateService("", "")

	us, err := gs.GetUser(&api.User{Id: pro.Id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	stripe.Key = os.Getenv("STRIPE_TEST_KEY")

	if us.StripeId == "" {

		params := &stripe.AccountParams{
			Type:  stripe.String(string(stripe.AccountTypeExpress)),
			Email: stripe.String(pro.Email),
		}
		acct, err := account.New(params)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		us.StripeId = acct.ID

		err = gs.UpdateUser(us)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	p := &stripe.AccountLinkParams{
		Account:    stripe.String(us.StripeId),
		RefreshURL: stripe.String("https://" + os.Getenv("WEB_HOST") + "/stripe/auth"),
		ReturnURL:  stripe.String("https://" + os.Getenv("WEB_HOST") + "/stripe/return"),
		Type:       stripe.String("account_onboarding"),
	}
	acc, err := accountlink.New(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, acc.URL, http.StatusFound)

}

func stripeReturnHandlerFunc(w http.ResponseWriter, r *http.Request) {

	pro, err := getProfFromCookie(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}

	data := &templateData{
		Profile: pro,
	}

	stripe.Key = os.Getenv("STRIPE_TEST_KEY")

	gs := gate.NewGateService("", "")

	us, err := gs.GetUser(&api.User{Id: pro.Id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	acct, err := account.GetByID(us.Id, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	data.Profile.StripeInfo.ChargesEnabled = acct.ChargesEnabled
	data.Profile.StripeInfo.DetailsSubmitted = acct.DetailsSubmitted

	err = tpl.ExecuteTemplate(w, "profile", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
