package main

import "encoding/json"

// Request -- запрос клиента к серверу.
type Request struct {
	// Поле Command может принимать три значения:
	// * "quit" - прощание с сервером (после этого сервер рвёт соединение);
	// * "add" - передача новых точек на сервер;
	// * "sq" - просьба посчитать площадь круга.
	Command string `json:"command"`

	// Если Command == "add", в поле Data должна лежать 2 точки в виде структуры TwoPoints
	// В противном случае, поле Data пустое.
	Data *json.RawMessage `json:"data"`
}

// Response -- ответ сервера клиенту.
type Response struct {
	// Поле Status может принимать три значения:
	// * "ok" - успешное выполнение команды "quit" или "add";
	// * "failed" - в процессе выполнения команды произошла ошибка;
	// * "result" - площадь.
	Status string `json:"status"`

	// Если Status == "failed", то в поле Data находится сообщение об ошибке.
	// Если Status == "result", в поле Data должно лежать число float64
	// В противном случае, поле Data пустое.
	Data *json.RawMessage `json:"data"`
}

type Point struct {
	//координата х (в десятичной системе, разрешён знак)
	CordX string `json:"x"`

	//координата х (в десятичной системе, разрешён знак)
	CordY string `json:"y"`
}

type TwoPoints struct {
	PointO Point
	Point1 Point
}
