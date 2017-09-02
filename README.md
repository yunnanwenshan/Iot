## Iot 

一个智能设备服务器端集群解决方案。
Iot内部由node和router组成，一个router管理多个node，并实现node负载均衡。
数据传输支持加密。

协议格式
------------
```
	bit        7        6        5        4        3        2        1       0
	byte1	EnCode(2bit)	Message-Type(6bit)
	byte2	Transaction  ID(4bytes)
	byte3	
	byte4	
	byte5	
	byte6	Body-len(4bytes 可选)
	byte7	
	byte8	
	byte9	
	byte...	消息体（可选）
```

消息类型
-------------
#### 消息类型	值	描述
	Reserved	0	保留
	REGISTER	1	注册，Device发起注册
	REGRESP	2	注册应答
	PING	3	客户端发起心跳消息，心跳间隔建议300S
	PONG	4	心跳应答
	IM	9	即时消息请求
	IMRESP	10	即时消息应答
EnCode
-------------
	指定消息体编码（加密）方式
	0x00：默认，无特殊编码
	0x01：按位取反
	0x02：字节逆序（两两互逆，若单字节最后一字节不做变换）
	0x03：环形异或
	
TID(Transaction ID)
-------------
	unsigned int类型的随机不重复事务ID，初始消息ID随机；请求消息与响应消息拥有相同的事务ID；
	事务默认超时时长：5秒

消息体长度
-------------
	可选，类型为unsigned int类型，当指定消息类型携带消息体时有效。

消息体
-------------
	可选， 当指定消息类型携带消息体时有效
