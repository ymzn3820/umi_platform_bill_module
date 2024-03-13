# 使用go1.20.7版本作为基础镜像
FROM golang:1.20.7

# 设置工作目录
WORKDIR /app

# 设置国内代理，加速go module下载
ENV GOPROXY=https://goproxy.cn,direct

# 复制go.mod和go.sum文件
COPY go.mod ./
COPY go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码到容器中
COPY . .

# 构建应用程序
RUN go build -o /server-bill

# 暴露28070端口
EXPOSE 28070
