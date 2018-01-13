// classify project classify.go
package classify

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"

	"golang.org/x/net/html"
)

type Arena struct {
	List []Node
}

func (a *Arena) Get(id int) Node {
	return a.List[id]
}

func (a *Arena) NodesByClass(className string) []Node {
	res := []Node{}
	for _, n := range a.List {
		if n.HasClass(className) {
			res = append(res, n)
		}
	}
	return res
}

func (a *Arena) IndexesByClass(className string) []int {
	res := []int{}
	for i, n := range a.List {
		if n.HasClass(className) {
			res = append(res, i)
		}
	}
	return res
}

func (a *Arena) NodesByAttr(k, v string) []Node {
	res := []Node{}
	for _, n := range a.List {
		for _, attr := range n.Attr {
			if attr.Key == k && attr.Val == v {
				res = append(res, n)
			}
		}
	}
	return res
}

func (a *Arena) IndexesByAttr(k string, v string) []int {
	res := []int{}
	for id, n := range a.List {
		for _, attr := range n.Attr {
			if attr.Key == k && attr.Val == v {
				res = append(res, id)
			}
		}
	}
	return res
}

func (a *Arena) AddChild(p int, c int) {
	a.List[p].Children = append(a.List[p].Children, c)
	a.List[c].Parent = p
}

func NewArena(root html.Node) *Arena {
	result := NewArenaRoot()
	result.transform(0, root)
	return result
}

func (a *Arena) Append(root html.Node) {
	a.transform(0, root)
}

func NewArenaRoot() *Arena {
	return &Arena{
		List: make([]Node, 0),
	}
}

func (a *Arena) HasParent(child, parent int) bool {
	n := a.Get(child)
	for n.Parent != 0 {
		if n.Parent == parent {
			return true
		}
		n = a.Get(n.Parent)
	}
	return false
}

func (a *Arena) transform(node_index int, n html.Node) {
	if n.Type == html.CommentNode ||
		n.Type == html.ErrorNode ||
		(n.Type == html.ElementNode && strings.ToLower(n.Data) == "noscript") ||
		(n.Type == html.ElementNode && strings.ToLower(n.Data) == "script") ||
		(n.Type == html.TextNode && strings.TrimSpace(n.Data) == "") {
		return
	}

	a.List = append(a.List, *NewNode(n))
	currentId := len(a.List) - 1
	a.AddChild(node_index, currentId)
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		a.transform(currentId, *c)
	}
}

func clone(src, dst *Arena, srcId, dstId int, dstParentId int) {
	dst.List[dstId].Type = src.List[srcId].Type
	dst.List[dstId].Data = src.List[srcId].Data
	dst.List[dstId].Attr = src.List[srcId].Attr
	dst.AddChild(dstParentId, dstId)
	for _, c := range src.List[srcId].Children {
		dst.List = append(dst.List, Node{})
		dstChildId := len(dst.List) - 1
		clone(src, dst, c, dstChildId, dstId)
		//dst.AddChild(dstParentId, dstChildId)
	}
}

func (a *Arena) Clone(srcId int) (dest *Arena) {
	dest = &Arena{make([]Node, 2)}
	clone(a, dest, srcId, 1, 0)
	return dest
}

// clone Node; dstId should aready exist in 'dest'
/*
func (a *Arena) Clone(srcId int, dstId int, dest *Arena) {
	for _, c := range a.List[srcId].Children {
		dest.list = append(dest.list, Node{
			Type: a.List[c].Type,
			Data: a.List[c].Data,
			Attr: a.List[c].Attr,
		})
		destChildId := len(dest.list) - 1
		dest.AddChild(dstId, destChildId)
		a.Clone(c, destChildId, dest)
	}
}
*/

func (a *Arena) PrintList() string {
	res := ""
	for _, n := range a.List {
		if n.Type == html.TextNode {
			res += "text:" + strings.TrimSpace(n.Data) + "\n"
		} else {
			res += n.Data + ":" + n.printAttr() + "\n"
		}
	}
	return res
}

func (a *Arena) CmpColumn(n1 int, n2 int) *CmpResult {
	var cnt float64
	curr_el1 := n1
	curr_el2 := n2
	p1 := a.Get(curr_el1).Parent
	p2 := a.Get(curr_el2).Parent
	next_rate := CmpShallow(a.Get(curr_el1), a.Get(curr_el2))

	if next_rate == nil {
		return next_rate
	}

	if p1 == 0 || p2 == 0 {
		return next_rate
	}

	cnt++

	var total CmpResult
	for p1 != p2 {
		total.Append(*next_rate)

		curr_el1 = a.Get(curr_el1).Parent
		curr_el2 = a.Get(curr_el2).Parent
		p1 = a.Get(curr_el1).Parent
		p2 = a.Get(curr_el2).Parent
		next_rate = CmpShallow(a.Get(curr_el1), a.Get(curr_el2))

		cnt++

		if next_rate == nil {
			return next_rate
		}

		if p1 == 0 || p2 == 0 {
			return next_rate
		}
	}

	total.Append(*next_rate)
	return &total
}

