#!/usr/bin/env bash

set -u

exec 3<>/dev/tcp/localhost/8888
echo "helloworld" >&3
timeout 1 bash -c "cat <&3"
exec 3>&-
