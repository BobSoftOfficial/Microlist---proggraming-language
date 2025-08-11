package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Instruction types for the virtual machine
type OpCode byte

const (
	OP_CREATE_LIST OpCode = iota
	OP_SET_VALUE
	OP_SET_STRING
	OP_OUTPUT_LIST
	OP_REFERENCE_LIST
	OP_HALT
)

// Instruction represents a single VM instruction
type Instruction struct {
	OpCode   OpCode
	ListID   int
	Value    string
	RefID    int
}

// MicroList represents a single numbered list with 8-bit values
type MicroList struct {
	ID       int
	Values   []byte
	CanOutput bool // Only ml. lists can output
}

// VM represents the MicroList virtual machine
type VM struct {
	Lists        map[int]*MicroList
	Instructions []Instruction
	PC           int // Program Counter
}

// Compiler compiles MicroList source to bytecode
type Compiler struct {
	Source       string
	Instructions []Instruction
	Lines        []string
	CurrentLine  int
}

func NewCompiler(source string) *Compiler {
	return &Compiler{
		Source:       source,
		Instructions: make([]Instruction, 0),
		Lines:        strings.Split(strings.TrimSpace(source), "\n"),
		CurrentLine:  0,
	}
}

func (c *Compiler) Compile() error {
	for c.CurrentLine < len(c.Lines) {
		line := strings.TrimSpace(c.Lines[c.CurrentLine])
		
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "//") {
			c.CurrentLine++
			continue
		}
		
		if err := c.parseLine(line); err != nil {
			return fmt.Errorf("line %d: %v", c.CurrentLine+1, err)
		}
		
		c.CurrentLine++
	}
	
	// Add halt instruction
	c.Instructions = append(c.Instructions, Instruction{OpCode: OP_HALT})
	return nil
}

func (c *Compiler) parseLine(line string) error {
	// Check for numbered ml (e.g., "1 ml" or "2 ml.")
	parts := strings.Fields(line)
	if len(parts) >= 2 {
		if listNum, err := strconv.Atoi(parts[0]); err == nil {
			if parts[1] == "ml" {
				// Regular ml (can't output)
				c.Instructions = append(c.Instructions, Instruction{
					OpCode: OP_CREATE_LIST,
					ListID: listNum,
				})
				return nil
			} else if parts[1] == "ml." {
				// Output ml (can output)
				c.Instructions = append(c.Instructions, Instruction{
					OpCode: OP_OUTPUT_LIST,
					ListID: listNum,
				})
				return nil
			}
		}
	}
	
	// Check for reference (e.g., "1[]")
	if strings.HasSuffix(line, "[]") {
		refStr := strings.TrimSuffix(line, "[]")
		if refID, err := strconv.Atoi(refStr); err == nil {
			c.Instructions = append(c.Instructions, Instruction{
				OpCode: OP_REFERENCE_LIST,
				RefID:  refID,
			})
			return nil
		}
	}
	
	// Check for string literal
	if strings.HasPrefix(line, "\"") && strings.HasSuffix(line, "\"") {
		str := line[1 : len(line)-1] // Remove quotes
		c.Instructions = append(c.Instructions, Instruction{
			OpCode: OP_SET_STRING,
			Value:  str,
		})
		return nil
	}
	
	// Check for numeric value
	if num, err := strconv.Atoi(line); err == nil {
		if num < 0 || num > 255 {
			return fmt.Errorf("value %d out of range (0-255)", num)
		}
		c.Instructions = append(c.Instructions, Instruction{
			OpCode: OP_SET_VALUE,
			Value:  strconv.Itoa(num),
		})
		return nil
	}
	
	return fmt.Errorf("unknown instruction: %s", line)
}

func (c *Compiler) GetInstructions() []Instruction {
	return c.Instructions
}

// VM Implementation
func NewVM() *VM {
	return &VM{
		Lists:        make(map[int]*MicroList),
		Instructions: make([]Instruction, 0),
		PC:           0,
	}
}

func (vm *VM) LoadProgram(instructions []Instruction) {
	vm.Instructions = instructions
	vm.PC = 0
}

