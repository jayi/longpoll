package longpoll

import (
	"fmt"
	"time"
)

func ExampleGetData() {

	// 更新数据时调用UpdateData通知
	DefaultManager.UpdateData("tag_test", "example data", "1")

	// 获取数据调用GetData, 版本不是最新时，立即返回最新版本
	data, version, err := DefaultManager.GetData("tag_test", "0", time.Minute)
	fmt.Println(data, version, err)

	// 版本是最新时，阻塞等待数据更新或超时
	data, version, err = DefaultManager.GetData("tag_test", version, time.Second*1)
	fmt.Println(data, version, err)

	go func() {
		time.Sleep(3)
		DefaultManager.UpdateData("tag_test", "example data v2", "2")
	}()
	// 此处查询会等上面数据更新后返回
	data, version, err = DefaultManager.GetData("tag_test", version, time.Second*5)
	fmt.Println(data, version, err)

	// 查询无效tag会报错
	data, version, err = DefaultManager.GetData("tag_test2", version, time.Second*1)
	fmt.Println(data, version, err)
	// Output:
	// example data 1 <nil>
	// example data 1 <nil>
	// example data v2 2 <nil>
	// <nil>  tag_test2 is not exist
}
