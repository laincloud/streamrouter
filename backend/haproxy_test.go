package backend

import (
	"testing"

	"github.com/laincloud/streamrouter/model"
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
  option tcp-check
  tcp-check send-binary ping\n
  tcp-check expect binary pong\n
  server testStreamProc1_8081_1 192.168.0.1:9081 check
  server testStreamProc1_8081_2 192.168.0.2:9081 check

listen testStreamApp_testStreamProc2_8180
  mode tcp
  bind :8180
  option tcp-check
  tcp-check send-binary GET\ /\ HTTP/1.0\r\n
  tcp-check expect binary (2..|3..)
  server testStreamProc2_8180_1 192.168.1.1:9180 check
  server testStreamProc2_8180_2 192.168.1.2:9180 check

listen testStreamApp_testStreamProc2_8181
  mode tcp
  bind :8181
  server testStreamProc2_8181_1 192.168.1.1:9181 check
  server testStreamProc2_8181_2 192.168.1.2:9181 check
`
	testApp := model.StreamApp{
		Name: "testStreamApp",
		StreamProcs: []model.StreamProc{
			{
				Name: "testStreamProc1",
				Upstreams: []model.StreamUpstream{
					{
						InstanceNo: 1,
						Host:       "192.168.0.1",
					},
					{
						InstanceNo: 2,
						Host:       "192.168.0.2",
					},
				},
				Services: []model.StreamService{
					{
						ListenPort:   8080,
						UpstreamPort: 9080,
					},
					{
						ListenPort:   8081,
						UpstreamPort: 9081,
						Send:         `ping\n`,
						Expect:       `pong\n`,
					},
				},
			},
			{
				Name: "testStreamProc2",
				Upstreams: []model.StreamUpstream{
					{
						InstanceNo: 1,
						Host:       "192.168.1.1",
					},
					{
						InstanceNo: 2,
						Host:       "192.168.1.2",
					},
				},
				Services: []model.StreamService{
					{
						ListenPort:   8180,
						UpstreamPort: 9180,
						Send:         `GET\ /\ HTTP/1.0\r\n`,
						Expect:       `(2..|3..)`,
					},
					{
						ListenPort:   8181,
						UpstreamPort: 9181,
						Send:         ``,
						Expect:       `pong\n`,
					},
				},
			},
		},
	}
	hb := HaproxyBackend{}
	result, err := hb.parseStreamApp(testApp)
	assert.Nil(t, err)
	assert.Equal(t, expected, string(result))
}
