# inql
# arjum
# https

target=demo.testfire.net
target_ip=192.168.30.225

fyp_pcrawler:
	curl pcrawl.hb1.workers.dev/?url=$(target) -s | tee self_pcrawler.txt
fyp_acrawler:
	go run aCrawler.go -p -u $(target) | tee self_acrawler.txt
fyp_afuzzer:
	go run aCrawler.go -w test.txt -p -u $(target) -pc self_pcrawler.txt| tee self_afuzzer1.txt
	uddup -u self_afuzzer1.txt -s | unew | sort -r | grep -v -E '\.war|\.gif|\.png|-POST-' | head -n 50 | tee self_afuzzer.txt 

fyp_parafuzzer:
# http://192.168.192.20/wavsep/active/Reflected-XSS/RXSS-Detection-Evaluation-GET/Case01-Tag2HtmlPageScope.jsp?userinput=123
# ~/go/bin/ffuf -ac -u $(target)?FUZZ=123 -w test.txt -s
# docker run -i --rm -v ${PWD}/self_afuzzer.txt:/afuzzer.txt wtwver/arjun /afuzzer.txt | tee parafuzzer.txt
	docker run -i --rm -v ${PWD}/self_afuzzer.txt:/afuzzer.txt wtwver/arjun /afuzzer.txt | tee parafuzzer.txt
fyp_parafuzzer2:
	docker run -i --rm -v ${PWD}/self_afuzzer2.txt:/afuzzer.txt wtwver/arjun /afuzzer.txt | tee parafuzzer.txt

fyp_self_scanner:
	python3 scanner.py -i parafuzzer.txt
fyp_dalfox:
	docker run -i --rm -v ${PWD}/parafuzzer.txt:/afuzzer.txt hahwul/dalfox /app/dalfox file /afuzzer.txt --silence
fyp_xsstrike:
	docker run --rm -i -v ${PWD}/self_afuzzer1.txt:/afuzzer.txt femtopixel/xsstrike --seeds /afuzzer.txt --params -t 100 --console-log-level VULN



# for framework scan
self_pcrawler:
	curl pcrawl.hb1.workers.dev/?url=$(target) -s | tee self_pcrawler.txt
self_acrawler:
	go run aCrawler.go -p -u $(target) | tee self_acrawler.txt
self_afuzzer:
	go run aCrawler.go -w test.txt -p -u $(target) -pc self_pcrawler.txt| tee self_afuzzer.txt
self_scanner:
	python3 scanner.py


nuclei:
	docker run --rm -i --entrypoint sh projectdiscovery/nuclei -c "nuclei -update-templates 2>/dev/null 1>&2 ; echo https://$(target) | nuclei -t / -silent"

zap:
	docker run --rm -i owasp/zap2docker-stable zap-cli quick-scan --spider -r --self-contained --start-options '-config api.disablekey=true' https://$(target)

nikto:
	docker run --rm kalo/nikto2 /usr/local/nikto/nikto.pl -host https://$(target)

inql:
	inql -t $(target)

iis_shortname_scanner:
	docker run --rm -it smarticu5/iis_shortname_scanner https://$(target)

wpscan:
	docker run --rm -it --rm wpscanteam/wpscan --url https://$(target)

bfac:
	bfac/bfac --url http://$(target)

dalfox:
	docker run -i --rm hahwul/dalfox /app/dalfox url http://$(target) --silence
dalfox_test:
	docker run -i --rm hahwul/dalfox /app/dalfox url "$(target)" --silence



gau:
	docker run -i --rm wtwver/gau gau $(target)

gospider:
	gospider -q -s $(target)

hakrawler:
	docker run --rm -it wtwver/hakrawler hakrawler -plain -url https://$(target) 

waybackurls:
	~/go/bin/waybackurls $(target)

headi:
	headi -u https://$(target)

ffuf_para:
		~/go/bin/ffuf -ac -u $(target)?FUZZ=123 -w test.txt -s
ffuf_header:
		~/go/bin/ffuf -ac -u https://$(target) -H "FUZZ: 127.0.0.1" -w test.txt
ffuf:
		~/go/bin/ffuf -ac -u https://$(target)/FUZZ -w test.txt

arjun:
	arjun -i /tmp/afuzzer.txt -o /tmp/aparam.json

xsstrike:
	docker run --rm -i femtopixel/xsstrike -u http://$(target) --fuzzer --params -t 100
xsstrike_test:
	docker run --rm -i femtopixel/xsstrike -u "$(target)" -t 100

# docker run --rm -ti -v ${PWD}/delme_xssstrike.txt:/afuzzer.txt femtopixel/xsstrike --seed /afuzzer.txt
# docker run --rm -ti femtopixel/xsstrike -u http://demo.testfire.net/search.jsp
# docker run --rm -ti femtopixel/xsstrike -u http://demo.testfire.net/search.jsp --params



