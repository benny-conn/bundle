package gate

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/internal/gate/grpc"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/webhook"
)

func checkoutCompleteHandlerFunc(w http.ResponseWriter, req *http.Request) {
	const MaxBodyBytes = int64(65536)
	req.Body = http.MaxBytesReader(w, req.Body, MaxBodyBytes)
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
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
		http.Error(w, err.Error(), http.StatusBadRequest) // Return a 400 error on a bad signature.
		return
	}

	if event.Type == "checkout.session.completed" {
		var session stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &session)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = handleCompletedCheckoutSession(session)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
}

func handleCompletedCheckoutSession(ses stripe.CheckoutSession) error {
	userClient := grpc.NewUserClient("", "")
	plClient := grpc.NewPluginClient("", "")
	cu := ses.CustomerEmail
	dbus, err := userClient.Get(&api.User{Email: cu})
	if err != nil {
		return err
	}

	item := ses.Object

	dbpl, err := plClient.Get(&api.Plugin{Name: item})
	if err != nil {
		return err
	}

	if dbus.Purchases == nil {
		return fmt.Errorf("user does not have any pending or non-pending purchases")
	}
	for _, v := range dbus.Purchases {
		if v.ObjectId == dbpl.Id {
			v.Complete = true
		}
	}

	err = userClient.Update(dbus)
	if err != nil {
		return err
	}
	return nil
}
