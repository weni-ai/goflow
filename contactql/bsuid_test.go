package contactql_test

import (
	"testing"

	"github.com/nyaruka/goflow/contactql"
	"github.com/stretchr/testify/assert"
)

func TestIsWhatsAppBSUIDPath(t *testing.T) {
	assert.True(t, contactql.IsWhatsAppBSUIDPath("BR.1583576196772655"))
	assert.True(t, contactql.IsWhatsAppBSUIDPath("br.1583576196772655"))
	assert.True(t, contactql.IsWhatsAppBSUIDPath("US.ENT.11815799212886844830"))
	assert.True(t, contactql.IsWhatsAppBSUIDPath("us.ent.11815799212886844830"))
	assert.True(t, contactql.IsWhatsAppBSUIDPath("MX.abc123XYZ"))

	assert.False(t, contactql.IsWhatsAppBSUIDPath("558231426933"))
	assert.False(t, contactql.IsWhatsAppBSUIDPath("+558231426933"))
	assert.False(t, contactql.IsWhatsAppBSUIDPath("BR1583576196772655"))
	assert.False(t, contactql.IsWhatsAppBSUIDPath(""))
	assert.False(t, contactql.IsWhatsAppBSUIDPath("B.123"))
}
