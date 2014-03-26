package bencode

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
)

func Marshal(v interface{}) (data []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("panic: " + fmt.Sprint(r))
		}
	}()

	m := new(marshaller)
	m.buffer = bytes.NewBuffer(nil)
	m.marshallers = make(map[reflect.Kind]func(reflect.Value) error)
	m.marshallers[reflect.Bool] = m.marshalBool
	m.marshallers[reflect.Int] = m.marshalInt
	m.marshallers[reflect.Int8] = m.marshalInt
	m.marshallers[reflect.Int16] = m.marshalInt
	m.marshallers[reflect.Int32] = m.marshalInt
	m.marshallers[reflect.Int64] = m.marshalInt
	m.marshallers[reflect.Uint] = m.marshalUint
	m.marshallers[reflect.Uint8] = m.marshalUint
	m.marshallers[reflect.Uint16] = m.marshalUint
	m.marshallers[reflect.Uint32] = m.marshalUint
	m.marshallers[reflect.Uint64] = m.marshalUint
	m.marshallers[reflect.Array] = m.marshalArray
	m.marshallers[reflect.Map] = m.marshalMap
	m.marshallers[reflect.Ptr] = m.marshalPtr
	m.marshallers[reflect.Slice] = m.marshalArray
	m.marshallers[reflect.String] = m.marshalString
	m.marshallers[reflect.Struct] = m.marshalStruct

	if err := m.marshal(reflect.ValueOf(v)); err != nil {
		return nil, err
	}

	return m.buffer.Bytes(), nil
}

type marshaller struct {
	marshallers map[reflect.Kind]func(reflect.Value) error
	buffer      *bytes.Buffer
}

func (m *marshaller) marshal(v reflect.Value) error {
	if marshaller, ok := m.marshallers[v.Kind()]; ok {
		if err := marshaller(v); err != nil {
			return err
		}
	} else {
		return errors.New("unsupported data type: " + v.Kind().String())
	}

	return nil
}

func (m *marshaller) marshalPtr(v reflect.Value) error {
	return m.marshal(v.Elem())
}

func (m *marshaller) marshalMap(v reflect.Value) error {
	if v.Type().Key().Kind() != reflect.String {
		return errors.New("can't marshal non-string keyed maps")
	}

	m.buffer.WriteString("d")

	for _, k := range v.MapKeys() {
		if err := m.marshalString(k); err != nil {
			return err
		}
		if err := m.marshal(v.MapIndex(k)); err != nil {
			return err
		}
	}

	m.buffer.WriteString("e")

	return nil
}

func (m *marshaller) marshalStruct(v reflect.Value) error {
	m.buffer.WriteString("d")

	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		m.marshalString(reflect.ValueOf(field.Name))
		if err := m.marshal(v.Field(i)); err != nil {
			return err
		}
	}

	m.buffer.WriteString("e")

	return nil
}

func (m *marshaller) marshalString(v reflect.Value) error {
	str := v.String()
	if len(str) > 0 {
		m.buffer.WriteString(fmt.Sprintf("%d:%s", len(str), str))
	}

	return nil
}

func (m *marshaller) marshalArray(v reflect.Value) error {
	m.buffer.WriteString("l")

	for i := 0; i < v.Len(); i++ {
		if err := m.marshal(v.Index(i)); err != nil {
			return err
		}
	}

	m.buffer.WriteString("e")

	return nil
}

func (m *marshaller) marshalUint(v reflect.Value) error {
	m.buffer.WriteString(fmt.Sprintf("i%de", v.Uint()))
	return nil
}

func (m *marshaller) marshalInt(v reflect.Value) error {
	m.buffer.WriteString(fmt.Sprintf("i%de", v.Int()))
	return nil
}

func (m *marshaller) marshalBool(v reflect.Value) error {
	if v.Bool() {
		m.buffer.Write([]byte("i1e"))
	} else {
		m.buffer.Write([]byte("i0e"))
	}

	return nil
}
