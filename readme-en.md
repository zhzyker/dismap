## 🌀 Dismap - Asset discovery and identification tool
<a href="https://github.com/zhzyker/dismap"><img alt="Release" src="https://img.shields.io/badge/golang-1.6+-9cf"></a>
<a href="https://github.com/zhzyker/dismap"><img alt="Release" src="https://img.shields.io/badge/dismap-0.1-ff69b4"></a>
<a href="https://github.com/zhzyker/dismap"><img alt="Release" src="https://img.shields.io/badge/LICENSE-GPL-important"></a>
![GitHub Repo stars](https://img.shields.io/github/stars/zhzyker/dismap?color=success)
![GitHub forks](https://img.shields.io/github/forks/zhzyker/dismap)  
[[中文readme点我]](https://github.com/zhzyker/dismap/blob/main/readme.md)  
Dismap positioning is an asset discovery and identification tool; its characteristic function is to quickly identify Web fingerprint information and locate asset types. Assist the red team to quickly locate the target asset information, and assist the blue team to find suspected vulnerabilities

Dismap has a comprehensive fingerprint rule library, so you can easily customize new recognition rules. With the help of golang's concurrency advantages, rapid asset detection and identification can be achieved

The scan results can be directly submitted to [vulmap](https://github.com/zhzyker/vulmap) (>=0.8) for vulnerability scanning. Introduction to rule base in [RuleLab](https://github.com/zhzyker/dismap#-rulelab)  

## 🏂 Run
Dismap is a binary file for Linux, MacOS, and Windows. Go to [Release](https://github.com/zhzyker/dismap/releases) to download the corresponding version to run:  
```Bash
# Linux and MacOS
zhzyker@debian:~$ chmod +x dismap
zhzyker@debian:~$ ./dismap -h

# Windows
C:\Users\zhzyker\Desktop> dismap.exe -h
```  
>  ![dismap1](https://github.com/zhzyker/zhzyker/blob/main/dd.png)
>  ![dismap2](https://github.com/zhzyker/zhzyker/blob/main/dd2.png)

## 🎡 Optons
```Python
-file string
    Select a URL file for batch identification
-ip string
    Network segment [e.g. -ip 192.168.1.0/24 or -ip 192.168.1.1-10]
-np
    Not use ICMP/PING to detect surviving hosts
-output string
    Save the scan results to the specified file (default "output.txt")
-port string
    Custom scan ports [e.g. -port 80,443 or -port 1-65535]
-thread int
    Number of concurrent threads, (adapted to two network segments 2x254) (default 508)
-timeout int
    Response timeout time, the default is 5 seconds (default 5)
-url string
    Specify a target URL [e.g. -url https://example.com]
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

## ⛪ Discussion
* Dismap bug feedback or new feature suggestion [click me](https://github.com/zhzyker/dismap/issues)
* Twitter: https://twitter.com/zhzyker

## 🌈 RuleLab
The entire rule base is a struct located in [rule.go](https://github.com/zhzyker/dismap/blob/main/config/rule.go)
Rough format：
```Golang
Rule:
  Name: name /* Define rule name */
  Type: header|body|ico  /* Support recognized types, header, body, ico can be any logical combination, ico is to request favicon.ico separately and calculate MD5*/
  Mode: and|or /* Type judgment logic */
  Rule
    InBody: str  /* Specify which str exists in the response body */
    InHeader: str  /* Specify which str exists in the response Header */
    InIcoMd5: str_md5  /* MD5 of favicon.ico */
  Http:
    ReqMethod: GET|POST  /* Custom request method, currently supports GET and POST */
    ReqPath: str  /* Custom request web path */
    ReqHeader: []str  /* Customize the header of the Http request */
    ReqBody: str  /* Customize the body of the POST request */
```
**Example1:**  

Whether the character `<flink-root></flink-root>` exists in the response body
```Golang
{"Apahce Flink", "body", "", InStr{"(<flink-root></flink-root>)", "", ""}, ReqHttp{"", "", nil, ""}},
```  

**Example2:**  

Customize the request path `/myportal/control/main`, and determine whether there are header characters and body characters in the result of the custom request  
It can be found that all support regular expressions  
```Golang
{"Apache OFBiz", "body|header", "or", InStr{"(Apache OFBiz|apache.ofbiz)", "(Set-Cookie: OFBiz.Visitor=(.*))", ""}, ReqHttp{"GET", "/myportal/control/main", nil, ""}},
```

**The logical relationship of header, body, ico can be combined at will, but cannot be combined repeatedly:**  

Can: `"body|header|ico", "or"` or `"body|header|ico", "or|and"` or `"body|ico", "and"`   
Can't: `"body|body", "or"`  
Repeated combination is not allowed to be specified by type, but it can be achieved through InBody to determine the character: `"body", "", InStr{"(str1|str2)"}`  
