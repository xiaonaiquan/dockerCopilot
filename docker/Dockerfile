FROM nginx:alpine

WORKDIR /app

# 拷贝后端
COPY ../dockerCopilot /app/dockerCopilot
COPY ../etc/. /app/etc
COPY ../start.sh /app/start.sh

# 拷贝前端静态资源
COPY ../front/. /usr/share/nginx/html

# 覆盖 nginx 配置
COPY ../nginx.conf /etc/nginx/nginx.conf

# 设置环境变量
ENV secretKey="" \
    DOCKER_HOST="unix:///var/run/docker.sock" \
    BACKUP_DIR="/data/backups" \
    TZ="Asia/Shanghai"
    UPLOAD_DIR="/data/uploads"

# 暴露端口
EXPOSE 80

VOLUME ["/data"]

# 启动后端 + nginx
CMD ["/bin/sh", "-c", "/app/start.sh & nginx -g 'daemon off;'"]