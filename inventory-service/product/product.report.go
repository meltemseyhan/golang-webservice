package product

import (
	"encoding/json"
	"log"
	"net/http"
	"path"
	"text/template"
)

type ProductReportFilter struct {
	NameFilter         string `json:"productName"`
	ManufacturerFilter string `json:"manufacturer"`
	SKUFilter          string `json:"sku"`
}

func productReportHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var productFilter ProductReportFilter
		err := json.NewDecoder(r.Body).Decode(&productFilter)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		products, err := searchForProductData(&productFilter)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if len(products) == 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		t := template.New("report.gotmpl").Funcs(template.FuncMap{"mod": func(i, j int) bool { return i%j == 0 }})
		t, err = t.ParseFiles(path.Join("templates", "report.gotmpl"))
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = t.Execute(w, products)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
