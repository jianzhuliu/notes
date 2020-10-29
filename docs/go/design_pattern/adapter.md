# 适配器模式 

## 意图

> 适配器模式将一个类的接口，转换成客户期望的另一个接口。适配器让原本接口不兼容的类可以合作无间

## 关键代码

适配器中持有旧接口对象，并实现新接口

## 应用实例

> 适配器适合用于解决新旧系统（或新旧接口）之间的兼容问题，而不建议在一开始就直接使用

## 代码实现

[adapter](/media/adapter/adapter.go ':include :type=code') 

## 参考代码

### package

[senghoo_adapter](/media/senghoo_design_pattern/02_adapter/adapter.go ':include :type=code')

### test

[senghoo_adapter](/media/senghoo_design_pattern/02_adapter/adapter_test.go ':include :type=code')
