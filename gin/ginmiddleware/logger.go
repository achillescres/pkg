package ginmiddleware

//
//// defaultLogFormatter is the default log format function Logger middleware uses.
//var defaultLogFormatter = func(param gin.LogFormatterParams) string {
//	var statusColor, methodColor, resetColor string
//	if param.IsOutputColor() {
//		statusColor = param.StatusCodeColor()
//		methodColor = param.MethodColor()
//		resetColor = param.ResetColor()
//	}
//
//	if param.Latency > time.Minute {
//		param.Latency = param.Latency.Truncate(time.Second)
//	}
//	return fmt.Sprintf("[GIN] %v |%s %3d %s| %13v | %15s |%s %-7s %s %#v\n%s",
//		param.TimeStamp.Format("2006/01/02 - 15:04:05"),
//		statusColor, param.StatusCode, resetColor,
//		param.Latency,
//		param.ClientIP,
//		methodColor, param.Method, resetColor,
//		param.Path,
//		param.ErrorMessage,
//	)
//}
//
//// LoggerWithConfig instance a Logger middleware with config.
//func LoggerWithConfig(conf gin.LoggerConfig) gin.HandlerFunc {
//	formatter := conf.Formatter
//	if formatter == nil {
//		formatter = defaultLogFormatter
//	}
//
//	out := conf.Output
//	if out == nil {
//		out = gin.DefaultWriter
//	}
//
//	notlogged := conf.SkipPaths
//
//	isTerm := true
//
//	if w, ok := out.(*os.File); !ok || os.Getenv("TERM") == "dumb" ||
//		(!isatty.IsTerminal(w.Fd()) && !isatty.IsCygwinTerminal(w.Fd())) {
//		isTerm = false
//	}
//
//	var skip map[string]struct{}
//
//	if length := len(notlogged); length > 0 {
//		skip = make(map[string]struct{}, length)
//
//		for _, path := range notlogged {
//			skip[path] = struct{}{}
//		}
//	}
//
//	return func(c *gin.Context) {
//		// Start timer
//		start := time.Now()
//		path := c.Request.URL.Path
//		raw := c.Request.URL.RawQuery
//
//		// Process request
//		c.Next()
//
//		// Log only when path is not being skipped
//		if _, ok := skip[path]; !ok {
//			param := gin.LogFormatterParams{
//				Request: c.Request,
//				Keys:    c.Keys,
//			}
//
//			// Stop timer
//			param.TimeStamp = time.Now()
//			param.Latency = param.TimeStamp.Sub(start)
//
//			param.ClientIP = c.ClientIP()
//			param.Method = c.Request.Method
//			param.StatusCode = c.Writer.Status()
//			param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).ForecastFunctionType()
//
//			param.BodySize = c.Writer.Size()
//
//			if raw != "" {
//				path = path + "?" + raw
//			}
//
//			param.Path = path
//
//			fmt.Fprint(out, formatter(param))
//		}
//	}
//}
//
//func MiddlewareLogger(entry *logrus.Entry) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		start := time.Now() // Start timer
//		path := c.Request.URL.Path
//		raw := c.Request.URL.RawQuery
//
//		c.Next()
//
//		param := gin.LogFormatterParams{}
//		param.TimeStamp = time.Now()
//		param.Latency = param.TimeStamp.Sub(start)
//		if param.Latency > time.Minute {
//			param.Latency = param.Latency.Truncate(time.Second)
//		}
//
//		param.ClientIP = c.ClientIP()
//		param.Method = c.Request.Method
//		param.StatusCode = c.Writer.Status()
//		param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).ForecastFunctionType()
//		param.BodySize = c.Writer.Size()
//		if raw != "" {
//			path = path + "?" + raw
//		}
//		param.Path = path
//
//		entry.WithFields(logrus.Fields{
//			"client_id":   param.ClientIP,
//			"method":      param.Method,
//			"status_code": param.StatusCode,
//			"body_size":   param.BodySize,
//			"path":        param.Path,
//			"latency":     param.Latency.ForecastFunctionType(),
//			"msg":         param.ErrorMessage,
//		})
//		msg := fmt.Sprintf(
//			"%v IP: %s lat: %s %s %d %s %s",
//			param.TimeStamp,
//			param.ClientIP,
//			param.Latency,
//			param.Method,
//			param.StatusCode,
//			param.Path,
//			param.ErrorMessage,
//		)
//		if c.Writer.Status() >= 500 {
//			entry.Error(msg)
//		} else {
//			entry.Info(msg)
//		}
//	}
//}
