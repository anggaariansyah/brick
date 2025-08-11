package main

import (
	"context"
	"encoding/csv"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

/*
Full-detail Tokopedia scraper (category: Handphone)
Outputs CSV: tokopedia_handphone_full_100.csv
Fields: Name, Description, ImageLink, Price, Rating, Merchant, URL
*/

func randSleep(minMs, maxMs int) {
	d := time.Duration(minMs+rand.Intn(maxMs-minMs)) * time.Millisecond
	time.Sleep(d)
}

func safeTrim(s interface{}) string {
	if s == nil {
		return ""
	}
	if str, ok := s.(string); ok {
		return strings.TrimSpace(str)
	}
	return ""
}

func main() {
	rand.Seed(time.Now().UnixNano())

	outFile := "tokopedia_handphone_full_100.csv"
	f, err := os.Create(outFile)
	if err != nil {
		log.Fatalf("create csv: %v", err)
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()
	if err := w.Write([]string{"Name", "Description", "ImageLink", "Price", "Rating", "Merchant", "URL"}); err != nil {
		log.Fatalf("write header: %v", err)
	}

	// chromedp allocator/options
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("disable-http2", true),
		chromedp.Flag("headless", false), // set true to run headless
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("start-maximized", true),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) "+
			"AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120 Safari/537.36"),
	)
	allocCtx, cancelAlloc := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancelAlloc()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// overall timeout for entire run (adjust if needed)
	ctx, cancel = context.WithTimeout(ctx, 30*time.Minute)
	defer cancel()

	categoryBase := "https://www.tokopedia.com/p/handphone-tablet/handphone?page="
	target := 100
	maxPages := 50        // safety cap (increase if needed)
	scrollsPerPage := 6   // more scrolls -> more links per page
	waitPerScroll := 1100 // ms

	links := make([]string, 0, target)
	seen := map[string]struct{}{}

	log.Println("Collecting product links from category pages...")

