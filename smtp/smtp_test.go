package smtp

import (
	"github.com/emersion/go-message/mail"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateSenderSlice(t *testing.T) {
	a := []*mail.Address{{Address: "a@example.com"}, {Address: "b@example.com"}}
	b := []*mail.Address{{Address: "a@example.com"}, {Address: "c@example.com"}}
	c := []*mail.Address{{Address: "a@example.com"}, {Address: "d@example.com"}}
	assert.Equal(t, []string{
		"a@example.com",
		"b@example.com",
		"a@example.com",
		"c@example.com",
		"a@example.com",
		"d@example.com",
	}, CreateSenderSlice(a, b, c))
}
