# MdDoc

通过项目的 README 生成简单的单页面项目文档

## 配置

```yaml
topnav:
    - Home

# show github star badge
starBadge: true
layoutFile: index.html
importCss:
    - path/to/some.css
importJs:
    - path/to/some.js
```

## 使用

```go
go get github.com/gookit/mddoc/mddoc
```

```bash
mddoc --file README.md --outdir output 
mddoc --file README.md --outdir output 
```

