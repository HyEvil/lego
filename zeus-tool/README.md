## 个人框架zeus的工具链

### 功能

解析zeus协议并转换成proto，生成swagger，protobuf代码;

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

