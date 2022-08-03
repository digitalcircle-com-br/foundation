package auth_test

//
// func Test_service_Login(t *testing.T) {
// 	core.Init("auth_test")
// 	err := dbmgr.DropRecreate("auth")
// 	assert.NoError(t, err)
// 	err = auth.Service.Setup()
// 	assert.NoError(t, err)
// 	res, err := auth.Service.Login(context.Background(), &auth.AuthRequest{Login: "root", Password: "root"})
// 	assert.NoError(t, err)
// 	assert.NotNil(t, res)
// 	// testmgr.Set(t, "session", res.SessionID)
// 	// testmgr.Set(t, "tenant", res.Tenant)
// }
// func Test_service_Check(t *testing.T) {
// 	testmgr.Init()
// 	sessid := testmgr.Get(t, "session")

// 	assert.NotEmpty(t, sessid)
// 	sess, err := sessionmgr.SessionLoad(sessid)
// 	assert.NoError(t, err)
// 	ctx := context.WithValue(context.Background(), model.COOKIE_SESSION, sess)
// 	ctx = context.WithValue(ctx, model.CTX_SESSION, sess)
// 	res, err := auth.Service.Check(ctx, &model.EMPTY{})
// 	assert.NoError(t, err)
// 	assert.True(t, res)
// }
// func Test_service_Logout(t *testing.T) {
// 	core.Init("auth_test")
// 	sessid := testmgr.Get(t, "session")
// 	assert.NotEmpty(t, sessid)
// 	sess, err := sessionmgr.SessionLoad(sessid)
// 	assert.NoError(t, err)
// 	ctx := context.WithValue(context.Background(), model.COOKIE_SESSION, sess)
// 	ctx = context.WithValue(ctx, model.CTX_SESSION, sess)
// 	res, err := auth.Service.Logout(ctx, &model.EMPTY{})
// 	assert.NoError(t, err)
// 	assert.True(t, res)
// }
