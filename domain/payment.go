package domain

import "strconv"

const url = `https://yoomoney.ru/quickpay/confirm?receiver=4100118475856202&quickpay-form=button&paymentType=AC&sum=2&label=`
const redirectUrl = `successURL=https://hooli-smotrim.ru/profile`

func Payment(user_id int) string {
	payment := url + strconv.Itoa(user_id) + "&" + redirectUrl
	return payment
}