func (vm *VM) Execute() error {
	var currentList *MicroList
	
	for vm.PC < len(vm.Instructions) {
		instr := vm.Instructions[vm.PC]
		
		switch instr.OpCode {
		case OP_CREATE_LIST:
			currentList = &MicroList{
				ID:        instr.ListID,
				Values:    make([]byte, 0),
				CanOutput: false,
			}
			vm.Lists[instr.ListID] = currentList
			
		case OP_OUTPUT_LIST:
			currentList = &MicroList{
				ID:        instr.ListID,
				Values:    make([]byte, 0),
				CanOutput: true,
			}
			vm.Lists[instr.ListID] = currentList
			
		case OP_SET_VALUE:
			if currentList == nil {
				return fmt.Errorf("no current list to add value to")
			}
			value, err := strconv.Atoi(instr.Value)
			if err != nil {
				return fmt.Errorf("invalid value: %s", instr.Value)
			}
			currentList.Values = append(currentList.Values, byte(value))
			
		case OP_SET_STRING:
			if currentList == nil {
				return fmt.Errorf("no current list to add string to")
			}
			for _, char := range instr.Value {
				if char > 255 {
					return fmt.Errorf("character %c out of 8-bit range", char)
				}
				currentList.Values = append(currentList.Values, byte(char))
			}
			
		case OP_REFERENCE_LIST:
			if currentList == nil {
				return fmt.Errorf("no current list to add reference to")
			}
			
			refList, exists := vm.Lists[instr.RefID]
			if !exists {
				return fmt.Errorf("referenced list %d does not exist", instr.RefID)
			}
			
			// Clear current list and copy values from referenced list (replace, don't append)
			currentList.Values = make([]byte, len(refList.Values))
			copy(currentList.Values, refList.Values)
			
		case OP_HALT:
			// Output all lists that can output (ml.) - clean output only
			for id := 1; id <= len(vm.Lists)+10; id++ { // Check in order
				if list, exists := vm.Lists[id]; exists && list.CanOutput {
					for _, value := range list.Values {
						if value >= 32 && value <= 126 {
							// Printable ASCII
							fmt.Printf("%c", value)
						} else {
							// Non-printable, show as number in brackets
							fmt.Printf("[%d]", value)
						}
					}
					fmt.Println() // New line after each output list
				}
			}
			return nil
			
		default:
			return fmt.Errorf("unknown opcode: %d", instr.OpCode)
		}
		
		vm.PC++
	}
	
	return nil
}

// Bytecode serialization for binary output
func (vm *VM) SerializeToBinary() []byte {
	var binary []byte
	
	// Magic number for MicroList bytecode
	binary = append(binary, 'M', 'L', 'I', 'S', 'T') // MLIST
	
	// Version
	binary = append(binary, 1, 0)
	
	// Number of instructions
	instrCount := len(vm.Instructions)
	binary = append(binary, byte(instrCount>>24), byte(instrCount>>16), 
		byte(instrCount>>8), byte(instrCount&0xFF))
	
	// Serialize instructions
	for _, instr := range vm.Instructions {
		binary = append(binary, byte(instr.OpCode))
		binary = append(binary, byte(instr.ListID>>8), byte(instr.ListID&0xFF))
		binary = append(binary, byte(instr.RefID>>8), byte(instr.RefID&0xFF))
		
		// String length and data
		strLen := len(instr.Value)
		binary = append(binary, byte(strLen))
		binary = append(binary, []byte(instr.Value)...)
	}
	
	return binary
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("MicroList Compiler v2.0")
		fmt.Println("Usage:")
		fmt.Println("  microlist <source.ml>           - Compile and run")
		fmt.Println("  microlist -c <source.ml>        - Compile to binary")
		fmt.Println("  microlist -i                    - Interactive mode")
		fmt.Println("")
		fmt.Println("MicroList Syntax Examples:")
		fmt.Println("  1 ml          // Create list 1 (no output)")
		fmt.Println("  \"hello\"       // Add string to current list")
		fmt.Println("  42            // Add number to current list")
		fmt.Println("  2 ml.         // Create list 2 (with output)")
		fmt.Println("  \"world\"       // Add string to list 2")
		fmt.Println("  3 ml.         // Create list 3 (with output)")
		fmt.Println("  1[]           // Reference list 1's content")
		fmt.Println("")
		fmt.Println("Only 'ml.' lists produce output, 'ml' lists are storage only.")
		return
	}
	
	// Interactive mode
	if os.Args[1] == "-i" {
		runInteractiveMode()
		return
	}
	
	// Compile mode
	if os.Args[1] == "-c" && len(os.Args) > 2 {
		compileToBinary(os.Args[2])
		return
	}
	
	// Run mode
	sourceFile := os.Args[1]
	content, err := os.ReadFile(sourceFile)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}
	
	// Compile
	compiler := NewCompiler(string(content))
	if err := compiler.Compile(); err != nil {
		fmt.Printf("Compilation error: %v\n", err)
		return
	}
	
	// Execute
	vm := NewVM()
	vm.LoadProgram(compiler.GetInstructions())
	
	if err := vm.Execute(); err != nil {
		fmt.Printf("Runtime error: %v\n", err)
		return
	}
}

