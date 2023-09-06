# ktoy-mysql(主从模式)

# 启动主节点容器ktoy_mysql_1

前后启动ktoy_mysql_1，ktoy_mysql_2， 启动一个ktoy_mysql_1 配置如下

```
cat > /opt/ktoy/mysql-1/conf <<END
[client]
default-character-set=utf8
[mysqld]
init_connect='SET collation_connection = utf8_unicode_ci'
init_connect='SET NAMES utf8'
character-set-server=utf8
collation-server=utf8_unicode_ci
skip-character-set-client-handshake

log-bin=mysql-bin
server-id=1
binlog-do-db=ktoy
binlog_cache_size=1M
binlog_format=row
binlog_expire_logs_seconds=2592000
replica_skip_errors=1062

[mysql]
default-character-set = utf8
END
```

启动容器 ktoy_mysql_1

```
docker run --name ktoy_mysql_1 -p 13306:3306 -v /opt/ktoy/mysql-1/data:/var/lib/mysql \
  -v /opt/ktoy/mysql-1/conf:/etc/mysql/conf.d -v /opt/ktoy/mysql-1/log:/var/log/mysql \ 
  -e MYSQL_ROOT_PASSWORD=zh08311534 --privileged=true -d icmp.harbor/starocean/mysql:8.0
```



# 启动从节点容器ktoy_mysql_2

ktoy_mysql_2启动配置文件为

```SHELL
[client]
default-character-set=utf8
[mysqld]
init_connect='SET collation_connection = utf8_unicode_ci'
init_connect='SET NAMES utf8'
character-set-server=utf8
collation-server=utf8_unicode_ci
skip-character-set-client-handshake

log-bin=mysql-bin
server-id=2
binlog-do-db=ktoy
binlog_cache_size=1M
binlog_format=row
binlog_expire_logs_seconds=2592000
replica_skip_errors=1062

relay_log=replicas-mysql-relay-bin
log_replica_updates=ON
read_only=ON

[mysql]
default-character-set = utf8
```

**read_only**对拥有super和connection_admin这两个权限的用户无效。

启动容器 ktoy_mysql_2

```SHELL
docker run --name ktoy_mysql_2 -p 13307:3306 -v /opt/ktoy/mysql-2/data:/var/lib/mysql \
 -v /opt/ktoy/mysql-2/conf:/etc/mysql/conf.d -v /opt/ktoy/mysql-2/log:/var/log/mysql \ 
 -e MYSQL_ROOT_PASSWORD=zh08311534 --privileged=true -d icmp.harbor/starocean/mysql:8.0
```



# 所有节点创建数据库

create database ktoy;  在所有节点上创建数据库



# 主节点ktoy_mysql_1配置

## 配置主从同步权限

```sql
CREATE USER 'repl'@'%' IDENTIFIED WITH 'mysql_native_password' BY 'zh08311534';
```


ktoy1 上创建 赋予权限 GRANT REPLICATION SLAVE ON *.* TO 'repl'@'%';

尽量用库名.*去赋权，而非*.*

```SHELL
use mysql;
select user,host,plugin,authentication_string from user;

+------------------+-----------+-----------------------+------------------------------------------------------------------------+
| user             | host      | plugin                | authentication_string                                                  |
+------------------+-----------+-----------------------+------------------------------------------------------------------------+
| repl             | %         | mysql_native_password | *7A1E725011C4D5468ED36DAE56C73DC26DB0CD75                              |
| root             | %         | caching_sha2_password | $A$005$"'R9M>wE/"i%,}"oiRAp3Is4D7otbbF8QAbdchA0P.7.JyQQXsi84nLUz64 |
| mysql.infoschema | localhost | caching_sha2_password | $A$005$THISISACOMBINATIONOFINVALIDSALTANDPASSWORDTHATMUSTNEVERBRBEUSED |
| mysql.session    | localhost | caching_sha2_password | $A$005$THISISACOMBINATIONOFINVALIDSALTANDPASSWORDTHATMUSTNEVERBRBEUSED |
| mysql.sys        | localhost | caching_sha2_password | $A$005$THISISACOMBINATIONOFINVALIDSALTANDPASSWORDTHATMUSTNEVERBRBEUSED |
| root             | localhost | caching_sha2_password | $A$005$(mf:FKaB 
                                                                        3tn0%5mHsj2mCVXYDVdAXPvyW4YQOV8eWQmNvZESL9ha8FgU4 |
+------------------+-----------+-----------------------+------------------------------------------------------------------------+
```

## 参数说明

MySQL创建授权命令：GRANT REPLICATION SLAVE ON database.table TO 'username'@'host'；
privilege ：指定授权的权限，比如create、drop等权限，具体有哪些权限，可查看官网文档
database：指定哪些数据库生效，*表示全部数据库生效
table：指定所在数据库的哪些数据表生效，*表示所在数据库的全部数据表生效
username：指定被授予权限的用户名
host：指定用户登录的主机ip，%表示任意主机都可远程登录

## 刷新权限生效

```sql
FLUSH PRIVILEGES; 
SHOW MASTER STATUS; 查看，获取关键信息
```

```shell
+------------------+----------+--------------+------------------+-------------------+
| File             | Position | Binlog_Do_DB | Binlog_Ignore_DB | Executed_Gtid_Set |
+------------------+----------+--------------+------------------+-------------------+
| mysql-bin.000001 |      1047 | ktoy         |                  |                   |
+------------------+----------+--------------+------------------+-------------------+
```




# 从节点ktoy_mysql_2配置

## 配置源启动

