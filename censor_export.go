package textcensor

type ITextCensorService interface {
	SetPunctuation(str string)
	InitWordsByPath(path string, caseSensitive bool) error
	InitWords(wordsArr []string, caseSensitive bool)
	CheckAndReplace(text string, strict bool, replaceCharacter rune) (pass bool, newText string)
	IsPass(text string, strict bool) bool
}

func NewTextCensorService() ITextCensorService {
	return CensorServiceConstructor()
}
