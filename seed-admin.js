// Connect to the users database
db = db.getSiblingDB('users');

// Clear all existing users
db.users.deleteMany({});
print("Cleared all existing users");

// Create hardcoded admin user
// Password: Admin123! (hashed with bcrypt)
const adminUser = {
    _id: ObjectId(),
    username: "admin",
    email: "admin@spotify.com",
    // Bcrypt hash for "Admin123!"
    password_hash: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy",
    first_name: "Admin",
    last_name: "User",
    role: "admin",
    email_verified: true,
    email_verification_token: "",
    email_verification_token_exp: new Date(0),
    password_reset_token: "",
    password_reset_token_exp: new Date(0),
    magic_link_token: "",
    magic_link_token_exp: new Date(0),
    password_changed_at: new Date(),
    created_at: new Date(),
    updated_at: new Date(),
    failed_login_attempts: 0,
    last_failed_login: new Date(0),
    locked_until: new Date(0)
};

db.users.insertOne(adminUser);

print("\n=== Admin user created successfully! ===");
print("Username: admin");
print("Email: admin@spotify.com");
print("Password: Admin123!");
print("Role: admin");
print("Email verified: true");
print("\nTotal users in database: " + db.users.countDocuments());
