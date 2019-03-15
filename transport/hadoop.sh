#!/bin/bash

DATE=$1
HOUR=$2
LOG_TYPE=$3

hadoop_log=motokawa_hadoop.log
if [ -e　${hadoop_log} ]; then
    rm ${hadoop_log}
fi

function main(){
    local TEMPORARY_DIR=/user/ops-cdh
    local SRC=/data/aladdin/logs_compressed/${LOG_TYPE}/$DATE/$HOUR/part-*.lzo
    #local DST_NAME=kojiro-honkawa一時保存よう
    local DST_NAME=K-Honkawa 
    local DST=$TEMPORARY_DIR/${DST_NAME}/${LOG_TYPE}/$DATE/$HOUR
    local MAP=mapper.py

    #-----HDFS上の出力ディレクトリ削除-----
    #hadoop fs -rm -r "$DST"
    
    #wait_by_num_jobs
    #-----Hadoop実行-----
    hadoop jar /opt/cloudera/parcels/CDH/lib/hadoop-mapreduce/hadoop-streaming.jar \
           -D mapreduce.job.reduces=0 \
           -D mapreduce.job.priority=LOW \
           -files "$MAP" \
           -input "$SRC" \
           -output "$DST" \
           -mapper "$MAP" >> ${hadoop_log}
    #hadoop fs -cat "$DST"/part-*
}

main
