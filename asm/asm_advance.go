package asm

import (
	"fmt"
	"github.com/4ra1n/swing-rce-inspector/classfile"
	"github.com/4ra1n/swing-rce-inspector/common"
	"github.com/4ra1n/swing-rce-inspector/global"
	"github.com/4ra1n/swing-rce-inspector/log"
	"os"
	"reflect"
	"strconv"
)

// *********************** GOTO_W ***********************

type GOTO_W struct {
	offset int
}

func (self *GOTO_W) FetchOperands(reader *common.BytecodeReader) {
	self.offset = int(reader.ReadInt32())
}
func (self *GOTO_W) GetOperands() []string {
	ret := make([]string, 1)
	ret[0] = "[not support]"
	return ret
}

// *********************** IINC ***********************

type IINC struct {
	Index uint
	Const int32
}

func (self *IINC) FetchOperands(reader *common.BytecodeReader) {
	self.Index = uint(reader.ReadUint8())
	self.Const = int32(reader.ReadInt8())
}

func (self *IINC) GetOperands() []string {
	ret := make([]string, 2)
	ret[0] = strconv.Itoa(int(self.Index))
	ret[1] = strconv.Itoa(int(self.Const))
	return ret
}

// *********************** INVOKEDYNAMIC ***********************

type INVOKEDYNAMIC struct{}

func (self INVOKEDYNAMIC) FetchOperands(reader *common.BytecodeReader) {
	reader.ReadInt8()
	reader.ReadInt8()
	reader.ReadInt8()
	reader.ReadInt8()
}

func (self *INVOKEDYNAMIC) GetOperands() []string {
	ret := make([]string, 1)
	ret[0] = "[not support]"
	return ret
}

// *********************** INVOKEINTERFACE ***********************

type INVOKEINTERFACE struct {
	index uint
}

func (self *INVOKEINTERFACE) FetchOperands(reader *common.BytecodeReader) {
	self.index = uint(reader.ReadUint16())
	reader.ReadUint8()
	reader.ReadUint8()
}

func (self *INVOKEINTERFACE) GetOperands() []string {
	name := global.CP.GetConstantInfo(uint16(self.index))
	typeName := reflect.TypeOf(name).String()

	var (
		className  string
		methodName string
		desc       string
	)

	switch typeName {
	case "*classfile.ConstantInterfaceMethodRefInfo":
		className = name.(*classfile.ConstantInterfaceMethodRefInfo).ClassName()
		methodName, desc = name.(*classfile.ConstantInterfaceMethodRefInfo).NameAndDescriptor()
	case "*classfile.ConstantMethodRefInfo":
		className = name.(*classfile.ConstantMethodRefInfo).ClassName()
		methodName, desc = name.(*classfile.ConstantMethodRefInfo).NameAndDescriptor()
	default:
		log.Error("error")
		os.Exit(-1)
	}

	ret := make([]string, 1)
	out := fmt.Sprintf("%s.%s %s", className, methodName, desc)
	ret[0] = out
	return ret
}

// *********************** INVOKESPECIAL ***********************

type INVOKESPECIAL struct{ Index16Instruction }

func (self *INVOKESPECIAL) GetOperands() []string {
	name := global.CP.GetConstantInfo(uint16(self.Index))
	typeName := reflect.TypeOf(name).String()

	var (
		className  string
		methodName string
		desc       string
	)

	switch typeName {
	case "*classfile.ConstantInterfaceMethodRefInfo":
		className = name.(*classfile.ConstantInterfaceMethodRefInfo).ClassName()
		methodName, desc = name.(*classfile.ConstantInterfaceMethodRefInfo).NameAndDescriptor()
	case "*classfile.ConstantMethodRefInfo":
		className = name.(*classfile.ConstantMethodRefInfo).ClassName()
		methodName, desc = name.(*classfile.ConstantMethodRefInfo).NameAndDescriptor()
	default:
		log.Error("error")
		os.Exit(-1)
	}

	ret := make([]string, 1)
	out := fmt.Sprintf("%s.%s %s", className, methodName, desc)
	ret[0] = out
	return ret
}

