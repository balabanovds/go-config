package httputils_test

import (
	"strings"
	"testing"

	"github.com/balabanovds/goutils/httputils"
	"github.com/stretchr/testify/require"
)

func TestSplitPath(t *testing.T) {
	tests := []struct {
		path     string
		wantTail string
		wantHead string
	}{
		{"", "/", ""},
		{"/", "/", ""},
		{"foo", "/", "foo"},
		{"/foo/bar", "/bar", "foo"},
		{"/foo/bar/", "/bar", "foo"},
	}

	for _, tst := range tests {
		t.Run(tst.path, func(t *testing.T) {
			head, tail := httputils.SplitPath(tst.path)
			require.Equal(t, tst.wantTail, tail, "tail check")
			require.Equal(t, tst.wantHead, head, "head check")
		})
	}
}

func TestParseID(t *testing.T) {
	tests := []struct {
		name          string
		path          string
		wantID        int
		wantTail      string
		wantErrSubstr string
	}{
		{
			name:          "simple ID",
			path:          "/123",
			wantID:        123,
			wantTail:      "",
			wantErrSubstr: "",
		},
		{
			name:          "error",
			path:          "/123a",
			wantID:        0,
			wantErrSubstr: "invalid syntax",
		},
		{
			name:          "middle path ID",
			path:          "/foo/123/bar",
			wantID:        123,
			wantTail:      "/bar",
			wantErrSubstr: "",
		},
	}

	for _, tst := range tests {
		t.Run(tst.name, func(t *testing.T) {
			id, tail, err := httputils.ParseIntID(tst.path)
			if tst.wantErrSubstr != "" && err != nil {
				require.True(t, strings.Contains(err.Error(), tst.wantErrSubstr))
				return
			}
			require.NoError(t, err)
			require.Equal(t, tst.wantID, id)
			require.Equal(t, tst.wantTail, tail)
		})
	}
}
