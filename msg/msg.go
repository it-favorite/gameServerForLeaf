package msg

import (
	"github.com/name5566/leaf/network/json"
	"server/datastruct"
)

const PC_Platform ="pc"  //pc端
const WX_Platform ="wx" //微信平台


var Processor = json.NewProcessor()

func init() {
	Processor.Register(&CS_UserLogin{})
	Processor.Register(&SC_UserLogin{})
	Processor.Register(&CS_PlayerMatching{})
	Processor.Register(&CS_EndlessModeMatching{})
	Processor.Register(&SC_PlayerMatching{})
	Processor.Register(&SC_PlayerAlreadyMatching{})
	Processor.Register(&SC_PlayerMatchingEnd{})
	Processor.Register(&CS_PlayerCancelMatching{})

	Processor.Register(&CS_PlayerJoinRoom{})

	Processor.Register(&SC_PlayerReMatch{})

	Processor.Register(&SC_InitRoomData{})
	Processor.Register(&SC_RoomFrameData{})
	Processor.Register(&CS_MoveData{})
	
	Processor.Register(&CS_EnergyExpended{})
	Processor.Register(&CS_PlayerDied{})

	Processor.Register(&CS_PlayerLeftRoom{})

	Processor.Register(&CS_PlayerRelive{})

	Processor.Register(&SC_GameOverData{})

	Processor.Register(&SC_PlayerInWaitRoom{})
	Processor.Register(&CS_InviteModeMatching{})
	Processor.Register(&CS_JoinInviteMode{})
	Processor.Register(&CS_LeaveInviteMode{})
	Processor.Register(&CS_MasterFirePlayer{})
	Processor.Register(&SC_NotifyMsg{})
	Processor.Register(&SC_PlayerIsFired{})
	Processor.Register(&CS_MasterStartGame{})
}

/*接收消耗的能量值*/
type CS_EnergyExpended struct {
	MsgHeader json.MsgHeader
	MsgContent CS_EnergyExpendedContent
}

type CS_EnergyExpendedContent struct {
	EnergyExpended int
}

/*客户端发送来完成注册*/
type CS_UserLogin struct {
	MsgHeader json.MsgHeader
	MsgContent CS_UserLoginContent
}

type CS_UserLoginContent struct {
	LoginName string //如果是微信发送过来就是微信code
	NickName string
	Avatar string
	Platform string //告知服务端是从哪家平台发送过来的,比如"微信","QQ"
}


/*服务端发送给客户端*/
type SC_UserLogin struct {
    MsgHeader json.MsgHeader
	MsgContent SC_UserLoginContent
}
type SC_UserLoginContent struct {
	Uid int //生成的用户id;为-1时,代表没登陆成功
}


/*玩家开始单人匹配模式*/
type CS_PlayerMatching struct {
	MsgHeader json.MsgHeader
}

/*无尽模式匹配*/
type CS_EndlessModeMatching struct {
	MsgHeader json.MsgHeader
}

/*好友模式*/
type CS_InviteModeMatching struct {
	MsgHeader json.MsgHeader
}

/*加入好友模式的房间*/
type CS_JoinInviteMode struct {
	MsgHeader json.MsgHeader
	MsgContent CS_JoinInviteModeContent
}
type CS_JoinInviteModeContent struct{
    RoomID string
}

/*主动离开好友模式的房间*/
type CS_LeaveInviteMode struct {
	MsgHeader json.MsgHeader
}

/*玩家取消匹配*/
type CS_PlayerCancelMatching struct {
	MsgHeader json.MsgHeader
}

/*发送正在匹配中*/
const SC_PlayerMatchingKey = "SC_PlayerMatching"
type SC_PlayerMatching struct {
	MsgHeader json.MsgHeader
	MsgContent SC_PlayerMatchingContent
}
type SC_PlayerMatchingContent struct {
	IsMatching bool
}


/*已在匹配中*/
const SC_PlayerAlreadyMatchingKey = "SC_PlayerAlreadyMatching"
type SC_PlayerAlreadyMatching struct {
	MsgHeader json.MsgHeader
}

/*发送等待房间的数据信息*/
const SC_PlayerInWaitRoomKey = "SC_PlayerInWaitRoom"
type SC_PlayerInWaitRoom struct {
	MsgHeader json.MsgHeader
	MsgContent SC_PlayerInWaitRoomContent
}
type SC_PlayerInWaitRoomContent struct {
	RoomID string
	IsMaster int
	State  datastruct.WaitRoomState
	Players []datastruct.PlayerInWaitRoom
}

/*发送通知信息*/
const SC_NotifyMsgKey = "SC_NotifyMsg"
type SC_NotifyMsg struct {
	MsgHeader json.MsgHeader
	MsgContent SC_NotifyMsgContent
}
type SC_NotifyMsgContent struct {
	Msg string
}

/*发送房主踢人消息*/
const CS_MasterFirePlayerKey = "CS_MasterFirePlayer"
type CS_MasterFirePlayer struct {
	MsgHeader json.MsgHeader
	MsgContent CS_MasterFirePlayerContent
}
type CS_MasterFirePlayerContent struct {
	Seat int //座位号
}

