package podcast

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
)

func GetRandomizedAds[T any](ads []T, r *http.Request) ([]T, int64) {
	randomisedAds := make([]T, len(ads))
	copy(randomisedAds, ads)
	randSeed := rand.Int63()

	seedCookie, err := r.Cookie("seed")
	if err == nil {
		fmt.Println("Seed cookie found:", seedCookie.Value)
		seed, err := strconv.ParseInt(seedCookie.Value, 10, 64)
		if err != nil {
			fmt.Println("Error parsing seed cookie: ", err)
		}
		randSeed = int64(seed)
	} else {
		fmt.Println("No seed cookie found")
	}

	rsrc := rand.New(rand.NewSource(randSeed))
	rsrc.Shuffle(len(randomisedAds), func(i, j int) {
		randomisedAds[i], randomisedAds[j] = randomisedAds[j], randomisedAds[i]
	})

	return randomisedAds, randSeed
}

func SetSeedCookie(w http.ResponseWriter, randSeed int64) {
	http.SetCookie(w, &http.Cookie{
		Name:  "seed",
		Value: fmt.Sprintf("%d", randSeed),
		Path:  "/",
	})
}
