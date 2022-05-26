package assert

import (
	"io"
	"testing"

	"github.com/jgroeneveld/schema"
	"github.com/stretchr/testify/assert"
)

func MatchesJSONSchema(t *testing.T, s schema.Map, rc io.ReadCloser) {
	defer rc.Close()

	if err := schema.MatchJSON(s, rc); err != nil {
		assert.FailNow(t, "contents of ReadCloser did not match JSON schema", err.Error())
	}
}
