# GO-AUTO
一个开箱即用，基于 robot-go 实现的录制键鼠操作，并回放键鼠操作工具。可以用于游戏、web检查等黑盒自动化测试场景。

# Setup
在release中有不同平台的可用安装包, 根据操作系统的不同下载并安装不同的包体即可

如果需要独自编译源码,请根据以下步骤：

- required
- git clone 
- go install wails
- go mod tidy
- wails dev

## Required

1. 安装gcc环境
   > 安装 llvm-mingw https://github.com/mstorsjo/llvm-mingw/releases/tag/20240619
   
   或者
    
   > 安装 mingw-w64 https://www.mingw-w64.org/
   
   1.1 配置环境变量
   > 获取llvm-mingw / mingw-w64 安装路径,配置到环境变量Path中。
   > 
   > example: C://User/llvm-mingw/bin
   
   1.2 验证是否安装成功
   > 打开cmd / powershell, 输入gcc -v 当看到类似以下的输出时则证明安装成功.
   
   ![gcc-cmd](docs/images/gcc-cmd.png)

## NSIS
