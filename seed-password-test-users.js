// Seed script to create multiple test users with different password states
// Useful for testing password expiry functionality

// Connect to the users database
db = db.getSiblingDB('users');

// Bcrypt hash for "Test123!"
const passwordHash = "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy";

// Helper function to create date X days ago
function daysAgo(days) {
    const date = new Date();
    date.setDate(date.getDate() - days);
    return date;
}

// Helper function to create date X minutes ago
function minutesAgo(minutes) {
    const date = new Date();
    date.setMinutes(date.getMinutes() - minutes);
    return date;
}

// Test users with different password states
const testUsers = [
    {
        username: "fresh_password_user",
        email: "fresh@test.com",
        first_name: "Fresh",
        last_name: "Password",
        password_changed_at: new Date(), // Just now - FRESH
        description: "Lozinka UPRAVO promenjena"
    },
    {
        username: "week_old_password",
        email: "weekold@test.com",
        first_name: "Week",
        last_name: "Old",
        password_changed_at: daysAgo(7), // 7 days ago - OK
        description: "Lozinka stara 7 dana - OK"
    },
    {
        username: "month_old_password",
        email: "monthold@test.com",
        first_name: "Month",
        last_name: "Old",
        password_changed_at: daysAgo(30), // 30 days ago - OK
        description: "Lozinka stara 30 dana - OK"
    },
    {
        username: "expiring_soon_user",
        email: "expiring@test.com",
        first_name: "Expiring",
        last_name: "Soon",
        password_changed_at: daysAgo(55), // 55 days ago - EXPIRING SOON (5 days left)
        description: "Lozinka istice za 5 dana - UPOZORENJE"
    },
    {
        username: "just_expired_user",
        email: "justexpired@test.com",
        first_name: "Just",
        last_name: "Expired",
        password_changed_at: daysAgo(61), // 61 days ago - JUST EXPIRED
        description: "Lozinka UPRAVO ISTEKLA (61 dan)"
    },
    {
        username: "long_expired_user",
        email: "longexpired@test.com",
        first_name: "Long",
        last_name: "Expired",
        password_changed_at: daysAgo(120), // 120 days ago - LONG EXPIRED
        description: "Lozinka DAVNO ISTEKLA (120 dana)"
    },
    {
        username: "demo_minutes_expired",
        email: "minutesexpired@test.com",
        first_name: "Minutes",
        last_name: "Expired",
        password_changed_at: minutesAgo(5), // 5 minutes ago - for demo mode testing
        description: "Za DEMO mod (PASSWORD_MAX_AGE_MINUTES=2) - istekla"
    }
];

print("\n============================================");
print("   CREATING PASSWORD TEST USERS");
print("============================================\n");

// Delete existing test users and insert new ones
testUsers.forEach(user => {
    // Delete if exists
    db.users.deleteOne({ username: user.username });

    // Create full user document
    const fullUser = {
        _id: ObjectId(),
        username: user.username,
        email: user.email,
        password_hash: passwordHash,
        first_name: user.first_name,
        last_name: user.last_name,
        role: "regular",
        email_verified: true,
        email_verification_token: "",
        email_verification_token_exp: new Date(0),
        password_reset_token: "",
        password_reset_token_exp: new Date(0),
        magic_link_token: "",
        magic_link_token_exp: new Date(0),
        password_changed_at: user.password_changed_at,
        created_at: user.password_changed_at,
        updated_at: user.password_changed_at,
        failed_login_attempts: 0,
        last_failed_login: new Date(0),
        locked_until: new Date(0)
    };

    db.users.insertOne(fullUser);

    // Calculate password age
    const ageMs = new Date() - user.password_changed_at;
    const ageDays = Math.floor(ageMs / (1000 * 60 * 60 * 24));
    const ageMinutes = Math.floor(ageMs / (1000 * 60));

    print("+ " + user.username);
    print("  Email:    " + user.email);
    print("  Password: Test123!");
    print("  Age:      " + ageDays + " dana (" + ageMinutes + " minuta)");
    print("  Status:   " + user.description);
    print("");
});

print("============================================");
print("   TEST SCENARIOS");
print("============================================\n");

print("SCENARIO 1: Normalan mod (PASSWORD_MAX_AGE_DAYS=60)");
print("------------------------------------------------");
print("  fresh_password_user     -> LOGIN USPEŠAN");
print("  week_old_password       -> LOGIN USPEŠAN");
print("  month_old_password      -> LOGIN USPEŠAN");
print("  expiring_soon_user      -> LOGIN USPEŠAN (ali blizu isteka)");
print("  just_expired_user       -> LOGIN BLOKIRAN (PASSWORD_EXPIRED)");
print("  long_expired_user       -> LOGIN BLOKIRAN (PASSWORD_EXPIRED)");
print("");

print("SCENARIO 2: Demo mod (PASSWORD_MAX_AGE_MINUTES=2)");
print("------------------------------------------------");
print("  fresh_password_user     -> LOGIN USPEŠAN");
print("  demo_minutes_expired    -> LOGIN BLOKIRAN (starija od 2 min)");
print("  Svi ostali              -> LOGIN BLOKIRAN");
print("");

print("============================================");
print("   HOW TO RUN");
print("============================================\n");

print("1. Pokreni MongoDB kontejner:");
print("   docker-compose up -d mongodb-users");
print("");
print("2. Izvrši ovu skriptu:");
print("   mongosh \"mongodb://admin:admin123@localhost:27017/users?authSource=admin\" seed-password-test-users.js");
print("");
print("3. Za DEMO mod, dodaj env varijablu u docker-compose.yml:");
print("   PASSWORD_MAX_AGE_MINUTES: 2");
print("");
print("4. Restartuj users-service:");
print("   docker-compose restart users-service");
print("");

print("============================================");
print("Total users in database: " + db.users.countDocuments());
print("============================================\n");
