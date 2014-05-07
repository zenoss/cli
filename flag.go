package cli

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
)

// This flag enables bash-completion for all commands and subcommands
var BashCompletionFlag = NewBoolFlag("generate-bash-completion", "")

// This flag prints the version for the application
var VersionFlag = NewBoolFlag("version, v", "print the version")

// This flag prints the help for all commands and subcommands
var HelpFlag = NewBoolFlag("help, h", "show help")

// Flag is a common interface related to parsing flags in cli.
// For more advanced flag parsing techniques, it is recomended that
// this interface be implemented.
type Flag interface {
	fmt.Stringer
	// Apply Flag settings to the given flag set
	Apply(*flag.FlagSet)
	getName() string
}

func flagSet(name string, flags []Flag) *flag.FlagSet {
	set := flag.NewFlagSet(name, flag.ContinueOnError)

	for _, f := range flags {
		f.Apply(set)
	}
	return set
}

func eachName(longName string, fn func(string)) {
	parts := strings.Split(longName, ",")
	for _, name := range parts {
		name = strings.Trim(name, " ")
		fn(name)
	}
}

type GenericWrapper struct {
	Generic  *Generic
	Original *[]string
}

type FlagWithOriginal struct {
	Flag         GenericFlag
	Original     []string
	ValueWrapper GenericWrapper
}

func (w GenericWrapper) Set(value string) error {
	*w.Original = append(*w.Original, value)
	return (*w.Generic).Set(value)
}

func (w GenericWrapper) String() string {
	return (*w.Generic).String()
}

func (f FlagWithOriginal) String() string {
	return f.Flag.String()
}

func (f FlagWithOriginal) getName() string {
	return f.Flag.getName()
}

func (f FlagWithOriginal) getValue() *Generic {
	return f.Flag.getValue()
}

func (f FlagWithOriginal) getUsage() string {
	return f.Flag.getUsage()
}

func (f FlagWithOriginal) Apply(set *flag.FlagSet) {
	eachName(f.Flag.getName(), func(name string) {
		set.Var(f.WrappedValue(), name, f.Flag.getUsage())
	})
}

func (f FlagWithOriginal) WrappedValue() Generic {
	if f.ValueWrapper.Generic == nil {
		f.ValueWrapper.Generic = f.Flag.getValue()
		f.ValueWrapper.Original = &f.Original
	}
	return f.ValueWrapper
}

// Generic is a generic parseable type identified by a specific flag
type Generic interface {
	Set(value string) error
	String() string
}

type GenericFlag interface {
	getName() string
	getValue() *Generic
	getUsage() string
	String() string
}

// GenericFlag is the flag type for types implementing Generic
type genericFlag struct {
	Name  string
	Value Generic
	Usage string
}

func NewGenericFlag(name string, value Generic, usage string) FlagWithOriginal {
	return FlagWithOriginal{
		Flag: genericFlag{
			Name:  name,
			Value: value,
			Usage: usage,
		},
	}
}

func (f genericFlag) String() string {
	return fmt.Sprintf("%s%s %v\t`%v` %s", prefixFor(f.Name), f.Name, f.Value, "-"+f.Name+" option -"+f.Name+" option", f.Usage)
}

func (f genericFlag) Apply(set *flag.FlagSet) {
	eachName(f.Name, func(name string) {
		set.Var(f.Value, name, f.Usage)
	})
}

func (f genericFlag) getName() string {
	return f.Name
}

func (f genericFlag) getValue() *Generic {
	return &f.Value
}

func (f genericFlag) getUsage() string {
	return f.Usage
}

func NewStringSlice(values ...string) StringSlice {
	return StringSlice{&values}
}

type StringSlice struct {
	value *[]string
}

func (f StringSlice) Set(value string) error {
	*f.value = append(*f.value, value)
	return nil
}

func (f StringSlice) String() string {
	return fmt.Sprintf("%s", *f.value)
}

func (f StringSlice) Value() []string {
	return *f.value
}

func NewStringSliceFlag(name string, value []string, usage string) FlagWithOriginal {
	return FlagWithOriginal{
		Flag: StringSliceFlag{genericFlag{
			Name:  name,
			Value: &StringSlice{value: &value},
			Usage: usage,
		}},
	}
}

type StringSliceFlag struct{ genericFlag }

func (f StringSliceFlag) String() string {
	firstName := strings.Trim(strings.Split(f.Name, ",")[0], " ")
	pref := prefixFor(firstName)
	return fmt.Sprintf("%s '%v'\t%v", prefixedNames(f.Name), pref+firstName+" option "+pref+firstName+" option", f.Usage)
}

type IntSlice struct {
	value *[]int
}

func (f IntSlice) Set(value string) error {

	tmp, err := strconv.Atoi(value)
	if err != nil {
		return err
	} else {
		*f.value = append(*f.value, tmp)
	}
	return nil
}

