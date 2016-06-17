package flaeg

import (
	"errors"
	"fmt"
	flag "github.com/ogier/pflag"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"sort"
	"strings"
	"text/template"
	"time"
)

// GetTypesRecursive links in flagmap a flag with its reflect.StructField
// You can whether provide objValue on a structure or a pointer to structure as first argument
// Flags are genereted from field name or from StructTag
func getTypesRecursive(objValue reflect.Value, flagmap map[string]reflect.StructField, key string) error {
	name := key
	switch objValue.Kind() {
	case reflect.Struct:

		for i := 0; i < objValue.NumField(); i++ {
			if objValue.Type().Field(i).Anonymous {
				if err := getTypesRecursive(objValue.Field(i), flagmap, name); err != nil {
					return err
				}
			} else if len(objValue.Type().Field(i).Tag.Get("description")) > 0 {
				fieldName := objValue.Type().Field(i).Name
				if !isExported(fieldName) {
					return fmt.Errorf("Flied %s is an unexported field", fieldName)
				}

				name += objValue.Type().Name()
				if tag := objValue.Type().Field(i).Tag.Get("long"); len(tag) > 0 {
					fieldName = tag
				}
				if len(key) == 0 {
					//Lower Camel Case
					//name = strings.ToLower(string(fieldName[0])) + fieldName[1:]
					name = strings.ToLower(fieldName)
				} else {
					name = key + "." + strings.ToLower(fieldName)
				}
				if _, ok := flagmap[name]; ok {
					return errors.New("Tag already exists: " + name)
				}
				flagmap[name] = objValue.Type().Field(i)

				if err := getTypesRecursive(objValue.Field(i), flagmap, name); err != nil {
					return err
				}
			}

		}
	case reflect.Ptr:
		if len(key) > 0 {
			field := flagmap[name]
			field.Type = reflect.TypeOf(false)
			flagmap[name] = field
		}
		typ := objValue.Type().Elem()
		inst := reflect.New(typ).Elem()
		if err := getTypesRecursive(inst, flagmap, name); err != nil {
			return err
		}
	default:
		return nil
	}
	return nil
}

//GetFlags returns flags
func GetFlags(config interface{}) ([]string, error) {
	flagmap := make(map[string]reflect.StructField)
	if err := getTypesRecursive(reflect.ValueOf(config), flagmap, ""); err != nil {
		return []string{}, err
	}
	flags := make([]string, 0, len(flagmap))
	for f := range flagmap {
		flags = append(flags, f)
	}
	return flags, nil
}

//loadParsers loads default parsers and custom parsers given as parameter. Return a map [reflect.Type]parsers
// bool, int, int64, uint, uint64, float64,
func loadParsers(customParsers map[reflect.Type]Parser) (map[reflect.Type]Parser, error) {
	parsers := map[reflect.Type]Parser{}

	var boolParser boolValue
	parsers[reflect.TypeOf(true)] = &boolParser

	var intParser intValue
	parsers[reflect.TypeOf(1)] = &intParser

	var int64Parser int64Value
	parsers[reflect.TypeOf(int64(1))] = &int64Parser

	var uintParser uintValue
	parsers[reflect.TypeOf(uint(1))] = &uintParser

	var uint64Parser uint64Value
	parsers[reflect.TypeOf(uint64(1))] = &uint64Parser

	var stringParser stringValue
	parsers[reflect.TypeOf("")] = &stringParser

	var float64Parser float64Value
	parsers[reflect.TypeOf(float64(1.5))] = &float64Parser

	var durationParser durationValue
	parsers[reflect.TypeOf(time.Second)] = &durationParser

	var timeParser timeValue
	parsers[reflect.TypeOf(time.Now())] = &timeParser

	for rType, parser := range customParsers {
		parsers[rType] = parser
	}
	return parsers, nil
}

