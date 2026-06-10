#!/bin/bash

MAX_DEPTH=4
WORKSPACE_DIR=/home/coder
install_extension() {
    local ext="$1"
    echo "Installing ${ext} extension..."
    code-server --install-extension "${ext}"
}

if find "$WORKSPACE_DIR" -maxdepth "$MAX_DEPTH" -type f \( -name "package.json" -o -name "tsconfig.json" -o -name "*.js" -o -name "*.ts" \) -print -quit | grep -q .; then
    install_extension "xabikos.JavaScriptSnippets"
fi

if find . -maxdepth "$MAX_DEPTH" -type f \( -name "requirements.txt" -o -name "pyproject.toml" \) -print -quit | grep -q .; then
    install_extension ms-python.python
fi

if find . -maxdepth "$MAX_DEPTH" -type f \( -name "pom.xml" -o -name "build.gradle" -o -name "build.gradle.kts" \) -print -quit | grep -q .; then
    install_extension redhat.java
fi

if find . -maxdepth 4 -type f -name "go.mod" -print -quit | grep -q .; then
    echo "Found go.mod"
    install_extension "golang.Go"
fi

if find . -maxdepth "$MAX_DEPTH" -type f -name "composer.json" -print -quit | grep -q .; then
    install_extension bmewburn.vscode-intelephense-client
fi

if find . -maxdepth "$MAX_DEPTH" -type f \( -name "*.yaml" -o -name "*.yml" \) -print -quit | grep -q .; then
    install_extension redhat.vscode-yaml
    install_extension ms-kubernetes-tools.vscode-kubernetes-tools
fi

echo "Extension Detection Completed"