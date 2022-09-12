# learning-tdengine

## Note

1. 因為 Taos 的驅動是 C 寫的，所以都要先安裝驅動 (不管你是用 jdbc or golang)
<https://www.taosdata.com/assets-download/3.0/TDengine-client-3.0.1.0-Linux-x64-Lite.tar.gz>
<https://www.taosdata.com/assets-download/3.0/TDengine-client-3.0.1.0-Windows-x64.exe>

1. tdengine 預設密碼 root:taosdata
1. Ddweaver 工具: <https://www.taosdata.com/engineering/12880.html>

## DBeaver

1. 先安裝 Tdengine windows client driver
<https://www.taosdata.com/assets-download/3.0/TDengine-client-3.0.1.0-Windows-x64.exe>

1. 對於 Windows 上的 JDBC, ODBC, Python, Go 等連接，確保C:\TDengine\driver\taos.dll在你的系統庫函數搜尋目錄裡 (建議taos.dll放在目錄 C:\Windows\System32)
1. Antifact

<dependency>
 <groupId>com.taosdata.jdbc</groupId>
 <artifactId>taos-jdbcdriver</artifactId>
 <version>3.0.0</version>
</dependency>

## Tools (db)

1. docker run -d --name tdengine --hostname="tdengine" -p 6030:6030 -p 6041:6041 -p 6043-6049:6043-6049 -p 6043-6049:6043-6049/udp tdengine/tdengine:3.0.1.0

2. 修改 host file, 加入你的 IP

   ```
   192.168.0.108 tdengine
   ```

3. 安装 windows drivers, 建立连线, 参考 [建立連接 | TDengine 文檔 | 濤思數據 (taosdata.com)](https://docs.taosdata.com/develop/connect/)

4. jdbc:TAOS-RS://localhost:6041
jdbc:TAOS://localhost:6030

## 写入资料

1. 對同一張表，如果新插入記錄的時間戳已經存在，默認情形下（UPDATE=0）新記錄將被直接拋棄，也就是說，在一張表裡，時間戳必須是唯一的
