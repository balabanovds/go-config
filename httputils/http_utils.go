package httputils

import (
	"errors"
	gopath "path"
	"regexp"
	"strconv"
	"strings"
)

var reDigits *regexp.Regexp

func init() {
	reDigits = regexp.MustCompile(`\d+/*`)
}

func SplitPath(path string) (head, tail string) {
	path = gopath.Clean("/" + path)
	i := strings.Index(path[1:], "/") + 1
	if i <= 0 {
		return path[1:], "/"
	}

	return path[1:i], path[i:]
}

func ParseIntID(path string) (value int, tail string, err error) {
	path = gopath.Clean("/" + path)
	loc := reDigits.FindStringIndex(path)
	if loc == nil {
		return 0, "", errors.New("int ID not found in path")
	}

	if len(loc) != 2 {
		return
	}

	value, err = strconv.Atoi(path[loc[0] : loc[1]-1])
	tail = path[loc[1]:]
	return
}
