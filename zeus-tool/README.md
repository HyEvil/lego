## 个人框架zeus的工具链(最近为了trade_bot造的轮子)

### 功能

解析zeus协议并转换成proto，生成protobuf代码，swagger，router register,packr打包等

生成service，dao，model等

#### zeus协议示例(类似protobuf，修改了rpc的语法，扩展了部分功能等，如中间件，路由，各种tag等)：

```package bot.api;

import "google/protobuf/any.proto";

@middleware(auth.RequireRole(model.RoleTypeDev))
@resource("/api/bot")
service  Bot{
    @get("/list")
    list(ListBotReq) ListBotResp;

    @get("/detail")
    detail(BotDetailReq) BotDetailResp;

    @post("/add")
    add(AddBotReq);

    @post("/delete")
    delete(DeleteBotReq);

    @post("/update")
    update(UpdateBotReq);

    @post("/updateDev")
    updateDev(UpdateBotDevReq);
}

@resource("/api/user")
service  User{
    @post("/login")
    login(LoginReq) LoginResp;
    @get("/current")
    currentUser() CurrentUserResp;
    @post("/logout")
    logout();

    @get("list")
    @middleware(auth.RequireRole(model.RoleTypeAdmin))
    list(ListUserReq) ListUserResp;

    @post("update")
    @middleware(auth.RequireRole(model.RoleTypeAdmin))
    update(UpdateUserReq);

    @post("add")
    @middleware(auth.RequireRole(model.RoleTypeAdmin))
    add(AddUserReq);

    @post("delete")
    @middleware(auth.RequireRole(model.RoleTypeAdmin))
    delete(DeleteUserReq);
}


@middleware(auth.RequireLogin)
@resource("/api/exchange")
service  Exchange{
    @get("/list")
    list(ListExchangeReq) ListExchangeResp;

    @post("/add")
    add(AddExchangeReq);

    @post("/delete")
    delete(DeleteExchangeReq);

    @post("/update")
    update(UpdateExchangeReq);

    @get("/names")
    nameList() NameListResp;

    @get("all")
    all() AllExchangeResp;
}

@middleware(auth.RequireLogin)
@resource("/api/market")
service  Market{
    @get("/list")
    list(ListMarketReq) ListMarketResp;
}

message ListBotReq {
    @tag(validator="required")
    int32 current;
    @tag(validator="required,lte=20")
    int32 pageSize;
    int64 id;
    string name;
    int32 botType;
}

message ListBotResp {
    message Item {
        int64 id;
        string name;
        int32 botType;
        string version;
        int32 updateTime;
    }
    bool success;
    int32 total;
    int32 current;
    int32 pageSize;
    repeated Item data;
}
message ListUserReq {
    @tag(validator="required")
    int32 current;
    @tag(validator="required,lte=20")
    int32 pageSize;
    int64 id;
    string user;
    string role;
}

message ListExchangeReq {
    @tag(validator="required")
    int32 current;
    @tag(validator="required,lte=20")
    int32 pageSize;
    int64 id;
    string alias;
    string name;
    string key;
}

message ListExchangeResp {
    message Item {
        int64 id;
        string alias;
        string name;
        string key;
        string secret;
        string pass;
        int32 updateTime;
    }
    bool success;
    int32 total;
    int32 current;
    int32 pageSize;
    repeated Item data;
}

message AddExchangeReq {
    @tag(validator="required")
    string alias;
    @tag(validator="required")
    string name;
    @tag(validator="required")
    string key;
    @tag(validator="required")
    string secret;
    string pass;
}

message DeleteExchangeReq {
    int64 id;
}

message UpdateExchangeReq {
    @tag(validator="required")
    int64 id;
    @tag(validator="required")
    string alias;
    @tag(validator="required")
    string name;
    @tag(validator="required")
    string key;
    @tag(validator="required")
    string secret;
    string pass;
}

message NameListResp {
    message Item {
        string label;
        string name;
    }
    repeated Item data;
}

message BotDetailReq {
     @tag(validator="required")
     int64 id;
}
message ListMarketReq {
     @tag(validator="required")
     int32 current;
     @tag(validator="required,lte=20")
     int32 pageSize;
     int64 id;
     string name;
     string developer;
}

message ListHostReq {
     @tag(validator="required")
     int32 current;
     @tag(validator="required,lte=20")
     int32 pageSize;
     int64 id;
     string name;
}

message HostChartResp {
    @tag(customtype="[]map[string]interface{}",nullable=false)
    bytes data;
}
```
### 生成的dao示例

