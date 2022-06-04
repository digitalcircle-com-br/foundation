package dbmgr_test

import (
	"log"
	"testing"

	"github.com/digitalcircle-com-br/foundation/lib/dbmgr"
	"github.com/stretchr/testify/assert"
)

func TestSome(t *testing.T) {
	ret, err := dbmgr.DBN("auth")
	assert.NoError(t, err)
	log.Printf("%#v", ret)

}
