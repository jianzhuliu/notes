# 命令模式

## 意图

>  命令模式本质是把某个对象的方法调用封装到对象中，方便传递、存储、调用。

示例中把主板单中的启动(start)方法和重启(reboot)方法封装为命令对象，再传递到主机(box)对象中。于两个按钮进行绑定：

* 第一个机箱(box1)设置按钮1(button1) 为开机按钮2(button2)为重启。
* 第二个机箱(box1)设置按钮2(button2) 为开机按钮1(button1)为重启。

从而得到配置灵活性。

## 应用实例

> * 批处理
> * 任务队列
> * undo, redo

## 参考代码

### package

[senghoo_command](/media/senghoo_design_pattern/11_command/command.go ':include :type=code')

### test

[senghoo_command](/media/senghoo_design_pattern/11_command/command_test.go ':include :type=code')
