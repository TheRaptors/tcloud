# 腾讯云 DNSPod & CVM 管理工具

基于 Go + 腾讯云 SDK 的命令行工具，实现竞价实例一键部署与 DNS 自动切换。

## 项目结构

```
my-project-golang/
├── cmd/tcloud/            # 入口
│   └── main.go
├── internal/
│   ├── config/            # 配置加载
│   │   └── config.go
│   ├── dnspod/            # DNS 域名操作
│   │   ├── describe_record_list.go
│   │   ├── describe_record.go
│   │   └── modify_record.go
│   └── cvm/               # 云服务器操作
│       ├── run_instances.go
│       ├── describe_instances.go
│       └── terminate_instances.go
├── config/
│   └── tencentcloud.json  # 配置文件（已加入 .gitignore）
├── go.mod
└── go.sum
```

## 配置文件

编辑 `config/tencentcloud.json`：

```json
{
    "secret_id": "你的 SecretId",
    "secret_key": "你的 SecretKey",
    "region": "ap-hongkong",
    "domain": "example.com",
    "subdomain": "cvm",
    "private_ip": "172.19.0.100",
    "zone": "ap-hongkong-2",
    "vpc_id": "vpc-xxxxxx",
    "subnet_id": "subnet-xxxxxx",
    "security_group_ids": ["sg-xxxxxx"],
    "instance_name": "未命名",
    "instance_type": "MA5.LARGE32",
    "image_id": "img-xxxxxx",
    "key_id": "skey-xxxxxx",
    "max_price": "1000"
}
```

| 字段 | 说明 |
|------|------|
| secret_id / secret_key | 腾讯云 API 密钥 |
| region | 地域，如 ap-hongkong |
| domain / subdomain | 域名与主机名，如 cvm.example.com |
| private_ip | 指定的内网 IP |
| zone | 可用区 |
| vpc_id / subnet_id | VPC 与子网 |
| security_group_ids | 安全组 ID 列表 |
| instance_name | 实例备注名 |
| instance_type | 实例规格 |
| image_id | 镜像 ID |
| key_id | SSH 密钥 ID |
| max_price | 竞价最高出价 |

## 编译

### Linux / macOS

```bash
GOOS=linux   GOARCH=amd64 go build -o tcloud ./cmd/tcloud
GOOS=darwin  GOARCH=amd64 go build -o tcloud ./cmd/tcloud
GOOS=darwin  GOARCH=arm64 go build -o tcloud ./cmd/tcloud
```

### Windows CMD

```cmd
set GOOS=windows
set GOARCH=amd64
go build -o tcloud.exe ./cmd/tcloud
```

交叉编译 Linux：

```cmd
set GOOS=linux
set GOARCH=amd64
go build -o tcloud ./cmd/tcloud
```

### Windows PowerShell

```powershell
$env:GOOS="linux"; $env:GOARCH="amd64"; go build -o tcloud ./cmd/tcloud
```

## 使用方法

### 开发运行

```bash
go run ./cmd/tcloud <命令>
```

### 编译后运行

```bash
./tcloud <命令>      # Linux/macOS
tcloud.exe <命令>    # Windows
```

> 配置文件 `config/tencentcloud.json` 需与可执行文件保持相对目录结构。

## 命令列表

| 命令 | 说明 |
|------|------|
| `list` | 获取 DNS 解析记录列表（提取 RecordId） |
| `list --detail` | 获取记录列表后查询第一条记录详情 |
| `describe` | 自动获取 RecordId 并查询记录详情 |
| `modify <IP>` | 自动获取 RecordId 并修改 A 记录指向的 IP |
| `run-instances` | 创建竞价 CVM 实例 |
| `deploy` | 一键部署：购买实例 → 获取公网IP → 修改 DNS |
| `destroy` | 按内网 IP 自动查找并销毁实例 |
| `undeploy` | 一键回收：查找实例 → 销毁 → DNS 还原为 0.0.0.0 |

### 典型流程

**部署：**

```bash
go run ./cmd/tcloud deploy
```

执行步骤：购买竞价实例 → 获取公网IP → 获取 RecordId → 查看修改前 DNS → 修改 A 记录 → 查看修改后 DNS

**回收：**

```bash
go run ./cmd/tcloud undeploy
```

执行步骤：按内网IP查找实例 → 获取 RecordId → 查看销毁前 DNS → 销毁实例 → DNS 还原为 0.0.0.0 → 查看修改后 DNS

## 自动构建（GitHub Actions）

推送 tag 时自动构建 6 个平台的二进制文件并创建 GitHub Release：

```bash
git tag v1.0.0
git push origin v1.0.0
```

构建产物：

| 文件 | 平台 |
|------|------|
| `tcloud_windows_amd64.exe` | Windows 64位 |
| `tcloud_windows_arm64.exe` | Windows ARM64 |
| `tcloud_linux_amd64` | Linux 64位 |
| `tcloud_linux_arm64` | Linux ARM64 |
| `tcloud_darwin_amd64` | macOS Intel |
| `tcloud_darwin_arm64` | macOS Apple Silicon |

## 注意事项

- 配置文件包含敏感信息，已加入 `.gitignore`，请勿提交至版本库
- 首次使用请复制 `config/tencentcloud.json.example` 为 `config/tencentcloud.json` 并填入真实配置
- 子账号使用 `deploy` 创建竞价实例时，需在 CAM 中开通**财务支付**权限
- `destroy` / `undeploy` 通过匹配配置中的 `private_ip` 自动查找实例，无需手动输入实例 ID
- 二进制文件中不包含密钥等敏感信息（运行时从配置文件读取）