// Seed script to create a test user with an expired password
// This user can be used to demonstrate the password expiry audit feature

// Connect to the users database
db = db.getSiblingDB('users');

// Calculate date 90 days in the past (password expired)
const ninetyDaysAgo = new Date();
ninetyDaysAgo.setDate(ninetyDaysAgo.getDate() - 90);

// Test user with expired password
// Password: Test123! (hashed with bcrypt)
const expiredPasswordUser = {
    _id: ObjectId(),
    username: "expired_user",
    email: "expired@test.com",
    // Bcrypt hash for "Test123!"
    password_hash: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy",
    first_name: "Test",
    last_name: "ExpiredPassword",
    role: "regular",
    email_verified: true,
    email_verification_token: "",
    email_verification_token_exp: new Date(0),
    password_reset_token: "",
    password_reset_token_exp: new Date(0),
    magic_link_token: "",
    magic_link_token_exp: new Date(0),
    // PASSWORD CHANGED 90 DAYS AGO - EXPIRED!
    password_changed_at: ninetyDaysAgo,
    created_at: ninetyDaysAgo,
    updated_at: ninetyDaysAgo,
    failed_login_attempts: 0,
    last_failed_login: new Date(0),
    locked_until: new Date(0)
};

// Check if user already exists and delete it
const existingUser = db.users.findOne({ username: "expired_user" });
if (existingUser) {
    db.users.deleteOne({ username: "expired_user" });
    print("Deleted existing expired_user");
}

// Insert the user
db.users.insertOne(expiredPasswordUser);

print("\n============================================");
print("=== Test user with EXPIRED password created ===");
print("============================================");
print("");
print("Username:           expired_user");
print("Email:              expired@test.com");
print("Password:           Test123!");
print("Role:               regular");
print("Email verified:     true");
print("");
print("PASSWORD STATUS:    EXPIRED");
print("Password changed:   " + ninetyDaysAgo.toISOString());
print("Password age:       90 days (max allowed: 60 days)");
print("");
print("When this user tries to login, they will receive:");
print("  - HTTP 403 Forbidden");
print("  - Error code: PASSWORD_EXPIRED");
print("  - Message: 'Lozinka je istekla. Molimo resetujte vašu lozinku.'");
print("");
print("============================================");
print("");
print("Total users in database: " + db.users.countDocuments());
