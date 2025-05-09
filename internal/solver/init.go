package solver

import (
	"github.com/fathurwithyou/Tubes2_BE_RecipePlayground/internal/model"
)

var elementsMapGlobal map[string]model.Element

func InitElementsMap(data model.Data) {
	elementsMapGlobal = make(map[string]model.Element)
	for _, el := range data.Elements {
		elementsMapGlobal[el.Name] = el
	}
}
