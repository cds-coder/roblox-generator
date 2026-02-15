package roblox_profile

import (
	"fmt"
	"math/rand"
	"strings"
)

var wordParts = []string{
	"Shadow", "Flame", "Wolf", "Tiger", "Dragon", "Phoenix", "Hunter", "Star", "Ghost",
	"Legend", "Galaxy", "Frost", "Sonic", "Crystal", "Silver", "Dark", "Power", "Magic", "Light",
	"Alpha", "King", "Queen", "Master", "Pro", "Hero", "Knight", "Beast", "Epic", "Ultra",
	"Fire", "Storm", "Blaze", "Ice", "Sky", "Thunder", "Raven", "Fox", "Lion", "Eagle",
	"Night", "Dawn", "Viper", "Blade", "Hawk", "Claw", "Venom", "Echo", "Bane", "Mystic",
	"Cyber", "Nova", "Orbit", "Pixel", "Glitch", "Byte", "Circuit", "Spark", "Neon", "Chase",
	"Rogue", "Stealth", "Fusion", "Prism", "Wraith", "Saber", "Pulse", "Zero", "Fury",
	"Builder", "Ninja", "Gamer", "Lava", "Stream", "Craft", "Miner", "Block", "Code", "Playz",
	"Panda", "Bear", "Slime", "Duck", "Bacon", "Cookie", "Rocket", "Moon", "Starry", "Turbo",
	"Lucky", "Flash", "Arrow", "Ace", "Omega", "Chaos", "Pixelated", "Void", "Rift", "Toxic",
	"Golden", "Blizzard", "Inferno", "Vortex", "Zoom", "Aqua", "Primal", "Chill", "Hyper", "Stormy",
	"Zap", "Flick", "Max", "Jelly", "Turbo", "Rift", "Blast", "Skater", "Dancer", "Builder",
	"Drift", "Hero", "Haze", "Craze", "Giga", "Sparkly", "Pixel", "Echo", "Rider",

	"Mirage", "Surge", "NovaCore", "Vanta", "Shifter", "Crypt", "Snare", "Ashen", "Fallen", "Warden",
	"Ion", "Delta", "EchoX", "Thorn", "RogueX", "Venator", "Sniper", "Titan", "Celestial", "Oblivion",
	"Mecha", "Warp", "Quantum", "PulseX", "Nightfall", "Bolt", "Scythe", "Striker", "Viral", "Twilight",
	"Nebula", "Dust", "Crimson", "Shade", "Flicker", "Rune", "Techno", "Smash", "Shred", "Core",
	"Switch", "Hollow", "Synth", "Blitz", "ByteStorm", "Glide", "Frostbite", "Rumble", "EchoBlade", "StormChaser",
	"Phantom", "NebulaX", "VortexX", "Specter", "Nebulon", "Aether", "Zenith", "Cinder", "FuryX", "RiftWalker", "QuantumLeap",
	"NovaStrike", "ShadowStrike", "BlazeRunner", "FrostRunner", "ThunderRunner", "MysticRunner", "CosmicRunner", "LunarRunner", "SolarRunner", "StellarRunner", "GalacticRunner",
	"EchoRunner", "PhantomRunner", "NebulaRunner", "VortexRunner", "SpectralRunner", "AetherRunner", "ZenithRunner", "CinderRunner", "FuryRunner", "RiftRunner", "QuantumRunner",
	"NovaHunter", "ShadowHunter", "BlazeHunter", "FrostHunter", "ThunderHunter", "MysticHunter", "CosmicHunter", "LunarHunter", "SolarHunter", "StellarHunter", "GalacticHunter",
	"EchoHunter", "PhantomHunter", "NebulaHunter", "VortexHunter", "SpectralHunter", "AetherHunter", "ZenithHunter", "CinderHunter", "FuryHunter", "RiftHunter", "QuantumHunter",
}

