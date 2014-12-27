package main

// Helper function to check if a string exists within a slice (array).
func stringInSlice(a string, list []string) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}
