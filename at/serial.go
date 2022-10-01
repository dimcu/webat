package at

import (
	"errors"
	"github.com/tarm/serial"
	"log"
	"time"
)

type Serial struct {
	cof  *serial.Config
	conn *serial.Port

	writeCh      chan []byte
	endCh        chan bool
	writeState   bool
	readState    bool
	onRead       OnRead
	onConnect    OnConnect
	onDisconnect OnDisconnect
}
type OnRead func([]byte)
type OnConnect func()
type OnDisconnect func()

func OpenAtSerial(name string, baud int) (*Serial, error) {

	atSerial := Serial{cof: &serial.Config{Name: name, Baud: baud}}

	go atSerial.run()

	return &atSerial, nil

}

func (at *Serial) open() error {

	//打开串口
	conn, err := serial.OpenPort(at.cof)
	if err != nil {
		return err
	}
	at.conn = conn
	return nil
}
func (at *Serial) close() error {

	if at.conn != nil {
		return at.conn.Close()
	}
	return nil
}

func (at *Serial) run() {
	for {
		at.writeCh = make(chan []byte, 100)
		at.endCh = make(chan bool, 1)

		err := at.open()
		if err != nil {
			//log.Println(err)
			time.Sleep(time.Second)
			continue
		}

		if at.onConnect != nil {
			at.onConnect()
		}

		go at.writeTask()

		at.SendAtCmd("ATE0\r\n")

		//阻塞读取
		at.readTask()
		//读取失败 通知写入结束
		at.writeEnd()

		if at.onDisconnect != nil {
			at.onDisconnect()
		}
		time.Sleep(time.Second)
	}

}

// 通知写入线程关闭
func (at *Serial) writeEnd() {
	at.close()

	if !at.writeState {
		return
	}

	select {
	case <-time.After(time.Second):
		break
	case at.endCh <- true:
		break
	}

}

// 读取线程
func (at *Serial) readTask() {
	log.Println("readTask run")
	at.readState = true
	defer func() {
		log.Println("readTask end")
		at.readState = false

	}()

	for {
		buf := make([]byte, 2000)
		n, err := at.conn.Read(buf)
		if err != nil {
			log.Println("readTask err", err)
			return
		}
		//log.Println(string(buf[:n]))

		log.Println("recv:<-", string(buf))

		if at.onRead != nil {
			at.onRead(buf[:n])
		}
	}
}

// 写入线程
func (at *Serial) writeTask() {

	log.Println("writeTask run")
	at.writeState = true
	defer func() {
		log.Println("writeTask end")
		at.writeState = false
	}()

	for {
		select {
		case <-at.endCh:
			return
		case buf := <-at.writeCh:
			log.Println("send:->", string(buf))
			_, err := at.conn.Write(buf)
			if err != nil {
				log.Println("writeTask err", err)
				return
			}
			time.Sleep(time.Millisecond * 100)
		}
	}
}

func (at *Serial) SendAtCmd(atStr string) error {

	select {
	case <-time.After(time.Second):
		return errors.New("time out")
	case at.writeCh <- []byte(atStr + "\r\n"):
		return nil
	}
}

func (at *Serial) SetCallback(read OnRead, connect OnConnect, disconnect OnDisconnect) {
	at.onRead = read
	at.onConnect = connect
	at.onDisconnect = disconnect
}
