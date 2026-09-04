FROM golang:1.24-bookworm

RUN apt-get update && apt-get install -y \
    nginx \
    curl \
    git \
    ca-certificates \
    gnupg \
    lsb-release \
    && rm -rf /var/lib/apt/lists/*

RUN install -m 0755 -d /etc/apt/keyrings && \
    curl -fsSL https://download.docker.com/linux/debian/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg && \
    chmod a+r /etc/apt/keyrings/docker.gpg && \
    echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/debian $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null && \
    apt-get update && apt-get install -y docker-ce-cli && \
    rm -rf /var/lib/apt/lists/*

RUN curl -L --output /usr/local/bin/cloudflared https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-amd64 && \
    chmod +x /usr/local/bin/cloudflared

RUN mkdir -p /etc/nginx/sites-available /etc/nginx/sites-enabled /mini-paas/deployments

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o mini-paas-server .

RUN echo '#!/bin/bash\n\
service nginx start\n\
./mini-paas-server\n\
' > /app/entrypoint.sh && chmod +x /app/entrypoint.sh

EXPOSE 8080 80

ENTRYPOINT ["/app/entrypoint.sh"]