//ParseArgs : parses args return valmap map[flag]Getter, using parsers map[type]Getter
//args must be formated as like as flag documentation. See https://golang.org/pkg/flag
func parseArgs(args []string, flagmap map[string]reflect.StructField, parsers map[reflect.Type]Parser) (map[string]Parser, error) {
	//Return var
	valmap := make(map[string]Parser)
	//Visitor in flag.Parse
	flagList := []*flag.Flag{}
	visitor := func(fl *flag.Flag) {
		flagList = append(flagList, fl)
	}
	newParsers := map[string]Parser{}
	flagSet := flag.NewFlagSet("flaeg.Load", flag.ContinueOnError)
	//Disable output
	flagSet.SetOutput(ioutil.Discard)
	var errMissingParser error
	for flag, structField := range flagmap {
		//for _, flag := range flags {
		//structField := flagmap[flag]
		if parser, ok := parsers[structField.Type]; ok {
			newparserValue := reflect.New(reflect.TypeOf(parser).Elem())
			newparserValue.Elem().Set(reflect.ValueOf(parser).Elem())
			newparser := newparserValue.Interface().(Parser)
			if short := structField.Tag.Get("short"); len(short) == 1 {
				// fmt.Printf("short : %s long : %s\n", short, flag)
				flagSet.VarP(newparser, flag, short, structField.Tag.Get("description"))
			} else {
				flagSet.Var(newparser, flag, structField.Tag.Get("description"))
			}
			newParsers[flag] = newparser
		} else {
			errMissingParser = fmt.Errorf("%s :No parser for type %s", flag, structField.Type)
			// return nil, fmt.Errorf("%s :No parser for type %s", flag, structField.Type)
		}
	}

	// prevents case sensitivity issue
	args = argsToLower(args)
	if err := flagSet.Parse(args); err != nil {
		return nil, err
	}

	//Fill flagList with parsed flags
	flagSet.Visit(visitor)
	//Return parsers on parsed flag
	for _, flag := range flagList {
		valmap[flag.Name] = newParsers[flag.Name]
	}

	return valmap, errMissingParser
}

func getDefaultValue(defaultValue reflect.Value, defaultPointersValue reflect.Value, defaultValmap map[string]reflect.Value, key string) error {
	if defaultValue.Type() != defaultPointersValue.Type() {
		return fmt.Errorf("Parameters defaultValue and defaultPointersValue must be the same struct. defaultValue type : %s is not defaultPointersValue type : %s", defaultValue.Type().String(), defaultPointersValue.Type().String())
	}
	name := key
	switch defaultValue.Kind() {
	case reflect.Struct:
		for i := 0; i < defaultValue.NumField(); i++ {
			if defaultValue.Type().Field(i).Anonymous {
				if err := getDefaultValue(defaultValue.Field(i), defaultPointersValue.Field(i), defaultValmap, name); err != nil {
					return err
				}
			} else if len(defaultValue.Type().Field(i).Tag.Get("description")) > 0 {
				name += defaultValue.Type().Name()
				fieldName := defaultValue.Type().Field(i).Name
				if tag := defaultValue.Type().Field(i).Tag.Get("long"); len(tag) > 0 {
					fieldName = tag
				}
				if len(key) == 0 {
					name = strings.ToLower(fieldName)
				} else {
					name = key + "." + strings.ToLower(fieldName)
				}
				if defaultValue.Field(i).Kind() != reflect.Ptr {
					// if _, ok := defaultValmap[name]; ok {
					// 	return errors.New("Tag already exists: " + name)
					// }
					defaultValmap[name] = defaultValue.Field(i)
					// fmt.Printf("%s: got default value %+v\n", name, defaultValue.Field(i))
				}
				if err := getDefaultValue(defaultValue.Field(i), defaultPointersValue.Field(i), defaultValmap, name); err != nil {
					return err
				}
			}
		}
	case reflect.Ptr:
		if !defaultPointersValue.IsNil() {
			if len(key) != 0 {
				//turn ptr fields to nil
				defaultPointersNilValue, err := setPointersNil(defaultPointersValue)
				if err != nil {
					return err
				}
				defaultValmap[name] = defaultPointersNilValue
				// fmt.Printf("%s: got default value %+v\n", name, defaultPointersNilValue)
			}
			if !defaultValue.IsNil() {
				if err := getDefaultValue(defaultValue.Elem(), defaultPointersValue.Elem(), defaultValmap, name); err != nil {
					return err
				}
			} else {
				if err := getDefaultValue(defaultPointersValue.Elem(), defaultPointersValue.Elem(), defaultValmap, name); err != nil {
					return err
				}
			}
		} else {
			instValue := reflect.New(defaultPointersValue.Type().Elem())
			if len(key) != 0 {
				defaultValmap[name] = instValue
				// fmt.Printf("%s: got default value %+v\n", name, instValue)
			}
			if !defaultValue.IsNil() {
				if err := getDefaultValue(defaultValue.Elem(), instValue.Elem(), defaultValmap, name); err != nil {
					return err
				}
			} else {
				if err := getDefaultValue(instValue.Elem(), instValue.Elem(), defaultValmap, name); err != nil {
					return err
				}
			}
		}
	default:
		return nil
	}
	return nil
}

