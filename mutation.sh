#!/bin/bash
curl -X POST -d 'mutation { user(id: 100){ id, address{ state, city }}}' http://localhost:8080/
