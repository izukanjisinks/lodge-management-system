#!/bin/bash
export $(cat .env | xargs)
./bin/server.exe
