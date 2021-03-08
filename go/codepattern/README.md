
# go语言编码模式

> - 接口编程

[interface](/interface.go ':include :type=code :fragment=interface')

[interface_cmd](/cmd/interface.go ':include :type=code :fragment=interface_cmd')

> - 错误处理，参考 bufio.Scanner

[error](/error.go ':include :type=code :fragment=error')

[error_cmd](/cmd/error.go ':include :type=code :fragment=error_cmd')

> - 函数编程

[functional_options](/functional_options.go ':include :type=code :fragment=functional_options')

[functional_options_cmd](/cmd/functional_options.go ':include :type=code :fragment=functional_options_cmd')

# 性能提示

> - 如果需要把数字转换成字符串，使用 strconv.Itoa() 比 fmt.Sprintf() 要快

> - 尽可能避免把String转成[]Byte ，这个转换会导致性能下降

> - 如果在 for-loop 里对某个 Slice 使用 append()，请先把 Slice 的容量扩充到位，这样可以避免内存重新分配但又用不到的情况，从而避免浪费内存

> - 使用 bytes.Buffer 或是 strings.Builder 来拼接字符串，性能会比使用 + 或 +=高

> - 尽可能使用并发的 goroutine，然后使用 sync.WaitGroup 来同步分片操作

> - 避免在热代码中进行内存分配，这样会导致 gc 很忙。尽可能使用 sync.Pool 来重用对象

> - 使用 lock-free 的操作，避免使用 sync.Mutex，尽可能使用 sync/atomic 包

> - 使用 I/O 缓冲，I/O 是个非常非常慢的操作，使用 bufio.NewWriter() 和 bufio.NewReader() 可以带来更高的性能

> - 对于在 for-loop 里的固定的正则表达式，一定要使用 regexp.Compile() 编译正则表达式。性能会提升两个数量级

> - 如果你需要更高性能的协议，就要考虑使用 protobuf 或 msgp 而不是 JSON，因为 JSON 的序列化和反序列化里使用了反射

> - 你在使用 Map 的时候，使用整型的 key 会比字符串的要快，因为整型比较比字符串比较要快

