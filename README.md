构建了一个使用Go、 redis 和 HTMX 实现的 URL 缩短工具，它非常快速、简单且实用。Go 和 Redis 是两种关键技术，使 Stream 对数百万终端用户来说既快速又可靠。如何使用它们来构建快速且可扩展的 Web 应用有了基本的了解。



### wrk 测试

wrk -t2 -c50 -d15s http://服务地址:端口

Running 15s test @ http://服务地址:端口

2 threads and 50 connections

Thread Stats   Avg      Stdev     Max   +/- Stdev

Latency     1.29ms  134.44us   5.78ms   72.77%

Req/Sec    19.27k   588.36    20.51k    78.67%

575114 requests in 15.00s, 601.67MB read

Requests/sec:  38339.45

Transfer/sec:     40.11MB

