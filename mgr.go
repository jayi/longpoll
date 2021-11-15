/*
Package longpoll 长轮询服务端实现

*/
package longpoll

import (
	"errors"
	"time"
)

// Manager 长轮询管理
type Manager struct {
	pollMap map[string]PollData
	timeout time.Duration
}

// SetTimeout 设置所传tag的超时时间。不传tag时，设置默认超时时间，只对新tag生效。
func (lpm *Manager) SetTimeout(timeout time.Duration, tag ...string) *Manager {
	if timeout < time.Second {
		return lpm
	}

	if len(tag) > 0 {
		for _, t := range tag {
			pd, ok := lpm.pollMap[t]
			if ok {
				pd.SetTimeout(timeout)
			}
		}
	} else {
		lpm.timeout = timeout
	}

	return lpm
}

// UpdateData 接收所传tag的数据更新
func (lpm *Manager) UpdateData(tag string, data interface{}, version ...string) {
	pd, ok := lpm.pollMap[tag]
	if !ok {
		pd = NewPollData().SetTimeout(lpm.timeout)
		lpm.pollMap[tag] = pd
	}
	pd.UpdateData(data, version...)
}

// GetData 根据所传tag和version获取最新数据及版本
// 版本不是最新时，立即返回最新版本
// 版本是最新时，阻塞等待数据更新或超时
func (lpm *Manager) GetData(tag string, version string, wait ...time.Duration) (interface{}, string, error) {
	pd, ok := lpm.pollMap[tag]
	if !ok {
		return nil, "", errors.New(tag + " is not exist")
	}

	data, version := pd.GetData(version, wait...)
	return data, version, nil
}

// GetPoolData 获取tag对应的数据管理实例
func (lpm *Manager) GetPoolData(tag string) (PollData, error) {
	pd, ok := lpm.pollMap[tag]
	if !ok {
		return nil, errors.New(tag + " is not exist")
	}

	return pd, nil
}

// NewManager 新建长轮询管理实例
func NewManager() *Manager {
	return &Manager{
		timeout: DefalutTimeout,
		pollMap: make(map[string]PollData),
	}
}

// DefaultManager 默认的长轮询管理实例
var DefaultManager *Manager = NewManager()
