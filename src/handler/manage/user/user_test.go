// Copyright 2021 HenryYee.
//
// Licensed under the AGPL, Version 3.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.gnu.org/licenses/agpl-3.0.en.html
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package user

import (
	"Yearning-go/src/handler/common"
	"Yearning-go/src/model"
	"Yearning-go/src/test"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func init() {
	model.DBNew("../../../../conf.toml")
	apis.NewTest()
}

var apis = test.Case{
	Method:  http.MethodPost,
	Uri:     "/api/v2/manage/user",
	Handler: SuperUserApi(),
}

func TestFetchUser(t *testing.T) {
	var Ref common.Resp
	apis.Put(`{"page":1,"find":{"valve":false}}`).Do().Unmarshal(&Ref)
	assert.NotEqual(t, nil, Ref.Payload)
	assert.Equal(t, 1200, Ref.Code)

	apis.Put(`1{"page":"0"}`).Do().Unmarshal(&Ref)
	assert.Equal(t, 1310, Ref.Code)
}
