# longpoll
--
    import "test/longpoll"

Package longpoll 长轮询服务端实现

## Usage

```go
const DefalutTimeout = time.Minute
```
DefalutTimeout 长轮询服务端默认超时时间

#### type Manager

```go
type Manager struct {
}
```

Manager 长轮询管理

```go
var DefaultManager *Manager = NewManager()
```
DefaultManager 默认的长轮询管理实例

#### func  NewManager

```go
func NewManager() *Manager
```
NewManager 新建长轮询管理实例

#### func (*Manager) GetData

```go
func (lpm *Manager) GetData(tag string, version string, wait ...time.Duration) (interface{}, string, error)
```
GetData 根据所传tag和version获取最新数据及版本

#### func (*Manager) GetPoolData

```go
func (lpm *Manager) GetPoolData(tag string) (PollData, error)
```
GetPoolData 获取tag对应的数据管理实例

#### func (*Manager) SetTimeout

```go
func (lpm *Manager) SetTimeout(timeout time.Duration, tag ...string) *Manager
```
SetTimeout 设置所传tag的超时时间。不传tag时，设置默认超时时间，只对新tag生效。

#### func (*Manager) UpdateData

```go
func (lpm *Manager) UpdateData(tag string, data interface{}, version ...string)
```
UpdateData 接收所传tag的数据更新

#### type PollData

```go
type PollData interface {

	// UpdateData 更新数据(data)和版本号（version)
	UpdateData(data interface{}, version ...string)

	// GetData 携带当前version获取最新数据和版本号，客户端等待时间为wait。
	GetData(version string, wait ...time.Duration) (interface{}, string)

	// SetTimeout 设置服务端默认等待时间为timeout
	SetTimeout(timeout time.Duration) PollData
}
```

PollData 长轮询实现接口

#### func  NewPollData

```go
func NewPollData() PollData
```
NewPollData 生成一个长轮询管理实例
