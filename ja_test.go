package test

import (
	"reflect"
	"testing"

	jalibs "pro2.jobarranger.info/jobarranger/test/src_t/jalibs_t"
)

type ja_test interface {
	TestJaz() bool
}

func TestCallTest(t *testing.T) {
	var tests []ja_test
	tests = append(tests, new(jalibs.Jalockutil_ja_test)) //add more test cases' object here

	for _, e := range tests {
		if !e.TestJaz() {
			t.Errorf("test failed for: %s", reflect.TypeOf(e))
		}
	}

}
