package languages

import (
	"github.com/codersrank-org/repo_info_extractor/v2/librarydetection"
	"golang.org/x/net/html"
	"strings"
)

// NewCAnalyzer constructor
func NewHTMLAnalyzer() librarydetection.Analyzer {
	return &htmlAnalyzer{}
}

type htmlAnalyzer struct{}

func (a *htmlAnalyzer) ExtractLibraries(contents string) ([]string, error) {
	doc, err := html.Parse(strings.NewReader(contents))
	if err != nil {
		return nil, err
	}

	return a.traverseNodes(doc), nil
}

func (a *htmlAnalyzer) traverseNodes(node *html.Node) []string {
	var res []string
	// JS <script>
	if node.Type == html.ElementNode && node.Data == "script" {
		for _, v := range node.Attr {
			if v.Key == "src" && strings.HasSuffix(v.Val, ".js") {
				bits := strings.Split(v.Val, "/")
				return []string{bits[len(bits)-1]}
			}
		}
	}

	// CSS <link>
	if node.Type == html.ElementNode && node.Data == "link" {
		for _, v := range node.Attr {
			if v.Key == "href" && strings.HasSuffix(v.Val, ".css") {
				bits := strings.Split(v.Val, "/")
				return []string{bits[len(bits)-1]} // last element is the lib name
			}
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		res = append(res, a.traverseNodes(child)...)
	}

	return res
}
