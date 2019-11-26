package parser

import (
	"golang.org/x/net/html"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	rawInput := `
<!doctype html>
<head>
	<link rel="stylesheet" href="/css/example.css">
	<link rel="canonical" href="https://example.com">
</head>
<body>
	<a href="/about">About Me</a>
	<a href="https://example.com/blog"></a>
	<a href="https://google.com">Search</a>
</body>
`
	input := strings.NewReader(rawInput)
	expectedDetails := PageDetails{
		Assets: []string{
			"https://example.com/css/example.css",
			"https://example.com",
		},
		ExternalLinks: []string{
			"https://google.com",
		},
		InternalLinks: []string{
			"https://example.com/about",
			"https://example.com/blog",
		},
	}
	rawurl := "https://example.com"
	url, err := url.Parse(rawurl)
	if err != nil {
		t.Fatalf("Couldn't parse url '%s'", rawurl)
	}
	details := Parse(url, input)

	if !reflect.DeepEqual(details, expectedDetails) {
		t.Fatalf("Expected %v, got %v", expectedDetails, details)
	}
}

func TestGetAttribute(t *testing.T) {
	testCases := []struct {
		Attr          []html.Attribute
		ExpectedValue string
	}{
		{
			Attr: []html.Attribute{
				html.Attribute{
					Namespace: "", Key: "href", Val: "https://example.com",
				},
			},
			ExpectedValue: "https://example.com",
		},
		{
			Attr: []html.Attribute{
				html.Attribute{
					Namespace: "", Key: "href", Val: "bar",
				},
			},
			ExpectedValue: "",
		},
	}

	for _, testCase := range testCases {
		value := getHref(testCase.Attr)
		if value != testCase.ExpectedValue {
			t.Fatalf("Expected %s. got %s", testCase.ExpectedValue, value)
		}
		if len(value) == 0 {
			t.Fatalf("Expected %s. got %s", testCase.ExpectedValue, value)
		}
	}
}
