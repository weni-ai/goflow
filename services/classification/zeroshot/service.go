package zeroshot

import (
	"net/http"

	"github.com/nyaruka/gocommon/httpx"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

// a classification service implementation for a zeroshot bot
type service struct {
	client     *Client
	classifier *flows.Classifier
	redactor   utils.Redactor
}

// NewService creates a new classification service
func NewService(httpClient *http.Client, httpRetries *httpx.RetryConfig, classifier *flows.Classifier, accessToken string, repository string) flows.ClassificationService {
	return &service{
		client:     NewClient(httpClient, httpRetries, accessToken, repository),
		classifier: classifier,
		redactor:   utils.NewRedactor(flows.RedactionMask, accessToken),
	}
}

func (s *service) Classify(session flows.Session, input string, logHTTP flows.HTTPLogCallback) (*flows.Classification, error) {
	response, trace, err := s.client.Predict(input)
	if trace != nil {
		logHTTP(flows.NewHTTPLog(trace, flows.HTTPStatusFromCode, s.redactor))
	}
	if err != nil {
		return nil, err
	}

	result := &flows.Classification{
		Intents: make([]flows.ExtractedIntent, 0),
	}
	result.Intents = append(result.Intents, flows.ExtractedIntent{Name: response.Text})

	return result, nil
}

var _ flows.ClassificationService = (*service)(nil)
