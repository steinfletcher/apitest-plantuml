package plantuml

import (
	"fmt"
	"github.com/steinfletcher/apitest"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

type writer struct {
	captured string
}

func (p *writer) Write(data []byte) (int, error) {
	p.captured = strings.TrimSpace(string(data))
	return -1, nil
}

func TestNewFormatter(t *testing.T) {
	t.SkipNow()

	recorder := aRecorder()
	capture := &writer{}

	NewFormatter(capture).Format(recorder)

	actual := string(capture.captured)
	expected := ``

	if !reflect.DeepEqual(expected, actual) {
		fmt.Printf("Expected '%s'\nReceived '%s'\n", expected, actual)
		t.Fail()
	}
}

func aRecorder() *apitest.Recorder {
	return apitest.NewTestRecorder().
		AddTitle("title").
		AddSubTitle("subTitle").
		AddHttpRequest(aRequest()).
		AddMessageRequest(apitest.MessageRequest{Header: "SQL Query", Body: "SELECT * FROM users", Source: "sut-a", Target: "a"}).
		AddMessageResponse(apitest.MessageResponse{Header: "SQL Result", Body: "Rows count: 122", Source: "a", Target: "sut-a"}).
		AddHttpResponse(aResponse()).
		AddMeta(map[string]interface{}{
			"path":   "/user",
			"name":   "some test",
			"host":   "example.com",
			"method": "GET",
		})
}

func aRequest() apitest.HttpRequest {
	req := httptest.NewRequest(http.MethodGet, "http://example.com/abcdef", nil)
	req.Header.Set("Content-Type", "application/json")
	return apitest.HttpRequest{Value: req, Source: "cli", Target: "sut-a"}
}

func aResponse() apitest.HttpResponse {
	return apitest.HttpResponse{
		Value: &http.Response{
			StatusCode:    http.StatusOK,
			ProtoMajor:    1,
			ProtoMinor:    1,
			ContentLength: 0,
		},
		Source: "sut-a",
		Target: "cli",
	}
}
