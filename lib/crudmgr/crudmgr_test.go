package crudmgr_test

import (
	"time"

	"github.com/digitalcircle-com-br/foundation/lib/model"
)

type Astr struct {
	model.BaseVO
	Atxt string
	Adt  time.Time
}

// func Test_Create(t *testing.T) {
// 	db, err := dbmgr.DBN("root")
// 	assert.NoError(t, err)
// 	err = db.AutoMigrate(&Astr{})
// 	assert.NoError(t, err)
// 	crudmgr.Handle(&Astr{})

// 	opts := &crudmgr.CrudOpts{Op: crudmgr.OP_C}
// 	opts.Data = &Astr{Atxt: "some", Adt: time.Now()}
// 	res := testmgr.NewInMemResponseWriter()
// 	req := testmgr.HttpNewAuthRequestO(t, http.MethodPost, "/crud/astr", opts, res)
// 	routemgr.Router().ServeHTTP(res, req)
// 	assert.True(t, res.Status() <= 399)

// }

// func Test_Retrieve(t *testing.T) {
// 	db, err := dbmgr.DBN("root")
// 	assert.NoError(t, err)
// 	err = db.AutoMigrate(&Astr{})
// 	assert.NoError(t, err)
// 	crudmgr.Handle(&Astr{})

// 	opts := &crudmgr.CrudOpts{Op: crudmgr.OP_R, Tb: "astrs"}
// 	res := testmgr.NewInMemResponseWriter()
// 	req := testmgr.HttpNewAuthRequestO(t, http.MethodPost, "/crud/astr", opts, res)
// 	routemgr.Router().ServeHTTP(res, req)
// 	assert.True(t, res.Status() <= 399)

// }

// func Test_Update(t *testing.T) {
// 	db, err := dbmgr.DBN("root")
// 	assert.NoError(t, err)
// 	err = db.AutoMigrate(&Astr{})
// 	assert.NoError(t, err)
// 	crudmgr.Handle(&Astr{})

// 	opts := &crudmgr.CrudOpts{Op: crudmgr.OP_U, Tb: "astrs", ID: 1}
// 	opts.Data = &Astr{Atxt: "somev2", Adt: time.Now()}
// 	res := testmgr.NewInMemResponseWriter()
// 	req := testmgr.HttpNewAuthRequestO(t, http.MethodPost, "/crud/astr", opts, res)
// 	routemgr.Router().ServeHTTP(res, req)
// 	assert.True(t, res.Status() <= 399)

// }

// func Test_Delete(t *testing.T) {
// 	db, err := dbmgr.DBN("root")
// 	assert.NoError(t, err)
// 	err = db.AutoMigrate(&Astr{})
// 	assert.NoError(t, err)
// 	crudmgr.Handle(&Astr{})

// 	opts := &crudmgr.CrudOpts{Op: crudmgr.OP_D, Tb: "astrs", ID: 1}
// 	opts.Data = &Astr{Atxt: "some", Adt: time.Now()}
// 	res := testmgr.NewInMemResponseWriter()
// 	req := testmgr.HttpNewAuthRequestO(t, http.MethodPost, "/crud/astr", opts, res)
// 	routemgr.Router().ServeHTTP(res, req)
// 	assert.True(t, res.Status() <= 399)

// }

// func Test_Create(t *testing.T) {
// 	db, err := dbmgr.DBN("test")
// 	assert.NoError(t, err)
// 	err = db.AutoMigrate(&Astr{})
// 	assert.NoError(t, err)
// 	ret, err := crudmgr.Create(&crudmgr.CrudOpts{
// 		Db:   "test",
// 		Tb:   "astrs",
// 		Data: &Astr{Atxt: "zxc", Adt: time.Now()},
// 	})
// 	assert.NoError(t, err)
// 	assert.NotEmpty(t, ret)
// }

// func Test_Update(t *testing.T) {
// 	db, err := dbmgr.DBN("test")
// 	assert.NoError(t, err)
// 	err = db.AutoMigrate(&Astr{})
// 	assert.NoError(t, err)
// 	ret, err := crudmgr.Update(&crudmgr.CrudOpts{
// 		ID:   1,
// 		Db:   "test",
// 		Tb:   "astrs",
// 		Data: &Astr{Atxt: "zxcv2", Adt: time.Now()},
// 	})
// 	assert.NoError(t, err)
// 	assert.Nil(t, ret)
// }
