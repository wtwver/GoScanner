printf "#!/bin/sh 
while read line; do 
null=\$(arjun -u \$line -oT temp.txt -c 450 -t 50) 
cat temp.txt 
done < \$1" > run.sh