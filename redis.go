package predis

import (
	"strconv"
	"net"
	"bufio"
	"time"
)

var (
	CRLF = "\r\n"
)


func Dial(network,address string,pwd ...string)(*Conn,error){
	conn,err := net.Dial(network,address)
	if err != nil{
		return nil,err
	}
	c := &Conn{
		conn :conn,
		dst:make([]byte,0),
		br:bufio.NewReader(conn),
		bw:bufio.NewWriter(conn),
		t:time.Now(),
	}
	if len(pwd) > 0 && pwd[0] != ""{
		err = c.Send("AUTH",pwd[0])
		if err != nil{
			return nil,err
		}
	}
	return c,nil
}


func Int(replay interface{},err error)(int,error){
	if err != nil{
		return 0,err
	}
	switch replay := replay.(type) {
	case []byte:
		return strconv.Atoi(string(replay))
	case nil:
		return 0,ErrNil
	case error:
		return 0,replay
	default:
		return 0,checkErrResponse(replay)
	}
}

func String(replay interface{},err error)(string,error){
	if err != nil{
		return "",err
	}
	switch replay := replay.(type) {
	case []byte:
		return string(replay),nil
	case string:
		return replay,nil
	case nil:
		return "",ErrNil
	case error:
		return "",replay
	default:
		return "",checkErrResponse(replay)
	}
}

func Bool(replay interface{},err error)(bool,error){
	if err != nil{
		return false,err
	}
	switch replay := replay.(type) {
	case []byte:
		return strconv.ParseBool(string(replay))
	case int64:
		return replay != 0,nil
	case nil:
		return false,ErrNil
	case error:
		return false,replay
	default:
		return false,checkErrResponse(replay)
	}
}

func StringMap(replay interface{},err error)(map[string]string,error){
	if err != nil{
		return nil,err
	}
	switch replay := replay.(type) {
	case []interface{}:
		result := make(map[string]string,0)
		for i := 0; i < len(replay);i+=2{
			if replay[i] == nil{
				continue
			}
			key := string(replay[i].([]byte))
			value := string(replay[i+1].([]byte))
			result[key] = value
		}
		return result,nil
	case nil:
		return nil,ErrNil
	case error:
		return nil,replay
	default:
		return nil,checkErrResponse(replay)
	}
}

func Strings(replay interface{},err error)([]string,error){
	if err != nil{
		return nil,err
	}
	switch replay := replay.(type) {
	case []interface{}:
		result := make([]string,len(replay))
		for i := range replay{
			if replay[i] == nil{
				continue
			}
			v := replay[i].([]byte)
			result[i] = string(v)
		}
		return result,nil
	case nil:
		return nil,ErrNil
	case error:
		return nil,replay
	default:
		return nil,checkErrResponse(replay)
	}
}