func (a *Arena) XPath(n int, stopper int) string {
	parent := a.Get(n).Parent
	if parent > 0 && parent != stopper {
		return a.XPath(parent, stopper) + a.Get(n).String()
	}
	return a.Get(n).String()
}

func (a *Arena) PathArray(n int) []int {
	init := make([]int, 0)
	return a.pathArray(init, n)
}

// iterate all nodes up to root
func (a *Arena) pathArray(init []int, n int) []int {
	parent := a.Get(n).Parent
	if parent > 0 {
		return a.pathArray(append(init, n), parent)
	}
	return append(init, n)
}

func (a *Arena) Stringify(nodeId int) string {
	n := a.Get(nodeId)
	res := n.String() + "\n"
	for _, c := range n.Children {
		res += "  " + a.Stringify(c)
	}
	return res
}

func (a *Arena) StringifyInformation(nodeId int) string {
	n := a.Get(nodeId)
	var res string
	if n.Type == html.TextNode {
		res = n.Data + "\n"
	}

	if n.Type == html.ElementNode && n.Data == "img" {
		res = n.GetAttr("src") + "\n"
	}

	for _, c := range n.Children {
		res += "  " + a.StringifyInformation(c)
	}
	return res
}

func (a *Arena) Rate(nodeId int) int {
	r := nodePoints
	n := a.Get(nodeId)
	for _, attr := range n.Attr {
		// doubled for compatibility
		//r += (attrKeyPoints + len(attr.Val)*attrValPoints) * 2
		r += attrKeyPoints + len(attr.Val)*attrValPoints
	}

	for _, c := range n.Children {
		r += a.Rate(c)
	}
	return r
}

// Function returns list of "informative" endings
func (a *Arena) GetInformative(nId int) []int {
	var r []int
	if a.List[nId].isInformative() {
		r = append(r, nId)
		return r
	}

	for _, id := range a.List[nId].Children {
		r = append(r, a.GetInformative(id)...)
	}

	return r
}

func (a *Arena) WithNonequalFields(nId int) []int {
	var r []int
	if a.List[nId].DataArray != nil {
		if !a.List[nId].DataEqual() {
			r = append(r, nId)
			return r
		}
	}

	for _, id := range a.List[nId].Children {
		r = append(r, a.WithNonequalFields(id)...)
	}

	return r
}

func (a *Arena) WithFields(nId int) []int {
	var r []int
	if a.List[nId].DataArray != nil {
		r = append(r, nId)
		return r
	}

	for _, id := range a.List[nId].Children {
		r = append(r, a.WithFields(id)...)
	}

	return r
}

func (a *Arena) RenderString(id int) (string, error) {
	var s string
	buf := bytes.NewBufferString(s)
	err := a.Render(buf, id)
	return buf.String(), err
}

func (a *Arena) Render(w io.Writer, id int) error {
	if x, ok := w.(writer); ok {
		return a.render(x, id)
	}
	buf := bufio.NewWriter(w)
	if err := a.render(buf, id); err != nil {
		return err
	}
	return buf.Flush()
}

type writer interface {
	io.Writer
	io.ByteWriter
	WriteString(string) (int, error)
}

var plaintextAbort = errors.New("html: internal error (plaintext abort)")

