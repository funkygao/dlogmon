#!/usr/bin/env python
#encoding=utf-8
'''mapper for dlog-amf'''

import sys

while True:
    line = sys.stdin.readline()
    if not line:
        break

    print 'got here:', line,
    sys.stdout.flush()

