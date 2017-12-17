#!/bin/bash
curl -X POST -H "Content-Type: application/json" -d '{ id(id: 100), name }' http://localhost:8080/
