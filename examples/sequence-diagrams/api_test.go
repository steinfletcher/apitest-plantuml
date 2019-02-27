package main

import (
	"github.com/steinfletcher/apitest"
	"github.com/steinfletcher/apitest-plantuml"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

func TestGetUser(t *testing.T) {
	apiTest("gets the user 1").
		Mocks(getPreferencesMock, getUserMock).
		Post("/user/search").
		Body(`{"name": "jon"}`).
		Expect(t).
		Status(http.StatusOK).
		Header("Content-Type", "application/json").
		Body(`{"name": "jon", "is_contactable": true}`).
		End()
}

var getPreferencesMock = apitest.NewMock().
	Get("http://preferences/api/preferences/12345").
	RespondWith().
	Body(`{"is_contactable": true}`).
	Status(http.StatusOK).
	End()

var getUserMock = apitest.NewMock().
	Get("http://users/api/user/12345").
	RespondWith().
	Body(`{"name": "jon", "id": "1234"}`).
	Status(http.StatusOK).
	End()

type fileWriter struct{}

func (p *fileWriter) Write(data []byte) (int, error) {
	err := os.MkdirAll(".sequence", os.ModePerm)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(".sequence/diagram.txt", data, os.ModePerm)
	if err != nil {
		panic(err)
	}
	return -1, nil
}

func apiTest(name string) *apitest.APITest {
	return apitest.New(name).
		Report(plantuml.NewFormatter(&fileWriter{})).
		Meta(map[string]interface{}{"app": "user-service"}).
		Handler(newApp().Router)
}
