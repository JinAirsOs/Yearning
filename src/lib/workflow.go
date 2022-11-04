package lib

import (
	"Yearning-go/src/model"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"github.com/cookieY/yee/logger"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const (
	WORKFLOW_API string = "https://rocket.nioint.com"
)

type WorkflowResponse struct {
	RequestID  string      `json:"request_id"`
	ResultCode string      `json:"result_code"`
	Data       interface{} `json:"data"`
}

type WorkflowRequest struct {
	Creator  string            `json:"creator"`
	FlowCode string            `json:"flow_code"`
	Context  map[string]string `json:"context"`
}

func NewWorkflowRequest() *WorkflowRequest {
	flowCode := os.Getenv("FLOWCODE")

	if len(flowCode) == 0 {
		flowCode = "22792e0b-5aa7-4c7c-b7de-c5fd6ca6fede"
	}

	return &WorkflowRequest{
		Creator:  "zhendong.pan",
		FlowCode: flowCode,
		Context:  make(map[string]string),
	}
}

func WorkflowSign(url, appID string, appSecret string) string {
	return url + "?operator_worker_user_id="
}

func SendWorkflowRequest(msg *WorkflowRequest, path string, username string) (*WorkflowResponse, error) {

	//创建一个请求
	url := WorkflowSign(WORKFLOW_API+path, "100678", "appSecret") + username
	body, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		logger.DefaultLogger.Errorf("request:", err)
		return nil, err
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}
	//设置请求头
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("X-Domain-Id", msg.Creator)
	//发送请求
	resp, err := client.Do(req)

	if err != nil {
		logger.DefaultLogger.Errorf("resp:", err)
		return nil, err
	}

	//关闭请求
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("http code != 200 error")
	}

	b, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var response WorkflowResponse
	err = json.Unmarshal(b, &response)

	if err != nil {
		return nil, err
	}

	return &response, nil
}

func CreateWorkflowInstance(order *model.CoreSqlOrder, user *Token) (workflowID string, err error) {
	req := NewWorkflowRequest()

	req.Creator = user.Username
	var orderType string
	switch order.Type {
	case 0:
		orderType = "DDL"
	case 1:
		orderType = "DML"
	case 2:
		orderType = "QUERY"
	default:
		orderType = "COMMON"
	}
	req.Context["flowType"] = orderType
	req.Context["flowID"] = order.WorkId
	req.Context["applier"] = user.Username
	var area string
	if strings.Contains(area, "欧洲") {
		area = "欧洲"
	} else {
		area = "国内"
	}
	req.Context["area"] = area

	var env string
	if strings.Contains(order.Source, "prod") {
		env = "prod"
	} else {
		env = order.Source
	}

	req.Context["env"] = env
	req.Context["text"] = order.Text

	req.Context["dataSource"] = order.Source
	req.Context["sql"] = SQLFormat(order.SQL)

	resp, err := SendWorkflowRequest(req, "/api/v1/instance/create", user.Username)

	if err != nil {
		logger.DefaultLogger.Errorf("send workflow request failed with:", err)
		return "", err
	}

	if resp.ResultCode != "success" {
		return "", errors.New("create workflow instance failed")
	}

	type CreateWorkflowInstanceResp struct {
		WorkflowInstanceID string `json:"flow_instance_id"`
	}

	data := make(map[string]string)

	v, _ := json.Marshal(resp.Data)

	err = json.Unmarshal(v, &data)
	if err != nil {
		return "", errors.New("typo: not type of CreateWorkflowInstanceResp")
	}

	return data["flow_instance_id"], nil
}

func SQLFormat(sql string) string {
	if len(sql) < 255 {
		return sql
	}

	sqlSlice := strings.Split(sql, "\n")
	if len(sqlSlice) > 15 {
		sqlSlice = sqlSlice[:15]
	}
	formatedSQL := strings.Join(sqlSlice, " \n")

	formatedSQL += "----------------------------------------\n 略...... \n"
	return formatedSQL
}

type FlowDetail struct {
	FlowName       string                 `json:"flow_name"`
	FlowInstanceID string                 `json:"flow_instance_id"`
	Status         string                 `json:"status"`
	Context        map[string]interface{} `json:"context"`

	Nodes []FlowNode `json:"nodes"`
}

type FlowNode struct {
	Name string `json:"name"`
	//"审批"
	FlowInstanceNodeTitle string `json:"flow_instance_node_title"`
	//eg. "or"
	FlowInstanceNodeRule string `json:"flow_instance_node_rule"`

	Tasks []FlowTask `json:"tasks"`
}

type FlowTask struct {
	Operator  FlowOperator `json:"operator"`
	Status    string       `json:"status"`
	ExtraInfo ExtraInfo    `json:"extra_info"`
}

type ExtraInfo struct {
	Comment string `json:"comment"`
}

type FlowOperator struct {
	UserName string `json:"user_name"`
	UserID   string `json:"worker_user_id"`
}

func CallBackWorkflowInstance(workflowInstanceID string, username string) (*FlowDetail, error) {

	url := WORKFLOW_API + "/api/v1/instance/" + workflowInstanceID
	logger.DefaultLogger.Infof("CallBackWorkflowInstance request: " + url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.DefaultLogger.Errorf("request:", err)
		return nil, err
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}
	//设置请求头
	req.Header.Set("X-Domain-Id", username)
	//发送请求
	resp, err := client.Do(req)

	if err != nil {
		logger.DefaultLogger.Errorf("err resp:", url, err)
		return nil, err
	}

	//关闭请求
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	logger.DefaultLogger.Infof("url resp:", url, string(body))
	type FlowDetailResp struct {
		RequestID  string     `json:"request_id"`
		ResultCode string     `json:"result_code"`
		Data       FlowDetail `json:"data"`
	}
	var flowDetail FlowDetailResp
	err = json.Unmarshal(body, &flowDetail)
	if err != nil {
		logger.DefaultLogger.Errorf("unmarshal err:", err)
		return nil, err
	}

	return &flowDetail.Data, nil
}

func RevokeWorkflow(workID string, username string) error {
	var order model.CoreSqlOrder
	model.DB().Where("work_id = ?", workID).Find(&order)
	url := WORKFLOW_API + "/api/v1/instance/" + order.WorkflowID + "/revoke"
	logger.DefaultLogger.Infof("RevokeWorkflow request: " + url)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		logger.DefaultLogger.Errorf("request:", err)
		return err
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}
	//设置请求头
	req.Header.Set("X-Domain-Id", username)
	//发送请求
	resp, err := client.Do(req)

	if err != nil {
		logger.DefaultLogger.Errorf("err resp:", url, err)
		return err
	}

	//关闭请求
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("err response")
	}

	return nil
}