/*发送房主开始游戏消息*/
const CS_MasterStartGameKey = "CS_MasterStartGame"
type CS_MasterStartGame struct {
	MsgHeader json.MsgHeader
}

/*发送被踢消息*/
const SC_PlayerIsFiredKey = "SC_PlayerIsFired"
type SC_PlayerIsFired struct {
	MsgHeader json.MsgHeader
}

/*发送匹配成功的信息*/
const SC_PlayerMatchingEndKey = "SC_PlayerMatchingEnd"
type SC_PlayerMatchingEnd struct {
	MsgHeader json.MsgHeader
	MsgContent SC_PlayerMatchingEndContent
}

type SC_PlayerMatchingEndContent struct {
	RoomID string
}

/*客户端发送来加入房间*/
type CS_PlayerJoinRoom struct {
	MsgHeader json.MsgHeader
	MsgContent CS_PlayerJoinRoomContent
}
type CS_PlayerJoinRoomContent struct {
	RoomID string
}

/*玩家加入房间无效*/
type SC_PlayerJoinInvalid struct {
	MsgHeader json.MsgHeader
}

/*重新开始匹配*/
type SC_PlayerReMatch struct {
	MsgHeader json.MsgHeader
}


type EnergyPointType int

//能量点类型
const (
	TypeA EnergyPointType= 1 +iota
    TypeB 
    TypeC 
    TypeD
)

type Point struct {
    X int
    Y int
}

type Quadrant struct {
    X_Min int
    X_Max int
    Y_Min int
    Y_Max int
}

/*发送给客户端房间初始化数据*/
type SC_InitRoomData struct {
	MsgHeader json.MsgHeader
	MsgContent SC_InitRoomDataContent
}

type SC_InitRoomDataContent struct {
	 MapHeight int//3000
	 MapWidth int//4000
	 CurrentFrameIndex int //游戏进行到当前多少帧,从0开始
	 Interval int //毫秒单位 比如50,代表50毫秒
	 PlayId int //分配给玩家在游戏中的id
	 GameTime int //以毫秒单位
	 GameMode int
}

/*接收客户端的帧数据*/
type CS_MoveData struct {
	MsgHeader json.MsgHeader
	MsgContent CS_MoveDataContent //{"Action":1,"Direction":{X:-1,Y:-2}
}

type CS_MoveDataContent struct {
	X int
	Y int
	Speed int
}

type CS_PlayerDied struct {
	MsgHeader json.MsgHeader
	MsgContent []datastruct.PlayerDiedData
}


type CS_PlayerLeftRoom struct { //玩家离开房间
	MsgHeader json.MsgHeader
}

/*玩家发送请求来复活*/
type CS_PlayerRelive struct {
	MsgHeader json.MsgHeader
}

/*发送给客户端当前帧数据*/
type SC_RoomFrameData struct {
	MsgHeader json.MsgHeader
	MsgContent *SC_RoomFrameDataContent
}

type SC_RoomFrameDataContent struct {
	 FramesData []FrameData
}

type FrameData struct {
	FrameIndex int
	PlayerFrameData []interface{}
	CreateEnergyPoints []datastruct.EnergyPoint
}


type ActionType int
const (
    Create ActionType = iota // value --> 0
    Move              // value --> 1
	Death            // value --> 2
	
	NullAction        
)

/*以下为玩家事件*/
type CreatePlayer struct {//玩家的创建
	 PlayerId int
	 PlayerName string
	 X int
	 Y int
	 AddEnergy int //默认值是0
	 Action ActionType
}

type PlayerIsDied struct {//玩家的死亡
	 PlayerId int
	 Action ActionType
}

var DefaultDirection = Point{X:0,Y:1}
var DefaultSpeed = 1
type PlayerMoved struct {//玩家的移动
	PlayerId int
	Action ActionType
	Speed int//默认速度 1
	X int
	Y int
}

type OfflinePlayerMoved struct {//离线玩家的移动
    Action PlayerMoved
    StartFrameIndex int //从多少帧开始记录的
    SpeedInterval int //加速的时间间隔
	StopSpeedFrameIndex int //持续到多少帧结束 
	DirectionInterval int //转向的时间间隔
}




// //测试
// var Test1Point= Point{X:400,Y:320}
// var Test2Point= Point{X:400,Y:120}

// var Num = 0



type PlayerRelive struct {//玩家的重生
    ReLiveFrameIndex int
	Action CreatePlayer
}


/*发送给客户端游戏结束数据*/
type SC_GameOverData struct {
	MsgHeader json.MsgHeader
	MsgContent *SC_GameOverDataContent
}
type SC_GameOverDataContent struct {
    RoomId string 
}



