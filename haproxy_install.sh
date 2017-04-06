#!/usr/bin/env bash

set -ex

cd /tmp
HAPROXY_VERSION=1.7.5
HAPROXY_MD5=ed84c80cb97852d2aa3161ed16c48a1c
curl -s -o haproxy.tar.gz http://www.haproxy.org/download/1.7/src/haproxy-$HAPROXY_VERSION.tar.gz
echo "$HAPROXY_MD5 haproxy.tar.gz" | md5sum -c
tar -zxf haproxy.tar.gz
cd /tmp/haproxy-$HAPROXY_VERSION
MAKE_OPTS='
TARGET=linux2628
USE_OPENSSL=1
USE_PCRE=1 PCREDIR=
USE_ZLIB=1
'
make all $MAKE_OPTS
make install-bin $MAKE_OPTS
mkdir -p /etc/haproxy
cp -R examples/errorfiles /etc/haproxy/errors
rm -rf /tmp/*
groupadd haproxy
useradd haproxy -g haproxy
