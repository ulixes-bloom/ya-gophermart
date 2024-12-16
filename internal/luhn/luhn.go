package luhn

// Валидация числа по алгоритму Луна
func ValidateNumber(number string) bool {
	total := 0
	isSecondDigit := false

	for i := len(number) - 1; i >= 0; i-- {
		digit := int(number[i] - '0')

		if isSecondDigit {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}

		total += digit
		isSecondDigit = !isSecondDigit
	}

	return total%10 == 0
}
