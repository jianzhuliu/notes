# 外观模式 

## 意图

> 提供了一个统一的接口，用来访问子系统中的一群接口

## 关键代码

> 外观层中依次调用子系统的接口

## 应用实例

> - 电脑开机时，点击开机按钮，但同时启动了 CPU，内存，硬盘等
> - 后台开发，初始化文件，比如目录，数据库连接，日志接口等

## 代码实现

[facade](/media/facade/facade.go ':include :type=code')

## 参考代码

### package

[senghoo_facade](/media/senghoo_design_pattern/01_facade/facade.go ':include :type=code')

### test

[senghoo_facade](/media/senghoo_design_pattern/01_facade/facade_test.go ':include :type=code')
