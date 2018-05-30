package Matching

import (
	"server/datastruct"
	"sync"
	"time"
	"github.com/name5566/leaf/gate"
	"server/db"
	"server/tools"
	"server/msg"
)

/*单人匹配*/
type SingleMatch struct {
	Times time.Duration
	MaxWaitTime time.Duration
	Pool_Capacity int
	ticker *time.Ticker
	isExistTicker bool
	rooms *Rooms
	onlinePlayers *datastruct.OnlinePlayers
	singleMatchPool *SingleMatchingPool
	actionPool *MatchActionPool
}

func NewSingleMatch()*SingleMatch{
	singleMatch:=new(SingleMatch)
	singleMatch.init()
	return singleMatch
}

func (match *SingleMatch)init(){
	match.isExistTicker = false
	match.Times = 1*time.Second //定时器多少时间执行一次
	match.MaxWaitTime = 5*time.Second//玩家最大等待时间多少秒
	match.Pool_Capacity = 10 //满足有多少个人就开始游戏
	match.onlinePlayers = datastruct.NewOnlinePlayers()
	match.singleMatchPool = newSingleMatchingPool(match.Pool_Capacity)
	match.actionPool = newMatchActionPool(match.Pool_Capacity)
	match.rooms = NewRooms()
}

func (match *SingleMatch)addPlayer(connUUID string,a gate.Agent,uid int){
	match.addOnlinePlayer(connUUID,a,uid)
	match.actionPool.AddInMatchActionPool(connUUID)
}

func (match *SingleMatch)RemovePlayer(connUUID string){
	match.onlinePlayers.Delete(connUUID)
	match.actionPool.RemoveFromMatchActionPool(connUUID)
}

func (match *SingleMatch)addOnlinePlayer(connUUID string,a gate.Agent,uid int){
	match.onlinePlayers.Lock.Lock()
	 defer match.onlinePlayers.Lock.Unlock()
	 v, ok := match.onlinePlayers.Bm[connUUID];
	 if !ok {
		 user:=db.Module.GetUserInfo(uid)
		 player:=datastruct.CreatePlayer(user)
		 player.Agent = a
		 match.onlinePlayers.Bm[connUUID]=*player
	 }else{ 
		 v.GameData.EnterType = datastruct.EmptyWay
		 v.GameData.RoomId = datastruct.NULLSTRING
	 }
}
func (match *SingleMatch)CheckActionPool(connUUID string) bool{
	  return match.actionPool.Check(connUUID)
}

func (match *SingleMatch)Matching(connUUID string, a gate.Agent,uid int){
	  match.addPlayer(connUUID,a,uid)
	  willEnterRoom:=false
	//willEnterRoom 是否将要加入了房间
	//r_id,willEnterRoom:=rooms.GetFreeRoomId()

	if !willEnterRoom{
	   match.singleMatchPool.Mutex.Lock()
	   defer match.singleMatchPool.Mutex.Unlock()
	   num:=len(match.singleMatchPool.Pool)
	   LeastPeople:=match.Pool_Capacity
	   if num<LeastPeople{
		match.singleMatchPool.Pool=append(match.singleMatchPool.Pool,connUUID)
		match.createTicker()
		if num == LeastPeople-1{
			//check player is online or offline
			//offline player is removed from pool
			//if all online create room
			removeIndex,_:=match.getOfflinePlayers()
			rm_num:=len(removeIndex)
			if rm_num<=0{//池中没有离线玩家,则创建房间
				match.cleanPoolAndCreateRoom()
			}else{
				match.removeOfflinePlayersInPool(removeIndex)
			}
		}
	   }
	}else{
		// player,tf:=onlinePlayers.GetAndUpdateState(p_uuid,datastruct.FreeRoom,r_id)
		// if tf{
		// 	player.Agent.WriteMsg(msg.GetMatchingEndMsg(r_id))
		// }
	}
}

