package protocol

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"strconv"
	"strings"
)

type RespType byte

const (
	SimpleStrings RespType = '+'
	Errors        RespType = '-'
	Integers      RespType = ':'
	BulkStrings   RespType = '$'
	Arrays        RespType = '*'
)

type Command struct {
	Type    RespType
	RawByte []byte
	Value   []interface{}
}

func (c *Command) Receiver(reader *bufio.Reader) ([]byte, error) {
	var rawData bytes.Buffer

	maxFor := 1
	for i := 0; i < maxFor; i++ {
		data, err := reader.ReadBytes('\n')
		if err != nil {
			return nil, err
		}
		rawData.Write(data)
		if i == 0 {
			c.Type = RespType(data[0])
		}
		switch RespType(data[0]) {
		case Arrays:
			arrCount, err := strconv.Atoi(string(data[1 : len(data)-2]))
			if err != nil {
				return nil, err
			}
			maxFor += arrCount
		case SimpleStrings:
			data, err := reader.ReadBytes('\n')
			if err != nil {
				return nil, err
			}
			rawData.Write(data)
		case Errors:
			data, err := reader.ReadBytes('\n')
			if err != nil {
				return nil, err
			}
			rawData.Write(data)
		case BulkStrings:
			textCount, err := strconv.Atoi(string(data[1 : len(data)-2]))
			if err != nil {
				return nil, err
			}
			data = make([]byte, textCount+2)
			if _, err := io.ReadFull(reader, data); err != nil {
				return nil, err
			}
			rawData.Write(data)
			c.Value = append(c.Value, string(data[:len(data)-2]))
		case Integers:
			data, err := reader.ReadBytes('\n')
			if err != nil {
				return nil, err
			}
			rawData.Write(data)
			if value, err := strconv.Atoi(string(data[:len(data)-2])); err == nil {
				c.Value = append(c.Value, value)
			} else {
				return nil, err
			}

		default:
			return nil, errors.New("解析错误")
		}
	}
	c.RawByte = rawData.Bytes()
	return c.RawByte, nil
}

func (c Command) Text(text string) []byte {
	if strings.Index(text, "\r\n") > 0 {
		return c.BulkStrings(text)
	} else {
		return c.SimpleStrings(text)
	}
}

func (c Command) SimpleStrings(text string) []byte {
	var result bytes.Buffer
	result.WriteByte(byte(SimpleStrings))
	result.Write([]byte(text))
	result.Write([]byte("\r\n"))
	return result.Bytes()
}

func (c Command) BulkStrings(text string) []byte {
	var result bytes.Buffer
	result.WriteByte(byte(BulkStrings))
	result.Write([]byte(strconv.Itoa(len(text))))
	result.Write([]byte("\r\n"))
	result.Write([]byte(text))
	result.Write([]byte("\r\n"))
	return result.Bytes()
}

func (c Command) Error(text string) []byte {
	var result bytes.Buffer
	result.WriteByte(byte(Errors))
	result.Write([]byte(text))
	result.Write([]byte("\r\n"))
	return result.Bytes()
}

func (c Command) Array(data []interface{}) []byte {
	var result bytes.Buffer
	arrLen := len(data)
	result.WriteByte(byte(Arrays))
	result.Write([]byte(strconv.Itoa(arrLen)))
	result.Write([]byte("\r\n"))
	for _, i := range data {
		switch value := i.(type) {
		case int:
			result.Write(c.Integers(value))
		case string:
			result.Write(c.Text(value))
		case []byte:
			result.WriteByte(byte(BulkStrings))
			result.Write([]byte(strconv.Itoa(len(value))))
			result.Write([]byte("\r\n"))
			result.Write(value)
			result.Write([]byte("\r\n"))
		}
	}
	return result.Bytes()
}

func (c Command) Integers(value int) []byte {
	var result bytes.Buffer
	result.WriteByte(byte(Integers))
	result.Write([]byte(strconv.Itoa(value)))
	result.Write([]byte("\r\n"))
	return result.Bytes()
}
