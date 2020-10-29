# 简单工厂模式

## 意图

> - 定义一个创建对象的接口，通过传递参数来决定由哪个工厂类来实现
> - go 语言没有构造函数一说，所以一般会定义NewXXX函数来初始化相关类
> - NewXXX 函数返回接口时就是简单工厂模式。

## 关键代码

> 返回一个接口 

## 应用实例

> 造车，可以是 A 厂制造，也可以是 B 厂制造，不需要关心是如何制造出来的

## 代码实现

[simplefactory](/media/factory/simplefactory.go ':include :type=code')

## 参考代码

### package

[senghoo_simplefactory](/media/senghoo_design_pattern/00_simple_factory/simple.go ':include :type=code')

### test

[senghoo_simplefactory](/media/senghoo_design_pattern/00_simple_factory/simple_test.go ':include :type=code')
