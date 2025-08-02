# CoraPanel MCP Server

**CoraPanel MCP Server** là triển khai giao thức Model Context Protocol (MCP) cho [CoraPanel](https://github.com/cloudora-vn/CoraPanel).

---

## Cách cài đặt


### ✅ Cách 1: Tải từ trang Release (Khuyến nghị)

1. Truy cập [trang Releases](https://github.com/cloudora-vn/mcp-corapanel/releases), tải file thực thi phù hợp với hệ thống của bạn.

2. Ví dụ cài đặt (với `amd64`):

```bash
chmod +x mcp-corapanel-linux-amd64
mv mcp-corapanel-linux-amd64 /usr/local/bin/mcp-corapanel
```

---

### 🛠️ Cách 2: Build từ mã nguồn

Đảm bảo đã cài đặt Go 1.23 hoặc cao hơn, sau đó thực hiện:

1. Clone repository:

```bash
git clone https://github.com/cloudora-vn/mcp-corapanel.git
cd mcp-corapanel
```

2. Build file thực thi:

```bash
make build
```

3. File thực thi được tạo tại: `./build/mcp-corapanel`, nên di chuyển vào thư mục PATH của hệ thống.

---

### 🚀 Cách 3: Cài đặt qua `go install`

Đảm bảo đã cài đặt Go 1.23 hoặc cao hơn:

```bash
go install github.com/cloudora-vn/mcp-corapanel@latest
```

---

### 🐳 Cách 4: Cài đặt qua Docker

Đảm bảo Docker đã được cài đặt và cấu hình đúng.

Image chính thức hỗ trợ các kiến trúc sau:

- `amd64`
- `arm64`
- `arm/v7`
- `s390x`
- `ppc64le`

---

## Cách sử dụng

CoraPanel MCP Server hỗ trợ hai chế độ: **stdio** và **sse**

---

### Chế độ 1: stdio (mặc định)

#### 📦 Sử dụng file binary

Trong file cấu hình của Cursor hoặc Windsurf, thêm nội dung sau:

```json
{
  "mcpServers": {
    "mcp-corapanel": {
      "command": "mcp-corapanel",
      "env": {
        "PANEL_ACCESS_TOKEN": "<token truy cập CoraPanel của bạn>",
        "PANEL_HOST": "ví dụ http://localhost:8080"
      }
    }
  }
}
```

#### 🐳 Sử dụng Docker

```json
{
  "mcpServers": {
    "mcp-corapanel": {
      "command": "docker",
      "args": [
        "run",
        "-i",
        "--rm",
        "-e",
        "PANEL_HOST",
        "-e",
        "PANEL_ACCESS_TOKEN",
        "cloudora-vn/corapanel-mcp-server"
      ],
      "env": {
        "PANEL_HOST": "ví dụ http://localhost:8080",
        "PANEL_ACCESS_TOKEN": "<token truy cập CoraPanel của bạn>"
      }
    }
  }
}
```

---

### Chế độ 2: sse

#### 🚀 Khởi động MCP Server:

```bash
mcp-corapanel -host http://localhost:8080 -token <token truy cập CoraPanel của bạn> -transport sse -addr http://localhost:8000
```

#### ⚙️ Cấu hình Cursor hoặc Windsurf:

```json
{
  "mcpServers": {
    "mcp-corapanel": {
      "url": "http://localhost:8000/sse"
    }
  }
}
```

---

### 🔧 Tham số dòng lệnh

- `-token`: Token truy cập CoraPanel
- `-host`: Địa chỉ CoraPanel, ví dụ: http://localhost:8080
- `-transport`: Phương thức truyền tải: `stdio` hoặc `sse`, mặc định là `stdio`
- `-addr`: Địa chỉ lắng nghe SSE server, mặc định là `http://localhost:8000`

---

## 🧰 Các công cụ khả dụng (Tools)

Danh sách các công cụ MCP Server cung cấp để tương tác với CoraPanel:

| Tên công cụ             | Phân loại   | Mô tả                            |
|-------------------------|-------------|----------------------------------|
| `get_dashboard_info`    | Hệ thống    | Lấy thông tin bảng điều khiển   |
| `get_system_info`       | Hệ thống    | Lấy thông tin hệ thống          |
| `list_websites`         | Website     | Liệt kê tất cả website          |
| `create_website`        | Website     | Tạo website mới                 |
| `list_ssls`             | Chứng chỉ   | Liệt kê tất cả chứng chỉ        |
| `create_ssl`            | Chứng chỉ   | Tạo chứng chỉ mới               |
| `list_installed_apps`   | Ứng dụng    | Liệt kê ứng dụng đã cài đặt     |
| `install_openresty`     | Ứng dụng    | Cài đặt OpenResty               |
| `install_mysql`         | Ứng dụng    | Cài đặt MySQL                   |
| `list_databases`        | Cơ sở dữ liệu| Liệt kê tất cả database        |
| `create_database`       | Cơ sở dữ liệu| Tạo database mới               |