// *********************** INVOKESTATIC ***********************

type INVOKESTATIC struct{ Index16Instruction }

func (self *INVOKESTATIC) GetOperands() []string {
	name := global.CP.GetConstantInfo(uint16(self.Index))
	typeName := reflect.TypeOf(name).String()

	var (
		className  string
		methodName string
		desc       string
	)

	switch typeName {
	case "*classfile.ConstantInterfaceMethodRefInfo":
		className = name.(*classfile.ConstantInterfaceMethodRefInfo).ClassName()
		methodName, desc = name.(*classfile.ConstantInterfaceMethodRefInfo).NameAndDescriptor()
	case "*classfile.ConstantMethodRefInfo":
		className = name.(*classfile.ConstantMethodRefInfo).ClassName()
		methodName, desc = name.(*classfile.ConstantMethodRefInfo).NameAndDescriptor()
	default:
		log.Error("error")
		os.Exit(-1)
	}

	ret := make([]string, 1)
	out := fmt.Sprintf("%s.%s %s", className, methodName, desc)
	ret[0] = out
	return ret
}

// *********************** INVOKEVIRTUAL ***********************

type INVOKEVIRTUAL struct{ Index16Instruction }

func (self *INVOKEVIRTUAL) GetOperands() []string {
	name := global.CP.GetConstantInfo(uint16(self.Index))
	typeName := reflect.TypeOf(name).String()

	var (
		className  string
		methodName string
		desc       string
	)

	switch typeName {
	case "*classfile.ConstantInterfaceMethodRefInfo":
		className = name.(*classfile.ConstantInterfaceMethodRefInfo).ClassName()
		methodName, desc = name.(*classfile.ConstantInterfaceMethodRefInfo).NameAndDescriptor()
	case "*classfile.ConstantMethodRefInfo":
		className = name.(*classfile.ConstantMethodRefInfo).ClassName()
		methodName, desc = name.(*classfile.ConstantMethodRefInfo).NameAndDescriptor()
	default:
		log.Error("error")
		os.Exit(-1)
	}

	ret := make([]string, 1)
	out := fmt.Sprintf("%s.%s %s", className, methodName, desc)
	ret[0] = out
	return ret
}

// *********************** BIPUSH ***********************

type BIPUSH struct {
	val int8
}

func (self *BIPUSH) FetchOperands(reader *common.BytecodeReader) {
	self.val = reader.ReadInt8()
}

func (self *BIPUSH) GetOperands() []string {
	ret := make([]string, 1)
	ret[0] = strconv.Itoa(int(self.val))
	return ret
}

// *********************** SIPUSH ***********************

type SIPUSH struct {
	val int16
}

func (self *SIPUSH) FetchOperands(reader *common.BytecodeReader) {
	self.val = reader.ReadInt16()
}

func (self *SIPUSH) GetOperands() []string {
	ret := make([]string, 1)
	ret[0] = strconv.Itoa(int(self.val))
	return ret
}

// *********************** JSR_W ***********************

type JSR_W struct{}

func (J JSR_W) FetchOperands(reader *common.BytecodeReader) {
	reader.ReadUint8()
	reader.ReadUint8()
	reader.ReadUint8()
	reader.ReadUint8()
}

func (self *JSR_W) GetOperands() []string {
	ret := make([]string, 1)
	ret[0] = "[not support]"
	return ret
}

// *********************** LOOKUPSWITCH ***********************

type LOOKUPSWITCH struct {
	defaultOffset int32
	npairs        int32
	matchOffsets  []int32
}

func (self *LOOKUPSWITCH) FetchOperands(reader *common.BytecodeReader) {
	reader.SkipPadding()
	self.defaultOffset = reader.ReadInt32()
	self.npairs = reader.ReadInt32()
	self.matchOffsets = reader.ReadInt32s(self.npairs * 2)
}

