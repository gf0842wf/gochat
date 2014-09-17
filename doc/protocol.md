# 聊天服务器协议 #

`所有时间为unix时间戳且为整数,单位ms,其中uid,sid是数字(在redis中是字符串),gid是字符串最大32bytes`  
`now = int(time.time()*1000)`  
`聊天消息内不能含连续的\xef\xff`

## 前端部署 ##
haproxy进行tcp负载均衡反向代理

	<haproxy.cfg>
	# this config needs haproxy-1.1.28 or haproxy-1.2.1
	
	global
		log 127.0.0.1	local0
		log 127.0.0.1	local1 notice
		#log loghost	local0 info
		maxconn 65535
		ulimit-n 131086
		#chroot /usr/share/haproxy
		user haproxy
		group haproxy
		daemon
		nbproc  5 # 五个并发进程
		pidfile /var/run/haproxy.pid
		#debug
		#quiet
	
	defaults
		#log	global
		mode	http
		option	httplog
		option	dontlognull
		retries	2
		option redispatch
		maxconn	4096
		contimeout	5000
		clitimeout	50000
		srvtimeout	50000
	
	########统计页面配置########
	listen admin_stats
		bind 0.0.0.0:9100               #监听端口 
		mode http                       #http的7层模式  
		option httplog                  #采用http日志格式 
		log 127.0.0.1 local0 err
		maxconn 10
		stats enable
		stats refresh 30s               #统计页面自动刷新时间  
		stats uri /                     #统计页面url  
		stats realm XingCloud\ Haproxy  #统计页面密码框上提示文本  
		stats auth admin:admin          #统计页面用户名和密码设置  
		stats hide-version              #隐藏统计页面上HAProxy的版本信息  
	
	########chat服务器配置############# 
	listen chat
		bind 0.0.0.0:9000
		mode tcp
		maxconn 100000
		log 127.0.0.1 local0 debug
		server s1 192.168.1.111:9001 weight 1 # 这个可以部署haproxy
		server s1 192.168.1.112:9001 weight 5
		server s1 192.168.1.113:9001 weight 5
		server s1 192.168.1.114:9001 weight 5
	########frontend配置############### 
`sudo haproxy -f /etc/haproxy/haproxy.cfg`


## 好友关系 ##
采用redis  

	分组用户好友关系(set) relation:`gid`:`uid`=>set("1","2","3","4","4","5") # gid:group id, uid:用户的id
    用户分组(set) relation:gids:`uid`=>set("0","1","2","同事") 
	注: 分组名不能出现:号,'0'是默认分组

	1.用户 36 添加分组'同事'  ('0'是默认分组)
	is_exist = 0
	if `SISMEMBER relation:gids:36 同事`:
		is_exist = 1
	else:
		`SADD relation:gids:36 同事`

	2.用户 36 删除分组'同事'  ('0'是默认分组)
	is_exist = 1
	if not `SISMEMBER relation:gids:36 同事`:
		is_exist = 0
	else:
		# 删除分组把好友移动到分组0
		`SUNIONSTORE relation:0:36 relation:0:36 relation:同事:36`
		`DEL relation:同事:36`
		`SREM relation:gids:36 同事`

	3.用户 36 添加好友 38 到分组 '同事'
	gid = '同事' or '0'
	gids = `SMEMBERS relation:gids:36`
	is_exist = 0
	for i in `gids`:
		if `SISMEMBER relation:`gid`:36 38`:
			if i != gid:
				`SREM relation:`gid`:36 38`
			is_exist = 1
	`SADD relation:`gid`:36 38`

	4.用户 36 (从分组 2) 删除好友 38
	考虑到客户端可能不发送分组名字2,需要遍历来删除
	for gid in `gid`:
		`SREM relation:`gid`:36 38` 

    5.用户 36 获取所有好友
    groups = []
    for gid in `gids`:
        group = {}
        group["gid"] = `gid` # gid-group id
        group["members"] = `SMEMBERS relation:`gid`:36`
        groups.append(group)

	#TODO: 用户1和2的共同好友,用户1的好友数等等    
	#TODO 聊天室功能

## 个人信息 ##
采用redis

    用户信息,负责登陆等(hash) 
	user:`uid`=>password=>`password`
	user:`uid`=>nickname=>`nickname`
	...
    1.用户36,密码'1123'登陆认证
    `GET user:36` == '1123'

## 状态 ##
采用redis key过期设置3天

	用户所在服务器(kv) status:uid:`uid`=>`sid` # sid:server id
	服务器所有用户(zset) status:sid:`sid`=>uid=>now # uid作为member, now作为score(建立连接时的时间戳)
	
	1.用户36在服务器2上线
	`SET status:uid:36 2`
	`ZADD status:sid:2 `now` 36`

	2.用户36所在服务器(是否在线)
	`GET status:uid:36`

	3.用户36下线
	sid = `GET status:uid:36`
	`ZREM status:sid:`sid` 36`
	`DEL status:uid:36`

	4.服务器2的用户数
	`ZCARD status:sid:2`
	
	5.服务器所有用户数
	`KEYS status:uid*`

	# TODO: 其他状态统计信息