```sql
mysql> CHANGE REPLICATION SOURCE TO SOURCE_HOST='10.11.17.202',SOURCE_PORT=13306,SOURCE_USER='repl',SOURCE_PASSWORD='zh08311534',SOURCE_LOG_FILE='mysql-bin.000001',SOURCE_LOG_POS=1047;
Query OK, 0 rows affected, 2 warnings (0.01 sec)

mysql> START REPLICA;
Query OK, 0 rows affected (0.01 sec)

同步成功：Slave_IO_Running/Replica_IO_Running和 Slave_SQL_Running/Replica_SQL_Running 为 Yes ，
以及Slave_IO_State/Replica_IO_State 为 Waiting for master to send event/Waiting for source to send event，


```



```SHELL
+----------------------------------+--------------+-------------+-------------+---------------+------------------+---------------------+---------------------------------+---------------+-----------------------+--------------------+---------------------+-----------------+---------------------+--------------------+------------------------+-------------------------+-----------------------------+------------+------------+--------------+---------------------+-----------------+-----------------+----------------+---------------+--------------------+--------------------+--------------------+-----------------+-------------------+----------------+-----------------------+-------------------------------+---------------+---------------+----------------+----------------+-----------------------------+------------------+--------------------------------------+-------------------------+-----------+---------------------+----------------------------------------------------------+--------------------+-------------+-------------------------+--------------------------+----------------+--------------------+--------------------+-------------------+---------------+----------------------+--------------+--------------------+------------------------+-----------------------+-------------------+
| Replica_IO_State                 | Source_Host  | Source_User | Source_Port | Connect_Retry | Source_Log_File  | Read_Source_Log_Pos | Relay_Log_File                  | Relay_Log_Pos | Relay_Source_Log_File | Replica_IO_Running | Replica_SQL_Running | Replicate_Do_DB | Replicate_Ignore_DB | Replicate_Do_Table | Replicate_Ignore_Table | Replicate_Wild_Do_Table | Replicate_Wild_Ignore_Table | Last_Errno | Last_Error | Skip_Counter | Exec_Source_Log_Pos | Relay_Log_Space | Until_Condition | Until_Log_File | Until_Log_Pos | Source_SSL_Allowed | Source_SSL_CA_File | Source_SSL_CA_Path | Source_SSL_Cert | Source_SSL_Cipher | Source_SSL_Key | Seconds_Behind_Source | Source_SSL_Verify_Server_Cert | Last_IO_Errno | Last_IO_Error | Last_SQL_Errno | Last_SQL_Error | Replicate_Ignore_Server_Ids | Source_Server_Id | Source_UUID                          | Source_Info_File        | SQL_Delay | SQL_Remaining_Delay | Replica_SQL_Running_State                                | Source_Retry_Count | Source_Bind | Last_IO_Error_Timestamp | Last_SQL_Error_Timestamp | Source_SSL_Crl | Source_SSL_Crlpath | Retrieved_Gtid_Set | Executed_Gtid_Set | Auto_Position | Replicate_Rewrite_DB | Channel_Name | Source_TLS_Version | Source_public_key_path | Get_Source_public_key | Network_Namespace |
+----------------------------------+--------------+-------------+-------------+---------------+------------------+---------------------+---------------------------------+---------------+-----------------------+--------------------+---------------------+-----------------+---------------------+--------------------+------------------------+-------------------------+-----------------------------+------------+------------+--------------+---------------------+-----------------+-----------------+----------------+---------------+--------------------+--------------------+--------------------+-----------------+-------------------+----------------+-----------------------+-------------------------------+---------------+---------------+----------------+----------------+-----------------------------+------------------+--------------------------------------+-------------------------+-----------+---------------------+----------------------------------------------------------+--------------------+-------------+-------------------------+--------------------------+----------------+--------------------+--------------------+-------------------+---------------+----------------------+--------------+--------------------+------------------------+-----------------------+-------------------+
| Waiting for source to send event | 10.11.17.202 | repl        |       13306 |            60 | mysql-bin.000001 |                1047 | replicas-mysql-relay-bin.000002 |           326 | mysql-bin.000001      | Yes                | Yes                 |                 |                     |                    |                        |                         |                             |          0 |            |            0 |                1047 |             545 | None            |                |             0 | No                 |                    |                    |                 |                   |                |                     0 | No                            |             0 |               |              0 |                |                             |                1 | fdb57683-4c52-11ee-b932-0242ac110004 | mysql.slave_master_info |         0 |                NULL | Replica has read all relay log; waiting for more updates |              86400 |             |                         |                          |                |                    |                    |                   |             0 |                      |              |                    |                        |                     0 |                   |
+----------------------------------+--------------+-------------+-------------+---------------+------------------+---------------------+---------------------------------+---------------+-----------------------+--------------------+---------------------+-----------------+---------------------+--------------------+------------------------+-------------------------+-----------------------------+------------+------------+--------------+---------------------+-----------------+-----------------+----------------+---------------+--------------------+--------------------+--------------------+-----------------+-------------------+----------------+-----------------------+-------------------------------+---------------+---------------+----------------+----------------+-----------------------------+------------------+--------------------------------------+-------------------------+-----------+---------------------+----------------------------------------------------------+--------------------+-------------+-------------------------+--------------------------+----------------+--------------------+--------------------+-------------------+---------------+----------------------+--------------+--------------------+------------------------+-----------------------+-------------------+
1 row in set (0.00 sec)
```



# 主节点ktoy_mysql_1查看连接状态

检查连接到该主库的从库信息

```sql
SHOW SLAVE HOSTS;
+-----------+------+------+-----------+--------------------------------------+
| Server_id | Host | Port | Master_id | Slave_UUID                           |
+-----------+------+------+-----------+--------------------------------------+
|         2 |      | 3306 |         1 | 1266f2b9-4c53-11ee-b8b9-0242ac110005 |
+-----------+------+------+-----------+--------------------------------------+
1 row in set, 1 warning (0.00 sec)
```