# 中介者模式

## 意图

> * 中介者模式封装对象之间互交，使依赖变的简单，并且使复杂互交简单化，封装在中介者中。
> * 例子中的中介者使用单例模式生成中介者。
> * 中介者的change使用switch判断类型。

## 参考代码

### package

[senghoo_mediator](/media/senghoo_design_pattern/08_mediator/mediator.go ':include :type=code')

### test

[senghoo_mediator](/media/senghoo_design_pattern/08_mediator/mediator_test.go ':include :type=code')

