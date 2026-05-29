FROM codercom/code-server:latest

USER root

ENV DEBIAN_FRONTEND=noninteractive
ENV HOME=/home/coder
ENV PORT=3000
ENV WORKSPACE=/home/coder/project

# Install required packages
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        curl \
        git \
        jq \
        ca-certificates && \
    apt-get clean && \
    rm -rf \
        /var/lib/apt/lists/* \
        /tmp/* \
        /var/tmp/*

# Create required directories
RUN mkdir -p \
    /home/coder/.local/share/code-server/User \
    /home/coder/.local/share/code-server/extensions \
    /home/coder/.config/code-server \
    /home/coder/project

# Copy settings
COPY vscode-settings.json \
    /home/coder/.local/share/code-server/User/settings.json

# Copy config
COPY config/config.yaml \
    /home/coder/.config/code-server/config.yaml

# Copy startup script
COPY scripts/entrypoint.sh /usr/local/bin/entrypoint.sh

RUN chmod +x /usr/local/bin/entrypoint.sh

# Fix permissions
RUN chown -R coder:coder /home/coder

USER coder

WORKDIR /home/coder

EXPOSE 3000

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]