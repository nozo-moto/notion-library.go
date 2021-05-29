package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/nozo-moto/notion-library/google_isbn"
	"github.com/nozo-moto/notion-library/notion"
	"go.uber.org/zap"
)

func main() {
	notionAccessToken := os.Getenv("NOTION_ACCESS_TOKEN")
	if notionAccessToken == "" {
		panic("set NOTION_ACCESS_TOKEN")
	}

	googleISBN := google_isbn.New()
	notion := notion.New(notionAccessToken)

	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		isbn := r.URL.Query().Get("isbn")
		if isbn == "" {
			w.WriteHeader(http.StatusBadRequest)
			logger.Info("bad request no isbn")
			return
		}

		notionPageId := r.URL.Query().Get("pageid")
		if isbn == "" {
			w.WriteHeader(http.StatusBadRequest)
			logger.Info("bad request no notion pageid")
			return
		}
		bookInfo, err := googleISBN.GetInfo(isbn)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Info("google isbn get info error", zap.String("isbn", isbn), zap.Error(err))
			return
		}

		bookInfoVolume := bookInfo.Items[0].Volumeinfo
		if err := notion.PostToDB(
			notion.NewBookShelf(
				notionPageId,
				isbn,
				bookInfoVolume.Title,
				fmt.Sprint(bookInfoVolume.Authors),
				bookInfoVolume.Publisheddate,
			),
		); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Info("notion post error", zap.String("isbn", isbn), zap.Error(err))
			return
		}

		w.WriteHeader(http.StatusCreated)
		return
	})
	http.ListenAndServe(":8080", nil)
}
