package bencode

import (
	"bytes"
	"errors"
	"reflect"
	"strconv"
	"strings"
	"fmt"
)

const (
	I     = byte('i')
	E     = byte('e')
	L     = byte('l')
	D     = byte('d')
	COLON = byte(':')
)

func Marshal(v interface{}) ([]byte, error) {
	return nil, nil
}

func Unmarshal(data []byte, v interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("panic: " + fmt.Sprint(r))
		}
	}()
	s := new(scanner)
	s.data = data
	s.position = 0
	s.unmarshallers = make(map[byte]func(reflect.Value) error)
	s.unmarshallers[I] = s.unmarshalNumber
	s.unmarshallers[L] = s.unmarshalArray
	s.unmarshallers[D] = s.unmarshalObject
	s.unmarshallers[0] = s.unmarshalString

	return s.unmarshal(v)
}

type scanner struct {
	data          []byte
	position      uint64
	unmarshallers map[byte]func(reflect.Value) error
}

func (s *scanner) unmarshal(v interface{}) error {
	rv := reflect.ValueOf(v)

	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New("invalid type: " + reflect.TypeOf(v).String())
	}

	indirect := reflect.Indirect(rv)
	return s.unmarshalValue(indirect)
}

func (s *scanner) unmarshalValue(indirect reflect.Value) error {
	if indirect.Kind() == reflect.Ptr {
		indirect.Set(reflect.New(indirect.Type().Elem()))
		indirect = reflect.Indirect(indirect)
	}

	var unmarshaller func(reflect.Value) error
	mark := s.data[s.position]
	if _, ok := s.unmarshallers[mark]; ok {
		unmarshaller = s.unmarshallers[mark]
	} else {
		unmarshaller = s.unmarshallers[0]
	}

	return unmarshaller(indirect)
}

func (s *scanner) unmarshalObject(indirect reflect.Value) error {
	switch indirect.Kind() {
	case reflect.Struct:
		return s.unmarshalStruct(indirect)
	case reflect.Map:
		return s.unmarshalMap(indirect)
	}

	return errors.New("invalid type: " + indirect.Kind().String())
}

func (s *scanner) unmarshalMap(indirect reflect.Value) error {
	s.position++

	if indirect.Kind() != reflect.Map {
		return errors.New("invalid type: " + indirect.Kind().String())
	}

	if indirect.Type().Key().Kind() != reflect.String {
		return errors.New("invalid key type: " + indirect.Type().Key().Kind().String())
	}

	indirect.Set(reflect.MakeMap(indirect.Type()))

	for s.data[s.position] != E {
		key := reflect.New(indirect.Type().Key()).Elem()
		val := reflect.New(indirect.Type().Elem()).Elem()
		if err := s.unmarshalString(key); err != nil {
			return err
		}
		if err := s.unmarshalValue(val); err != nil {
			return err
		}
		indirect.SetMapIndex(key, val)
	}

	s.position++

	return nil
}

func (s *scanner) unmarshalStruct(indirect reflect.Value) error {
	s.position++

	if indirect.Kind() != reflect.Struct {
		return errors.New("invalid type: " + indirect.Kind().String())
	}

	indirect.Set(reflect.New(indirect.Type()).Elem())

	for s.data[s.position] != E {
		key := reflect.New(reflect.TypeOf("")).Elem()
		if err := s.unmarshalString(key); err != nil {
			return err
		}

		field := indirect.FieldByNameFunc(func(f string) bool {
			k := key.String()
			k = strings.Replace(k, " ", "", -1)
			k = strings.Replace(k, "-", "", -1)
			return strings.ToLower(f) == strings.ToLower(k)
		})

		if !field.IsValid() {
			return errors.New("invalid struct key: " + key.String())
		}

		val := reflect.New(field.Type()).Elem()
		if err := s.unmarshalValue(val); err != nil {
			return err
		}
		field.Set(val)
	}

	s.position++

	return nil
}

func (s *scanner) unmarshalArray(indirect reflect.Value) error {
	s.position++

	switch indirect.Kind() {
	case reflect.Array:
		indirect.Set(reflect.New(indirect.Type()).Elem())
	case reflect.Slice:
		indirect.Set(reflect.MakeSlice(indirect.Type(), 0, 0))
	default:
		return errors.New("array or slice expected, got: " + indirect.Kind().String())
	}

	for i := 0; s.data[s.position] != E; i++ {
		if indirect.Cap() < i {
			newcap := indirect.Cap() + indirect.Cap()/2
			if newcap < 4 {
				newcap = 4
			}
			indirect.SetCap(newcap)
		}
		val := reflect.New(indirect.Type().Elem()).Elem()
		if err := s.unmarshalValue(val); err != nil {
			return err
		}
		indirect.Set(reflect.Append(indirect, val))
	}

	s.position++

	return nil
}

func (s *scanner) unmarshalString(indirect reflect.Value) error {
	chunk := s.scanWhile(COLON)
	uint, err := strconv.ParseUint(string(chunk), 10, 64)
	if err != nil {
		return err
	}
	str := string(s.data[s.position : s.position+uint])
	s.position += uint
	indirect.SetString(str)

	return nil
}

func (s *scanner) unmarshalNumber(indirect reflect.Value) error {
	s.position++
	chunk := s.scanWhile(E)
	switch indirect.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		int, err := strconv.ParseInt(string(chunk), 10, 64)
		if err != nil {
			return err
		}
		indirect.SetInt(int)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uint, err := strconv.ParseUint(string(chunk), 10, 64)
		if err != nil {
			return err
		}
		indirect.SetUint(uint)
	default:
		return errors.New("cannot unmarshal a number into: " + indirect.Kind().String())
	}

	return nil
}

func (s *scanner) scanWhile(b byte) []byte {
	buf := bytes.NewBuffer(nil)
	for ; s.data[s.position] != b; s.position++ {
		buf.WriteByte(s.data[s.position])
	}

	s.position++

	return buf.Bytes()
}
