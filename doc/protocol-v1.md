# 聊天服务器协议 #

`所有时间为unix时间戳且为整数,单位ms,其中uid(user id),gid(group id),sid(server id)是数字(在redis中是字符串),0是默认分组`  
`now = int(time.time()*1000)`  
`聊天消息不能含连续的\a\r\n`  
<pre>聊天消息中 [` `] 里面为超链接 </pre>
`st 发送时间,是服务器时间,不是客户端时间,值等于`now``

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

	分组用户好友关系(set) relation:`uid`:`gid`=>set("1","2","3","4","4","5")
    用户分组(hash) group:`uid`=>`gid`=>gname('同事')
    gids = `HKEYS group:`uid``
    gnames = `HVALS group:`uid``    

	1.用户 36 添加分组 2:'同事' 
	is_exist = 0
	if `HEXISTS group:36 2`:
		is_exist = 1
	else:
		`HSET group:`uid` 2`

	2.用户 36 删除分组 2
	is_exist = 1
	if not `HEXISTS group:36 2`:
		is_exist = 0
	else:
		# 删除分组把好友移动到分组0
		`SUNIONSTORE relation:36:0 relation:36:0 relation:36:2`
        `DEL relation:36:2`
		`HDEL group:36 2`

	3.用户 36 添加好友 38 到分组 2
	gid = 2 or 0
	`SADD relation:36:`gid` 38`

	4.用户 36 (从分组 2) 删除好友 38
    if not gid:
    	for _id in `gids`:
    		`SREM relation:36:`_id` 38`
    else:
        `SREM relation:36:2 38`

    5.用户 36 获取所有好友
    groups = []
    for gid in `gids`:
        group = {}
        group["gid"] = `gid`
        group["gname"] = `HGET group:36 `gid``
        group["members"] = `SMEMBERS relation:36:`gid``
        groups.append(group)

	#TODO: 用户1和2的共同好友,用户1的好友数等等    
	#TODO 聊天室功能,用订阅/发布实现

## 个人信息 ##
采用redis key过期设置15天 持久化数据在mysql/mongodb里

    用户信息,负责登陆等(hash) 
	user:`uid`=>password=>`password`
	user:`uid`=>nickname=>`nickname`
	...
    1.用户36,密码'1123'登陆认证
    `GET user:36` == '1123'

## 状态 ##
采用redis key过期设置3天

	用户所在服务器(kv) status_uid:`uid`=>`sid`
	服务器所有用户(zset) status_sid:`sid`=>uid=>now # uid作为member, now作为score(建立连接时的时间戳)
	
	1.用户36在服务器2上线
	`SET status_uid:36 2`
	`ZADD status_sid:2 `now` 36`

	2.用户36所在服务器(是否在线)
	`GET status_uid:36`

	3.用户36下线
	sid = `GET status_uid:36`
	`ZREM status_sid:`sid` 36`
	`DEL status_uid:36`

	4.服务器2的用户数
	`ZCARD status_sid:2`
	
	5.服务器所有用户数
	`KEYS status_uid*`

	# TODO: 其他状态统计信息

## 离线消息 ##
采用redis
	
	(list) msg_offline:`to_uid`=>`message`
	过期时间1年
	定期检查,最大1000条
	message:(json),下面是必须字段
	{
		"from_id":uint32,
		"line":byte, # 0-离线消息, 1-在线消息
        "gid":uint32,
		"st":uint64,
		"ctx":"", # 消息内容
	}

## 在线消息 ##
存储在mysql,定期删除


## TCP消息协议 ##
中文采用utf8编码  
消息长度(4bytes,表示后续数据长度,不包括自己的4bytes)+消息类型(byte)+消息体
采用大端对齐,下面写法虽然是字典,但是是按顺序的二进制编码


- 握手 shake  加密只加密消息类型和消息体,握手消息不加密
        
        消息类型: 0

        <=客户端发送握手准备消息 {"subtype":0(byte)}
        =>服务器发送key给客户端 {"subtype":1(byte), "key":uint32}
        <=客户端响应握手成功 {"subtype":2(byte)}

- 心跳 heartbeat

        消息类型: 1

        <=客户端发送心跳 消息体为空
        =>服务端回复心跳 消息体为空 (向hub服务器广播该用户在线,心跳60s, 服务端等待300s断开客户端连接)

- 认证 auth  认证时要检查是否是已经登录过了,踢掉原连接,清除状态(先下线后上线)

        消息类型: 2

        <=客户端发送用户名密码 {"uid":36, "password":""}
        =>服务器响应认证结果 {"code":byte(0-ok,其它-ng)}

