# learning-tdengine

## Note

1. 因為 Taos 的驅動是 C 寫的，所以都要先安裝驅動 (不管你是用 jdbc or golang)
<https://www.taosdata.com/assets-download/3.0/TDengine-client-3.0.1.0-Linux-x64-Lite.tar.gz>

1. tdengine 預設密碼 root:taosdata
1. Ddweaver 工具: <https://www.taosdata.com/engineering/12880.html>

## Tools (db)

1. docker run -d --name tdengine --hostname="tdengine" -p 6030-6049:6030-6049 -p 6030-6049:6030-6049/udp tdengine/tdengine:3.0.0.0

2. 修改 host file, 加入你的 IP

   ```
   192.168.0.108 tdengine
   ```

3. 安装 windows drivers, 建立连线, 参考 [建立連接 | TDengine 文檔 | 濤思數據 (taosdata.com)](https://docs.taosdata.com/develop/connect/)

4. jdbc:TAOS-RS://localhost:6041

## 写入资料

1. 對同一張表，如果新插入記錄的時間戳已經存在，默認情形下（UPDATE=0）新記錄將被直接拋棄，也就是說，在一張表裡，時間戳必須是唯一的
