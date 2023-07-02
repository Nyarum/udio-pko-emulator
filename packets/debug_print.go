package packets

import (
	"fmt"
	"little/types"
	"reflect"

	"github.com/fatih/structs"
	"github.com/gookit/color"
)

func renderHex(buf []byte, v interface{}, field string, t string) {
	color.Redp(fmt.Sprintf("%# x", buf))
	fmt.Print(" | ")
	color.Greenp(fmt.Sprintf("%v", v))
	fmt.Print(" | ")
	color.Cyanp(field)
	fmt.Print(" | ")
	color.Bluep(t)
	fmt.Print("\n")
}

func DebugPrint(v interface{}) {
	color.Blueln("Debug print of packet:")
	color.Blueln("------------")

	debugPrint(v)

	color.Blueln("------------")
	fmt.Println()
}

func debugPrint(v interface{}) {
	f := structs.Fields(v)

	for _, v := range f {
		if !v.IsExported() {
			fmt.Println(v.Name())
			continue
		}

		pullContext, ok := v.Value().(types.PullContext)
		if !ok {
			if v.Kind() == reflect.Array || v.Kind() == reflect.Slice {
				s := reflect.ValueOf(v.Value())

				for i := 0; i < s.Len(); i++ {
					index := s.Index(i)

					internalPullContext, ok := index.Interface().(types.PullContext)
					if ok {
						getValue := internalPullContext.GetValue()

						if len(internalPullContext.GetBuf()) == 0 {
							continue
						}

						switch typeConvValue := getValue.(type) {
						case string:
							renderHex(internalPullContext.GetBuf(), internalPullContext.GetValue(), v.Name(), internalPullContext.GetType()+fmt.Sprintf(" (%v)", len(typeConvValue)))
						default:
							renderHex(internalPullContext.GetBuf(), internalPullContext.GetValue(), v.Name(), internalPullContext.GetType())
						}
						continue
					}

					debugPrint(index.Interface())
				}
				continue
			}

			debugPrint(v.Value())
			continue
		}

		getValue := pullContext.GetValue()

		if len(pullContext.GetBuf()) == 0 {
			continue
		}

		switch typeConvValue := getValue.(type) {
		case string:
			renderHex(pullContext.GetBuf(), pullContext.GetValue(), v.Name(), pullContext.GetType()+fmt.Sprintf(" (%v)", len(typeConvValue)))
		default:
			renderHex(pullContext.GetBuf(), pullContext.GetValue(), v.Name(), pullContext.GetType())
		}
	}
}
