## 🌀 Dismap - Asset discovery and identification tool
<a href="https://github.com/zhzyker/dismap"><img alt="Release" src="https://img.shields.io/badge/golang-1.6+-9cf"></a>
<a href="https://github.com/zhzyker/dismap"><img alt="Release" src="https://img.shields.io/badge/dismap-0.2-ff69b4"></a>
<a href="https://github.com/zhzyker/dismap"><img alt="Release" src="https://img.shields.io/badge/LICENSE-GPL-important"></a>
![GitHub Repo stars](https://img.shields.io/github/stars/zhzyker/dismap?color=success)
![GitHub forks](https://img.shields.io/github/forks/zhzyker/dismap)  
[[中文 Readme]](https://github.com/zhzyker/dismap/blob/main/readme-zh.md)  
Dismap positioning is an asset **discovery** and **identification** tool. It can quickly identify protocols and fingerprint information such as web/tcp/udp, locate asset types, and is suitable for internal and external networks. It assists red team personnel to quickly locate potential risk asset information, and assist blue team personnel to detect Suspected Fragile Assets

Dismap has a complete fingerprint rule base, currently including tcp/udp/tls protocol fingerprints and **4500+ web fingerprint rules**, which can identify favicon, body, header, etc. The introduction to the rule base is located at [RuleLab](https://github.com/zhzyker/dismap#-rulelab)

~~Scan results can be directly sent to [vulmap](https://github.com/zhzyker/vulmap)(>=0.8) for vulnerability scanning.~~ In version 0.3, the text result has been changed, the json file result has been added, and vulmap will support linkage in >= 1.0

## 🏂 Run
Dismap is a binary file for Linux, MacOS, and Windows. Go to [Release](https://github.com/zhzyker/dismap/releases) to download the corresponding version to run:
```Bash
# Linux or MacOS
zhzyker@debian:~$ chmod +x dismap-0.3-linux-amd64
zhzyker@debian:~$ ./dismap-0.3-linux-amd64 -h

# Windows
C:\Users\zhzyker\Desktop> dismap-0.3-windows-amd64.exe -h
```  
>  ![dismap](https://github.com/zhzyker/zhzyker/blob/main/dismap-images/dismap-0.3.png)


## 🎡 Optons
```Python
  -f, --file string     Parse the target from the specified file for batch recognition
  -h, --help            Show help
  -i, --ip string       Network segment [e.g. -i 192.168.1.0/24 or -i 192.168.1.1-10]
  -j, --json string     Scan result in json format [e.g. -j r.json]
  -l, --level int       Specify log level (0:Fatal 1:Error 2:Info 3:Warning 4:Debug 5:Verbose) (default 3)
  -m, --mode string     Specify the protocol [e.g. -m mysql/-m http]
      --nc              Do not print character colors
      --np              Not use ICMP/PING to detect surviving hosts
  -o, --output string   Save the scan results to the specified file (default "output.txt")
  -p, --port string     Custom scan ports [e.g. -p 80,443 or -p 1-65535]
      --proxy string    Use proxy scan, support http/socks5 protocol [e.g. --proxy socks5://127.0.0.1:1080]
  -t, --thread int      Number of concurrent threads (default 500)
      --timeout int     Response timeout time, the default is 5 seconds (default 5)
      --type string     Specify the type [e.g. --type tcp/--type udp]
  -u, --uri string      Specify a target URI [e.g. -u https://example.com]
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
* Dismap bug feedback or new feature suggestion [click me](https://github.com/zhzyker/dismap/issues)
* Twitter: https://twitter.com/zhzyker

## 🌈 RuleLab
The entire rule base is a struct located in [rule.go](https://github.com/zhzyker/dismap/blob/main/configs/rule.go)
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
