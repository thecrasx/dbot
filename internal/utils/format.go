package utils

import "fmt"

func EmojiToUnicodeString(str string) string {

	// Array like [U+0033 U+FE0F U+20E3]
	uni_arr := fmt.Sprintf("%U", []rune(str))

	uni := ""
	for _, c := range uni_arr {
		if c == '+' || c == ' ' || c == '[' || c == ']' {
			continue
		} else if c == 'U' {
			uni += "\\u"
		} else {
			uni += string(c)
		}
	}
	return uni
}

func CreateReactionMessageString(strs ...string) []string {
	out := []string{}
	for i, str := range strs {
		out = append(
			out,
			fmt.Sprintf("d.Session.MessageReactionAdd(msg.ChannelID, msg.ID, \"%s\") // %d", str, i+1),
		)
	}
	return out
}
