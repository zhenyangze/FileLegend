过期文件删除工具
-----
主要解决文件缓存失效不清理导致的磁盘占用量上升

**使用场景**

- laravel 文件缓存
- 页面静态化
- 生成的临时下载文件

**配置文件**

config.toml

```toml
pidpath ="/tmp/filelegend.pid"
#如果匹配不到是否要删除掉
isforce = "0"
showlog = "1"
[items]
    [[items.node]]
        rootdir = "/storage/framework/cache/data/"
        reg     = "^(\\d{10})"
    [[items.node]]
        rootdir = "/public/page-cache"
        reg     = "<!--\\((\\d{10})\\)-->"
```

