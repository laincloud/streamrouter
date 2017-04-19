package backend

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNginxBackend_ParseStreamApp(t *testing.T) {
	const expected = `
server {
    listen 8080;
    proxy_pass testStreamApp_testStreamProc1_8080_9080;
    access_log /var/log/nginx/testStreamApp.log stream;
}

upstream testStreamApp_testStreamProc1_8080_9080 {
    server 192.168.0.1:9080;
    server 192.168.0.2:9080;
    check interval=3000 rise=2 fall=5 timeout=1000 type=http;
}

server {
    listen 8081;
    proxy_pass testStreamApp_testStreamProc1_8081_9081;
    access_log /var/log/nginx/testStreamApp.log stream;
}

upstream testStreamApp_testStreamProc1_8081_9081 {
    server 192.168.0.1:9081;
    server 192.168.0.2:9081;
    check interval=3000 rise=2 fall=5 timeout=1000 type=http;
}

server {
    listen 8180;
    proxy_pass testStreamApp_testStreamProc2_8180_9180;
    access_log /var/log/nginx/testStreamApp.log stream;
}

upstream testStreamApp_testStreamProc2_8180_9180 {
    server 192.168.1.1:9180;
    server 192.168.1.2:9180;
    check interval=3000 rise=2 fall=5 timeout=1000 type=http;
}

server {
    listen 8181;
    proxy_pass testStreamApp_testStreamProc2_8181_9181;
    access_log /var/log/nginx/testStreamApp.log stream;
}

upstream testStreamApp_testStreamProc2_8181_9181 {
    server 192.168.1.1:9181;
    server 192.168.1.2:9181;
    check interval=3000 rise=2 fall=5 timeout=1000 type=http;
}
`
	nb := NginxBackend{}
	result, err := nb.parseStreamApp(testApp)
	assert.Nil(t, err)
	assert.Equal(t, expected, string(result))
}
