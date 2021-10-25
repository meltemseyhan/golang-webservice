package receipt

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/meltemseyhan/inventoryservice/cors"
)

func receiptsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		receiptList, err := GetReceipts()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		j, err := json.Marshal(receiptList)
		if err != nil {
			log.Fatal(err)
		}
		_, err = w.Write(j)
		if err != nil {
			log.Fatal(err)
		}
	case http.MethodPost:
		r.ParseMultipartForm(5 << 20)                //5Mb
		uFile, uHeader, err := r.FormFile("receipt") //Uploaded file
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer uFile.Close()

		//File on server disk
		f, err := os.OpenFile(filepath.Join(ReceiptDirectory, uHeader.Filename), os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer f.Close()

		// Copy uploaded file to disk
		io.Copy(f, uFile)
		w.WriteHeader(http.StatusCreated)
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	urlPathSegments := strings.Split(r.URL.Path, fmt.Sprintf("%s/", receiptBasePath))
	if len(urlPathSegments) > 2 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fileName := urlPathSegments[1]
	file, err := os.Open(filepath.Join(ReceiptDirectory, fileName))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer file.Close()
	fHeader := make([]byte, 512)
	file.Read(fHeader)
	fContentType := http.DetectContentType(fHeader)
	fInfo, err := file.Stat()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	w.Header().Set("Content-Type", fContentType)
	w.Header().Set("Content-Length", strconv.FormatInt(fInfo.Size(), 10))
	file.Seek(0, 0)
	io.Copy(w, file)
}

const receiptBasePath = "receipts"

func SetupRoutes(apiBasePath string) {
	http.Handle(fmt.Sprintf("%s/%s", apiBasePath, receiptBasePath), cors.Middleware(http.HandlerFunc(receiptsHandler)))
	http.Handle(fmt.Sprintf("%s/%s/", apiBasePath, receiptBasePath), cors.Middleware(http.HandlerFunc(downloadHandler)))
}
