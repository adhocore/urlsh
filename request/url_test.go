package request

import (
	"errors"
	"testing"

	"github.com/adhocore/urlsh/common"
)

func tester(input URLInput, expect error, t *testing.T) {
	if actual := input.Validate(); !errors.Is(actual, expect) {
		t.Helper()
		t.Errorf("wanted %v error, got %v", expect, actual)
	}
}

func TestURLInput_Validate(t *testing.T) {
	t.Run("validate empty", func(t *testing.T) {
		input := URLInput{URL: ""}

		tester(input, common.ErrInvalidURLLen, t)
	})

	t.Run("validate empty", func(t *testing.T) {
		input := URLInput{URL: "xyz://x/y/z"}

		tester(input, common.ErrInvalidURL, t)
	})

	t.Run("invalid url", func(t *testing.T) {
		input := URLInput{URL: "http://local\\host"}

		tester(input, common.ErrInvalidURL, t)
	})

	t.Run("blacklist url", func(t *testing.T) {
		input := URLInput{URL: "http://localhost/xxx"}

		tester(input, common.ErrBlacklistedURL, t)
	})

	t.Run("keywords", func(t *testing.T) {
		input := URLInput{URL: "http://localhost", Keywords: make([]string, 11)}

		tester(input, common.ErrKeywordsCount, t)
	})

	t.Run("keyword length", func(t *testing.T) {
		input := URLInput{URL: "http://localhost", Keywords: []string{"x"}}

		tester(input, common.ErrKeywordLength, t)
	})

	t.Run("invalid date", func(t *testing.T) {
		input := URLInput{URL: "http://localhost", ExpiresOn: "2030x01x01x00x00x00"}

		tester(input, common.ErrInvalidDate, t)
	})

	t.Run("validate OK", func(t *testing.T) {
		input := URLInput{URL: "http://localhost", ExpiresOn: "2030-01-01 00:00:00"}

		if input.Validate() != nil {
			t.Errorf("valid data should not give error")
		}
	})
}

func TestURLInput_ValidateExpiry(t *testing.T) {
	t.Run("invalid expiry", func(t *testing.T) {
		input := URLInput{URL: "http://localhost", ExpiresOn: ""}

		tester(input, nil, t)
	})

	t.Run("invalid expiry", func(t *testing.T) {
		input := URLInput{URL: "http://localhost", ExpiresOn: "2020-01-01"}

		tester(input, common.ErrInvalidDate, t)
	})

	t.Run("past expiry", func(t *testing.T) {
		input := URLInput{URL: "http://localhost", ExpiresOn: "2020-01-01 00:00:00"}

		tester(input, common.ErrPastExpiration, t)
	})
}

func TestURLInput_GetExpiresOn(t *testing.T) {
	t.Run("get expires_on empty", func(t *testing.T) {
		input := URLInput{ExpiresOn: ""}

		if _, err := input.GetExpiresOn(); err != nil {
			t.Errorf("empty date should not give error")
		}
	})

	t.Run("get expires_on invalid", func(t *testing.T) {
		input := URLInput{ExpiresOn: "Jan 2020"}

		if _, err := input.GetExpiresOn(); err == nil {
			t.Errorf("invalid date should give error")
		}
	})

	t.Run("get expires_on ok", func(t *testing.T) {
		input := URLInput{ExpiresOn: "2020-01-01 00:00:00"}

		if _, err := input.GetExpiresOn(); err != nil {
			t.Errorf("valid date should not give error")
		}
	})
}

func TestURLFilter_GetOffset(t *testing.T) {
	t.Run("get offset - empty", func(t *testing.T) {
		if actual := (URLFilter{Page: ""}).GetOffset(2); actual != 0 {
			t.Errorf("offset wanted 0, got %v", actual)
		}
	})

	t.Run("get offset - less than 2", func(t *testing.T) {
		if actual := (URLFilter{Page: "1"}).GetOffset(2); actual != 0 {
			t.Errorf("offset wanted 0, got %v", actual)
		}
	})

	t.Run("get offset", func(t *testing.T) {
		if actual := (URLFilter{Page: "2"}).GetOffset(50); actual != 50 {
			t.Errorf("offset wanted 100, got %v", actual)
		}

		if actual := (URLFilter{Page: "3"}).GetOffset(50); actual != 100 {
			t.Errorf("offset wanted 100, got %v", actual)
		}
	})
}
