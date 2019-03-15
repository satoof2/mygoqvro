#!/bin/bash

PHYBBIT_PROJECT=mojaco-1208
# 2019/03/05~07は一時保存しよう
dates=2019030{5..7}

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
# $1:date is YYYYmmdd
# $2:category is imps or clicks
function upload(){
    local bucket=spideraf-exchange-geniee

    local date=$1
    local category=$2
    
    local logtype=$(getLogType $category)
    for hour in $(seq -f %02g 00 23)
    do
	echo "upload ${date}${hour}"
	hadoop2binary $date "$hour" $logtype | gsutil cp - gs://${bucket}/${category}/${date}_$hour.avro
    done
}

function _save(){
    local bucket=spideraf-exchange-geniee

    local date=$1
    local category=$2
    local logtype=$(getLogType $category)
    echo $date
    for hour in $(seq -f %02g 00 23)
    do
	ssh ops2 "/home/ops-cdh/hadoop.sh $date $hour $logtype"
    done
}

function main(){
    GCSCertification
    scp_hadoop
    for date in dates
    do
	echo "upload imps"
	upload $date imps
	echo "upload clicks"
	upload $date clicks
    done
}

function test(){
    GCSCertification
    scp_hadoop
    #upload
}

function save(){
    scp_hadoop
    ssh ops2 "hadoop fs -rm -r /user/ops-cdh/K-Honkawa"
    for date in 2019030{5..7}
    do
	echo $date
	echo "save imps"
	_save $date imps
	echo "save clicks"
	_save $date clicks
    done

}
save
#test
