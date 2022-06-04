package f8n

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type CrudOpts struct {
	Name        string
	Ops         string
	DefaultPerm string
	CPerm       string
	RPerm       string
	UPerm       string
	DPerm       string
}

func (c *CrudOpts) PermFor(op string) string {
	switch op {
	case "C":
		if c.CPerm != "" {
			return c.CPerm
		}
	case "R":
		if c.RPerm != "" {
			return c.RPerm
		}
	case "U":
		if c.UPerm != "" {
			return c.UPerm
		}
	case "D":
		if c.DPerm != "" {
			return c.DPerm
		}
	default:
	}
	if c.DefaultPerm != "" {
		return c.DefaultPerm
	}

	return c.Name
}

func Crud[T any](opts *CrudOpts) {
	if opts.Ops == "" {
		opts.Ops = "C,R,U,D"
	}
	ops := strings.Split(strings.ToUpper(opts.Ops), ",")

	for _, op := range ops {

		switch op {

		case "R":
			HttpHandle(fmt.Sprintf("/crud/%s", opts.Name), http.MethodGet, PermDef(opts.PermFor("R")), func(ctx context.Context, t *EMPTY_TYPE) (ret *[]T, err error) {
				db, err := CtxDB(ctx)
				if err != nil {
					return
				}
				ret = new([]T)

				req := CtxReq(ctx)
				rawWhere := req.URL.Query().Get("where")
				skip := req.URL.Query().Get("skip")
				limit := req.URL.Query().Get("limit")
				order := req.URL.Query().Get("order")
				tx := db
				if rawWhere != "" {
					parts := strings.Split(rawWhere, ";")
					if len(parts) == 1 {
						tx = tx.Where(parts[0])
					} else {
						is := make([]interface{}, len(parts)-1)

						for _, v := range parts[1:] {
							is = append(is, v)
						}

						tx = tx.Where(parts[0], is...)
					}

				}
				if order != "" {
					tx = tx.Order(order)
				}
				if skip != "" {
					skipi, err := strconv.ParseInt(skip, 10, 32)
					if err == nil {
						tx = tx.Offset(int(skipi))
					}
				}

				if limit != "" {
					limiti, err := strconv.ParseInt(skip, 10, 32)
					if err == nil {
						tx = tx.Limit(int(limiti))
					}
				}
				err = tx.Find(ret).Error
				return
			})

			HttpHandle(fmt.Sprintf("/crud/%s/{id}", opts.Name), http.MethodGet, PermDef(opts.PermFor("R")), func(ctx context.Context, t *EMPTY_TYPE) (ret *T, err error) {
				db, err := CtxDB(ctx)
				id := CtxVars(ctx)["id"]
				ret = new(T)
				if err != nil {
					return
				}
				err = db.Where("id = ?", id).First(ret).Error
				return
			})

		case "C":
			HttpHandle(fmt.Sprintf("/crud/%s", opts.Name), http.MethodPost, PermDef(opts.PermFor("C")), func(ctx context.Context, t *T) (out *T, err error) {
				db, err := CtxDB(ctx)
				if err != nil {
					return
				}
				err = db.Create(t).Error
				out = t
				return
			})
		case "U":
			HttpHandle(fmt.Sprintf("/crud/%s/{id}", opts.Name), http.MethodPut, PermDef(opts.PermFor("U")), func(ctx context.Context, t *T) (out *T, err error) {
				db, err := CtxDB(ctx)
				id := CtxVars(ctx)["id"]
				if err != nil {
					return
				}
				err = db.Where("id = ?", id).Updates(t).Error
				if err != nil {
					return
				}
				err = db.Where("id =?", id).First(t).Error
				out = t
				return
			})

		case "D":
			HttpHandle(fmt.Sprintf("/crud/%s/{id}", opts.Name), http.MethodDelete, PermDef(opts.PermFor("D")), func(ctx context.Context, t *EMPTY_TYPE) (out *EMPTY_TYPE, err error) {
				db, err := CtxDB(ctx)
				id := CtxVars(ctx)["id"]
				if err != nil {
					return
				}
				in := new(T)
				err = db.Where("id = ?", id).Delete(in).Error
				return
			})
		}
	}

}