## 离线消息 ##
采用redis

	用户到xx的离线消息(hash) msg:offline:`to_uid`=>`from_uid`=>`message`, 过期时间1年
	to_uid:见下面的消息协议说明
	message:(json)
	{
		"line":0,
		"st":0, # 发送时间,unix时间戳(单位ms, uint64)
		"flg1":byte,
		"flg2":byte,
		"flg3":byte,
		"flg4":byte,
		"ctx":"", # 消息内容
	}

## 在线消息 ##
暂不保存

## 消息协议 ##
中文采用utf8编码  
消息长度(4bytes)+消息类型(byte)+消息体
采用大端对齐,下面写法虽然是字典,但是是按顺序的二进制编码


- 握手 shake  加密只加密消息类型和消息体,握手消息不加密
        
        消息类型: 0

        <=客户端发送握手准备消息 {"subtype":0(byte)}
        =>服务器发送key给客户端 {"subtype":1(byte), "key":uint32}
        <=客户端响应握手成功 {"subtype":2(byte)}

- 心跳 heartbeat

        消息类型: 1

        <=客户端发送心跳 消息体为空
        服务端不用回应(向hub服务器广播该用户在线,心跳90s, 服务端等待300s断开客户端连接)

- 认证 auth  认证时要检查是否是已经登录过了,踢掉原连接,清除状态(先下线后上线)

        消息类型: 2

        <=客户端发送用户名密码 {"uid":36, "password":""}
        =>服务器响应认证结果 {"code":byte(0-ok,其它-ng)}

- 请求离线消息 check offline message (同时实现http)

		消息类型: 3 

		<=客户端请求离线消息 请求最大条数uint32, 0-全部请求 (服务器是队列存储,请求过就删除)
		=>服务端发送离线消息(group by from_id, sort by send_time)
		to(uin32)+<from1+消息1(见redis离线消息字段)>+"\xef\xff"+<from2+消息2(见redis离线消息字段)>..., 每条消息用 \xef\xff 分割

- 聊天 chat (发送消息同时实现http))

        消息类型: 4

        服务端负责转发消息的,如果目标不在线/发送失败,要把消息的online改为0存储待转发

    	{
			"to":38, # target_uid    uint32, [1, 50)系统预留, [50, 150)聊天室, [150, 1000)预留, [1000, +∞)用户
			"from":36, # uid    uint32, [1, 50)系统预留, [50, 150)聊天室, [150, 1000)预留, [1000, +∞)用户
			"line":byte, # 0-离线消息, 1-在线消息
			"st":0, # 发送时间,unix时间戳(单位ms, uint64)
			"flg1":byte,
			"flg2":byte,
			"flg3":byte,
			"flg4":byte,
			"ctx":"", # 消息内容
    	}
        
- ×通知    `添加好友会有通知消息发给对方,等等` 
	
		消息类型: 6

		消息体同聊天消息

- 命令 cmd ((同时实现http))
    
        消息类型: 5
        
        添加分组
        <=
        {
            "type": "group",
            "group":"同事",
            "action":"add"
        }
        =>服务器结果 {"code":byte(0-ok,其它-ng)}
        
        删除分组
        <=
        {
            "type":"group",
            "group":"同事",
            "action":"del"
        }
        =>服务器结果 {"code":byte(0-ok,其它-ng)}
        
        添加用户到分组
        <=
        {
            "type":"user",
            "user":38,
            "group":"同事",
            "action":"add"
        }
        =>服务器结果 {"code":byte(0-ok,其它-ng)}
        
        删除用户
        <=
        {
            "type":"user",
            "user":38,
            "action":"add"
        }
        =>服务器结果 {"code":byte(0-ok,其它-ng)}
        
        获取所有好友及在线状态 (在线状态可以不获取,定时广播)
        <=
        {
            "type":"user",
            "action":"getall"
        }
        =>
        {
            "code":byte, (0-ok,其它-ng)
            "groups":
            [
				{
                	"gid":"同事",
                	"members":[],
				}
            ]
        }

- (续)命令 二进制形式

		{
			"type":byte,
			"action":byte,
			"user":uint32,
			"group":string(固定32)
		}
		type:0-group, 1-user
		action:0-add, 1-del,2-getall
		
		getall的响应消息
		code(byte)+gid(32bytes)+uid1+uid2+..+"\xef\xff"+gid(32bytes)+...

- 其他消息一律过滤

## IPC消息 ##

- 用户的上线下线在hub服务器完成,用户上线后广播到各服, 认证部分在hub服务器完成
- 跨服聊天消息转发