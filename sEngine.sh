#!/bin/bash

# target="http://demo.testfire.net"
target=$2

run_nuclei(){   
    # echo $target | nuclei -t / -update-templates #-silent
    ~/go/bin/nuclei -update-templates  2> /dev/null
    ~/go/bin/nuclei -u https://$target -t technologies/ -silent 
}

all(){
    echo "runing all"
    run_nuclei
}

setup(){
    GO111MODULE=on go get github.com/projectdiscovery/nuclei/v2/cmd/nuclei
    go get github.com/lc/gau 
    go get github.com/hakluke/hakrawler
    go get github.com/tomnomnom/hacks/waybackurls
    go get -u github.com/jaeles-project/gospider
    GO111MODULE=on go get -v github.com/hahwul/dalfox/v2
    go get github.com/mlcsec/headi
    go get -u github.com/ffuf/ffuf
    # go get -u github.com/hahwul/jwt-hack
    # GO111MODULE=on go get -u github.com/sw33tLie/sns

    cd ~
    apt install -y python3 python3-pip git
    pip3 install arjun
    # $(cd Arjun) || git clone --depth 1 https://github.com/s0md3v/Arjun.git
    # python3 ~/Arjun/setup.py install
    # chmod +x ~/Arjun/*

    $(cd XSStrike) || git clone --depth 1 https://github.com/s0md3v/XSStrike.git
    pip3 install -r XSStrike/requirements.txt
    chmod +x XSStrike/*

    $(cd OpenRedireX) || git clone --depth 1 https://github.com/devanshbatham/OpenRedireX
    pip3 install aiohttp


    # git clone --depth 1 https://github.com/D35m0nd142/LFISuite.git
    # python2

    $(cd FavFreak) ||  git clone --depth 1 https://github.com/devanshbatham/FavFreak
    pip3 install -r FavFreak/requirements.txt

    $(cd smuggler) || git clone --depth 1 https://github.com/defparam/smuggler.git

    $(cd bfac) || git clone --depth 1 https://github.com/mazen160/bfac.git
    # sudo python bfac/setup.py install
    pip3 install colorama requests requests[socks]

    $(cd inql) || git clone --depth 1 https://github.com/doyensec/inql.git
    pip3 install -r inql/requirements.txt
    cd inql && python3 setup.py install && cd ../
}

# run_jwt-hack(){
#     jwt-hack payload $target
# }
run_zap(){
    docker run --rm -i owasp/zap2docker-stable zap-cli quick-scan --spider -r --self-contained --start-options '-config api.disablekey=true' https://$target
}
run_nikto2(){
    docker run --rm kalo/nikto2 /usr/local/nikto/nikto.pl -host https://$target
}

run_inql(){
    inql -t $target
}
run_iis_shortname_scanner(){
    docker run --rm -it smarticu5/iis_shortname_scanner https://$target
}
run_wpscan(){
    docker run --rm -it --rm wpscanteam/wpscan --url https://$target
}
run_bfac(){
    bfac/bfac --url http://$target
}
run_dalfox(){
    ~/go/bin/dalfox url https://$target
}

run_gau(){
    # ~/go/bin/gau $target
    docker run -i --rm wtwver/gau gau $target
}
run_gospider(){
    gospider -q -s $target
}
run_hakrawler(){
    ~/go/bin/hakrawler -url $target
}
run_waybackurls(){
    ~/go/bin/waybackurls $target
}
run_headi(){
    headi -u https://$target
}

run_ffuf(){
    if [[ $1 = "para" ]]
    then
        ~/go/bin/ffuf -ac -u https://$target/?FUZZ=123 -w test.txt
    elif [[ $1 = "header" ]]
    then
        ~/go/bin/ffuf -ac -u https://$target -H "FUZZ: 127.0.0.1" -w test.txt
    else
        ~/go/bin/ffuf -ac -u https://$target/FUZZ -w test.txt
    fi
}
run_arjun(){
    /Arjun/arjun.py -i /tmp/afuzzer_clean.txt -o /tmp/aparam.json
}
run_xsstrike(){
    /XSStrike/xsstrike.py -u $target
}
# run_openredirex(){
#     python3 ~/OpenRedireX/openredirex.py -u "https://$target/?url=FUZZ" -p payloads.txt --keyword FUZZ
# }
run_smuggler(){
    smuggler/smuggler.py -u https://$target
}
run_intrigue-ident(){
    docker run --rm -t intrigueio/intrigue-ident -v -u https://$target
}
run_whatweb(){
    docker run --rm guidelacour/whatweb:0.5.5-alpine ./whatweb -a 3 https://$target
}
run_wafw00f(){
    docker run -it --rm jonaslejon/wafw00f:2.1.0 $target
}

run_favfreak(){
    cat urls.txt | python3 /FavFreak/favfreak.py -o favfreak.txt
}


if [[ $1 = "all" ]]
then
    all
elif [[ $1 = "setup" ]]
then
    setup
else
    echo run_$1
    run_$1
    # run_nuclei
# elif [[ $1 = "gau" ]]
# then
#     run_gau
# elif [[ $1 = "gospider" ]]
# then
#     run_gospider
# elif [[ $1 = "hakrawler" ]]
# then
#     run_hakrawler
# elif [[ $1 = "waybackurls" ]]
# then
#     run_waybackurls
# elif [[ $1 = "dalfox" ]]
# then
#     run_dalfox
# elif [[ $1 = "ffuf" ]]
# then
#     run_ffuf
fi