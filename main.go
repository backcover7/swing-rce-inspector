package main

import (
	"bytes"
	"fmt"
	"github.com/4ra1n/swing-rce-inspector/asm"
	"github.com/4ra1n/swing-rce-inspector/classfile"
	"github.com/4ra1n/swing-rce-inspector/common"
	"github.com/4ra1n/swing-rce-inspector/files"
	"github.com/4ra1n/swing-rce-inspector/global"
	"github.com/4ra1n/swing-rce-inspector/inherit"
	"github.com/4ra1n/swing-rce-inspector/log"
	"os"
	"os/signal"
	"runtime"
	"strings"
)

// analysis data
var (
	discoveryClass     []*asm.Class
	discoveryClassFile []*savedClassFile
	classMap           map[string]*asm.Class
)

// conditions
var (
	IsComponentSubclass bool
	HasNoParamInit      bool
	HasStringSetMethod  bool
)

// save analysis data
type (
	savedClassFile struct {
		ClassName string
		ClassFile *classfile.ClassFile
		ClassRef  *asm.Class
	}
	resultData struct {
		StringData   string
		InstListData []string
	}
)

// startAnalysis is the core function of project
// cannot use goroutine because of constant pool
func startAnalysis(saved *savedClassFile, resultChan chan *resultData) {
	IsComponentSubclass = false
	HasNoParamInit = false
	HasStringSetMethod = false
	cf := saved.ClassFile
	cl := saved.ClassRef
	// is subclass of component
	if inherit.IsSubclassOf(cl.Name(), global.ComponentName) {
		IsComponentSubclass = true
	}
	// refresh constant pool
	// IMPORTANT
	global.CP = cf.ConstantPool()
	for _, method := range cf.Methods() {
		// has no parameter init method
		if IsComponentSubclass && method.Name() == global.InitMethod {
			if method.Descriptor() == global.NoParamVoidMethod {
				HasNoParamInit = true
				// OK
				continue
			}
		}
		// has field with set method
		// only one string parameter
		var continueFlag bool
		if strings.HasPrefix(method.Name(), global.SetPrefix) {
			// setField -> Field
			temp := method.Name()[3:]
			var flag bool
			for _, field := range cl.Fields() {
				if strings.ToLower(field.Name()) == strings.ToLower(temp) {
					flag = true
					break
				}
			}
			// IMPORTANT: refresh flag
			HasStringSetMethod = false
			if flag && method.Descriptor() == global.StringParamSetVoidMethod {
				HasStringSetMethod = true
			}
			// check conditions
			if IsComponentSubclass && HasStringSetMethod && HasNoParamInit {
				continueFlag = true
			}
		}
		if !continueFlag {
			continue
		}
		var invokeFlag bool

		// **************************************************
		// ************ Analysis JVM Instruction ************
		// **************************************************
		codeAttr := method.CodeAttribute()
		if codeAttr == nil {
			// interface or abstract
			continue
		}
		bytecode := codeAttr.Code()
		// virtual program counter
		var pc int
		reader := &common.BytecodeReader{}
		// save all instructions to struct
		instSet := &common.InstructionSet{}
		instSet.ClassName = cl.Name()
		instSet.MethodName = method.Name()
		instSet.Desc = method.Descriptor()

		var instList []string
		// simple taint analysis
		var taint bool

		for {
			// read finish
			if pc >= len(bytecode) {
				break
			}
			// offset
			reader.Reset(bytecode, pc)
			// read instruction
			opcode := reader.ReadUint8()
			inst := asm.NewInstruction(opcode)
			// read operands of the instruction
			inst.FetchOperands(reader)
			ops := inst.GetOperands()
			instEntry := common.InstructionEntry{
				Instrument: getInstructionName(inst),
				Operands:   ops,
			}

			// (1) set(param) init -> param in LOCAL VARIABLE ARRAYS[1]
			// (2) load arrays[1] -> param on top of stack
			// (3) INVOKE ANY -> POP -> taint
			//
			// |   LOCAL VARIABLES   |  | OPERAND STACK TOP |
			// [ this | string param ]  [       OTHERS      ]
			// AFTER LOAD INST
			// [ this | string param ]->[   string param    ]
			// AFTER INVOKE ANY INST
			// [ this | string param ]  [       OTHERS      ] -> POP
			setLoadInst := &asm.ALOAD_1{}
			if instEntry.Instrument == getInstructionName(setLoadInst) {
				taint = true
			}

			instSet.InstArray = append(instSet.InstArray, instEntry)
			// offset++
			// read next
			pc = reader.PC()
			// INVOKE ANY
			if strings.HasPrefix(instEntry.Instrument, global.Invoke) {
				invokeFlag = true
				// do not show desc info
				temp := strings.Split(ops[0], " ")[0]
				// now top of stack is taint
				if taint {
					temp = temp + " (taint)"
					// clean
					taint = false
				}
				instList = append(instList, temp)
			}
		}
		if invokeFlag {
			fmt.Println(instSet.ClassName, instSet.MethodName)
			for _, i := range instList {
				fmt.Println("->", i)
			}
			fmt.Println()

			s := fmt.Sprintf("%s %s", instSet.ClassName, instSet.MethodName)
			data := &resultData{
				StringData:   s,
				InstListData: instList,
			}
			resultChan <- data
			// clean cache
			instList = nil
		}
	}
}

