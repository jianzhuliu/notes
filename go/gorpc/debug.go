/*
支持 http 协议的好处在于， rpc 服务仅仅使用了监听端口的 /gorpc路径
在其他路径上，可以提供诸如日志，统计等更为丰富的功能
*/
package gorpc

import (
	"fmt"
	"html/template"
	"net/http"
)

const debugTmpl = `<html>
<header><title>gorpc services</title></header>
<body>
{{ range .}}
<hr>
Service {{.Name}}
<hr>
	<table>
	<th align="center">Method</th>
	<th align="center">Calls</th>
	
	{{range $name, $mtype := .Method}}
		<tr>
		<td align="left" font="fixed">{{$name}}({{$mtype.ArgType}}, {{$mtype.ReplyType}}) error</td>
		<td align="center">{{$mtype.NumCalls}}</td>
		</tr>
	{{end}}
	</table>
{{end}}
</body></html>`

var debug = template.Must(template.New("gorpc").Parse(debugTmpl))

type debugHTTP struct {
	*Server
}

type debugService struct {
	Name   string
	Method map[string]*methodType
}

func (d debugHTTP) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	var services []debugService
	d.serviceMap.Range(func(namei, svci interface{}) bool {
		svc := svci.(*service)
		services = append(services, debugService{
			Name:   namei.(string),
			Method: svc.method,
		})

		return true
	})

	err := debug.Execute(rw, services)
	if err != nil {
		_, _ = fmt.Fprintln(rw, "rpc:error executing template:", err.Error())
	}
}
