package json

import (
	"encoding/json"
	"sort"
	"strconv"
	"strings"
	"testing"

	"gotest.tools/assert/cmp"
)

type pointerTest struct {
	// Json Document
	doc string

	// JsonPath
	path string

	// List of expected pointers
	expected []string
}

type pointerExceptionTest struct {
	// Json Document
	doc string

	// JsonPath
	path string

	// expected exception
	expected string
}

const case1_simple string = `{
			"a": "b"
		}`

const case2_nested string = `{
			"a": "b",
			"v": [{
				"x": "test1",
				"y": "hello"
			},
			{
				"x": "test2",
				"y": "world"
			}],
			"f":{
				"w": "hi",
				"q": "welcome",
				"y": "ciao"
			},
			"y": "c"
		}`

const case3_route string = `{
  "name": "default/eg/http",
  "virtual_hosts": [
   {
      "name": "default/eg/http/www_test_com",
      "domains": [
        "www.test.com"
      ],
      "routes": [
        {
          "name": "httproute/default/backend/rule/0/match/0/www_test_com",
          "match": {
            "prefix": "/"
          },
          "route": {
            "cluster": "httproute/default/backend/rule/0",
            "upgrade_configs": [
              {
                "upgrade_type": "websocket"
              }
            ]
          }
        }
      ]
    },
    {
      "name": "default/eg/http/www_example_com",
      "domains": [
        "www.example.com"
      ],
      "routes": [
        {
          "name": "httproute/default/backend/rule/1/match/1/www_example_com",
          "match": {
            "prefix": "/"
          },
          "route": {
            "cluster": "httproute/default/backend/rule/1",
            "upgrade_configs": [
              {
                "upgrade_type": "websocket"
              }
            ]
          }
        }
      ]
    }
  ],
  "ignore_port_in_host_matching": true
}`

var pointerTests = []pointerTest{
	{
		doc:  case1_simple,
		path: "$.a",
		expected: []string{
			"/a",
		},
	},
	{
		doc:  case2_nested,
		path: "$.v[?(@.x=='test2')]",
		expected: []string{
			"/v/1",
		},
	},
	{
		doc:  case2_nested,
		path: "..v[?(@.x=='test1')].y",
		expected: []string{
			"/v/0/y",
		},
	},
	{
		doc:  case2_nested,
		path: "$.v[?(@.x=='test2')].y",
		expected: []string{
			"/v/1/y",
		},
	},
	{
		doc:  case2_nested,
		path: "$.v[?(@.x=='test1')].y",
		expected: []string{
			"/v/0/y",
		},
	},
	{
		doc:  case2_nested,
		path: "$.v[*].y",
		expected: []string{
			"/v/0/y",
			"/v/1/y",
		},
	},
	{
		doc:      case2_nested,
		path:     "$.v[?(@.x=='UNKNOWN')].y",
		expected: []string{},
	},
	{
		doc:  case1_simple,
		path: ".a",
		expected: []string{
			"/a",
		},
	},
	{
		doc:  case1_simple,
		path: "a",
		expected: []string{
			"/a",
		},
	},
	{
		doc:  case2_nested,
		path: "f.w",
		expected: []string{
			"/f/w",
		},
	},
	{
		doc:  case2_nested,
		path: "f.*",
		expected: []string{
			"/f/w",
			"/f/q",
			"/f/y",
		},
	},
	{
		doc:  case2_nested,
		path: "v.*",
		expected: []string{
			"/v/0",
			"/v/1",
		},
	},
	{
		doc:  case2_nested,
		path: "v.**",
		expected: []string{
			"/v/0/x",
			"/v/0/y",
			"/v/1/x",
			"/v/1/y",
		},
	},
	{
		doc:  case2_nested,
		path: "$..y",
		expected: []string{
			"/f/y",
			"/v/0/y",
			"/v/1/y",
			"/y",
		},
	},
	{
		doc:  case3_route,
		path: "..routes[?(@.name =~ 'www_example_com')]",
		expected: []string{
			"/virtual_hosts/1/routes/0",
		},
	},
	{
		doc:  case3_route,
		path: "..routes[?(@.name =~ 'www_test_com')]",
		expected: []string{
			"/virtual_hosts/0/routes/0",
		},
	},
	{
		doc:  case3_route,
		path: "..routes[?(@.name =~ 'www')]",
		expected: []string{
			"/virtual_hosts/0/routes/0",
			"/virtual_hosts/1/routes/0",
		},
	},
}

func Test(t *testing.T) {
	for i, test := range pointerTests {

		testCasePrefix := "TestCase " + strconv.Itoa(i+1)
		pointers, err := ConvertPathToPointers([]byte(test.doc), test.path)
		if err != nil {
			t.Error(testCasePrefix + ": Error during conversion:\n" + err.Error())
			continue
		}

		expectedAsString := asString(test.expected)
		pointersAsString := asString(pointers)

		areEqual := cmp.Equal(expectedAsString, pointersAsString)()

		if !areEqual.Success() {
			var obj map[string]interface{}
			json.Unmarshal([]byte(test.doc), &obj)
			prettyJson, _ := json.MarshalIndent(obj, "", "    ")
			t.Error(testCasePrefix + ": Compare failed!\nJsonDoc:\n" + string(prettyJson) + "\n\nJsonPath: '" + test.path + "'\n\nExpected pointers (" + strconv.Itoa(len(test.expected)) + "):\n" + expectedAsString + "\nbut found (" + strconv.Itoa(len(pointers)) + "):\n" + pointersAsString)
		}
	}
}

var pointerExceptionTests = []pointerExceptionTest{
	{
		doc:      case1_simple,
		path:     ".$",
		expected: "Error during parsing jpath",
	},
	{
		doc:      case1_simple,
		path:     "$",
		expected: "only Root",
	},
}

func TestException(t *testing.T) {
	for i, test := range pointerExceptionTests {

		testCasePrefix := "TestCase " + strconv.Itoa(i+1)
		_, err := ConvertPathToPointers([]byte(test.doc), test.path)
		if err == nil {
			t.Error(testCasePrefix + ": Error expected, but no error found!")
			continue
		}

		if !strings.Contains(err.Error(), test.expected) {
			t.Error(testCasePrefix + ": Compare failed!\nExpected Error should contain:\n'" + test.expected + "'\n\nbut error was:\n'" + err.Error() + "'\n")
		}
	}
}

func asString(values []string) string {
	var buf []byte

	sort.Strings(values)
	for _, v := range values {
		buf = append(buf, []byte("* '")...)
		buf = append(buf, []byte(v)...)
		buf = append(buf, []byte("'\n")...)
	}

	return string(buf)
}
