#!/bin/bash
vite-ssg build

if [ "$?" -gt 0 ]; then
    echo "Prevent the fucking error temporarily!"
fi
