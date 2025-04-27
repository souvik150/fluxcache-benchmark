package mixed

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/souvik150/fluxcache/pkg/fluxcache"

	"fluxcache-test/userpb"
)

var setSuccessCount int32
var getSuccessCount int32
var errorCount int32

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func MixedSetGetBenchmark(ctx context.Context, cache *fluxcache.FluxCache, totalOps, concurrency int) {
	fmt.Println("Starting MIXED SET/GET benchmark...")

	start := time.Now()

	var wg sync.WaitGroup
	opCh := make(chan int, totalOps)

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for id := range opCh {
				action := rand.Intn(2) // 0 = SET, 1 = GET

				key := fmt.Sprintf("user:%d", rand.Intn(totalOps)) // random user id

				if action == 0 {
					// SET operation
					user := &userpb.User{
						Id:                 fmt.Sprintf("user-%d", id),
						Name:               fmt.Sprintf("UserName-%d", id),
						Age:                int32(rand.Intn(100)),
						Email:              fmt.Sprintf("user%d@example.com", id),
						Phone:              fmt.Sprintf("+123456789%03d", id),
						Address:            randomString(20),
						City:               "CityName",
						State:              "StateName",
						Country:            "CountryName",
						ZipCode:            fmt.Sprintf("%05d", rand.Intn(99999)),
						Gender:             "Other",
						Username:           fmt.Sprintf("username%d", id),
						PasswordHash:       randomString(32),
						ProfilePictureUrl:  "https://example.com/profile.jpg",
						Bio:                randomString(100),
						IsActive:           true,
						IsVerified:         true,
						IsBanned:           false,
						FollowersCount:     int32(rand.Intn(10000)),
						FollowingCount:     int32(rand.Intn(5000)),
						PostsCount:         int32(rand.Intn(200)),
						LikesCount:         int32(rand.Intn(10000)),
						CommentsCount:      int32(rand.Intn(5000)),
						SharesCount:        int32(rand.Intn(1000)),
						WebsiteUrl:         "https://userwebsite.com",
						Company:            "TestCorp",
						JobTitle:           "Engineer",
						Education:          "Bachelor's Degree",
						Hometown:           "Hometown",
						RelationshipStatus: "Single",
						Interests:          "Coding, Reading",
						LanguagesSpoken:    "English, Spanish",
						Timezone:           "UTC",
						CreatedAt:          time.Now().Format(time.RFC3339),
						UpdatedAt:          time.Now().Format(time.RFC3339),
						Skills:             []string{"Go", "Redis", "Docker"},
						Hobbies:            []string{"Gaming", "Cooking"},
						FavoriteBooks:      []string{"Book1", "Book2"},
						FavoriteMovies:     []string{"Movie1", "Movie2"},
						FavoriteMusic:      []string{"Rock", "Jazz"},
						FavoriteFoods:      []string{"Pizza", "Burger"},
						VisitedCountries:   []string{"USA", "Japan"},
						Certifications:     []string{"Cert1", "Cert2"},
						HasPremium:         true,
						TwoFactorEnabled:   true,
						LastLogin:          time.Now().Format(time.RFC3339),
						LastIp:             "192.168.1.1",
						LoginAttempts:      int32(rand.Intn(10)),
						EmailNotificationsEnabled: true,
						SmsNotificationsEnabled:   false,
						PushNotificationsEnabled:  true,
						ThemePreference:    "dark",
						Locale:             "en-US",
						Currency:           "USD",
						AccountBalance:     rand.Float64() * 10000,
						LifetimeEarnings:   rand.Float64() * 50000,
						LifetimeSpent:      rand.Float64() * 25000,
						Devices:            []string{"iPhone", "MacBook"},
						SessionTokens:      []string{randomString(20)},
						SecurityQuestions:  []string{"Mother's maiden name?", "First pet?"},
						BlockedUsers:       []string{"user999", "user888"},
						MutedUsers:         []string{"user777", "user666"},
						FavoritePosts:      []string{"post123", "post456"},
						SavedPosts:         []string{"post789", "post101"},
						MarketingOptIn:     true,
						BetaProgramOptIn:   false,
						ReferralCode:       randomString(8),
						ReferredBy:         "user100",
						ReferredUsers:      []string{"user101", "user102"},
						AdminAccess:        false,
						ModeratorAccess:    false,
						AssignedGroups:     []string{"group1", "group2"},
						CreatedGroups:      []string{"group3"},
						JoinedGroups:       []string{"group1", "group2", "group3"},
						EventsAttending:    []string{"event1", "event2"},
						EventsHosting:      []string{"event3"},
						Awards:             []string{"Award1", "Award2"},
						Badges:             []string{"Badge1"},
						RecentSearches:     []string{"search1", "search2"},
						RecentViews:        []string{"view1", "view2"},
						WishlistItems:      []string{"item1", "item2"},
						PurchasedItems:     []string{"item3", "item4"},
						CartItems:          []string{"item5"},
						OrderHistory:       []string{"order1", "order2"},
						Subscriptions:      []string{"sub1"},
						Playlists:          []string{"playlist1", "playlist2"},
						FavoriteArtists:    []string{"artist1", "artist2"},
						FavoriteAuthors:    []string{"author1"},
						FavoritePlaces:     []string{"place1", "place2"},
						Pets:               []string{"dog", "cat"},
						FamilyMembers:      []string{"father", "mother"},
						EmergencyContacts:  []string{"contact1"},
						InsurancePolicies:  []string{"policy1"},
						MedicalRecords:     []string{"record1"},
						Vehicles:           []string{"car1", "bike1"},
						PropertiesOwned:    []string{"house1", "apartment1"},
						Investments:        []string{"investment1"},
						BankAccounts:       []string{"account1"},
						CreditCards:        []string{"card1", "card2"},
						Loans:              []string{"loan1"},
					}
					err := cache.SetProto(key, user)
					if err != nil {
						fmt.Printf("Worker %d: Error setting user %d: %v\n", workerID, id, err)
						atomic.AddInt32(&errorCount, 1)
						continue
					}
					atomic.AddInt32(&setSuccessCount, 1)
				} else {
					var user userpb.User
					err := cache.GetProto(key, &user)
					if err != nil {
						if errors.Is(err, redis.Nil) {
							continue
						}
						atomic.AddInt32(&errorCount, 1)
						continue
					}
					if user.Id == "" {
						fmt.Printf("Worker %d: Empty ID after get for user %d!\n", workerID, id)
						atomic.AddInt32(&errorCount, 1)
					} else {
						atomic.AddInt32(&getSuccessCount, 1)
					}
				}
			}
		}(i)
	}

	// Feed operations
	for i := 0; i < totalOps; i++ {
		opCh <- i
	}
	close(opCh)

	wg.Wait()

	duration := time.Since(start)
	fmt.Println("MIXED benchmark finished")
	fmt.Printf("Total ops: %d in %v (%.2f ops/sec)\n", totalOps, duration, float64(totalOps)/duration.Seconds())
	fmt.Printf("Set successes: %d, Get successes: %d, Errors: %d\n", setSuccessCount, getSuccessCount, errorCount)
}
