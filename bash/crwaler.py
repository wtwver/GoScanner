import  os
# docker inside docker to skip installing?
# should build docker by ourself to enhacne security

def exec(cmd):
    os.system(cmd) 

def create_vol():
    exec("docker volume create result ")

target = "juiceshop321.herokuapp.com"
crawl_depth = 3
fuzz_recur = True
wordlist = "common-api-endpoints-mazen160.txt"

acrawler = "docker run -it --rm -v result:/tmp snowsecurity/photon -u {} --wayback -l {} -t 10 -o /tmp/".format(target,crawl_depth) 
pcrawler =" docker run -it --rm -v result:/tmp heywoodlh/tomnomnom-tools:latest ash -c \"echo {} | waybackurls | sort -u | tee /tmp/2.txt \"".format(target)
afuzzer="docker run -it --rm -v result:/tmp --entrypoint sh trickest/ffuf -c \"wget https://raw.githubusercontent.com/danielmiessler/SecLists/master/Discovery/Web-Content/{} && ffuf -u http://{}/FUZZ -w {} | tee /tmp/3.txt\"".format(wordlist,target,wordlist)

heuristic = "docker run -it --rm -v result:/tmp wtwver/parth1 bash -c \"python3 parth.py -put {} && mv params-* /tmp/ && cat /tmp/params*\"".format(target)

CVE = "docker run --rm -it -v result:/tmp --entrypoint sh projectdiscovery/nuclei -c \"echo http://{} | ./nuclei -t / -silent -update-templates\"".format(target)

result = "docker run --rm -it -v result:/tmp ubuntu bash -c \"cat /tmp/*\" "

create_vol()
exec(acrawler) 
exec(pcrawler) 
exec(afuzzer) 

exec(result)

exec(heuristic)
exec(CVE)

