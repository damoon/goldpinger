#!/bin/bash

CompileDaemon \
    -pattern "(.+\\.go|.+\\.css|.+\\.yaml|.+\\.yml)$" \
    -build="make deploy"