func (a *Arena) render(w writer, id int) error {
	n := a.Get(id)

	// Render non-element nodes; these are the easy cases.
	switch n.Type {
	case html.ErrorNode:
		return errors.New("html: cannot render an ErrorNode node")
	case html.TextNode:
		return escape(w, n.Data)
	case html.DocumentNode:
		for _, c := range n.Children {
			if err := a.render(w, c); err != nil {
				return err
			}
		}
		return nil
	case html.ElementNode:
		// No-op.
	case html.CommentNode:
		if _, err := w.WriteString("<!--"); err != nil {
			return err
		}
		if _, err := w.WriteString(n.Data); err != nil {
			return err
		}
		if _, err := w.WriteString("-->"); err != nil {
			return err
		}
		return nil
	case html.DoctypeNode:
		if _, err := w.WriteString("<!DOCTYPE "); err != nil {
			return err
		}
		if _, err := w.WriteString(n.Data); err != nil {
			return err
		}
		if n.Attr != nil {
			var p, s string
			for _, a := range n.Attr {
				switch a.Key {
				case "public":
					p = a.Val
				case "system":
					s = a.Val
				}
			}
			if p != "" {
				if _, err := w.WriteString(" PUBLIC "); err != nil {
					return err
				}
				if err := writeQuoted(w, p); err != nil {
					return err
				}
				if s != "" {
					if err := w.WriteByte(' '); err != nil {
						return err
					}
					if err := writeQuoted(w, s); err != nil {
						return err
					}
				}
			} else if s != "" {
				if _, err := w.WriteString(" SYSTEM "); err != nil {
					return err
				}
				if err := writeQuoted(w, s); err != nil {
					return err
				}
			}
		}
		return w.WriteByte('>')
	default:
		return errors.New("html: unknown node type")
	}

	// Render the <xxx> opening tag.
	if err := w.WriteByte('<'); err != nil {
		return err
	}
	if _, err := w.WriteString(n.Data); err != nil {
		return err
	}
	for _, a := range n.Attr {
		if err := w.WriteByte(' '); err != nil {
			return err
		}
		if a.Namespace != "" {
			if _, err := w.WriteString(a.Namespace); err != nil {
				return err
			}
			if err := w.WriteByte(':'); err != nil {
				return err
			}
		}
		if _, err := w.WriteString(a.Key); err != nil {
			return err
		}
		if _, err := w.WriteString(`="`); err != nil {
			return err
		}
		if err := escape(w, a.Val); err != nil {
			return err
		}
		if err := w.WriteByte('"'); err != nil {
			return err
		}
	}
	if voidElements[n.Data] {
		if len(n.Children) > 0 {
			return fmt.Errorf("html: void element <%s> has child nodes", n.Data)
		}
		_, err := w.WriteString("/>")
		return err
	}
	if err := w.WriteByte('>'); err != nil {
		return err
	}

	// Add initial newline where there is danger of a newline beging ignored.
	for _, indx := range n.Children {
		c := a.Get(indx)
		if c.Type == html.TextNode && strings.HasPrefix(c.Data, "\n") {
			switch n.Data {
			case "pre", "listing", "textarea":
				if err := w.WriteByte('\n'); err != nil {
					return err
				}
			}
		}
	}

	// Render any child nodes.
	switch n.Data {
	case "iframe", "noembed", "noframes", "noscript", "plaintext", "script", "style", "xmp":
		for _, indx := range n.Children {
			c := a.Get(indx)
			if c.Type == html.TextNode {
				if _, err := w.WriteString(c.Data); err != nil {
					return err
				}
			} else {
				if err := a.render(w, indx); err != nil {
					return err
				}
			}
		}
		if n.Data == "plaintext" {
			// Don't render anything else. <plaintext> must be the
			// last element in the file, with no closing tag.
			return plaintextAbort
		}
	default:
		for _, indx := range n.Children {
			if err := a.render(w, indx); err != nil {
				return err
			}
		}
	}

	// Render the </xxx> closing tag.
	if _, err := w.WriteString("</"); err != nil {
		return err
	}
	if _, err := w.WriteString(n.Data); err != nil {
		return err
	}
	return w.WriteByte('>')
}

// lower lower-cases the A-Z bytes in b in-place, so that "aBc" becomes "abc".
func lower(b []byte) []byte {
	for i, c := range b {
		if 'A' <= c && c <= 'Z' {
			b[i] = c + 'a' - 'A'
		}
	}
	return b
}

const escapedChars = "&'<>\"\r"

func escape(w writer, s string) error {
	i := strings.IndexAny(s, escapedChars)
	for i != -1 {
		if _, err := w.WriteString(s[:i]); err != nil {
			return err
		}
		var esc string
		switch s[i] {
		case '&':
			esc = "&amp;"
		case '\'':
			// "&#39;" is shorter than "&apos;" and apos was not in HTML until HTML5.
			esc = "&#39;"
		case '<':
			esc = "&lt;"
		case '>':
			esc = "&gt;"
		case '"':
			// "&#34;" is shorter than "&quot;".
			esc = "&#34;"
		case '\r':
			esc = "&#13;"
		default:
			panic("unrecognized escape character")
		}
		s = s[i+1:]
		if _, err := w.WriteString(esc); err != nil {
			return err
		}
		i = strings.IndexAny(s, escapedChars)
	}
	_, err := w.WriteString(s)
	return err
}

// writeQuoted writes s to w surrounded by quotes. Normally it will use double
// quotes, but if s contains a double quote, it will use single quotes.
// It is used for writing the identifiers in a doctype declaration.
// In valid HTML, they can't contain both types of quotes.
func writeQuoted(w writer, s string) error {
	var q byte = '"'
	if strings.Contains(s, `"`) {
		q = '\''
	}
	if err := w.WriteByte(q); err != nil {
		return err
	}
	if _, err := w.WriteString(s); err != nil {
		return err
	}
	if err := w.WriteByte(q); err != nil {
		return err
	}
	return nil
}

// Section 12.1.2, "Elements", gives this list of void elements. Void elements
// are those that can't have any contents.
var voidElements = map[string]bool{
	"area":    true,
	"base":    true,
	"br":      true,
	"col":     true,
	"command": true,
	"embed":   true,
	"hr":      true,
	"img":     true,
	"input":   true,
	"keygen":  true,
	"link":    true,
	"meta":    true,
	"param":   true,
	"source":  true,
	"track":   true,
	"wbr":     true,
}