//objValue a reflect.Value of a not-nil pointer on a struct
func setPointersNil(objValue reflect.Value) (reflect.Value, error) {
	if objValue.Kind() != reflect.Ptr {
		return objValue, fmt.Errorf("Parameters objValue must be a not-nil pointer on a struct, not a %s", objValue.Kind().String())
	} else if objValue.IsNil() {
		return objValue, fmt.Errorf("Parameters objValue must be a not-nil pointer")
	} else if objValue.Elem().Kind() != reflect.Struct {
		// fmt.Printf("Parameters objValue must be a not-nil pointer on a struct, not a pointer on a %s\n", objValue.Elem().Kind().String())
		return objValue, nil
	}
	//Clone
	starObjValue := objValue.Elem()
	nilPointersObjVal := reflect.New(starObjValue.Type())
	starNilPointersObjVal := nilPointersObjVal.Elem()
	starNilPointersObjVal.Set(starObjValue)

	for i := 0; i < nilPointersObjVal.Elem().NumField(); i++ {
		if field := nilPointersObjVal.Elem().Field(i); field.Kind() == reflect.Ptr && field.CanSet() {
			field.Set(reflect.Zero(field.Type()))
		}
	}
	return nilPointersObjVal, nil
}

//FillStructRecursive initialize a value of any taged Struct given by reference
func fillStructRecursive(objValue reflect.Value, defaultPointerValmap map[string]reflect.Value, valmap map[string]Parser, key string) error {
	name := key
	switch objValue.Kind() {
	case reflect.Struct:

		for i := 0; i < objValue.Type().NumField(); i++ {
			if objValue.Type().Field(i).Anonymous {
				if err := fillStructRecursive(objValue.Field(i), defaultPointerValmap, valmap, name); err != nil {
					return err
				}
			} else if len(objValue.Type().Field(i).Tag.Get("description")) > 0 {
				name += objValue.Type().Name()
				fieldName := objValue.Type().Field(i).Name
				if tag := objValue.Type().Field(i).Tag.Get("long"); len(tag) > 0 {
					fieldName = tag
				}
				if len(key) == 0 {
					name = strings.ToLower(fieldName)
				} else {
					name = key + "." + strings.ToLower(fieldName)
				}
				// fmt.Println(name)
				if objValue.Field(i).Kind() != reflect.Ptr {

					if val, ok := valmap[name]; ok {
						// fmt.Printf("%s : set def val\n", name)
						if err := setFields(objValue.Field(i), val); err != nil {
							return err
						}
					}
				}
				if err := fillStructRecursive(objValue.Field(i), defaultPointerValmap, valmap, name); err != nil {
					return err
				}
			}
		}

	case reflect.Ptr:
		if len(key) == 0 && !objValue.IsNil() {
			if err := fillStructRecursive(objValue.Elem(), defaultPointerValmap, valmap, name); err != nil {
				return err
			}
			return nil
		}
		contains := false
		for flag := range valmap {
			// TODO replace by regexp
			if strings.Contains(flag, name+".") {
				contains = true
				break
			}
		}
		needDefault := false
		if _, ok := valmap[name]; ok {
			needDefault = valmap[name].Get().(bool)
		}
		if contains && objValue.IsNil() {
			needDefault = true
		}

		if needDefault {
			if defVal, ok := defaultPointerValmap[name]; ok {
				//set default pointer value
				// fmt.Printf("%s  : set default value %+v\n", name, defVal)
				objValue.Set(defVal)
			} else {
				return fmt.Errorf("flag %s default value not provided", name)
			}
		}
		if !objValue.IsNil() && contains {
			if objValue.Type().Elem().Kind() == reflect.Struct {
				if err := fillStructRecursive(objValue.Elem(), defaultPointerValmap, valmap, name); err != nil {
					return err
				}
			}
		}
	default:
		return nil
	}
	return nil
}

