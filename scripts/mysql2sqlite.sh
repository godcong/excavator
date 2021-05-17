#!/bin/bash

mysql2sqlite -f ./exc.db -d excavator -u root -p 111111 -h "localhost" -V --use-buffered-cursors
