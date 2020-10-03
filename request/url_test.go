package request

import (
    "errors"
    "github.com/adhocore/urlsh/common"
    "testing"
)

func tester(input UrlInput, expect error, t *testing.T) {
    if actual := input.Validate(); !errors.Is(actual, expect) {
        t.Helper()
        t.Errorf("wanted %v error, got %v", expect, actual)
    }
}

func TestUrlInput_Validate(t *testing.T) {
    t.Run("validate empty", func(t *testing.T) {
        input := UrlInput{Url:""}

        tester(input, common.ErrInvalidUrlLen, t)
    })

    t.Run("validate empty", func(t *testing.T) {
        input := UrlInput{Url: "xyz://x/y/z"}

        tester(input, common.ErrInvalidUrl, t)
    })

    t.Run("invalid url", func(t *testing.T) {
        input := UrlInput{Url: "http://local\\host"}

        tester(input, common.ErrInvalidUrl, t)
    })

    t.Run("blacklist url", func(t *testing.T) {
        input := UrlInput{Url: "http://localhost/xxx"}

        tester(input, common.ErrBlacklistedUrl, t)
    })

    t.Run("keywords", func(t *testing.T) {
        input := UrlInput{Url: "http://localhost", Keywords: make([]string, 11)}

        tester(input, common.ErrKeywordsCount, t)
    })

    t.Run("keyword length", func(t *testing.T) {
        input := UrlInput{Url: "http://localhost", Keywords: []string{"x"}}

        tester(input, common.ErrKeywordLength, t)
    })

    t.Run("invalid date", func(t *testing.T) {
        input := UrlInput{Url: "http://localhost", ExpiresOn: "2030x01x01x00x00x00"}

        tester(input, common.ErrInvalidDate, t)
    })

    t.Run("validate OK", func(t *testing.T) {
        input := UrlInput{Url: "http://localhost", ExpiresOn: "2030-01-01 00:00:00"}

        if input.Validate() != nil {
            t.Errorf("valid data should not give error")
        }
    })
}

func TestUrlInput_ValidateExpiry(t *testing.T) {
    t.Run("invalid expiry", func(t *testing.T) {
        input := UrlInput{Url: "http://localhost", ExpiresOn: ""}

        tester(input, nil, t)
    })

    t.Run("invalid expiry", func(t *testing.T) {
        input := UrlInput{Url: "http://localhost", ExpiresOn: "2020-01-01"}

        tester(input, common.ErrInvalidDate, t)
    })

    t.Run("past expiry", func(t *testing.T) {
        input := UrlInput{Url: "http://localhost", ExpiresOn: "2020-01-01 00:00:00"}

        tester(input, common.ErrPastExpiration, t)
    })
}

func TestUrlInput_GetExpiresOn(t *testing.T) {
    t.Run("get expires_on empty", func(t *testing.T) {
        input := UrlInput{ExpiresOn: ""}

        if _, err := input.GetExpiresOn(); err != nil {
            t.Errorf("empty date should not give error")
        }
    })

    t.Run("get expires_on invalid", func(t *testing.T) {
        input := UrlInput{ExpiresOn: "Jan 2020"}

        if _, err := input.GetExpiresOn(); err == nil {
            t.Errorf("invalid date should give error")
        }
    })

    t.Run("get expires_on ok", func(t *testing.T) {
        input := UrlInput{ExpiresOn: "2020-01-01 00:00:00"}

        if _, err := input.GetExpiresOn(); err != nil {
            t.Errorf("valid date should not give error")
        }
    })
}
