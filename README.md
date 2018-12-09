大促类运营活动的开发套路（本工程用beego+k8做微服务化）  
这个工程的主要模块是  
1）业务模块  
  1.1）接入层  
  1.2）service层  
2）初始化模块  
  2.1）go writehandler()  
  2.2) go readhandlerZ()  
3)持久层  
   3.1）热数据持久化技术选型redis  
   3.2）数据一致性，redis做分布式锁  
   3.3）配置信息持久化技术选型etcd，etcd的优势是不断watch  

其中：  
1）业务模块的持久层用redis，数据一致性用redis分布式锁  
2）业务模块运营活动请求req到reqchan协程的管道    
3）初始化模块的writeHandler协程读出reqchan的信息，进行封装，lpush到redis的queue中  
4）初始化模块的readhandler协程从queue中取req做逻辑处理，然后放入resultchan管道中  
5）业务模块非阻塞的读取resultchan中的返回信号，返回接入层  
  
总结：  
本工程的业务主要是做了首页的逻辑，没有把整个大促活动逻辑放到工程里面，主要关注点的是运营活动的处理思路。  
