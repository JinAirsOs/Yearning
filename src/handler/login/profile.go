package login

import (
	"Yearning-go/src/handler/common"
	"Yearning-go/src/model"
	"github.com/cookieY/yee"
	"net/http"
)

func UserReqSwitch(c yee.Context) (err error) {
	return c.JSON(http.StatusOK, common.SuccessPayload(map[string]interface{}{"reg": model.GloOther.Register}))
}
