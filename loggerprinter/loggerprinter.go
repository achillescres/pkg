package loggerprinter

type Printer interface {
	Printf(format string, v ...any)
}

type Logger interface {
	Logf(format string, v ...any)
}

type LoggerPrinter interface {
	Logger
	Printer
}

type fromLogger struct {
	Logger
}

func (logger fromLogger) Printf(format string, v ...any) {
	logger.Logf(format, v...)
}

func FromLogger(logger Logger) LoggerPrinter {
	return fromLogger{logger}
}

type fromPrinter struct {
	Printer
}

func (printer fromPrinter) Logf(format string, v ...any) {
	printer.Printf(format, v...)
}

func FromPrinter(printer Printer) LoggerPrinter {
	return fromPrinter{printer}
}
