package textcensor

import (
	"io/ioutil"
	"strings"
)

type _TextCensorService struct {
	defaultPunctuation   string
	defaultCaseSensitive bool
	tree                 *sTree
	punctuations         map[rune]bool
}

func CensorServiceConstructor() *_TextCensorService {
	return &_TextCensorService{
		defaultPunctuation:   " !\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~，。？；：”’￥（）——、！……",
		defaultCaseSensitive: false,
		tree:                 &sTree{&runeNode{false, make(map[rune]*runeNode, 1000)}},
		punctuations:         getPunctuationMap(defaultPunctuation),
	}
}

func (service _TextCensorService) SetPunctuation(str string) {
	service.punctuations = getPunctuationMap(str)
}

func (service _TextCensorService) initOneWord(str string, caseSensitive bool) {
	str = strings.TrimSpace(str)
	l := len(str)
	if l <= 0 {
		return
	}
	if !caseSensitive {
		str = strings.ToLower(str)
	}
	runeArr := []rune(str)
	l = len(runeArr)
	node := service.tree.Root
	for i := 0; i < l; i++ {
		v := runeArr[i]
		if v == bomHead {
			//fmt.Printf("bomHead = %+v\n", bomHead)
			continue
		}
		next := node.find(v)
		if next == nil {
			next = &runeNode{}
			node.add(v, next)
		}
		node = next
		if i == l-1 {
			node.isEnd = true
		}
	}
}

func (service _TextCensorService) InitWordsByPath(path string, caseSensitive bool) error {
	words, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	str := string(words)
	wordsArr := strings.Split(str, "\n")
	InitWords(wordsArr, caseSensitive)
	return nil
}

func (service _TextCensorService) InitWords(wordsArr []string, caseSensitive bool) {
	//tree = &STree{&Node{false, make(map[rune]*Node, 1000)}}
	service.defaultCaseSensitive = caseSensitive
	for _, v := range wordsArr {
		service.initOneWord(v, caseSensitive)
	}
}

func (service _TextCensorService) CheckAndReplace(text string, strict bool, replaceCharacter rune) (pass bool, newText string) {
	text = strings.TrimSpace(text)
	if len(text) < 1 {
		return true, text
	}
	originText := text
	if !service.defaultCaseSensitive {
		text = strings.ToLower(text)
	}
	runeArr := []rune(text)
	originRuneArr := []rune(originText)
	l := len(runeArr)

	pass = true
	for i := 0; i < l; i++ {
		cWord := runeArr[i]
		node := service.tree.Root.find(cWord)

		if node == nil {
			continue
		}
		if node.isEnd {
			originRuneArr[i] = replaceCharacter
			pass = false
			continue
		}

		for j := i + 1; j < l; j++ {
			//如果是严格模式，将所有的标点忽略掉
			nextNode := node.find(runeArr[j])
			if nextNode == nil && strict && service.punctuations[runeArr[j]] {
				continue
			}
			node = nextNode
			if node == nil {
				break
			}
			if node.isEnd {
				for ri := i; ri <= j; ri++ {
					originRuneArr[ri] = replaceCharacter
					pass = false
				}
			}
		}

	}

	return pass, string(originRuneArr)
}

func (service _TextCensorService) IsPass(text string, strict bool) bool {
	text = strings.TrimSpace(text)
	if len(text) < 1 {
		return true
	}
	if !service.defaultCaseSensitive {
		text = strings.ToLower(text)
	}
	runeArr := []rune(text)
	l := len(runeArr)

	for i := 0; i < l; i++ {
		cWord := runeArr[i]
		node := service.tree.Root.find(cWord)

		if node == nil {
			continue
		}
		if node.isEnd {
			return false
		}

		for j := i + 1; j < l; j++ {
			//如果是严格模式，将所有的标点忽略掉
			nextNode := node.find(runeArr[j])
			if nextNode == nil && strict && service.punctuations[runeArr[j]] {
				continue
			}
			node = nextNode
			if node == nil {
				break
			}
			if node.isEnd {
				return false
			}
		}

	}
	return true
}
