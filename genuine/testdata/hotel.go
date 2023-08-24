package testdata

import "time"

type Hotels struct {
	ID            int
	Name          string
	City          string
	Rating        float64
	PricePerNight float64
	Description   string
	Facilities    string
	ContactEmail  string
	Phone         string
	Website       string
	CreatedAt     time.Time
}

var HotelExample = []Hotels{
	{1, "Luxury Resort", "Miami", 4.500, 250.000, "A luxurious beachside resort", "Pool, Spa, Private Beach", "info@luxuryresort.com", "+1-123-456-7890", "https://www.luxuryresort.com", time.Date(2025, 3, 14, 9, 23, 45, 0, time.UTC)},
	{2, "Cozy Inn", "New York", 3.800, 120.000, "A charming inn in the heart of the city", "Free Wi-Fi, Breakfast, Lounge", "info@cozyinn.com", "+1-987-654-3210", "https://www.cozyinn.com", time.Date(1984, 7, 9, 9, 12, 32, 0, time.UTC)},
	{3, "Seaside Lodge", "Los Angeles", 4.200, 180.000, "A cozy lodge with ocean views", "Ocean View, Fireplace, Restaurant", "info@seasidelodge.com", "+1-555-123-4567", "https://www.seasidelodge.com", time.Date(2012, 2, 27, 9, 56, 8, 0, time.UTC)},
	{4, "Mountain Retreat", "Denver", 4.000, 150.000, "A retreat in the mountains", "Hiking Trails, Spa, Scenic Views", "info@mountainretreat.com", "+1-888-567-8901", "https://www.mountainretreat.com", time.Date(2001, 11, 2, 9, 34, 17, 0, time.UTC)},
	{5, "Urban Hotel", "Chicago", 3.700, 200.000, "A modern hotel in the city center", "Fitness Center, Rooftop Bar", "info@urbanhotel.com", "+1-333-456-7890", "https://www.urbanhotel.com", time.Date(1993, 10, 15, 10, 42, 19, 0, time.UTC)},
	{6, "Sunny Getaway", "Miami", 4.600, 280.000, "A sunny paradise by the beach", "Beachfront, Poolside Bar", "info@sunnygetaway.com", "+1-555-789-1234", "https://www.sunnygetaway.com", time.Date(2017, 5, 2, 10, 18, 41, 0, time.UTC)},
	{7, "Downtown Suites", "New York", 4.300, 180.000, "Luxury suites in the heart of the city", "Spa, Concierge, Sky Lounge", "info@downtownsuites.com", "+1-999-888-7777", "https://www.downtownsuites.com", time.Date(2029, 12, 31, 10, 3, 52, 0, time.UTC)},
	{8, "Beachfront Resort", "Los Angeles", 4.800, 320.000, "A beachfront oasis with stunning views", "Private Beach, Oceanfront Dining", "info@beachfrontresort.com", "+1-444-555-6666", "https://www.beachfrontresort.com", time.Date(2015, 9, 4, 10, 24, 37, 0, time.UTC)},
	{9, "Mountain Chalet", "Denver", 4.400, 210.000, "Charming chalets nestled in the mountains", "Ski Access, Fireplace, Spa", "info@mountainchalet.com", "+1-222-333-4444", "https://www.mountainchalet.com", time.Date(1990, 1, 23, 11, 55, 13, 0, time.UTC)},
	{10, "City Center Hotel", "Chicago", 3.900, 190.000, "Conveniently located hotel in the city center", "Business Center, On-site Restaurant", "info@citycenterhotel.com", "+1-777-888-9999", "https://www.citycenterhotel.com", time.Date(1991, 6, 30, 11, 11, 22, 0, time.UTC)},
	{11, "Lakeview Lodge", "Seattle", 4.700, 260.000, "Lakeside retreat with stunning lake views", "Boating, Fishing, Lakeside Dining", "info@lakeviewlodge.com", "+1-111-222-3333", "https://www.lakeviewlodge.com", time.Date(1980, 4, 8, 11, 9, 45, 0, time.UTC)},
	{12, "Riverside Inn", "San Francisco", 4.100, 170.000, "Charming inn along the riverside", "Scenic Views, Garden Patio", "info@riversideinn.com", "+1-444-333-2222", "https://www.riversideinn.com", time.Date(2022, 1, 1, 11, 59, 57, 0, time.UTC)},
	{13, "Historic Mansion", "Boston", 4.500, 300.000, "Elegant mansion with a rich history", "Antique Furnishings, Ballroom", "info@historicmansion.com", "+1-555-666-7777", "https://www.historicmansion.com", time.Date(2000, 8, 13, 12, 15, 0, 0, time.UTC)},
	{14, "Desert Oasis", "Phoenix", 3.600, 140.000, "A tranquil oasis in the desert", "Spa, Desert Gardens, Pool", "info@desertoasis.com", "+1-777-888-9999", "https://www.desertoasis.com", time.Date(2021, 11, 11, 12, 38, 2, 0, time.UTC)},
	{15, "Skyline Tower", "Las Vegas", 4.400, 240.000, "Modern tower with breathtaking city views", "Rooftop Pool, Casino, Nightclub", "info@skylinetower.com", "+1-111-222-3333", "https://www.skylinetower.com", time.Date(1996, 12, 31, 12, 17, 36, 0, time.UTC)},
}
