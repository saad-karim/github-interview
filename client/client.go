package client

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/pkg/errors"
)

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

type Client struct {
	HTTPClient HTTPClient
	RetryLimit int
}

type Metadata struct {
	Year       int    `json:"year"`
	Length     int    `json:"length"`
	Title      string `json:"title"`
	Subject    string `json:"subject"`
	Action     string `json:"action"`
	Actor      string `json:"actor"`
	Actress    string `json:"actress"`
	Director   string `json:"director"`
	Popularity int    `json:"popularity"`
	Awards     string `json:"awards"`
	Image      string `json:"image"`
}

func (c *Client) PushMetadata(csvFile string) error {
	allMetadata, err := c.readInCSVFile(csvFile)
	if err != nil {
		return err
	}

	for _, metadata := range allMetadata {
		fmt.Printf("Sending metdata: %+v\n", metadata)
		b, err := json.Marshal(metadata)
		if err != nil {
			return errors.Wrap(err, "failed to serialize data")
		}

		var resp *http.Response
		for i := 0; i <= c.RetryLimit; i++ {
			req, err := http.NewRequest("POST", "http://localhost:9009/movies", bytes.NewBuffer(b))
			if err != nil {
				return errors.Wrap(err, "failed to create request")
			}

			resp, err = c.HTTPClient.Do(req)
			if err != nil {
				return errors.Wrap(err, "failed request to server")
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusServiceUnavailable {
				break
			}
		}

		if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
			fmt.Println("failed for metadata: ", string(b))
			return fmt.Errorf("request failed with status code: %d", resp.StatusCode)
		}
	}

	return nil
}

func (c *Client) readInCSVFile(csvFile string) ([]Metadata, error) {
	csvfile, err := os.Open(csvFile)
	if err != nil {
		return nil, err
	}

	reader := csv.NewReader(csvfile)
	reader.Comma = ';'

	rows := 0
	moviesMetadata := []Metadata{}
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Println("ERROR: invalied row: ", row)
		}

		if rows == 0 {
			rows++
			continue // skip header row
		}

		year, err := strconv.Atoi(row[0])
		if err != nil {
			fmt.Println("ERROR: invalied number specified for year: ", row)
			continue
		}

		length, err := strconv.Atoi(row[1])
		if err != nil {
			fmt.Println("ERROR: invalied number specified for length: ", row)
			continue
		}

		pop, err := strconv.Atoi(row[7])
		if err != nil {
			fmt.Println("ERROR: invalied number specified for popularity: ", row)
			continue
		}

		metadata := Metadata{
			Year:       year,
			Length:     length,
			Title:      row[2],
			Subject:    row[3],
			Actor:      row[4],
			Actress:    row[5],
			Director:   row[6],
			Popularity: pop,
			Awards:     row[7],
			Image:      row[8],
		}

		moviesMetadata = append(moviesMetadata, metadata)
	}

	return moviesMetadata, nil
}
