package util

import (
	"github.com/smilingpoplar/translate/util"
)

func ApplyTranslationFixes(texts []string, fixesMap map[string]string) {
	fixes := []util.FixTransform{}
	for k, v := range fixesMap {
		fixes = append(fixes, util.FixTransform{From: k, To: v})
	}
	util.ApplyTranslationFixes(texts, fixes)
}
