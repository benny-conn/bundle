package web

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/internal"
	"github.com/bennycio/bundle/internal/gate"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/checkout/session"
)

const perPageCount = 15

func pluginsHandlerFunc(w http.ResponseWriter, req *http.Request) {

	data := templateData{}
	user, err := getProfFromCookie(req)

	if err == nil {
		data.Profile = user
	}

	gs := gate.NewGateService("", "")

	if user.Id != "" {
		bu, err := gs.GetBundle(&api.Bundle{UserId: user.Id})
		if err == nil {
			data.Bundle = bu
		}
	}

	data.Contains = internal.Contains

	switch req.Method {

	case http.MethodGet:
		err := req.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		pluginName := req.FormValue("plugin")

		if pluginName == "" {
			page := req.FormValue("page")
			search := req.FormValue("search")
			sort := req.FormValue("sort")
			category := req.FormValue("category")

			if page == "" {
				page = "1"
			}
			if sort == "" {
				sort = "0"
			}
			if category == "" {
				category = "0"
			}

			pageNumber, err := strconv.Atoi(page)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			sortNumber, err := strconv.Atoi(sort)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			categoryNum, err := strconv.Atoi(category)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}

			data.Page = pageNumber

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
			data.Math = math

			date := func(i int64) string {
				d := time.Unix(i, 0)
				return d.Format("Mon Jan 2 15:04:05 MST")
			}
			data.Date = date

			req := &api.PaginatePluginsRequest{
				Count:    perPageCount,
				Search:   search,
				Page:     int32(pageNumber),
				Sort:     api.Sort(sortNumber),
				Category: api.Category(categoryNum),
			}

			plugins, err := gs.PaginatePlugins(req)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}

			data.Plugins = plugins

		} else {
			req := &api.Plugin{
				Name: pluginName,
			}
			plugin, err := gs.GetPlugin(req)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}

			readme, err := gs.GetReadme(req)

			if err == nil {
				output := blackfriday.Run([]byte(readme.Text))
				data.Readme = string(bluemonday.UGCPolicy().SanitizeBytes(output))
			}

			data.Plugin = plugin
		}
	}

	err = tpl.ExecuteTemplate(w, "plugins", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func thumbnailHandlerFunc(w http.ResponseWriter, req *http.Request) {

	if req.Method != http.MethodPost {
		http.Redirect(w, req, req.Header.Get("Referer"), http.StatusFound)
		return
	}
	gs := gate.NewGateService("", "")
	err := req.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user := req.FormValue("user")
	plugin := req.FormValue("plugin")

	thumbnail, h, err := req.FormFile("thumbnail")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if h.Size > (1 << 20) {
		http.Error(w, "file too large", http.StatusBadRequest)
		return
	}

	if user == "" {
		http.Error(w, "no user specified", http.StatusBadRequest)
		return
	}
	if plugin == "" {
		http.Error(w, "no user specified", http.StatusBadRequest)
		return
	}

	u := &api.User{
		Id: user,
	}
	p := &api.Plugin{
		Id: plugin,
	}

	err = gs.UploadThumbnail(u, p, thumbnail)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, req, req.Header.Get("Referer"), http.StatusFound)
}

func bundlerHandlerFunc(w http.ResponseWriter, req *http.Request) {

	if req.Method != http.MethodPost {
		http.Redirect(w, req, req.Header.Get("Referer"), http.StatusFound)
		return
	}
	req.ParseForm()

	user, err := getProfFromCookie(req)
	if err != nil {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}

	gs := gate.NewGateService("", "")

	r := &api.Bundle{UserId: user.Id}
	bu, err := gs.GetBundle(&api.Bundle{UserId: user.Id})
	if err != nil {
		bu = r
	}

	plname := req.FormValue("plugin")

	if plname == "" {
		http.Error(w, "no plugin name found", http.StatusBadRequest)
		return
	}

	bu.Plugins = append(bu.Plugins, plname)

	err = gs.UpdateBundle(bu)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, req, req.Header.Get("Referer"), http.StatusFound)
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
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		if dbpl.Author == nil {
			http.Error(w, "author is nil", http.StatusNotFound)
			return
		}

		if dbpl.Author.StripeId == "" {
			http.Error(w, "author is not striped up", http.StatusNotFound)
			return
		}

		stripe.Key = os.Getenv("STRIPE_TEST_KEY")
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
			SuccessURL: stripe.String("https://" + os.Getenv("WEB_HOST") + "plugins/payment_success"),
			CancelURL:  stripe.String("https://" + os.Getenv("WEB_HOST") + "plugins/payment_cancel"),
		}

		session, err := session.New(params)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data.PurchaseSession = session.ID

		err = tpl.ExecuteTemplate(w, "purchase", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

}

func paymentSuccessHandlerFunc(w http.ResponseWriter, r *http.Request) {

}
