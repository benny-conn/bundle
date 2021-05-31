package web

import (
	"fmt"
	"net/http"
	"os"

	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/internal/gate"
	"github.com/bennycio/bundle/logger"
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
		logger.ErrLog.Print(err.Error())
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	stripe.Key = os.Getenv("STRIPE_KEY")

	if us.StripeId == "" {

		params := &stripe.AccountParams{
			Type:         stripe.String(string(stripe.AccountTypeExpress)),
			Email:        stripe.String(pro.Email),
			BusinessType: stripe.String("individual"),
		}
		acct, err := account.New(params)

		if err != nil {
			logger.ErrLog.Print(err.Error())
			handleError(w, err, http.StatusInternalServerError)
			return
		}
		us.StripeId = acct.ID
		err = gs.UpdateUser(us)

		if err != nil {
			logger.ErrLog.Print(err.Error())
			handleError(w, err, http.StatusInternalServerError)
			return
		}
	}

	p := &stripe.AccountLinkParams{
		Account:    stripe.String(us.StripeId),
		RefreshURL: stripe.String(fmt.Sprintf("https://%s:%s/stripe/auth", os.Getenv("WEB_HOST"), os.Getenv("WEB_PORT"))),
		ReturnURL:  stripe.String(fmt.Sprintf("https://%s:%s/stripe/return", os.Getenv("WEB_HOST"), os.Getenv("WEB_PORT"))),
		Type:       stripe.String("account_onboarding"),
	}
	acc, err := accountlink.New(p)
	if err != nil {
		logger.ErrLog.Print(err.Error())
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, acc.URL, http.StatusFound)

}

func stripeReturnHandlerFunc(w http.ResponseWriter, r *http.Request) {

	pro, err := getProfFromCookie(r)
	if err != nil {
		logger.ErrLog.Print(err.Error())
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	data := &templateData{
		Profile: pro,
	}

	stripe.Key = os.Getenv("STRIPE_KEY")

	gs := gate.NewGateService("", "")

	us, err := gs.GetUser(&api.User{Id: pro.Id})
	if err != nil {
		logger.ErrLog.Print(err.Error())
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	acct, err := account.GetByID(us.StripeId, nil)
	if err != nil {
		logger.ErrLog.Print(err.Error())
		handleError(w, err, http.StatusNotFound)
		return
	}
	data.Profile.StripeInfo.ChargesEnabled = acct.ChargesEnabled
	data.Profile.StripeInfo.DetailsSubmitted = acct.DetailsSubmitted

	fmt.Println(acct)

	err = tpl.ExecuteTemplate(w, "profile", data)
	if err != nil {
		logger.ErrLog.Print(err.Error())
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

		data.Plugin = dbpl

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

		dbUser, err := gs.GetUser(&api.User{Id: pro.Id})

		if err != nil {
			logger.ErrLog.Print(err.Error())
			handleError(w, err, http.StatusInternalServerError)
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
			CustomerEmail: stripe.String(dbUser.Email),
			PaymentIntentData: &stripe.CheckoutSessionPaymentIntentDataParams{
				ApplicationFeeAmount: stripe.Int64(int64((dbpl.Premium.Price / 15))),
				TransferData: &stripe.CheckoutSessionPaymentIntentDataTransferDataParams{
					Destination: stripe.String(dbpl.Author.StripeId),
				},
			},
			Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
			SuccessURL: stripe.String(fmt.Sprintf("https://%s:%s/plugins?plugin=%s", os.Getenv("WEB_HOST"), os.Getenv("WEB_PORT"), dbpl.Name)),
			CancelURL:  stripe.String(fmt.Sprintf("https://%s:%s/plugins?plugin=%s", os.Getenv("WEB_HOST"), os.Getenv("WEB_PORT"), dbpl.Name)),
		}

		session, err := session.New(params)
		if err != nil {
			logger.ErrLog.Print(err.Error())
			handleError(w, err, http.StatusInternalServerError)
			return
		}
		data.PurchaseSession = session.ID

		pur := &api.Purchase{
			ObjectId: dbpl.Id,
			Complete: false,
		}

		if dbUser.Purchases != nil {
			dbUser.Purchases = append(dbUser.Purchases, pur)
		} else {
			dbUser.Purchases = []*api.Purchase{pur}
		}

		err = gs.UpdateUser(dbUser)
		if err != nil {
			logger.ErrLog.Print(err.Error())
			handleError(w, err, http.StatusInternalServerError)
			return
		}

		data = fillFunctions(data)
		err = tpl.ExecuteTemplate(w, "purchase", data)
		if err != nil {
			logger.ErrLog.Print(err.Error())
		}
	}

}
