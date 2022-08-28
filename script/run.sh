baseDir=$(cd $(dirname $0); cd ..; pwd)
# echo $baseDir
output=$baseDir/output
# echo $output



if [ ! -d $output ];then
    mkdir $output
fi

cp -rf $baseDir/config $output
cp -rf $baseDir/file $output

go build -o $output/rankland $baseDir 