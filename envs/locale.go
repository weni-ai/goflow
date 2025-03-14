package envs

import (
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/text/language"
)

// Locale is the combination of a language and country, e.g. US English, Brazilian Portuguese
type Locale struct {
	Language Language
	Country  Country
}

// NewLocale creates a new locale
func NewLocale(language Language, country Country) Locale {
	return Locale{Language: language, Country: country}
}

// ToBCP47 returns the BCP47 code, e.g. en-US, pt, pt-BR
func (l Locale) ToBCP47() string {
	if l == NilLocale {
		return ""
	}

	lang, err := language.ParseBase(string(l.Language))
	if err != nil {
		return ""
	}
	code := lang.String()

	// not all languages have a 2-letter code
	if len(code) != 2 {
		return ""
	}

	if l.Country != NilCountry {
		code += "-" + string(l.Country)
	}
	return code
}

// FromBCP47 creates a Locale from a BCP47 code, e.g. "en-US" becomes English and US
func FromBCP47(code string) (Locale, error) {
	if code == "" {
		return NilLocale, nil
	}

	parts := strings.Split(code, "-")
	if len(parts) > 2 {
		return NilLocale, fmt.Errorf("invalid BCP47 code: %s", code)
	}

	tag, err := language.Parse(parts[0])
	if err != nil {
		return NilLocale, fmt.Errorf("invalid language code: %s", parts[0])
	}
	base, _ := tag.Base()
	lang := Language(base.ISO3())

	var country Country
	if len(parts) == 2 {
		country = Country(strings.ToUpper(parts[1]))
		if !country.IsValid() {
			return NilLocale, fmt.Errorf("invalid country code: %s", parts[1])
		}
	}

	return NewLocale(lang, country), nil
}

var NilLocale = Locale{}

// IsValid returns whether the country code is valid
func (c Country) IsValid() bool {
	if c == "" {
		return true
	}
	matched, _ := regexp.MatchString(`^[A-Z]{2}$`, string(c))
	return matched
}
