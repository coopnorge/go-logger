package logger

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"testing"
)

var testBuffer bytes.Buffer
var testTable = []struct {
	name   string
	msg    string
	fields Fields
}{
	{
		"Test 0 - normal message, testfield",
		"Test message number1",
		Fields{
			"testField": "127.0.0.1",
		},
	},
	{
		"Test 1 - number message, no testfield",
		"123",
		Fields{},
	},
	{
		"Test 2 - empty message, no testfield",
		"",
		Fields{},
	},
}

func LogAndAssertJSON(t *testing.T, log func(*LogrusFacade), assertions func(fields Fields)) {
	fields := make(map[string]interface{})

	testBuffer.Reset()

	logger := Logrus(&testBuffer)

	log(logger)

	if err := json.Unmarshal(testBuffer.Bytes(), &fields); err != nil {
		t.Errorf("LogAndAssertJSON failed with error: %v", err)
	}

	assertions(fields)
}

func TestWithFieldsInfof(t *testing.T) {
	for _, test := range testTable {
		LogAndAssertJSON(t, func(log *LogrusFacade) {
			log.LogWithFields(test.fields).Infof(test.msg)
		}, func(fields Fields) {
			wantedMsg := test.msg
			gotMsg := fields["msg"]
			if gotMsg != wantedMsg {
				t.Errorf("Test message not correct want %v got %v", wantedMsg, gotMsg)
			}
			wantedLevel := "info"
			gotLevel := fields["level"]
			if gotLevel != wantedLevel {
				t.Errorf("Test level not correct want %v got %v", wantedLevel, gotLevel)
			}
			wantedTestField := test.fields["testField"]
			gotTestField := fields["testField"]
			if gotTestField != wantedTestField {
				t.Errorf("Test field not correct want %v got %v", wantedTestField, gotTestField)
			}
		})
	}
}

func TestWithFieldsPrintf(t *testing.T) {
	for _, test := range testTable {
		LogAndAssertJSON(t, func(log *LogrusFacade) {
			log.LogWithFields(test.fields).Printf(test.msg)
		}, func(fields Fields) {
			wantedMsg := test.msg
			gotMsg := fields["msg"]
			if gotMsg != wantedMsg {
				t.Errorf("Test message not correct want %v got %v", wantedMsg, gotMsg)
			}
			wantedLevel := "info"
			gotLevel := fields["level"]
			if gotLevel != wantedLevel {
				t.Errorf("Test level not correct want %v got %v", wantedLevel, gotLevel)
			}
			wantedTestField := test.fields["testField"]
			gotTestField := fields["testField"]
			if gotTestField != wantedTestField {
				t.Errorf("Test field not correct want %v got %v", wantedTestField, gotTestField)
			}
		})
	}
}

func TestWithFieldsWarnf(t *testing.T) {
	for _, test := range testTable {
		LogAndAssertJSON(t, func(log *LogrusFacade) {
			log.LogWithFields(test.fields).Warnf(test.msg)
		}, func(fields Fields) {
			wantedMsg := test.msg
			gotMsg := fields["msg"]
			if gotMsg != wantedMsg {
				t.Errorf("Test message not correct want %v got %v", wantedMsg, gotMsg)
			}
			wantedLevel := "warning"
			gotLevel := fields["level"]
			if gotLevel != wantedLevel {
				t.Errorf("Test level not correct want %v got %v", wantedLevel, gotLevel)
			}
			wantedTestField := test.fields["testField"]
			gotTestField := fields["testField"]
			if gotTestField != wantedTestField {
				t.Errorf("Test field not correct want %v got %v", wantedTestField, gotTestField)
			}
		})
	}
}

func TestWithFieldsErrorf(t *testing.T) {
	for _, test := range testTable {
		LogAndAssertJSON(t, func(log *LogrusFacade) {
			log.LogWithFields(test.fields).Errorf(test.msg)
		}, func(fields Fields) {
			wantedMsg := test.msg
			gotMsg := fields["msg"]
			if gotMsg != wantedMsg {
				t.Errorf("Test message not correct want %v got %v", wantedMsg, gotMsg)
			}
			wantedLevel := "error"
			gotLevel := fields["level"]
			if gotLevel != wantedLevel {
				t.Errorf("Test level not correct want %v got %v", wantedLevel, gotLevel)
			}
			wantedTestField := test.fields["testField"]
			gotTestField := fields["testField"]
			if gotTestField != wantedTestField {
				t.Errorf("Test field not correct want %v got %v", wantedTestField, gotTestField)
			}
		})
	}
}

func TestPrint(t *testing.T) {
	for _, test := range testTable {
		LogAndAssertJSON(t, func(log *LogrusFacade) {
			log.LogPrint(test.msg)
		}, func(fields Fields) {
			wantedMsg := test.msg
			gotMsg := fields["msg"]
			if gotMsg != wantedMsg {
				t.Errorf("Test message not correct want %v got %v", wantedMsg, gotMsg)
			}
			wantedLevel := "info"
			gotLevel := fields["level"]
			if gotLevel != wantedLevel {
				t.Errorf("Test level not correct want %v got %v", wantedLevel, gotLevel)
			}
		})
	}
}

