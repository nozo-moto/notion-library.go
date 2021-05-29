package main

import (
	"context"
	"fmt"
	"os"

	"github.com/nozo-moto/notion-library/google_isbn"
	"github.com/nozo-moto/notion-library/notion"
)

func main() {
	notionAccessToken := os.Getenv("NOTION_ACCESS_TOKEN")
	if notionAccessToken == "" {
		panic("set NOTION_ACCESS_TOKEN")
	}
	notionPageId := os.Getenv("NOTION_PAGE_ID")
	if notionPageId == "" {
		panic("set NOTION_PAGE_ID")
	}
	if len(os.Args) < 2 {
		panic("set isbn")
	}
	isbn := os.Args[1]
	googleISBN := google_isbn.New()
	notion := notion.New(notionAccessToken)

	bookInfo, err := googleISBN.GetInfo(context.Background(), isbn)
	if err != nil {
		panic(err)
	}

	bookInfoVolume := bookInfo.Items[0].Volumeinfo
	if err := notion.PostToDB(
		context.Background(),
		notion.NewBookShelf(
			notionPageId,
			isbn,
			bookInfoVolume.Title,
			fmt.Sprint(bookInfoVolume.Authors),
			bookInfoVolume.Publisheddate,
		),
	); err != nil {
		panic(err)
	}
}
