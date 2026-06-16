package auth

import (
	"strings"
	"testing"
)

func TestPasswordHashing(t *testing.T) {
	t.Run("HashPassword", func(t *testing.T) {
		password := "testpassword123"
		
		hash, err := HashPassword(password)
		if err != nil {
			t.Fatalf("Failed to hash password: %v", err)
		}
		
		if hash == "" {
			t.Error("Hash should not be empty")
		}
		
		if hash == password {
			t.Error("Hash should not be the same as password")
		}
		
		// Check that hash contains expected Argon2id format
		if !strings.HasPrefix(hash, "$argon2id$") {
			t.Error("Hash should use Argon2id format")
		}
	})

	t.Run("VerifyPassword", func(t *testing.T) {
		password := "testpassword123"
		
		hash, err := HashPassword(password)
		if err != nil {
			t.Fatalf("Failed to hash password: %v", err)
		}
		
		// Correct password should verify
		valid, err := VerifyPassword(password, hash)
		if err != nil {
			t.Fatalf("Failed to check password hash: %v", err)
		}
		if !valid {
			t.Error("Correct password should verify against hash")
		}
		
		// Wrong password should not verify
		valid, err = VerifyPassword("wrongpassword", hash)
		if err != nil {
			t.Fatalf("Failed to check wrong password hash: %v", err)
		}
		if valid {
			t.Error("Wrong password should not verify against hash")
		}
		
		// Empty password should not verify
		valid, err = VerifyPassword("", hash)
		if err != nil {
			t.Fatalf("Failed to check empty password hash: %v", err)
		}
		if valid {
			t.Error("Empty password should not verify against hash")
		}
	})

	t.Run("HashUniqueness", func(t *testing.T) {
		password := "testpassword123"
		
		hash1, err := HashPassword(password)
		if err != nil {
			t.Fatalf("Failed to hash password first time: %v", err)
		}
		
		hash2, err := HashPassword(password)
		if err != nil {
			t.Fatalf("Failed to hash password second time: %v", err)
		}
		
		// Hashes should be different due to salt
		if hash1 == hash2 {
			t.Error("Two hashes of the same password should be different due to salt")
		}
		
		// Both should verify correctly
		valid, err := VerifyPassword(password, hash1)
		if err != nil {
			t.Fatalf("Failed to check first hash: %v", err)
		}
		if !valid {
			t.Error("First hash should verify correctly")
		}
		
		valid, err = VerifyPassword(password, hash2)
		if err != nil {
			t.Fatalf("Failed to check second hash: %v", err)
		}
		if !valid {
			t.Error("Second hash should verify correctly")
		}
	})

	t.Run("EmptyPassword", func(t *testing.T) {
		hash, err := HashPassword("")
		if err != nil {
			t.Fatalf("Failed to hash empty password: %v", err)
		}
		
		// Empty password should still produce a hash
		if hash == "" {
			t.Error("Empty password should still produce a hash")
		}
		
		// Empty password should verify against its hash
		valid, err := VerifyPassword("", hash)
		if err != nil {
			t.Fatalf("Failed to check empty password: %v", err)
		}
		if !valid {
			t.Error("Empty password should verify against its own hash")
		}
		
		// Non-empty password should not verify against empty password hash
		valid, err = VerifyPassword("nonempty", hash)
		if err != nil {
			t.Fatalf("Failed to check non-empty password: %v", err)
		}
		if valid {
			t.Error("Non-empty password should not verify against empty password hash")
		}
	})

	t.Run("LongPassword", func(t *testing.T) {
		// Test with very long password
		longPassword := strings.Repeat("a", 1000)
		
		hash, err := HashPassword(longPassword)
		if err != nil {
			t.Fatalf("Failed to hash long password: %v", err)
		}
		
		valid, err := VerifyPassword(longPassword, hash)
		if err != nil {
			t.Fatalf("Failed to check long password: %v", err)
		}
		if !valid {
			t.Error("Long password should verify against its hash")
		}
	})

	t.Run("SpecialCharacters", func(t *testing.T) {
		// Test with special characters
		specialPassword := "!@#$%^&*()_+-=[]{}|;':\",./<>?`~"
		
		hash, err := HashPassword(specialPassword)
		if err != nil {
			t.Fatalf("Failed to hash password with special characters: %v", err)
		}
		
		valid, err := VerifyPassword(specialPassword, hash)
		if err != nil {
			t.Fatalf("Failed to check special password: %v", err)
		}
		if !valid {
			t.Error("Password with special characters should verify against its hash")
		}
	})

	t.Run("UnicodePassword", func(t *testing.T) {
		// Test with Unicode characters
		unicodePassword := "пароль123密码🔐"
		
		hash, err := HashPassword(unicodePassword)
		if err != nil {
			t.Fatalf("Failed to hash Unicode password: %v", err)
		}
		
		valid, err := VerifyPassword(unicodePassword, hash)
		if err != nil {
			t.Fatalf("Failed to check Unicode password: %v", err)
		}
		if !valid {
			t.Error("Unicode password should verify against its hash")
		}
	})

	t.Run("InvalidHashFormat", func(t *testing.T) {
		password := "testpassword"
		
		// Test with invalid hash formats
		invalidHashes := []string{
			"",
			"invalid",
			"$argon2id$",
			"$argon2id$v=19$m=65536$t=3$p=2$",
			"$argon2id$v=19$m=65536$t=3$p=2$salt",
			"$md5$invalid",
			"plaintext",
		}
		
		for _, invalidHash := range invalidHashes {
			valid, err := VerifyPassword(password, invalidHash)
			if err == nil && valid {
				t.Errorf("Password should not verify against invalid hash: %s", invalidHash)
			}
		}
	})

	t.Run("DecodeHashFunction", func(t *testing.T) {
		password := "testpassword123"
		hash, err := HashPassword(password)
		if err != nil {
			t.Fatalf("Failed to hash password: %v", err)
		}
		
		// Test decodeHash function
		_, salt, hashedPassword, err := decodeHash(hash)
		if err != nil {
			t.Fatalf("Failed to decode hash: %v", err)
		}
		
		if len(salt) == 0 {
			t.Error("Salt should not be empty")
		}
		
		if len(hashedPassword) == 0 {
			t.Error("Hashed password should not be empty")
		}
		
		// Test with invalid hash
		_, _, _, err = decodeHash("invalid")
		if err == nil {
			t.Error("decodeHash should fail with invalid hash")
		}
		
		// Test with wrong number of parts
		_, _, _, err = decodeHash("$argon2id$v=19$m=65536")
		if err == nil {
			t.Error("decodeHash should fail with incomplete hash")
		}
	})

	t.Run("HashConsistency", func(t *testing.T) {
		// Test that the same password always verifies correctly
		password := "consistencytest"
		
		hash, err := HashPassword(password)
		if err != nil {
			t.Fatalf("Failed to hash password: %v", err)
		}
		
		// Verify multiple times
		for i := 0; i < 10; i++ {
			valid, err := VerifyPassword(password, hash)
			if err != nil {
				t.Fatalf("Failed to check password on iteration %d: %v", i, err)
			}
			if !valid {
				t.Errorf("Password verification failed on iteration %d", i)
			}
		}
	})

	t.Run("CaseSensitivity", func(t *testing.T) {
		password := "CaseSensitive123"
		
		hash, err := HashPassword(password)
		if err != nil {
			t.Fatalf("Failed to hash password: %v", err)
		}
		
		// Correct case should verify
		valid, err := VerifyPassword(password, hash)
		if err != nil {
			t.Fatalf("Failed to check correct case password: %v", err)
		}
		if !valid {
			t.Error("Correct case password should verify")
		}
		
		// Different case should not verify
		valid, err = VerifyPassword("casesensitive123", hash)
		if err != nil {
			t.Fatalf("Failed to check lowercase password: %v", err)
		}
		if valid {
			t.Error("Different case password should not verify")
		}
		
		valid, err = VerifyPassword("CASESENSITIVE123", hash)
		if err != nil {
			t.Fatalf("Failed to check uppercase password: %v", err)
		}
		if valid {
			t.Error("Different case password should not verify")
		}
	})
}