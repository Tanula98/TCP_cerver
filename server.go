package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
)

import (
	"log"
	"os"
	"strconv"
)

// Client - состояние клиента.
type Client struct {
	logger log.Logger    // Объект для печати логов
	conn   *net.TCPConn  // Объект TCP-соединения
	enc    *json.Encoder // Объект для кодирования и отправки сообщений
	S      float64       // Площадь

}

// NewClient - конструктор клиента, принимает в качестве параметра
// объект TCP-соединения.
func NewClient(conn *net.TCPConn) *Client {
	return &Client{
		logger: *log.New(os.Stdout, fmt.Sprintf("client %s", conn.RemoteAddr().String()), 0),
		conn:   conn,
		enc:    json.NewEncoder(conn),
		S:      float64(0.0),
	}
}

// serve - метод, в котором реализован цикл взаимодействия с клиентом.
// Подразумевается, что метод serve будет вызаваться в отдельной go-программе.
func (client *Client) serve() {
	defer client.conn.Close()
	decoder := json.NewDecoder(client.conn)
	for {
		var req Request
		if err := decoder.Decode(&req); err != nil {
			client.logger.Print("cannot decode message ", "reason ", err)
			break
		} else {
			client.logger.Print("received command ", "command ", req.Command)
			if client.handleRequest(&req) {
				client.logger.Print("shutting down connection ")
				break
			}
		}
	}
}

// handleRequest - метод обработки запроса от клиента. Он возвращает true,
// если клиент передал команду "quit" и хочет завершить общение.
func (client *Client) handleRequest(req *Request) bool {
	switch req.Command {
	case "quit":
		client.respond("ok", nil)
		return true
	case "add":
		errorMsg := ""
		if req.Data == nil {
			errorMsg = "data field is absent"
		} else {
			var O TwoPoints
			if err := json.Unmarshal(*req.Data, &O); err != nil {
				errorMsg = "malformed data field"
			} else {
				if _, ee1 := strconv.Atoi(O.PointO.CordX);ee1!= nil {
					client.respond("faled", "Not digit")
					errorMsg = "malformed data field"
					return false
				}
				a,_ :=strconv.Atoi(O.PointO.CordX)

				if _, ee2 := strconv.Atoi(O.PointO.CordY);ee2!= nil {
					client.respond("faled", "Not digit")
					errorMsg = "malformed data field"
					return false
				}
				b,_ :=strconv.Atoi(O.PointO.CordY)
				if _, ee3 := strconv.Atoi(O.Point1.CordX);ee3!= nil {
					client.respond("faled", "Not digit")
					errorMsg = "malformed data field"
					return false
				}
				c,_ :=strconv.Atoi(O.Point1.CordX)
				if _, ee4 := strconv.Atoi(O.Point1.CordY);ee4!= nil {
					client.respond("faled", "Not digit")
					errorMsg = "malformed data field"
					return false
				}
				d,_ :=strconv.Atoi(O.Point1.CordY)

				client.S = 3.14 * float64((a-c)*(a-c) + (b-d)*(b-d))
			}
		}
		if errorMsg == "" {
			client.respond("ok", nil)
		} else {
			client.logger.Print("addition failed", "reason", errorMsg)
			client.respond("failed", errorMsg)
		}

	case "sq":
		if client.S == 0 {
			client.respond("failed", "division by zero")
		} else {

			client.respond("result", client.S)
		}

	default:
		client.logger.Print("unknown command")
		client.respond("failed", "unknown command")
	}
	return false
}

// respond - вспомогательный метод для передачи ответа с указанным статусом
// и данными. Данные могут быть пустыми (data == nil).
func (client *Client) respond(status string, data interface{}) {
	var raw json.RawMessage
	raw, _ = json.Marshal(data)
	client.enc.Encode(&Response{status, &raw})
}

func main() {
	// Работа с командной строкой, в которой может указываться необязательный ключ -addr.
	var addrStr string
	flag.StringVar(&addrStr, "addr", "127.0.0.1:6000", "specify ip address and port")
	flag.Parse()

	// Разбор адреса, строковое представление которого находится в переменной addrStr.
	if addr, err := net.ResolveTCPAddr("tcp", addrStr); err != nil {
		log.Print("address resolution failed address ", addrStr)
	} else {
		log.Print("resolved TCP address ", addr.String())

		// Инициация слушания сети на заданном адресе.
		if listener, err := net.ListenTCP("tcp", addr); err != nil {
			log.Print("listening failed reason ", err)
		} else {
			// Цикл приёма входящих соединений.
			for {
				if conn, err := listener.AcceptTCP(); err != nil {
					log.Print("cannot accept connection reason ", err)
				} else {
					log.Print("accepted connection address ", conn.RemoteAddr().String())

					// Запуск go-программы для обслуживания клиентов.
					go NewClient(conn).serve()
				}
			}
		}
	}
}