pageLoop:
	for p := 1; p <= maxPages && len(links) < target; p++ {
		pageURL := categoryBase + strconv.Itoa(p)
		log.Printf("-> load page %d : %s", p, pageURL)

		// navigate and small waits
		err := chromedp.Run(ctx,
			chromedp.Navigate(pageURL),
			chromedp.Sleep(1200*time.Millisecond),
			chromedp.WaitReady("body", chromedp.ByQuery),
		)
		if err != nil {
			log.Printf("navigate page %d error: %v (continue)", p, err)
		}

		// scroll several times to trigger lazy loading
		for s := 0; s < scrollsPerPage; s++ {
			if err := chromedp.Run(ctx, chromedp.Evaluate(`window.scrollBy(0, window.innerHeight);`, nil)); err != nil {
				log.Printf("scroll error: %v", err)
			}
			time.Sleep(time.Duration(waitPerScroll) * time.Millisecond)
		}

		// collect anchors (prefer stable data-testid)
		var hrefs []string
		js := `
		(function(){
			let out = [];
			let nodes = document.querySelectorAll('a[data-testid="lnkProductContainer"]');
			if(nodes.length === 0){
				nodes = document.querySelectorAll('a[href*="/p/"], a[href*="/product/"]');
			}
			nodes.forEach(n => {
				try {
					let href = n.href || n.getAttribute('href');
					if(!href) return;
					if(href.startsWith('/')) href = location.origin + href;
					out.push(href.split('?')[0]);
				} catch(e){}
			});
			return Array.from(new Set(out));
		})()
		`
		if err := chromedp.Run(ctx, chromedp.Evaluate(js, &hrefs)); err != nil {
			log.Printf("collect links eval error page %d: %v", p, err)
			continue
		}

		for _, h := range hrefs {
			if len(links) >= target {
				break pageLoop
			}
			if h == "" {
				continue
			}
			if _, ok := seen[h]; ok {
				continue
			}
			seen[h] = struct{}{}
			links = append(links, h)
		}
		log.Printf("After page %d: collected total links = %d", p, len(links))
		// polite pause between pages
		randSleep(500, 1200)
	}

	if len(links) == 0 {
		log.Fatal("no product links found; try increase scrolls or debug HTML")
	}
	if len(links) > target {
		links = links[:target]
	}
	log.Printf("Collected %d links. Begin full-detail extraction...", len(links))

	// js snippet: multiple selector fallbacks to extract fields from product page
	jsExtract := `(function(){
		function tx(sel){ try { let e=document.querySelector(sel); return e ? e.innerText.trim() : ""; } catch(e){ return ""; } }
		function attr(sel, a){ try { let e=document.querySelector(sel); return e ? (e.getAttribute(a) || "") : ""; } catch(e){ return ""; } }

		let o = {};
		// name
		o.name = tx('h1[data-testid="lblPDPDetailProductName"]') || tx('h1') || tx('[data-testid="lblPDPTitle"]') || tx('[data-testid="spnSRPProdName"]') || "";
		// description (detailed)
		o.desc = tx('div[data-testid="lblPDPDesc"]') || tx('[data-testid="lblPDPDetailProductDesc"]') || tx('#description') || tx('div[class*="description"]') || "";
		// image
		o.image = attr('img[data-testid="PDPImageMain"]','src') || attr('img[data-testid="PDPHeroImage"]','src') || attr('img','src') || "";
		// price
		o.price = tx('div[data-testid="lblPDPDetailProductPrice"]') || tx('[data-testid="spnSRPProdPrice"]') || tx('div[class*="price"]') || "";
		// rating
		o.rating = tx('span[data-testid="lblPDPDetailProductRatingNumber"]') || tx('[data-testid="lblPDPDetailProductRating"]') || tx('[aria-label*="rating"]') || "";
		// merchant
		o.merchant = tx('a[data-testid="llbPDPFooterShopName"]') || tx('a[data-testid="lnkPDPFooterShopName"]') || tx('[data-testid="spnSRPProdShopName"]') || tx('a[href*="/shop/"]') || tx('div[class*="shop"]') || "";
		return o;
	})()`

	// iterate each link and extract detail
	for i, url := range links {
		log.Printf("Processing %d/%d : %s", i+1, len(links), url)

		// per-product timeout/context
		prodCtx, cancel := context.WithTimeout(ctx, 45*time.Second)
		var out map[string]interface{}
		err := chromedp.Run(prodCtx,
			chromedp.Navigate(url),
			chromedp.Sleep(1200*time.Millisecond),
			// wait for likely product element; don't fail if not present quickly
			chromedp.WaitReady("body", chromedp.ByQuery),
			// inside product page may require scroll to load description/images
			chromedp.Evaluate(`window.scrollTo(0, 200)`, nil),
			chromedp.Sleep(800*time.Millisecond),
			chromedp.Evaluate(`window.scrollTo(0, document.body.scrollHeight/3)`, nil),
			chromedp.Sleep(1000*time.Millisecond),
			chromedp.Evaluate(jsExtract, &out),
		)
		cancel()
		if err != nil {
			log.Printf("error extracting %s : %v (writing partial row)", url, err)
			// write a fallback row with URL so you can retry later
			_ = w.Write([]string{"", "N/A", "", "", "", "", url})
			w.Flush()
			// polite delay and continue
			randSleep(700, 1500)
			continue
		}

		name := safeTrim(out["name"])
		desc := safeTrim(out["desc"])
		image := safeTrim(out["image"])
		price := safeTrim(out["price"])
		rating := safeTrim(out["rating"])
		merchant := safeTrim(out["merchant"])

		// normalize rating: try to extract leading numeric like "4.8" else set "0"
		rating = strings.TrimSpace(rating)
		if rating == "" {
			rating = "0"
		} else {
			// remove parentheses, non-digit trailing text
			rating = strings.ReplaceAll(rating, "(", "")
			rating = strings.ReplaceAll(rating, ")", "")
			// take first token if multiple
			parts := strings.Fields(rating)
			if len(parts) > 0 {
				rating = parts[0]
			}
		}

		if desc == "" {
			desc = "N/A"
		}
		if name == "" {
			name = "N/A"
		}
		if price == "" {
			price = "N/A"
		}
		if merchant == "" {
			merchant = "N/A"
		}
		if image == "" {
			image = "N/A"
		}

		// write CSV row
		if err := w.Write([]string{name, desc, image, price, rating, merchant, url}); err != nil {
			log.Printf("csv write error: %v", err)
		}
		w.Flush()

		log.Printf("Saved %d: %s | %s | %s", i+1, name, price, merchant)

		// polite random delay
		randSleep(900, 1900)
	}

	log.Printf("Done — CSV: %s", outFile)
}
