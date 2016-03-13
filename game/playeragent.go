package game
import (
	"encoding/json"
	"chess/com"
	"github.com/lkj01010/log"
	"chess/dao"
)

type playerAgent struct {
	connId   int
	id       string
	c        chan string
	sendFunc func(int, string)

	toGame chan string		// 发往game的channel

	info com.UserInfo
}

func NewPlayerAgent(connId int, sf func(string)) *playerAgent {
	return &playerAgent{
		connId: connId,
		c: make(chan string, 10),
		sendFunc: sf,
		toGame: make(chan string, 10),
	}
}

func (p *playerAgent)Go() {
	// get user info
	var reply dao.User_InfoReply
	modelInst.dao.Call("User.GetInfo", &dao.Args{Id: p.id}, &reply)
	p.info = reply.Info

	for {
		select {
		case rcv := <-p.c:
			p.handle(rcv)
		}
	}
}

func (p *playerAgent)Receive(msg string) {
	p.c <- msg
}

func (p *playerAgent)Send(msg string) {
	p.sendFunc(p.connId, msg)
}

func (p *playerAgent)handle(msg string) (err error) {
	var msg com.Msg
	if err = json.Unmarshal([]byte(msg.Content), &msg); err != nil {
		return
	}

	switch msg.Cmd {
	case Cmd_Game_EnterReq:
		err = p.handleEnterReq(msg.Content)
	default:
		p.toGame <-msg.Content
	}
	return
}

func (p *playerAgent)handleEnterReq(content string) (err error) {
	_, e := modelInst.playerGames[p.connId]
	if e == nil {
		// 已经在游戏中，报错
		err = com.ErrAlreadyInGame
		log.Warning("game:model:handle:player enter req, err=", err.Error())
		return
	} else {
		var content EnterGameReq
		if err = json.Unmarshal([]byte(content), &content); err != nil {
			return
		}
		p.id = content.Id

		//判断钱是否够
		isCoinEnough := (p.info.Coin >= RoomEnterCoin[content.RoomType])

		// 够入场费
		if isCoinEnough {
			// 塞进房间
			game := modelInst.GetFreeGameByType(content.RoomType)
			game.PlayerEnter(content.Id, p.info, p.toGame, p.Send)
		}else {
			// 不够入场费
			// 返回response
			resp := com.MakeMsgString(Cmd_Game_EnterResp, com.E_CoinNotEnough, nil)
			p.Send(resp)
			return
		}
	}
	return
}