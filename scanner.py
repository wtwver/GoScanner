# add transformation check
# add custom payload and ART

import requests
import os
from urllib.parse import urlparse, parse_qs
import random
import argparse

special_chars = ['!', '@', '#', '$', '%', '%',
                 '^', '&', '*', '*', '*', '(', ')', '_', '+','\'','"']

proxyDict = {
    "http": "http://127.0.0.1:1000",
    "https": "http://127.0.0.1:1000",
    "ftp": "http://127.0.0.1:1000"
}

def inject(url):
    print("[*] injecting: ", url)
    # print( dir(requests.get(url)))
    print("orginal resp: ", len(requests.get(url, verify=False).text))
    print("===================")

    # params = []
    a = urlparse(url)
    base_url = url.split("?")[0]
    base_resp = requests.get(base_url, verify=False)

    print("base resp: ", len(base_resp.text))
    print("===================")
    params = parse_qs(a.query)

    # print("params: ",params)

    for param, value in params.items():
        
        temp = params.copy()
        rnd = str(random.randint(0,100))
        temp[param] = value[0] + rnd
        reflect_resp = requests.get(base_url, params=temp, verify=False)
        if rnd in reflect_resp.text:
            print("reflection detected in __ ", param,
                  "__ with resp: ", len(reflect_resp.text))
            print("===================")

        for special_char in special_chars:
            copy = params.copy()
            copy[param] = value[0] + special_char
            # print(copy[param])

            # inject_resp = requests.get(base_url, params=copy, proxies=proxyDict)
            inject_resp = requests.get(base_url, params=copy, verify=False)
            if interesting(inject_resp, reflect_resp):
                print("injected", special_char ," with resp: ",len(inject_resp.text) )
                print("===================")

# inejcted leng not same
# char count
# word count
# line count


def interesting(resp1, resp2):
    return ( len(resp1.text) != len(resp2.text) and 
    len(resp1.text.split()) != len(resp2.text.split()) )  
    # len(resp1.text.split("\n")) != len(resp2.text.split("\n")) )
    

# input = []

parser = argparse.ArgumentParser()
parser.add_argument("-i","--input")
args = parser.parse_args()

# if args.a == 'magic.name':
#     print 'You nailed it!'

with open(args.input, "r") as f:
    for line in f:
        # input.append(line.strip())
        inject(line)

# parameters = []
# for line in input:
#     a = urlparse(line)
#     parameters.append(parse_qs(a.query))

# for i in parameters:
#     print(i)

# inject("http://192.168.192.1/wavsep/active/SQL-Injection/SInjection-Detection-Evaluation-GET-500Error/Case01-InjectionInLogin-String-LoginBypass-WithErrors.jsp?a=1&username=textvalue")

# inject("http://192.168.192.20/wavsep/active/Reflected-XSS/RXSS-Detection-Evaluation-GET/Case01-Tag2HtmlPageScope.jsp?userinput=textvalue")




# inject("http://192.168.192.20/wavsep/active/Reflected-XSS/RXSS-Detection-Evaluation-GET/Case01-Tag2HtmlPageScope.jsp?userinput=textvalue&b=1")

