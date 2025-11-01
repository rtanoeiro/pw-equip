package api

import (
	"net/http"
	"net/url"
	"os"
	"testing"
)

func TestParseEmailAndHWID(t *testing.T) {
	testCases := []struct {
		urlString               string
		expectedUserCredentials UserCredentials
		valid                   bool
	}{
		{
			urlString: "email=testemail@example.com&hwid=1234567890",
			expectedUserCredentials: UserCredentials{
				Email: "testemail@example.com",
				Hwid:  "1234567890",
			},
			valid: true,
		},
		{
			urlString: "email=anothertest@example.com&hwid=1234567890",
			expectedUserCredentials: UserCredentials{
				Email: "anothertest@example.com",
				Hwid:  "1234567890",
			},
			valid: true,
		},

		{
			urlString: "email=&hwid=1234567890",
			expectedUserCredentials: UserCredentials{
				Email: "",
				Hwid:  "1234567890",
			},
			valid: false,
		},
	}

	for _, testCase := range testCases {
		request := &http.Request{
			URL: &url.URL{
				RawQuery: testCase.urlString,
			},
		}
		userCredentials, valid := parseEmailAndHWID(request)
		if userCredentials.Email != testCase.expectedUserCredentials.Email {
			t.Errorf("Expected email %s, got %s", testCase.expectedUserCredentials.Email, userCredentials.Email)
		}
		if userCredentials.Hwid != testCase.expectedUserCredentials.Hwid {
			t.Errorf("Expected HWID %s, got %s", testCase.expectedUserCredentials.Hwid, userCredentials.Hwid)
		}
		if valid != testCase.valid {
			t.Errorf("Expected valid %t, got %t", testCase.valid, valid)
		}
	}
}

func TestGetEnvVar(t *testing.T) {
	testCases := []struct {
		key          string
		defaultValue string
		expected     string
	}{
		{key: "TEST_KEY", defaultValue: "default", expected: "default"},
	}

	for _, testCase := range testCases {
		env := GetEnvVar(testCase.key, testCase.defaultValue)
		if env != testCase.expected {
			t.Errorf("Expected %s, got %s", testCase.expected, env)
		}
	}
}

func TestLoadAPIConfig(t *testing.T) {
	config := EquipConfig{}
	testConfig := config.LoadEquipConfig()
	if testConfig.ApiPort != "8989" {
		t.Errorf("Expected port %s, got %s", "8989", testConfig.ApiPort)
	}
	if testConfig.MySQLPort != "3306" {
		t.Errorf("Expected port %s, got %s", "3306", testConfig.MySQLPort)
	}
}

func TestGetEnvVarValid(t *testing.T) {
	_ = os.Setenv("VAR1", "Test1")
	_ = os.Setenv("VAR2", "Test2")
	_ = os.Setenv("VAR3", "Test3")
	testCases := []struct {
		key          string
		defaultValue string
		expected     string
	}{
		{key: "VAR1", defaultValue: "TestFailed1", expected: "Test1"},
		{key: "VAR2", defaultValue: "TestFailed2", expected: "Test2"},
		{key: "VAR3", defaultValue: "TestFailed3", expected: "Test3"},
	}

	for _, testCase := range testCases {
		env := GetEnvVar(testCase.key, testCase.defaultValue)
		if env != testCase.expected {
			t.Errorf("Expected %s, got %s", testCase.expected, env)
		}
	}
}

func TestGetEnvVarInvalid(t *testing.T) {
	_ = os.Setenv("VAR1", "Test1")
	_ = os.Setenv("VAR2", "Test2")
	_ = os.Setenv("VAR3", "Test3")
	testCases := []struct {
		key          string
		defaultValue string
		expected     string
	}{
		{key: "NOTAVAR1", defaultValue: "TestFailed1", expected: "TestFailed1"},
		{key: "NOTAVAR2", defaultValue: "TestFailed2", expected: "TestFailed2"},
		{key: "NOTAVAR3", defaultValue: "TestFailed3", expected: "TestFailed3"},
	}
	for _, testCase := range testCases {
		env := GetEnvVar(testCase.key, testCase.defaultValue)
		if env != testCase.expected {
			t.Errorf("Expected %s, got %s", testCase.expected, env)
		}
	}
}
