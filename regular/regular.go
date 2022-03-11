package regular

import (
	"bufio"
	"fmt"
	"math/bits"
	"strconv"
	"strings"

	xt "github.com/free-utils-go/xorm_type_assist"
	"github.com/godcong/excavator"
	"github.com/godcong/fate"
	"xorm.io/xorm"
)

var RegularList = map[int]int{} //key:unicode,value:stroke

func init() {
	file_regular, err := excavator.DataFiles.Open("data/regular.txt")
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(file_regular)

	for scanner.Scan() {
		line := scanner.Text()
		regular_pair := strings.Split(line, ":")
		regular_stroke, err := strconv.ParseInt(strings.TrimSpace(regular_pair[0]), 10, bits.UintSize)
		if err != nil {
			panic(err)
		}

		regular_arr_stage1 := strings.ReplaceAll(regular_pair[1], ",", "")

		regular_arr_stage2 := strings.ReplaceAll(regular_arr_stage1, "`", "")

		for _, regular_char := range strings.Split(strings.TrimSpace(regular_arr_stage2), "、") {
			RegularList[int(([]rune(regular_char))[0])] = int(regular_stroke)
		}

	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

type regular struct {
	exc     *excavator.Excavator
	total   int
	fixed   int
	unfixed int
}

// Run ...
func (r *regular) Run() {

	for ch_uni := range RegularList {
		ch := string(rune(ch_uni))

		r.total++
		fmt.Printf("character %s is fixing\n", ch)
		if fixRegular(r.exc.DbFate, ch) {
			r.fixed++
		} else {
			r.unfixed++
		}
	}
	fmt.Printf("fix regular finished(total:%d,fixed:%d,unfixed:%d)\n", r.total, r.fixed, r.unfixed)
}

// Regular ...
type Regular interface {
	Run()
}

// New ...
func New(excavator *excavator.Excavator) Regular {
	return &regular{
		exc: excavator,
	}
}

func fixRegular(db *xorm.Engine, ch string) bool {
	my_char := fate.Character{
		Ch: ch,
	}
	has, err := db.Get(&my_char)
	if err != nil {
		fmt.Printf("failed get char(%s) with error (%v)\n", ch, err)
		return false
	} else {
		if !has {
			fmt.Printf("%s not found\n", ch)
			return false
		}
		if my_char.IsRegular == xt.TRUE {
			return false
		}
	}

	my_char.IsRegular = xt.TRUE
	e := excavator.InsertOrUpdate(db, &my_char)

	if e != nil {
		fmt.Printf("failed update char(%s) with error (%v)\n", ch, e)
		return false
	}

	return true
}
