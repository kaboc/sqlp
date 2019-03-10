package placeholder

import (
	"reflect"
	"testing"
)

func TestReplace(t *testing.T) {
	src := `SELECT * FROM /*tbl1 --aaa
				*/ ` + "`tbl2`" + ` --bbb
				WHERE col1 = 'c\'cc'# /* ddd */
					OR col2 = 'eee
eee' /*
					OR col3 = '#fff'*/
					OR col4 = 'g/*g*/g' # :hhh ii/*ii*/ii
					OR col5 = 'iii'/*jjj*/#kkk
					OR col6 = ''
					OR col7 = :col7 -- :lll
					OR col8 = :col8# :mmm
					OR col9 = :col9/* :nnn */
				ORDER BY id--ooo
				LIMIT 10#'ppp'`

	dest := "SELECT * FROM /**SQLP_REPLACE**/ /**SQLP_REPLACE**//**SQLP_REPLACE**/" +
		"\t\t\t\tWHERE col1 = /**SQLP_REPLACE**//**SQLP_REPLACE**/" +
		"\t\t\t\t\tOR col2 = /**SQLP_REPLACE**/ /**SQLP_REPLACE**/\n" +
		"\t\t\t\t\tOR col4 = /**SQLP_REPLACE**/ /**SQLP_REPLACE**/" +
		"\t\t\t\t\tOR col5 = /**SQLP_REPLACE**//**SQLP_REPLACE**//**SQLP_REPLACE**/" +
		"\t\t\t\t\tOR col6 = /**SQLP_REPLACE**/\n" +
		"\t\t\t\t\tOR col7 = :col7/**SQLP_REPLACE**/" +
		"\t\t\t\t\tOR col8 = :col8/**SQLP_REPLACE**/" +
		"\t\t\t\t\tOR col9 = :col9/**SQLP_REPLACE**/\n" +
		"\t\t\t\tORDER BY id--ooo\n" +
		"\t\t\t\tLIMIT 10/**SQLP_REPLACE**/"

	q, r := replace(src)
	if q != dest {
		t.Fatalf("got: %#v\nwant: %#v", q, dest)
	}

	q = restore(q, r)
	if q != src {
		t.Fatalf("got: %#v\nwant: %#v", q, src)
	}
}

func TestIsNamed(t *testing.T) {
	{
		src := []interface{}{1, 2, 3}

		if isNamed(src) {
			t.Fatal("unnamed placeholder was recognized as named")
		}
	}

	{
		src := map[string]interface{}{
			"id": []interface{}{1, 2, 3},
		}

		if !isNamed(src) {
			t.Fatal("named placeholder was recognized as unnamed")
		}
	}
}

func TestSimplifyUnnamed(t *testing.T) {
	// Whether IN is all in capital letters does not matter
	{
		srcQ := "SELECT * FROM user WHERE id iN ?[3] AND age > ? AND name LIKE ? LIMIT ?# ?"
		srcV1 := []interface{}{1, 2, 3}
		srcV2 := 10
		srcV3 := []interface{}{"J%", 100}

		destQ := "SELECT * FROM user WHERE id iN (?,?,?) AND age > ? AND name LIKE ? LIMIT ?# ?"
		destV := []interface{}{1, 2, 3, 10, "J%", 100}

		q, b, err := Convert(srcQ, srcV1, srcV2, srcV3)
		if err != nil {
			t.Fatal(err)
		}

		if q != destQ {
			t.Fatalf("got: %s\nwant: %s", q, destQ)
		} else if !reflect.DeepEqual(b, destV) {
			t.Fatalf("got: %v\nwant: %v", b, destV)
		}
	}

	// Having mismatched numbers of placeholders and binding values gives no error
	{
		srcQ := "SELECT * FROM user WHERE id IN ?[3]"
		srcV := []interface{}{1, 2, 3, 4}

		destQ := "SELECT * FROM user WHERE id IN (?,?,?)"
		destV := []interface{}{1, 2, 3, 4}

		q, b, err := Convert(srcQ, srcV)
		if err != nil {
			t.Fatal(err)
		}

		if q != destQ {
			t.Fatalf("got: %s\nwant: %s", q, destQ)
		} else if !reflect.DeepEqual(b, destV) {
			t.Fatalf("got: %v\nwant: %v", b, destV)
		}
	}

	// 'IN ?[N]' inside quotes or comments are ignored
	{
		srcQ := "SELECT * FROM user WHERE id IN ?[2] AND name = 'id IN ?[2] ' AND age IN ?[2]"
		srcV := []interface{}{1, 2}

		destQ := "SELECT * FROM user WHERE id IN (?,?) AND name = 'id IN ?[2] ' AND age IN (?,?)"
		destV := []interface{}{1, 2}

		q, b, err := Convert(srcQ, srcV)
		if err != nil {
			t.Fatal(err)
		}

		if q != destQ {
			t.Fatalf("got: %s\nwant: %s", q, destQ)
		} else if !reflect.DeepEqual(b, destV) {
			t.Fatalf("got: %v\nwant: %v", b, destV)
		}
	}
}

