package unsplash

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	// ImageNotFoundURL is to be used in case the methods return error
	ImageNotFoundURL string = "https://www.gannett-cdn.com/-mm-/331318330f2cfeaacf892e0b15cf38f07a8c0e09/c=278-0-3092-3752/local/-/media/2017/04/21/USATODAY/USATODAY/636284056090492726-Vin-Diesel.jpg"
	unsplashEndpoint string = "https://api.unsplash.com/photos"
)

type photosResponse struct {
	ID             string        `json:"id"`
	CreatedAt      string        `json:"created_at"`
	UpdatedAt      string        `json:"updated_at"`
	PromotedAt     string        `json:"promoted_at"`
	Width          int64         `json:"width"`
	Height         int64         `json:"height"`
	Color          string        `json:"color"`
	BlurHash       string        `json:"blur_hash"`
	Description    string        `json:"description"`
	AltDescription string        `json:"alt_description"`
	Urls           urls          `json:"urls"`
	Links          links         `json:"links"`
	Categories     []interface{} `json:"categories"`
}

type links struct {
	Self             string `json:"self"`
	HTML             string `json:"html"`
	Download         string `json:"download"`
	DownloadLocation string `json:"download_location"`
}

type urls struct {
	Raw     string `json:"raw"`
	Full    string `json:"full"`
	Regular string `json:"regular"`
	Small   string `json:"small"`
	Thumb   string `json:"thumb"`
}

// GetRandomPhoto gets the url of a random photo
func GetRandomPhoto(accessKey string, keywords []string, collectionID string) (string, error) {
	query := strings.Join(keywords, ",")
	url := generateUnsplashURL(accessKey, query, collectionID)

	fmt.Println("calling url: ", url)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("Got unsuccessful unsplash response: " + string(body))
	}

	var photosResp photosResponse
	err = json.Unmarshal(body, &photosResp)
	if err != nil {
		return "", err
	}

	return photosResp.Urls.Full, nil
}

func generateUnsplashURL(accessKey string, query string, collectionID string) string {
	var sb strings.Builder
	sb.WriteString(unsplashEndpoint)
	sb.WriteString("/random?client_id=")
	sb.WriteString(accessKey)
	if query != "" {
		sb.WriteString("&query=")
		sb.WriteString(query)
	}
	if collectionID != "" {
		sb.WriteString("&collectionId=")
		sb.WriteString(collectionID)
	}
	return sb.String()
}
