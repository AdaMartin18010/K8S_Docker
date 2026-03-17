# eBPF Go 示例

## 依赖

```bash
go get github.com/cilium/ebpf
```

## 安装工具

```bash
# 安装 bpf2go
go install github.com/cilium/ebpf/cmd/bpf2go@latest

# 安装 clang
apt-get install clang llvm
```

## 编译 eBPF 程序

```bash
go generate
```

## 运行

```bash
sudo go run main.go
```

## 参考

- [cilium/ebpf](https://github.com/cilium/ebpf)
- [eBPF.io](https://ebpf.io/)