func runInteractiveMode() {
	fmt.Println("MicroList Interactive Mode v2.0")
	fmt.Println("Enter your code line by line.")
	fmt.Println("Commands: 'run' to execute, 'clear' to reset, 'quit' to exit")
	fmt.Println("Example:")
	fmt.Println("  1 ml")
	fmt.Println("  \"hello\"")
	fmt.Println("  2 ml.")
	fmt.Println("  1[]")
	fmt.Println("  run")
	fmt.Println()
	
	scanner := bufio.NewScanner(os.Stdin)
	var sourceLines []string
	
	for {
		fmt.Print("ML> ")
		if !scanner.Scan() {
			break
		}
		
		line := strings.TrimSpace(scanner.Text())
		
		if line == "quit" || line == "exit" {
			break
		}
		
		if line == "run" {
			// Compile and run
			source := strings.Join(sourceLines, "\n")
			compiler := NewCompiler(source)
			
			if err := compiler.Compile(); err != nil {
				fmt.Printf("❌ Compilation error: %v\n", err)
				continue
			}
			
			vm := NewVM()
			vm.LoadProgram(compiler.GetInstructions())
			
			if err := vm.Execute(); err != nil {
				fmt.Printf("❌ Runtime error: %v\n", err)
			}
			fmt.Println()
			continue
		}
		
		if line == "clear" {
			sourceLines = nil
			fmt.Println("✅ Source cleared.")
			continue
		}
		
		if line == "help" {
			fmt.Println("MicroList Syntax:")
			fmt.Println("  NUM ml        - Create storage list")
			fmt.Println("  NUM ml.       - Create output list")
			fmt.Println("  \"text\"        - Add string to current list")
			fmt.Println("  123           - Add number (0-255) to current list")
			fmt.Println("  NUM[]         - Reference another list's content")
			fmt.Println("  // comment    - Comment line")
			continue
		}
		
		sourceLines = append(sourceLines, line)
		fmt.Printf("    [%d lines entered]\n", len(sourceLines))
	}
}

func compileToBinary(sourceFile string) {
	content, err := os.ReadFile(sourceFile)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}
	
	// Compile
	compiler := NewCompiler(string(content))
	if err := compiler.Compile(); err != nil {
		fmt.Printf("Compilation error: %v\n", err)
		return
	}
	
	// Generate binary
	vm := NewVM()
	vm.LoadProgram(compiler.GetInstructions())
	binary := vm.SerializeToBinary()
	
	// Write binary file
	outputFile := strings.TrimSuffix(sourceFile, ".ml") + ".mlist"
	if err := os.WriteFile(outputFile, binary, 0644); err != nil {
		fmt.Printf("Error writing binary file: %v\n", err)
		return
	}
	
	fmt.Printf("✅ Compiled to binary: %s (%d bytes)\n", outputFile, len(binary))
	fmt.Printf("📊 Instructions: %d\n", len(compiler.GetInstructions()))
}
