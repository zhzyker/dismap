## 🌀 Dismap - Asset discovery and identification tool
Dismap 定位是一个资产发现和识别工具；其特色功能在于快速识别 Web 指纹信息，定位资产类型。辅助红队快速定位目标资产信息，辅助蓝队发现疑似脆弱点  
Dismap 拥有完善的指纹规则库，可轻松自定义新识别规则。借助于 golang 并发优势，即可实现快速资产探测与识别

## 🏂 Run
Dismap 对 Linux、MacOS、Windows 均提供了二进制可执行文件，前往 [Release](https://github.com/zhzyker/dismap/releases) 下载对应版本即可运行:
```Bash
# Linux and MacOS
zhzyker@debian:~$ chmod +x dismap
zhzyker@debian:~$ ./dismap -h

# Windows
C:\Users\zhzyker\Desktop> dismap.exe -h
```

## 🎡 Optons
```Python
-file string
    Select a URL file for batch identification
    # 从文件中读取 Url 进行批量识别
-ip string
    Network segment [e.g. -ip 192.168.1.0/24 or -ip 192.168.1.1-10]
    # 指定一个网段,格式示例: 192.168.1.1/24  192.168.1.1-100  192.168.1.1-192.168.1.254
-np
    Not use ICMP/PING to detect surviving hosts
    # 不进行主机存活检测,跳过存活检测直接识别 Url
-output string
    Save the scan results to the specified file (default "output.txt")
    # 自定义识别结果输出文件,默认追加到 output.txt 中
-port string
    Custom scan ports [e.g. -port 80,443 or -port 1-65535]
    # 自定义需要扫描的 Web 端口,默认端口在 /config/config.go 中
-thread int
    Number of concurrent threads, (adapted to two network segments 2x254) (default 508)
    # 多线程数量,默认508(两个C段的数量),线程越高存活和识别丢失率可能越高,不建议超过2000
-timeout int
    Response timeout time, the default is 5 seconds (default 5)
    # 主机存活探测和 Http 超时时间,默认均为5秒
-url string
    Specify a target URL [e.g. -url https://example.com]
    # 识别单个 Url 时用该选项指定
```

## 🎨 Examples
```Bash
zhzyker@debian:~$ ./dismap -ip 192.168.1.1/24
zhzyker@debian:~$ ./dismap -ip 192.168.1.1/24 -output result.txt
zhzyker@debian:~$ ./dismap -ip 192.168.1.1/24 -np -timeout 10
zhzyker@debian:~$ ./dismap -ip 192.168.1.1/24 -thread 1000
zhzyker@debian:~$ ./dismap -url https://github.com/zhzyker/dismap
zhzyker@debian:~$ ./dismap -ip 192.168.1.1/24 -port 1-65535
```

## 🌈 RuleLab
整个规则库是一个 struct 位于 [rule.go](https://github.com/zhzyker/dismap/blob/main/config/rule.go)
大致格式如下：
```Golang
Rule:
  Name: name /* 定义规则名称 */
  Type: header, body, ico  /* 支持识别的类型, header 和 body 均为响应 Body 中, ico 为单独请求 favicon.ico 并计算 MD5*/
  Mode: and, or /* 类型的判断逻辑关系 */
  Rule
    InBody: str  /* 需要指定响应 Body 中存在 str 则命中 */
    InHeader: str  /* 需要指定响应 Hedaer 中存在 str 则命中 */
    InIcoMd5: str_md5  /* favicon.ico 的 MD5 值 */
  Http:
    ReqMethod GET, POST  /* 自定义请求方法,目前支持 GET 和 POST */
    ReqPath str  /* 自定义请求 Web 路径 */
    ReqHeader []str  /* 自定义 Http 请求的 Header */
    ReqBody str  /* 自定义 POST 请求时的 Body */
```
** 规则库示例1: **

即在响应Body中检查是否存在字符`<flink-root></flink-root>`
```Golang
{"Apahce Flink", "body", "", InStr{"(<flink-root></flink-root>)", "", ""}, ReqHttp{"", "", nil, ""}},
```  

** 规则库示例2: ** 

自定义请求访问`/myportal/control/main`,判断自定义请求的结果中是否存在指定的 header 字符和 body 字符  
可以发现均支持正则表达式  
```Golang
{"Apache OFBiz", "body|header", "or", InStr{"(Apache OFBiz|apache.ofbiz)", "(Set-Cookie: OFBiz.Visitor=(.*))", ""}, ReqHttp{"GET", "/myportal/control/main", nil, ""}},
```

** header, body, ico 的逻辑关系可以随意组合,但不可重复组合：**   

允许: `"body|header|ico", "or"` or `"body|header|ico", "or|and"` or `"body|ico", "and"`   
不允许: `"body|body", "or"`  
重复组合不允许通过类型指定,但可通过 InBody 判断字符内实现: `"body", "", InStr{"(str1|str2)"}`  
