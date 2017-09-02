#!/bin/bash
#按天分割日志
date=` date +%Y-%m-%d`
# 需要分割的日志文件名绝对路径，可以配置多个
#targets="/data/log/ebike-factory-api"
targets="/data/log/ebike-factory-api/logrus.log"
for target in $targets
do
    file=`basename $target`
    basedir=`dirname $target`
    cd $basedir
    echo $basedir/$file
    cat $file >> $file.$date
    cat /dev/null > $file
done
