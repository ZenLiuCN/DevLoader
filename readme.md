# loader - Program Environment loader
## Ideal
It's comes from that I need to work on different workstations for serveral program enveriments.  
That's really a troubling thing to install and config every thing,then I realized most program to be run only needs some environment variables to be configuated.  
So I create this tool to keep my programming system portable.
I had used bat/C++/nodejs/autoit and some other language or tool to create loader, finally I choose use golang to rewrite it.
## What can it do
With configuations and some other softwares,it can create a development environment contains Java/Go/NodeJs/Shell/git/RUST/Docker...
## Build
*only for windows*
```shell
1. install golang >=1.8+
2. go get gopkg.in/yaml.v2
3. go get github.com/axgle/mahonia
4. build.bat
```
## usage
rename and move executeable to some folder (name must be alphabet);

excute it will generate a template configuration file;

edit configuration file (yaml format)

## configuration
```yaml
# File Encoding utf-8
# debug             output debug log; DEFAULT:false; when set true will create NAME-TIMESTAMP.log
# runbefore         command to execute before set Environment Variables
# env               Environment Variables need to be set
# runafter          command to execute after set Environment Variables
# runquit           command to execute when LOADER quit
# *.x86             settings will be used when is X86 os or X64 not configuated
# *.x64             settings will be used when is X64 os
# run*.x*.path      command path or name (with relative path of LOADER should be start with '/'
# run*.x*.param     command parameters (list of string)
# run*.x*.show      show execute window (DEFALUT: false)
# run*.x*.must      panic when failed (DEFAULT:false)
# run*.x*.fixparam  fix parameter path (DEFAULT:false) will fix parameter begin with './' as relative path of LOADER
# run*.x*.wait      will wait command finished (DEFAULT:false)
# run*.x*.sleep     sleep time after excute by seconds;DEFAULT:0 ;if not zero ,path can be empty
# env.*.name        Environment Variable name
# env.*.type        setting mode:
#                               0 append as path (DEFAULT)
#                               1 preappend as path
#                               2 replace as path
#                               3 append as value
#                               4 preappend as value
#                               5 replace as value
#                               6 create as value (only when variable is not set)
#                               7 create as path (only when variable is not set)
debug: false
runbefore:
    x86:
        -   path: notepad
            param:
                - execute1
            show: true
            must: true
            fixparam: true
            wait: false
            sleep: 2
        -   path: notepad
            param:
                - execute2
            show: true
            must: true
            fixparam: true
            wait: false
            sleep: 2
    x64:
        -   path: notepad
            param:
                - execute64-1
            show: true
            must: true
            fixparam: true
            wait: false
            sleep: 2
        -   path: notepad
            param:
                - execute64-2
            show: true
            must: true
            fixparam: true
            wait: false
            sleep: 2
env:
    -   name: path
        type: 0
        x86:
            - /dir
        x64:
            - /dir
runafter:
    x86:
        -   path: notepad
            param:
                - abcd
            show: true
            must: true
            fixparam: true
            wait: false
            sleep: 2
```

# loader - 编程环境加载器
## 点子来源
由于工作和兴趣需要，我经常要在不同的工作电脑上使用多种开发平台；而每次安装这些平台是个很烦人的事情。最后我意识到大多数开发环境需要的只是几个环境变量，于是这个工具产生了。
这个工具曾经由 bat脚本、c++、nodejs、autoit 等等语言或工具实现过，最后我选择用go语言来再一次整合实现。
## 能做什么
通过配置文件和工具软件配合，该启动器可以创建一套包含 Java、Go、Nodejs、Mingw、git、shell等的便携开发环境。
## 编译
*仅windows可用*
```shell
1. install golang >=1.8+
2. go get gopkg.in/yaml.v2
3. go get github.com/axgle/mahonia
4. build.bat
```
## 使用
将编译获得的可执行文件命名为需要的名字(英文)  
运行可执行文件将自动生成配置模版  
编辑配置文件(yaml格式),保存,重新启动  