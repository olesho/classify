package cluster

import (
	"fmt"
	"strings"
	"testing"

	"github.com/olesho/classify/arena"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

func TestComparator(t *testing.T) {
	text := `
<html>
	<body>
		
		<div class="text">
		   <p>Hacker
			  Noon elevates tech writing far and wide across the interwebs. With 25k+
			  stories in our library, we’ve cultivated rare publishing and technology
			  expertise. We believe tech infiltrates all aspects of life, so writers 
			  from all industries are welcome here. To the contributing writer, we 
			  pledge quality distribution, editorial expertise and community access.
		   </p>
		   <div class="cta">
			  <a href="https://community.hackernoon.com/session/sso?return_path=https%3A%2F%2Fapp.hackernoon.com" class="button">Submit a Tech Story</a>
		   </div>
		</div>
		
		<div class="text">
		   <p>Hacker
			  Noon partners with companies that build cool products and employ people
			  worth publishing. To sponsors, we currently offer sitewide billboard, 
			  brand-as-author stories, podcast placements, and  events. Sponsors keep 
			  Hacker Noon free for readers forever and add to the reading experience 
			  rather than distracting from it. Sponsoring Hacker Noon is great way to 
			  make a positive impression on tech’s brightest minds.
		   </p>
		   <div class="cta">
			  <a href="https://sponsor.hackernoon.com/" class="button">Sponsor Hacker Noon</a>
		   </div>
		</div>
		
		<div class="story-card long-title">
		   <a href="https://hackernoon.com/product-is-the-king-how-ukrainian-engineers-creating-home-security-products-jwbm3y5j">
			  <div class="img" style="background-image: url('https://hackernoon.com/drafts/731953yhh.png')"></div>
		   </a>
		   <div class="excerpt">
			  <div class="title">
				 <a href="https://hackernoon.com/product-is-the-king-how-ukrainian-engineers-creating-home-security-products-jwbm3y5j">Product is the King: How Ukrainian Engineers Are Creating Home Security Products</a>
			  </div>
			  <div class="bio">
				 <div class="flex">
					<div class="avatar" style="background-image: url(https://hackernoon.com/avatars/SxcJoG3rFNPmC7X5zjBbNJ6KrXA2.png)"></div>
					<div>
					   <a href="https://hackernoon.com/@amina" class="name">Amina</a>
					   <div class="published">March 28</div>
					</div>
				 </div>
			  </div>
		   </div>
		   <a class="tag" href="https://hackernoon.com/tagged/home-security">Home Security</a>
		</div>
		
		<div class="story-card">
		   <a href="https://hackernoon.com/blockchain-use-cases-cutting-through-the-hype-489l328r">
			  <div class="img" style="background-image: url('https://hackernoon.com/drafts/6i31338w.png')"></div>
		   </a>
		   <div class="excerpt">
			  <div class="title">
				 <a href="https://hackernoon.com/blockchain-use-cases-cutting-through-the-hype-489l328r">Blockchain Use Cases: Cutting Through the Hype</a>
			  </div>
			  <div class="bio">
				 <div class="flex">
					<div class="avatar" style="background-image: url(https://hackernoon.com/avatars/EBNEt8ZDkqPFfK65oD0aZyPRqBk1.png)"></div>
					<div>
					   <a href="https://hackernoon.com/@Carl%20Lang" class="name">Carl Lang</a>
					   <div class="published">March 26</div>
					</div>
				 </div>
			  </div>
		   </div>
		   <a class="tag" href="https://hackernoon.com/tagged/blockchain">Blockchain</a>
		</div>

	</body>
</html>
`
	a := assert.New(t)
	reader := strings.NewReader(text)
	n, err := html.Parse(reader)
	a.NoError(err)

	arena := arena.NewArena()
	arena.Load(*n)
	Init(arena)
	c := NewDefaultComparator(arena)

	idx1 := 4
	idx2 := 10
	idx3 := 16
	idx4 := 33

	node1 := arena.Get(idx1)
	node2 := arena.Get(idx2)
	node3 := arena.Get(idx3)
	node4 := arena.Get(idx4)

	fmt.Println(c.Cmp(node1, node2))
	fmt.Println(c.Cmp(node3, node4))
	fmt.Println(c.Cmp(node1, node4))
}
