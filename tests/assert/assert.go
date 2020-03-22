package assert

import (
	"fmt"
	"strings"
	"testing"

	"github.com/go-pg/pg/v9"
	ta "github.com/stretchr/testify/assert"

	"vobook/database"
	"vobook/services/fs"
	"vobook/utils"
)

type NotNilT bool

var NotNil = NotNilT(true)

func DatabaseCount(t *testing.T, table string, data utils.M) int {
	args := []interface{}{}
	wheres := []string{}

	args = append(args, pg.Ident(table))

	for col, val := range data {
		if val == nil {
			wheres = append(wheres, fmt.Sprintf("%s IS NULL", col))
			continue
		}
		switch val.(type) {
		case NotNilT:
			wheres = append(wheres, fmt.Sprintf("%s IS NOT NULL", col))
			continue
		}
		args = append(args, val)
		wheres = append(wheres, fmt.Sprintf(`"%s"=?`, col))
	}

	whereExpr := strings.Join(wheres, " AND ")
	query := fmt.Sprintf("SELECT COUNT(*) FROM ? WHERE %s", whereExpr)
	rows := 0
	_, err := database.ORM().Query(pg.Scan(&rows), query, args...)
	if err != nil {
		t.Fatal(err)
	}

	return rows
}

func DatabaseHas(t *testing.T, table string, data utils.M) {
	count := DatabaseCount(t, table, data)
	if count == 0 {
		t.Fatal(fmt.Sprintf("Table %s missing row with data %+v", table, data))
	}
}

func DatabaseHasDeleted(t *testing.T, table string, ids ...string) {
	for _, id := range ids {
		data := utils.M{
			"id":         id,
			"deleted_at": NotNil,
		}
		count := DatabaseCount(t, table, data)
		if count == 0 {
			t.Fatal(fmt.Sprintf("Table %s missing row with data %+v", table, data))
		}
	}
}

func DatabaseMissing(t *testing.T, table string, data utils.M) {
	count := DatabaseCount(t, table, data)
	if count != 0 {
		t.Fatal(fmt.Sprintf("Table %s has row with data %+v", table, data))
	}
}

func NotError(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func Equals(t *testing.T, expected, actual interface{}, msgAndArgs ...interface{}) {
	if !ta.Equal(t, expected, actual, msgAndArgs...) {
		t.FailNow()
	}
}

func True(t *testing.T, value bool, msgAndArgs ...interface{}) {
	if !ta.True(t, value, msgAndArgs...) {
		t.FailNow()
	}
}

func FileExists(t *testing.T, path string, msgAndArgs ...interface{}) {
	path = fs.FullPath(path)
	if !ta.FileExists(t, path, msgAndArgs...) {
		t.FailNow()
	}
}

func FileNotExists(t *testing.T, path string, msgAndArgs ...interface{}) {
	path = fs.FullPath(path)
	if !ta.NoFileExists(t, path, msgAndArgs...) {
		t.FailNow()
	}
}