- 请求离线消息 check offline message

		消息类型: 3 

		<=客户端请求离线消息 请求最大条数uint32, 0-全部请求 (服务器是队列存储,请求过就删除)
		=>服务端发送离线消息(group by from_id, sort by send_time)
		to(uin32)+<消息1>+"\a\r\n"+<消息2>..., 每条消息用 \a\r\n 分割

- 聊天 chat

        消息类型: 4

        服务端负责转发消息的,如果目标不在线/发送失败,要把消息的online改为0存储待转发

    	{
			"to":38, # target_uid    uint32, [1, 50)系统预留, [50, 150)聊天室, [150, 1000)预留, [1000, +∞)用户
			"from":36, # uid    uint32, [1, 50)系统预留, [50, 150)聊天室, [150, 1000)预留, [1000, +∞)用户
			"line":byte, # 0-离线消息, 1-在线消息
            "gid":uint32,
			"st":uint64,
			"ctx":"", # 消息内容
    	}

- 命令 cmd 暂不实现tcp的

        消息类型: 5

- 订阅发布
采用管道 把管道作为map的key, 字典的value是一个set(里面存conn/user), 订阅此管道后,发布者发布时,服务器只需要遍历该key的set,然后循环发送给订阅者即可

- 其他消息一律过滤


## HTTP消息协议 ##
中文采用utf8编码  
http服务器同时连到hub服务器  
http服务器接受客户端发送的消息 发送给 hub服务器和redis存储在msg_web_online  
http服务器推送给客户端的消息来源 msg_offline和msg_web_online  
`-H X-GOCHAT-TOKEN: uuid` uuid-其它接口登陆后分配,所有http请求头都需要带此token,token有效期为3小时  
`-H X-GOCHAT-UID: uid`

	code:
	1-uid/password 错误
	2-token失效,如果token失效,需要调用登陆接口,重新换取token,简单的办法是2小时登陆一次,刷新token

- 登陆

		POST "/chat/login"
		<=
		{
			"uid":uint32,
			"password":"",
		}
		=>
		{
			"code":byte(0-ok,其它-ng),
			"token":"26c0ac40-4222-11e4-861a-b8ee657d7c26", # uuid
		}

- 命令 cmd

        1)添加分组
        POST "/chat/groups"
        <=
        {
            "gname":"同事",
        }
        =>
        {
            "code":byte,(0-ok,其它-ng)
            "gid":2
        }
        
        2)删除分组
        DELETE "/chat/groups/gid"
        =>{"code":byte(0-ok,其它-ng)}
        
        3)添加好友到分组
        POST "/chat/friends"
        <=
        {
            "uid":38,
            "gid":2,
        }
        =>{"code":byte(0-ok,其它-ng)}
        
        4)从分组删除好友
        DELETE "/chat/friends/uid"
        <=
        {
            "gid":2,
        }
        =>{"code":byte(0-ok,其它-ng)}
        
        5)获取所有好友
        GET "/chat/friends"
        =>
        {
            "code":byte, (0-ok,其它-ng)
            "groups":
            [
				{
                	"gid":2,
                    "gname":'同事',
                	"members":[],
				}
            ]
        }

- send(客户端发送) 

        POST "/chat/message/send"
        <=
        {
            "to":38,
            "from":36,
            "gid":2, # 发往对方哪个分组的
			"st":0, # 客户端可以传0值,最好传`now`
            "ctx":"hello",
        }
        =>{"code":byte(0-ok,其它-ng)}


- push(服务器推送)(轮询调用此接口拉取消息)

        POST "/chat/message/push"
		<=
		{
			"max_num":0, # 请求最大条数, 0为不限制
		}
        =>
        {
			"code":byte,(0-ok,其它-ng)
			"interval":2000, # 下次轮询时间间隔, 由服务端动态决定,客户端及时刷新此值,单位ms

            "to":38,
            "msgs":
            [
                {
                    "from":36,
					"line":byte, # 0-离线消息, 1-在线消息
                    "gid":2,
                    "st":`now`,
                    "ctx":"hello",
                },
            ]
        }

- 转存图片

		注意headers
		客户端转存成功后获得转存的url后,此url嵌入到聊天消息, eg: 嗨,看这个图[`/static/chat/img/xxx.jpg`]图中ooxx是那个意思

		POST "/chat/image/upload", name=>image
		=>
		{
			"code": 0, 
			"image":"",
			"thumb":"",
		}

- 转存语音

		注意headers

		POST "/chat/sound/upload", name=>sound
		=>
		{
			"code": 0, 
			"sound":"",
		}

- 自定义表情

		表情存在服务器端,以url方式提供,客户端最好缓存
		表情映射表
		/cy:static/chat/expression/cy.jpg

## IPC消息 ##

- 用户的上线下线在hub服务器完成,用户上线后广播到各服, 认证部分在hub服务器完成
- 跨服聊天消息转发