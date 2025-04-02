package logger

// import (
// 	"log/slog"
// 	"os"
// )

// type Attr = slog.Attr

// type Logger interface {
// 	Error(message string, attrs ...Attr)
// 	Info(message string, attrs ...Attr)
// 	Warn(message string, attrs ...Attr)
// 	Debug(message string, attrs ...Attr)
// 	Fatal(message string, attrs ...Attr)
// 	With(attrs ...Attr) Logger
// }

// type SlogLogger struct {
// 	logger *slog.Logger
// }

// func NewSlogLogger(level slog.Level) *SlogLogger {
// 	opts := &slog.HandlerOptions{
// 		Level: level,
// 	}
// 	handler := slog.NewJSONHandler(os.Stdout, opts)
// 	return &SlogLogger{
// 		logger: slog.New(handler),
// 	}
// }

// func (l *SlogLogger) Error(msg string, attrs ...Attr) {
// 	l.logger.Error(msg, convertAttrs(attrs)...)
// }

// func (l *SlogLogger) Info(msg string, attrs ...Attr) {
// 	l.logger.Info(msg, convertAttrs(attrs)...)
// }

// func (l *SlogLogger) Warn(msg string, attrs ...Attr) {
// 	l.logger.Warn(msg, convertAttrs(attrs)...)
// }

// func (l *SlogLogger) Debug(msg string, attrs ...Attr) {
//     l.logger.Debug(msg, convertAttrs(attrs)...)
// }

// func (l *SlogLogger) Fatal(msg string, attrs ...Attr) {
//     l.logger.Error(msg, convertAttrs(attrs)...)
//     os.Exit(1)
// }

// func (l *SlogLogger) With(attrs ...Attr) Logger {
//     return &SlogLogger{
//         logger: l.logger.With(convertAttrs(attrs)...),
//     }
// }

// func convertAttrs(attrs []Attr) []any {
//     anyAttrs := make([]any, len(attrs))
//     for i, attr := range attrs {
//         anyAttrs[i] = attr
//     }
//     return anyAttrs
// }