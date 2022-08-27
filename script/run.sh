baseDir = $(cd $(dirname $0);pwd)/..
output = $baseDir/output

if [! -d output]; then
    mkdir output

cp -rf $baseDir/config $output/config
cp -rf $baseDir/file $output/file

go build -o $output/rankland $baseDir
sh $output/rankland