package domain

type Wagon struct {
	ID     int
	Number int    // شماره واگن
	Type   string // نوع واگن
	Axles  int    // تعداد محور

	// --- Weights ---
	WeightEmpty       float64 // وزن واگن خالی
	WeightLoaded      float64 // وزن واگن با بار
	MaxCapacity       float64 // ظرفیت بارگیری
	LoadVolume        float64 // حجم بارگیری
	BrakeWeightEmpty  float64 // وزن ترمز بی بار
	BrakeWeightLoaded float64 // وزن ترمز با بار

	// --- Dimensions ---
	Length             float64 // طول واگن
	LoadLength         float64 // طول بارگیری
	LoadWidth          float64 // عرض بارگیری
	FloorHeight        float64 // ارتفاع از ریل تا کف
	InternalHeight     float64 // ارتفاع از کف به بالا
	BogiePivotDistance float64 // فاصله مرکز دو بوژی
	WheelDiameter      float64 // قطر چرخ

	// --- Technical Details ---
	RIVCode           string  // حرف RIV
	Manufacturer      string  // کشور سازنده
	Year              string  // سال ورود
	BogieType         string  // نوع بوژی
	BearingType       string  // نوع جعبه یاتاقان
	SpringType        string  // نوع فنر
	HandBrakeType     string  // نوع ترمز دستی
	HandBrakeWeight   float64 // وزن ترمز دستی
	ControlValveType  string  // نوع سوپاپ سه قلو (KE1CSL...)
	BrakeCylinderType string  // نوع خودکار ترمز (Cylinder)
	CouplingType      string  // نوع قلاب
}

// SelectedWagon represents a wagon added to the train composition with specific user inputs (Dynamic Data).
type SelectedWagon struct {
	WagonSpec Wagon // Embeds the static data

	// User Inputs
	IsMainBrakeHealthy   bool   // Status of the main air brake
	IsHandBrakeHealthy   bool   // Status of the hand brake
	IsBrakeHandleHealthy bool   // Status of the brake handle/valve
	IsLoaded             bool   // True if loaded, False if empty
	HasDangerousGoods    bool   // True if carrying dangerous goods
	DangerousGoodsCode   string // The code of dangerous goods (e.g., "2a", "3b")

	// Computed Values for Calculation
	EffectiveWeight      float64 // Final weight based on Load Status
	EffectiveBrakeWeight float64 // Final brake weight based on Brake Health & Load
}