func (self *LOOKUPSWITCH) GetOperands() []string {
	ret := make([]string, 1)
	ret[0] = "[not support]"
	return ret
}

// *********************** MULTIANEWARRAY ***********************

type MULTIANEWARRAY struct {
	index      uint16
	dimensions uint8
}

func (self *MULTIANEWARRAY) FetchOperands(reader *common.BytecodeReader) {
	self.index = reader.ReadUint16()
	self.dimensions = reader.ReadUint8()
}

func (self *MULTIANEWARRAY) GetOperands() []string {
	ret := make([]string, 1)
	ret[0] = "[not support]"
	return ret
}

// *********************** NEWARRAY ***********************

/*
T_BOOLEAN	4
T_CHAR		5
T_FLOAT		6
T_DOUBLE	7
T_BYTE		8
T_SHORT		9
T_INT		10
T_LONG		11
*/

type NEWARRAY struct {
	arrayType uint8
}

func (self *NEWARRAY) FetchOperands(reader *common.BytecodeReader) {
	self.arrayType = reader.ReadUint8()
}

func (self *NEWARRAY) GetOperands() []string {
	ret := make([]string, 1)
	ret[0] = string(self.arrayType)
	return ret
}

// *********************** TABLESWITCH ***********************

type TABLESWITCH struct {
	defaultOffset int32
	low           int32
	high          int32
	jumpOffsets   []int32
}

func (self *TABLESWITCH) FetchOperands(reader *common.BytecodeReader) {
	reader.SkipPadding()
	self.defaultOffset = reader.ReadInt32()
	self.low = reader.ReadInt32()
	self.high = reader.ReadInt32()
	jumpOffsetsCount := self.high - self.low + 1
	self.jumpOffsets = reader.ReadInt32s(jumpOffsetsCount)
}

func (self *TABLESWITCH) GetOperands() []string {
	ret := make([]string, 1)
	ret[0] = "[not support]"
	return ret
}

// *********************** WIDE ***********************

type WIDE struct {
	modifiedInstruction Instruction
}

func (self *WIDE) FetchOperands(reader *common.BytecodeReader) {
	opcode := reader.ReadUint8()
	switch opcode {
	case 0x15:
		inst := &ILOAD{}
		inst.Index = uint(reader.ReadUint16())
		self.modifiedInstruction = inst
	case 0x16:
		inst := &LLOAD{}
		inst.Index = uint(reader.ReadUint16())
		self.modifiedInstruction = inst
	case 0x17:
		inst := &FLOAD{}
		inst.Index = uint(reader.ReadUint16())
		self.modifiedInstruction = inst
	case 0x18:
		inst := &DLOAD{}
		inst.Index = uint(reader.ReadUint16())
		self.modifiedInstruction = inst
	case 0x19:
		inst := &ALOAD{}
		inst.Index = uint(reader.ReadUint16())
		self.modifiedInstruction = inst
	case 0x36:
		inst := &ISTORE{}
		inst.Index = uint(reader.ReadUint16())
		self.modifiedInstruction = inst
	case 0x37:
		inst := &LSTORE{}
		inst.Index = uint(reader.ReadUint16())
		self.modifiedInstruction = inst
	case 0x38:
		inst := &FSTORE{}
		inst.Index = uint(reader.ReadUint16())
		self.modifiedInstruction = inst
	case 0x39:
		inst := &DSTORE{}
		inst.Index = uint(reader.ReadUint16())
		self.modifiedInstruction = inst
	case 0x3a:
		inst := &ASTORE{}
		inst.Index = uint(reader.ReadUint16())
		self.modifiedInstruction = inst
	case 0x84:
		inst := &IINC{}
		inst.Index = uint(reader.ReadUint16())
		inst.Const = int32(reader.ReadInt16())
		self.modifiedInstruction = inst
	case 0xa9:
		log.Error("unsupported opcode: 0xa9")
		os.Exit(-1)
	}
}

func (self *WIDE) GetOperands() []string {
	ret := make([]string, 1)
	ret[0] = "[not support]"
	return ret
}
