package util

import (
	"bytes"
	"fiber-boilerplate/internal/defs"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type stringBlock struct{}

func (s stringBlock) FromInt(value interface{}) string {
	switch value.(type) {
	case int:
		return strconv.Itoa(value.(int))
	case int64:
		return strconv.FormatInt(value.(int64), 10)
	case time.Month:
		return s.FromInt(int(value.(time.Month)))

	case *int:
		return s.FromInt(*(value.(*int)))
	case *int64:
		return s.FromInt(*(value.(*int64)))
	case *time.Month:
		return s.FromInt(*(value.(*time.Month)))

	default:
		panic(defs.ErrNotImplemented)
	}
}

func (s stringBlock) Join(separator string, values ...interface{}) string {
	var buf bytes.Buffer

	last := len(values) - 1
	for i, v := range values {
		switch v.(type) {
		case int, int64, *int:
			buf.WriteString(s.FromInt(v))
		case bool:
			buf.WriteString(Bool.ToString(v.(bool)))
		case string, *string:
			buf.WriteString(v.(string))
		case []string:
			buf.WriteString(strings.Join(v.([]string), separator))

		case []interface{}:
			innerValues := v.([]interface{})
			innerLast := len(innerValues)
			for ii, vv := range innerValues {
				buf.WriteString(s.Join(separator, vv))
				if ii < innerLast {
					buf.WriteString(separator)
				}
			}

		default:
			panic(defs.ErrNotImplemented)
		}

		if i < last {
			buf.WriteString(separator)
		}
	}

	return buf.String()
}

func (s stringBlock) Concat(values ...interface{}) string {
	return s.Join("", values...)
}

func (s stringBlock) Words(values ...interface{}) string {
	return s.Join(" ", values...)
}

func (s stringBlock) List(values ...interface{}) string {
	return s.Join(", ", values...)
}

func (s stringBlock) SingleQuote(values ...interface{}) string {
	for i, v := range values {
		switch v.(type) {
		case string:
			values[i] = "'" + v.(string) + "'"
		case *string:
			values[i] = "'" + *(v.(*string)) + "'"
		default:
			panic(defs.ErrNotImplemented)
		}
	}
	return s.Join(", ", values...)
}

func (s stringBlock) DoubleQuote(values ...interface{}) string {
	for i, v := range values {
		switch v.(type) {
		case string:
			values[i] = "\"" + v.(string) + "\""
		case *string:
			values[i] = "\"" + *(v.(*string)) + "\""
		default:
			panic(defs.ErrNotImplemented)
		}
	}
	return s.Join(", ", values...)
}

func (s stringBlock) SliceTrim(sl []string, trim string) []string {
	for i, v := range sl {
		sl[i] = strings.Trim(v, trim)
	}
	return sl
}

func (s stringBlock) ObjectStringToSlice(str string) []string {
	sl := strings.Split(str, ",")
	for i := range sl {
		sl[i] = strings.TrimPrefix(strings.TrimSuffix(sl[i], "}"), "{")
		sl[i] = strings.Trim(sl[i], "\"")
	}
	return sl
}

func (s stringBlock) EmojiRemove(str string) string {
	emojiRegex := regexp.MustCompile(`[^\p{L}\p{N}\p{P}\p{Z}\p{Cf}\p{Cs}\s]`)
	return emojiRegex.ReplaceAllString(str, "")
}

func (s stringBlock) MatchPattern(pattern, str string) bool {
	patternRegex := regexp.MustCompile(pattern)
	return patternRegex.MatchString(str)
}
