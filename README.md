# 🛠️ GitHub Webhook 自动部署脚本（Go 实现）

用于接收 GitHub Webhook 请求，并在验证签名成功后执行本地脚本（如自动部署）。

---

## 功能特性

- 🔒 安全的 HMAC-SHA256 签名验证
- ⚡️ 并发执行脚本
- 🛠️ 调试模式便于排查问题
- 📝 YAML 配置文件
- 🚀 支持多个 Webhook 端点

## 快速开始

### 环境要求

- Go 1.16+ 环境
- Bash 环境

### 安装步骤

1. 克隆本仓库

2. 编译可执行文件：

   ```bash
   go build -o webhook-server
   ```

### 配置说明

创建 `config.yaml` 配置文件（参考下方示例），并将脚本放在指定位置

示例 `config.yaml`：

```yaml
port: 5599
debug: false

service:
  - name: "博客"
    path: "/hook/blog" 
    secret: "./secret/blog_secret.txt"
    script: "./script/blog_script.sh"
```

### 运行服务

```bash
./webhook-server
```

## 配置选项

|  字段   |            说明             |
| :-----: | :-------------------------: |
|  port   |     服务端口 (1-65535)      |
|  debug  |      是否启用调试日志       |
| service | 需要处理的 Webhook 端点列表 |

### 服务配置

每个服务需要配置：

- `name`: 服务名称
- `path`: Webhook 端点路径
- `secret`: 存放 GitHub Webhook 密钥的文件路径
- `script`: 验证通过后要执行的脚本路径

## 安全机制

服务端会使用 GitHub 的 HMAC-SHA256 签名验证所有收到的 Webhook

## 示例脚本

将可执行脚本放在配置指定的位置。示例 (`blog_script.sh`)：

```bash
#!/bin/bash
cd /home/jeanhua/my_blog/Blog/ &&
git fetch origin main &&
git reset --hard origin/main &&
npm run build
```

## 许可证

MIT 开源协议
