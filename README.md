# apitest-plantuml

[![Build Status](https://travis-ci.org/steinfletcher/apitest-plantuml.svg?branch=master)](https://travis-ci.org/steinfletcher/apitest-plantuml)

Formats the results of an [apitest](https://github.com/steinfletcher/apitest) run as plant uml markup.

## Installation

```bash
go get -u github.com/steinfletcher/apitest-plantuml
```

## Example

```go
apitest.New("search user").
    Handler(myHandler).
	Report(plantuml.NewFormatter(myWriter)).
	Mocks(getPreferencesMock, getUserMock).
	Post("/user/search").
	Body(`{"name": "jon"}`).
	Expect(t).
	Status(http.StatusOK).
	Header("Content-Type", "application/json").
	Body(`{"name": "jon", "is_contactable": true}`).
	End()
```

![Diagram](/testdata/plantuml.png?raw=true "Sequence Diagram")

[SVG version](https://www.plantuml.com/plantuml/svg/fPFBJiCm44NtynMZhAX4IPEsI6I1AY6W4aYjrEoHSQUDIs97jhC0nBzZkpoIbaSHMJapetlld3WJOvcsJLM2UH2oPffLA9MbAoNjGZmH9achKocfUA5LHMXrGn3nKaJzyyWqDihmAEdXVBR8CMuCwTWGqxn0y7AenRgmiD-TvlayJauIc2fZCsHrNGhEh50Iu1dGFP5a5JdrQADa12z0SXaIGd1rvbxEUFkqXzUNEHRMrbaSN4pCQh7rIzBXQDm9usTt-pjnWWP0otDhrl_SUTZ3T31u4ovfPU5T8zHdDt3XK9Aq_LkIQrjac9vzbFB7cZfBnnlppUR7sv9O-e8F-oMBMjEAD4bEWSYeGwJL37kttt_0ip-sGwarq67L3jCYDluxiT6Xn8IvuswISSlkIyykNDCzK2wClxakpEaXOIHnRvKXBvMX_tozt-B1n1tzt_WA)