func TestError(t *testing.T) {
	for _, test := range testTable {
		LogAndAssertJSON(t, func(log *LogrusFacade) {
			log.LogError(test.msg)
		}, func(fields Fields) {
			wantedMsg := test.msg
			gotMsg := fields["msg"]
			if gotMsg != wantedMsg {
				t.Errorf("Test message not correct want %v got %v", wantedMsg, gotMsg)
			}
			wantedLevel := "error"
			gotLevel := fields["level"]
			if gotLevel != wantedLevel {
				t.Errorf("Test level not correct want %v got %v", wantedLevel, gotLevel)
			}
		})
	}
}

func TestInfof(t *testing.T) {
	for _, test := range testTable {
		LogAndAssertJSON(t, func(log *LogrusFacade) {
			log.LogInfof(test.msg)
		}, func(fields Fields) {
			wantedMsg := test.msg
			gotMsg := fields["msg"]
			if gotMsg != wantedMsg {
				t.Errorf("Test message not correct want %v got %v", wantedMsg, gotMsg)
			}
			wantedLevel := "info"
			gotLevel := fields["level"]
			if gotLevel != wantedLevel {
				t.Errorf("Test level not correct want %v got %v", wantedLevel, gotLevel)
			}
		})
	}
}

func TestPrintf(t *testing.T) {
	for _, test := range testTable {
		LogAndAssertJSON(t, func(log *LogrusFacade) {
			log.LogPrintf(test.msg)
		}, func(fields Fields) {
			wantedMsg := test.msg
			gotMsg := fields["msg"]
			if gotMsg != wantedMsg {
				t.Errorf("Test message not correct want %v got %v", wantedMsg, gotMsg)
			}
			wantedLevel := "info"
			gotLevel := fields["level"]
			if gotLevel != wantedLevel {
				t.Errorf("Test level not correct want %v got %v", wantedLevel, gotLevel)
			}
		})
	}
}

func TestWarnf(t *testing.T) {
	for _, test := range testTable {
		LogAndAssertJSON(t, func(log *LogrusFacade) {
			log.LogWarnf(test.msg)
		}, func(fields Fields) {
			wantedMsg := test.msg
			gotMsg := fields["msg"]
			if gotMsg != wantedMsg {
				t.Errorf("Test message not correct want %v got %v", wantedMsg, gotMsg)
			}
			wantedLevel := "warning"
			gotLevel := fields["level"]
			if gotLevel != wantedLevel {
				t.Errorf("Test level not correct want %v got %v", wantedLevel, gotLevel)
			}
		})
	}
}

func TestErrorf(t *testing.T) {
	for _, test := range testTable {
		LogAndAssertJSON(t, func(log *LogrusFacade) {
			log.LogErrorf(test.msg)
		}, func(fields Fields) {
			wantedMsg := test.msg
			gotMsg := fields["msg"]
			if gotMsg != wantedMsg {
				t.Errorf("Test message not correct want %v got %v", wantedMsg, gotMsg)
			}
			wantedLevel := "error"
			gotLevel := fields["level"]
			if gotLevel != wantedLevel {
				t.Errorf("Test level not correct want %v got %v", wantedLevel, gotLevel)
			}
		})
	}
}

func TestRequest(t *testing.T) {
	wantedMsg := "Test message for TestRequest"
	req := &http.Request{
		Method: http.MethodGet,
		URL:    &url.URL{},
	}
	LogAndAssertJSON(t, func(log *LogrusFacade) {
		reqFields := Fields{
			"remote_ip":  req.RemoteAddr,
			"user_agent": req.Header.Get("User-Agent"),
			"method":     req.Method,
		}
		log.LogWithFields(reqFields).Infof(wantedMsg)
	}, func(fields Fields) {
		gotMsg := fields["msg"]
		if gotMsg != wantedMsg {
			t.Errorf("Test message not correct want %v got %v", wantedMsg, gotMsg)
		}
		wantedMethod := req.Method
		gotMethod := fields["method"]
		if gotMethod != wantedMethod {
			t.Errorf("Request method not correct want %v got %v", wantedMethod, gotMethod)
		}
	})
}

func TestNilRequest(t *testing.T) {
	wantedMsg := "Test message for TestNilRequest"
	req := &http.Request{
		URL: &url.URL{},
	}
	LogAndAssertJSON(t, func(log *LogrusFacade) {
		reqFields := Fields{
			"remote_ip":  req.RemoteAddr,
			"user_agent": req.Header.Get("User-Agent"),
			"method":     req.Method,
		}
		log.LogWithFields(reqFields).Infof(wantedMsg)
	}, func(fields Fields) {
		gotMsg := fields["msg"]
		if gotMsg != wantedMsg {
			t.Errorf("Test message not correct want %v got %v", wantedMsg, gotMsg)
		}
		wantedMethod := req.Method
		gotMethod := fields["method"]
		if gotMethod != wantedMethod {
			t.Errorf("Request method not correct want %v got %v", wantedMethod, gotMethod)
		}
	})
}
