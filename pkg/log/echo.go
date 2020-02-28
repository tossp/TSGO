package log

//var (
//	verServer = fmt.Sprintf("PHP/%s", reverseString(strings.Replace(runtime.Version(), "go", "", -1)))
//)

//func reverseString(s string) string {
//	runes := []rune(s)
//
//	for from, to := 0, len(runes)-1; from < to; from, to = from+1, to-1 {
//		runes[from], runes[to] = runes[to], runes[from]
//	}
//
//	return string(runes)
//}
//
//func Echo() echo.MiddlewareFunc {
//	return func(next echo.HandlerFunc) echo.HandlerFunc {
//		return func(c echo.Context) (err error) {
//			req := c.Request()
//			res := c.Response()
//			start := time.Now()
//			c.Response().Header().Set(echo.HeaderServer, verServer)
//
//			// Request
//			reqBody := []byte("")
//			if c.Request().Body != nil { // Read
//				reqBody, _ = ioutil.ReadAll(c.Request().Body)
//			}
//			c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(reqBody)) // Reset
//
//			if err = next(c); err != nil {
//				c.Error(err)
//			}
//			//_ = next(c)
//			fields, entry, n := makelog(c, reqBody, req, res, start)
//
//			switch {
//			case n >= 500:
//				entry.Fatalf("HTTP(%s) %s %s:[%s]%s\n", fields["latencyHuman"], fields["remote"], fields["status"], req.Method, fields["path"])
//			case n >= 400:
//				entry.Errorf("HTTP(%s) %s %s:[%s]%s\n", fields["latencyHuman"], fields["remote"], fields["status"], req.Method, fields["path"])
//			//case n >= 300:
//			//	entry.Infof("HTTP(%s) %s %s:[%s]%s\n", fields["latencyHuman"], fields["remote"], fields["status"], req.Method, fields["path"])
//			default:
//				entry.Infof("HTTP(%s) %s %s:[%s]%s\n", fields["latencyHuman"], fields["remote"], fields["status"], req.Method, fields["path"])
//			}
//			return
//		}
//	}
//}
//
//func makelog(c echo.Context, reqBody []byte, req *http.Request, res *echo.Response, start time.Time) (map[string]interface{}, *logrus.Entry, int) {
//	l := time.Now().Sub(start)
//
//	id := req.Header.Get(echo.HeaderXRequestID)
//	if id == "" {
//		id = res.Header().Get(echo.HeaderXRequestID)
//	}
//	p := req.URL.Path
//	if p == "" {
//		p = "/"
//	}
//	cl := req.Header.Get(echo.HeaderContentLength)
//	if cl == "" {
//		cl = "0"
//	}
//
//	fields := make(map[string]interface{})
//	fields["id"] = id
//	fields["remote"] = c.RealIP()
//	fields["host"] = req.Host
//	fields["uri"] = req.RequestURI
//	fields["method"] = req.Method
//	fields["path"] = p
//	fields["referer"] = req.Referer()
//	fields["userAgent"] = req.UserAgent()
//	fields["status"] = strconv.FormatInt(int64(res.Status), 10)
//	fields["latency"] = strconv.FormatInt(int64(l), 10)
//	fields["latencyHuman"] = l.String()
//	fields["bytesIn"] = cl
//	fields["bytesOut"] = strconv.FormatInt(res.Size, 10)
//	fields["header"] = c.Request().Header
//	fields["query"] = c.QueryParams()
//	//fields["form"], fields["formErr"] = c.FormParams()
//	fields["cookie"] = c.Cookies()
//	bodylen := len(reqBody)
//	fields["bodySize"] = bodylen
//	limitbody := 1 << 30
//	if bodylen != 0 && bodylen < limitbody {
//		fields["body"] = string(reqBody)
//	}
//	//="{\r\n  UserName:zh,\r\n  PassWord:\r\n}"
//	return fields, log.WithFields(fields), res.Status
//}
//
//func makelog2(c echo.Context, req *http.Request, res *echo.Response, start time.Time) (map[string]interface{}, *logrus.Entry, int) {
//	l := time.Now().Sub(start)
//
//	id := req.Header.Get(echo.HeaderXRequestID)
//	if id == "" {
//		id = res.Header().Get(echo.HeaderXRequestID)
//	}
//	p := req.URL.Path
//	if p == "" {
//		p = "/"
//	}
//	cl := req.Header.Get(echo.HeaderContentLength)
//	if cl == "" {
//		cl = "0"
//	}
//
//	fields := make(map[string]interface{})
//	fields["id"] = id
//	fields["remote"] = c.RealIP()
//	fields["host"] = req.Host
//	fields["uri"] = req.RequestURI
//	fields["method"] = req.Method
//	fields["path"] = p
//	fields["referer"] = req.Referer()
//	fields["userAgent"] = req.UserAgent()
//	fields["status"] = strconv.FormatInt(int64(res.Status), 10)
//	fields["latency"] = strconv.FormatInt(int64(l), 10)
//	fields["latencyHuman"] = l.String()
//	fields["bytesIn"] = cl
//	fields["bytesOut"] = strconv.FormatInt(res.Size, 10)
//	fields["header"] = c.Request().Header
//	fields["query"] = c.QueryParams()
//	//fields["form"], fields["formErr"] = c.FormParams()
//	fields["cookie"] = c.Cookies()
//
//	if bd, err := c.Request().GetBody(); err != nil {
//		log.Warn("获取Body失败", err)
//	} else {
//		if body, err := ioutil.ReadAll(bd); err != nil {
//			log.Warn("读取Body失败", err)
//		} else {
//			j := make(map[string]interface{})
//			if err = json.Unmarshal(body, &j); err != nil {
//				fields["body"] = j
//			} else {
//				fields["body"] = string(body)
//			}
//
//		}
//	}
//
//	return fields, log.WithFields(fields), res.Status
//}
