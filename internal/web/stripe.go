package web

import (
	"fmt"
	"net/http"
	"os"

	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/internal/gate"
	"github.com/bennycio/bundle/internal/logger"
	stripe "github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/account"
	"github.com/stripe/stripe-go/v72/accountlink"
	"github.com/stripe/stripe-go/v72/checkout/session"
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
		handleError(w, err, http.StatusInternalServerError)
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
			handleError(w, err, http.StatusInternalServerError)
			return
		}
		us.StripeId = acct.ID

		err = gs.UpdateUser(us)

		if err != nil {
			handleError(w, err, http.StatusInternalServerError)
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
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, acc.URL, http.StatusFound)

}

func stripeReturnHandlerFunc(w http.ResponseWriter, r *http.Request) {

	pro, err := getProfFromCookie(r)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	data := &templateData{
		Profile: pro,
	}

	stripe.Key = os.Getenv("STRIPE_TEST_KEY")

	gs := gate.NewGateService("", "")

	us, err := gs.GetUser(&api.User{Id: pro.Id})
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	acct, err := account.GetByID(us.Id, nil)
	if err != nil {
		handleError(w, err, http.StatusNotFound)
		return
	}
	data.Profile.StripeInfo.ChargesEnabled = acct.ChargesEnabled
	data.Profile.StripeInfo.DetailsSubmitted = acct.DetailsSubmitted

	err = tpl.ExecuteTemplate(w, "profile", data)
	if err != nil {
		logger.ErrLog.Panic(err.Error())
	}
}

func purchasePluginHandlerFunc(w http.ResponseWriter, r *http.Request) {

	pro, err := getProfFromCookie(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}

	data := templateData{
		Profile: pro,
	}

	switch r.Method {

	case http.MethodGet:

		r.ParseForm()

		plugin := r.FormValue("plugin")

		if plugin == "" {
			http.Redirect(w, r, "/plugins", http.StatusNotFound)
			return
		}

		gs := gate.NewGateService("", "")

		dbpl, err := gs.GetPlugin(&api.Plugin{Id: plugin})

		if err != nil {
			logger.ErrLog.Print(err.Error())
			handleError(w, err, http.StatusNotFound)
			return
		}

		if dbpl.Author == nil {
			err = fmt.Errorf("author is nil")
			logger.ErrLog.Print(err.Error())
			handleError(w, err, http.StatusNotFound)
			return
		}

		if dbpl.Author.StripeId == "" {
			err = fmt.Errorf("author is not striped up")
			logger.ErrLog.Print(err.Error())
			handleError(w, err, http.StatusNotFound)
			return
		}

		stripe.Key = os.Getenv("STRIPE_KEY")
		params := &stripe.CheckoutSessionParams{
			PaymentMethodTypes: stripe.StringSlice([]string{
				"card",
			}),
			LineItems: []*stripe.CheckoutSessionLineItemParams{
				{
					Name:     stripe.String(dbpl.Name),
					Amount:   stripe.Int64(int64(dbpl.Premium.Price)),
					Currency: stripe.String(string(stripe.CurrencyUSD)),
					Quantity: stripe.Int64(1),
				},
			},
			PaymentIntentData: &stripe.CheckoutSessionPaymentIntentDataParams{
				ApplicationFeeAmount: stripe.Int64(int64(dbpl.Premium.Price / 20)),
				TransferData: &stripe.CheckoutSessionPaymentIntentDataTransferDataParams{
					Destination: stripe.String(dbpl.Author.StripeId),
				},
			},
			Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
			SuccessURL: stripe.String(fmt.Sprintf("https://%splugins?plugin=%s", os.Getenv("WEB_HOST"), plugin)),
			CancelURL:  stripe.String(fmt.Sprintf("https://%splugins?plugin=%s", os.Getenv("WEB_HOST"), plugin)),
		}

		session, err := session.New(params)
		if err != nil {
			logger.ErrLog.Print(err.Error())
			handleError(w, err, http.StatusInternalServerError)
			return
		}
		data.PurchaseSession = session.ID

		err = tpl.ExecuteTemplate(w, "purchase", data)
		if err != nil {
			logger.ErrLog.Panic(err.Error())
		}
	}

}
