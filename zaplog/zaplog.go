package zaplog

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path"
	"path/filepath"
	"time"
)

var (
	sp = string(filepath.Separator)
)

type Options struct {
	LogFileDir  string        // 文件保存地方
	AppName     string        // 日志文件前缀
	Level       zapcore.Level // 日志等级
	Development bool          // 是否开发模式
	MaxSize     int           // 日志文件小大（M）
	MaxBackups  int           // 最多存在多少个切片文件
	MaxAge      int           // 保存的最大天数
}

type OptionsFunc func(options *Options)

type Logger struct {
	*zap.Logger

	zapConfig zap.Config

	opts *Options
}

func NewLogger(opts ...OptionsFunc) *Logger {
	l := &Logger{}
	l.defaultOptions()

	for _, fn := range opts {
		fn(l.opts)
	}

	l.setConfig()
	var err error
	l.Logger, err = l.zapConfig.Build(l.getCores())
	if err != nil {
		panic(err)
		return nil
	}
	defer l.Logger.Sync()
	return l
}

func (l *Logger) defaultOptions() {
	dir, _ := filepath.Abs(filepath.Dir(filepath.Join(".")))
	dir = path.Join(dir, "logs")
	l.opts = &Options{
		LogFileDir: dir,
		AppName:    "app",
		MaxSize:    100,
		MaxAge:     30,
		MaxBackups: 60,
		Level:      zap.DebugLevel,
	}
}

func (l *Logger) setConfig() {
	l.zapConfig = zap.NewProductionConfig()
	l.zapConfig.Level.SetLevel(l.opts.Level)
	l.zapConfig.OutputPaths = []string{"stdout"}
	l.zapConfig.ErrorOutputPaths = []string{"stderr"}
	l.zapConfig.EncoderConfig.EncodeTime = timeEncoder
	l.zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
}

func (l *Logger) getCores() zap.Option {
	fileEncoder := zapcore.NewJSONEncoder(l.zapConfig.EncoderConfig)
	normalFunc := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level >= l.zapConfig.Level.Level()
	})
	errorFunc := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level >= zapcore.WarnLevel
	})

	f := func(fn string) zapcore.WriteSyncer {
		return l.getWriter(fn)
	}

	cores := []zapcore.Core{
		zapcore.NewCore(fileEncoder, f(".log"), normalFunc),
		zapcore.NewCore(fileEncoder, f(".err.log"), errorFunc),
	}
	if l.opts.Development {
		encoderConfig := zap.NewDevelopmentEncoderConfig()
		encoderConfig.EncodeTime = timeEncoder
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
		cores = append(cores, []zapcore.Core{
			zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), normalFunc),
			zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stderr), errorFunc),
		}...)
	}
	return zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return zapcore.NewTee(cores...)
	})
}

func (l *Logger) getWriter(name string) zapcore.WriteSyncer {
	fileName := l.opts.LogFileDir + sp + l.opts.AppName + name
	hook := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    l.opts.MaxSize,
		MaxBackups: l.opts.MaxBackups,
		MaxAge:     l.opts.MaxAge,
		Compress:   true,
	}
	return zapcore.AddSync(hook)
}

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

func SetLogFileDir(dir string) OptionsFunc {
	return func(opt *Options) {
		opt.LogFileDir = dir
	}
}

func SetAppName(name string) OptionsFunc {
	return func(opt *Options) {
		opt.AppName = name
	}
}

func SetDevelopment(dev bool) OptionsFunc {
	return func(opt *Options) {
		opt.Development = dev
	}
}

func SetLevel(lv string) OptionsFunc {
	return func(opt *Options) {
		switch lv {
		case "debug":
			opt.Level = zap.DebugLevel
		case "info":
			opt.Level = zap.InfoLevel
		case "warn":
			opt.Level = zap.WarnLevel
		case "error":
			opt.Level = zap.ErrorLevel
		}
	}
}
