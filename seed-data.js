// Connect to the content database
db = db.getSiblingDB('content');

// Clear existing data
db.genres.deleteMany({});
db.artists.deleteMany({});
db.albums.deleteMany({});
db.songs.deleteMany({});

print("Cleared existing data");

// Insert Genres
const genres = [
    { _id: ObjectId(), name: "Hip Hop", description: "Hip hop and rap music", created_at: new Date() },
    { _id: ObjectId(), name: "Rock", description: "Rock music", created_at: new Date() },
    { _id: ObjectId(), name: "Pop", description: "Pop music", created_at: new Date() },
    { _id: ObjectId(), name: "Electronic", description: "Electronic and EDM", created_at: new Date() },
    { _id: ObjectId(), name: "R&B", description: "Rhythm and Blues", created_at: new Date() }
];

db.genres.insertMany(genres);
print("Inserted " + genres.length + " genres");

// Insert Artists
const artists = [
    {
        _id: ObjectId(),
        name: "Eminem",
        biography: "Marshall Bruce Mathers III, known professionally as Eminem, is an American rapper, songwriter, and record producer. He is one of the best-selling music artists of all time.",
        genres: [genres[0]._id], // Hip Hop
        created_at: new Date(),
        updated_at: new Date()
    },
    {
        _id: ObjectId(),
        name: "The Beatles",
        biography: "The Beatles were an English rock band formed in Liverpool in 1960. They are regarded as the most influential band of all time.",
        genres: [genres[1]._id, genres[2]._id], // Rock, Pop
        created_at: new Date(),
        updated_at: new Date()
    },
    {
        _id: ObjectId(),
        name: "Daft Punk",
        biography: "Daft Punk were a French electronic music duo formed in 1993 in Paris. They achieved popularity in the late 1990s as part of the French house movement.",
        genres: [genres[3]._id], // Electronic
        created_at: new Date(),
        updated_at: new Date()
    },
    {
        _id: ObjectId(),
        name: "Beyoncé",
        biography: "Beyoncé Giselle Knowles-Carter is an American singer, songwriter, and actress. She rose to fame in the late 1990s as the lead singer of Destiny's Child.",
        genres: [genres[2]._id, genres[4]._id], // Pop, R&B
        created_at: new Date(),
        updated_at: new Date()
    },
    {
        _id: ObjectId(),
        name: "Kendrick Lamar",
        biography: "Kendrick Lamar Duckworth is an American rapper and songwriter. He is often cited as one of the most influential rappers of his generation.",
        genres: [genres[0]._id], // Hip Hop
        created_at: new Date(),
        updated_at: new Date()
    }
];

db.artists.insertMany(artists);
print("Inserted " + artists.length + " artists");

// Insert Albums
const albums = [
    {
        _id: ObjectId(),
        name: "The Eminem Show",
        date: new Date("2002-05-26"),
        genre: genres[0]._id, // Hip Hop
        artists: [artists[0]._id], // Eminem
        created_at: new Date(),
        updated_at: new Date()
    },
    {
        _id: ObjectId(),
        name: "Abbey Road",
        date: new Date("1969-09-26"),
        genre: genres[1]._id, // Rock
        artists: [artists[1]._id], // The Beatles
        created_at: new Date(),
        updated_at: new Date()
    },
    {
        _id: ObjectId(),
        name: "Random Access Memories",
        date: new Date("2013-05-17"),
        genre: genres[3]._id, // Electronic
        artists: [artists[2]._id], // Daft Punk
        created_at: new Date(),
        updated_at: new Date()
    },
    {
        _id: ObjectId(),
        name: "Lemonade",
        date: new Date("2016-04-23"),
        genre: genres[4]._id, // R&B
        artists: [artists[3]._id], // Beyoncé
        created_at: new Date(),
        updated_at: new Date()
    },
    {
        _id: ObjectId(),
        name: "DAMN.",
        date: new Date("2017-04-14"),
        genre: genres[0]._id, // Hip Hop
        artists: [artists[4]._id], // Kendrick Lamar
        created_at: new Date(),
        updated_at: new Date()
    }
];

db.albums.insertMany(albums);
print("Inserted " + albums.length + " albums");

