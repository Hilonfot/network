package catch

import (
	"fmt"
	"github.com/hilonfot/network/utils/log"
	"os"
	"runtime/debug"
	"time"
)

func PanicHandler() {
	exeName := os.Args[0] // 获取程序名称

	now := time.Now()  // 获取当前时间
	pid := os.Getpid() // 获取进程ID

	timeStr := now.Format("20060102150405")
	fName := fmt.Sprintf("%s-%d-%s-dump.log", exeName, pid, timeStr)
	log.Error("dump to file ", fName)

	f, err := os.Create(fName)
	if err != nil {
		return
	}
	defer f.Close()

	if err := recover(); err != nil {
		f.WriteString(fmt.Sprintf("%v\r\n", err)) // 输出panic信息
		f.WriteString("========\r\n")
	}

	f.WriteString(string(debug.Stack())) //输出堆栈信息
}
