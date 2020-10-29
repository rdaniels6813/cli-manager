package nodeman_test

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"testing"

	"github.com/rdaniels6813/cli-manager/internal/nodeman"
	"github.com/stretchr/testify/assert"
)

type clientMock struct {
}

func (c *clientMock) Do(req *http.Request) (*http.Response, error) {
	return &http.Response{
		Body: ioutil.NopCloser(bytes.NewReader(defaultBody)),
	}, nil
}

var tests = map[string]string{
	">=10.x <14.x": "12.19.0",
	"10.x":         "10.23.0",
	">=10.x":       "14.15.0",
	"12.x":         "12.19.0",
}

func getNodeSchedule() []byte {
	defaultBody, err := ioutil.ReadFile("./node-schedule.json")
	if err != nil {
		log.Fatal(err)
	}
	return defaultBody
}

var defaultBody = getNodeSchedule()

func TestGetNodeVersionByRangeOrLTS(t *testing.T) {
	for k := range tests {
		input := k
		expected := tests[k]
		t.Run(input, func(t *testing.T) {
			actual, err := nodeman.GetNodeVersionByRangeOrLTS(input, &clientMock{})
			assert.Nil(t, err)
			assert.Equal(t, expected, actual)
		})
	}
}
