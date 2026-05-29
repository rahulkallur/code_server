#!/bin/bash

set -e

export HOME=/home/coder
export WORKSPACE=/home/coder/project

mkdir -p "$WORKSPACE"

echo "Installing Visual Studio Code extensions..."

# Install Visual Studio Code extensions
extensions=(
ms-python.python
redhat.java
golang.Go
)

for ext in "${extensions[@]}";
do 
    echo "Installing ${ext} extension..."
    code-server --install-extension "${ext}"
done


echo "Starting code-server..."
exec code-server --bind-addr 0.0.0.0:3000 --auth none "$WORKSPACE"