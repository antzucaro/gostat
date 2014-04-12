package qstr

import (
	"fmt"
	"html"
	"html/template"
	"regexp"
	"strconv"
	"strings"
)

var allColors = regexp.MustCompile(`\^(\d|x[\dA-Fa-f]{3})`)
var decColors = regexp.MustCompile(`\^(\d)`)
var hexColors = regexp.MustCompile(`\^x([\dA-Fa-f])([\dA-Fa-f])([\dA-Fa-f])`)

func hexToRGBSpanStr(r string, g string, b string) string {

	red, _   := strconv.ParseInt(fmt.Sprintf("%s%s", r, r), 16, 0)
	green, _ := strconv.ParseInt(fmt.Sprintf("%s%s", g, g), 16, 0)
	blue, _  := strconv.ParseInt(fmt.Sprintf("%s%s", b, b), 16, 0)

	return fmt.Sprintf("<span style=\"color:rgb(%d,%d,%d)\">", red, green, blue)
}

type QStr string

func (s *QStr) Stripped() string {
	return allColors.ReplaceAllString(string(*s), "")
}

func (s *QStr) HTML() template.HTML {
    // color representation by key for the "^n" format, where n is 0-9
    var decimalSpans = map[string]string{
        "^0": "<span style='color:rgb(128,128,128)'>",
        "^1": "<span style='color:rgb(255,0,0)'>",
        "^2": "<span style='color:rgb(51,255,0)'>",
        "^3": "<span style='color:rgb(255,255,0)'>",
        "^4": "<span style='color:rgb(51,102,255)'>",
        "^5": "<span style='color:rgb(51,255,255)'>",
        "^6": "<span style='color:rgb(255,51,102)'>",
        "^7": "<span style='color:rgb(255,255,255)'>",
        "^8": "<span style='color:rgb(153,153,153)'>",
        "^9": "<span style='color:rgb(128,128,128)'>",
    }

    // cast once to the string representation 'r'
    r := string(*s)

    // remove HTMl special characters
    r = html.EscapeString(r)

    // substitute matches of the form ^n, with n in 0..9
	matchedDecStrings := decColors.FindAllStringSubmatch(r, -1)
	for _, v := range matchedDecStrings {
		r = strings.Replace(r, v[0], decimalSpans[v[0]], 1)
	}

    // substitute matches of the form ^xrgb
    // with r, g, and b being hexadecimal digits
	matchedHexStrings := hexColors.FindAllStringSubmatch(r, -1)
	for _, v := range matchedHexStrings {
		r = strings.Replace(r, v[0], hexToRGBSpanStr(v[1], v[2], v[3]), 1)
	}

    // add the appropriate amount of closing spans
    for i := 0; i < (len(matchedDecStrings) + len(matchedHexStrings)); i++ {
        r = fmt.Sprintf("%s%s", r, "</span>")
    }

	return template.HTML(r)
}
