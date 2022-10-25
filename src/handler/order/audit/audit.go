package audit

import (
	"Yearning-go/src/handler/common"
	"Yearning-go/src/lib"
	"Yearning-go/src/model"
	"github.com/cookieY/yee"
	"github.com/cookieY/yee/logger"
	"golang.org/x/net/websocket"
	"net/http"
	"time"
)

const QueryField = "work_id, username, text, backup, date, real_name, `status`, `type`, `delay`, `source`, `source_id`,`id_c`,`data_base`,`table`,`execute_time`,assigned,current_step,relevant"

func AuditOrderState(c yee.Context) (err error) {
	u := new(Confirm)
	user := new(lib.Token).JwtParse(c)
	if err = c.Bind(u); err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusOK, common.ERR_REQ_BIND)
	}

	switch u.Tp {
	case "undo":
		lib.MessagePush(u.WorkId, 6, "")
		model.DB().Model(model.CoreSqlOrder{}).Where("work_id =?", u.WorkId).Updates(&model.CoreSqlOrder{Status: 6})
		return c.JSON(http.StatusOK, common.SuccessPayLoadToMessage(common.ORDER_IS_UNDO))
	case "agree":
		return c.JSON(http.StatusOK, MultiAuditOrder(u, user.Username))
	case "reject":
		return c.JSON(http.StatusOK, RejectOrder(u, user.Username))
	default:
		return c.JSON(http.StatusOK, common.ERR_REQ_FAKE)
	}
}

type WorkflowCallBackParam struct {
	Context map[string]interface{} `json:"context"`
	//"deny" "success"
	Status           string `json:"status"`
	CurrentNodeName  string `json:"current_node_name"`
	PreviousNodeName string `json:"previous_node_name"`
	FlowInstanceID   string `json:"flow_instance_id"`
}

//type WorkflowCallBackParam map[string]string

//返回
//{
//    "request_id": "6bd85fce-8de9-4116-bf2a-acb891f443f2",
//    "result_code": "success",
//    "data": {
//        "context": {
//            "day": "123"
//        }
//    }
//}

func OpenAuditOrderState(c yee.Context) (err error) {
	c.Logger().Info("workflow callback................")
	u := new(WorkflowCallBackParam)

	if err = c.Bind(u); err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusOK, common.ERR_REQ_BIND)
	}

	confirm := new(Confirm)

	confirm.WorkId = u.Context["flowID"].(string)

	username := u.Context["applier"].(string)

	c.Logger().Infof("parsed param ", confirm.WorkId, username)
	flowDetail, err := lib.CallBackWorkflowInstance(u.FlowInstanceID, username)
	if err != nil {
		c.Logger().Error("call get workflow instance failed" + err.Error())
		return c.JSON(http.StatusBadRequest, lib.WorkflowResponse{
			RequestID:  "uuid",
			ResultCode: "success",
			Data:       map[string]string{},
		})
	}

	var operators []lib.FlowOperator
	for _, node := range flowDetail.Nodes {
		if node.FlowInstanceNodeTitle == "审批" {
			for _, task := range node.Tasks {
				switch task.Status {
				case "deny", "success":
					operators = append(operators, task.Operator)
				default:
				}
			}
		}
	}

	auditUser := username

	if len(operators) > 0 {
		auditUser = operators[0].UserName
	}

	logger.DefaultLogger.Errorf("audituser: " + auditUser + "status:" + u.Status)
	switch u.Status {
	case "success":
		OpenAuditOrder(confirm, auditUser)
		return c.JSON(http.StatusBadRequest, lib.WorkflowResponse{
			RequestID:  "uuid",
			ResultCode: "success",
			Data:       map[string]string{},
		})
	case "deny":
		RejectOrder(confirm, auditUser)
		return c.JSON(http.StatusBadGateway, lib.WorkflowResponse{
			RequestID:  "uuid",
			ResultCode: "success",
			Data:       map[string]string{},
		})
	default:
		return c.JSON(http.StatusBadRequest, lib.WorkflowResponse{
			RequestID:  "uuid",
			ResultCode: "success",
			Data:       map[string]string{},
		})
	}
}

// DelayKill will stop delay order
func DelayKill(c yee.Context) (err error) {
	u := new(Confirm)
	if err = c.Bind(u); err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusOK, common.ERR_REQ_BIND)
	}
	user := new(lib.Token).JwtParse(c)
	model.DB().Create(&model.CoreWorkflowDetail{
		WorkId:   u.WorkId,
		Username: user.Username,
		Time:     time.Now().Format("2006-01-02 15:04"),
		Action:   ORDER_KILL_STATE,
	})
	return c.JSON(http.StatusOK, common.SuccessPayLoadToMessage(delayKill(u.WorkId)))
}

func FetchAuditOrder(c yee.Context) (err error) {
	u := new(common.PageList[[]model.CoreSqlOrder])
	if err = c.Bind(u); err != nil {
		c.Logger().Error(err.Error())
		return
	}
	user := new(lib.Token).JwtParse(c)
	u.Paging().Select(QueryField).Query(common.AccordingToAllOrderState(u.Expr.Status),
		common.AccordingToAllOrderType(u.Expr.Type),
		common.AccordingToRelevant(user.Username),
		common.AccordingToText(u.Expr.Text),
		common.AccordingToUsernameEqual(u.Expr.Username),
		common.AccordingToDatetime(u.Expr.Picker))
	return c.JSON(http.StatusOK, u.ToMessage())
}

func FetchOSCAPI(c yee.Context) (err error) {
	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()
		workId := c.QueryParam("work_id")
		var msg string
		for {
			if workId != "" {
				var osc model.CoreSqlOrder
				model.DB().Model(model.CoreSqlOrder{}).Where("work_id =?", workId).Find(&osc)
				err := websocket.Message.Send(ws, osc.OSCInfo)
				if err != nil {
					c.Logger().Error(err)
					break
				}
			}
			if err := websocket.Message.Receive(ws, &msg); err != nil {
				break
			}
			if msg == common.CLOSE {
				break
			}
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}

func AuditOrderApis(c yee.Context) (err error) {
	switch c.Params("tp") {
	case "state":
		return AuditOrderState(c)
	case "kill":
		return DelayKill(c)
	default:
		return c.JSON(http.StatusOK, common.ERR_REQ_FAKE)
	}
}

func OpenAuditOrderApis(c yee.Context) (err error) {
	return OpenAuditOrderState(c)
}

func AuditOrRecordOrderFetchApis(c yee.Context) (err error) {
	switch c.Params("tp") {
	case "list":
		return FetchAuditOrder(c)
	//case "record":
	//	return FetchRecord(c)
	default:
		return c.JSON(http.StatusOK, common.ERR_REQ_FAKE)
	}
}

func AuditOSCFetchAndKillApis(c yee.Context) (err error) {
	switch c.Params("tp") {
	case "osc":
		return FetchOSCAPI(c)
	case "kill":
		return nil
	default:
		return c.JSON(http.StatusOK, common.ERR_REQ_FAKE)
	}
}
