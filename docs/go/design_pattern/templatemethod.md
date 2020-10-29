# 模版方法模式

## 意图

> 模版方法模式使用继承机制，把通用步骤和通用方法放到父类中，把具体实现延迟到子类中实现。使得实现符合开闭原则。
> 因为Golang不提供继承机制，需要使用匿名组合模拟实现继承。

## 关键代码

> * 因为父类需要调用子类方法，所以子类需要匿名组合父类的同时，父类需要持有子类的引用。
> * 通用步骤在抽象类中实现，变化的步骤在具体的子类中实现

## 应用场景

> * 如实例代码中通用步骤在父类中实现（`准备`、`下载`、`保存`、`收尾`）下载和保存的具体实现留到子类中，并且提供 `保存`方法的默认实现
> * 做饭，打开煤气，开火，（做饭）， 关火，关闭煤气。除了做饭其他步骤都是相同的，抽到抽象类中实现

## 代码实现

[templatemethod](/media/templatemethod/templatemethod.go ':include :type=code')


## 参考代码

### package

[senghoo_templatemethod](/media/senghoo_design_pattern/14_template_method/templatemethod.go ':include :type=code')

### test

[senghoo_templatemethod](/media/senghoo_design_pattern/14_template_method/templatemethod_test.go ':include :type=code')

