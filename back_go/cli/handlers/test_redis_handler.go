package handlers

import (
	"csv-importer/api/cache"
	"database/sql"
	"fmt"
	"time"
)

func HandleTestRedis(db *sql.DB) {
	fmt.Println("🧪 Testing Redis connection...")

	config := cache.CacheConfig{
		Host:     "localhost",
		Port:     "6379",
		Password: "",
		DB:       0,
	}

	redisCache := cache.NewRedisCache(config)

	fmt.Print("📡 Testing connection... ")
	if err := redisCache.Ping(); err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
		return
	}
	fmt.Println("✅ Connected!")

	fmt.Print("📝 Testing SET... ")
	if err := redisCache.Set("test_key", "Hello Redis!", time.Minute*10); err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
		return
	}
	fmt.Println("✅ Set success!")

	fmt.Print("📖 Testing GET... ")
	var result string
	if err := redisCache.Get("test_key", &result); err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
		return
	}
	fmt.Printf("✅ Got: %s\n", result)

	fmt.Print("📋 Testing JSON SET... ")
	testData := map[string]interface{}{
		"entity_number": "123456789",
		"name":          "Test Company",
		"nace_code":     "62020",
		"contacts": map[string]string{
			"email": "test@example.be",
			"phone": "+32123456789",
		},
	}

	if err := redisCache.Set("search_session:test", testData, time.Hour); err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
		return
	}
	fmt.Println("✅ JSON set success!")

	fmt.Print("📊 Testing JSON GET... ")
	var retrievedData map[string]interface{}
	if err := redisCache.Get("search_session:test", &retrievedData); err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
		return
	}
	fmt.Printf("✅ Retrieved company: %s\n", retrievedData["name"])

	fmt.Print("🔍 Testing EXISTS... ")
	exists, err := redisCache.Exists("search_session:test")
	if err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
		return
	}
	fmt.Printf("✅ Key exists: %t\n", exists)

	fmt.Print("🗝️  Testing KEYS... ")
	keys, err := redisCache.GetKeys("*")
	if err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
		return
	}
	fmt.Printf("✅ Found %d keys\n", len(keys))

	fmt.Print("🗑️  Testing DELETE... ")
	if err := redisCache.Delete("test_key"); err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
		return
	}
	fmt.Println("✅ Deleted!")

	fmt.Println("\n🎉 All Redis tests passed!")
	fmt.Println("🚀 Ready for search funnel implementation!")

	redisCache.Close()
}
