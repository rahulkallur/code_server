#!/bin/bash

SCRIPT="/usr/local/bin/install-extensions.sh"

run_script() {
    echo "Running extension detection script..."
    "$SCRIPT"
}

run_script

while true; do
    inotifywait -e create -e modify -e delete -e moved_to .

    echo "Extension Detection Script Updated, Re-running..."
    sleep 300

    run_script
done