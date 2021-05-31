package gate

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/bennycio/bundle/internal"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/webhook"
)

func checkoutCompleteHandlerFunc(w http.ResponseWriter, req *http.Request) {
	const MaxBodyBytes = int64(65536)
	req.Body = http.MaxBytesReader(w, req.Body, MaxBodyBytes)
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		internal.HttpError(w, err, http.StatusServiceUnavailable)
		return
	}

	// If you are testing your webhook locally with the Stripe CLI you
	// can find the endpoint's secret by running `stripe listen`
	// Otherwise, find your endpoint's secret in your webhook settings
	// in the Developer Dashboard
	endpointSecret := os.Getenv("STRIPE_WBH_SECRET")

	// Verify webhook signature and extract the event.
	// See https://stripe.com/docs/webhooks/signatures for more information.
	event, err := webhook.ConstructEvent(body, req.Header.Get("Stripe-Signature"), endpointSecret)

	if err != nil {
		internal.HttpError(w, err, http.StatusBadRequest) // Return a 400 error on a bad signature.
		return
	}

	if event.Type == "checkout.session.completed" {
		var session stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &session)
		if err != nil {
			internal.HttpError(w, err, http.StatusBadRequest)
			return
		}
		handleCompletedCheckoutSession(session)
	}
}

func handleCompletedCheckoutSession(ses stripe.CheckoutSession) {

}