// SetFields sets value to fieldValue using tag as key in valmap
func setFields(fieldValue reflect.Value, val Parser) error {
	if fieldValue.CanSet() {
		fieldValue.Set(reflect.ValueOf(val).Elem().Convert(fieldValue.Type()))
	} else {
		return errors.New(fieldValue.Type().String() + " is not settable.")
	}
	return nil
}

//PrintHelp generates and prints command line help
func PrintHelp(flagmap map[string]reflect.StructField, defaultValmap map[string]reflect.Value, parsers map[reflect.Type]Parser) error {
	// Define a templates
	// Using POSXE STD : http://pubs.opengroup.org/onlinepubs/9699919799/
	const helper = `
Usage: {{.ProgName}} [--flag=flag_argument] [-f[flag_argument]] ...     set flag_argument to flag(s)
   or: {{.ProgName}} [--flag[=true|false| ]] [-f[true|false| ]] ...     set true/false to boolean flag(s)   

Flags:{{range $j, $flag := .Flags}}{{$description:= index $.Descriptions $j}}{{$defaultValues := index $.DefaultValues $j}}
{{printf "\t%-50s %s (default \"%s\")" $flag $description $defaultValues}}{{end}}`

	// Preprocess data

	// Sort alphabetically & Delete unparsable flags in a slice
	flags := []string{}
	for flag, field := range flagmap {
		if _, ok := parsers[field.Type]; ok {
			flags = append(flags, flag)
		}
	}
	sort.Strings(flags)

	// Process data
	printDescriptions := []string{}
	printDefaultValues := []string{}
	printFlags := []string{}
	for _, flag := range flags {
		field := flagmap[flag]
		if short := field.Tag.Get("short"); len(short) == 1 {
			printFlags = append(printFlags, "-"+short+", --"+flag)
		} else {
			printFlags = append(printFlags, "--"+flag)
		}
		printDescriptions = append(printDescriptions, field.Tag.Get("description"))
		//flag on pointer ?
		if defVal, ok := defaultValmap[flag]; ok {
			if defVal.Kind() != reflect.Ptr {
				// Set defaultValue on parsers
				parsers[field.Type].SetValue(defaultValmap[flag].Interface())
			}
			printDefaultValues = append(printDefaultValues, parsers[field.Type].String())
		} /*else {
			printDefaultValues = append(printDefaultValues, "N/A")
		}*/
	}

	// Use a struct to give data to template
	type TempStruct struct {
		ProgName      string
		Flags         []string
		Descriptions  []string
		DefaultValues []string
	}
	tempStruct := TempStruct{
		Flags:         printFlags,
		Descriptions:  printDescriptions,
		DefaultValues: printDefaultValues,
	}
	_, tempStruct.ProgName = path.Split(os.Args[0])

	//Run Template
	tmplHelper, err := template.New("helper").Parse(helper)
	if err != nil {
		return err
	}
	err = tmplHelper.Execute(os.Stdout, tempStruct)
	if err != nil {
		return err
	}
	//And footer
	//FIXME --help ?
	fmt.Fprintf(os.Stdout, "\n\t%-50s %s\n", "-h, --help", "Print Help (this message) and exit")
	return nil
}

//PrintError takes a not nil error and prints command line help
func PrintError(err error, flagmap map[string]reflect.StructField, defaultValmap map[string]reflect.Value, parsers map[reflect.Type]Parser) error {
	if err != flag.ErrHelp {
		fmt.Printf("Error : %s\n", err)
	}
	if !strings.Contains(err.Error(), ":No parser for type") {
		PrintHelp(flagmap, defaultValmap, parsers)
	}
	return err
}

