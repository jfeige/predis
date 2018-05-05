package predis

import (
	"errors"
	"strconv"
	"bytes"
	"fmt"
	"net"
	"bufio"
	"time"
)


type Conn struct {
	address string
	conn    net.Conn
	dst     []byte
	br 		*bufio.Reader
	bw 		*bufio.Writer
	pending int
	t time.Time
}


func (c *Conn) Send(command string, args ...interface{})error {

	c.pending += 1
	c.writeCommand(command,args)

	return c.bw.Flush()
}

func (c *Conn) Cmd(command string,args ...interface{})(replay interface{},err error){

	pending := c.pending
	c.pending = 0

	c.writeCommand(command,args)
	if err := c.bw.Flush();err != nil{
		return nil,err
	}

	for i := 0; i <= pending;i++{
		//解析返回值
		replay,err = c.readReply()
		if err != nil{
			return nil,err
		}
	}
	return replay,nil
}


func (c *Conn)readReply()(replay interface{},err error){

	line,err := c.readLine()
	if err != nil{
		return nil,err
	}
	if len(line) == 0{
		return nil,errors.New("response line is nil")
	}

	switch line[0] {
	case '+':
		switch {
		case len(line) == 3 && line[1] == 'O' && line[2] == 'K':
			return "OK", nil
		case len(line) == 5 && line[1] == 'P' && line[2] == 'O' && line[3] == 'N' && line[4] == 'G':
			return "PONG", nil
		default:
			return line[1:], nil
		}
	case '-':
		return errors.New(string(line[1:])), nil
	case ':':
		return line[1:],nil
	case '$': //$5\r\nlifei\r\n
		length, err := strconv.Atoi(string(line[1:]))
		if length < 0 || err != nil {
			return nil, err
		}
		tmp_line, err := c.readLine()
		if err != nil {
			return nil, err
		}
		if len(tmp_line) == 0 {
			return nil, ErrNil
		}
		return tmp_line, nil
	case '*':
		length, err := strconv.Atoi(string(line[1]))
		if length < 0 || err != nil {
			return nil, err
		}
		ret := make([]interface{}, length)
		for i := range ret{
			replay, err := c.readReply()
			if err != nil {
				return nil, err
			}
			ret[i] = replay
		}
		return ret, nil
	}
	return nil, errors.New("unexpected response line")
}

func (c *Conn)readLine()([]byte,error){
	line,err := c.br.ReadSlice('\n')
	if err != nil{
		return nil,err
	}
	i := len(line)-2
	if i < 0 || line[i] != '\r'{
		return nil,errors.New("bad response line terminator")
	}
	return line[:i],nil
}




func (c *Conn) writeCommand(command string,args []interface{}){
	c.writeLen(args)
	c.writeString(command)

	for _, arg := range args {
		switch arg := arg.(type) {
		case string:
			c.writeString(arg)
		case []byte:
			c.writeBytes(arg)
		case int:
			c.writeInt(arg)
		case int64:
			c.writeInt64(arg)
		case float64:
			c.writeFloat64(arg)
		default:
			var buf bytes.Buffer
			fmt.Fprint(&buf, arg)
			c.writeBytes(buf.Bytes())
		}
	}
}

func (c *Conn) writeLen(args []interface{}){
	c.bw.WriteString("*" + (strconv.Itoa(1+len(args))))
	c.bw.WriteString(CRLF)
}



func (c *Conn) writeInt(arg int) {
	c.writeBytes(strconv.AppendInt(c.dst[:0], int64(arg), 10))
}

func (c *Conn) writeFloat64(arg float64) {
	c.writeBytes(strconv.AppendFloat(c.dst[:0], arg, 'g', -1, 64))
}

func (c *Conn) writeInt64(arg int64) {
	c.writeBytes(strconv.AppendInt(c.dst[:0], arg, 10))
}

func (c *Conn) writeString(arg string) {
	c.bw.WriteString("$" + strconv.Itoa(len(arg)))
	c.bw.WriteString(CRLF)
	c.bw.WriteString(arg)
	c.bw.WriteString(CRLF)
}

func (c *Conn) writeBytes(arg []byte) {
	c.writeString(string(arg))
}


func (c *Conn) Flush()error{
	return c.bw.Flush()
}


