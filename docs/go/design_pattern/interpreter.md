# 解释器模式

## 意图

> * 解释器模式定义一套语言文法，并设计该语言解释器，使用户能使用特定文法控制解释器行为。
> * 解释器模式的意义在于，它分离多种复杂功能的实现，每个功能只需关注自身的解释。
> * 对于调用者不用关心内部的解释器的工作，只需要用简单的方式组合命令就可以。

## 参考代码

### package

[senghoo_interpreter](/media/senghoo_design_pattern/19_interpreter/interpreter.go ':include :type=code')

### test

[senghoo_interpreter](/media/senghoo_design_pattern/19_interpreter/interpreter_test.go ':include :type=code')

