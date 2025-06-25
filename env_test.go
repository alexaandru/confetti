//nolint:testpackage // ok
package confetti

import "testing"

func TestCamelToUpperSnake(t *testing.T) {
	t.Parallel()

	cases := []struct {
		in, want string
	}{
		{"AWSRegion", "AWS_REGION"},
		{"HTTPRequest", "HTTP_REQUEST"},
		{"OAuthToken", "O_AUTH_TOKEN"},
		{"MyID", "MY_ID"},
		{"MyId", "MY_ID"},
		{"UserID", "USER_ID"},
		{"SimpleTest", "SIMPLE_TEST"},
		{"XMLParser", "XML_PARSER"},
		{"JSONData", "JSON_DATA"},
		{"ABTest", "AB_TEST"},
		{"TestA", "TEST_A"},
		{"TestAB", "TEST_AB"},
		{"TestABC", "TEST_ABC"},
		{"TestABCTest", "TEST_ABC_TEST"},
		{"test", "TEST"},
		{"A", "A"},
		{"ID", "ID"},
		{"", ""},
		{"lowercase", "LOWERCASE"},
		{"CamelCase", "CAMEL_CASE"},
		{"CamelCASE", "CAMEL_CASE"},
		{"CamelCaseX", "CAMEL_CASE_X"},
		{"CamelCASETest", "CAMEL_CASE_TEST"},
		{"HTTPRequestID", "HTTP_REQUEST_ID"},
		{"HTTPRequestId", "HTTP_REQUEST_ID"},
		{"HTTPRequestIDTest", "HTTP_REQUEST_ID_TEST"},
	}

	for _, c := range cases {
		got := camelToUpperSnake(c.in)
		if got != c.want {
			t.Errorf("camelToUpperSnake(%q) = %q; want %q", c.in, got, c.want)
		}
	}
}
