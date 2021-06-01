#!/bin/bash

curl \
    -d '{"id":1, "ran":"10.0.10.1", "upf":"10.0.5.1", "slice_match": {"uteid": 0, "dteid": 1}, "qos_match": {"dscp":16}}'  \
    http://localhost:9090/ntn-session/new-session
echo ""