//LoadWithParsers initializes config : struct fields given by reference, with args : arguments.
//Some custom parsers may be given.
func LoadWithParsers(config interface{}, defaultValue interface{}, args []string, customParsers map[reflect.Type]Parser) error {
	parsers, err := loadParsers(customParsers)
	if err != nil {
		return err
	}

	// for typ, parser := range parsers {
	// 	fmt.Printf("%s : %+v\n", typ.Name(), parser)
	// }

	tagsmap := make(map[string]reflect.StructField)
	if err := getTypesRecursive(reflect.ValueOf(config), tagsmap, ""); err != nil {
		return err
	}
	defaultValmap := make(map[string]reflect.Value)
	if err := getDefaultValue(reflect.ValueOf(config), reflect.ValueOf(defaultValue), defaultValmap, ""); err != nil {
		return err
	}
	// for flag := range defaultValmap {
	// 	fmt.Println(flag)
	// }
	valmap, err := parseArgs(args, tagsmap, parsers)
	if err != nil {
		return PrintError(err, tagsmap, defaultValmap, parsers)
	}
	// for flag, val := range valmap {
	// 	fmt.Printf("%s : %+s (default : %+v)\n", flag, val, defaultValmap[flag])
	// }
	if err := fillStructRecursive(reflect.ValueOf(config), defaultValmap, valmap, ""); err != nil {
		return err
	}

	return nil
}

//Load initializes config : struct fields given by reference, with args : arguments.
//Some custom parsers may be given.
func Load(config interface{}, defaultValue interface{}, args []string) error {
	parsers, err := loadParsers(nil)
	if err != nil {
		return err
	}
	tagsmap := make(map[string]reflect.StructField)
	if err := getTypesRecursive(reflect.ValueOf(config), tagsmap, ""); err != nil {
		return err
	}
	defaultValmap := make(map[string]reflect.Value)
	if err := getDefaultValue(reflect.ValueOf(config), reflect.ValueOf(defaultValue), defaultValmap, ""); err != nil {
		return err
	}
	valmap, errParseArgs := parseArgs(args, tagsmap, parsers)
	if errParseArgs != nil && !strings.Contains(errParseArgs.Error(), "No parser for type") {
		return PrintError(errParseArgs, tagsmap, defaultValmap, parsers)
	}
	if err := fillStructRecursive(reflect.ValueOf(config), defaultValmap, valmap, ""); err != nil {
		return err
	}

	return errParseArgs
}

// Command structure contains program/command information (command name and description)
// Config must be a pointer on the configuration struct to parse (it contains default values of field)
// DefaultPointersConfig contains default pointers values: those values are set on pointers fields if their flags are called
// It must be the same type(struct) as Config
// Run is the func which launch the program using initialized configuration structure
type Command struct {
	Name                  string
	Description           string
	Config                interface{}
	DefaultPointersConfig interface{} //TODO:case DefaultPointersConfig is nil
	Run                   func() error
}

//LoadWithCommand initializes config : struct fields given by reference, with args : arguments.
//Some custom parsers and some subCommand may be given.
func LoadWithCommand(cmd *Command, cmdArgs []string, customParsers map[reflect.Type]Parser, subCommand []*Command) error {

	parsers, err := loadParsers(customParsers)
	if err != nil {
		return err
	}

	tagsmap := make(map[string]reflect.StructField)
	if err := getTypesRecursive(reflect.ValueOf(cmd.Config), tagsmap, ""); err != nil {
		return err
	}
	defaultValmap := make(map[string]reflect.Value)
	if err := getDefaultValue(reflect.ValueOf(cmd.Config), reflect.ValueOf(cmd.DefaultPointersConfig), defaultValmap, ""); err != nil {
		return err
	}

	valmap, err := parseArgs(cmdArgs, tagsmap, parsers)
	if err != nil {
		return PrintErrorWithCommand(err, tagsmap, defaultValmap, parsers, cmd, subCommand)
	}

	if err := fillStructRecursive(reflect.ValueOf(cmd.Config), defaultValmap, valmap, ""); err != nil {
		return err
	}

	return nil
}

