> # go-plugin基础示例

- 为注释做了汉化，方便理解学习
- 修改为从配置文件读取信息

项目仓库：https://github.com/hashicorp/go-plugin

改自：https://github.com/hashicorp/go-plugin/tree/main/examples/basic

## Windows:

1. 编译插件

```
go build -o ./plugin/greeter/main.exe ./plugin/greeter/main.go
```

2. 编译主程序

```
go build -o basic.exe
```

3. 启动主程序

```
./basic.exe
```



## Linux:

1. 编译插件

```
go build -o ./plugin/greeter/main ./plugin/greeter/main.go
```

2. 编译主程序

```
go build -o basic .
```

3. 启动主程序

```
./basic
```
