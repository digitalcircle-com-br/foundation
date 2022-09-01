package crudxlsmgr

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var db *gorm.DB

func Setup(ndb *gorm.DB) {
	db = ndb
}

type fieldHeader struct {
	fname  string
	ftype  string
	format string
}

const (
	OP_INSERT   = "i"
	OP_UPDATE   = "u"
	OP_DELETE   = "d"
	TYPE_STRING = "string"
	TYPE_INT    = "int"
	TYPE_FLOAT  = "float"
	TYPE_DATE   = "date"
	TYPE_BOOL   = "bool"
	COL_OP      = 0
	COL_ID      = 1
)

func prepareHeader(cols []string) []fieldHeader {
	ret := make([]fieldHeader, 0)
	for i := 2; i < len(cols); i++ {
		parts := strings.Split(cols[i], ",")
		var fh = fieldHeader{
			fname: parts[0],
			ftype: TYPE_STRING,
		}
		if len(parts) > 1 {
			fh.ftype = parts[1]
		}
		if len(parts) > 2 {
			fh.format = parts[2]
		}
		ret = append(ret, fh)
	}
	return ret
}

func prepareUpsertVal(h []fieldHeader, row []string) (map[string]interface{}, error) {
	var ret = map[string]interface{}{}
	for i := 2; i < len(h)+2; i++ {
		oneHeader := h[i-2]
		switch oneHeader.ftype {
		case TYPE_STRING:
			ret[oneHeader.fname] = row[i]
			continue
		case TYPE_BOOL:
			b, err := strconv.ParseBool(row[i])
			if err != nil {
				return nil, fmt.Errorf("field: %s - %s", oneHeader.fname, err.Error())
			}
			ret[oneHeader.fname] = b
			continue
		case TYPE_FLOAT:
			b, err := strconv.ParseFloat(row[i], 64)
			if err != nil {
				return nil, fmt.Errorf("field: %s - %s", oneHeader.fname, err.Error())
			}
			ret[oneHeader.fname] = b
			continue
		case TYPE_INT:
			b, err := strconv.ParseInt(row[i], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("field: %s - %s", oneHeader.fname, err.Error())
			}
			ret[oneHeader.fname] = b
			continue
		case TYPE_DATE:
			b, err := time.Parse(oneHeader.format, row[i])
			if err != nil {
				return nil, fmt.Errorf("field: %s - %s", oneHeader.fname, err.Error())
			}
			ret[oneHeader.fname] = b
		}

	}
	return ret, nil
}

func processRow(h []fieldHeader, row []string, tb string) error {
	switch row[0] {
	case OP_INSERT:
		upsert, err := prepareUpsertVal(h, row)
		if err != nil {
			return err
		}
		return db.Table(tb).Create(upsert).Error
	case OP_UPDATE:
		upsert, err := prepareUpsertVal(h, row)
		if err != nil {
			return err
		}
		return db.Table(tb).Where("id = ?", row[COL_ID]).Updates(upsert).Error
	case OP_DELETE:
		return db.Table(tb).Where("id = ?", row[COL_ID]).Delete(nil).Error
	}

	return nil
}

func processXls(xls *excelize.File, tb string) error {
	sheetName := xls.GetSheetName(xls.GetActiveSheetIndex())
	rows, err := xls.Rows(sheetName)
	if err != nil {
		return err
	}
	defer rows.Close()
	counter := 0
	var headers []fieldHeader
	if rows.Next() {
		counter++
		cols, err := rows.Columns()
		if err != nil {
			return err
		}
		headers = prepareHeader(cols)
	}
	for rows.Next() {
		counter++
		cols, err := rows.Columns()
		if err != nil {
			return err
		}
		err = processRow(headers, cols, tb)
		if err != nil {
			return fmt.Errorf("error processing row %v - %s", counter, err.Error())
		}
	}
	return nil
}

func UploadExcel(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(2 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	for _, fs := range r.MultipartForm.File {
		for _, f := range fs {
			fdesc, err := f.Open()
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			xls, err := excelize.OpenReader(fdesc)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			tb := r.URL.Query().Get("table")
			err = processXls(xls, tb)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
	}
}
