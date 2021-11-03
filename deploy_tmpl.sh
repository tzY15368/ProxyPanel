# # docker
# curl -fsSL https://get.docker.com | bash -s docker

# sudo systemctl enable docker
# sudo systemctl daemon-reload
# sudo systemctl restart docker

# # 允许运行普通账号直接执行 docker 命令
# sudo groupadd docker
# sudo usermod -aG docker $USER
# newgrp docker

# # docker-compose
# sudo curl -L "https://github.com/docker/compose/releases/download/1.23.1/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
# sudo chmod +x /usr/local/bin/docker-compose
# docker-compose --version

# docker-compose -f lazarus-worker.yml up

sudo apt install v2ray nginx cron

# v2ray
sudo systemctl enable v2ray
sudo systemctl restart v2ray

# nginx
sudo systemctl enable nginx
sudo systemctl restart nginx

# acme
echo "sleeping for 61 seconds, waiting for DNS TTL to come to effect"
sleep 61
echo "end of sleep"


curl -sL https://get.acme.sh | sh -s email={{certEmail}}
source ~/.bashrc
~/.acme.sh/acme.sh --upgrade --auto-upgrade
~/.acme.sh/acme.sh --set-default-ca --server letsencrypt
~/.acme.sh/acme.sh --issue -d {{domain}} --keylength ec-256 --pre-hook "systemctl stop nginx" --post-hook "systemctl restart nginx" --standalone
CERT_FILE="/etc/v2ray/{{domain}}.pem"
KEY_FILE="/etc/v2ray/{{domain}}.key"
~/.acme.sh/acme.sh --install-cert -d {{domain}} --ecc \
    --key-file $KEY_FILE \
    --fullchain-file $CERT_FILE \
    --reloadcmd "service nginx reload"
