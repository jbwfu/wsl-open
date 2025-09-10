# wsl-open

[![Go Report Card](https://goreportcard.com/badge/github.com/jbwfu/wsl-open)](https://goreportcard.com/report/github.com/jbwfu/wsl-open)
[![Build Status](https://github.com/jbwfu/wsl-open/actions/workflows/go.yml/badge.svg)](https://github.com/jbwfu/wsl-open/actions/workflows/go.yml)
[![Latest release](https://img.shields.io/github/v/release/jbwfu/wsl-open)](https://github.com/jbwfu/wsl-open/releases/latest)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

[**English**](./README.md)

一个简单、快速且健壮的实用工具，用于在 WSL 中通过其默认的 Windows 应用程序打开文件、目录和 URL。

### 功能特性

-   **快速 & 轻量：** 作为编译后的原生二进制文件可即时启动，没有 shell 开销。
-   **安全 & 健壮：** 能正确处理包含空格和特殊字符的复杂路径。使用经过清理的命令执行方法来防止 shell 注入漏洞。
-   **简单 & 专注：** 只做一件事并把它做好。命令行界面极简且直观。
-   **开箱即用：** 无需安装额外的库或运行时。它通过利用其内置工具，在任何现代 WSL2 环境中都能工作。

---

## 安装

推荐的方式是从 [**Releases**](https://github.com/jbwfu/wsl-open/releases/latest) 页面下载最新的预编译二进制文件。

```sh
# 将 vX.Y.Z 替换为最新的版本号，例如：v1.0.0
wget https://github.com/jbwfu/wsl-open/releases/download/vX.Y.Z/wsl-open_linux_amd64.tar.gz
tar -xvf wsl-open_linux_amd64.tar.gz
sudo mv wsl-open /usr/local/bin/
```

<details>
<summary>其它安装方式</summary>

**使用 `go install`：**
```sh
go install github.com/jbwfu/wsl-open@latest
```

**从源码构建：**
```sh
git clone https://github.com/jbwfu/wsl-open.git
cd wsl-open
go build
sudo mv wsl-open /usr/local/bin/
```
</details>

---

## 使用方法

提供一个文件路径、目录路径或 URL 作为参数。

<details>
<summary><strong>关于 WSL1 的支持说明</strong></summary>

本工具已为 **WSL2** 进行了全面优化。由于文件系统架构的根本差异，在 **WSL1** 上打开文件和目录的支持是有限的。

-   打开 URL：在 WSL1 和 WSL2 上均可正常工作。
-   打开 Windows 磁盘分区的文件/目录（例如 `/mnt/c/...`）：在 WSL1 和 WSL2 上均可正常工作。
-   打开原生 Linux 文件系统的文件/目录（例如 `~/file.txt`）：仅在 WSL2 上得到可靠支持。

没有为 WSL1 添加特定兼容性方案的计划，强烈建议使用 WSL2 以获得最佳体验。
</details>

```sh
# 在其默认的 Windows 应用程序中打开一个文件
wsl-open "My Documents/report.docx"

# 在 Windows 文件资源管理器中打开当前目录
wsl-open .

# 在默认的 Windows 浏览器中打开一个 URL
wsl-open https://github.com
```

使用 `-x` 标志可以查看将要执行的命令，而不会实际运行它。
要获取更多细节，请运行 `wsl-open --help`。

#### 高级技巧：作为 `xdg-open` 的备用方案

对于可能未包含 `xdg-open` 的最小化 WSL 环境，您可以将 `wsl-open` 链接为系统级的默认命令。这使得那些调用 `xdg-open` 的工具（如 `git browse`）能够无缝工作。

首先，检查 `xdg-open` 是否已安装：
```sh
command -v xdg-open
```

如果以上命令没有任何输出，您就可以创建一个符号链接：
```sh
# 这会找到 wsl-open 的位置，并将其链接到 /usr/local/bin/xdg-open
sudo ln -s "$(command -v wsl-open)" /usr/local/bin/xdg-open
```
---

## TODO

-   [ ] 完整的 `xdg-open` 集成。
-   [ ] 原生路径转换（移除对 `wslpath` 的依赖）。
---

## 贡献

欢迎任何形式的贡献！如果您发现 Bug 或有功能请求，请随时创建 Issue 或提交 Pull Request。

## 许可证

本项目采用 **MIT 许可证** 授权。
