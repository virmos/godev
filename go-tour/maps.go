import (
	"code.google.com/p/go-tour/wc"
	"strings"
)

func WordCount(s string) map[string]int {
	res := make(map[string]int)
	for _, w := range strings.Fields(s) {
			res[w]++
					
	}
	return res
}

func main() {
	wc.Test(WordCount)
}