func (match *SingleMatch)createTicker(){
    if !match.isExistTicker {
		match.isExistTicker = true
		match.ticker = time.NewTicker(match.Times)
        go match.selectTicker()
    } 
}

func (match *SingleMatch)stopTicker(){
    if match.ticker != nil{
	   match.ticker.Stop() 
       match.isExistTicker = false
    }
}

func (match *SingleMatch)selectTicker(){
     for {
        select {
         case <-match.ticker.C:
            match.computeMatchingTime()
        }
    }
}

func (match *SingleMatch)getOfflinePlayers() ([]int, map[string]datastruct.Player){
    tmp_map:=match.onlinePlayers.Items()
	LeastPeople:=match.Pool_Capacity
	
    online_map:=make(map[string]datastruct.Player)
    
    removeIndex:=make([]int,0,LeastPeople)
    
    var online_player datastruct.Player
    online_key:=datastruct.NULLSTRING
    
    for index,v := range match.singleMatchPool.Pool{
        isOnline:=false
        for key,player :=range tmp_map{
            if key == v{
                isOnline=true
                online_key = key
                online_player = player
                break
            }
        }
        if isOnline{
            online_map[online_key]= online_player
            delete(tmp_map, online_key)//移除对比过的数据,减少空间复杂度
        }else{
            removeIndex=append(removeIndex,index)//保存离线玩家
        }
    }
    return removeIndex,online_map
}

func (match *SingleMatch)removeOfflinePlayersInPool(removeIndex []int){
    rm_count:=0
    for index,v := range removeIndex {
        if index!=0{
           v = v-rm_count
        }
        match.singleMatchPool.Pool=append(match.singleMatchPool.Pool[:v], match.singleMatchPool.Pool[v+1:]...)
        rm_count++;
    }
}

func (match *SingleMatch)cleanPoolAndCreateRoom(){
	match.stopTicker()
    arr:=make([]string,len(match.singleMatchPool.Pool))
    copy(arr,match.singleMatchPool.Pool)
    match.singleMatchPool.Pool=match.singleMatchPool.Pool[:0]//clean pool
    go match.createMatchingTypeRoom(arr)
}

func (match *SingleMatch)createMatchingTypeRoom(playerUUID []string){
    r_uuid:=tools.UniqueId()
    players:=match.onlinePlayers.GetsAndUpdateState(playerUUID,datastruct.FromMatchingPool,r_uuid)
    room:=createRoom(playerUUID,Matching,r_uuid)
    rooms.Set(r_uuid,room)
    for _,play := range players{
        play.Agent.WriteMsg(msg.GetMatchingEndMsg(r_uuid))
    }
}


func (match *SingleMatch)computeMatchingTime(){
	match.singleMatchPool.Mutex.Lock()
    defer  match.singleMatchPool.Mutex.Unlock()
    num:=len(match.singleMatchPool.Pool)
    if num >0{
        removeIndex,online_map:=match.getOfflinePlayers()
        rm_num:=len(removeIndex)
        if rm_num>0{//删除池中离线玩家
			match.removeOfflinePlayersInPool(removeIndex)
        }
        now_t := time.Now()
        for _,player := range online_map{
            rs_sub:=now_t.Sub(player.GameData.StartMatchingTime)
            if rs_sub>=match.MaxWaitTime{
                match.cleanPoolAndCreateRoom()
                break
            }
        }
    }else{
	   match.stopTicker()
    }
}


/*单人匹配池*/
type SingleMatchingPool struct {
	 Mutex *sync.RWMutex //读写互斥量
	 Pool  []string //存放玩家uuid
}

func newSingleMatchingPool(poolCapacity int)*SingleMatchingPool{
	singleMatchingPool:=new(SingleMatchingPool)
	singleMatchingPool.init(poolCapacity)
	return singleMatchingPool
}

func (pool *SingleMatchingPool)init(poolCapacity int){
	  pool.Mutex = new(sync.RWMutex)
	  pool.Pool = make([]string,0,poolCapacity)
}













