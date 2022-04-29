package main

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/lemon-mint/lvm2"
	"github.com/lemon-mint/lvm2/asm"
	"github.com/lemon-mint/lvm2/binf"
)

type Instruction struct {
	Name     string     `parser:"@Ident"`
	Operands []*Operand `parser:"(@@ ','?)*"`
}

type Register struct {
	Name string `parser:"\"%\" @Ident"`
}

type Sign bool

func (b *Sign) Capture(values []string) error {
	*b = values[0] == "-"
	return nil
}

type Operand struct {
	Register *Register `parser:"  @@ |"`
	Sign     Sign      `parser:" @('-' | '+')?"`
	Int      *int64    `parser:" @Int"`
	String   *string   `parser:"| @String"`
	Variable *string   `parser:"| \"@\"@Ident"`
}

type File struct {
	Instructions []*Instruction `parser:"@@*"`
}

func main() {
	flags := map[string]string{}
	for i := 1; i < len(os.Args); i++ {
		if strings.HasPrefix(os.Args[i], "--") {
			if i+1 < len(os.Args) && !strings.HasPrefix(os.Args[i+1], "--") && !strings.HasPrefix(os.Args[i+1], "-") {
				flags[os.Args[i][2:]] = os.Args[i+1]
				i++
			} else {
				flags[os.Args[i][2:]] = ""
			}
		} else if strings.HasPrefix(os.Args[i], "-") {
			if i+1 < len(os.Args) && !strings.HasPrefix(os.Args[i+1], "--") && !strings.HasPrefix(os.Args[i+1], "-") {
				flags[os.Args[i][1:]] = os.Args[i+1]
				i++
			} else {
				flags[os.Args[i][1:]] = ""
			}
		} else {
			flags["__INPUT__"] = os.Args[i]
		}
	}

	var err error

	if v, ok := flags["__INPUT__"]; !ok || v == "" {
		log.Fatalln("No input file specified")
	}

	flags["__INPUT__"] = filepath.Clean(flags["__INPUT__"])
	flags["__INPUT__"], err = filepath.Abs(flags["__INPUT__"])
	if err != nil {
		log.Fatalln("Failed to get absolute path of input file:", err)
	}

	if v, ok := flags["o"]; !ok || v == "" {
		flags["o"] = filepath.Join(filepath.Dir(flags["__INPUT__"]), filepath.Base(flags["__INPUT__"])+".clvm2")
	} else {
		flags["o"], err = filepath.Abs(flags["o"])
		if err != nil {
			log.Fatalln("Failed to get absolute path of output file:", err)
		}
	}

	f, err := os.Open(flags["__INPUT__"])
	if err != nil {
		log.Fatalln("Failed to open input file:", err)
	}
	defer f.Close()

	parser, err := participle.Build(&File{})
	if err != nil {
		log.Fatalln("Failed to build parser:", err)
	}

	var file File
	err = parser.Parse(flags["__INPUT__"], f, &file)
	if err != nil {
		log.Fatalln("Failed to parse input file:", err)
	}

	e := asm.NewEncoder()

	var variables = map[string]uint64{}

	//repr.Println(file)

	var entryPoint uint64

	var ops []asm.Operand
	for _, instr := range file.Instructions {
		ops = ops[:0]
		instr.Name = strings.ToUpper(instr.Name)
		//fmt.Printf("%s", instr.Name)
		for i, operand := range instr.Operands {
			if operand.Register != nil {
				//fmt.Printf(" %%%s", operand.Register.Name)
				if id, ok := lvm2.Registers[operand.Register.Name]; ok {
					if i > 0 {
						ops = append(ops, asm.OPREG(id))
					} else {
						ops = append(ops, asm.OPCONST(id))
					}
				} else {
					log.Fatalln("Unknown register:", operand.Register.Name)
				}
			} else if operand.Int != nil {
				if operand.Sign {
					*operand.Int = -*operand.Int
				}
				//fmt.Printf(" %d", *operand.Int)
				ops = append(ops, asm.OPCONST(uint64(*operand.Int)))
			} else if operand.String != nil {
				*operand.String, err = strconv.Unquote(*operand.String)
				if err != nil {
					log.Fatalln("Failed to unquote string:", err)
				}
				//fmt.Printf(" %v", []byte(*operand.String))

				if instr.Name != "DATA" {
					log.Fatalln("string operand only allowed in DATA instruction")
				}
			} else if operand.Variable != nil {
				//fmt.Printf(" @%s", *operand.Variable)
				if instr.Name == "DATA" || instr.Name == "LABEL" {
					continue
				}
				if offset, ok := variables[*operand.Variable]; ok {
					ops = append(ops, asm.OPCONST(offset))
				}
			}
		}
		//fmt.Printf("\n")

		switch instr.Name {
		case "DATA":
			if len(instr.Operands) != 2 {
				log.Fatalln("DATA instruction must have exactly two operands")
			}
			if instr.Operands[0].Variable == nil {
				log.Fatalln("DATA instruction's first operand must be a variable")
			}
			if instr.Operands[1].String == nil {
				log.Fatalln("DATA instruction's second operand must be a string")
			}

			varName := *instr.Operands[0].Variable
			data := []byte(*instr.Operands[1].String)
			offset := e.Encode(asm.DATA(data))
			variables[varName] = offset
		case "LABEL":
			if len(instr.Operands) != 1 {
				log.Fatalln("LABEL instruction must have exactly one operand")
			}
			if instr.Operands[0].Variable == nil {
				log.Fatalln("LABEL instruction's operand must be a variable")
			}

			varName := *instr.Operands[0].Variable
			e.Encode(asm.INST(lvm2.InstructionType_NOP))
			offset := e.Encode(asm.LABEL(varName))
			variables[varName] = offset
		default:
			opcode, ok := lvm2.Instructions[instr.Name]
			if !ok {
				log.Fatalln("Invalid instruction:", instr.Name)
			}
			pc := e.Encode(asm.INST(opcode, ops...))
			if entryPoint == 0 {
				entryPoint = pc
			}
		}
	}

	e = asm.NewEncoder()

	//repr.Println(file)

	entryPoint = 0

	for _, instr := range file.Instructions {
		ops = ops[:0]
		instr.Name = strings.ToUpper(instr.Name)
		//fmt.Printf("%s", instr.Name)
		for i, operand := range instr.Operands {
			if operand.Register != nil {
				//fmt.Printf(" %%%s", operand.Register.Name)
				if id, ok := lvm2.Registers[operand.Register.Name]; ok {
					if i > 0 {
						ops = append(ops, asm.OPREG(id))
					} else {
						ops = append(ops, asm.OPCONST(id))
					}
				} else {
					log.Fatalln("Unknown register:", operand.Register.Name)
				}
			} else if operand.Int != nil {
				if operand.Sign {
					*operand.Int = -*operand.Int
				}
				//fmt.Printf(" %d", *operand.Int)
				ops = append(ops, asm.OPCONST(uint64(*operand.Int)))
			} else if operand.Variable != nil {
				//fmt.Printf(" @%s", *operand.Variable)
				if instr.Name == "DATA" || instr.Name == "LABEL" {
					continue
				}
				if offset, ok := variables[*operand.Variable]; ok {
					ops = append(ops, asm.OPCONST(offset))
				} else {
					log.Fatalln("Undefined variable:", *operand.Variable)
				}
			}
		}
		//fmt.Printf("\n")

		switch instr.Name {
		case "DATA":
			if len(instr.Operands) != 2 {
				log.Fatalln("DATA instruction must have exactly two operands")
			}
			if instr.Operands[0].Variable == nil {
				log.Fatalln("DATA instruction's first operand must be a variable")
			}
			if instr.Operands[1].String == nil {
				log.Fatalln("DATA instruction's second operand must be a string")
			}

			varName := *instr.Operands[0].Variable
			data := []byte(*instr.Operands[1].String)
			offset := e.Encode(asm.DATA(data))
			variables[varName] = offset
		case "LABEL":
			if len(instr.Operands) != 1 {
				log.Fatalln("LABEL instruction must have exactly one operand")
			}
			if instr.Operands[0].Variable == nil {
				log.Fatalln("LABEL instruction's operand must be a variable")
			}

			varName := *instr.Operands[0].Variable
			e.Encode(asm.INST(lvm2.InstructionType_NOP))
			offset := e.Encode(asm.LABEL(varName))
			variables[varName] = offset
		default:
			opcode, ok := lvm2.Instructions[instr.Name]
			if !ok {
				log.Fatalln("Invalid instruction:", instr.Name)
			}
			pc := e.Encode(asm.INST(opcode, ops...))
			if entryPoint == 0 {
				entryPoint = pc
			}
		}
	}

	//fmt.Println(e.Bytes())

	if pc, ok := variables["ENTRYPOINT"]; ok {
		entryPoint = pc
	}

	// vm := lvm2.VM{
	// 	Memory: lvm2.NewMemory(),
	// 	Files: map[uint64]lvm2.VMFile{
	// 		0: os.Stdin,
	// 		1: os.Stdout,
	// 		2: os.Stderr,
	// 	},
	// 	FileCounter: 3,
	// }
	// vm.Memory.SetProgram(e.Bytes())

	// // Set Initial Stack Pointer
	// vm.Registers[lvm2.REGISTER_SP] = vm.Memory.MaxAddress
	// vm.Registers[lvm2.REGISTER_SB] = vm.Memory.MaxAddress

	// // Set Entry Point
	// vm.SetProgramCounter(entryPoint)
	// ret, err := vm.Run()
	// if err != nil {
	// 	panic(err)
	// }
	// os.Exit(int(ret))

	prog := binf.New_Program(binf.EncodingType_RAW, binf.New_Header(0x01, entryPoint), e.Bytes())

	of, err := os.Create(flags["o"])
	if err != nil {
		log.Fatalln("Failed to open output file:", err)
	}
	defer of.Close()

	_, err = of.Write(prog)
	if err != nil {
		log.Fatalln("Failed to write output file:", err)
	}

}
