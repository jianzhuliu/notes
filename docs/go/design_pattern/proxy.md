# 代理模式

## 意图

> 代理模式用于延迟处理操作或者在进行实际操作前后进行其它处理

## 关键代码

> 代理类和被代理类实现同一接口，代理类中持有被代理类对象

## 应用实例

> * 火车票的代理售票点。代售点就是代理，它拥有被代理对象的部分功能 — 售票功能
> * 明星的经纪人，经纪人就是代理，负责为明星处理一些事务。
> * 虚代理
> * COW代理
> * 远程代理
> * 保护代理
> * Cache 代理
> * 防火墙代理
> * 同步代理
> * 智能指引

## 代码实现

[proxy](/media/proxy/proxy.go ':include :type=code') 

## 参考代码

### package

[senghoo_proxy](/media/senghoo_design_pattern/09_proxy/proxy.go ':include :type=code')

### test

[senghoo_proxy](/media/senghoo_design_pattern/09_proxy/proxy_test.go ':include :type=code')
