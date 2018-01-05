package middlewares

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractBearer(t *testing.T) {
	req, err := http.NewRequest("GET", "/tokens", nil)
	if err != nil {
		t.Fatal(err)
	}
	raw := "thisisabearertokenthatshouldbefound"
	req.Header.Add("Authentication", fmt.Sprintf("Bearer %s", raw))
	tval, err := bearer(req.Header)
	assert.Nil(t, err, "should be nil")
	assert.Equal(t, raw, tval, "should be equal")

	t.Run("Empty bearer", func(t *testing.T) {
		req.Header.Set("Authentication", "Bearer ")
		tval, err := bearer(req.Header)
		assert.NotNil(t, err, "should not be nil")
		assert.Empty(t, tval, "should be empty")
		assert.Equal(t, "empty bearer token", err.Error(), "should be equal")
	})

	t.Run("Empty authentication", func(t *testing.T) {
		req.Header.Del("Authentication")
		tval, err := bearer(req.Header)
		assert.NotNil(t, err, "should not be nil")
		assert.Empty(t, tval, "should be empty")
		assert.Equal(t, "no authentication header found", err.Error(), "should be equal")
	})
}
