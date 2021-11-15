package longpoll

import (
	"time"
)

// DefalutTimeout 长轮询服务端默认超时时间
const DefalutTimeout = time.Minute

// PollData 长轮询实现接口
type PollData interface {

	// UpdateData 更新数据(data)和版本号（version)
	UpdateData(data interface{}, version ...string)

	// GetData 携带当前version获取最新数据和版本号，客户端等待时间为wait。
	GetData(version string, wait ...time.Duration) (interface{}, string)

	// SetTimeout 设置服务端默认等待时间为timeout
	SetTimeout(timeout time.Duration) PollData
}

type dataDetail struct {
	data    interface{}
	version string
}

type pollDataImp struct {
	dataDetail *dataDetail
	ch         chan struct{}
	timeout    time.Duration
}

// UpdateData 更新数据(data)和版本号（version)
// 如果没带版本号，自动填充时间，格式为RFC3339-"2006-01-02T15:04:05Z07:00"
func (pd *pollDataImp) UpdateData(data interface{}, version ...string) {
	var ver string
	if len(version) == 0 {
		ver = time.Now().Format(time.RFC3339)
	} else {
		ver = version[0]
	}
	pd.dataDetail.version = ver
	pd.dataDetail.data = data

	close(pd.ch) // 广播更新信号
	pd.ch = make(chan struct{})
}

// GetData 携带当前version获取数据，客户端等待时间为wait
// 版本不是最新时，立即返回最新版本
// 版本是最新时，阻塞等待数据更新或超时
func (pd *pollDataImp) GetData(version string, wait ...time.Duration) (interface{}, string) {
	if version != pd.dataDetail.version {
		return pd.dataDetail.data, pd.dataDetail.version
	}

	var timeout *time.Timer
	if len(wait) > 0 && wait[0] >= time.Second && wait[0] < pd.timeout {
		// 客户端自定义的超时时间
		timeout = time.NewTimer(wait[0])
	} else {
		// 服务端默认超时时间
		timeout = time.NewTimer(pd.timeout)
	}

	select {
	case <-timeout.C: // 超时
		break
	case <-pd.ch: // 收到更新信号
		break
	}
	return pd.dataDetail.data, pd.dataDetail.version
}

// SetTimeout 设置服务端默认等待时间为timeout
func (pd *pollDataImp) SetTimeout(timeout time.Duration) PollData {
	if timeout < time.Second {
		return pd
	}

	pd.timeout = timeout
	return pd
}

// NewPollData 生成一个长轮询管理实例
func NewPollData() PollData {
	return &pollDataImp{
		dataDetail: &dataDetail{},
		timeout:    DefalutTimeout,
		ch:         make(chan struct{}),
	}
}
