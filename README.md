# gojvm

在 go 环境下加载 jvm 的字节码并允许调用。
 
- - -

使用方法:

一、先确定 classpath, 存放要调用的 java 类，可以是 class 文件，可以是 jar 等。

二、在 go 代码中启动虚拟机，然后可以直接在 Go 内对 java 类进行调用。

```go
// 启动虚拟机
jvm := gojvm.NewJVM(".", "128m", "512m", "384m", "512k")
// 拿到依附于当前线程的 jvm 环境上下文
env := jvm.Attach()

// 新建一个 java 类
c0 := env.NewObject("com.example.MyClass")
// 从 java 类中获取成员变量
s0 := c0.GetString("value")
fmt.Printf("s0 = %v\n", s0)

// 向 java 类中写入新值
c0.SetString("value", "hello world")
s1 := c0.GetString("value")
fmt.Printf("s1 = %v\n", s1)

// 将 java 类释放掉
c0.Free()

// 从当前线程卸载 jvm 环境
env.Detach()
// 释放虚拟机
jvm.Free()
```