func startDiscovery(class string) {
	data, err := os.ReadFile(class)
	cf, err := classfile.Parse(data)
	if err != nil {
		log.Error(err.Error())
		os.Exit(-1)
	}
	cl := asm.NewClass(cf)
	s := &savedClassFile{
		ClassName: cl.Name(),
		ClassFile: cf,
		ClassRef:  cl,
	}
	discoveryClassFile = append(discoveryClassFile, s)
	classMap[s.ClassName] = s.ClassRef
	discoveryClass = append(discoveryClass, s.ClassRef)
}

func getInstructionName(instruction asm.Instruction) string {
	// type name -> instruction name
	i := fmt.Sprintf("%T", instruction)
	return strings.Split(i, ".")[1]
}

func collect(resultChan chan *resultData, finishChan chan bool) {
	buf := bytes.Buffer{}
	for {
		select {
		case r := <-resultChan:
			buf.Write([]byte(r.StringData))
			buf.Write([]byte("\n"))
			for _, v := range r.InstListData {
				buf.Write([]byte("->"))
				buf.Write([]byte(v))
				buf.Write([]byte("\n"))
			}
			buf.Write([]byte("\n"))
			break
		case <-finishChan:
			goto WRITE
		}
	}
WRITE:
	os.WriteFile("result.txt", buf.Bytes(), 0644)
	close(finishChan)
	close(resultChan)
}

func handle() {
	if err := recover(); err != nil {
		message := fmt.Sprintf("%s", err)
		var pcs [32]uintptr
		n := runtime.Callers(3, pcs[:])
		var str strings.Builder
		temp := fmt.Sprintf("%s\nTraceback:", message)
		str.WriteString(temp)
		for _, pc := range pcs[:n] {
			fn := runtime.FuncForPC(pc)
			file, line := fn.FileLine(pc)
			s := fmt.Sprintf("\n\t%s:%d", file, line)
			str.WriteString(s)
		}
		log.Error("%s\n\n", str.String())
		os.Exit(-1)
	}
}

// swing-rce-inspector
// author: 4ra1n of Chaitin Tech
// create: 2022/10/02
// update: 2022/10/10
// line: 2663
func main() {
	// set log level
	log.SetLevel(log.InfoLevel)
	// recover panic
	defer handle()
	log.Info("start swing-rce-inspector")
	log.Info("delete last analysis data")
	// delete last data
	files.RemoveTempFiles()
	log.Info("unzip jar files")
	files.UnzipJars("jars")
	// class file path array
	classes := files.ReadAllClasses()
	// (1) discovery
	// all class -> discoveryClass
	log.Info("start discovery")
	classMap = make(map[string]*asm.Class)
	for _, c := range classes {
		startDiscovery(c)
	}
	log.Info("finish discovery")
	// (2) IMPORTANT: build inheritance
	// class A extends B
	// class B extends C
	// class C implements D,E
	// A is subclass of B,C,D,E
	// B is subclass of C,D,E
	// C is subclass of D,E
	log.Info("build inherit data")
	inherit.Init(discoveryClass, classMap)
	// (3) start look up
	//   <1> must be a subclass of java/awt/Component
	//   <2> must have a construction with no parameters
	//   <3> must have a field with set method
	//   <4> the set method has only one string parameters
	//   <5> have method invoke in the set method
	//   <6> simple taint analysis
	log.Info("start analysis")
	resultChan := make(chan *resultData)
	finishChan := make(chan bool, 1)
	go collect(resultChan, finishChan)
	for _, i := range discoveryClassFile {
		// IMPORTANT
		// do not use goroutine
		// because of static ConstantPool
		startAnalysis(i, resultChan)
	}
	// make sure task finish
	finishChan <- true
	// delete temp data
	log.Info("delete temp data")
	files.RemoveTempFiles()
	log.Info("finish")
	log.Info("press ctrl+c exit")
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill)
	for {
		s := <-c
		fmt.Println(s)
		os.Exit(0)
	}
}
