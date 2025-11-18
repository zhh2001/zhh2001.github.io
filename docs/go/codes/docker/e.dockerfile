# 基础镜像
FROM openjdk:11.0-jre-buster

# 设定时区
ENV TZ=Asia/Shanghai

RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# 拷贝 jar 包
COPY docker-demo.jar /app.jar

# 入口
ENTRYPOINT ["java", "-jar", "/app.jar"]