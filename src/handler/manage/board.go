package manage

import (
	"Yearning-go/src/handler/common"
	"Yearning-go/src/lib"
	"Yearning-go/src/model"
	"github.com/cookieY/yee"
	"net/http"
	"strings"
)

type board struct {
	Board string `json:"board"`
}

const BOARD_MESSAGE_SAVE = "公告已保存"

func GeneralPostBoard(c yee.Context) (err error) {
	req := new(board)
	if err = c.Bind(req); err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusOK, err.Error())
	}
	model.DB().Model(model.CoreGlobalConfiguration{}).Where("1=1").Updates(&model.CoreGlobalConfiguration{Board: req.Board})
	return c.JSON(http.StatusOK, common.SuccessPayLoadToMessage(BOARD_MESSAGE_SAVE))
}

func GeneralGetBoard(c yee.Context) (err error) {
	var board model.CoreGlobalConfiguration
	model.DB().Select("board").First(&board)
	return c.JSON(http.StatusOK, common.SuccessPayload(board.Board))
}

type Result struct {
	Sql string
}

func GetQueryRecords(c yee.Context) (err error) {
	t := new(lib.Token).JwtParse(c)

	var result []Result
	model.DB().Raw("select * from core_query_orders a left join core_query_records b on a.work_id=b.work_id where a.username= ? and b.sql is not NULL order by a.id desc limit 1000", t.Username).Scan(&result)
	m := make(map[string]bool)
	var s []string
	for _, v := range result {
		trimedSql := trim(v.Sql)
		if m[trimedSql] {
			continue
		}
		m[trimedSql] = true
		s = append(s, trimedSql)
	}

	return c.JSON(http.StatusOK, common.SuccessPayload(s))
}

func trim(s string) string {
	return strings.Trim(s, "\n")
}