//PrintHelpWithCommand generates and prints command line help for a Command
func PrintHelpWithCommand(flagmap map[string]reflect.StructField, defaultValmap map[string]reflect.Value, parsers map[reflect.Type]Parser, cmd *Command, subCmd []*Command) error {
	// Define a templates
	// Using POSXE STD : http://pubs.opengroup.org/onlinepubs/9699919799/
	const helper = `{{.ProgDescription}}
	
Usage: {{.ProgName}} [--flag=flag_argument] [-f[flag_argument]] ...     set flag_argument to flag(s)
   or: {{.ProgName}} [--flag[=true|false| ]] [-f[true|false| ]] ...     set true/false to boolean flag(s)
{{if .SubCommands}}
Available Commands:{{range $subCmdName, $subCmdDesc := .SubCommands}}
{{printf "\t%-50s %s" $subCmdName $subCmdDesc}}{{end}}
Use "{{.ProgName}} [command] --help" for more information about a command.
{{end}}
Flags:{{range $j, $flag := .Flags}}{{$description:= index $.Descriptions $j}}{{$defaultValues := index $.DefaultValues $j}}
{{printf "\t%-50s %s (default \"%s\")" $flag $description $defaultValues}}{{end}}`

	// Preprocess data

	// Sort alphabetically & Delete unparsable flags in a slice
	flags := []string{}
	for flag, field := range flagmap {
		if _, ok := parsers[field.Type]; ok {
			flags = append(flags, flag)
		}
	}
	sort.Strings(flags)

	// Process data
	printDescriptions := []string{}
	printDefaultValues := []string{}
	printFlags := []string{}
	for _, flag := range flags {
		field := flagmap[flag]
		if short := field.Tag.Get("short"); len(short) == 1 {
			printFlags = append(printFlags, "-"+short+", --"+flag)
		} else {
			printFlags = append(printFlags, "--"+flag)
		}
		printDescriptions = append(printDescriptions, field.Tag.Get("description"))
		//flag on pointer ?
		if defVal, ok := defaultValmap[flag]; ok {
			if defVal.Kind() != reflect.Ptr {
				// Set defaultValue on parsers
				parsers[field.Type].SetValue(defaultValmap[flag].Interface())
			}
			printDefaultValues = append(printDefaultValues, parsers[field.Type].String())
		} /*else {
			printDefaultValues = append(printDefaultValues, "N/A")
		}*/
	}

	// Use a struct to give data to template
	type TempStruct struct {
		ProgName        string
		ProgDescription string
		SubCommands     map[string]string
		Flags           []string
		Descriptions    []string
		DefaultValues   []string
	}
	tempStruct := TempStruct{
		Flags:         printFlags,
		Descriptions:  printDescriptions,
		DefaultValues: printDefaultValues,
	}
	if cmd != nil {
		tempStruct.ProgName = cmd.Name
		tempStruct.ProgDescription = cmd.Description
		tempStruct.SubCommands = map[string]string{}
		if len(subCmd) > 1 && cmd == subCmd[0] {
			for _, c := range subCmd[1:] {
				tempStruct.SubCommands[c.Name] = c.Description
			}
		}
	} else {
		_, tempStruct.ProgName = path.Split(os.Args[0])
		tempStruct.ProgDescription = "N/A"
	}

	//Run Template
	tmplHelper, err := template.New("helper").Parse(helper)
	if err != nil {
		return err
	}
	err = tmplHelper.Execute(os.Stdout, tempStruct)
	if err != nil {
		return err
	}
	//And footer
	fmt.Fprintf(os.Stdout, "\n\t%-50s %s\n", "-h, --help", "Print Help (this message) and exit")
	return nil
}

//PrintErrorWithCommand takes a not nil error and prints command line help
func PrintErrorWithCommand(err error, flagmap map[string]reflect.StructField, defaultValmap map[string]reflect.Value, parsers map[reflect.Type]Parser, cmd *Command, subCmd []*Command) error {
	if err != flag.ErrHelp {
		fmt.Printf("Error : %s\n", err)
	}
	PrintHelpWithCommand(flagmap, defaultValmap, parsers, cmd, subCmd)
	return err
}

//Flaeg struct contains commands (at least the root one)
//and row arguments (command and/or flags)
//a map of custom parsers could be use
type Flaeg struct {
	calledCommand *Command
	commands      []*Command ///rootCommand is th fist one in this slice
	args          []string
	commmandArgs  []string
	customParsers map[reflect.Type]Parser
}

