# 备忘录模式

## 意图

> * 备忘录模式用于保存程序内部状态到外部，又不希望暴露内部状态的情形。
> * 程序内部状态使用窄接口船体给外部进行存储，从而不暴露程序实现细节。
> * 备忘录模式同时可以离线保存内部状态，如保存到数据库，文件等。

## 参考代码

### package

[senghoo_memento](/media/senghoo_design_pattern/17_memento/memento.go ':include :type=code')

### test

[senghoo_memento](/media/senghoo_design_pattern/17_memento/memento_test.go ':include :type=code')

