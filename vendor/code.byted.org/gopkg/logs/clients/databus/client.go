package databus

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/golang/protobuf/proto"
)

// 1. 尽量复用DatabusCollector对象
// 2. 测试情况下BufferedDatabusCollector吞吐会比DatabusCollector吞吐低
// 3. 如果DatabusCollector对象太多会导致domain socket连接特别多

const (
	DEFAULT_SOCKET_PATH  = "/opt/tmp/sock/databus_collector.seqpacket.sock"
	PACKET_SIZE_LIMIT    = 60 * (1 << 10) // 60KB
	DEFAULT_TIMEOUT      = 100 * time.Millisecond
	ERROR_TOLERATE       = 50
	DEFAULT_MAX_CONN_NUM = 5
)

var ErrAlreadyClosed = errors.New("client already closed")

type conn struct {
	nc       *net.UnixConn
	config   *DatabusCollector
	errTimes int
}

func (cn *conn) Write(b []byte) error {
	// 发送失败次数达到阈值重建链接

	if cn.nc == nil || cn.errTimes > ERROR_TOLERATE {
		cn.errTimes = 0
		con, err := net.DialUnix("unixpacket", nil, &net.UnixAddr{cn.config.socketPath, "unixpacket"})
		if err != nil {
			return err
		}
		cn.nc = con
	}
	cn.nc.SetDeadline(time.Now().Add(cn.config.writeTimeout))
	_, err := cn.nc.Write(b)
	if err != nil {
		cn.errTimes += 1
		return err
	}
	return nil
}

func (cn *conn) release() {
	if cn.nc != nil {
		cn.nc.Close()
	}
}

type DatabusCollector struct {
	socketPath   string
	pool         chan *conn
	writeTimeout time.Duration
}

func NewDefaultCollector() *DatabusCollector {
	return NewCollector(DEFAULT_SOCKET_PATH, DEFAULT_TIMEOUT, DEFAULT_MAX_CONN_NUM)
}

func NewCollectorWithTimeout(timeout time.Duration) *DatabusCollector {
	return NewCollector(DEFAULT_SOCKET_PATH, timeout, DEFAULT_MAX_CONN_NUM)
}

func NewCollector(socketPath string, timeout time.Duration, maxConn int) *DatabusCollector {

	collector := &DatabusCollector{
		socketPath:   socketPath,
		writeTimeout: timeout,
		pool:         make(chan *conn, maxConn),
	}
	return collector
}

// 从pool中取出一个连接
func (this *DatabusCollector) borrow() (*conn, error) {
	var cl *conn
	var ok bool
	select {
	case cl, ok = <-this.pool:
		if !ok {
			return nil, ErrAlreadyClosed
		}
	default:
		cl = this.newConn()
		fmt.Println("pool empty create a databus connection")
	}
	return cl, nil
}

// 归还连接
func (this *DatabusCollector) putBack(cl *conn) {
	select {
	case this.pool <- cl:
	default:
		cl.release()
		fmt.Println("pool full destory a databus connection")
	}
}

func (this *DatabusCollector) newConn() *conn {
	con := &conn{
		nc:       nil,
		config:   this,
		errTimes: 0,
	}
	return con
}

func (this *DatabusCollector) CollectArray(channel string, messages []*ApplicationMessage) error {
	var payload RequestPayload
	payload.Channel = &channel
	payload.Messages = messages
	buf, err := proto.Marshal(&payload)
	if err != nil {
		return err
	}
	if len(buf) > PACKET_SIZE_LIMIT {
		return fmt.Errorf("CollectArray Err: Packet too large.")
	}
	con, err := this.borrow()
	if err != nil {
		return err
	}
	err = con.Write(buf)
	this.putBack(con)
	if err != nil {
		return err
	}
	return nil
}

func (this *DatabusCollector) Collect(channel string, value []byte, key []byte, codec int32) error {
	var message ApplicationMessage
	message.Key = key
	message.Value = value
	message.Codec = &codec
	var payload RequestPayload
	payload.Channel = &channel
	payload.Messages = []*ApplicationMessage{&message}
	buf, err := proto.Marshal(&payload)
	if err != nil {
		return err
	}
	con, err := this.borrow()
	if err != nil {
		return err
	}
	err = con.Write(buf)
	this.putBack(con)
	if err != nil {
		return err
	}
	return nil
}

// 内部实现，不要调用
func (this *DatabusCollector) SendProto(b []byte) error {
	con, err := this.borrow()
	if err != nil {
		return err
	}
	err = con.Write(b)
	this.putBack(con)
	if err != nil {
		return err
	}
	return nil
}

func (this *DatabusCollector) Close() error {
	//连接等待gc关闭, close queue 会触发panic

	return nil
}
