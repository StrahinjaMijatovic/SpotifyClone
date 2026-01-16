// Seed script to create a demo user for password expiry demonstration
// This creates a user whose password is about to expire (for live demo)

// Connect to the users database
db = db.getSiblingDB('users');

// Calculate date 3 minutes in the past
// When using PASSWORD_MAX_AGE_MINUTES=2, this user's password will be expired
const threeMinutesAgo = new Date();
threeMinutesAgo.setMinutes(threeMinutesAgo.getMinutes() - 3);

// Demo user for simulation
// Password: Demo123! (same hash as Admin123!)
const demoUser = {
    _id: ObjectId(),
    username: "demo_user",
    email: "demo@test.com",
    // Bcrypt hash for "Test123!"
    password_hash: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy",
    first_name: "Demo",
    last_name: "User",
    role: "regular",
    email_verified: true,
    email_verification_token: "",
    email_verification_token_exp: new Date(0),
    password_reset_token: "",
    password_reset_token_exp: new Date(0),
    magic_link_token: "",
    magic_link_token_exp: new Date(0),
    // PASSWORD CHANGED 3 MINUTES AGO - for demo
    password_changed_at: threeMinutesAgo,
    created_at: threeMinutesAgo,
    updated_at: threeMinutesAgo,
    failed_login_attempts: 0,
    last_failed_login: new Date(0),
    locked_until: new Date(0)
};

// Check if user already exists and delete it
const existingUser = db.users.findOne({ username: "demo_user" });
if (existingUser) {
    db.users.deleteOne({ username: "demo_user" });
    print("Deleted existing demo_user");
}

// Insert the user
db.users.insertOne(demoUser);

print("\n============================================");
print("=== Demo user for password expiry created ===");
print("============================================");
print("");
print("Username:           demo_user");
print("Email:              demo@test.com");
print("Password:           Test123!");
print("Role:               regular");
print("Email verified:     true");
print("");
print("Password changed:   " + threeMinutesAgo.toISOString());
print("Password age:       ~3 minutes");
print("");
print("============================================");
print("   DEMO INSTRUCTIONS");
print("============================================");
print("");
print("1. Set environment variable for short expiry period:");
print("   PASSWORD_MAX_AGE_MINUTES=2");
print("");
print("2. Restart the users-service");
print("");
print("3. Try to login with demo_user/Test123!");
print("   -> Login will be BLOCKED (password older than 2 min)");
print("");
print("4. The system will return PASSWORD_EXPIRED error");
print("");
print("============================================");
print("");
print("Total users in database: " + db.users.countDocuments());