```
package repository

import (
	"yym/trade_bot/common/apiec"
	"yym/trade_bot/dao/model"
	"yym/trade_bot/dao/model/table"

	"go.yym.plus/zeus/pkg/sql"
	"xorm.io/builder"
	"xorm.io/xorm"
)

var (
	Bot = bot{}
)

type bot struct {
}

func (self *bot) session(session ...*xorm.Session) *xorm.Session {
	if len(session) > 0 {
		return session[0].Table(table.Bot.Model())
	} else {
		return table.Bot.Session()
	}
}

func (self *bot) OneById(id int64, session ...*xorm.Session) (*model.Bot, error) {
	s := self.session(session...)
	m := model.Bot{}
	ok, err := s.ID(id).Get(&m)
	if err != nil {
		return nil, apiec.DB.WithError(err)
	}
	if !ok {
		return nil, nil
	}
	return &m, nil
}

func (self *bot) UpdateById(id int64, bean interface{}, session ...*xorm.Session) *sql.AffectedResult {
	s := self.session(session...)
	affected, err := s.ID(id).Update(bean)
	if err != nil {
		err = apiec.DB.WithError(err)
	}
	return sql.NewAffectedResult(affected, err)
}

func (self *bot) DeleteById(id int64, session ...*xorm.Session) *sql.AffectedResult {
	s := self.session(session...)
	affected, err := s.ID(id).Delete(table.Bot.Model())
	if err != nil {
		err = apiec.DB.WithError(err)
	}
	return sql.NewAffectedResult(affected, err)
}

func (self *bot) Add(m *model.Bot, session ...*xorm.Session) *sql.InsertResult {
	s := self.session(session...)
	affected, err := s.Insert(m)
	if err != nil {
		err = apiec.DB.WithError(err)
	}
	return sql.NewInsertResult(affected, err)
}

func (self *bot) Exist(cond builder.Cond, session ...*xorm.Session) (bool, error) {
	s := self.session(session...)
	ok, err := s.Exist(cond)
	if err != nil {
		return false, apiec.DB.WithError(err)
	}
	return ok, nil
}

func (self *bot) One(cond builder.Cond, session ...*xorm.Session) (*model.Bot, error) {
	s := self.session(session...)
	m := model.Bot{}
	ok, err := s.Where(cond).Get(&m)
	if err != nil {
		return nil, apiec.DB.WithError(err)
	}
	if !ok {
		return nil, nil
	}
	return &m, nil
}

func (self *bot) List(cond builder.Cond, session ...*xorm.Session) ([]*model.Bot, *sql.FindResult) {
	s := self.session(session...)
	list := []*model.Bot{}
	count, err := s.Where(cond).FindAndCount(&list)
	if err != nil {
		return nil, sql.NewFindResult(0, apiec.DB.WithError(err))
	}
	return list, sql.NewFindResult(count, nil)
}

func (self *bot) Find(cond builder.Cond, extCond ...sql.Cond) ([]*model.Bot, *sql.FindResult) {
	s := self.session()
	list := []*model.Bot{}
	s = s.Where(cond)
	for _, c := range extCond {
		err := c.Apply(s)
		if err != nil {
			return nil, sql.NewFindResult(0, apiec.Internal.WithError(err))
		}
	}

	count, err := s.FindAndCount(&list)
	if err != nil {
		return nil, sql.NewFindResult(0, apiec.DB.WithError(err))
	}
	return list, sql.NewFindResult(count, nil)
}

func (self *bot) FindEx(cond builder.Cond, session *xorm.Session, extCond ...sql.Cond) ([]*model.Bot, *sql.FindResult) {
	s := self.session(session)
	list := []*model.Bot{}
	s = s.Where(cond)
	for _, c := range extCond {
		err := c.Apply(s)
		if err != nil {
			return nil, sql.NewFindResult(0, apiec.Internal.WithError(err))
		}
	}

	count, err := s.FindAndCount(&list)
	if err != nil {
		return nil, sql.NewFindResult(0, apiec.DB.WithError(err))
	}
	return list, sql.NewFindResult(count, nil)
}

func (self *bot) Update(cond builder.Cond, bean interface{}, session ...*xorm.Session) *sql.AffectedResult {
	s := self.session(session...)
	affected, err := s.Where(cond).Update(bean)
	if err != nil {
		err = apiec.DB.WithError(err)
	}
	return sql.NewAffectedResult(affected, err)
}

func (self *bot) Delete(cond builder.Cond, session ...*xorm.Session) *sql.AffectedResult {
	s := self.session(session...)
	affected, err := s.Where(cond).Delete(table.Bot.Model())
	if err != nil {
		err = apiec.DB.WithError(err)
	}
	return sql.NewAffectedResult(affected, err)
}

```

