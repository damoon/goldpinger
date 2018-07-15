#!/bin/bash

CompileDaemon \
    -pattern "(.+\\.go|.+\\.elm|.+\\.css|.+\\.yaml|.+\\.yml)$" \
    -build="make deploy"
