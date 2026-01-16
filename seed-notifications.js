/**
 * Seed script za popunjavanje notifikacija u Cassandri
 *
 * Instalacija: npm install cassandra-driver
 * Pokretanje: node seed-notifications.js <user_id>
 *
 * Primer: node seed-notifications.js 507f1f77bcf86cd799439011
 */

const cassandra = require('cassandra-driver');

const client = new cassandra.Client({
  contactPoints: ['localhost:9042'],
  localDataCenter: 'datacenter1',
  keyspace: 'notifications'
});

async function seedNotifications(userId) {
  if (!userId) {
    console.error('Greška: Morate proslediti user_id kao argument!');
    console.log('Upotreba: node seed-notifications.js <user_id>');
    console.log('');
    console.log('Da biste dobili user_id, pokrenite:');
    console.log('  docker exec -it mongodb-users mongosh -u admin -p admin123 --authenticationDatabase admin --eval "use users; db.users.find({}, {_id: 1, username: 1}).pretty()"');
    process.exit(1);
  }

  const notifications = [
    {
      id: cassandra.types.Uuid.random(),
      user_id: userId,
      message: 'Dobrodošli na Spotify Clone! Istražite našu kolekciju muzike.',
      type: 'system',
      read: false,
      created_at: new Date()
    },
    {
      id: cassandra.types.Uuid.random(),
      user_id: userId,
      message: 'Eminem je objavio novi album "The Death of Slim Shady"!',
      type: 'new_album',
      read: false,
      created_at: new Date(Date.now() - 1000 * 60 * 60) // pre 1 sat
    },
    {
      id: cassandra.types.Uuid.random(),
      user_id: userId,
      message: 'Nova pesma "Houdini" je dostupna za slušanje.',
      type: 'new_song',
      read: false,
      created_at: new Date(Date.now() - 1000 * 60 * 60 * 2) // pre 2 sata
    },
    {
      id: cassandra.types.Uuid.random(),
      user_id: userId,
      message: 'Vaša ocena za pesmu "Lose Yourself" je sačuvana.',
      type: 'rating',
      read: true,
      created_at: new Date(Date.now() - 1000 * 60 * 60 * 24) // pre 1 dan
    },
    {
      id: cassandra.types.Uuid.random(),
      user_id: userId,
      message: 'Preporučujemo vam album "Recovery" na osnovu vaših preferencija.',
      type: 'system',
      read: true,
      created_at: new Date(Date.now() - 1000 * 60 * 60 * 24 * 2) // pre 2 dana
    },
    {
      id: cassandra.types.Uuid.random(),
      user_id: userId,
      message: 'Novi umetnik "Dr. Dre" je dodat u našu biblioteku!',
      type: 'new_album',
      read: false,
      created_at: new Date(Date.now() - 1000 * 60 * 30) // pre 30 minuta
    }
  ];

  console.log(`Dodajem ${notifications.length} notifikacija za korisnika: ${userId}`);

  const query = `
    INSERT INTO notifications (id, user_id, message, type, read, created_at)
    VALUES (?, ?, ?, ?, ?, ?)
  `;

  for (const n of notifications) {
    try {
      await client.execute(query, [n.id, n.user_id, n.message, n.type, n.read, n.created_at], { prepare: true });
      console.log(`  ✓ Dodato: "${n.message.substring(0, 40)}..."`);
    } catch (err) {
      console.error(`  ✗ Greška: ${err.message}`);
    }
  }

  console.log('\nGotovo!');
}

async function main() {
  const userId = process.argv[2];

  try {
    await client.connect();
    console.log('Povezan na Cassandru\n');

    await seedNotifications(userId);
  } catch (err) {
    console.error('Greška pri povezivanju:', err.message);
  } finally {
    await client.shutdown();
  }
}

main();
