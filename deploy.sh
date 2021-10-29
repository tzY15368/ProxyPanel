# docker
curl -fsSL https://get.docker.com | bash -s docker

sudo systemctl enable docker
sudo systemctl daemon-reload
sudo systemctl restart docker

# 允许运行普通账号直接执行 docker 命令
sudo groupadd docker
sudo usermod -aG docker $USER
newgrp docker

# docker-compose
sudo curl -L "https://github.com/docker/compose/releases/download/1.23.1/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
docker-compose --version

docker-compose -f lazarus-worker.yml up