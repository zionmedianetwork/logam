package logam

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	// EnvDevelopment represents dev environment
	EnvDevelopment = "development"
	// EnvProduction is a live production environment
	EnvProduction = "production"
	// EnvStaging is QA or staging environment
	EnvStaging = "staging"
	// EnvWorkstation developper's workstation
	EnvWorkstation = "workstation"
)

type Logger interface {
	Errorf(format string, args ...interface{})
	Error(args ...interface{})
	Errorw(string, ...interface{})
	Fatalf(format string, args ...interface{})
	Fatal(args ...interface{})
	Infof(format string, args ...interface{})
	Info(args ...interface{})
	Infow(string, ...interface{})
	Warnf(format string, args ...interface{})
	Warn(args ...interface{})
	Warnw(string, ...interface{})
	Debugf(format string, args ...interface{})
	Debug(args ...interface{})
	Debugw(string, ...interface{})
	Printf(format string, args ...interface{})
	Print(args ...interface{})
	Tracef(string, ...interface{})
}

type dxLogger struct {
	cfg    Config
	logger *zap.SugaredLogger
}

type Config struct {
	LogLevel    string
	LogFormat   string
	Environment string
}

func NewLogger(cfg Config) Logger {
	logger := newLogger(cfg)
	logger.initLogger()

	return logger
}

func newLogger(cfg Config) *dxLogger {
	return &dxLogger{cfg: cfg}
}

var levelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

func (s *dxLogger) getLogLevel() zapcore.Level {
	level, ok := levelMap[strings.ToLower(s.cfg.LogLevel)]
	if !ok {
		level = zapcore.DebugLevel
	}

	return level
}

func (s *dxLogger) initLogger() {
	logLevel := s.getLogLevel()
	logWriter := zapcore.AddSync(os.Stderr)

	// Set encoderconfig
	var encoderCfg zapcore.EncoderConfig
	switch s.cfg.Environment {
	case EnvWorkstation, EnvDevelopment:
		encoderCfg = zap.NewDevelopmentEncoderConfig()
	default:
		encoderCfg = zap.NewProductionEncoderConfig()
	}

	encoderCfg.LevelKey = "LEVEL"
	encoderCfg.MessageKey = "MESSAGE"
	encoderCfg.CallerKey = "CALLER"
	encoderCfg.TimeKey = "TIME"
	encoderCfg.NameKey = "NAME"
	encoderCfg.EncodeTime = zapcore.RFC3339TimeEncoder

	// set encoder
	var encoder zapcore.Encoder
	if s.cfg.LogFormat == "console" {
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	}
	// instantiate core logger
	core := zapcore.NewCore(encoder, logWriter, zap.NewAtomicLevelAt(logLevel))
	logger := zap.New(core, zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel))

	s.logger = logger.Sugar()
	if err := s.logger.Sync(); err != nil {
		s.logger.Warn(err)
	}
}

func (s *dxLogger) Errorf(format string, args ...interface{}) {
	s.logger.Errorf(format, args...)
}

func (s *dxLogger) Error(args ...interface{}) {
	s.logger.Error(args...)
}

func (s *dxLogger) Fatalf(format string, args ...interface{}) {
	s.logger.Fatalf(format, args...)
}

func (s *dxLogger) Fatal(args ...interface{}) {
	s.logger.Fatal(args...)
}

func (s *dxLogger) Infof(format string, args ...interface{}) {
	s.logger.Infof(format, args...)
}

func (s *dxLogger) Info(args ...interface{}) {
	s.logger.Info(args...)
}

func (s *dxLogger) Warnf(format string, args ...interface{}) {
	s.logger.Warnf(format, args...)
}

func (s *dxLogger) Warn(args ...interface{}) {
	s.logger.Warn(args...)
}

func (s *dxLogger) Debugf(format string, args ...interface{}) {
	s.logger.Debugf(format, args...)
}

func (s *dxLogger) Debug(args ...interface{}) {
	s.logger.Debug(args...)
}

func (s *dxLogger) Printf(format string, args ...interface{}) {
	s.logger.Infof(format, args...)
}

func (s *dxLogger) Print(args ...interface{}) {
	s.logger.Info(args...)
}

func (s *dxLogger) Infow(message string, keyAndValues ...interface{}) {
	s.logger.Infow(message, keyAndValues...)
}

func (s *dxLogger) Errorw(message string, keyAndValues ...interface{}) {
	s.logger.Errorw(message, keyAndValues...)
}

func (s *dxLogger) Debugw(message string, keyAndValues ...interface{}) {
	s.logger.Debugw(message, keyAndValues...)
}

func (s *dxLogger) Warnw(message string, keyAndValues ...interface{}) {
	s.logger.Warnw(message, keyAndValues...)
}

func (s *dxLogger) Tracef(format string, args ...interface{}) {
	s.logger.Infof(format, args...)
}
