package common

import "regexp"

const (
    // DateLayout is standard yyyy-mm-dd format
    DateLayout = "2006-01-02 15:04:05"

    // ShortCodeLength is length of short code to be generated (5-12 chars)
    ShortCodeLength = 6
)

// ShortCodePattern is regex to check if a string looks like short code
var ShortCodeRegex, _ = regexp.Compile("^[a-zA-Z0-9]{4,12}$")
