package generator

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShortGenerator_GenerateShortLink(t *testing.T) {
	t.Parallel()
	sg := NewShortGenerator()

	tests := []struct {
		name string
	}{
		{name: "valid"},
		{name: "valid again"},
		{name: "valid third"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := sg.GenerateShortLink()
			require.NoError(t, err)
			assert.Len(t, got, ShortLength)
			for _, r := range got {
				assert.True(t, strings.ContainsRune(charset, r), "char %q not in charset", r)
			}
		})
	}
}
