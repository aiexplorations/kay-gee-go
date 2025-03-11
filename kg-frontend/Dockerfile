FROM nginx:alpine

# Copy static files
COPY public/ /usr/share/nginx/html/

# Copy nginx configuration
RUN echo 'server { \
    listen 8080; \
    server_name localhost; \
    root /usr/share/nginx/html; \
    index index.html; \
    \
    # Add CORS headers \
    add_header "Access-Control-Allow-Origin" "*"; \
    add_header "Access-Control-Allow-Methods" "GET, POST, OPTIONS"; \
    add_header "Access-Control-Allow-Headers" "DNT,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Range,Authorization"; \
    add_header "Access-Control-Expose-Headers" "Content-Length,Content-Range"; \
    \
    location / { \
        try_files $uri $uri/ /index.html; \
    } \
    \
    # Proxy for Neo4j Browser API \
    location /db/ { \
        proxy_pass http://neo4j:7474/db/; \
        proxy_http_version 1.1; \
        proxy_set_header Host $host; \
        proxy_set_header X-Real-IP $remote_addr; \
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for; \
        proxy_set_header X-Forwarded-Proto $scheme; \
    } \
}' > /etc/nginx/conf.d/default.conf

# Expose port
EXPOSE 8080

# Start nginx
CMD ["nginx", "-g", "daemon off;"] 