# 中间件容器启动

项目所用到的容器启动说明，参考下面的方式用命令行手动启动，也可以使用 compose 启动。

## mysql5.7
创建数据目录后启动 mysql，注意修改数据库密码
```bash
mkdir -p ~/data/mysql/data

podman run --name mysql8 \
        -p 3306:3306 \
        -v ~/data/mysql/data:/var/lib/mysql \
        -e MYSQL_ROOT_PASSWORD="xxxxxxx" \
        -e MYSQL_DATABASE="go-mall" \
        -e MYSQL_USER="go-mall" \
        -e MYSQL_PASSWORD="xxxxxx" \
        -d mysql:5.7
```

## redis
创建数据目录和配置文件后启动 redis
```bash
mkdir -p ~/data/redis/{data,conf}

cat > ~/data/redis/conf/redis.conf <<EOF
protected-mode no
aclfile /etc/redis/acl.conf
appendonly yes
EOF

echo 'user wuvikr on >[your_password] ~* +@all' > ~/data/redis/conf/acl.conf

podman run --name redis7 -d \
    -p 6379:6379 \
    -v ~/data/redis/data:/data \
    -v ~/data/redis/conf:/etc/redis \
    redis:7 redis-server /etc/redis/redis.conf
```