func GetCreatePlayerAction(p_id int,x int,y int,reLiveFrameIndex int,playerName string,addEnergy int) *PlayerRelive{
	  relive:=new(PlayerRelive)
	  relive.ReLiveFrameIndex = reLiveFrameIndex
	  
	  var action CreatePlayer
	  action.Action = Create
	  action.PlayerId = p_id
	  action.PlayerName = playerName
	  action.AddEnergy = addEnergy
	//   switch Num{
	//   case 0:
	// 	action.X = Test1Point.X
	// 	action.Y = Test1Point.Y
	//   case 1:
	// 	action.X = Test2Point.X
	// 	action.Y = Test2Point.Y
	//   default:
	// 	action.X = x
	// 	action.Y = y
	//   }
	//   Num++
	   action.X = x
	   action.Y = y
	 
	  relive.Action = action
	  return relive
}

func GetCreatePlayerMoved(p_id int,x int,y int,speed int) *PlayerMoved{
	action:=new(PlayerMoved)
	action.Action = Move
	action.PlayerId = p_id
	action.X = x
	action.Y = y
	action.Speed = speed
	return action
}

func UpdatePlayerMoved(move *PlayerMoved,x int,y int,speed int){
	move.X = x
	move.Y = y
	move.Speed = speed
}

func GetNotifyMsg(str string) *SC_NotifyMsg{
	var msgHeader json.MsgHeader
    msgHeader.MsgName = SC_NotifyMsgKey
	var msgContent SC_NotifyMsgContent
	msgContent.Msg = str
	return &SC_NotifyMsg{
		MsgHeader:msgHeader,
		MsgContent:msgContent,
	}
}

func GetIsFiredMsg() *SC_PlayerIsFired{
	var msgHeader json.MsgHeader
    msgHeader.MsgName = SC_PlayerIsFiredKey
	return &SC_PlayerIsFired{
		MsgHeader:msgHeader,
	}
}

func GetInWaitRoomStateMsg(state datastruct.WaitRoomState,r_id string) *SC_PlayerInWaitRoom{
	var msgHeader json.MsgHeader
    msgHeader.MsgName = SC_PlayerInWaitRoomKey
    var msgContent SC_PlayerInWaitRoomContent
	msgContent.RoomID = r_id
	msgContent.State = state
    return &SC_PlayerInWaitRoom{
		MsgHeader:msgHeader,
		MsgContent:msgContent,
	}
}

func GetInWaitRoomMsg(state datastruct.WaitRoomState,r_id string,isMaster int,players []datastruct.PlayerInWaitRoom) *SC_PlayerInWaitRoom{
	var msgHeader json.MsgHeader
    msgHeader.MsgName = SC_PlayerInWaitRoomKey
    var msgContent SC_PlayerInWaitRoomContent
	msgContent.RoomID = r_id
	msgContent.State = state
	msgContent.IsMaster = isMaster
    msgContent.Players = players
    return &SC_PlayerInWaitRoom{
		MsgHeader:msgHeader,
		MsgContent:msgContent,
	}
}

func GetMatchingEndMsg(r_id string) *SC_PlayerMatchingEnd{
	var msgHeader json.MsgHeader
    msgHeader.MsgName = SC_PlayerMatchingEndKey
    var msgContent SC_PlayerMatchingEndContent
    msgContent.RoomID =r_id
    return &SC_PlayerMatchingEnd{
		MsgHeader:msgHeader,
		MsgContent:msgContent,
	}
}

func GetReMatchMsg() *SC_PlayerReMatch{
	var msgHeader json.MsgHeader
    msgHeader.MsgName = "SC_PlayerReMatch"
    return &SC_PlayerReMatch{
		MsgHeader:msgHeader,
	}
}

func GetJoinInvalidMsg() *SC_PlayerJoinInvalid{
	var msgHeader json.MsgHeader
    msgHeader.MsgName = "SC_PlayerJoinInvalid"
    return &SC_PlayerJoinInvalid{
		MsgHeader:msgHeader,
	}
}

func GetInitRoomDataMsg(content SC_InitRoomDataContent) *SC_InitRoomData{
	var msgHeader json.MsgHeader
    msgHeader.MsgName = "SC_InitRoomData"
    return &SC_InitRoomData{
		MsgHeader:msgHeader,
		MsgContent:content,
	}
}

func GetRoomFrameDataMsg(content *SC_RoomFrameDataContent) *SC_RoomFrameData{
	var msgHeader json.MsgHeader
    msgHeader.MsgName = "SC_RoomFrameData"
    return &SC_RoomFrameData{
		MsgHeader:msgHeader,
		MsgContent:content,
	}
}

func GetGameOverMsg(content *SC_GameOverDataContent)*SC_GameOverData{
	var msgHeader json.MsgHeader
    msgHeader.MsgName = "SC_GameOverData"
    return &SC_GameOverData{
		MsgHeader:msgHeader,
		MsgContent:content,
	} 
}

func GetPower(e_type EnergyPointType) int {
	 power:=0
	 switch e_type{
	   case TypeA:
		 power=10
	   case TypeB:
		 power=20
	   case TypeC:
		 power=40
	   case TypeD: 
	 }
	 return power
}





