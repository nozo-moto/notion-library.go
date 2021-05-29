package google_isbn

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	BaseURL = "https://www.googleapis.com/books/v1/volumes?q=isbn:"
)

type GoogleISBN struct {
}

func New() *GoogleISBN {
	return &GoogleISBN{}
}

func (g *GoogleISBN) GetInfo(isbn string) (bookInfo *BookInfo, err error) {
	resp, err := http.Get(fmt.Sprintf("%s%s", BaseURL, isbn))
	if err != nil {
		return nil, err
	}
	resBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(resBytes, &bookInfo); err != nil {
		return nil, err
	}

	return
}

type BookInfo struct {
	Kind       string `json:"kind"`
	Totalitems int    `json:"totalItems"`
	Items      []struct {
		Kind       string `json:"kind"`
		ID         string `json:"id"`
		Etag       string `json:"etag"`
		Selflink   string `json:"selfLink"`
		Volumeinfo struct {
			Title               string   `json:"title"`
			Subtitle            string   `json:"subtitle"`
			Authors             []string `json:"authors"`
			Publisheddate       string   `json:"publishedDate"`
			Description         string   `json:"description"`
			Industryidentifiers []struct {
				Type       string `json:"type"`
				Identifier string `json:"identifier"`
			} `json:"industryIdentifiers"`
			Readingmodes struct {
				Text  bool `json:"text"`
				Image bool `json:"image"`
			} `json:"readingModes"`
			Pagecount           int    `json:"pageCount"`
			Printtype           string `json:"printType"`
			Maturityrating      string `json:"maturityRating"`
			Allowanonlogging    bool   `json:"allowAnonLogging"`
			Contentversion      string `json:"contentVersion"`
			Panelizationsummary struct {
				Containsepubbubbles  bool `json:"containsEpubBubbles"`
				Containsimagebubbles bool `json:"containsImageBubbles"`
			} `json:"panelizationSummary"`
			Imagelinks struct {
				Smallthumbnail string `json:"smallThumbnail"`
				Thumbnail      string `json:"thumbnail"`
			} `json:"imageLinks"`
			Language            string `json:"language"`
			Previewlink         string `json:"previewLink"`
			Infolink            string `json:"infoLink"`
			Canonicalvolumelink string `json:"canonicalVolumeLink"`
		} `json:"volumeInfo"`
		Saleinfo struct {
			Country     string `json:"country"`
			Saleability string `json:"saleability"`
			Isebook     bool   `json:"isEbook"`
		} `json:"saleInfo"`
		Accessinfo struct {
			Country                string `json:"country"`
			Viewability            string `json:"viewability"`
			Embeddable             bool   `json:"embeddable"`
			Publicdomain           bool   `json:"publicDomain"`
			Texttospeechpermission string `json:"textToSpeechPermission"`
			Epub                   struct {
				Isavailable bool `json:"isAvailable"`
			} `json:"epub"`
			Pdf struct {
				Isavailable bool `json:"isAvailable"`
			} `json:"pdf"`
			Webreaderlink       string `json:"webReaderLink"`
			Accessviewstatus    string `json:"accessViewStatus"`
			Quotesharingallowed bool   `json:"quoteSharingAllowed"`
		} `json:"accessInfo"`
		Searchinfo struct {
			Textsnippet string `json:"textSnippet"`
		} `json:"searchInfo"`
	} `json:"items"`
}
