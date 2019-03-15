#!/usr/bin/env python
# coding: utf-8

import sys
import json

zoneids = [
]

for line in sys.stdin:
    try:
        line = line.rstrip()
        fields = line.split("\t")
        logs = json.loads(fields[2])
        print '%s' %(line)

    except:
        import traceback
        traceback.print_exc()
        continue
