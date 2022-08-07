package files_test

// func Test_service_Upload(t *testing.T) {
// 	testmgr.Init()
// 	sessid := testmgr.Get(t, "session")
// 	// tenant, err := redis.Get("test:tenant")
// 	// assert.NoError(t, err)
// 	err := files.Service.Setup()
// 	assert.NoError(t, err)

// 	buf := &bytes.Buffer{}
// 	mpw := multipart.NewWriter(buf)
// 	fw, _ := mpw.CreateFormFile("f1", "f1.txt")
// 	fw.Write([]byte("f1 content"))
// 	fw, _ = mpw.CreateFormFile("f1", "f2.txt")
// 	fw.Write([]byte("f2 content"))
// 	mpw.Close()

// 	w := testmgr.NewInMemResponseWriter()

// 	r, _ := http.NewRequest(http.MethodPost, "/upload", buf)

// 	r.Header = http.Header{}
// 	r.Header.Set("Cookie", fmt.Sprintf("%s=%s", fmodel.COOKIE_SESSION, sessid))
// 	r.Header.Set("Content-Type", mpw.FormDataContentType())
// 	nctx := context.WithValue(r.Context(), fmodel.CTX_REQ, r)
// 	nctx = context.WithValue(nctx, fmodel.CTX_RES, w)

// 	r = r.WithContext(nctx)
// 	files.Service.Upload(w, r)
// 	resp := make([]files.UploadResponseEntry, 0)
// 	err = json.NewDecoder(w).Decode(&resp)
// 	assert.NoError(t, err)
// 	assert.Equal(t, len(resp), 2)
// }

// func Test_service_Download(t *testing.T) {
// 	testmgr.Init()
// 	sessid := testmgr.Get(t, "session")

// 	err := files.Service.Setup()
// 	assert.NoError(t, err)

// 	w := testmgr.NewInMemResponseWriter()

// 	r, _ := http.NewRequest(http.MethodGet, "/download?f=1", nil)

// 	r.Header = http.Header{}
// 	r.Header.Set("Cookie", fmt.Sprintf("%s=%s", fmodel.COOKIE_SESSION, sessid))

// 	nctx := context.WithValue(r.Context(), fmodel.CTX_REQ, r)
// 	nctx = context.WithValue(nctx, fmodel.CTX_RES, w)

// 	r = r.WithContext(nctx)
// 	files.Service.Download(w, r)
// 	assert.Equal(t, w.Status(), http.StatusOK)
// }

// func Test_service_List(t *testing.T) {
// 	testmgr.Init()
// 	sessid := testmgr.Get(t, "session")
// 	err := files.Service.Setup()
// 	assert.NoError(t, err)

// 	w := testmgr.NewInMemResponseWriter()

// 	r, _ := http.NewRequest(http.MethodGet, "/list?name=%f%", nil)

// 	r.Header = http.Header{}
// 	r.Header.Set("Cookie", fmt.Sprintf("%s=%s", fmodel.COOKIE_SESSION, sessid))

// 	nctx := context.WithValue(r.Context(), fmodel.CTX_REQ, r)
// 	nctx = context.WithValue(nctx, fmodel.CTX_RES, w)
// 	r = r.WithContext(nctx)
// 	files.Service.List(w, r)
// 	assert.Equal(t, w.Status(), http.StatusOK)
// 	ret := make([]fmodel.File, 0)
// 	err = json.NewDecoder(w).Decode(&ret)
// 	assert.NoError(t, err)
// 	assert.Greater(t, len(ret), 0)
// }

// func Test_service_List_DtLimit(t *testing.T) {
// 	testmgr.Init()
// 	sessid := testmgr.Get(t, "session")

// 	err := files.Service.Setup()
// 	assert.NoError(t, err)

// 	w := testmgr.NewInMemResponseWriter()

// 	r, _ := http.NewRequest(http.MethodGet, "/list?name=%f%&dtini=2022-01-02", nil)

// 	r.Header = http.Header{}
// 	r.Header.Set("Cookie", fmt.Sprintf("%s=%s", fmodel.COOKIE_SESSION, sessid))

// 	nctx := context.WithValue(r.Context(), fmodel.CTX_REQ, r)
// 	nctx = context.WithValue(nctx, fmodel.CTX_RES, w)
// 	r = r.WithContext(nctx)
// 	files.Service.List(w, r)
// 	assert.Equal(t, w.Status(), http.StatusOK)
// 	ret := make([]fmodel.File, 0)
// 	err = json.NewDecoder(w).Decode(&ret)
// 	assert.NoError(t, err)
// 	assert.Greater(t, len(ret), 0)
// }
