package match

import (
    "sync"
    "server/datastruct"
    "github.com/name5566/leaf/gate"
    "server/msg"
)


/*匹配动作池*/ //收到匹配消息的时候加入池，主动离开和自动离开在池中删除，完成匹配后，在池中删除
type MatchActionPool struct {
	Mutex *sync.RWMutex //读写互斥量
	Pool  []string //存放玩家uuid
}

func NewMatchActionPool(poolCapacity int)*MatchActionPool{
	matchActionPool:=new(MatchActionPool)
	matchActionPool.init(poolCapacity)
	return matchActionPool
}

func (actionPool *MatchActionPool)init(poolCapacity int){
	actionPool.Mutex = new(sync.RWMutex)
	if poolCapacity > 0{
		actionPool.Pool = make([]string,0,poolCapacity)
	}else{
		actionPool.Pool = make([]string,0)
	}
}

func (actionPool *MatchActionPool)RemoveFromMatchActionPool(p_uuid string){
    actionPool.Mutex.Lock()
    defer actionPool.Mutex.Unlock()
    rm_index:=-1
    for index,v := range actionPool.Pool{
        if v==p_uuid{
            rm_index = index
            break
        }
    }
    if rm_index>=0{
        actionPool.Pool=append(actionPool.Pool[:rm_index], actionPool.Pool[rm_index+1:]...)
    }
}

func (actionPool *MatchActionPool)AddInMatchActionPool(p_uuid string){
    actionPool.Mutex.Lock()
    actionPool.Pool=append(actionPool.Pool,p_uuid)
    actionPool.Mutex.Unlock()
}

func (actionPool *MatchActionPool)Check(p_uuid string) bool{
    tf:=false
    actionPool.Mutex.RLock()
    defer  actionPool.Mutex.RUnlock()
    for _,v:=range actionPool.Pool{
        if v==p_uuid{
            tf = true
            break
        }
    }
    return tf
}

type ParentMatch interface {
     RemoveRoomWithID(r_id string)
     GetOnlinePlayersPtr() *datastruct.OnlinePlayers
     Matching(connUUID string, a gate.Agent,uid int) string
     CheckActionPool(connUUID string) bool
     PlayerLeftRoom(r_id string,connUUID string)
     PlayersDied(r_id string,values []datastruct.PlayerDiedData)
     EnergyExpended(expended int,agentUserData datastruct.AgentUserData)
     PlayerMoved(r_id string,play_id int,moveData *msg.CS_MoveData)
     RemovePlayer(connUUID string)
     PlayerJoin(connUUID string,joinData *msg.CS_PlayerJoinRoom)
}


