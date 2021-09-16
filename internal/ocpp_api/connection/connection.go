package connection

import (
	"errors"
	"fmt"
	"log"
	"sync"

	msg "github.com/Beep-Technologies/beepbeep3-ocpp/internal/ocpp_api/messaging"
	ops "github.com/Beep-Technologies/beepbeep3-ocpp/internal/ocpp_api/operations"
	"github.com/gorilla/websocket"
)

type Connection struct {
	cpId   string
	conn   *websocket.Conn
	state  *connectionOperation
	logger *log.Logger
}

type connectionOperation struct {
	m sync.Mutex
	v ops.OCPP16OperationType
}

func NewConnection(cpId string, conn *websocket.Conn, l *log.Logger) *Connection {
	return &Connection{
		cpId: cpId,
		conn: conn,
		state: &connectionOperation{
			v: ops.INITIAL,
		},
		logger: l,
	}
}

func (c *Connection) GetCurrentOperation() ops.OCPP16OperationType {
	c.state.m.Lock()
	defer c.state.m.Unlock()

	return c.state.v
}

func (c *Connection) SetCurrentOperation(o ops.OCPP16OperationType) error {
	c.state.m.Lock()
	defer c.state.m.Unlock()

	if c.state.v != ops.NONE {
		return ErrorOngoingOperation
	}

	c.state.v = o

	return nil
}

func (c *Connection) ServeOCPP16() error {
	// the first operation should be a BootNotification
	// initiated by the Charge Point
	c.SetCurrentOperation(ops.BOOT_NOTIFICATION)
	_, p, err := c.conn.ReadMessage()
	if err != nil {
		return err
	}

	c.logger.Printf("[CALL: FROM %s] %s", c.cpId, p)

	req, err := msg.ParseOCPP16Call(p)
	if err != nil {
		return err
	}

	if req.Action != "BootNotification" {
		return ErrorFirstRequest
	}

	res, resErr := c.HandleRequest(req)
	if resErr != nil {
		resErrJson, err := msg.UnparseOCPP16CallError(resErr)
		if err != nil {
			return err
		}

		c.logger.Printf("[CALLERROR: TO %s] %s", c.cpId, resErrJson)
		c.conn.WriteMessage(websocket.TextMessage, resErrJson)

		errorMsg := fmt.Sprintf("%+v", resErr)
		return errors.New(errorMsg)
	}

	resJson, err := msg.UnparseOCPP16CallResult(res)
	if err != nil {
		return err
	}

	c.logger.Printf("[CALLRESULT: TO %s] %s", c.cpId, resJson)
	c.conn.WriteMessage(websocket.TextMessage, resJson)

	return nil
}
