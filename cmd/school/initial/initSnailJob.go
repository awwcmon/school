package initial

import (
	"school/internal/config"
	"time"

	"github.com/go-dev-frame/sponge/pkg/logger"
	snailjob "github.com/open-snail/snail-job-go"
	"github.com/open-snail/snail-job-go/dto"
	"github.com/open-snail/snail-job-go/job"
	"github.com/sirupsen/logrus"
)

func InitSnailJob() {
	// 配置Options参数
	snail := config.Get().Snailjob
	exec := snailjob.NewSnailJobManager(&dto.Options{
		ServerHost:   snail.ServerHost,
		ServerPort:   snail.ServerPort,
		HostIP:       snail.HostIP,
		HostPort:     snail.HostPort,
		Namespace:    snail.Namespace,
		GroupName:    snail.GroupName,
		Token:        snail.Token,
		Level:        logrus.InfoLevel,
		ReportCaller: true,
	})
	logger.Info("snail is regersting")
	// 注册执行器
	exec.Register("testJobExecutor", func() job.IJobExecutor {
		return &TestJobExecutor{}
	})

	// 初始化环境
	if nil == exec.Init() {
		// 启动客户端
		exec.Run()
		logger.Info("启动snailjob客户端成功")
	}
	logger.Info("hahahahaha")
}

type TestJobExecutor struct {
	job.BaseJobExecutor
}

func (executor *TestJobExecutor) DoJobExecute(jobArgs dto.IJobArgs) dto.ExecuteResult {
	executor.RemoteLogger.Infof("TestJobExecutor 执行结束 DoJobExecute. jobId: [%d] now:[%s]", jobArgs.GetJobId(), time.Now().String())
	return *dto.Success().WithMessage("hello 这是go客户端")
}
