FROM golang:1.24.1-alpine3.21

ARG APP_NAME=""
ARG ENV_TYPE=""

# 设置工作目录
WORKDIR /home/www

# 拷贝项目到工作目录
COPY  . /home/www

RUN echo ${APP_NAME} ${ENV_TYPE} && go version && go mod tidy

# 编译构建
RUN go build -o /home/www/target/server /home/www/app/${APP_NAME}/cmd/server && chmod +x /home/www/target/server

# 复制配置文件
RUN cp /home/www/app/${APP_NAME}/configs/config_${ENV_TYPE}.toml /home/www/target/config.toml


EXPOSE 80

CMD ["/home/www/target/server", "-conf", "/home/www/target/config.toml"]