var humanNames = []string{
	"Liam", "Emma", "Noah", "Olivia", "Oliver", "Ava", "Elijah", "Sophia", "Lucas", "Isabella",
	"Mason", "Mia", "Ethan", "Charlotte", "Logan", "Amelia", "Aiden", "Harper", "James", "Evelyn",
	"Jayden", "Abigail", "Henry", "Ella", "Sebastian", "Aria", "Jackson", "Scarlett", "Alexander", "Grace",
	"Mateo", "Chloe", "Michael", "Victoria", "Daniel", "Zoe", "William", "Luna", "Levi", "Hannah",
	"Gabriel", "Addison", "Carter", "Willow", "Wyatt", "Nora", "Isaac", "Layla", "Eli", "Hazel",
	"Samuel", "Ellie", "Jack", "Paisley", "Owen", "Aurora", "Luke", "Brooklyn", "Julian", "Savannah",
	"Jaxon", "Grayson", "Hunter", "Zayden", "Ezra", "Kaylee", "Aubrey", "Riley", "Brooklynn", "Asher",

	"Kai", "Nova", "Milo", "Arlo", "Luca", "Rhett", "Remy", "Sienna", "Zara", "Ari",
	"Kinsley", "Jude", "Kairo", "Nico", "Ember", "River", "Skye", "Oakley", "Rowan", "Everly",
	"Zion", "Alina", "Callum", "Elias", "Finn", "Adeline", "Nina", "Freya", "Leo", "Noelle",
	"Ayden", "Josie", "Talia", "Amara", "Reign", "Indie", "Koda", "Selah", "Maeve", "Ezri",

	"Ahmad", "Omar", "Yusuf", "Hassan", "Ali", "Tariq", "Khalid", "Nabil", "Faris", "Samir",
	"Zayd", "Bilal", "Ibrahim", "Amir", "Rami", "Adnan", "Fahad", "Jamal", "Karim", "Sami",
	"Anas", "Tamer", "Waleed", "Mahmoud", "Zain", "Nasser", "Rashid", "Suleiman", "Ayman", "Mounir",
	"Hani", "Basim", "Saif", "Mustafa", "Nasir", "Imran", "Qasim", "Najib", "Hakeem", "Rafiq",
	"Yahya", "Majid", "Kareem", "Alaa", "Fadl", "Idris", "Munir", "Ridha", "Salim", "Bashir",
	"Abbas", "Harith", "Mazin", "Suhail", "Wahid", "Tawfiq", "Riad", "Nizar", "Zaher", "Thamer",
	"Ashraf", "Issam", "Younes", "Luay", "Jaber", "Raed", "Hamza", "Sherif", "Hilal", "Malik",
	"Habib", "Barakat", "Ghassan", "Haitham", "Najm", "Rafat", "Talal", "Marwan", "Khalil", "Adel",
	"Qays", "Ehab", "Fathi", "Mahdi", "Sabir", "Hatem", "Zaher", "Murad", "Aref", "Rashad",
	"Abed", "Saber", "Labib", "Safi", "Jalil", "Abdul", "Zaher", "Farouq", "Wajdi", "Hatim",

	"Layla", "Fatima", "Aisha", "Zainab", "Mariam", "Noor", "Rania", "Amira", "Huda", "Nadia",
	"Yasmin", "Salma", "Lina", "Dina", "Sara", "Reem", "Alaa", "Nour", "Rim", "Jumana",
	"Hanan", "Samar", "Nisreen", "Lubna", "Sana", "Rabab", "Bushra", "Najwa", "Farah", "Ghada",
	"Lamia", "Samira", "Shaimaa", "Iman", "Ayah", "Rula", "Muna", "Khadija", "Maha", "Tamara",
	"Nada", "Widad", "Ibtisam", "Dalal", "Souad", "Rasha", "Hanadi", "Wafaa", "Mays", "Mariam",
	"Afaf", "Warda", "Razan", "Thuraya", "Amani", "Hiba", "Lubna", "Zahra", "Najah", "Nahla",
	"Ikram", "Manal", "Tahani", "Asma", "Nouran", "Nourah", "Balqis", "Rawan", "Jamila", "Sherine",
	"Maysaa", "Niveen", "Lujain", "Shahed", "Rowan", "Leen", "Suha", "Amal", "Rabab", "Mira",
	"Soumaya", "Hoor", "Bushra", "Marwa", "Batoul", "Wissam", "Nihal", "Arwa", "Thuraya", "Siba",
	"Ruba", "Malak", "Huda", "Inas", "Feryal", "Nourhan", "Jouri", "Rita", "Areej", "Raneem",

	"James", "John", "Robert", "Michael", "William", "David", "Richard", "Joseph", "Thomas", "Charles",
	"Christopher", "Daniel", "Matthew", "Anthony", "Mark", "Donald", "Steven", "Paul", "Andrew", "Joshua",
	"Kenneth", "Kevin", "Brian", "George", "Edward", "Ronald", "Timothy", "Jason", "Jeffrey", "Ryan",
	"Jacob", "Gary", "Nicholas", "Eric", "Stephen", "Jonathan", "Larry", "Justin", "Scott", "Brandon",
	"Frank", "Benjamin", "Gregory", "Samuel", "Raymond", "Patrick", "Alexander", "Jack", "Dennis", "Jerry",
	"Tyler", "Aaron", "Henry", "Douglas", "Jose", "Peter", "Adam", "Zachary", "Nathan", "Walter",
	"Harold", "Kyle", "Carl", "Arthur", "Gerald", "Roger", "Keith", "Lawrence", "Terry", "Sean",
	"Christian", "Albert", "Joe", "Ethan", "Austin", "Jesse", "Willie", "Billy", "Bryan", "Bruce",
	"Jordan", "Ralph", "Roy", "Noah", "Dylan", "Eugene", "Wayne", "Alan", "Juan", "Louis",
	"Russell", "Gabriel", "Randy", "Philip", "Harry", "Vincent", "Bobby", "Johnny", "Logan", "Cody",

	"Mary", "Patricia", "Jennifer", "Linda", "Elizabeth", "Barbara", "Susan", "Jessica", "Sarah", "Karen",
	"Nancy", "Lisa", "Margaret", "Betty", "Sandra", "Ashley", "Dorothy", "Kimberly", "Emily", "Donna",
	"Michelle", "Carol", "Amanda", "Melissa", "Deborah", "Stephanie", "Rebecca", "Laura", "Sharon", "Cynthia",
	"Kathleen", "Amy", "Shirley", "Angela", "Helen", "Anna", "Brenda", "Pamela", "Nicole", "Emma",
	"Samantha", "Katherine", "Christine", "Debra", "Rachel", "Catherine", "Carolyn", "Janet", "Ruth", "Maria",
	"Heather", "Diane", "Virginia", "Julie", "Joyce", "Victoria", "Olivia", "Kelly", "Christina", "Lauren",
	"Joan", "Evelyn", "Judith", "Megan", "Cheryl", "Andrea", "Hannah", "Martha", "Jacqueline", "Frances",
	"Gloria", "Ann", "Teresa", "Kathryn", "Sara", "Janice", "Jean", "Alice", "Madison", "Doris",
	"Abigail", "Julia", "Judy", "Grace", "Denise", "Amber", "Marilyn", "Beverly", "Danielle", "Theresa",
	"Sophia", "Marie", "Diana", "Brittany", "Natalie", "Isabella", "Charlotte", "Rose", "Alexis", "Kayla",
}

