package plantuml

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/steinfletcher/apitest"
	"io"
	"strconv"
)

const requestOperation = "->"
const responseOperation = "-->"

type Formatter struct {
	writer io.Writer
}

type DSL struct {
	count int
	data  bytes.Buffer
}

func (r *DSL) AddRequestRow(source, target, description, body string) {
	r.addRow(requestOperation, source, target, description, body)
}

func (r *DSL) AddResponseRow(source, target, description, body string) {
	r.addRow(responseOperation, source, target, description, body)
}

func (r *DSL) addRow(operation, source, target, description, body string) {
	var notePosition = "left"
	if operation == requestOperation {
		notePosition = "right"
	}

	var note string
	if body != "" {
		note = fmt.Sprintf("\nnote %s\n%s\nend note", notePosition, escape(body))
	}

	r.count += 1
	r.data.WriteString(fmt.Sprintf("%s%s%s: (%d) %s %s\n",
		source,
		operation,
		target,
		r.count,
		description,
		note))
}

func (r *DSL) ToString() string {
	return fmt.Sprintf("@startuml\nskinparam noteFontSize 11\nskinparam monochrome true\n%s\n@enduml", r.data.String())
}

func (r *Formatter) Format(recorder *apitest.Recorder) {
	markup, err := buildMarkup(recorder)
	if err != nil {
		panic(err)
	}

	_, err = r.writer.Write([]byte(markup))
	if err != nil {
		panic(err)
	}
}

func NewFormatter(writer io.Writer) apitest.ReportFormatter {
	return &Formatter{writer: writer}
}

func buildMarkup(r *apitest.Recorder) (string, error) {
	if len(r.Events) == 0 {
		return "", errors.New("no events are defined")
	}

	dsl := &DSL{}
	for _, event := range r.Events {
		switch v := event.(type) {
		case apitest.HttpRequest:
			httpReq := v.Value
			entry, err := apitest.NewHttpRequestLogEntry(httpReq)
			if err != nil {
				return "", err
			}
			entry.Timestamp = v.Timestamp
			dsl.AddRequestRow(v.Source, v.Target, fmt.Sprintf("%s %s", httpReq.Method, httpReq.URL), formatNote(entry))
		case apitest.HttpResponse:
			entry, err := apitest.NewHttpResponseLogEntry(v.Value)
			if err != nil {
				return "", err
			}
			entry.Timestamp = v.Timestamp
			dsl.AddResponseRow(v.Source, v.Target, strconv.Itoa(v.Value.StatusCode), formatNote(entry))
		case apitest.MessageRequest:
			dsl.AddRequestRow(v.Source, v.Target, v.Header, v.Body)
		case apitest.MessageResponse:
			dsl.AddResponseRow(v.Source, v.Target, v.Header, v.Body)
		default:
			panic("received unknown event type")
		}
	}

	return dsl.ToString(), nil
}

func escape(in string) string {
	return in
	//in = strings.Replace(in,"\n","\\n",-1)
	//return strings.Replace(in,"\t","\\t",-1)
}

func formatNote(entry apitest.LogEntry) string {
	return fmt.Sprintf("%s%s", entry.Header, entry.Body)
	//return fmt.Sprintf("%s%s", escape(entry.Header), escape(entry.Body))
}