func (f IntSlice) String() string {
	return fmt.Sprintf("%d", *f.value)
}

func (f IntSlice) Value() []int {
	return *f.value
}

func NewIntSliceFlag(name string, value []int, usage string) FlagWithOriginal {
	return FlagWithOriginal{
		Flag: IntSliceFlag{genericFlag{
			Name:  name,
			Value: &IntSlice{value: &value},
			Usage: usage,
		}},
	}
}

type IntSliceFlag struct{ genericFlag }

func (f IntSliceFlag) String() string {
	firstName := strings.Trim(strings.Split(f.Name, ",")[0], " ")
	pref := prefixFor(firstName)
	return fmt.Sprintf("%s '%v'\t%v", prefixedNames(f.Name), pref+firstName+" option "+pref+firstName+" option", f.Usage)
}

type Bool struct {
	value *bool
}

func (f Bool) Set(value string) error {
	tmp, err := strconv.ParseBool(value)
	if err != nil {
		return err
	} else {
		*f.value = tmp
	}
	return nil
}

func (f Bool) String() string {
	return strconv.FormatBool(*f.value)
}

func (f Bool) Value() bool {
	return *f.value
}

func NewBoolFlag(name string, usage string) FlagWithOriginal {
	value := false
	return FlagWithOriginal{
		Flag: BoolFlag{genericFlag{
			Name:  name,
			Value: &Bool{value: &value},
			Usage: usage,
		}},
	}
}

func NewBoolTFlag(name string, usage string) FlagWithOriginal {
	value := true
	return FlagWithOriginal{
		Flag: BoolFlag{genericFlag{
			Name:  name,
			Value: &Bool{value: &value},
			Usage: usage,
		}},
	}
}

type BoolFlag struct{ genericFlag }

func (f BoolFlag) String() string {
	return fmt.Sprintf("%s\t%v", prefixedNames(f.Name), f.Usage)
}

type String struct {
	value *string
}

func (f String) Set(value string) error {
	*f.value = value
	return nil
}

func (f String) String() string {
	return *f.value
}

func (f String) Value() string {
	return *f.value
}

func NewStringFlag(name, value, usage string) FlagWithOriginal {
	return FlagWithOriginal{
		Flag: StringFlag{genericFlag{
			Name:  name,
			Value: &String{value: &value},
			Usage: usage,
		}},
	}
}

type StringFlag struct{ genericFlag }

func (f StringFlag) String() string {
	var fmtString string
	fmtString = "%s %v\t%v"

	if len(f.Value.String()) > 0 {
		fmtString = "%s '%v'\t%v"
	} else {
		fmtString = "%s %v\t%v"
	}

	return fmt.Sprintf(fmtString, prefixedNames(f.Name), f.Value, f.Usage)
}

type Int struct {
	value *int
}

func (f Int) Set(value string) error {
	tmp, err := strconv.Atoi(value)
	if err != nil {
		return err
	} else {
		*f.value = tmp
	}
	return nil
}

func (f Int) String() string {
	return fmt.Sprintf("%d", *f.value)
}

func (f Int) Value() int {
	return *f.value
}

func NewIntFlag(name string, value int, usage string) FlagWithOriginal {
	return FlagWithOriginal{
		Flag: IntFlag{genericFlag{
			Name:  name,
			Value: &Int{value: &value},
			Usage: usage,
		}},
	}
}

type IntFlag struct{ genericFlag }

func (f IntFlag) String() string {
	return fmt.Sprintf("%s '%v'\t%v", prefixedNames(f.Name), f.Value, f.Usage)
}

type Float64 struct {
	value *float64
}

func (f Float64) Set(value string) error {
	tmp, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return err
	} else {
		*f.value = tmp
	}
	return nil
}

func (f Float64) String() string {
	return fmt.Sprintf("%f", *f.value)
}

func (f Float64) Value() float64 {
	return *f.value
}

func NewFloat64Flag(name string, value float64, usage string) FlagWithOriginal {
	return FlagWithOriginal{
		Flag: Float64Flag{genericFlag{
			Name:  name,
			Value: &Float64{value: &value},
			Usage: usage,
		}},
	}
}

type Float64Flag struct{ genericFlag }

func (f Float64Flag) String() string {
	return fmt.Sprintf("%s '%v'\t%v", prefixedNames(f.Name), f.Value, f.Usage)
}

func prefixFor(name string) (prefix string) {
	if len(name) == 1 {
		prefix = "-"
	} else {
		prefix = "--"
	}

	return
}

func prefixedNames(fullName string) (prefixed string) {
	parts := strings.Split(fullName, ",")
	for i, name := range parts {
		name = strings.Trim(name, " ")
		prefixed += prefixFor(name) + name
		if i < len(parts)-1 {
			prefixed += ", "
		}
	}
	return
}