//New creats and initialize a pointer on Flaeg
func New(rootCommand *Command, args []string) *Flaeg {
	var f Flaeg
	f.commands = []*Command{rootCommand}
	f.args = args
	f.customParsers = map[reflect.Type]Parser{}
	return &f
}

//AddCommand adds sub-command to the root command
func (f *Flaeg) AddCommand(command *Command) {
	f.commands = append(f.commands, command)
}

//AddParser adds custom parser for a type to the map of custom parsers
func (f *Flaeg) AddParser(typ reflect.Type, parser Parser) {
	f.customParsers[typ] = parser
}

// Run calls the command with flags given as agruments
func (f *Flaeg) Run() error {
	if f.calledCommand == nil {
		if _, _, err := f.findCommandWithCommandArgs(); err != nil {
			return err
		}
	}
	if _, err := f.Parse(f.calledCommand); err != nil {
		return err
	}
	return f.calledCommand.Run()
}

// Parse calls Flaeg Load Function end returns the parsed command structure (by reference)
// It returns nil and a not nil error if it fails
func (f *Flaeg) Parse(cmd *Command) (*Command, error) {
	if f.calledCommand == nil {
		f.commmandArgs = f.args
	}
	if err := LoadWithCommand(cmd, f.commmandArgs, f.customParsers, f.commands); err != nil {
		return nil, err
	}
	return cmd, nil
}

//splitArgs takes args (type []string) and return command ("" if rootCommand) and command's args
func splitArgs(args []string) (string, []string) {
	if len(args) >= 1 && len(args[0]) >= 1 && string(args[0][0]) != "-" {
		if len(args) == 1 {
			return strings.ToLower(args[0]), []string{}
		}
		return strings.ToLower(args[0]), args[1:]
	}
	return "", args
}

// findCommandWithCommandArgs returns the called command (by reference) and command's args
// the error returned is not nil if it fails
func (f *Flaeg) findCommandWithCommandArgs() (*Command, []string, error) {
	commandName := ""
	commandName, f.commmandArgs = splitArgs(f.args)
	if len(commandName) > 0 {
		for _, command := range f.commands {
			if commandName == command.Name {
				f.calledCommand = command
				return f.calledCommand, f.commmandArgs, nil
			}
		}
		return nil, []string{}, fmt.Errorf("Command %s not found", commandName)
	}

	f.calledCommand = f.commands[0]
	return f.calledCommand, f.commmandArgs, nil
}

// GetCommand splits args and returns the called command (by reference)
// It returns nil and a not nil error if it fails
func (f *Flaeg) GetCommand() (*Command, error) {
	if f.calledCommand == nil {
		_, _, err := f.findCommandWithCommandArgs()
		return f.calledCommand, err
	}
	return f.calledCommand, nil
}

//isExported return true is the field (from fieldName) is exported,
//else false
func isExported(fieldName string) bool {
	if len(fieldName) < 1 {
		return false
	}
	if string(fieldName[0]) == strings.ToUpper(string(fieldName[0])) {
		return true
	}
	return false
}

func argToLower(inArg string) string {
	if len(inArg) < 2 {
		return strings.ToLower(inArg)
	}
	var outArg string
	dashIndex := strings.Index(inArg, "--")
	if dashIndex == -1 {
		if dashIndex = strings.Index(inArg, "-"); dashIndex == -1 {
			return inArg
		}
		//-fValue
		outArg = strings.ToLower(inArg[dashIndex:dashIndex+2]) + inArg[dashIndex+2:]
		return outArg
	}
	//--flag
	if equalIndex := strings.Index(inArg, "="); equalIndex != -1 {
		//--flag=value
		outArg = strings.ToLower(inArg[dashIndex:equalIndex]) + inArg[equalIndex:]
	} else {
		//--boolflag
		outArg = strings.ToLower(inArg[dashIndex:])
	}

	return outArg
}

func argsToLower(inArgs []string) []string {
	outArgs := make([]string, len(inArgs), len(inArgs))
	for i, inArg := range inArgs {
		outArgs[i] = argToLower(inArg)
	}
	return outArgs
}
