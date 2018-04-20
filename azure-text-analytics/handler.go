package function

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/openfaas-incubator/go-function-sdk"
)

// Handle a function invocation
func Handle(req handler.Request) (handler.Response, error) {
	var err error

	apiKeyBytes, err := ioutil.ReadFile("/run/secrets/" + os.Getenv("API_KEY_NAME"))
	if err != nil {
		return handler.Response{
			Body:       createErrorResponseBody(err.Error()),
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	url := os.Getenv("TEXT_ANALYZE_BASE_URL")

	service := &AzureClient{
		apiKey:  string(apiKeyBytes),
		baseUrl: url,
	}

	var request AnalyzeRequest
	err = json.Unmarshal(req.Body, &request)
	if err != nil {
		return handler.Response{
			Body:       createErrorResponseBody(err.Error()),
			StatusCode: http.StatusBadRequest,
		}, err
	}
	log.Printf("REQUEST:%+v", request)

	response, err := internalHandle(&request, service)
	if err != nil {
		return handler.Response{
			Body:       createErrorResponseBody(response.ErrorMessage),
			StatusCode: response.StatusCode,
		}, err
	}

	body, _ := json.Marshal(response.Result)
	return handler.Response{
		Body:       body,
		StatusCode: response.StatusCode,
	}, err
}

func internalHandle(req *AnalyzeRequest, service TextAnalyticsService) (*AnalyzeResponse, error) {

	var err error
	var result AnalyzeResult

	if len(req.Language) == 0 {
		lr, err := service.FetchLanguage(req.Text)
		if err != nil {
			return &AnalyzeResponse{
				ErrorMessage: fmt.Sprintf("Analyzing language failed: %s", err.Error()),
				StatusCode:   http.StatusBadRequest,
			}, err
		}

		req.Language = lr.ISO6391Name
	}
	result.Language = req.Language

	score, err := service.FetchSentiment(req)
	if err != nil {
		return &AnalyzeResponse{
			ErrorMessage: fmt.Sprintf("Analyzing sentiment failed: %s", err.Error()),
			StatusCode:   http.StatusInternalServerError,
		}, err
	}
	result.SentimentScore = score

	kp, err := service.FetchKeyPhrases(req)
	if err != nil {
		return &AnalyzeResponse{
			ErrorMessage: fmt.Sprintf("Analyzing key phrases failed: %s", err.Error()),
			StatusCode:   http.StatusInternalServerError,
		}, err
	}
	result.KeyPhrases = kp

	return &AnalyzeResponse{
		Result:     &result,
		StatusCode: http.StatusOK,
	}, err
}