var replacements = map[string]string{
	"golden": "Auric",
	"toxic":  "Viral",
	"haz":    "Hz",
}

func GetUsername() string {

	var part1, part2 string

	if rand.Float64() < 0.5 {
		part1 = humanNames[rand.Intn(len(humanNames))]
		part2 = wordParts[rand.Intn(len(wordParts))]
	} else {
		part1 = wordParts[rand.Intn(len(wordParts))]
		part2 = wordParts[rand.Intn(len(wordParts))]
		for strings.EqualFold(part1, part2) {
			part2 = wordParts[rand.Intn(len(wordParts))]
		}
	}

	var username string
	if rand.Float64() < 0.5 {
		username = part1 + "_" + part2
	} else {
		username = part1 + part2
	}

	username += fmt.Sprintf("%d", 1+rand.Intn(9999))

	if rand.Float64() < 0.3 {
		username = strings.ReplaceAll(username, "o", "0")
	}

	if rand.Float64() < 0.3 {
		username = strings.ReplaceAll(username, "e", "3")
	}

	lower := strings.ToLower(username)
	for bad, good := range replacements {
		if strings.Contains(lower, bad) {
			username = strings.ReplaceAll(strings.ToLower(username), bad, good)
		}
	}

	for year := 1900; year <= 2024; year++ {
		yearStr := fmt.Sprintf("%d", year)
		if strings.Contains(username, yearStr) {
			username = strings.ReplaceAll(username, yearStr, fmt.Sprintf("%d", 100+rand.Intn(900)))
		}
	}

	username = strings.Trim(username, "_")
	username = strings.ReplaceAll(username, "__", "_")

	if len(username) < 3 {
		username += fmt.Sprintf("%d", 10+rand.Intn(90))
	}

	if len(username) > 20 {
		username = username[:20]
	}

	return username
}

func GetPassword() string {
	lower := "abcdefghijklmnopqrstuvwxyz"
	upper := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits := "0123456789"
	special := "@!"
	all := lower + upper + digits + special

	length := 8 + rand.Intn(8)
	pass := make([]byte, length)

	for i := 0; i < length; i++ {
		pass[i] = all[rand.Intn(len(all))]
	}

	return string(pass)
}

func GetBirthDay() string {
	day := fmt.Sprintf("%02d", 1+rand.Intn(28))
	month := fmt.Sprintf("%02d", 1+rand.Intn(12))
	year := 1990 + rand.Intn(11)
	return fmt.Sprintf("%d-%s-%sT23:00:00.000Z", year, month, day)
}

func GetGender() int {
	return 1 + rand.Intn(2)
}
