package function

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

type AzureClient struct {
	apiKey  string
	baseUrl string
}

type DetectLanguageResponse struct {
	Documents []struct {
		DetectedLanguages []struct {
			Name        string  `json:"name"`
			ISO6391Name string  `json:"iso6391name"`
			Score       float64 `json:"score"`
		} `json:"detectedLanguages`
	} `json:"documents"`
}

type DetectKeyPhrasesResponse struct {
	Documents []struct {
		KeyPhrases []string `json:"keyPhrases`
	} `json:"documents"`
}

type DetectSentimentResponse struct {
	Documents []struct {
		Score float64 `json:"score`
	} `json:"documents"`
}

func NewAzureClient(apiKey string, baseUrl string) TextAnalyticsService {

	return &AzureClient{
		apiKey:  apiKey,
		baseUrl: baseUrl,
	}
}

func (c *AzureClient) newPostRequest(url string, json string) (*http.Request, error) {
	log.Printf("Creating Post Request: %s, %s\n", url, json)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(json)))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Ocp-Apim-Subscription-Key", c.apiKey)
	return req, nil
}

func doHTTP(req *http.Request) ([]byte, error) {

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (c *AzureClient) FetchLanguage(text string) (*LanguageResult, error) {
	url := c.baseUrl + "/languages"
	jsonStr := `{"documents":[{"id": 0,"text": "` + text + `"}]}`
	req, err := c.newPostRequest(url, jsonStr)

	body, err := doHTTP(req)
	if err != nil {
		return nil, err
	}

	var dlResponse DetectLanguageResponse
	err = json.Unmarshal(body, &dlResponse)
	if err != nil {
		return nil, err
	}

	result := dlResponse.Documents[0].DetectedLanguages[0]
	return &LanguageResult{
		Name:        result.Name,
		ISO6391Name: result.ISO6391Name,
		Score:       result.Score,
	}, nil
}

func (c *AzureClient) FetchKeyPhrases(r *AnalyzeRequest) ([]string, error) {

	url := c.baseUrl + "/keyPhrases"
	jsonStr := `{"documents":[{"id": 0,"language":"` + r.Language + `","text": "` + r.Text + `"}]}`
	req, err := c.newPostRequest(url, jsonStr)

	body, err := doHTTP(req)
	if err != nil {
		return nil, err
	}

	var dkResponse DetectKeyPhrasesResponse
	err = json.Unmarshal(body, &dkResponse)
	if err != nil {
		return nil, err
	}

	return dkResponse.Documents[0].KeyPhrases, nil
}

func (c *AzureClient) FetchSentiment(r *AnalyzeRequest) (float64, error) {

	url := c.baseUrl + "/sentiment"
	jsonStr := `{"documents": [{"id": 0,"language": "` + r.Language + `","text": "` + r.Text + `"}]}`
	req, err := c.newPostRequest(url, jsonStr)

	body, err := doHTTP(req)
	if err != nil {
		return -1.0, err
	}

	var dsResponse DetectSentimentResponse
	err = json.Unmarshal(body, &dsResponse)
	if err != nil {
		return -1.0, err
	} else if len(dsResponse.Documents) == 0 {
		return -1.0, errors.New("could not analyze text")
	}

	return dsResponse.Documents[0].Score, nil
}