func TestSimplifyNamed(t *testing.T) {
	// Whether IN is all in capital letters does not matter
	{
		srcQ := "SELECT * FROM user WHERE id iN :id[3] AND age > :age AND name LIKE :name LIMIT :limit -- :dummy"
		srcV := map[string]interface{}{
			"id":    []interface{}{1, 2, 3},
			"age":   10,
			"name":  "J%",
			"limit": 100,
		}

		destQ := "SELECT * FROM user WHERE id iN (?,?,?) AND age > ? AND name LIKE ? LIMIT ? -- :dummy"
		destV := []interface{}{1, 2, 3, 10, "J%", 100}

		q, b, err := Convert(srcQ, srcV)
		if err != nil {
			t.Fatal(err)
		}

		if q != destQ {
			t.Fatalf("got: %s\nwant: %s", q, destQ)
		} else if !reflect.DeepEqual(b, destV) {
			t.Fatalf("got: %v\nwant: %v", b, destV)
		}
	}

	// 'IN :placeholder[N]' inside quotes or comments are ignored
	{
		srcQ := "SELECT * FROM user WHERE id IN :id[2] AND name = 'id IN :id[2] ' AND age IN :age[2]/* :dummy */"
		srcV := map[string]interface{}{
			"id":  []interface{}{1, 2},
			"age": []interface{}{21, 22},
		}

		destQ := "SELECT * FROM user WHERE id IN (?,?) AND name = 'id IN :id[2] ' AND age IN (?,?)/* :dummy */"
		destV := []interface{}{1, 2, 21, 22}

		q, b, err := Convert(srcQ, srcV)
		if err != nil {
			t.Fatal(err)
		}

		if q != destQ {
			t.Fatalf("got: %s\nwant: %s", q, destQ)
		} else if !reflect.DeepEqual(b, destV) {
			t.Fatalf("got: %v\nwant: %v", b, destV)
		}
	}
}

func TestSimplifyNamedErr(t *testing.T) {
	{
		srcQ := "SELECT * FROM user WHERE id iN :id[1]"
		srcV := map[string]interface{}{"id": 1}

		if _, _, err := convertNamed(srcQ, srcV); err == nil {
			t.Fatal("error must be given if a slice is not used for IN statement")
		}
	}

	{
		srcQ := "SELECT * FROM user WHERE id = :id"
		srcV := map[string]interface{}{"id": []interface{}{1}}

		if _, _, err := convertNamed(srcQ, srcV); err == nil {
			t.Fatal("error must be given if a slice is used not for IN statement")
		}
	}

	{
		srcQ := "SELECT * FROM user WHERE id = :id"
		srcV := map[string]interface{}{
			"id":  1,
			"age": 10,
		}

		if _, _, err := convertNamed(srcQ, srcV); err == nil {
			t.Fatal("error must be given if binding values are set for unknown named placeholders")
		}
	}

	{
		srcQ := "SELECT * FROM user WHERE id = :id AND age > :age"
		srcV := map[string]interface{}{"id": 1}

		if _, _, err := convertNamed(srcQ, srcV); err == nil {
			t.Fatal("error must be given if binding values are insufficient")
		}
	}

	{
		srcQ := "SELECT * FROM user WHERE id IN :id[4]"
		srcV := map[string]interface{}{
			"id": []interface{}{1, 2, 3},
		}

		if _, _, err := convertNamed(srcQ, srcV); err == nil {
			t.Fatal("error must be given if a slice does not have N number of elements for :placeholder[N]")
		}
	}
}

func TestConvertSql(t *testing.T) {
	{
		src := "SELECT * FROM user WHERE id iN ?[3] -- :dummy"
		dest := "SELECT * FROM user WHERE id iN (?,?,?) -- :dummy"

		q, err := ConvertSQL(src)
		if err != nil {
			t.Fatal(err)
		}

		if q != dest {
			t.Fatalf("got: %s\nwant: %s", q, dest)
		}
	}

	{
		src := "SELECT * FROM user WHERE id iN :id[3] -- :dummy"
		dest := "SELECT * FROM user WHERE id iN (?,?,?) -- :dummy"

		q, err := ConvertSQL(src)
		if err != nil {
			t.Fatal(err)
		}

		if q != dest {
			t.Fatalf("got: %s\nwant: %s", q, dest)
		}
	}
}
