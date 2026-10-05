# Java 25 JRE，基础系统为 Ubuntu 24.04
FROM eclipse-temurin:25-jre-noble

WORKDIR /app
ENV TZ=Asia/Shanghai
RUN ln -snf "/usr/share/zoneinfo/$TZ" /etc/localtime \
    && printf '%s\n' "$TZ" > /etc/timezone

# JAR 需提前构建，并放在构建上下文中
COPY docker-demo.jar app.jar

# 应用本身需要监听此端口
EXPOSE 8080
ENTRYPOINT ["java", "-jar", "app.jar"]
