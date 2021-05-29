package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/nozo-moto/notion-library/google_isbn"
	"github.com/nozo-moto/notion-library/notion"
	_ "go.uber.org/automaxprocs"
	"go.uber.org/zap"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

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

	handler := &Handler{logger: logger, googleISBN: googleISBN, notion: notion}
	server := &http.Server{Addr: ":8080", Handler: handler}
	go func() {
		<-ctx.Done()
		server.Shutdown(ctx)
	}()

	logger.Info("listen and serve", zap.Error(server.ListenAndServe()))
}

type Handler struct {
	logger     *zap.Logger
	googleISBN *google_isbn.GoogleISBN
	notion     *notion.Notion
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	isbn := r.URL.Query().Get("isbn")
	if isbn == "" {
		w.WriteHeader(http.StatusBadRequest)
		h.logger.Info("bad request no isbn")
		return
	}

	notionPageId := r.URL.Query().Get("pageid")
	if isbn == "" {
		w.WriteHeader(http.StatusBadRequest)
		h.logger.Info("bad request no notion pageid")
		return
	}
	bookInfo, err := h.googleISBN.GetInfo(ctx, isbn)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Info("google isbn get info error", zap.String("isbn", isbn), zap.Error(err))
		return
	}

	bookInfoVolume := bookInfo.Items[0].Volumeinfo
	if err := h.notion.PostToDB(
		ctx,
		h.notion.NewBookShelf(
			notionPageId,
			isbn,
			bookInfoVolume.Title,
			fmt.Sprint(bookInfoVolume.Authors),
			bookInfoVolume.Publisheddate,
		),
	); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Info("notion post error", zap.String("isbn", isbn), zap.Error(err))
		return
	}

	w.WriteHeader(http.StatusCreated)
	return
}