# openredirex:
#	python3 ~/OpenRedireX/openredirex.py -u "https://$(target)/?url=FUZZ" -p payloads.txt --keyword FUZZ
# 
smuggler:
	smuggler/smuggler.py -u https://$(target)

intrigue-ident:
	docker run --rm -t intrigueio/intrigue-ident -v -u https://$(target)

whatweb:
	docker run --rm guidelacour/whatweb:0.5.5-alpine ./whatweb -a 3 https://$(target)

wafw00f:
	docker run -it --rm jonaslejon/wafw00f:2.1.0 $(target)


favfreak:
	cat urls.txt | python3 /FavFreak/favfreak.py -o favfreak.txt




test :
	docker run --rm alpine echo $(target)

setup: 
	sudo bash sEngine.sh setup

scand :
# docker run --rm -it --rm wpscanteam/wpscan --url https://$(target)
	docker run --rm -it wtwver/hakrawler hakrawler -plain -url https://$(target) 
	docker run -it smarticu5/iis_shortname_scanner https://$(target)  
	docker run --rm -it --entrypoint sh projectdiscovery/nuclei -c "nuclei -update-templates 2>/dev/null 1>&2 ; echo https://$(target) | nuclei -t / -silent"
	docker run --rm -t intrigueio/intrigue-ident -v -u https://$(target)
	docker run --rm guidelacour/whatweb:0.5.5-alpine ./whatweb -a 3 https://$(target)
	docker run -it --rm jonaslejon/wafw00f:2.1.0 $(target)
# docker run -it --entrypoint sh trickest/ffuf -c "wget https://raw.githubusercontent.com/danielmiessler/SecLists/master/Discovery/Web-Content/big.txt && ffuf -u https://$(target)/FUZZ -w big.txt -ac -recursion"
# docker run -it --entrypoint sh trickest/ffuf -c "wget https://raw.githubusercontent.com/danielmiessler/SecLists/master/Discovery/Web-Content/big.txt && ffuf -u https://$(target)/?FUZZ=123 -w big.txt -ac "
# docker run -it --entrypoint sh trickest/ffuf -c "wget https://raw.githubusercontent.com/danielmiessler/SecLists/master/Discovery/Web-Content/big.txt && ffuf -u https://$(target) -w big.txt -H "FUZZ: 127.0.0.1" -ac "
	docker run -it --rm rustscan/rustscan:2.0.0 -a $(target)
	docker run --rm kalo/nikto2 /usr/local/nikto/nikto.pl -host https://$(target)
	docker run --rm -i owasp/zap2docker-stable zap-cli quick-scan --spider -r --self-contained --start-options '-config api.disablekey=true' https://$(target)
	docker run --rm blairy/sslscan $(target)
	docker run --rm wtwver/nmap $(target)
	docker run -it --rm hahwul/dalfox /app/dalfox url https://$(target) --silence
	docker run --rm -it wtwver/scripthunter sudo ./scripthunter https://$(target)
# docker run -it --rm wtwver/feroxbuster -u https://$(target)
scan : 
# inql -t https://$(target)
	~/bfac/bfac --url http://$(target)
	~/go/bin/gau $(target)
	~/go/bin/gospider -q -s https://$(target)   
	~/go/bin/waybackurls $(target)
	~/go/bin/headi -u https://$(target)/
# /Arjun/arjun.py -i /tmp/afuzzer_clean.txt -o /tmp/aparam.json
	~/XSStrike/xsstrike.py -u $(target)
	~/smuggler/smuggler.py -q -u https://$(target)
	echo https://$(target) | python3 ~/FavFreak/favfreak.py -o favfreak.txt

xss:
	docker run -it --rm hahwul/dalfox /app/dalfox url http://$(target) --silence
# docker run --rm -ti femtopixel/xsstrike -u http://$(target) --params --fuzzer -t 100
xss1:
	docker run --rm -ti femtopixel/xsstrike -u http://$(target) --fuzzer --params -t 100

finger:
	docker run --rm -t intrigueio/intrigue-ident -v -u https://$(target)
	docker run --rm guidelacour/whatweb:0.5.5-alpine ./whatweb -a 3 https://$(target)
	docker run -it --rm jonaslejon/wafw00f:2.1.0 $(target)
	docker run --rm -it wtwver/hakrawler hakrawler -plain -url https://$(target) 
	docker run --rm blairy/sslscan $(target)

port: 
# docker run --rm wtwver/nmap $(target_ip)
	docker run -it --rm rustscan/rustscan:2.0.0 -a $(target_ip)