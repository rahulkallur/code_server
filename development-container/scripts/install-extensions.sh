#!/bin/bash

# This script is used to install VS Code extensions in the development container.
# It is executed during the Docker build process.
# Make sure to update the list of extensions as needed.
# The extensions are installed using the `code-server` command line tool.
# All the extensions insalled are open source and can be found on the Visual Studio Code Marketplace.

# List of extension to install

extensions=(
ms-python.python
redhat.java
)

for ext in "${extensions[@]}";
do 
    echo "Installing ${ext} extension..."
    code-server --install-extension "${ext}"
done