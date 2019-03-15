#!/bin/bash

PHYBBIT_PROJECT=mojaco-1208
dates=2019030{5..14}

function GCSCertification(){
    local readonly keyFile=spideraf-exchange-geniee.json
    echo GCSCertification
    gcloud auth activate-service-account --key-file $keyFile
}

function scp_hadoop(){
    scp hadoop.sh mapper.py ops2:/home/ops-cdh
}

function hadoop2binary(){
    ssh ops2 "/home/ops-cdh/hadoop.sh $1 $2 $3" | ./phybbit
}

function getLogType(){
    if [ "$1" = "imps" ];then
	echo impression
    elif [ "$1" = "clicks" ];then
	echo click
    else
	echo failed log type
	exit 1
    fi
}

# inpupt $1:date $2:category $3:logtype
# date is YYYYmmdd
# category is imps or clicks
function upload(){
    local bucket=spideraf-exchange-geniee

    local date=$1
    # clicks or imps
    local category=$2
    local logtype=$(getLogType $category)
    for hour in `seq -w 00 23`
    do
	hadoop2binary $date $hour $logtype | gsutil cp - gs://${bucket}/${category}_$date_$hour.avro
    done
}

function main(){
    GCSCertification
    scp_hadoop
    for date in dates
    do
	upload $date clicks
	#upload $date imps
    done
}

function test(){
    tmpdir=tmp
    tmpFile=${tmpdir}/test.log
    mkdir -p $tmpdir
    
    scp_hadoop
    date=20190307
    hour="00"
    logtype="click"
    hadoop2binary $date "${hour}" $logtype > $tmpFile
}

#main
test
