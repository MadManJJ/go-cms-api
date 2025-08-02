package helpers

import (
	"math/rand"
	"strconv"
	"time"
)

func RandomState() string {
    // Create a new random source with current time as seed
    r := rand.New(rand.NewSource(time.Now().UnixNano()))
    
    // Generate a random number and convert to base 36
    num := r.Float64()
    str := strconv.FormatFloat(num, 'f', -1, 64)[2:] // Remove "0."
    
    // Convert to something more like base36
    const charset = "0123456789abcdefghijklmnopqrstuvwxyz"
    result := ""
    for i := 0; i < 13 && i < len(str); i++ {
        if str[i] >= '0' && str[i] <= '9' {
            result += string(charset[str[i]-'0'])
        }
    }
    
    // Fill remaining with random chars if needed
    for len(result) < 13 {
        result += string(charset[r.Intn(36)])
    }
    
    return result
}
