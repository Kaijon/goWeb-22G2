#!/bin/sh

storage=$1

if [ "$#" -lt 1 ]; then
    logger -s "Wrong Arguments, please specify getac or flash"
    exit 1
fi

if [ "$storage" = "getac" ]; then
    logger -s "cleanup /tmp/tmp_daemon"
    rm -rf /tmp/tmp_daemon
    exit 0
fi

if [ "$storage" = "flash" ]; then
    logger -s "cleanup /tmp/tmp_flash"
    rm -rf /tmp/tmp_flash
    exit 0
fi

logger -s "Arguments not supportted"

