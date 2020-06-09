#### [https://trade.yym.plus/ ](https://trade.yym.plus/ "") 

测试账户：test  密码：testtest



#### 机器人demo:

```javascript
var config = host.config()

var exchangeList = config.exchanges.map(exId => exchange.get(exId))

function onInit()
{

}

function onTick()
{
    var tickerTasks = exchangeList.map(ex => ex.tickerEx("BTC/USDT"))

    var tickers = utils.wait(tickerTasks, "5s")
    
    for (var i = 0; i < tickers.length; i++)
    {
        var ticker = tickers[i]
        var ex = exchangeList[i]
        chart.save("tickers",
        {
            price: ticker.last
        },
        {
            ex: ex.name(),
            coin: "BTC/USDT"
        })
    }
    
    utils.sleep("3s")
}

function onStop()
{
    console.log("onStop")
}
```
