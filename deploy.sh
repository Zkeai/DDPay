#!/bin/bash

# 设置颜色
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # 无颜色

echo -e "${GREEN}DDPay 部署脚本${NC}"
echo "=============================="

# 检查Docker是否运行
if ! docker info > /dev/null 2>&1; then
  echo -e "${YELLOW}错误: Docker未运行，请先启动Docker${NC}"
  exit 1
fi

# 检查docker-compose是否安装
if ! command -v docker-compose &> /dev/null; then
  echo -e "${YELLOW}错误: docker-compose未安装${NC}"
  exit 1
fi

# 检查配置文件
if [ ! -f "etc/config.yaml" ]; then
  echo -e "${YELLOW}错误: 未找到配置文件 etc/config.yaml${NC}"
  exit 1
fi

# 检查环境变量文件
if [ ! -f ".env" ] && [ -f ".env.example" ]; then
  echo -e "${YELLOW}未找到.env文件，将复制.env.example${NC}"
  cp .env.example .env
  echo -e "${GREEN}已创建.env文件，请根据需要修改${NC}"
fi

# 创建必要的目录
mkdir -p log
echo -e "${GREEN}创建日志目录: log/${NC}"

# 构建镜像
echo -e "${GREEN}开始构建Docker镜像...${NC}"
docker-compose build

# 启动服务
echo -e "${GREEN}启动服务...${NC}"
docker-compose up -d

# 显示服务状态
echo -e "${GREEN}服务状态:${NC}"
docker-compose ps

echo -e "${GREEN}部署完成!${NC}"
echo -e "您可以使用以下命令查看日志:"
echo -e "${YELLOW}docker-compose logs -f http-server${NC} - 查看HTTP服务日志"
echo -e "${YELLOW}docker-compose logs -f cron-server${NC} - 查看定时任务服务日志"
echo -e "${YELLOW}docker-compose logs -f wallet-server${NC} - 查看钱包服务日志" 