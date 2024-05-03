package test

import (
	"reflect"
	"testing"

	jalibs "pro2.jobarranger.info/jobarranger/test/src_t/jalibs_t"
	jobarg_agentd "pro2.jobarranger.info/jobarranger/test/src_t/jobarg_agentd"
)

type ja_test interface {
	TestJaz() bool
}

func TestJobarranger(t *testing.T) {
	var tests []ja_test
	tests = append(tests, new(jalibs.Jalockutil_ja_test)) //add more test cases' object here
	tests = append(tests, new(jobarg_agentd.Jajobfile))

	for _, e := range tests {
		if !e.TestJaz() {
			t.Errorf("test failed for: %s", reflect.TypeOf(e))
		}
	}
}
