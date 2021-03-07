package service

import (
	"testing"
)

// arg1 means argument 1 and arg2 means argument 2, and the expected stands for the 'result we expect'
type addTestFillProto struct {
	arg1, expected string
}

var addTestsFillProto = []addTestFillProto{
	addTestFillProto{"www.cat.com", "http://www.cat.com"},
	addTestFillProto{"http://dog.com", "http://dog.com"},
	addTestFillProto{"rat.com", "http://rat.com"},
	addTestFillProto{"https://fish.net", "https://fish.net"},
}

func TestFillProto(t *testing.T) {
	t.Run("it should return URL with HTTP", func(t *testing.T) {
		for _, test := range addTestsFillProto {
			if output := FillProto(test.arg1); output != test.expected {
				t.Errorf("Output %q not equal to expected %q", output, test.expected)
			}
		}
	})

}

type addTestFetchHTTP struct {
	arg1, expected string
}

var addTestsFetchHTTP = []addTestFetchHTTP{
	addTestFetchHTTP{"http://facebook.com", "success"},
	addTestFetchHTTP{"http://google.com", "success"},
	addTestFetchHTTP{"http://rattttss.com", "fail"},
	addTestFetchHTTP{"http://fishhhhh.net", "fail"},
}

func TestFetchHTTP(t *testing.T) {
	t.Run("it should return right response tatus", func(t *testing.T) {
		for _, test := range addTestsFetchHTTP {
			if output := FetchHTTP(test.arg1); output != test.expected {
				t.Errorf("Output %q not equal to expected %q", output, test.expected)
			}
		}
	})

}
