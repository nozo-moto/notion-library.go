package notion

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Notion struct {
	accesstoken string
}

const (
	BASE_URL = "https://api.notion.com/v1/"
)

func New(accesstoken string) *Notion {
	return &Notion{
		accesstoken,
	}
}

func (n *Notion) NewBookShelf(databaseID, isbn, title, author, publishedDate string) *BookShelf {
	return &BookShelf{
		Parent: struct {
			DatabaseID string "json:\"database_id\""
		}{
			DatabaseID: databaseID,
		},
		Properties: struct {
			Title []struct {
				Text struct {
					Content string "json:\"content\""
				} "json:\"text\""
			} "json:\"Title\""
			Author []struct {
				Text struct {
					Content string "json:\"content\""
				} "json:\"text\""
			} "json:\"Author\""
			Publisheddate []struct {
				Text struct {
					Content string "json:\"content\""
				} "json:\"text\""
			} "json:\"PublishedDate\""
			ISBN []struct {
				Text struct {
					Content string "json:\"content\""
				} "json:\"text\""
			} "json:\"ISBN\""
		}{
			Title: []struct {
				Text struct {
					Content string "json:\"content\""
				} "json:\"text\""
			}{
				{
					Text: struct {
						Content string "json:\"content\""
					}{
						Content: title,
					},
				},
			},
			Author: []struct {
				Text struct {
					Content string "json:\"content\""
				} "json:\"text\""
			}{
				{
					Text: struct {
						Content string "json:\"content\""
					}{
						Content: author,
					},
				},
			},
			Publisheddate: []struct {
				Text struct {
					Content string "json:\"content\""
				} "json:\"text\""
			}{
				{
					Text: struct {
						Content string "json:\"content\""
					}{
						Content: publishedDate,
					},
				},
			},
			ISBN: []struct {
				Text struct {
					Content string "json:\"content\""
				} "json:\"text\""
			}{
				{
					Text: struct {
						Content string "json:\"content\""
					}{
						Content: isbn,
					},
				},
			},
		},
	}
}

type BookShelf struct {
	Parent struct {
		DatabaseID string `json:"database_id"`
	} `json:"parent"`
	Properties struct {
		Title []struct {
			Text struct {
				Content string `json:"content"`
			} `json:"text"`
		} `json:"Title"`
		Author []struct {
			Text struct {
				Content string `json:"content"`
			} `json:"text"`
		} `json:"Author"`
		Publisheddate []struct {
			Text struct {
				Content string `json:"content"`
			} `json:"text"`
		} `json:"PublishedDate"`
		ISBN []struct {
			Text struct {
				Content string `json:"content"`
			} `json:"text"`
		} `json:"ISBN"`
	} `json:"properties"`
}

func (n *Notion) PostToDB(bookShelfInfo *BookShelf) error {
	bookShelfBytes, err := json.Marshal(bookShelfInfo)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s%s", BASE_URL, "pages"),
		bytes.NewBuffer([]byte(bookShelfBytes)),
	)
	if err != nil {
		return err
	}
	fmt.Println(string(bookShelfBytes))

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Notion-Version", "2021-05-13")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", n.accesstoken))
	fmt.Println(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(res))
	defer resp.Body.Close()
	fmt.Println(resp.StatusCode)

	return nil
}
