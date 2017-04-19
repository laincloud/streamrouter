package backend

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHaproxyBackend_ParseStreamApp(t *testing.T) {
	const expected = `
listen testStreamApp_testStreamProc1_8080
  mode tcp
  bind :8080
  server testStreamProc1_8080_1 192.168.0.1:9080 check
  server testStreamProc1_8080_2 192.168.0.2:9080 check

listen testStreamApp_testStreamProc1_8081
  mode tcp
  bind :8081
  server testStreamProc1_8081_1 192.168.0.1:9081 check
  server testStreamProc1_8081_2 192.168.0.2:9081 check

listen testStreamApp_testStreamProc2_8180
  mode tcp
  bind :8180
  server testStreamProc2_8180_1 192.168.1.1:9180 check
  server testStreamProc2_8180_2 192.168.1.2:9180 check

listen testStreamApp_testStreamProc2_8181
  mode tcp
  bind :8181
  server testStreamProc2_8181_1 192.168.1.1:9181 check
  server testStreamProc2_8181_2 192.168.1.2:9181 check
`
	hb := HaproxyBackend{}
	result, err := hb.parseStreamApp(testApp)
	assert.Nil(t, err)
	assert.Equal(t, expected, string(result))
}
