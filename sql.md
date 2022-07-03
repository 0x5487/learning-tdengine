
## 每秒 K line

select sum(size) as size, first(price) as open, max(price) as high, min(price) as low, last(price) as close
from gate_trades_btc_usdt
interval(1s) sliding(1s)
order by ts desc
limit 10;


## k-line
create table k_line_1s_btc_usdt as select sum(size) as size, first(price) as open, max(price) as high, min(price) as low, last(price) as close
from gate_trades_eth_usdt interval(1s) sliding(1s);