// Insert Songs
const songs = [
    // The Eminem Show
    { _id: ObjectId(), name: "Without Me", duration: 290, genre: genres[0]._id, album: albums[0]._id, artists: [artists[0]._id], audio_url: "https://www.soundhelix.com/examples/mp3/SoundHelix-Song-1.mp3", created_at: new Date(), updated_at: new Date() },
    { _id: ObjectId(), name: "Cleanin' Out My Closet", duration: 297, genre: genres[0]._id, album: albums[0]._id, artists: [artists[0]._id], audio_url: "https://www.soundhelix.com/examples/mp3/SoundHelix-Song-2.mp3", created_at: new Date(), updated_at: new Date() },
    { _id: ObjectId(), name: "Sing for the Moment", duration: 339, genre: genres[0]._id, album: albums[0]._id, artists: [artists[0]._id], audio_url: "https://www.soundhelix.com/examples/mp3/SoundHelix-Song-3.mp3", created_at: new Date(), updated_at: new Date() },

    // Abbey Road
    { _id: ObjectId(), name: "Come Together", duration: 259, genre: genres[1]._id, album: albums[1]._id, artists: [artists[1]._id], audio_url: "https://www.soundhelix.com/examples/mp3/SoundHelix-Song-4.mp3", created_at: new Date(), updated_at: new Date() },
    { _id: ObjectId(), name: "Something", duration: 182, genre: genres[1]._id, album: albums[1]._id, artists: [artists[1]._id], audio_url: "https://www.soundhelix.com/examples/mp3/SoundHelix-Song-5.mp3", created_at: new Date(), updated_at: new Date() },
    { _id: ObjectId(), name: "Here Comes the Sun", duration: 185, genre: genres[1]._id, album: albums[1]._id, artists: [artists[1]._id], audio_url: "https://www.soundhelix.com/examples/mp3/SoundHelix-Song-6.mp3", created_at: new Date(), updated_at: new Date() },

    // Random Access Memories
    { _id: ObjectId(), name: "Get Lucky", duration: 368, genre: genres[3]._id, album: albums[2]._id, artists: [artists[2]._id], audio_url: "https://www.soundhelix.com/examples/mp3/SoundHelix-Song-7.mp3", created_at: new Date(), updated_at: new Date() },
    { _id: ObjectId(), name: "Instant Crush", duration: 337, genre: genres[3]._id, album: albums[2]._id, artists: [artists[2]._id], audio_url: "https://www.soundhelix.com/examples/mp3/SoundHelix-Song-8.mp3", created_at: new Date(), updated_at: new Date() },
    { _id: ObjectId(), name: "Lose Yourself to Dance", duration: 353, genre: genres[3]._id, album: albums[2]._id, artists: [artists[2]._id], audio_url: "https://www.soundhelix.com/examples/mp3/SoundHelix-Song-9.mp3", created_at: new Date(), updated_at: new Date() },

    // Lemonade
    { _id: ObjectId(), name: "Formation", duration: 205, genre: genres[4]._id, album: albums[3]._id, artists: [artists[3]._id], audio_url: "https://www.soundhelix.com/examples/mp3/SoundHelix-Song-10.mp3", created_at: new Date(), updated_at: new Date() },
    { _id: ObjectId(), name: "Sorry", duration: 232, genre: genres[4]._id, album: albums[3]._id, artists: [artists[3]._id], audio_url: "https://www.soundhelix.com/examples/mp3/SoundHelix-Song-11.mp3", created_at: new Date(), updated_at: new Date() },
    { _id: ObjectId(), name: "Hold Up", duration: 221, genre: genres[4]._id, album: albums[3]._id, artists: [artists[3]._id], audio_url: "https://www.soundhelix.com/examples/mp3/SoundHelix-Song-12.mp3", created_at: new Date(), updated_at: new Date() },

    // DAMN.
    { _id: ObjectId(), name: "HUMBLE.", duration: 177, genre: genres[0]._id, album: albums[4]._id, artists: [artists[4]._id], audio_url: "https://www.soundhelix.com/examples/mp3/SoundHelix-Song-13.mp3", created_at: new Date(), updated_at: new Date() },
    { _id: ObjectId(), name: "DNA.", duration: 185, genre: genres[0]._id, album: albums[4]._id, artists: [artists[4]._id], audio_url: "https://www.soundhelix.com/examples/mp3/SoundHelix-Song-14.mp3", created_at: new Date(), updated_at: new Date() },
    { _id: ObjectId(), name: "LOYALTY.", duration: 227, genre: genres[0]._id, album: albums[4]._id, artists: [artists[4]._id], audio_url: "https://www.soundhelix.com/examples/mp3/SoundHelix-Song-15.mp3", created_at: new Date(), updated_at: new Date() }
];

db.songs.insertMany(songs);
print("Inserted " + songs.length + " songs");

print("\n=== Database seeded successfully! ===");
print("Genres: " + db.genres.countDocuments());
print("Artists: " + db.artists.countDocuments());
print("Albums: " + db.albums.countDocuments());
print("Songs: " + db.songs.countDocuments());
