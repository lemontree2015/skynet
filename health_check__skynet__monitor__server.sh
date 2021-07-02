#! /bin/sh  

proc_name="skynet_monitor_server" #[进程名]
  
proc_num()                        #查询进程数量
{
  num=`ps -ef | grep $proc_name | grep -v grep | wc -l`
  return $num
}
 
proc_num
number=$?                         #获取进程数量
echo $number `date "+%Y-%m-%d %H:%M:%S"`
if [ $number -eq 0 ]              #如果进程数量为0
then                              #重新启动服务器, 或者扩展其它内容
  # crontab.template
  # */1 * * * * /home/liqiang/health_check__skynet__monitor__server.sh >> /home/liqiang/health_check__skynet__monitor__server.log 2>&1
  cd /home/liqiang/mygo_projs/src/bitbucket.org/gimcloud/skynet/; ./skynet_ctl.sh monitor_server_nohup # [重启服务的Command]
fi

