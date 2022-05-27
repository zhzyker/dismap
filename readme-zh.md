## 🌀 Dismap - Asset discovery and identification tool
<a href="https://github.com/zhzyker/dismap"><img alt="Release" src="https://img.shields.io/badge/golang-1.6+-9cf"></a>
<a href="https://github.com/zhzyker/dismap"><img alt="Release" src="https://img.shields.io/badge/dismap-0.3-ff69b4"></a>
<a href="https://github.com/zhzyker/dismap"><img alt="Release" src="https://img.shields.io/badge/LICENSE-GPL-important"></a>
![GitHub Repo stars](https://img.shields.io/github/stars/zhzyker/dismap?color=success)
![GitHub forks](https://img.shields.io/github/forks/zhzyker/dismap)
![GitHub all release](https://img.shields.io/github/downloads/zhzyker/dismap/total?color=blueviolet)  
  
Dismap 定位是一个资产发现和识别工具，他可以快速识别 Web/tcp/udp 等协议和指纹信息，定位资产类型，适用于内外网，辅助红队人员快速定位潜在风险资产信息，辅助蓝队人员探测疑似脆弱资产

Dismap 拥有完善的指纹规则库，目前包括 tcp/udp/tls 协议指纹和 4500+ Web 指纹规则，可以识别包括 favicon、body、header 等，   对于规则库的简介位于 [RuleLab](https://github.com/zhzyker/dismap#-rulelab)

~~扫描结果可直接丢给 [Vulmap](https://github.com/zhzyker/vulmap)(>=0.8) 进行漏洞扫描。~~, 0.3 版本中改变了文本结果，新增了 json 文件结果，vulmap 将在 >= 1.0 支持联动

## 🏂 Run
Dismap 对 Linux、MacOS、Windows 均提供了二进制可执行文件，前往 [Release](https://github.com/zhzyker/dismap/releases) 下载对应版本即可运行:
```Bash
# Linux or MacOS
zhzyker@debian:~$ chmod +x dismap-0.3-linux-amd64
zhzyker@debian:~$ ./dismap-0.3-linux-amd64 -h

# Windows
C:\Users\zhzyker\Desktop> dismap-0.3-windows-amd64.exe -h
```  
>  ![dismap](https://github.com/zhzyker/zhzyker/blob/main/dismap-images/dismap-0.3.png)



## 🎡 Options
```Bash
  -f, --file string     从文件中解析目标进行批量识别
  -h, --help            查看帮助说明
  -i, --ip string       指定一个网段 [示例 -i 192.168.1.0/24 or -i 192.168.1.1-10]
  -j, --json string     扫描结果保存到 json 格式文件
  -l, --level int       指定日志等级 (0:Fatal 1:Error 2:Info 3:Warning 4:Debug 5:Verbose) (默认 3)
  -m, --mode string     指定要识别的协议 [e.g. -m mysql/-m http]
      --nc              不打印字符颜色
      --np              不使用 ICMP/PING 检测存活主机
  -o, --output string   将扫描结果保存到指定文件 (默认 "output.txt")
  -p, --port string     自定义要识别的端口 [示例 -p 80,443 or -p	 1-65535]
      --proxy string    使用代理进行扫描, 支持 http/socks5 协议代理 [示例 --proxy socks5://127.0.0.1:1080]
  -t, --thread int      并发线程数量 (默认 500)
      --timeout int     超时时间 (默认 5)
      --type string     指定扫描类型 [示例 --type tcp/--type udp]
  -u, --uri string      指定目标地址 [示例 -u https://example.com]

```

## 🎨 Examples
```Bash
zhzyker@debian:~$ ./dismap -i 192.168.1.1/24
zhzyker@debian:~$ ./dismap -i 192.168.1.1/24 -o result.txt -j result.json
zhzyker@debian:~$ ./dismap -i 192.168.1.1/24 --np --timeout 10
zhzyker@debian:~$ ./dismap -i 192.168.1.1/24 -t 1000
zhzyker@debian:~$ ./dismap -u https://github.com/zhzyker/dismap
zhzyker@debian:~$ ./dismap -u mysql://192.168.1.1:3306
zhzyker@debian:~$ ./dismap -i 192.168.1.1/24 -p 1-65535
```

## ⛪ Discussion
* Dismap Bug 反馈或新功能建议[点我](https://github.com/zhzyker/dismap/issues)
* Twitter: https://twitter.com/zhzyker
* WeChat: 扫码滴滴我入群聊
<p>
    <img alt="QR-code" src="https://github.com/zhzyker/zhzyker/blob/main/my-wechat.jpg" width="20%" height="20%" style="max-width:100%;">
</p>

## 🌈 RuleLab
整个规则库是一个 struct 位于 [rule.go](https://github.com/zhzyker/dismap/blob/main/config/rule.go)
大致格式如下：
```Golang
Rule:
Name: name /* 定义规则名称 */
Type: header|body|ico  /* 支持识别的类型, header、body、ico 可任意逻辑组合, ico 为单独请求 favicon.ico 并计算 MD5*/
Mode: and|or /* 类型的判断逻辑关系 */
Rule
InBody: str  /* 需要指定响应 Body 中存在 str 则命中 */
InHeader: str  /* 需要指定响应 Hedaer 中存在 str 则命中 */
InIcoMd5: str_md5  /* favicon.ico 的 MD5 值 */
Http:
ReqMethod: GET|POST  /* 自定义请求方法,目前支持 GET 和 POST */
ReqPath: str  /* 自定义请求 Web 路径 */
ReqHeader: []str  /* 自定义 Http 请求的 Header */
ReqBody: str  /* 自定义 POST 请求时的 Body */
```
**规则库示例1:**

即在响应Body中检查是否存在字符`<flink-root></flink-root>`
```Golang
{"Apahce Flink", "body", "", InStr{"(<flink-root></flink-root>)", "", ""}, ReqHttp{"", "", nil, ""}},
```  

**规则库示例2:**

自定义请求访问`/myportal/control/main`,判断自定义请求的结果中是否存在指定的 header 字符和 body 字符  
可以发现均支持正则表达式
```Golang
{"Apache OFBiz", "body|header", "or", InStr{"(Apache OFBiz|apache.ofbiz)", "(Set-Cookie: OFBiz.Visitor=(.*))", ""}, ReqHttp{"GET", "/myportal/control/main", nil, ""}},
```

**header, body, ico 的逻辑关系可以随意组合,但不可重复组合:**

允许: `"body|header|ico", "or"` or `"body|header|ico", "or|and"` or `"body|ico", "and"`   
不允许: `"body|body", "or"`  
重复组合不允许通过类型指定,但可通过 InBody 判断字符内实现: `"body", "", InStr{"(str1|str2)"}`  
