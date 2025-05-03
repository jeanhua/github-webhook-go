# 🛠️ GitHub Webhook 自动部署脚本（Go 实现）

用于接收 GitHub Webhook 请求，并在验证签名成功后执行本地脚本（如自动部署）。

---

## ✅ 功能特性

- 接收 GitHub 的 Webhook 请求
- 验证请求签名以确保安全性（HMAC SHA-256）
- 若签名正确，运行本地 shell 脚本（默认为 `./action.sh`）
- 支持自定义端口启动服务
- 支持 debug 模式输出日志与错误请求体

---

## 📦 依赖文件结构

```bash
github-webhook-go:.
│  action.sh		# 要执行的脚本
│  build_script.bat	# go构建脚本，用于自动构建linux程序
│  main.go			# 主程序代码
│  README.md
│
└─build				# 构建结果
        action.sh
        github-hook	# linux可执行程序
```

---

## ⚙️ 使用方法

### 1. 准备 `action.sh`

创建一个名为 `action.sh` 的脚本文件，内容为你希望在接收到推送事件后执行的操作。例如：

```bash
#!/bin/bash
echo "Pulling latest code..."
cd /path/to/repo && git pull origin main
echo "Restarting service..."
systemctl restart myapp
```

记得赋予执行权限：

```bash
chmod +x action.sh
```

---

### 2. 获取并设置 GitHub Webhook Secret

在 GitHub 仓库的 **Settings > Webhooks** 页面中添加一个新的 Webhook：

- **Payload URL**: `http://yourdomain.com:5599/hook`
- **Content type**: `application/json`
- **Secret**: 输入你将在程序运行时设置的密钥（建议复杂度 ≥ 10 字符）

---

### 3. 运行服务

#### 正常运行：

```bash
go run main.go -p 5599
```

或者编译后运行：

```bash
go build -o webhook-go main.go
./webhook-go -p 8080
```

> 如果使用linux环境，可以在win下运行 `build_script.bat`自动构建

#### 开启调试模式：

```bash
./webhook-go -p 8080 debug=true
```

> 在调试模式下，非法请求的 body 会被写入 `request.log`，方便排查问题。

---

## 🔐 安全性说明

- 必须配置 `.secret` 文件或在第一次运行时输入 GitHub Webhook 的 Secret。
- 所有请求都会验证签名，防止伪造请求攻击。
- 不推荐在公网直接暴露该服务，应通过 Nginx 或反向代理加 HTTPS 并限制访问来源。

---

## 🧪 测试本地 Webhook

你可以使用 `curl` 来测试本地服务是否正常工作：

```bash
curl -X POST http://localhost:5599/hook \
     -H "X-Hub-Signature-256: sha256=your-signature" \
     -H "X-GitHub-Delivery: test-guid" \
     -d '{"test":"data"}'
```

⚠️ 注意：签名需真实有效才能触发脚本。

---

## 📌 常见问题

### Q：为什么没有执行 action.sh？

A：请确认：
- `action.sh` 是否存在且可执行
- Webhook 的 secret 是否一致
- 请求签名是否正确（调试模式可查看日志）
