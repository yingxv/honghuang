###
 # @Author: fuRan NgeKaworu@gmail.com
 # @Date: 2023-03-19 01:16:04
 # @LastEditors: fuRan NgeKaworu@gmail.com
 # @LastEditTime: 2023-11-03 20:55:07
 # @FilePath: /honghuang/docker/boot.sh
 # @Description: 
 # 
 # Copyright (c) 2023 by ${git_name_email}, All Rights Reserved. 
### 
docker compose pull
docker compose --env-file ~/.env down
docker compose --env-file ~/.env up -d
docker restart furan.xyz