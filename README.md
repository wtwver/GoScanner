## Known issue
1. [affuzer] cannot handle server with rate limiting -> fix using distributed? -> fix by retry null request
2. [affuzer] testphp.vulnweb.com 404 differnt content length in response 16 & 153
git config credential.helper store

## DEMO
http://206.189.90.86/dashboard.html
http://206.189.90.86/api/gitpull

# Dependence win10
```
1. https://chocolatey.org/install
2. choco install make
```

## How to Run
Native / VPS 
```
docker run -d --restart always --memory="0.4g" --name fypmysql -p 1000:3306 -e MYSQL_ROOT_PASSWORD=password@123 -e MYSQL_DATABASE=test mysql
// no idea why password env doesnt work
docker run -d --restart=always --publish 6379:6379 --env 'REDIS_PASSWORD=redispassword' redis    

<!-- go get github.com/gocolly/redisstorage
cd ~/go/src/github.com/go-redis/redis
git checkout tags/v6.14.2 -->

go install .

./sEngine.sh setup
go run ./database.go
```

Vargant
```
Vargant up
```

<!-- Docker compose
```
git clone https://github.com/JP-934/Web-API-Vulnerability-Scanner.git
docker build -t backend .
docker-compose up 
``` -->


# Scanning engine

Docker in docker is a must for ruby tools
- intrigue-ident


Consideration
1. non performaance tools can docker 
2. container size <50MB ?
3. no exisiting container avaible
4. python dependency crash
5. security
6. latest commit wont work

docker for 
1. go tools, all go tools have binary release? 
2. pytohn tools
3. all tools = scnaning engine docker

python2 pip2 can apk add 
- http://dl-cdn.alpinelinux.org/alpine/v3.11/main/x86_64/py2-pip-18.1-r0.apk

golang , python3 native install

other language
ruby : docker
python2: docker for secuirty reason
javascript


# Crawling engine
crawl recursion on 200, out 200 403 500 
fuzz recursion on 200, out 200 403 500 


crawl(juiceshop.com) => crawlout{unqiuewords, deadabsolute, deadrelative}

fuzz(crawlout, big.txt)
- pathsplit(crawlout) , fuzz()
- ColletUnqiuewords during fuzz
- crawl found endpoints => crawlout => pathsplit(crawlout) => fuzz

=> fuzz out

TODO
1. acrawler json output
2. pathsplit
3. import function from acrawler
4. Recursion

# Frontend

![](https://i.imgur.com/TjkqJmC.png)
![](https://github.com/plenumlab/lazyrecon/raw/dev/upgrade/report.gif)


1. select tools
2. send request to backend 

```
POST http://backend/api/commandHandler

{
    "seed":"juiceshop.com"
    "pcrawler": "A"
    "acrawler": "B"
    "afuzzer": "self"
}
```
3. View discovered endpoints
```
GET http://backend/api/getAll
GET http://backend/api/getDomains
GET http://backend/api/getFolders
GET http://backend/api/getEndpoints
```


# API backend 
1.  JSON to command 

```
./pcrawler.sh A google.com | ./acrawler.sh B | ./afuzzer.sh self 
```

2. Run command and store result in database


# Potential function to add

Docker command new tools
Schedule scan
Sitemap (present only, have time sin )