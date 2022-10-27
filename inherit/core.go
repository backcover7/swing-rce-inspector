package inherit

import (
	"github.com/4ra1n/swing-rce-inspector/asm"
	"github.com/4ra1n/swing-rce-inspector/common"
)

var (
	subClassMap         map[string][]string
	implicitInheritance map[string][]string
	classMap            map[string]*asm.Class
)

func Init(discoveryClass []*asm.Class, clazzMap map[string]*asm.Class) {
	classMap = clazzMap
	implicitInheritance = make(map[string][]string)
	for _, class := range discoveryClass {
		var allParents []string
		getAllParents(class, &allParents)
		allParents, _ = common.RemoveDup(allParents)
		implicitInheritance[class.Name()] = allParents
	}
	subClassMap = make(map[string][]string)
	for k, v := range implicitInheritance {
		child := k
		for _, parent := range v {
			if subClassMap[parent] != nil {
				var tempSet []string
				tempSet = append(tempSet, parent)
			} else {
				subClassMap[parent] = append(subClassMap[parent], child)
			}
		}
	}
}

func getAllParents(class *asm.Class, allPatents *[]string) {
	var parents []string
	if class.SuperClassName() != "" {
		parents = append(parents, class.SuperClassName())
	}
	for _, i := range class.InterfaceNames() {
		parents = append(parents, i)
	}
	for _, immediateParent := range parents {
		parentClassReference := classMap[immediateParent]
		if parentClassReference == nil {
			continue
		}
		*allPatents = append(*allPatents, parentClassReference.Name())
		getAllParents(parentClassReference, allPatents)
	}
}

func IsSubclassOf(class string, superClass string) bool {
	parents := implicitInheritance[class]
	if parents == nil {
		return false
	}
	var ok bool
	for _, clazz := range parents {
		if clazz == superClass {
			ok = true
		}
	}
	return ok
}
