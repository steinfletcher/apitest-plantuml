package plantuml

import (
	"fmt"
	"github.com/steinfletcher/apitest"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
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

func normalize(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func TestNewFormatter(t *testing.T) {
	recorder := aRecorder()
	capture := &writer{}

	NewFormatter(capture).Format(recorder)

	actual := capture.captured
	expected, _ := ioutil.ReadFile("testdata/snapshot.txt")

	if normalize(string(expected)) != normalize(actual) {
		fmt.Printf("Expected '%s'\nReceived '%s'\n", string(expected), actual)
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
