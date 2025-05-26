#!/bin/bash

# 设置颜色
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # 无颜色

echo -e "${GREEN}DDPay 环境配置向导${NC}"
echo "=============================="
echo -e "${YELLOW}注意：此脚本将帮助您设置环境变量，用于Docker部署${NC}"
echo ""

# 检查是否存在.env文件
if [ -f .env ]; then
  echo -e "${YELLOW}检测到已存在.env文件，是否要覆盖？(y/n)${NC}"
  read -r overwrite
  if [[ "$overwrite" != "y" && "$overwrite" != "Y" ]]; then
    echo "保留现有.env文件，退出配置。"
    exit 0
  fi
fi

# 复制示例文件
cp .env.example .env

echo -e "${GREEN}请按照提示输入配置信息：${NC}"

# 服务器配置
read -p "服务器端口 (默认: 2900): " SERVER_PORT
SERVER_PORT=${SERVER_PORT:-2900}

# 数据库配置
read -p "数据库主机 (默认: mysql): " DB_HOST
DB_HOST=${DB_HOST:-mysql}

read -p "数据库端口 (默认: 3306): " DB_PORT
DB_PORT=${DB_PORT:-3306}

read -p "数据库用户名 (默认: ddpay): " DB_USER
DB_USER=${DB_USER:-ddpay}

read -p "数据库密码: " DB_PASSWORD
if [ -z "$DB_PASSWORD" ]; then
  echo -e "${YELLOW}警告: 数据库密码为空${NC}"
fi

read -p "数据库名称 (默认: ddpay): " DB_NAME
DB_NAME=${DB_NAME:-ddpay}

# Redis配置
read -p "Redis主机 (默认: redis): " REDIS_HOST
REDIS_HOST=${REDIS_HOST:-redis}

read -p "Redis端口 (默认: 6379): " REDIS_PORT
REDIS_PORT=${REDIS_PORT:-6379}

read -p "Redis密码 (默认为空): " REDIS_PASSWORD

# Telegram配置
read -p "Telegram Bot Token: " TELEGRAM_BOT_TOKEN

# 区块链配置
read -p "EVM助记词: " EVM_MNEMONIC
read -p "EVM RPC地址: " EVM_RPC
read -p "Solana GRPC地址 (默认: https://solana-yellowstone-grpc.publicnode.com:443): " SOLANA_GRPC
SOLANA_GRPC=${SOLANA_GRPC:-https://solana-yellowstone-grpc.publicnode.com:443}

# 更新.env文件
sed -i '' "s|SERVER_PORT=.*|SERVER_PORT=$SERVER_PORT|g" .env
sed -i '' "s|DB_HOST=.*|DB_HOST=$DB_HOST|g" .env
sed -i '' "s|DB_PORT=.*|DB_PORT=$DB_PORT|g" .env
sed -i '' "s|DB_USER=.*|DB_USER=$DB_USER|g" .env
sed -i '' "s|DB_PASSWORD=.*|DB_PASSWORD=$DB_PASSWORD|g" .env
sed -i '' "s|DB_NAME=.*|DB_NAME=$DB_NAME|g" .env
sed -i '' "s|REDIS_HOST=.*|REDIS_HOST=$REDIS_HOST|g" .env
sed -i '' "s|REDIS_PORT=.*|REDIS_PORT=$REDIS_PORT|g" .env
sed -i '' "s|REDIS_PASSWORD=.*|REDIS_PASSWORD=$REDIS_PASSWORD|g" .env
sed -i '' "s|TELEGRAM_BOT_TOKEN=.*|TELEGRAM_BOT_TOKEN=$TELEGRAM_BOT_TOKEN|g" .env
sed -i '' "s|EVM_MNEMONIC=.*|EVM_MNEMONIC=$EVM_MNEMONIC|g" .env
sed -i '' "s|EVM_RPC=.*|EVM_RPC=$EVM_RPC|g" .env
sed -i '' "s|SOLANA_GRPC=.*|SOLANA_GRPC=$SOLANA_GRPC|g" .env

echo -e "${GREEN}环境变量配置完成！${NC}"
echo "您可以通过编辑.env文件进一步调整配置。"
echo -e "${YELLOW}接下来可以运行 'docker-compose up -d' 来启动服务。${NC}"

# 设置执行权限
chmod +x setup-env.sh 