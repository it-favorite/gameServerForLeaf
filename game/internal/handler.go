package internal
import (
	"server/db"
    "reflect"  
    "server/msg"
    "server/datastruct"
    "github.com/name5566/leaf/gate"  
    "github.com/name5566/leaf/log"
    "github.com/name5566/leaf/network/json"
    "server/tools"
    "server/game/internal/match"
)

// 异步处理  
func handleMsg(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func init() {
    handleMsg(&msg.CS_PlayerMatching{}, handleSinglePersonMatching)
    handleMsg(&msg.CS_EndlessModeMatching{},handleEndlessModeMatching)
    handleMsg(&msg.CS_InviteModeMatching{},handleInviteModeMatching)
    handleMsg(&msg.CS_JoinInviteMode{},handleJoinInviteMode)
    handleMsg(&msg.CS_LeaveInviteMode{},handleLeaveInviteMode)
    handleMsg(&msg.CS_MasterFirePlayer{},handleMasterFirePlayer)
    handleMsg(&msg.CS_MasterStartGame{},handleMasterStartGame)

    handleMsg(&msg.CS_PlayerCancelMatching{}, handleCancelMatching)
    handleMsg(&msg.CS_PlayerJoinRoom{}, handlePlayerJoinRoom)
    handleMsg(&msg.CS_MoveData{}, handlePlayerMoveData)
    handleMsg(&msg.CS_EnergyExpended{}, handleEnergyExpended)
    handleMsg(&msg.CS_PlayerDied{}, handlePlayersDied)
    handleMsg(&msg.CS_PlayerLeftRoom{}, handlePlayerLeftRoom)
    handleMsg(&msg.CS_PlayerRelive{}, handlePlayerRelive)

    handleMsg(&msg.CS_GameOver1{}, handlePlayerGameOver1)
    
    handleMsg(&msg.CS_GameOverSinglePersonMode{}, handlePlayerGameOverSinglePersonMode)
    handleMsg(&msg.CS_GameOverInviteMode{}, handlePlayerGameOverInviteMode)

    handleMsg(&msg.CS_GetSnakeLength{}, handleGetSnakeLength)
    handleMsg(&msg.CS_GetKillNum{}, handleGetKillNum)
}


func getParentMatch(mode datastruct.GameModeType) match.ParentMatch{
    var match match.ParentMatch
    match = nil
    switch mode{
     case  datastruct.SinglePersonMode:
         match = ptr_singleMatch
     case datastruct.EndlessMode:
         match = ptr_endlessModeMatch
     case datastruct.InviteMode:
         match = ptr_inviteModeMatch 
    }
    return match
}

func handleGetSnakeLength(args []interface{}){
    a := args[1].(gate.Agent)
    if !tools.IsValid(a){
       return
    }
    m := args[0].(*msg.CS_GetSnakeLength)
    agentUserData := tools.GetUserData(a)
    mode:=m.MsgContent.GameMode
    currentPlayerData,arr:=db.Module.GetSnakeLengthRank(agentUserData.Uid,mode,m.MsgContent.RankStart,m.MsgContent.RankEnd)
    currentPlayerData.Avatar = agentUserData.Extra.Avatar
    currentPlayerData.Name = agentUserData.Extra.PlayName
    a.WriteMsg(msg.GetSnakeLengthMsg(mode,currentPlayerData,arr))
}

func handleGetKillNum(args []interface{}){
    a := args[1].(gate.Agent)
    if !tools.IsValid(a){
       return
    }
    m := args[0].(*msg.CS_GetKillNum)
    agentUserData := tools.GetUserData(a)
    mode:=m.MsgContent.GameMode
    currentPlayerData,arr:=db.Module.GetKillNumRank(agentUserData.Uid,mode,m.MsgContent.RankStart,m.MsgContent.RankEnd)
    currentPlayerData.Avatar = agentUserData.Extra.Avatar
    currentPlayerData.Name = agentUserData.Extra.PlayName
    a.WriteMsg(msg.GetKillNumMsg(mode,currentPlayerData,arr))

}

func handlePlayerGameOverInviteMode(args []interface{}){
    a := args[1].(gate.Agent)
    if !tools.IsValid(a){
       return
    }
    m := args[0].(*msg.CS_GameOverInviteMode)
    tf:=tools.EnableSettle(m.MsgContent.RoomID,a)
    agentUserData := tools.GetUserData(a)
    if tf&&agentUserData.GameMode==datastruct.InviteMode{
        uid:=agentUserData.Uid
        score:=m.MsgContent.Score
        killNum:=m.MsgContent.KillNum
        integral:=tools.GetGameIntegral(m.MsgContent.Ranking)
        db.Module.AddGameIntegral(uid,integral)
        maxScore,maxKillNum:=db.Module.GetMaxScoreInviteMode(uid)
        isUpdate:=false
        if score>maxScore{
           maxScore = score
           isUpdate = true
        }
        if killNum>maxKillNum{
           maxKillNum = killNum
           isUpdate = true
        }
        if isUpdate{
           db.Module.UpdateMaxScoreInviteMode(uid,maxScore,maxKillNum)
        }
    }
}

func handlePlayerGameOverSinglePersonMode(args []interface{}){
    a := args[1].(gate.Agent)

    if !tools.IsValid(a){
       return
    }
    
    m := args[0].(*msg.CS_GameOverSinglePersonMode)
    tf:=tools.EnableSettle(m.MsgContent.RoomID,a)
    agentUserData := tools.GetUserData(a)
    if tf&&agentUserData.GameMode==datastruct.SinglePersonMode{
        uid:=agentUserData.Uid
        score:=m.MsgContent.Score
        killNum:=m.MsgContent.KillNum
        fragmentNum:=tools.GetFragmentNum(m.MsgContent.Ranking)
        db.Module.AddFragmentNum(uid,fragmentNum)
        maxScore,maxKillNum:=db.Module.GetMaxScoreInSinglePersonMode(uid)
        isUpdate:=false
        if score>maxScore{
           maxScore = score
           isUpdate = true
        }
        if killNum>maxKillNum{
           maxKillNum = killNum
           isUpdate = true
        }
        if isUpdate{
           db.Module.UpdateMaxScoreInSinglePersonMode(uid,maxScore,maxKillNum)
        }
    }
}

func handlePlayerGameOver1(args []interface{}){
    a := args[1].(gate.Agent)
    if !tools.IsValid(a){
       return
    }
    agentUserData := tools.GetUserData(a)
    m := args[0].(*msg.CS_GameOver1)
    tf:=tools.EnableSettle(m.MsgContent.RoomID,a)
    if tf&&agentUserData.GameMode==datastruct.EndlessMode{
        uid:=agentUserData.Uid
        score:=m.MsgContent.Score
        killNum:=m.MsgContent.KillNum
        maxScore,maxKillNum:=db.Module.GetMaxScoreInEndlessMode(uid)
        isUpdate:=false
        if score>maxScore{
            maxScore = score
            isUpdate = true
        }
        if killNum>maxKillNum{
            maxKillNum = killNum
            isUpdate = true
        }
        a.WriteMsg(msg.GetGameOver1Msg(maxScore,m.MsgContent.Score,m.MsgContent.KillNum))
        if isUpdate{
            db.Module.UpdateMaxScoreInEndlessMode(uid,maxScore,maxKillNum)
        }
    }
}

func handlePlayerRelive(args []interface{}){
    a := args[1].(gate.Agent)
    if !tools.IsValid(a){
       return
    }
    agentUserData := tools.GetUserData(a)
    switch agentUserData.GameMode{
       case datastruct.EndlessMode:
        ptr_endlessModeMatch.PlayerRelive(agentUserData.Extra.RoomID,agentUserData.PlayId,agentUserData.Extra.PlayName,agentUserData.Extra.Avatar)
    }
}

func handlePlayerLeftRoom(args []interface{}){
    a := args[1].(gate.Agent)
    if !tools.IsValid(a){
       return
    }
    agentUserData := tools.GetUserData(a)
    playerLeftRoom(agentUserData.ConnUUID,agentUserData.GameMode,agentUserData.Extra.RoomID)
}

func handlePlayersDied(args []interface{}){
    a := args[1].(gate.Agent)
    if !tools.IsValid(a){
       return
    }
    agentUserData := tools.GetUserData(a)
    m := args[0].(*msg.CS_PlayerDied)
    
    //接收玩家死亡坐标,生成指定范围能量点
    //指定某一帧复活
    match:=getParentMatch(agentUserData.GameMode)
    match.PlayersDied(agentUserData.Extra.RoomID,m.MsgContent)
}

func handleEnergyExpended(args []interface{}){
    a := args[1].(gate.Agent)
    if !tools.IsValid(a){
       return
    }
    agentUserData := tools.GetUserData(a)
    m := args[0].(*msg.CS_EnergyExpended)
    expended:=m.MsgContent.EnergyExpended
    if expended>0{
        match:=getParentMatch(agentUserData.GameMode)
        match.EnergyExpended(expended,*agentUserData)
    }
}

func handlePlayerMoveData(args []interface{}){
    //测试
    //msg.Num = 0
    a := args[1].(gate.Agent)
    if !tools.IsValid(a){
       return
    }
    agentUserData := tools.GetUserData(a)
    r_id:=agentUserData.Extra.RoomID
    m := args[0].(*msg.CS_MoveData)
    match:=getParentMatch(agentUserData.GameMode)
    match.PlayerMoved(r_id,agentUserData.PlayId,m)
}

func handlePlayerJoinRoom(args []interface{}){
    a := args[1].(gate.Agent)
    if !tools.IsValid(a){
       return
    }
    agentUserData := tools.GetUserData(a)
    connUUID:=agentUserData.ConnUUID
    m := args[0].(*msg.CS_PlayerJoinRoom)
    match:=getParentMatch(agentUserData.GameMode)
    match.PlayerJoin(connUUID,m)
}

func handleCancelMatching(args []interface{}){
    a := args[1].(gate.Agent)
    if !tools.IsValid(a){
       return
    }
    agentUserData := tools.GetUserData(a)
    connUUID:=agentUserData.ConnUUID

    switch agentUserData.GameMode{
    case datastruct.SinglePersonMode:
         ptr_singleMatch.RemovePlayerFromMatchingPool(connUUID)
    case datastruct.EndlessMode:
         ptr_endlessModeMatch.RemovePlayer(connUUID)
    }
}
 

//收到单人匹配消息的时候加入池，主动离开和自动离开在池中删除，
//完成单人匹配后，在池中删除
func handleSinglePersonMatching(args []interface{}) {
     startMatching(args,datastruct.SinglePersonMode)   
}

func handleEndlessModeMatching(args []interface{}){
     startMatching(args,datastruct.EndlessMode)
}

func handleInviteModeMatching(args []interface{}){
     startMatching(args,datastruct.InviteMode)
}


func removePlayerFromOtherMatchs(connUUID string,mode datastruct.GameModeType){
     switch mode{
      case datastruct.SinglePersonMode:
           removePlayer(connUUID,datastruct.EndlessMode)
           removePlayer(connUUID,datastruct.InviteMode)
      case datastruct.EndlessMode:
           removePlayer(connUUID,datastruct.SinglePersonMode)
           removePlayer(connUUID,datastruct.InviteMode)
      case datastruct.InviteMode:
           removePlayer(connUUID,datastruct.SinglePersonMode)
           removePlayer(connUUID,datastruct.EndlessMode)
     }
}

func removePlayer(key string,mode datastruct.GameModeType){
    match:=getParentMatch(mode)
    match.RemovePlayer(key)
}

func playerLeftRoom(connUUID string,mode datastruct.GameModeType,r_id string){
    match:=getParentMatch(mode)
    if match != nil{
       match.PlayerLeftRoom(r_id,connUUID) 
    }
}

func startMatching(args []interface{},mode datastruct.GameModeType){
    a := args[1].(gate.Agent)
    if !tools.IsValid(a){
       return
    }
    agentUserData := tools.GetUserData(a)
    
    uid:=agentUserData.Uid
    if uid <= 0{
        log.Error("Uid error : %v",uid)
        return
    }
    connUUID:=agentUserData.ConnUUID
    
    
    //重置UserData
    tools.ReSetAgentUserData(uid,mode,datastruct.NULLID,a,connUUID,tools.ReSetExtraRoomID(agentUserData.Extra))
    removePlayerFromOtherMatchs(connUUID,mode)
    
    if checkActionPool(connUUID,mode,a){
       return
    }
    
    matchingChanRPC(mode,connUUID,a,uid)
    
    if mode != datastruct.InviteMode {
      var msgHeader json.MsgHeader
      msgHeader.MsgName = msg.SC_PlayerMatchingKey
    
      var msgContent msg.SC_PlayerMatchingContent
      msgContent.IsMatching =true
    
      a.WriteMsg(&msg.SC_PlayerMatching{  
        MsgHeader:msgHeader,
        MsgContent:msgContent,
      })
    }
}

func checkActionPool(connUUID string,mode datastruct.GameModeType,a gate.Agent) bool {
    isMatching:=false
    var match match.ParentMatch
    switch mode{
     case  datastruct.SinglePersonMode:
         match = ptr_singleMatch
     case datastruct.EndlessMode:
         match = ptr_endlessModeMatch
     case datastruct.InviteMode:
         match = ptr_inviteModeMatch 
    }
    if match.CheckActionPool(connUUID){//已在匹配中
       isMatching = true
    }
    if isMatching{
        var msgHeader json.MsgHeader
        msgHeader.MsgName = msg.SC_PlayerAlreadyMatchingKey
        a.WriteMsg(&msg.SC_PlayerAlreadyMatching{
            MsgHeader:msgHeader,
        })
    }
    return isMatching
}

func matchingChanRPC(mode datastruct.GameModeType,connUUID string,a gate.Agent,uid int){
    var match match.ParentMatch
    switch mode{
     case  datastruct.SinglePersonMode:
         match = ptr_singleMatch
     case datastruct.EndlessMode:
         match = ptr_endlessModeMatch
     case datastruct.InviteMode:
         match = ptr_inviteModeMatch
    }
    ChanRPC.Go(MatchingKey,match,connUUID,a,uid)//玩家匹配
}

func handleJoinInviteMode(args []interface{}){
   
    a := args[1].(gate.Agent)
    if !tools.IsValid(a){
       return
    }
    agentUserData := tools.GetUserData(a)
    
    uid:=agentUserData.Uid
    if uid <= 0{
        log.Error("Uid error : %v",uid)
        return
    }
    connUUID:=agentUserData.ConnUUID
     
    tools.ReSetAgentUserData(uid,datastruct.InviteMode,datastruct.NULLID,a,connUUID,tools.ReSetExtraRoomID(agentUserData.Extra))
    removePlayerFromOtherMatchs(connUUID,datastruct.InviteMode)
    
    if checkActionPool(connUUID,datastruct.InviteMode,a){
       return
    }
    
    m := args[0].(*msg.CS_JoinInviteMode)
    w_id:=m.MsgContent.RoomID
    ptr_inviteModeMatch.JoinWaitRoom(w_id,a,agentUserData.Uid,connUUID)
}

func handleLeaveInviteMode(args []interface{}){
    a := args[1].(gate.Agent)
    if !tools.IsValid(a){
       return
    }
    agentUserData:= tools.GetUserData(a)
    ptr_inviteModeMatch.LeftWaitRoom(agentUserData.Extra.WaitRoomID,agentUserData.ConnUUID)
}

func handleMasterFirePlayer(args []interface{}){
    a := args[1].(gate.Agent)
    if !tools.IsValid(a){
       return
    }
    agentUserData:= tools.GetUserData(a)
    m := args[0].(*msg.CS_MasterFirePlayer)
    seat := m.MsgContent.Seat
    ptr_inviteModeMatch.MasterFirePlayer(agentUserData.Extra.WaitRoomID,agentUserData.ConnUUID,seat)
}

func handleMasterStartGame(args []interface{}){
    a := args[1].(gate.Agent)
    if !tools.IsValid(a){
       return
    }
    agentUserData:= tools.GetUserData(a)
    ptr_inviteModeMatch.StartGame(agentUserData.Extra.WaitRoomID,agentUserData.ConnUUID)
}

func leaveWaitRoom(w_id string,connUUID string){
    ptr_inviteModeMatch.LeftWaitRoom(w_id,connUUID)
}

func sendInviteQRCode(r_id string,qrcode string){
    ptr_inviteModeMatch.SendInviteQRCode(r_id,qrcode)
}

func relogin(loginName string,a gate.Agent){
    onlinePlayersData.Mutex.Lock()
    defer onlinePlayersData.Mutex.Unlock()
    agent,tf:=onlinePlayersData.LoginNames[loginName]
    onlinePlayersData.LoginNames[loginName]=a
    if tf && agent.RemoteAddr().String() != a.RemoteAddr().String(){
        agent.Close()
        agent.Destroy()
    }
}

func deleteOnlinePlayersData(ip_str string){
    onlinePlayersData.Mutex.Lock()
    defer onlinePlayersData.Mutex.Unlock()
    for k,v := range onlinePlayersData.LoginNames{
        if v.RemoteAddr().String() == ip_str{
            delete(onlinePlayersData.LoginNames,k)
            break
        }
    }
}