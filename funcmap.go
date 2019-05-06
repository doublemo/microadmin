package msadmin

import (
	"html/template"

	"github.com/doublemo/msadmin/config"
)

func FuncMap(r *config.Registry) template.FuncMap {
	var funcmap template.FuncMap
	{
		funcmap = make(template.FuncMap)
		funcmap["add"] = func(a int, b ...int) int {
			for _, i := range b {
				a += i
			}
			return a
		}

		funcmap["sub"] = func(a int, b ...int) int {
			for _, i := range b {
				a -= i
			}
			return a
		}

		funcmap["mul"] = func(a int, b ...int) int {
			for _, i := range b {
				a *= i
			}

			return a
		}

		funcmap["div"] = func(a int, b ...int) int {
			for _, i := range b {
				if i == 0 {
					return a
				}

				a /= i
			}
			return a
		}
	}
	return funcmap
}
