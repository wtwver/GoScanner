# input 
# -i self_afuzzer.txt
# -u http://192.168.192.20/wavsep/active/Reflected-XSS/RXSS-Detection-Evaluation-GET/Case01-Tag2HtmlPageScope.jsp?userinput=123
# dos2unix $1 
while read line; do
# echo $line
null=$(arjun -u $line -oT temp.txt -c 450 -t 50)
cat temp.txt
done < $1