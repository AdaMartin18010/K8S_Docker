# OCI 开放容器标准

> 容器标准化规范详解

---

## 什么是 OCI？

OCI (Open Container Initiative) 是 Linux 基金会主导的开源项目，旨在制定容器技术的开放标准。

---

## OCI 三大规范

| 规范 | 说明 | 相关工具 |
|------|------|----------|
| **Runtime Spec** | 容器运行时标准 | runc, crun, kata-runtime |
| **Image Spec** | 镜像格式标准 | Docker, containerd |
| **Distribution Spec** | 镜像分发的标准 | Docker Registry, Harbor |

---

## OCI Runtime Spec

定义了容器运行时的标准接口：

```json
{
  "ociVersion": "1.0.2",
  "process": {
    "terminal": false,
    "user": {"uid": 0, "gid": 0},
    "args": ["sh", "-c", "echo hello"],
    "env": ["PATH=/usr/local/bin"]
  },
  "root": {
    "path": "rootfs",
    "readonly": false
  },
  "hostname": "mycontainer",
  "linux": {
    "namespaces": [
      {"type": "pid"},
      {"type": "network"},
      {"type": "ipc"},
      {"type": "uts"},
      {"type": "mount"}
    ],
    "cgroupsPath": "/docker/mycontainer",
    "resources": {
      "cpu": {"shares": 1024},
      "memory": {"limit": 536870912}
    }
  }
}
```

---

## OCI Image Spec

定义了镜像的格式和布局：

```
my-image/
├── blobs/
│   └── sha256/
│       ├── abc123... (config)
│       ├── def456... (layer1)
│       └── ghi789... (layer2)
├── index.json
├── manifest.json
└── oci-layout
```

### 关键概念

| 概念 | 说明 |
|------|------|
| **Layer** | 镜像层，只读的文件系统变更 |
| **Manifest** | 镜像清单，描述镜像内容和元数据 |
| **Config** | 镜像配置，包含环境变量、入口点等 |
| **Index** | 多平台镜像的索引 |

---

## OCI Distribution Spec

定义了镜像仓库的 API 标准：

```
# 拉取镜像
GET /v2/<name>/manifests/<reference>

# 上传镜像
PUT /v2/<name>/manifests/<reference>

# 下载层
GET /v2/<name>/blobs/<digest>

# 上传层
POST /v2/<name>/blobs/uploads/
```

---

## 为什么要标准化？

1. **互操作性**: 不同工具可以处理相同的镜像
2. **可移植性**: 镜像可以在不同平台运行
3. **生态繁荣**: 降低创新门槛

---

## 符合 OCI 标准的工具

| 类型 | 工具 |
|------|------|
| 运行时 | runc, crun, youki |
| 镜像构建 | Docker, Buildah, kaniko |
| 镜像仓库 | Docker Registry, Harbor, distribution |
