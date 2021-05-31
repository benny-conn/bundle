package web

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/internal/gate"
	"github.com/bennycio/bundle/logger"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

const perPageCount = 15

func pluginsHandlerFunc(w http.ResponseWriter, req *http.Request) {

	data := templateData{}
	user, err := getProfFromCookie(req)

	if err == nil {
		data.Profile = user
	}

	gs := gate.NewGateService("", "")

	switch req.Method {

	case http.MethodGet:
		err := req.ParseForm()
		if err != nil {
			logger.ErrLog.Print(err.Error())
			handleError(w, err, http.StatusBadRequest)
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
				logger.ErrLog.Print(err.Error())
				handleError(w, err, http.StatusNotFound)
				return
			}
			sortNumber, err := strconv.Atoi(sort)
			if err != nil {
				logger.ErrLog.Print(err.Error())
				handleError(w, err, http.StatusNotFound)
				return
			}
			categoryNum, err := strconv.Atoi(category)
			if err != nil {
				logger.ErrLog.Print(err.Error())
				handleError(w, err, http.StatusNotFound)
				return
			}

			data.Page = pageNumber

			req := &api.PaginatePluginsRequest{
				Count:    perPageCount,
				Search:   search,
				Page:     int32(pageNumber),
				Sort:     api.Sort(sortNumber),
				Category: api.Category(categoryNum),
			}

			plugins, err := gs.PaginatePlugins(req)
			if err != nil {
				logger.ErrLog.Print(err.Error())
				handleError(w, err, http.StatusNotFound)
				return
			}

			data.Plugins = plugins

		} else {
			req := &api.Plugin{
				Name: pluginName,
			}
			plugin, err := gs.GetPlugin(req)
			if err != nil {
				logger.ErrLog.Print(err.Error())
				handleError(w, err, http.StatusNotFound)
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

	data = fillFunctions(data)
	err = tpl.ExecuteTemplate(w, "plugins", data)
	if err != nil {
		logger.ErrLog.Print(err.Error())
	}

}

func thumbnailHandlerFunc(w http.ResponseWriter, req *http.Request) {

	prof, err := getProfFromCookie(req)
	if err != nil {
		logger.ErrLog.Print(err.Error())
		handleError(w, err, http.StatusUnauthorized)
	}
	if req.Method != http.MethodPost {
		err := fmt.Errorf("only method post allowed")
		logger.ErrLog.Print(err.Error())
		handleError(w, err, http.StatusBadRequest)
		return
	}
	gs := gate.NewGateService("", "")
	err = req.ParseMultipartForm(32 << 20)
	if err != nil {
		logger.ErrLog.Print(err.Error())
		handleError(w, err, http.StatusBadRequest)
		return
	}
	plugin := req.FormValue("plugin")

	thumbnail, h, err := req.FormFile("thumbnail")
	if err != nil {
		logger.ErrLog.Print(err.Error())
		handleError(w, err, http.StatusBadRequest)
		return
	}

	if h.Size > (1 << 20) {
		err = fmt.Errorf("file too large")
		logger.ErrLog.Print(err.Error())
		handleError(w, err, http.StatusBadRequest)
		return
	}

	if prof.Id == "" {
		err = fmt.Errorf("no user specified")
		logger.ErrLog.Print(err.Error())
		handleError(w, err, http.StatusBadRequest)
		return
	}
	if plugin == "" {
		err = fmt.Errorf("no plugin specified")
		logger.ErrLog.Print(err.Error())
		handleError(w, err, http.StatusBadRequest)
		return
	}

	u := &api.User{
		Id: prof.Id,
	}
	p := &api.Plugin{
		Id: plugin,
	}

	dbpl, err := gs.GetPlugin(p)
	if err != nil {
		logger.ErrLog.Print(err.Error())
		handleError(w, err, http.StatusBadRequest)
		return
	}
	if dbpl.Author.Id != u.Id {
		err = fmt.Errorf("must be plugin author")
		logger.ErrLog.Print(err.Error())
		handleError(w, err, http.StatusBadRequest)
		return
	}

	err = gs.UploadThumbnail(u, p, thumbnail)
	if err != nil {
		logger.ErrLog.Print(err.Error())
		handleError(w, err, http.StatusBadRequest)
		return
	}

	http.Redirect(w, req, req.Header.Get("Referer"), http.StatusFound)
}

func premiumHandlerFunc(w http.ResponseWriter, req *http.Request) {

	prof, err := getProfFromCookie(req)
	if err != nil {
		logger.ErrLog.Print(err.Error())
		handleError(w, err, http.StatusUnauthorized)
	}

	if req.Method != http.MethodPost {
		err := fmt.Errorf("only method post allowed")
		logger.ErrLog.Print(err.Error())
		handleError(w, err, http.StatusBadRequest)
		return
	}
	gs := gate.NewGateService("", "")
	err = req.ParseForm()
	if err != nil {
		logger.ErrLog.Print(err.Error())
		handleError(w, err, http.StatusBadRequest)
		return
	}
	plugin := req.FormValue("plugin")
	price := req.FormValue("price")

	if prof.Id == "" || plugin == "" || price == "" {
		err = fmt.Errorf("user, plugin, and price must be specified")
		logger.ErrLog.Print(err.Error())
		handleError(w, err, http.StatusBadRequest)
		return
	}

	spl := strings.Split(price, ".")

	finalPrice := price
	if len(spl) > 1 {
		dec := spl[1]
		if len(dec) > 2 {
			dec = dec[:2]
		}
		finalPrice = strings.Join([]string{spl[0], dec}, ".")
	}

	priceNum, err := strconv.ParseFloat(finalPrice, 64)
	if err != nil {
		p, err := strconv.Atoi(price)
		if err != nil {
			logger.ErrLog.Print(err.Error())
			handleError(w, err, http.StatusBadRequest)
			return
		}
		priceNum = float64(p)
	}

	prem := &api.Premium{
		Price: int32(priceNum * 100),
	}

	p := &api.Plugin{
		Id: plugin,
	}
	dbpl, err := gs.GetPlugin(p)
	if err != nil {
		logger.ErrLog.Print(err.Error())
		handleError(w, err, http.StatusBadRequest)
		return
	}
	if dbpl.Author.Id != prof.Id {
		err = fmt.Errorf("must be plugin author")
		logger.ErrLog.Print(err.Error())
		handleError(w, err, http.StatusBadRequest)
		return
	}

	p.Premium = prem

	err = gs.UpdatePlugin(p)
	if err != nil {
		logger.ErrLog.Print(err.Error())
		handleError(w, err, http.StatusBadRequest)
		return
	}

	if dbpl.Author.StripeId == "" {
		http.Redirect(w, req, "/stripe/auth", http.StatusFound)
	} else {
		http.Redirect(w, req, req.Header.Get("Referer"), http.StatusFound)
	}
}
