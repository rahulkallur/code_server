#!/bin/bash

exts=(
    ms-python.python
)

for ext in "${exts[@]}";
do 
    echo "Installing ${ext} extension..."
    code --install-extension "${ext}"
done

# file:///home/coder/.local/share/code-server/extensions/extensions.json