# recipes_api

## MongoDB
```bash
docker run -d --name mongodb -v /Users/<user>/db:/data/db -e MONGO_INITDB_ROOT_USERNAME=<uname> -e MONGO_INITDB_ROOT_PASSWORD=<password> -p 27017:27017 mongo
```

## Redis
```bash
docker run -d -v $PWD/conf:/usr/local/etc/redis --name redis -p 6379:6379 redis
```
To connect to redis-cli:
```bash
docker exec -it redis redis-cli
```

For GUI access to redis:
```bash
docker run -d --name redis-commander -p 8081:8081 -e REDIS_HOSTS=localhost:6379 rediscommander/redis-commander
```
or

```bash
docker run -d --name redisinsight --link redis -p 8001:8001 redislabs/redisinsight
```

## Application
```bash
MONGO_URI="mongodb://<uname>:<password>@localhost:27017/test?authSource=admin" MONGO_DATABASE=demo go run main.go
```

## Benchmarking
```bash
ab -n 2000 -c 100 -g with-cache.data http://localhost:8080/recipes
```

```text
Server Software:        
Server Hostname:        localhost
Server Port:            8080

Document Path:          /recipes
Document Length:        722318 bytes

Concurrency Level:      100
Time taken for tests:   25.155 seconds
Complete requests:      2000
Failed requests:        0
Total transferred:      1444842000 bytes
HTML transferred:       1444636000 bytes
Requests per second:    79.51 [#/sec] (mean)
Time per request:       1257.726 [ms] (mean)
Time per request:       12.577 [ms] (mean, across all concurrent requests)
Transfer rate:          56092.45 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    1   1.2      0      10
Processing:    54 1241 526.7   1203    3657
Waiting:       50 1083 488.4   1028    3623
Total:         54 1242 526.7   1204    3657

Percentage of the requests served within a certain time (ms)
  50%   1204
  66%   1411
  75%   1583
  80%   1681
  90%   1924
  95%   2216
  98%   2447
  99%   2625
 100%   3657 (longest request)
```

To generate a graph:
```bash
brew install gnuplot
gnuplot apache-benchmark.plot
```
