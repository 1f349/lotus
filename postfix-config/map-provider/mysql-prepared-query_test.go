package map_provider

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	testQuery       = "SELECT aliasMap.goto FROM aliasMap,aliasdomainMap WHERE aliasdomainMap.domain='%d' AND aliasMap.address = CONCAT('%u', '@', aliasdomainMap.goto) AND aliasMap.active > 0 AND aliasdomainMap.active > 0"
	testQueryRaw    = "SELECT aliasMap.goto FROM aliasMap,aliasdomainMap WHERE aliasdomainMap.domain='' AND aliasMap.address = CONCAT('', '@', aliasdomainMap.goto) AND aliasMap.active > 0 AND aliasdomainMap.active > 0"
	testQueryFormat = "SELECT aliasMap.goto FROM aliasMap,aliasdomainMap WHERE aliasdomainMap.domain='example.com' AND aliasMap.address = CONCAT('test', '@', aliasdomainMap.goto) AND aliasMap.active > 0 AND aliasdomainMap.active > 0"
)

func TestNewPreparedQuery(t *testing.T) {
	query, err := NewPreparedQuery(testQuery)
	assert.NoError(t, err)
	assert.Equal(t, PreparedQuery{
		raw: testQueryRaw,
		params: map[int]byte{
			79:  'd',
			112: 'u',
		},
	}, *query)
}

func TestPreparedQuery_Format(t *testing.T) {
	query := &PreparedQuery{
		raw: testQueryRaw,
		params: map[int]byte{
			79:  'd',
			112: 'u',
		},
	}
	format, err := query.Format(map[byte]string{
		'd': "example.com",
		'u': "test",
	})
	assert.NoError(t, err)
	assert.Equal(t, testQueryFormat, format)
}
