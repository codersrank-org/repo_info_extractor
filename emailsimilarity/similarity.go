package emailsimilarity

import (
	"math"
	"strings"

	"github.com/iancoleman/orderedmap"
)

var filterDomains = []string{"@gmail.com", ".local"}

// FindSimilarEmails returns similar emails to seed, from a given email list
func FindSimilarEmails(seeds, emails []string) []string {

	similarEmails := make([]string, 0, len(emails))

	for _, seed := range seeds {
		for _, email := range emails {
			similarity := calculateSimilarity(seed, email)
			if similarity > 0.55 {
				similarEmails = append(similarEmails, email)
			}
		}
	}

	return similarEmails
}

func calculateSimilarity(seed, email string) float64 {

	seedNgrams := getNgrams(seed, 2)
	emailNgrams := getNgrams(email, 2)

	ngramList := getNgramFullList(seedNgrams, emailNgrams)

	seedVec := getNgramVector(seedNgrams, ngramList)
	emailVec := getNgramVector(emailNgrams, ngramList)

	similarity := dotProduct(seedVec, emailVec) / (math.Sqrt(sum(square(seedVec))) * math.Sqrt(sum(square(emailVec))))
	return similarity
}

func getNgrams(word string, n int) *orderedmap.OrderedMap {
	ngrams := orderedmap.New()
	ngrams.SetEscapeHTML(false) // Default is true for serialization, but we don't need it

	for _, f := range filterDomains {
		word = strings.TrimRight(word, f)
	}

	for i := 0; i < len(word)-n+1; i++ {
		ngram := word[i : i+n]

		if val, ok := ngrams.Get(ngram); ok {
			ngrams.Set(ngram, val.(int)+1)
		} else {
			ngrams.Set(ngram, 1)
		}
	}

	return ngrams
}

// Just merging two maps, we won't use values here, only keys
func getNgramFullList(ngram1, ngram2 *orderedmap.OrderedMap) *orderedmap.OrderedMap {
	fullList := orderedmap.New()
	fullList.SetEscapeHTML(false)

	for _, v := range ngram1.Keys() {
		fullList.Set(v, 0)
	}

	for _, v := range ngram2.Keys() {
		fullList.Set(v, 0)
	}

	return fullList
}

func getNgramVector(ngram, ngramList *orderedmap.OrderedMap) []int {
	vector := make([]int, 0, len(ngramList.Keys()))
	for _, n := range ngramList.Keys() {
		if val, ok := ngram.Get(n); ok {
			vector = append(vector, val.(int))
		} else {
			vector = append(vector, 0)
		}
	}
	return vector
}

func dotProduct(x, y []int) float64 {
	total := 0.0
	for i := 0; i < len(x) && i < len(y); i++ {
		total += float64(x[i] * y[i])
	}

	return total
}

func square(x []int) []int {
	result := make([]int, len(x))
	for i := range x {
		result[i] = x[i] * x[i]
	}
	return result
}

func sum(x []int) float64 {
	sum := 0
	for i := range x {
		sum += x[i]
	}
	return float64(sum)
}
