# 观察者模式

## 意图

> - 定义对象间的一种一对多的依赖关系
> - 当一个对象的状态发生改变时，所有依赖于它的对象都得到通知并被自动更新
> - 而此对象无需关心连动对象的具体实现。

## 关键代码

> 被观察者持有了集合存放观察者 (收通知的为观察者)

## 应用场景

> - 报纸订阅，报社为被观察者，订阅的人为观察者
> - MVC 模式，当 model 改变时，View 视图会自动改变，model 为被观察者，View 为观察者
> - 创建订单之后，主业务完成，后续关联了弱关联业务,短信通知用户下单成功,通知物流系统

## 代码实现

[observer](/media/observer/observer.go ':include :type=code')

## 参考代码

### package

[senghoo_observer](/media/senghoo_design_pattern/10_observer/obserser.go ':include :type=code')

### test

[senghoo_observer](/media/senghoo_design_pattern/10_observer/obserser_test.go ':include :type=code')

