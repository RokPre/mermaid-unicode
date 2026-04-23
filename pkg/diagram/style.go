package diagram

import (
	"fmt"
	"strings"

	"github.com/gookit/color"
)

const (
	StyleTypeCLI  = "cli"
	StyleTypeHTML = "html"
)

type StyleMap map[string]string

type StyleClass struct {
	Name   string
	Styles StyleMap
}

func NewStyleClass(name, styles string) StyleClass {
	return StyleClass{Name: name, Styles: ParseStyleMap(styles)}
}

func ParseStyleMap(styles string) StyleMap {
	styleMap := make(StyleMap)
	for _, style := range strings.Split(styles, ",") {
		style = strings.TrimSpace(strings.TrimSuffix(style, ";"))
		if style == "" {
			continue
		}
		kv := strings.SplitN(style, ":", 2)
		if len(kv) != 2 {
			continue
		}
		key := strings.TrimSpace(kv[0])
		value := strings.TrimSpace(strings.TrimSuffix(kv[1], ";"))
		if key != "" && value != "" {
			styleMap[key] = value
		}
	}
	return styleMap
}

func ResolveStyle(defaultStyle, classStyle, directStyle StyleMap) StyleMap {
	resolved := make(StyleMap)
	for _, style := range []StyleMap{defaultStyle, classStyle, directStyle} {
		for key, value := range style {
			if NormalizeStyleColor(value) != "" {
				resolved[key] = value
			}
		}
	}
	return resolved
}

func NormalizeStyleColor(c string) string {
	c = strings.TrimSpace(c)
	c = strings.TrimSpace(strings.TrimSuffix(c, ";"))
	if c == "" || strings.EqualFold(c, "none") || strings.EqualFold(c, "transparent") {
		return ""
	}
	return c
}

func WrapTextInStyle(text, fg, bg, styleType string) string {
	fg = NormalizeStyleColor(fg)
	bg = NormalizeStyleColor(bg)
	if fg == "" && bg == "" {
		return text
	}
	switch styleType {
	case StyleTypeHTML:
		style := []string{}
		if fg != "" {
			style = append(style, fmt.Sprintf("color: %s", fg))
		}
		if bg != "" {
			style = append(style, fmt.Sprintf("background-color: %s", bg))
		}
		return fmt.Sprintf("<span style='%s'>%s</span>", strings.Join(style, "; "), text)
	case StyleTypeCLI:
		return color.HEXStyle(fg, bg).Sprint(text)
	default:
		return text
	}
}

func WrapTextInColor(text, c, styleType string) string {
	return WrapTextInStyle(text, c, "", styleType)
}
