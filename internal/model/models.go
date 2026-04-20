package model

type Role struct {
	ID    int    `json:"id"`
	Code  string `json:"code"`
	Label string `json:"label"`
}

type OrderStatus struct {
	ID    int    `json:"id"`
	Label string `json:"label"`
}

type PaymentType struct {
	ID    int    `json:"id"`
	Label string `json:"label"`
}

type User struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	Pass      string `json:"-"` // md5-хеш, никогда не отдаём клиенту
	Fio       string `json:"fio"`
	IDRole    int    `json:"id_role"`
	IsBlocked int    `json:"is_blocked"`

	// для удобства join запросов
	Role Role `json:"Role"`
}

type Client struct {
	ID    int    `json:"id"`
	Label string `json:"label"`
}

type Order struct {
	ID       int     `json:"id"`
	Label    string  `json:"label"`
	IDStatus int     `json:"id_status"`
	IDClient int     `json:"id_client"`
	Amount   float64 `json:"amount"`
}

type OrderItem struct {
	ID      int     `json:"id"`
	Label   string  `json:"label"`
	IDOrder int     `json:"id_order"`
	Amount  float64 `json:"amount"`
}

type Payment struct {
	ID            int     `json:"id"`
	IDOrder       int     `json:"id_order"`
	IDPaymentType int     `json:"id_payment_type"`
	Amount        float64 `json:"amount"`
}
