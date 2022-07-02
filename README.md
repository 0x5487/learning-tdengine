# learning-tdengine

## Note

1. 因為 Taos 的驅動是 C 寫的，所以都要先安裝驅動 (不管你是用 jdbc or golang)

## Tools (db)

## SQL

```sql
create database demo keep 365 precision "us";
CREATE STABLE trades(ts timestamp, order_id NCHAR(40), price DOUBLE, size DOUBLE, side TINYINT) TAGS (market NCHAR(16), vendor NCHAR(16));
CREATE TABLE gate_trades_btc_usdt USING trades TAGS ("BTC/USDT", "gate");

create table btc_usdt(ts timestamp, order_id NCHAR(30), price DOUBLE, size DOUBLE, side TINYINT);
```

## 安装

1. 启动 Tdengine (docker)

   ```shell
   docker run -d --name tdengine --hostname="tdengine" -p 6030-6049:6030-6049 -p 6030-6049:6030-6049/udp tdengine/tdengine
   ```

2. 修改 host file, 加入你的 IP

   ```
   192.168.0.108 tdengine
   ```

3. 安装 windows drivers, 建立连线, 参考 [建立連接 | TDengine 文檔 | 濤思數據 (taosdata.com)](https://docs.taosdata.com/develop/connect/)

4.

## 写入资料

1. 對同一張表，如果新插入記錄的時間戳已經存在，默認情形下（UPDATE=0）新記錄將被直接拋棄，也就是說，在一張表裡，時間戳必須是唯一的
