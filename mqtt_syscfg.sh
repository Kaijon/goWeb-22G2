#!/bin/sh

# Preparing infomation page
info=""
#fw_ver=$(cat /tmp/fw_version)
fw_ver="1.2.3.4"
#serial_no=$(grep "SysSerialNumber" /mnt/flash/etc/system.cfg | cut -d'=' -f2)
#ipaddr=$(ifconfig eth0 | grep 'inet addr:' | cut -d: -f2 | awk '{ print $1}')
#mac_address=$(ifconfig eth0 | grep 'HWaddr' | awk '{ print $5}')

info="FW Version=${fw_ver} <br>"
#info="${info}Serial Number=${serial_no} <br>"
#info="${info}IP Address=${ipaddr} <br>"
#info="${info}MAC Address=${mac_address} <br>"

# Read values from system.cfg for factory page
#model_name=$(grep "ModelName" /mnt/flash/etc/system.cfg | cut -d'=' -f2 | tr -d '"')
#mac_address=$(grep "MacAddress" /mnt/flash/etc/system.cfg | cut -d'=' -f2 | tr -d '"')
#pcba_serial_number=$(grep "PCBASerialNumber" /mnt/flash/etc/system.cfg | cut -d'=' -f2 | tr -d '"')
#sys_serial_number=$(grep "SysSerialNumber" /mnt/flash/etc/system.cfg | cut -d'=' -f2 | tr -d '"')
#sku_name=$(grep "SkuName" /mnt/flash/etc/system.cfg | cut -d'=' -f2 | tr -d '"')

# publish
echo "$info" | mosquitto_pub -h 127.0.0.1 -p 1883 -t info -l -r
echo $info

# Publish factory values to MQTT topics
#echo "$mac_address" | mosquitto_pub -h 127.0.0.1 -p 1883 -t "factory/info/MAC" -l -r
#echo "$sys_serial_number" | mosquitto_pub -h 127.0.0.1 -p 1883 -t "factory/info/SysSN" -l -r
#echo "$pcba_serial_number" | mosquitto_pub -h 127.0.0.1 -p 1883 -t "factory/info/PCBASN" -l -r
#echo "$model_name" | mosquitto_pub -h 127.0.0.1 -p 1883 -t "factory/info/ModelName" -l -r
#echo "$sku_name" | mosquitto_pub -h 127.0.0.1 -p 1883 -t "factory/info/SkuName" -l -r

