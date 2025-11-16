# GitPack Core

![License](https://img.shields.io/badge/License-MIT-dark_green)

这是[GitPack](https://github.com/Zhoucheng133/GitPack)的一部分，你也可以单独使用

## 如果你想要自行打包成动态库   

使用下面的命令来生成动态库

```bash
# 对于Windows系统
go build -o build/core.dll -buildmode=c-shared .
# 对于macOS系统
go build -o build/core.dylib -buildmode=c-shared .
```