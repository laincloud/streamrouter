#!/usr/bin/env bash

set -ex

cd /tmp

NGINX_VERSION=1.11.4
curl -s -o nginx-$NGINX_VERSION.tar.gz http://nginx.org/download/nginx-$NGINX_VERSION.tar.gz
tar -zxf nginx-$NGINX_VERSION.tar.gz
git clone https://github.com/cybercom-finland/ngx_stream_upstream_check_module.git
cd /tmp/ngx_stream_upstream_check_module
git checkout fixes
cd /tmp/nginx-$NGINX_VERSION
patch -p0 < /tmp/ngx_stream_upstream_check_module/patch-1.11.x.patch

./configure --conf-path=/etc/nginx/nginx.conf --pid-path=/var/run/nginx.pid --sbin-path=/usr/local/sbin/nginx --error-log-path=/var/log/nginx/error.log --http-log-path=/var/log/nginx/access.log --with-stream --add-module=/tmp/ngx_stream_upstream_check_module
make
make install

groupadd nginx
useradd nginx -g nginx