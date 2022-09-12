## SQL

```sql
create database demo keep 365 precision "us";
CREATE STABLE trades(ts timestamp, order_id NCHAR(40), price DOUBLE, size DOUBLE, vol DOUBLE, side TINYINT) TAGS (market NCHAR(16), vendor NCHAR(16));
CREATE TABLE trades_btc_usdt_gate USING trades TAGS ("BTC/USDT", "gate");

```

## 每秒 K line

select _wstart, first(price) as open, max(price) as high, min(price) as low, last(price) as close, sum(size) as base_vol, sum(vol) as quote_vol, first(price) - last(price) as delta, (first(price) / last(price) - 1) % 100 as percentage
from trades
interval(5m)
order by _wstart desc
limit 10;

## k-line (1s)

create table k_line_1s_btc_usdt as select sum(size) as size, first(price) as open, max(price) as high, min(price) as low, last(price) as close
from gate_trades_btc_usdt
interval(1s) sliding(1s);

## ticker
