# 策略模式

## 意图

> 定义一系列算法，让这些算法在运行时可以互换，使得分离算法，符合开闭原则。

## 关键代码

> 实现同一个接口

## 应用场景

> * 支付方式的使用
> * 主题的更换，每个主题都是一种策略
> * 旅行的出游方式，选择骑自行车、坐汽车，每一种旅行方式都是一个策略

## 代码实现

[strategy](/media/strategy/strategy.go ':include :type=code')

## 参考代码

### package

[senghoo_strategy](/media/senghoo_design_pattern/15_strategy/strategy.go ':include :type=code')

### test

[senghoo_strategy](/media/senghoo_design_pattern/15_strategy/strategy_test.go ':include :type=code')

