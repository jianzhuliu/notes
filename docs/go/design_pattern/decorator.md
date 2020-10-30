# 装饰模式

## 意图

> * 装饰模式使用对象组合的方式动态改变或增加对象行为。
> * Go语言借助于匿名组合和非入侵式接口可以很方便实现装饰模式
> * 使用匿名组合，在装饰器中不必显式定义转调原对象方法。

## 关键代码

> 装饰器和被装饰对象实现同一个接口，装饰器中使用了被装饰对象

## 应用实例

> 各种软件框架的中间件

## 代码实现

[decorator](/media/decorator/decorator.go ':include :type=code')

## 参考代码

### package

[senghoo_decorator](/media/senghoo_design_pattern/20_decorator/decorator.go ':include :type=code')

### test

[senghoo_decorator](/media/senghoo_design_pattern/20_decorator/decorator_test.go ':include :type=code')

